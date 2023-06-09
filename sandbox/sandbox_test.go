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

		tSandbox.IterateAccounts(func(_ crypto.Address, _ *AccountStatus) {
			panic("should be empty")
		})
	})

	t.Run("Retrieve an account from store and update it", func(t *testing.T) {
		acc, signer := account.GenerateTestAccount(888)
		addr := signer.Address()
		bal := acc.Balance()
		seq := acc.Sequence()
		tStore.UpdateAccount(addr, acc)

		sbAcc1 := tSandbox.Account(addr)
		assert.Equal(t, acc, sbAcc1)

		sbAcc1.IncSequence()
		sbAcc1.AddToBalance(1)

		assert.False(t, tSandbox.accounts[addr].Updated)
		assert.Equal(t, tSandbox.Account(addr).Balance(), bal)
		assert.Equal(t, tSandbox.Account(addr).Sequence(), seq)
		tSandbox.UpdateAccount(addr, sbAcc1)
		assert.True(t, tSandbox.accounts[addr].Updated)
		assert.Equal(t, tSandbox.Account(addr).Balance(), bal+1)
		assert.Equal(t, tSandbox.Account(addr).Sequence(), seq+1)

		t.Run("Update the same account again", func(t *testing.T) {
			sbAcc2 := tSandbox.Account(addr)
			sbAcc2.IncSequence()
			sbAcc2.AddToBalance(1)

			assert.True(t, tSandbox.accounts[addr].Updated, "it is updated before")
			assert.Equal(t, tSandbox.Account(addr).Balance(), bal+1)
			assert.Equal(t, tSandbox.Account(addr).Sequence(), seq+1)
			tSandbox.UpdateAccount(addr, sbAcc2)
			assert.True(t, tSandbox.accounts[addr].Updated)
			assert.Equal(t, tSandbox.Account(addr).Balance(), bal+2)
			assert.Equal(t, tSandbox.Account(addr).Sequence(), seq+2)
		})

		t.Run("Should be iterated", func(t *testing.T) {
			tSandbox.IterateAccounts(func(a crypto.Address, as *AccountStatus) {
				assert.Equal(t, addr, a)
				assert.True(t, as.Updated)
				assert.Equal(t, as.Account.Balance(), bal+2)
			})
		})
	})

	t.Run("Make new account", func(t *testing.T) {
		addr := crypto.GenerateTestAddress()
		acc := tSandbox.MakeNewAccount(addr)

		acc.IncSequence()
		acc.AddToBalance(1)

		tSandbox.UpdateAccount(addr, acc)
		sbAcc := tSandbox.Account(addr)
		assert.Equal(t, acc, sbAcc)

		t.Run("Should be iterated", func(t *testing.T) {
			tSandbox.IterateAccounts(func(a crypto.Address, as *AccountStatus) {
				if a == addr {
					assert.True(t, as.Updated)
					assert.Equal(t, as.Account.Balance(), int64(1))
				}
			})
		})
	})
}

func TestValidatorChange(t *testing.T) {
	setup(t)

	t.Run("Should returns nil for invalid address", func(t *testing.T) {
		invAddr := crypto.GenerateTestAddress()
		assert.Nil(t, tSandbox.Validator(invAddr))

		tSandbox.IterateValidators(func(_ *ValidatorStatus) {
			panic("should be empty")
		})
	})

	t.Run("Retrieve an validator from store and update it", func(t *testing.T) {
		val, _ := validator.GenerateTestValidator(888)
		addr := val.Address()
		stk := val.Stake()
		seq := val.Sequence()
		tStore.UpdateValidator(val)

		sbVal1 := tSandbox.Validator(addr)
		assert.Equal(t, val.Hash(), sbVal1.Hash())

		sbVal1.IncSequence()
		sbVal1.AddToStake(1)

		assert.False(t, tSandbox.validators[addr].Updated)
		assert.Equal(t, tSandbox.Validator(addr).Stake(), stk)
		assert.Equal(t, tSandbox.Validator(addr).Sequence(), seq)
		tSandbox.UpdateValidator(sbVal1)
		assert.True(t, tSandbox.validators[sbVal1.Address()].Updated)
		assert.Equal(t, tSandbox.Validator(addr).Stake(), stk+1)
		assert.Equal(t, tSandbox.Validator(addr).Sequence(), seq+1)

		t.Run("Update the same validator again", func(t *testing.T) {
			sbVal2 := tSandbox.Validator(addr)
			sbVal2.IncSequence()
			sbVal2.AddToStake(1)

			assert.True(t, tSandbox.validators[addr].Updated, "it is updated before")
			assert.Equal(t, tSandbox.Validator(addr).Stake(), stk+1)
			assert.Equal(t, tSandbox.Validator(addr).Sequence(), seq+1)
			tSandbox.UpdateValidator(sbVal2)
			assert.True(t, tSandbox.validators[sbVal1.Address()].Updated)
			assert.Equal(t, tSandbox.Validator(addr).Stake(), stk+2)
			assert.Equal(t, tSandbox.Validator(addr).Sequence(), seq+2)
		})

		t.Run("Should be iterated", func(t *testing.T) {
			tSandbox.IterateValidators(func(vs *ValidatorStatus) {
				assert.True(t, vs.Updated)
				assert.Equal(t, vs.Validator.Stake(), stk+2)
			})
		})
	})

	t.Run("Make new validator", func(t *testing.T) {
		pub, _ := bls.GenerateTestKeyPair()
		val := tSandbox.MakeNewValidator(pub)

		val.IncSequence()
		val.AddToStake(1)

		tSandbox.UpdateValidator(val)
		sbVal := tSandbox.Validator(val.Address())
		assert.Equal(t, val, sbVal)

		t.Run("Should be iterated", func(t *testing.T) {
			tSandbox.IterateValidators(func(vs *ValidatorStatus) {
				if vs.Validator.PublicKey() == pub {
					assert.True(t, vs.Updated)
					assert.Equal(t, vs.Validator.Stake(), int64(1))
				}
			})
		})
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
