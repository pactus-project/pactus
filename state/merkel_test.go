package state

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zarbchain/zarb-go/account"
	"github.com/zarbchain/zarb-go/validator"
)

func TestChangeAcc(t *testing.T) {
	st1 := setupStatewithOneValidator(t)
	st2 := setupStatewithOneValidator(t)

	require.Equal(t, st1.store.TotalAccounts(), 1)

	acc1, _ := account.GenerateTestAccount(1)
	acc2, _ := account.GenerateTestAccount(2)
	acc3, _ := account.GenerateTestAccount(3)
	acc4, _ := account.GenerateTestAccount(4)

	st1.store.UpdateAccount(acc1)
	st1.store.UpdateAccount(acc2)
	st1.store.UpdateAccount(acc3)
	st1.store.UpdateAccount(acc4)
	root1 := st1.accountsMerkleRootHash()

	// Change an account state
	acc3.IncSequence()

	st2.store.UpdateAccount(acc2)
	st2.store.UpdateAccount(acc3)
	st2.store.UpdateAccount(acc1)
	st2.store.UpdateAccount(acc4)
	root2 := st2.accountsMerkleRootHash()

	assert.NotEqual(t, root1, root2)
}

func TestChangeVal(t *testing.T) {
	st1 := setupStatewithOneValidator(t)
	st2 := setupStatewithOneValidator(t)

	require.Equal(t, st1.store.TotalValidators(), 1)

	val1, _ := validator.GenerateTestValidator(1)
	val2, _ := validator.GenerateTestValidator(2)
	val3, _ := validator.GenerateTestValidator(3)
	val4, _ := validator.GenerateTestValidator(4)

	st1.store.UpdateValidator(val1)
	st1.store.UpdateValidator(val2)
	st1.store.UpdateValidator(val3)
	st1.store.UpdateValidator(val4)
	root1 := st1.validatorsMerkleRootHash()

	// Change an account state
	val3.IncSequence()

	st2.store.UpdateValidator(val2)
	st2.store.UpdateValidator(val3)
	st2.store.UpdateValidator(val1)
	st2.store.UpdateValidator(val4)
	root2 := st2.validatorsMerkleRootHash()

	assert.NotEqual(t, root1, root2)
}
