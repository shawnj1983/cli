package installed

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/cli/cli/v2/internal/skills/frontmatter"
	"github.com/cli/cli/v2/internal/skills/registry"
	"github.com/cli/cli/v2/internal/skills/source"
)

// Skill is a locally installed skill parsed from its SKILL.md frontmatter.
type Skill struct {
	Name        string
	RepoHost    string
	Owner       string
	Repo        string
	TreeSHA     string
	Pinned      string
	SourcePath  string
	Dir         string
	Host        *registry.AgentHost
	Scope       registry.Scope
	MetadataErr error
}

// SourceRepo returns owner/repo, or host/owner/repo when not on github.com.
func (s Skill) SourceRepo() string {
	if s.Owner == "" || s.Repo == "" {
		return ""
	}
	if s.RepoHost != "" && s.RepoHost != source.SupportedHost {
		return s.RepoHost + "/" + s.Owner + "/" + s.Repo
	}
	return s.Owner + "/" + s.Repo
}

// AgentID returns the skill host ID, or empty when the install directory
// was not associated with a known host.
func (s Skill) AgentID() string {
	if s.Host == nil {
		return ""
	}
	return s.Host.ID
}

// ScanAll walks every registered agent skill directory (project and user
// scope) and collects installed skills. Shared install roots are scanned
// only once.
func ScanAll(gitRoot, homeDir string) []Skill {
	scannedDirs := make(map[string]bool)
	var all []Skill

	for i := range registry.Agents {
		host := &registry.Agents[i]
		for _, scope := range []registry.Scope{registry.ScopeProject, registry.ScopeUser} {
			dir, err := host.InstallDir(scope, gitRoot, homeDir)
			if err != nil {
				continue
			}
			if scannedDirs[dir] {
				continue
			}
			scannedDirs[dir] = true
			skills, err := ScanDir(dir, host, scope)
			if err != nil {
				continue
			}
			all = append(all, skills...)
		}
	}

	return all
}

// ScanDir reads SKILL.md files in a skills directory and extracts GitHub
// metadata from their frontmatter. It handles both flat layouts
// ({dir}/{name}/SKILL.md) and namespaced layouts
// ({dir}/{namespace}/{name}/SKILL.md).
func ScanDir(skillsDir string, host *registry.AgentHost, scope registry.Scope) ([]Skill, error) {
	entries, err := os.ReadDir(skillsDir)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("could not read skills directory: %w", err)
	}

	var skills []Skill
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}

		skillFile := filepath.Join(skillsDir, e.Name(), "SKILL.md")
		if data, readErr := os.ReadFile(skillFile); readErr == nil {
			if s, ok := parseSkill(data, e.Name(), filepath.Join(skillsDir, e.Name()), host, scope); ok {
				skills = append(skills, s)
				continue
			}
		}

		subEntries, subErr := os.ReadDir(filepath.Join(skillsDir, e.Name()))
		if subErr != nil {
			continue
		}
		for _, sub := range subEntries {
			if !sub.IsDir() {
				continue
			}
			subSkillFile := filepath.Join(skillsDir, e.Name(), sub.Name(), "SKILL.md")
			if data, readErr := os.ReadFile(subSkillFile); readErr == nil {
				installName := e.Name() + "/" + sub.Name()
				if s, ok := parseSkill(data, installName, filepath.Join(skillsDir, e.Name(), sub.Name()), host, scope); ok {
					skills = append(skills, s)
				}
			}
		}
	}

	return skills, nil
}

func parseSkill(data []byte, name, dir string, host *registry.AgentHost, scope registry.Scope) (Skill, bool) {
	result, err := frontmatter.Parse(string(data))
	if err != nil {
		return Skill{
			Name:        name,
			Dir:         dir,
			Host:        host,
			Scope:       scope,
			MetadataErr: fmt.Errorf("invalid SKILL.md: %w", err),
		}, true
	}

	s := Skill{
		Name:  name,
		Dir:   dir,
		Host:  host,
		Scope: scope,
	}

	if result.Metadata.Meta != nil {
		repoInfo, ok, repoErr := source.ParseMetadataRepo(result.Metadata.Meta)
		if repoErr != nil {
			s.MetadataErr = repoErr
		} else if ok {
			if err := source.ValidateSupportedHost(repoInfo.RepoHost()); err != nil {
				s.MetadataErr = err
			} else {
				s.RepoHost = repoInfo.RepoHost()
				s.Owner = repoInfo.RepoOwner()
				s.Repo = repoInfo.RepoName()
			}
		}
		s.TreeSHA, _ = result.Metadata.Meta["github-tree-sha"].(string)
		s.Pinned, _ = result.Metadata.Meta["github-pinned"].(string)
		s.SourcePath, _ = result.Metadata.Meta["github-path"].(string)
	}

	return s, true
}
