package state

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zarbchain/zarb-go/types/account"
	"github.com/zarbchain/zarb-go/types/validator"
)

func TestChangeAcc(t *testing.T) {
	setup(t)

	require.Equal(t, tState1.store.TotalAccounts(), int32(1))

	acc1, _ := account.GenerateTestAccount(1)
	acc2, _ := account.GenerateTestAccount(2)
	acc3, _ := account.GenerateTestAccount(3)
	acc4, _ := account.GenerateTestAccount(4)

	tState1.store.UpdateAccount(acc1)
	tState1.store.UpdateAccount(acc2)
	tState1.store.UpdateAccount(acc3)
	tState1.store.UpdateAccount(acc4)
	root1 := tState1.accountsMerkleRoot()

	// Change an account
	acc3.IncSequence()

	tState2.store.UpdateAccount(acc2)
	tState2.store.UpdateAccount(acc3)
	tState2.store.UpdateAccount(acc1)
	tState2.store.UpdateAccount(acc4)
	root2 := tState2.accountsMerkleRoot()

	assert.NotEqual(t, root1, root2)
}

func TestChangeVal(t *testing.T) {
	setup(t)

	require.Equal(t, tState1.store.TotalValidators(), int32(4))

	val1, _ := validator.GenerateTestValidator(4)
	val2, _ := validator.GenerateTestValidator(5)
	val3, _ := validator.GenerateTestValidator(6)
	val4, _ := validator.GenerateTestValidator(7)

	tState1.store.UpdateValidator(val1)
	tState1.store.UpdateValidator(val2)
	tState1.store.UpdateValidator(val3)
	tState1.store.UpdateValidator(val4)
	root1 := tState1.validatorsMerkleRoot()

	// Change a validtor
	val3.IncSequence()

	tState2.store.UpdateValidator(val2)
	tState2.store.UpdateValidator(val3)
	tState2.store.UpdateValidator(val1)
	tState2.store.UpdateValidator(val4)
	root2 := tState2.validatorsMerkleRoot()

	assert.NotEqual(t, root1, root2)
}

func TestCalculatingGenesisState(t *testing.T) {
	setup(t)

	r := tState1.calculateGenesisStateRootFromGenesisDoc()
	assert.Equal(t, tState1.stateRoot(), r)
}
