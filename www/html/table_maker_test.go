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
