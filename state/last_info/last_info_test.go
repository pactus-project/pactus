package last_info

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/sortition"
	"github.com/zarbchain/zarb-go/store"
	"github.com/zarbchain/zarb-go/tx"
	"github.com/zarbchain/zarb-go/util"
	"github.com/zarbchain/zarb-go/validator"
)

func TestRestore(t *testing.T) {
	path := util.TempDirPath()
	li1 := NewLastInfo(path)
	li2 := NewLastInfo(path)

	val1 := validator.NewValidator(crypto.GenerateTestSigner().PublicKey(), 10, 20)
	val2 := validator.NewValidator(crypto.GenerateTestSigner().PublicKey(), 18, 28)
	val3 := validator.NewValidator(crypto.GenerateTestSigner().PublicKey(), 2, 12)
	val4 := validator.NewValidator(crypto.GenerateTestSigner().PublicKey(), 6, 16)

	trx1, _ := tx.GenerateTestBondTx()
	trx2, _ := tx.GenerateTestSendTx()
	trx3, newValSigner := tx.GenerateTestSortitionTx()
	hash := crypto.GenerateTestHash()
	lastCertificate := block.GenerateTestCertificate(hash)
	txIDs := block.NewTxIDs()
	txIDs.Append(trx1.ID())
	txIDs.Append(trx2.ID())
	txIDs.Append(trx3.ID())
	lastSortitionSeed := sortition.GenerateRandomSeed()
	lastBlock := block.MakeBlock(1, util.Now(), txIDs,
		hash,
		crypto.GenerateTestHash(),
		lastCertificate, lastSortitionSeed, val1.Address())
	lastBlockHeight := 111
	lastBlockHash := lastBlock.Hash()
	ctrxs := []*tx.CommittedTx{}
	for _, trx := range []*tx.Tx{trx1, trx2, trx3} {
		ctrx := &tx.CommittedTx{
			Tx:      trx,
			Receipt: trx.GenerateReceipt(tx.Ok, lastBlockHash),
		}
		ctrxs = append(ctrxs, ctrx)
	}

	li1.SetSortitionSeed(lastSortitionSeed)
	li1.SetBlockHeight(lastBlockHeight)
	li1.SetBlockHash(lastBlockHash)
	li1.SetCertificate(lastCertificate)
	li1.SetBlockTime(lastBlock.Header().Time())
	assert.NoError(t, li1.SaveLastInfo())

	store := store.MockingStore()
	_, err := li2.RestoreLastInfo(store)
	assert.Error(t, err)

	store.SaveBlock(lastBlockHeight, lastBlock)
	_, err = li2.RestoreLastInfo(store)
	assert.Error(t, err)

	for _, ctrx := range ctrxs {
		store.SaveTransaction(ctrx)
	}
	_, err = li2.RestoreLastInfo(store)
	assert.Error(t, err)

	val := validator.NewValidator(newValSigner.PublicKey(), 54, 45)
	store.UpdateValidator(val)
	_, err = li2.RestoreLastInfo(store)
	assert.Error(t, err)

	store.UpdateValidator(val1)
	store.UpdateValidator(val2)
	store.UpdateValidator(val3)
	store.UpdateValidator(val4)

	_, err = li2.RestoreLastInfo(store)
	assert.NoError(t, err)

	assert.Equal(t, li1.SortitionSeed(), li2.SortitionSeed())
	assert.Equal(t, li1.BlockHeight(), li2.BlockHeight())
	assert.Equal(t, li1.BlockHash(), li2.BlockHash())
	assert.Equal(t, li1.Certificate().Hash(), li2.Certificate().Hash())
	assert.Equal(t, li1.BlockTime(), li2.BlockTime())
}
