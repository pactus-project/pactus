package sandbox

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/committee"
	"github.com/zarbchain/zarb-go/store"
	"github.com/zarbchain/zarb-go/types/account"
	"github.com/zarbchain/zarb-go/types/block"
	"github.com/zarbchain/zarb-go/types/crypto"
	"github.com/zarbchain/zarb-go/types/crypto/bls"
	"github.com/zarbchain/zarb-go/types/crypto/hash"
	"github.com/zarbchain/zarb-go/types/param"
	"github.com/zarbchain/zarb-go/types/validator"
	"github.com/zarbchain/zarb-go/util"
)

var tSigners []crypto.Signer
var tStore *store.MockStore
var tSandbox *sandbox

func setup(t *testing.T) {
	tStore = store.MockingStore()
	params := param.DefaultParams()
	params.TransactionToLiveInterval = 64

	committee, signers := committee.GenerateTestCommittee(21)
	acc := account.NewAccount(crypto.TreasuryAddress, 0)
	acc.AddToBalance(21 * 1e14)

	tStore.UpdateAccount(acc)
	for _, val := range committee.Validators() {
		// For testing purpose we create some test accounts here
		// Account number is validator number plus one.
		// Since account #0 is Treasury account, so it makes scene.
		acc := account.NewAccount(val.Address(), val.Number()+1)
		tStore.UpdateValidator(val)
		tStore.UpdateAccount(acc)
	}

	tSigners = signers
	tSandbox = NewSandbox(tStore, params, committee).(*sandbox)

	assert.Equal(t, tSandbox.CurrentHeight(), int32(1))
	lastHeight := util.RandInt32(144) + 21
	for i := int32(1); i < lastHeight; i++ {
		b := block.GenerateTestBlock(nil, nil)
		c := block.GenerateTestCertificate(b.Hash())
		tStore.SaveBlock(i, b, c)
	}
	assert.Equal(t, tSandbox.CurrentHeight(), lastHeight)
	assert.Equal(t, tSandbox.FeeFraction(), params.FeeFraction)
	assert.Equal(t, tSandbox.MinFee(), params.MinimumFee)
	assert.Equal(t, tSandbox.TransactionToLiveInterval(), params.TransactionToLiveInterval)
	assert.Equal(t, tSandbox.CommitteeSize(), params.CommitteeSize)
	assert.Equal(t, tSandbox.BondInterval(), params.BondInterval)
	assert.Equal(t, tSandbox.UnbondInterval(), params.UnbondInterval)
}

func TestAccountChange(t *testing.T) {
	setup(t)

	t.Run("Should returns nil for invalid address", func(t *testing.T) {
		invAddr := crypto.GenerateTestAddress()
		assert.Nil(t, tSandbox.Account(invAddr))
	})

	t.Run("Retrieve an account from store, modify it and commit it", func(t *testing.T) {
		acc1, _ := account.GenerateTestAccount(888)
		tStore.UpdateAccount(acc1)

		acc1a := tSandbox.Account(acc1.Address())
		assert.Equal(t, acc1, acc1a)

		acc1a.IncSequence()
		acc1a.AddToBalance(acc1a.Balance() + 1)

		assert.False(t, tSandbox.accounts[acc1a.Address()].Updated)
		tSandbox.UpdateAccount(acc1a)
		assert.True(t, tSandbox.accounts[acc1a.Address()].Updated)
	})

	t.Run("Make new account", func(t *testing.T) {
		addr := crypto.GenerateTestAddress()
		acc2 := tSandbox.MakeNewAccount(addr)

		acc2.IncSequence()
		acc2.AddToBalance(acc2.Balance() + 1)

		tSandbox.UpdateAccount(acc2)
		acc22 := tSandbox.Account(acc2.Address())
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
		acc, _ := account.GenerateTestAccount(999)
		tSandbox.UpdateAccount(acc)
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

	assert.Equal(t, tSandbox.BlockHeightByStamp(hash.GenerateTestStamp()), int32(0))

	height, _ := tStore.LastCertificate()
	hash := tSandbox.store.BlockHash(height)
	assert.Equal(t, tSandbox.BlockHeightByStamp(hash.Stamp()), height)
}

func TestBlockHashByStamp(t *testing.T) {
	setup(t)

	assert.True(t, tSandbox.BlockHashByStamp(hash.GenerateTestStamp()).IsUndef())

	height, _ := tStore.LastCertificate()
	hash := tSandbox.store.BlockHash(height)
	assert.True(t, tSandbox.BlockHashByStamp(hash.Stamp()).EqualsTo(hash))
}
