package list

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/MakeNowJust/heredoc"
	"github.com/cli/cli/v2/git"
	"github.com/cli/cli/v2/pkg/cmdutil"
	"github.com/cli/cli/v2/pkg/iostreams"
	"github.com/cli/cli/v2/pkg/jsonfieldstest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJSONFields(t *testing.T) {
	jsonfieldstest.ExpectCommandToSupportJSONFields(t, NewCmdList, []string{
		"name",
		"agent",
		"scope",
		"source",
		"pinned",
		"path",
	})
}

func TestNewCmdList(t *testing.T) {
	tests := []struct {
		name     string
		args     string
		wantOpts ListOptions
		wantErr  string
	}{
		{
			name:     "no flags",
			args:     "",
			wantOpts: ListOptions{},
		},
		{
			name: "agent filter",
			args: "--agent cursor",
			wantOpts: ListOptions{
				Agent: "cursor",
			},
		},
		{
			name: "scope filter",
			args: "--scope user",
			wantOpts: ListOptions{
				Scope: "user",
			},
		},
		{
			name: "custom dir",
			args: "--dir /tmp/skills",
			wantOpts: ListOptions{
				Dir: "/tmp/skills",
			},
		},
		{
			name:    "unknown agent",
			args:    "--agent not-a-host",
			wantErr: "valid values are",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &cmdutil.Factory{}
			var gotOpts *ListOptions
			cmd := NewCmdList(f, func(opts *ListOptions) error {
				gotOpts = opts
				return nil
			})

			var argv []string
			if tt.args != "" {
				argv = strings.Fields(tt.args)
			}
			cmd.SetArgs(argv)
			cmd.SetOut(io.Discard)
			cmd.SetErr(io.Discard)

			_, err := cmd.ExecuteC()
			if tt.wantErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.wantOpts.Agent, gotOpts.Agent)
			assert.Equal(t, tt.wantOpts.Scope, gotOpts.Scope)
			assert.Equal(t, tt.wantOpts.Dir, gotOpts.Dir)
		})
	}
}

func TestListRun(t *testing.T) {
	writeSkill := func(t *testing.T, dir, name, body string) {
		t.Helper()
		skillDir := filepath.Join(dir, name)
		require.NoError(t, os.MkdirAll(skillDir, 0o755))
		require.NoError(t, os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte(body), 0o644))
	}

	t.Run("lists skills from --dir", func(t *testing.T) {
		dir := t.TempDir()
		writeSkill(t, dir, "git-commit", heredoc.Doc(`
			---
			name: git-commit
			metadata:
			  github-repo: https://github.com/monalisa/awesome-copilot
			  github-tree-sha: abc123
			  github-pinned: v1.0.0
			---
			body
		`))
		writeSkill(t, dir, "unknown", heredoc.Doc(`
			---
			name: unknown
			---
			body
		`))

		ios, _, stdout, _ := iostreams.Test()
		err := listRun(&ListOptions{
			IO:  ios,
			Dir: dir,
		})
		require.NoError(t, err)
		out := stdout.String()
		assert.Contains(t, out, "git-commit")
		assert.Contains(t, out, "monalisa/awesome-copilot")
		assert.Contains(t, out, "v1.0.0")
		assert.Contains(t, out, "unknown")
	})

	t.Run("json export", func(t *testing.T) {
		dir := t.TempDir()
		writeSkill(t, dir, "git-commit", heredoc.Doc(`
			---
			name: git-commit
			metadata:
			  github-repo: https://github.com/monalisa/awesome-copilot
			  github-tree-sha: abc123
			---
			body
		`))

		ios, _, stdout, _ := iostreams.Test()
		exporter := cmdutil.NewJSONExporter()
		exporter.SetFields([]string{"name", "source", "pinned"})
		err := listRun(&ListOptions{
			IO:       ios,
			Dir:      dir,
			Exporter: exporter,
		})
		require.NoError(t, err)
		assert.JSONEq(t, `[{"name":"git-commit","source":"monalisa/awesome-copilot","pinned":""}]`, strings.TrimSpace(stdout.String()))
	})

	t.Run("empty dir is no results", func(t *testing.T) {
		ios, _, _, _ := iostreams.Test()
		err := listRun(&ListOptions{
			IO:  ios,
			Dir: t.TempDir(),
		})
		require.EqualError(t, err, "no installed skills found")
	})

	t.Run("scans agent directories and filters by agent", func(t *testing.T) {
		repoDir := t.TempDir()
		homeDir := t.TempDir()
		t.Setenv("HOME", homeDir)
		t.Setenv("USERPROFILE", homeDir)

		writeSkill(t, filepath.Join(repoDir, ".claude", "skills"), "code-review", heredoc.Doc(`
			---
			name: code-review
			metadata:
			  github-repo: https://github.com/monalisa/octocat-skills
			  github-tree-sha: def456
			---
			body
		`))
		writeSkill(t, filepath.Join(repoDir, ".agents", "skills"), "git-commit", heredoc.Doc(`
			---
			name: git-commit
			metadata:
			  github-repo: https://github.com/monalisa/octocat-skills
			  github-tree-sha: abc123
			---
			body
		`))

		ios, _, stdout, _ := iostreams.Test()
		err := listRun(&ListOptions{
			IO:        ios,
			GitClient: &git.Client{RepoDir: repoDir},
			Agent:     "claude-code",
		})
		require.NoError(t, err)
		out := stdout.String()
		assert.Contains(t, out, "code-review")
		assert.NotContains(t, out, "git-commit")
	})
}

func TestNewCmdList_jsonFields(t *testing.T) {
	f := &cmdutil.Factory{}
	cmd := NewCmdList(f, func(*ListOptions) error { return nil })
	var stdout bytes.Buffer
	cmd.SetOut(&stdout)
	cmd.SetErr(&stdout)
	cmd.SetArgs([]string{"--json"})
	_, err := cmd.ExecuteC()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "Specify one or more comma-separated fields")
}
