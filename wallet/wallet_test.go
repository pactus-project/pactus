package wallet

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/util"
)

func TestSaveToFile(t *testing.T) {
	w, err := GenerateWallet("test", "test")
	assert.NoError(t, err)
	assert.NoError(t, w.SaveToFile(util.TempFilePath()))
}
