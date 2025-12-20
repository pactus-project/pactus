//go:build gtk

package gtkutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSmallGray_EscapesMarkup(t *testing.T) {
	got := SmallGray(`a<b&"c"`)
	assert.NotEmpty(t, got)
	assert.NotContains(t, got, "<b", "expected markup to be escaped")
	assert.NotContains(t, got, "&\"", "expected markup to be escaped")
	assert.NotContains(t, got, `"`, "expected markup to be escaped")
	assert.Contains(t, got, "&lt;")
	assert.Contains(t, got, "&amp;")
}
