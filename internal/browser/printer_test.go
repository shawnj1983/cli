package browser

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPrinterBrowse(t *testing.T) {
	var buf bytes.Buffer
	b := NewPrinter(&buf)

	require.NoError(t, b.Browse("https://github.com/cli/cli"))
	assert.Equal(t, "https://github.com/cli/cli\n", buf.String())
}

func TestForInvoker(t *testing.T) {
	tests := []struct {
		name         string
		agent        string
		ghBrowserSet bool
		wantPrint    bool
	}{
		{name: "no agent opens a browser", agent: "", wantPrint: false},
		{name: "cursor IDE opens a browser", agent: "cursor", wantPrint: false},
		{name: "cursor-cloud prints the URL", agent: "cursor-cloud", wantPrint: true},
		{name: "cursor-cli prints the URL", agent: "cursor-cli", wantPrint: true},
		{name: "claude-code prints the URL", agent: "claude-code", wantPrint: true},
		{name: "GH_BROWSER keeps a real launcher", agent: "cursor-cloud", ghBrowserSet: true, wantPrint: false},
		{name: "replit is not driving", agent: "replit", wantPrint: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stdout bytes.Buffer
			var stderr bytes.Buffer
			b := ForInvoker(tt.agent, tt.ghBrowserSet, &stdout, &stderr)

			_, isPrinter := b.(*Printer)
			assert.Equal(t, tt.wantPrint, isPrinter)
		})
	}
}
