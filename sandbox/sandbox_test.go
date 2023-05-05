package sandbox

import (
	"testing"

	"github.com/pactus-project/pactus/committee"
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/store"
	"github.com/pactus-project/pactus/types/account"
	"github.com/pactus-project/pactus/types/block"
	"github.com/pactus-project/pactus/types/param"
	"github.com/pactus-project/pactus/types/validator"
	"github.com/pactus-project/pactus/util"
	"github.com/stretchr/testify/assert"
)

var tSigners []crypto.Signer
var tStore *store.MockStore
var tSandbox *sandbox

func setup(t *testing.T) {
	tStore = store.MockingStore()
	params := param.DefaultParams()
	params.TransactionToLiveInterval = 64

	committee, signers := committee.GenerateTestCommittee(21)
	acc := account.NewAccount(0)
	acc.AddToBalance(21 * 1e14)
	tStore.UpdateAccount(crypto.TreasuryAddress, acc)

	for _, val := range committee.Validators() {
		// For testing purpose we create some test accounts first.
		// Account number is the validator number plus one,
		// since account #0 is the Treasury account.
		acc := account.NewAccount(val.Number() + 1)
		tStore.UpdateValidator(val)
		tStore.UpdateAccount(val.Address(), acc)
	}

	tSigners = signers
	tSandbox = NewSandbox(tStore, params, committee).(*sandbox)

	assert.Equal(t, tSandbox.CurrentHeight(), uint32(1))
	lastHeight := util.RandUint32(144) + 21
	for i := uint32(1); i < lastHeight; i++ {
		b := block.GenerateTestBlock(nil, nil)
		c := block.GenerateTestCertificate(b.Hash())
		tStore.SaveBlock(i, b, c)
	}
	assert.Equal(t, tSandbox.CurrentHeight(), lastHeight)
	assert.Equal(t, tSandbox.Params(), params)
}

func TestAccountChange(t *testing.T) {
	setup(t)

	t.Run("Should returns nil for invalid address", func(t *testing.T) {
		invAddr := crypto.GenerateTestAddress()
		assert.Nil(t, tSandbox.Account(invAddr))
	})

	t.Run("Retrieve an account from store, modify it and commit it", func(t *testing.T) {
		acc1, signer1 := account.GenerateTestAccount(888)
		tStore.UpdateAccount(signer1.Address(), acc1)

		acc1a := tSandbox.Account(signer1.Address())
		assert.Equal(t, acc1, acc1a)

		acc1a.IncSequence()
		acc1a.AddToBalance(acc1a.Balance() + 1)

		assert.False(t, tSandbox.accounts[signer1.Address()].Updated)
		tSandbox.UpdateAccount(signer1.Address(), acc1a)
		assert.True(t, tSandbox.accounts[signer1.Address()].Updated)
	})

	t.Run("Make new account", func(t *testing.T) {
		addr2 := crypto.GenerateTestAddress()
		acc2 := tSandbox.MakeNewAccount(addr2)

		acc2.IncSequence()
		acc2.AddToBalance(acc2.Balance() + 1)

		tSandbox.UpdateAccount(addr2, acc2)
		acc22 := tSandbox.Account(addr2)
		assert.Equal(t, acc2, acc22)
	})
}

func TestValidatorChange(t *testing.T) {
	setup(t)

	t.Run("Should returns nil for invalid address", func(t *testing.T) {
		invAddr := crypto.GenerateTestAddress()
		assert.Nil(t, tSandbox.Validator(invAddr))
	})

	t.Run("Retrieve an validator from store, modify it and commit it", func(t *testing.T) {
		val1, _ := validator.GenerateTestValidator(888)
		tStore.UpdateValidator(val1)

		val1a := tSandbox.Validator(val1.Address())
		assert.Equal(t, val1.Hash(), val1a.Hash())

		val1a.IncSequence()
		val1a.AddToStake(val1a.Stake() + 1)

		assert.False(t, tSandbox.validators[val1a.Address()].Updated)
		tSandbox.UpdateValidator(val1a)
		assert.True(t, tSandbox.validators[val1a.Address()].Updated)
	})

	t.Run("Make new validator", func(t *testing.T) {
		pub, _ := bls.GenerateTestKeyPair()
		val2 := tSandbox.MakeNewValidator(pub)

		val2.IncSequence()
		val2.AddToStake(val2.Stake() + 1)

		tSandbox.UpdateValidator(val2)
		val22 := tSandbox.Validator(val2.Address())
		assert.Equal(t, val2, val22)
	})
}

func TestTotalAccountCounter(t *testing.T) {
	setup(t)

	t.Run("Should update total account counter", func(t *testing.T) {
		assert.Equal(t, tStore.TotalAccounts(), int32(len(tSigners)+1))

		addr1 := crypto.GenerateTestAddress()
		addr2 := crypto.GenerateTestAddress()
		acc := tSandbox.MakeNewAccount(addr1)
		assert.Equal(t, acc.Number(), int32(tSandbox.Committee().Size()+1))
		acc2 := tSandbox.MakeNewAccount(addr2)
		assert.Equal(t, acc2.Number(), int32(tSandbox.Committee().Size()+2))
		assert.Equal(t, acc2.Balance(), int64(0))
	})
}

func TestTotalValidatorCounter(t *testing.T) {
	setup(t)

	t.Run("Should update total validator counter", func(t *testing.T) {
		assert.Equal(t, tStore.TotalValidators(), int32(tSandbox.Committee().Size()))

		pub, _ := bls.GenerateTestKeyPair()
		pub2, _ := bls.GenerateTestKeyPair()
		val1 := tSandbox.MakeNewValidator(pub)
		val1.UpdateLastBondingHeight(tSandbox.CurrentHeight())
		assert.Equal(t, val1.Number(), int32(tSandbox.Committee().Size()))
		assert.Equal(t, val1.LastBondingHeight(), tSandbox.CurrentHeight())

		val2 := tSandbox.MakeNewValidator(pub2)
		val2.UpdateLastBondingHeight(tSandbox.CurrentHeight() + 1)
		assert.Equal(t, val2.Number(), int32(tSandbox.Committee().Size()+1))
		assert.Equal(t, val2.LastBondingHeight(), tSandbox.CurrentHeight()+1)
		assert.Equal(t, val2.Stake(), int64(0))
	})
}

func TestCreateDuplicated(t *testing.T) {
	setup(t)

	t.Run("Try creating duplicated account, Should panic", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("The code did not panic")
			}
		}()
		addr := crypto.TreasuryAddress
		tSandbox.MakeNewAccount(addr)
	})

	t.Run("Try creating duplicated validator, Should panic", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("The code did not panic")
			}
		}()
		pub := tSigners[3].PublicKey()
		tSandbox.MakeNewValidator(pub.(*bls.PublicKey))
	})
}

func TestUpdateFromOutsideTheSandbox(t *testing.T) {
	setup(t)

	t.Run("Try update an account from outside the sandbox, Should panic", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("The code did not panic")
			}
		}()
		acc, signer := account.GenerateTestAccount(999)
		tSandbox.UpdateAccount(signer.Address(), acc)
	})

	t.Run("Try update a validator from outside the sandbox, Should panic", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("The code did not panic")
			}
		}()
		val, _ := validator.GenerateTestValidator(999)
		tSandbox.UpdateValidator(val)
	})
}

func TestDeepCopy(t *testing.T) {
	setup(t)

	addr := crypto.GenerateTestAddress()
	pub, _ := bls.GenerateTestKeyPair()
	acc1 := tSandbox.MakeNewAccount(addr)
	val1 := tSandbox.MakeNewValidator(pub)

	acc2 := tSandbox.Account(addr)
	val2 := tSandbox.Validator(pub.Address())

	assert.Equal(t, acc1, acc2)
	assert.Equal(t, val1.Hash(), val2.Hash())

	acc1.IncSequence()
	val1.IncSequence()

	acc2.AddToBalance(1)
	val2.AddToStake(1)

	assert.NotEqual(t, acc1.Hash(), acc2.Hash())
	assert.NotEqual(t, val1.Hash(), val2.Hash())

	acc3 := tSandbox.accounts[addr]
	val3 := tSandbox.validators[pub.Address()]

	assert.NotEqual(t, acc1.Hash(), acc3.Account.Hash())
	assert.NotEqual(t, val1.Hash(), val3.Validator.Hash())

	assert.NotEqual(t, acc2.Hash(), acc3.Account.Hash())
	assert.NotEqual(t, val2.Hash(), val3.Validator.Hash())
}

func TestBlockHeightByStamp(t *testing.T) {
	setup(t)

	height, ok := tSandbox.FindBlockHeightByStamp(hash.GenerateTestStamp())
	assert.Zero(t, height)
	assert.False(t, ok)

	lastHeight, _ := tStore.LastCertificate()
	lastHash := tSandbox.store.BlockHash(lastHeight)
	height, ok = tSandbox.FindBlockHeightByStamp(lastHash.Stamp())
	assert.Equal(t, height, lastHeight)
	assert.True(t, ok)
}

func TestBlockHashByStamp(t *testing.T) {
	setup(t)

	h, ok := tSandbox.FindBlockHashByStamp(hash.GenerateTestStamp())
	assert.True(t, h.IsUndef())
	assert.False(t, ok)

	lastHeight, _ := tStore.LastCertificate()
	lastHash := tSandbox.store.BlockHash(lastHeight)
	h, ok = tSandbox.FindBlockHashByStamp(lastHash.Stamp())
	assert.True(t, h.EqualsTo(lastHash))
	assert.True(t, ok)
}
