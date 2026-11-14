//go:build gtk

package gtkutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEscapeMarkup(t *testing.T) {
	assert.Equal(t, "&lt;a href=&quot;x&quot;&gt;link&lt;/a&gt; &amp; more",
		EscapeMarkup(`<a href="x">link</a> & more`))
	assert.Equal(t, "plain text", EscapeMarkup("plain text"))
}
