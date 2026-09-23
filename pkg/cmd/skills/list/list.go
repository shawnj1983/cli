package list

import (
	"fmt"
	"sort"
	"strings"

	"github.com/MakeNowJust/heredoc"
	"github.com/cli/cli/v2/git"
	"github.com/cli/cli/v2/internal/skills/installed"
	"github.com/cli/cli/v2/internal/skills/installer"
	"github.com/cli/cli/v2/internal/skills/registry"
	"github.com/cli/cli/v2/internal/tableprinter"
	"github.com/cli/cli/v2/pkg/cmdutil"
	"github.com/cli/cli/v2/pkg/iostreams"
	"github.com/spf13/cobra"
)

// SkillListFields defines the set of fields available for --json output.
var SkillListFields = []string{
	"name",
	"agent",
	"scope",
	"source",
	"pinned",
	"path",
}

type ListOptions struct {
	IO        *iostreams.IOStreams
	GitClient *git.Client
	Exporter  cmdutil.Exporter

	Agent string
	Scope string
	Dir   string
}

type skillRow struct {
	Name   string `json:"name"`
	Agent  string `json:"agent"`
	Scope  string `json:"scope"`
	Source string `json:"source"`
	Pinned string `json:"pinned"`
	Path   string `json:"path"`
}

func (s skillRow) ExportData(fields []string) map[string]interface{} {
	return cmdutil.StructExportData(s, fields)
}

// NewCmdList creates the "skills list" command.
func NewCmdList(f *cmdutil.Factory, runF func(*ListOptions) error) *cobra.Command {
	opts := &ListOptions{
		IO:        f.IOStreams,
		GitClient: f.GitClient,
	}

	cmd := &cobra.Command{
		Use:   "list [flags]",
		Short: "List installed skills (preview)",
		Long: heredoc.Docf(`
			List skills already installed for known agent hosts.

			Scans project and user skill directories for Copilot, Claude, Cursor,
			and other registered hosts. Use %[1]s--agent%[1]s or %[1]s--scope%[1]s to
			narrow the results, or %[1]s--dir%[1]s to scan a custom directory.

			This command is useful for coding agents that need to see which
			skills are already present before installing or updating more.
		`, "`"),
		Example: heredoc.Doc(`
			# List every installed skill
			$ gh skill list

			# List skills for one host
			$ gh skill list --agent cursor

			# List only user-scope skills
			$ gh skill list --scope user

			# Machine-readable output
			$ gh skill list --json name,agent,source
		`),
		Aliases: []string{"ls"},
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if runF != nil {
				return runF(opts)
			}
			return listRun(opts)
		},
	}

	cmdutil.StringEnumFlag(cmd, &opts.Agent, "agent", "", "", registry.AgentIDs(), "Filter by target agent")
	cmdutil.StringEnumFlag(cmd, &opts.Scope, "scope", "", "", []string{"project", "user"}, "Filter by installation scope")
	cmd.Flags().StringVar(&opts.Dir, "dir", "", "Scan a custom directory for installed skills")
	cmdutil.AddJSONFlags(cmd, &opts.Exporter, SkillListFields)
	cmdutil.DisableAuthCheck(cmd)

	return cmd
}

func listRun(opts *ListOptions) error {
	gitRoot := installer.ResolveGitRoot(opts.GitClient)
	homeDir := installer.ResolveHomeDir()

	var skills []installed.Skill
	if opts.Dir != "" {
		scanned, err := installed.ScanDir(opts.Dir, nil, "")
		if err != nil {
			return fmt.Errorf("could not scan directory: %w", err)
		}
		skills = scanned
	} else {
		skills = installed.ScanAll(gitRoot, homeDir)
	}

	if opts.Agent != "" {
		var filtered []installed.Skill
		for _, s := range skills {
			if s.AgentID() == opts.Agent {
				filtered = append(filtered, s)
			}
		}
		skills = filtered
	}

	if opts.Scope != "" {
		var filtered []installed.Skill
		for _, s := range skills {
			if string(s.Scope) == opts.Scope {
				filtered = append(filtered, s)
			}
		}
		skills = filtered
	}

	sort.Slice(skills, func(i, j int) bool {
		if skills[i].Name != skills[j].Name {
			return strings.ToLower(skills[i].Name) < strings.ToLower(skills[j].Name)
		}
		if skills[i].AgentID() != skills[j].AgentID() {
			return skills[i].AgentID() < skills[j].AgentID()
		}
		return skills[i].Scope < skills[j].Scope
	})

	rows := make([]skillRow, 0, len(skills))
	for _, s := range skills {
		rows = append(rows, skillRow{
			Name:   s.Name,
			Agent:  s.AgentID(),
			Scope:  string(s.Scope),
			Source: s.SourceRepo(),
			Pinned: s.Pinned,
			Path:   s.Dir,
		})
	}

	if opts.Exporter != nil {
		return opts.Exporter.Write(opts.IO, rows)
	}

	if len(rows) == 0 {
		return cmdutil.NewNoResultsError("no installed skills found")
	}

	return renderTable(opts.IO, rows)
}

func renderTable(io *iostreams.IOStreams, rows []skillRow) error {
	table := tableprinter.New(io, tableprinter.WithHeader("NAME", "AGENT", "SCOPE", "SOURCE", "PIN"))
	for _, r := range rows {
		table.AddField(r.Name)
		table.AddField(emptyDash(r.Agent))
		table.AddField(emptyDash(r.Scope))
		table.AddField(emptyDash(r.Source))
		table.AddField(emptyDash(r.Pinned))
		table.EndRow()
	}
	return table.Render()
}

func emptyDash(s string) string {
	if s == "" {
		return "-"
	}
	return s
}
