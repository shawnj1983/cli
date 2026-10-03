package browser

import (
	"fmt"
	"io"

	"github.com/cli/cli/v2/internal/agents"
)

// Printer writes destination URLs instead of launching a GUI browser.
type Printer struct {
	w io.Writer
}

// NewPrinter returns a Browser that prints each URL to w.
func NewPrinter(w io.Writer) Browser {
	return &Printer{w: w}
}

func (p *Printer) Browse(url string) error {
	_, err := fmt.Fprintln(p.w, url)
	return err
}

// ForInvoker returns a Browser appropriate for the process invoking gh.
// Driving agents get a printer so commands such as `gh browse` and
// `gh pr view --web` do not hang trying to launch a GUI browser.
// An explicit GH_BROWSER keeps the real launcher.
func ForInvoker(agent string, ghBrowserSet bool, stdout, stderr io.Writer) Browser {
	if agents.IsDriving(agents.AgentName(agent)) && !ghBrowserSet {
		return NewPrinter(stdout)
	}
	return New("", stdout, stderr)
}
