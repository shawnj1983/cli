package agents

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func lookup(vars map[string]string) func(string) (string, bool) {
	return func(key string) (string, bool) {
		v, ok := vars[key]
		return v, ok
	}
}

func TestParseAgentName(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    AgentName
		wantErr bool
	}{
		{name: "valid lowercase", input: "my-agent", want: "my-agent"},
		{name: "valid with underscore", input: "my_agent_v2", want: "my_agent_v2"},
		{name: "valid uppercase", input: "MyAgent", want: "MyAgent"},
		{name: "valid numbers", input: "agent123", want: "agent123"},
		{name: "spaces rejected", input: "my agent", wantErr: true},
		{name: "newline rejected", input: "my\nagent", wantErr: true},
		{name: "carriage return rejected", input: "my\ragent", wantErr: true},
		{name: "null byte rejected", input: "my\x00agent", wantErr: true},
		{name: "dot rejected", input: "my.agent", wantErr: true},
		{name: "slash rejected", input: "my/agent", wantErr: true},
		{name: "empty rejected", input: "", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseAgentName(tt.input)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestDetectWith(t *testing.T) {
	tests := []struct {
		name      string
		env       map[string]string
		wantAgent AgentName
	}{
		{
			name:      "clean environment",
			env:       map[string]string{},
			wantAgent: "",
		},
		{
			name:      "empty var is not detected",
			env:       map[string]string{"GEMINI_CLI": ""},
			wantAgent: "",
		},
		{
			name:      "AGENT=amp detected as amp",
			env:       map[string]string{"AGENT": "amp"},
			wantAgent: "amp",
		},
		{
			name:      "AGENT with non-amp value is ignored",
			env:       map[string]string{"AGENT": "other"},
			wantAgent: "",
		},
		{
			name:      "AI_AGENT returns value as agent name",
			env:       map[string]string{"AI_AGENT": "some-agent"},
			wantAgent: "some-agent",
		},
		{
			name:      "AI_AGENT with invalid characters is ignored",
			env:       map[string]string{"AI_AGENT": "bad\nagent"},
			wantAgent: "",
		},
		{
			name:      "AI_AGENT with spaces is ignored",
			env:       map[string]string{"AI_AGENT": "bad agent"},
			wantAgent: "",
		},
		{
			name:      "AI_AGENT takes priority over AGENT",
			env:       map[string]string{"AGENT": "amp", "AI_AGENT": "other"},
			wantAgent: "other",
		},
		{
			name:      "CODEX_SANDBOX",
			env:       map[string]string{"CODEX_SANDBOX": "seatbelt"},
			wantAgent: "codex",
		},
		{
			name:      "CODEX_CI",
			env:       map[string]string{"CODEX_CI": "1"},
			wantAgent: "codex",
		},
		{
			name:      "CODEX_THREAD_ID",
			env:       map[string]string{"CODEX_THREAD_ID": "abc"},
			wantAgent: "codex",
		},
		{
			name:      "GEMINI_CLI",
			env:       map[string]string{"GEMINI_CLI": "1"},
			wantAgent: "gemini-cli",
		},
		{
			name:      "COPILOT_CLI",
			env:       map[string]string{"COPILOT_CLI": "1"},
			wantAgent: "copilot-cli",
		},
		{
			name:      "OPENCODE",
			env:       map[string]string{"OPENCODE": "1"},
			wantAgent: "opencode",
		},
		{
			name:      "CLAUDECODE",
			env:       map[string]string{"CLAUDECODE": "1"},
			wantAgent: "claude-code",
		},
		{
			name:      "AGENT=amp takes priority over CLAUDECODE",
			env:       map[string]string{"AGENT": "amp", "CLAUDECODE": "1"},
			wantAgent: "amp",
		},
		{
			name:      "invalid AI_AGENT falls through to tool-specific detection",
			env:       map[string]string{"AI_AGENT": "bad agent", "GEMINI_CLI": "1"},
			wantAgent: "gemini-cli",
		},
		{
			name:      "ANTIGRAVITY_AGENT",
			env:       map[string]string{"ANTIGRAVITY_AGENT": "1"},
			wantAgent: "antigravity",
		},
		{
			name:      "AUGMENT_AGENT",
			env:       map[string]string{"AUGMENT_AGENT": "1"},
			wantAgent: "augment-cli",
		},
		{
			name:      "REPL_ID",
			env:       map[string]string{"REPL_ID": "abc123"},
			wantAgent: "replit",
		},
		{
			name:      "GOOSE_PROVIDER",
			env:       map[string]string{"GOOSE_PROVIDER": "anthropic"},
			wantAgent: "goose",
		},
		{
			name:      "claude-code takes priority over goose",
			env:       map[string]string{"GOOSE_PROVIDER": "anthropic", "CLAUDECODE": "1"},
			wantAgent: "claude-code",
		},
		{
			name:      "kiro takes priority over goose",
			env:       map[string]string{"GOOSE_PROVIDER": "anthropic", "TERM_PROGRAM": "kiro"},
			wantAgent: "kiro",
		},
		{
			name:      "CLAUDE_CODE_IS_COWORK detected as cowork",
			env:       map[string]string{"CLAUDE_CODE_IS_COWORK": "1"},
			wantAgent: "cowork",
		},
		{
			name:      "cowork takes priority over CLAUDECODE",
			env:       map[string]string{"CLAUDE_CODE_IS_COWORK": "1", "CLAUDECODE": "1"},
			wantAgent: "cowork",
		},
		{
			name:      "CLAUDE_CODE",
			env:       map[string]string{"CLAUDE_CODE": "1"},
			wantAgent: "claude-code",
		},
		{
			name:      "CURSOR_AGENT_SOCKET detected as cursor-cloud",
			env:       map[string]string{"CURSOR_AGENT_SOCKET": "/run/cursor/api.sock"},
			wantAgent: "cursor-cloud",
		},
		{
			name:      "CURSOR_CONVERSATION_ID with bc- prefix detected as cursor-cloud",
			env:       map[string]string{"CURSOR_CONVERSATION_ID": "bc-76853ea2-d731-495e-8b09-5b171936f133"},
			wantAgent: "cursor-cloud",
		},
		{
			name:      "CURSOR_CONVERSATION_ID without bc- prefix is ignored",
			env:       map[string]string{"CURSOR_CONVERSATION_ID": "local-thread-1"},
			wantAgent: "",
		},
		{
			name:      "cursor-cloud takes priority over CURSOR_AGENT",
			env:       map[string]string{"CURSOR_AGENT_SOCKET": "/run/cursor/api.sock", "CURSOR_AGENT": "1"},
			wantAgent: "cursor-cloud",
		},
		{
			name:      "CURSOR_AGENT detected as cursor-cli",
			env:       map[string]string{"CURSOR_AGENT": "1"},
			wantAgent: "cursor-cli",
		},
		{
			name:      "CURSOR_EXTENSION_HOST_ROLE agent-exec detected as cursor-cli",
			env:       map[string]string{"CURSOR_EXTENSION_HOST_ROLE": "agent-exec"},
			wantAgent: "cursor-cli",
		},
		{
			name:      "CURSOR_EXTENSION_HOST_ROLE with other value is ignored",
			env:       map[string]string{"CURSOR_EXTENSION_HOST_ROLE": "worker"},
			wantAgent: "",
		},
		{
			name:      "CURSOR_TRACE_ID detected as cursor",
			env:       map[string]string{"CURSOR_TRACE_ID": "abc"},
			wantAgent: "cursor",
		},
		{
			name:      "CURSOR_AGENT takes priority over CURSOR_TRACE_ID",
			env:       map[string]string{"CURSOR_TRACE_ID": "abc", "CURSOR_AGENT": "1"},
			wantAgent: "cursor-cli",
		},
		{
			name:      "TERM_PROGRAM kiro detected as kiro",
			env:       map[string]string{"TERM_PROGRAM": "kiro"},
			wantAgent: "kiro",
		},
		{
			name:      "TERM_PROGRAM with kiro as a substring is ignored",
			env:       map[string]string{"TERM_PROGRAM": "kirostudio"},
			wantAgent: "",
		},
		{
			name:      "PATH containing .pi/agent detected as pi",
			env:       map[string]string{"PATH": "/usr/bin:/home/user/.pi/agent/bin"},
			wantAgent: "pi",
		},
		{
			name:      "PATH with .pi/agent not on a path boundary is ignored",
			env:       map[string]string{"PATH": "/usr/bin:/home/user/x.pi/agent"},
			wantAgent: "",
		},
		{
			name:      "PATH with Windows .pi\\agent separators detected as pi",
			env:       map[string]string{"PATH": `C:\Windows;C:\Users\user\.pi\agent\bin`},
			wantAgent: "pi",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := detectWith(lookup(tt.env))
			assert.Equal(t, tt.wantAgent, got)
		})
	}
}

func TestSkillHostID(t *testing.T) {
	tests := []struct {
		name     string
		agent    AgentName
		wantHost string
	}{
		{name: "amp", agent: "amp", wantHost: "amp"},
		{name: "claude-code", agent: "claude-code", wantHost: "claude-code"},
		{name: "cowork maps to claude-code", agent: "cowork", wantHost: "claude-code"},
		{name: "codex", agent: "codex", wantHost: "codex"},
		{name: "copilot-cli maps to github-copilot", agent: "copilot-cli", wantHost: "github-copilot"},
		{name: "gemini-cli", agent: "gemini-cli", wantHost: "gemini-cli"},
		{name: "opencode", agent: "opencode", wantHost: "opencode"},
		{name: "antigravity", agent: "antigravity", wantHost: "antigravity"},
		{name: "augment-cli maps to augment", agent: "augment-cli", wantHost: "augment"},
		{name: "replit", agent: "replit", wantHost: "replit"},
		{name: "goose", agent: "goose", wantHost: "goose"},
		{name: "cursor", agent: "cursor", wantHost: "cursor"},
		{name: "cursor-cli maps to cursor", agent: "cursor-cli", wantHost: "cursor"},
		{name: "cursor-cloud maps to cursor", agent: "cursor-cloud", wantHost: "cursor"},
		{name: "kiro maps to kiro-cli", agent: "kiro", wantHost: "kiro-cli"},
		{name: "pi", agent: "pi", wantHost: "pi"},
		{name: "unknown agent", agent: "some-custom-agent", wantHost: ""},
		{name: "empty agent", agent: "", wantHost: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.wantHost, SkillHostID(tt.agent))
		})
	}
}
