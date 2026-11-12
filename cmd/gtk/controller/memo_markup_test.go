//go:build gtk

package controller

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// The memo of a transaction is shown in Pango markup: it comes from the user or from the network,
// and it must be escaped wherever it is shown, so it can't add markup like a link.
func TestMemoIsEscapedInMarkup(t *testing.T) {
	files, err := filepath.Glob("*.go")
	require.NoError(t, err)

	memo := regexp.MustCompile(`[A-Za-z_.()]*\.Memo\(\)`)
	checked := 0
	for _, file := range files {
		if strings.HasSuffix(file, "_test.go") {
			continue
		}
		data, err := os.ReadFile(file)
		require.NoError(t, err)

		for _, line := range strings.Split(string(data), "\n") {
			if !memo.MatchString(line) || strings.Contains(line, "!=") || strings.Contains(line, "==") {
				continue
			}
			checked++
			assert.Contains(t, line, "EscapeMarkup(", "%s: %s", file, strings.TrimSpace(line))
		}
	}
	assert.Positive(t, checked)
}
