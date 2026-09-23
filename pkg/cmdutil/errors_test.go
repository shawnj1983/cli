package cmdutil

import (
	"testing"

	"github.com/cli/cli/v2/pkg/iostreams"
	"github.com/stretchr/testify/assert"
)

func TestNonInteractiveHint(t *testing.T) {
	assert.Equal(t, "must specify a name", NonInteractiveHint(nil, "must specify a name"))

	ios, _, _, _ := iostreams.Test()
	assert.Equal(t, "must specify a name", NonInteractiveHint(ios, "must specify a name"))

	ios.SetNeverPrompt(true)
	ios.SetNeverPromptReason("cursor-cloud is driving the CLI")
	assert.Equal(t,
		"must specify a name (cursor-cloud is driving the CLI; prompts disabled)",
		NonInteractiveHint(ios, "must specify a name"))
}
