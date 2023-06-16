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

		tSandbox.IterateAccounts(func(_ crypto.Address, _ *account.Account, _ bool) {
			panic("should be empty")
		})
	})

	t.Run("Retrieve an account from store and update it", func(t *testing.T) {
		acc, signer := account.GenerateTestAccount(util.RandInt32(0))
		addr := signer.Address()
		bal := acc.Balance()
		seq := acc.Sequence()
		tStore.UpdateAccount(addr, acc)

		sbAcc1 := tSandbox.Account(addr)
		assert.Equal(t, acc, sbAcc1)

		sbAcc1.IncSequence()
		sbAcc1.AddToBalance(1)

		assert.False(t, tSandbox.accounts[addr].updated)
		assert.Equal(t, tSandbox.Account(addr).Balance(), bal)
		assert.Equal(t, tSandbox.Account(addr).Sequence(), seq)
		tSandbox.UpdateAccount(addr, sbAcc1)
		assert.True(t, tSandbox.accounts[addr].updated)
		assert.Equal(t, tSandbox.Account(addr).Balance(), bal+1)
		assert.Equal(t, tSandbox.Account(addr).Sequence(), seq+1)

		t.Run("Update the same account again", func(t *testing.T) {
			sbAcc2 := tSandbox.Account(addr)
			sbAcc2.IncSequence()
			sbAcc2.AddToBalance(1)

			assert.True(t, tSandbox.accounts[addr].updated, "it is updated before")
			assert.Equal(t, tSandbox.Account(addr).Balance(), bal+1)
			assert.Equal(t, tSandbox.Account(addr).Sequence(), seq+1)
			tSandbox.UpdateAccount(addr, sbAcc2)
			assert.True(t, tSandbox.accounts[addr].updated)
			assert.Equal(t, tSandbox.Account(addr).Balance(), bal+2)
			assert.Equal(t, tSandbox.Account(addr).Sequence(), seq+2)
		})

		t.Run("Should be iterated", func(t *testing.T) {
			tSandbox.IterateAccounts(func(a crypto.Address, acc *account.Account, updated bool) {
				assert.Equal(t, addr, a)
				assert.True(t, updated)
				assert.Equal(t, acc.Balance(), bal+2)
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
			tSandbox.IterateAccounts(func(a crypto.Address, acc *account.Account, updated bool) {
				if a == addr {
					assert.True(t, updated)
					assert.Equal(t, acc.Balance(), int64(1))
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

		tSandbox.IterateValidators(func(_ *validator.Validator, _ bool) {
			panic("should be empty")
		})
	})

	t.Run("Retrieve an validator from store and update it", func(t *testing.T) {
		val, _ := validator.GenerateTestValidator(util.RandInt32(0))
		addr := val.Address()
		stk := val.Stake()
		seq := val.Sequence()
		tStore.UpdateValidator(val)

		sbVal1 := tSandbox.Validator(addr)
		assert.Equal(t, val.Hash(), sbVal1.Hash())

		sbVal1.IncSequence()
		sbVal1.AddToStake(1)

		assert.False(t, tSandbox.validators[addr].updated)
		assert.Equal(t, tSandbox.Validator(addr).Stake(), stk)
		assert.Equal(t, tSandbox.Validator(addr).Sequence(), seq)
		tSandbox.UpdateValidator(sbVal1)
		assert.True(t, tSandbox.validators[sbVal1.Address()].updated)
		assert.Equal(t, tSandbox.Validator(addr).Stake(), stk+1)
		assert.Equal(t, tSandbox.Validator(addr).Sequence(), seq+1)

		t.Run("Update the same validator again", func(t *testing.T) {
			sbVal2 := tSandbox.Validator(addr)
			sbVal2.IncSequence()
			sbVal2.AddToStake(1)

			assert.True(t, tSandbox.validators[addr].updated, "it is updated before")
			assert.Equal(t, tSandbox.Validator(addr).Stake(), stk+1)
			assert.Equal(t, tSandbox.Validator(addr).Sequence(), seq+1)
			tSandbox.UpdateValidator(sbVal2)
			assert.True(t, tSandbox.validators[sbVal1.Address()].updated)
			assert.Equal(t, tSandbox.Validator(addr).Stake(), stk+2)
			assert.Equal(t, tSandbox.Validator(addr).Sequence(), seq+2)
		})

		t.Run("Should be iterated", func(t *testing.T) {
			tSandbox.IterateValidators(func(val *validator.Validator, updated bool) {
				assert.True(t, updated)
				assert.Equal(t, val.Stake(), stk+2)
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
			tSandbox.IterateValidators(func(val *validator.Validator, updated bool) {
				if val.PublicKey() == pub {
					assert.True(t, updated)
					assert.Equal(t, val.Stake(), int64(1))
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
		acc, signer := account.GenerateTestAccount(util.RandInt32(0))
		tSandbox.UpdateAccount(signer.Address(), acc)
	})

	t.Run("Try update a validator from outside the sandbox, Should panic", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("The code did not panic")
			}
		}()
		val, _ := validator.GenerateTestValidator(util.RandInt32(0))
		tSandbox.UpdateValidator(val)
	})
}

func TestAccountDeepCopy(t *testing.T) {
	setup(t)

	t.Run("non existing account", func(t *testing.T) {
		addr := crypto.GenerateTestAddress()
		acc := tSandbox.MakeNewAccount(addr)
		acc.IncSequence()

		assert.NotEqual(t, tSandbox.Account(addr), acc)
	})

	t.Run("existing account", func(t *testing.T) {
		addr := crypto.TreasuryAddress
		acc := tSandbox.Account(addr)
		acc.IncSequence()

		assert.NotEqual(t, tSandbox.Account(addr), acc)
	})

	t.Run("sandbox account", func(t *testing.T) {
		addr := crypto.TreasuryAddress
		acc := tSandbox.Account(addr)
		acc.IncSequence()

		assert.NotEqual(t, tSandbox.Account(addr), acc)
		assert.NotEqual(t, acc.Sequence(), 1)
	})
}

func TestValidatorDeepCopy(t *testing.T) {
	setup(t)

	t.Run("non existing account", func(t *testing.T) {
		pub, _ := bls.GenerateTestKeyPair()
		acc := tSandbox.MakeNewValidator(pub)
		acc.IncSequence()

		assert.NotEqual(t, tSandbox.Validator(pub.Address()), acc)
	})

	val0, _ := tStore.ValidatorByNumber(0)
	addr := val0.Address()
	t.Run("existing validator", func(t *testing.T) {
		acc := tSandbox.Validator(addr)
		acc.IncSequence()

		assert.NotEqual(t, tSandbox.Validator(addr), acc)
	})

	t.Run("sandbox validator", func(t *testing.T) {
		acc := tSandbox.Validator(addr)
		acc.IncSequence()

		assert.NotEqual(t, tSandbox.Validator(addr), acc)
		assert.NotEqual(t, acc.Sequence(), 1)
	})
}

func TestRecentBlockByStamp(t *testing.T) {
	setup(t)

	h, b := tSandbox.RecentBlockByStamp(hash.GenerateTestStamp())
	assert.Zero(t, h)
	assert.Nil(t, b)

	lastHeight, _ := tStore.LastCertificate()
	lastHash := tSandbox.store.BlockHash(lastHeight)
	h, b = tSandbox.RecentBlockByStamp(lastHash.Stamp())
	assert.Equal(t, h, lastHeight)
	assert.Equal(t, b.Hash(), lastHash)
}
