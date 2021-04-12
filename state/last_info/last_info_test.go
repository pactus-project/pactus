package last_info

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/sortition"
	"github.com/zarbchain/zarb-go/util"
)

func TestLastInfoAccessors(t *testing.T) {
	li := NewLastInfo(util.TempDirPath())

	lastSortitionSeed := sortition.GenerateRandomSeed()
	lastBlockHeight := 111
	lastBlockHash := crypto.GenerateTestHash()
	lastReceiptsHash := crypto.GenerateTestHash()
	lastCertificate := block.GenerateTestCertificate(lastBlockHash)
	lastBlockTime := time.Now()

	li.SetSortitionSeed(lastSortitionSeed)
	li.SetBlockHeight(lastBlockHeight)
	li.SetBlockHash(lastBlockHash)
	li.SetReceiptsHash(lastReceiptsHash)
	li.SetCertificate(lastCertificate)
	li.SetBlockTime(lastBlockTime)

	assert.Equal(t, li.SortitionSeed(), lastSortitionSeed)
	assert.Equal(t, li.BlockHeight(), lastBlockHeight)
	assert.Equal(t, li.BlockHash(), lastBlockHash)
	assert.Equal(t, li.ReceiptsHash(), lastReceiptsHash)
	assert.Equal(t, li.Certificate(), lastCertificate)
	assert.Equal(t, li.BlockTime(), lastBlockTime)
}
