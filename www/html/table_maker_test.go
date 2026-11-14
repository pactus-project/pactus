package html

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTableMakerEscapesValues(t *testing.T) {
	script := `<script src=//evil.example/x.js></script>`

	tmk := newTableMaker()
	tmk.addRowString("Moniker", script)
	tmk.addRowStrings("Protocols", []string{script, `"><img src=x>`})
	tmk.addRowBlockHash("Hash", `"><script>alert(1)</script>`)
	tmk.addRowAccAddress("Account", script)
	tmk.addRowValAddress("Validator", script)
	tmk.addRowTxID("Tx", script)
	out := tmk.html()

	assert.NotContains(t, out, "<script")
	assert.NotContains(t, out, "<img")
	assert.Contains(t, out, "&lt;script")
}

func TestTableMakerMoniker(t *testing.T) {
	t.Run("Plain moniker", func(t *testing.T) {
		tmk := newTableMaker()
		tmk.addRowMoniker("Alice")
		out := tmk.html()

		assert.Contains(t, out, "<tr><td>Moniker</td><td>Alice</td></tr>")
		assert.NotContains(t, out, "Moniker Details")
	})

	t.Run("Structured moniker", func(t *testing.T) {
		tmk := newTableMaker()
		tmk.addRowMoniker("name=alice;web=https://example.com")
		out := tmk.html()

		assert.Contains(t, out, "<tr><td>Moniker</td><td>alice</td></tr>")
		assert.Contains(t, out, "<tr><td>Moniker Details</td><td>name=alice;web=https://example.com</td></tr>")
		assert.NotContains(t, out, "<a ", "the website must not be a link")
	})

	t.Run("Structured moniker without a name", func(t *testing.T) {
		tmk := newTableMaker()
		tmk.addRowMoniker("web=https://example.com")
		out := tmk.html()

		assert.Contains(t, out, "<tr><td>Moniker</td><td>web=https://example.com</td></tr>")
		assert.NotContains(t, out, "Moniker Details")
	})

	t.Run("Unsafe characters", func(t *testing.T) {
		tmk := newTableMaker()
		tmk.addRowMoniker("<b>Ali\x1bce\n</b>")
		out := tmk.html()

		assert.Contains(t, out, "<tr><td>Moniker</td><td>&lt;b&gt;Alice&lt;/b&gt;</td></tr>")
	})
}
