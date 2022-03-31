package wallet

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/util"
)

func TestSaveToFile(t *testing.T) {
	w, err := NewWallet(util.TempFilePath(), "1")
	assert.NoError(t, err)
	assert.NoError(t, w.SaveToFile())
}
