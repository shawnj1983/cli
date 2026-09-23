package installed

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/MakeNowJust/heredoc"
	"github.com/cli/cli/v2/internal/skills/registry"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestScanDir(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(t *testing.T, dir string)
		verify func(t *testing.T, skills []Skill, err error)
	}{
		{
			name: "happy path with metadata, no metadata, and pinned skills",
			setup: func(t *testing.T, dir string) {
				t.Helper()

				skillDir := filepath.Join(dir, "git-commit")
				require.NoError(t, os.MkdirAll(skillDir, 0o755))
				content := heredoc.Doc(`
					---
					name: git-commit
					description: Git commit helper
					metadata:
					  github-repo: https://github.com/monalisa/awesome-copilot
					  github-tree-sha: abc123
					  github-path: skills/git-commit
					---
					Body content
				`)
				require.NoError(t, os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte(content), 0o644))

				noMetaDir := filepath.Join(dir, "unknown-skill")
				require.NoError(t, os.MkdirAll(noMetaDir, 0o755))
				require.NoError(t, os.WriteFile(filepath.Join(noMetaDir, "SKILL.md"), []byte(heredoc.Doc(`
					---
					name: unknown-skill
					---
					No metadata here
				`)), 0o644))

				pinnedDir := filepath.Join(dir, "pinned-skill")
				require.NoError(t, os.MkdirAll(pinnedDir, 0o755))
				require.NoError(t, os.WriteFile(filepath.Join(pinnedDir, "SKILL.md"), []byte(heredoc.Doc(`
					---
					name: pinned-skill
					metadata:
					  github-repo: https://github.com/octocat/hubot-skills
					  github-tree-sha: def456
					  github-pinned: v1.0.0
					---
					Pinned content
				`)), 0o644))
			},
			verify: func(t *testing.T, skills []Skill, err error) {
				t.Helper()
				require.NoError(t, err)
				assert.Len(t, skills, 3)

				byName := make(map[string]Skill)
				for _, s := range skills {
					byName[s.Name] = s
				}

				gc := byName["git-commit"]
				assert.Equal(t, "monalisa", gc.Owner)
				assert.Equal(t, "awesome-copilot", gc.Repo)
				assert.Equal(t, "github.com", gc.RepoHost)
				assert.Equal(t, "abc123", gc.TreeSHA)
				assert.Equal(t, "skills/git-commit", gc.SourcePath)
				assert.Empty(t, gc.Pinned)
				assert.Equal(t, "monalisa/awesome-copilot", gc.SourceRepo())

				us := byName["unknown-skill"]
				assert.Empty(t, us.Owner)
				assert.Empty(t, us.Repo)
				assert.Empty(t, us.SourceRepo())

				ps := byName["pinned-skill"]
				assert.Equal(t, "github.com", ps.RepoHost)
				assert.Equal(t, "v1.0.0", ps.Pinned)
			},
		},
		{
			name: "unsupported host metadata returns error",
			setup: func(t *testing.T, dir string) {
				t.Helper()
				skillDir := filepath.Join(dir, "enterprise-skill")
				require.NoError(t, os.MkdirAll(skillDir, 0o755))
				require.NoError(t, os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte(heredoc.Doc(`
					---
					name: enterprise-skill
					metadata:
					  github-repo: https://acme.ghes.com/monalisa/octocat-skills
					  github-tree-sha: abc123
					---
					body
				`)), 0o644))
			},
			verify: func(t *testing.T, skills []Skill, err error) {
				t.Helper()
				require.NoError(t, err)
				require.Len(t, skills, 1)
				require.Error(t, skills[0].MetadataErr)
				assert.Contains(t, skills[0].MetadataErr.Error(), "does not currently support GitHub Enterprise Server")
			},
		},
		{
			name: "non-existent directory returns nil",
			verify: func(t *testing.T, skills []Skill, err error) {
				t.Helper()
				require.NoError(t, err)
				assert.Nil(t, skills)
			},
		},
		{
			name: "corrupted YAML is skipped gracefully",
			setup: func(t *testing.T, dir string) {
				t.Helper()
				skillDir := filepath.Join(dir, "corrupt")
				require.NoError(t, os.MkdirAll(skillDir, 0o755))
				require.NoError(t, os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte(heredoc.Doc(`
					---
					not: valid: yaml: [broken
					---
					body
				`)), 0o644))
			},
			verify: func(t *testing.T, skills []Skill, err error) {
				t.Helper()
				require.NoError(t, err)
				require.Len(t, skills, 1)
				assert.Equal(t, "corrupt", skills[0].Name)
				assert.ErrorContains(t, skills[0].MetadataErr, "invalid SKILL.md")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := filepath.Join(t.TempDir(), "skills")
			if tt.setup != nil {
				require.NoError(t, os.MkdirAll(dir, 0o755))
				tt.setup(t, dir)
			}

			skills, err := ScanDir(dir, nil, "")
			tt.verify(t, skills, err)
		})
	}
}

func TestScanAllDeduplicatesSharedProjectDirs(t *testing.T) {
	repoDir := t.TempDir()
	homeDir := t.TempDir()

	sharedSkillDir := filepath.Join(repoDir, ".agents", "skills", "git-commit")
	require.NoError(t, os.MkdirAll(sharedSkillDir, 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(sharedSkillDir, "SKILL.md"), []byte(heredoc.Doc(`
		---
		name: git-commit
		metadata:
		  github-repo: https://github.com/monalisa/octocat-skills
		  github-tree-sha: abc123
		---
		Body
	`)), 0o644))

	claudeSkillDir := filepath.Join(repoDir, ".claude", "skills", "code-review")
	require.NoError(t, os.MkdirAll(claudeSkillDir, 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(claudeSkillDir, "SKILL.md"), []byte(heredoc.Doc(`
		---
		name: code-review
		metadata:
		  github-repo: https://github.com/monalisa/octocat-skills
		  github-tree-sha: def456
		---
		Body
	`)), 0o644))

	skills := ScanAll(repoDir, homeDir)
	require.Len(t, skills, 2)

	byName := make(map[string]Skill)
	for _, skill := range skills {
		byName[skill.Name] = skill
	}

	assert.Equal(t, registry.ScopeProject, byName["git-commit"].Scope)
	assert.Equal(t, registry.ScopeProject, byName["code-review"].Scope)
	// Shared .agents/skills is first claimed by github-copilot in the registry order.
	assert.Equal(t, "github-copilot", byName["git-commit"].AgentID())
	assert.Equal(t, "claude-code", byName["code-review"].AgentID())
}
