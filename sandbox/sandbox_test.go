package sandbox

import (
	"testing"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/sortition"
	"github.com/pactus-project/pactus/store"
	"github.com/pactus-project/pactus/types/account"
	"github.com/pactus-project/pactus/types/param"
	"github.com/pactus-project/pactus/types/validator"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
)

type testData struct {
	*testsuite.TestSuite

	signers []crypto.Signer
	store   *store.MockStore
	sandbox *sandbox
}

func setup(t *testing.T) *testData {
	ts := testsuite.NewTestSuite(t)
	store := store.MockingStore(ts)
	params := param.DefaultParams()
	params.TransactionToLiveInterval = 64

	committee, signers := ts.GenerateTestCommittee(21)
	acc := account.NewAccount(0)
	acc.AddToBalance(21 * 1e14)
	store.UpdateAccount(crypto.TreasuryAddress, acc)

	totalPower := int64(0)
	for _, val := range committee.Validators() {
		// For testing purpose we create some test accounts first.
		// Account number is the validator number plus one,
		// since account #0 is the Treasury account.
		acc := account.NewAccount(val.Number() + 1)
		store.UpdateValidator(val)
		store.UpdateAccount(val.Address(), acc)

		totalPower += val.Power()
	}

	sandbox := NewSandbox(store, params, committee, totalPower).(*sandbox)

	assert.Equal(t, sandbox.CurrentHeight(), uint32(1))
	lastHeight := ts.RandUint32(144) + 21
	for i := uint32(1); i < lastHeight; i++ {
		b := ts.GenerateTestBlock(nil, nil)
		c := ts.GenerateTestCertificate(b.Hash())
		store.SaveBlock(i, b, c)
	}
	assert.Equal(t, sandbox.CurrentHeight(), lastHeight)
	assert.Equal(t, sandbox.Params(), params)

	return &testData{
		TestSuite: ts,
		signers:   signers,
		store:     store,
		sandbox:   sandbox,
	}
}

func TestAccountChange(t *testing.T) {
	td := setup(t)

	t.Run("Should returns nil for invalid address", func(t *testing.T) {
		invAddr := td.RandomAddress()
		assert.Nil(t, td.sandbox.Account(invAddr))

		td.sandbox.IterateAccounts(func(_ crypto.Address, _ *account.Account, _ bool) {
			panic("should be empty")
		})
	})

	t.Run("Retrieve an account from store and update it", func(t *testing.T) {
		acc, signer := td.GenerateTestAccount(td.RandInt32(10000))
		addr := signer.Address()
		bal := acc.Balance()
		seq := acc.Sequence()
		td.store.UpdateAccount(addr, acc)

		sbAcc1 := td.sandbox.Account(addr)
		assert.Equal(t, acc, sbAcc1)

		sbAcc1.IncSequence()
		sbAcc1.AddToBalance(1)

		assert.False(t, td.sandbox.accounts[addr].updated)
		assert.Equal(t, td.sandbox.Account(addr).Balance(), bal)
		assert.Equal(t, td.sandbox.Account(addr).Sequence(), seq)
		td.sandbox.UpdateAccount(addr, sbAcc1)
		assert.True(t, td.sandbox.accounts[addr].updated)
		assert.Equal(t, td.sandbox.Account(addr).Balance(), bal+1)
		assert.Equal(t, td.sandbox.Account(addr).Sequence(), seq+1)

		t.Run("Update the same account again", func(t *testing.T) {
			sbAcc2 := td.sandbox.Account(addr)
			sbAcc2.IncSequence()
			sbAcc2.AddToBalance(1)

			assert.True(t, td.sandbox.accounts[addr].updated, "it is updated before")
			assert.Equal(t, td.sandbox.Account(addr).Balance(), bal+1)
			assert.Equal(t, td.sandbox.Account(addr).Sequence(), seq+1)
			td.sandbox.UpdateAccount(addr, sbAcc2)
			assert.True(t, td.sandbox.accounts[addr].updated)
			assert.Equal(t, td.sandbox.Account(addr).Balance(), bal+2)
			assert.Equal(t, td.sandbox.Account(addr).Sequence(), seq+2)
		})

		t.Run("Should be iterated", func(t *testing.T) {
			td.sandbox.IterateAccounts(func(a crypto.Address, acc *account.Account, updated bool) {
				assert.Equal(t, addr, a)
				assert.True(t, updated)
				assert.Equal(t, acc.Balance(), bal+2)
			})
		})
	})

	t.Run("Make new account", func(t *testing.T) {
		addr := td.RandomAddress()
		acc := td.sandbox.MakeNewAccount(addr)

		acc.IncSequence()
		acc.AddToBalance(1)

		td.sandbox.UpdateAccount(addr, acc)
		sbAcc := td.sandbox.Account(addr)
		assert.Equal(t, acc, sbAcc)

		t.Run("Should be iterated", func(t *testing.T) {
			td.sandbox.IterateAccounts(func(a crypto.Address, acc *account.Account, updated bool) {
				if a == addr {
					assert.True(t, updated)
					assert.Equal(t, acc.Balance(), int64(1))
				}
			})
		})
	})
}

func TestValidatorChange(t *testing.T) {
	td := setup(t)

	t.Run("Should returns nil for invalid address", func(t *testing.T) {
		invAddr := td.RandomAddress()
		assert.Nil(t, td.sandbox.Validator(invAddr))

		td.sandbox.IterateValidators(func(_ *validator.Validator, _ bool) {
			panic("should be empty")
		})
	})

	t.Run("Retrieve an validator from store and update it", func(t *testing.T) {
		val, _ := td.GenerateTestValidator(td.RandInt32(10000))
		addr := val.Address()
		stk := val.Stake()
		seq := val.Sequence()
		td.store.UpdateValidator(val)

		sbVal1 := td.sandbox.Validator(addr)
		assert.Equal(t, val.Hash(), sbVal1.Hash())

		sbVal1.IncSequence()
		sbVal1.AddToStake(1)

		assert.False(t, td.sandbox.validators[addr].updated)
		assert.Equal(t, td.sandbox.Validator(addr).Stake(), stk)
		assert.Equal(t, td.sandbox.Validator(addr).Sequence(), seq)
		td.sandbox.UpdateValidator(sbVal1)
		assert.True(t, td.sandbox.validators[sbVal1.Address()].updated)
		assert.Equal(t, td.sandbox.Validator(addr).Stake(), stk+1)
		assert.Equal(t, td.sandbox.Validator(addr).Sequence(), seq+1)

		t.Run("Update the same validator again", func(t *testing.T) {
			sbVal2 := td.sandbox.Validator(addr)
			sbVal2.IncSequence()
			sbVal2.AddToStake(1)

			assert.True(t, td.sandbox.validators[addr].updated, "it is updated before")
			assert.Equal(t, td.sandbox.Validator(addr).Stake(), stk+1)
			assert.Equal(t, td.sandbox.Validator(addr).Sequence(), seq+1)
			td.sandbox.UpdateValidator(sbVal2)
			assert.True(t, td.sandbox.validators[sbVal1.Address()].updated)
			assert.Equal(t, td.sandbox.Validator(addr).Stake(), stk+2)
			assert.Equal(t, td.sandbox.Validator(addr).Sequence(), seq+2)
		})

		t.Run("Should be iterated", func(t *testing.T) {
			td.sandbox.IterateValidators(func(val *validator.Validator, updated bool) {
				assert.True(t, updated)
				assert.Equal(t, val.Stake(), stk+2)
			})
		})
	})

	t.Run("Make new validator", func(t *testing.T) {
		pub, _ := td.RandomBLSKeyPair()
		val := td.sandbox.MakeNewValidator(pub)

		val.IncSequence()
		val.AddToStake(1)

		td.sandbox.UpdateValidator(val)
		sbVal := td.sandbox.Validator(val.Address())
		assert.Equal(t, val, sbVal)

		t.Run("Should be iterated", func(t *testing.T) {
			td.sandbox.IterateValidators(func(val *validator.Validator, updated bool) {
				if val.PublicKey() == pub {
					assert.True(t, updated)
					assert.Equal(t, val.Stake(), int64(1))
				}
			})
		})
	})
}

func TestTotalAccountCounter(t *testing.T) {
	td := setup(t)

	t.Run("Should update total account counter", func(t *testing.T) {
		assert.Equal(t, td.store.TotalAccounts(), int32(len(td.signers)+1))

		addr1 := td.RandomAddress()
		addr2 := td.RandomAddress()
		acc := td.sandbox.MakeNewAccount(addr1)
		assert.Equal(t, acc.Number(), int32(td.sandbox.Committee().Size()+1))
		acc2 := td.sandbox.MakeNewAccount(addr2)
		assert.Equal(t, acc2.Number(), int32(td.sandbox.Committee().Size()+2))
		assert.Equal(t, acc2.Balance(), int64(0))
	})
}

func TestTotalValidatorCounter(t *testing.T) {
	td := setup(t)

	t.Run("Should update total validator counter", func(t *testing.T) {
		assert.Equal(t, td.store.TotalValidators(), int32(td.sandbox.Committee().Size()))

		pub, _ := td.RandomBLSKeyPair()
		pub2, _ := td.RandomBLSKeyPair()
		val1 := td.sandbox.MakeNewValidator(pub)
		val1.UpdateLastBondingHeight(td.sandbox.CurrentHeight())
		assert.Equal(t, val1.Number(), int32(td.sandbox.Committee().Size()))
		assert.Equal(t, val1.LastBondingHeight(), td.sandbox.CurrentHeight())

		val2 := td.sandbox.MakeNewValidator(pub2)
		val2.UpdateLastBondingHeight(td.sandbox.CurrentHeight() + 1)
		assert.Equal(t, val2.Number(), int32(td.sandbox.Committee().Size()+1))
		assert.Equal(t, val2.LastBondingHeight(), td.sandbox.CurrentHeight()+1)
		assert.Equal(t, val2.Stake(), int64(0))
	})
}

func TestCreateDuplicated(t *testing.T) {
	td := setup(t)

	t.Run("Try creating duplicated account, Should panic", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("The code did not panic")
			}
		}()
		addr := crypto.TreasuryAddress
		td.sandbox.MakeNewAccount(addr)
	})

	t.Run("Try creating duplicated validator, Should panic", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("The code did not panic")
			}
		}()
		pub := td.signers[3].PublicKey()
		td.sandbox.MakeNewValidator(pub.(*bls.PublicKey))
	})
}

func TestUpdateFromOutsideTheSandbox(t *testing.T) {
	td := setup(t)

	t.Run("Try update an account from outside the sandbox, Should panic", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("The code did not panic")
			}
		}()
		acc, signer := td.GenerateTestAccount(td.RandInt32(0))
		td.sandbox.UpdateAccount(signer.Address(), acc)
	})

	t.Run("Try update a validator from outside the sandbox, Should panic", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("The code did not panic")
			}
		}()
		val, _ := td.GenerateTestValidator(td.RandInt32(0))
		td.sandbox.UpdateValidator(val)
	})
}

func TestAccountDeepCopy(t *testing.T) {
	td := setup(t)

	t.Run("non existing account", func(t *testing.T) {
		addr := td.RandomAddress()
		acc := td.sandbox.MakeNewAccount(addr)
		acc.IncSequence()

		assert.NotEqual(t, td.sandbox.Account(addr), acc)
	})

	t.Run("existing account", func(t *testing.T) {
		addr := crypto.TreasuryAddress
		acc := td.sandbox.Account(addr)
		acc.IncSequence()

		assert.NotEqual(t, td.sandbox.Account(addr), acc)
	})

	t.Run("sandbox account", func(t *testing.T) {
		addr := crypto.TreasuryAddress
		acc := td.sandbox.Account(addr)
		acc.IncSequence()

		assert.NotEqual(t, td.sandbox.Account(addr), acc)
		assert.NotEqual(t, acc.Sequence(), 1)
	})
}

func TestValidatorDeepCopy(t *testing.T) {
	td := setup(t)

	t.Run("non existing validator", func(t *testing.T) {
		pub, _ := td.RandomBLSKeyPair()
		acc := td.sandbox.MakeNewValidator(pub)
		acc.IncSequence()

		assert.NotEqual(t, td.sandbox.Validator(pub.Address()), acc)
	})

	val0, _ := td.store.ValidatorByNumber(0)
	addr := val0.Address()
	t.Run("existing validator", func(t *testing.T) {
		acc := td.sandbox.Validator(addr)
		acc.IncSequence()

		assert.NotEqual(t, td.sandbox.Validator(addr), acc)
	})

	t.Run("sandbox validator", func(t *testing.T) {
		acc := td.sandbox.Validator(addr)
		acc.IncSequence()

		assert.NotEqual(t, td.sandbox.Validator(addr), acc)
		assert.NotEqual(t, acc.Sequence(), 1)
	})
}

func TestRecentBlockByStamp(t *testing.T) {
	td := setup(t)

	h, b := td.sandbox.RecentBlockByStamp(td.RandomStamp())
	assert.Zero(t, h)
	assert.Nil(t, b)

	lastHeight, _ := td.store.LastCertificate()
	lastHash := td.sandbox.store.BlockHash(lastHeight)
	h, b = td.sandbox.RecentBlockByStamp(lastHash.Stamp())
	assert.Equal(t, h, lastHeight)
	assert.Equal(t, b.Hash(), lastHash)
}

func TestPowerDelta(t *testing.T) {
	td := setup(t)

	assert.Zero(t, td.sandbox.PowerDelta())
	td.sandbox.UpdatePowerDelta(1)
	assert.Equal(t, td.sandbox.PowerDelta(), int64(1))
	td.sandbox.UpdatePowerDelta(-1)
	assert.Zero(t, td.sandbox.PowerDelta())
}

func TestVerifyProof(t *testing.T) {
	td := setup(t)

	lastHeight, _ := td.store.LastCertificate()
	vals := td.sandbox.committee.Validators()

	// Try to evaluate a valid sortition
	var validProof sortition.Proof
	var validStamp hash.Stamp
	var validVal *validator.Validator
	for i := lastHeight; i > 0; i-- {
		block := td.store.Blocks[i]
		for i, signer := range td.signers {
			ok, proof := sortition.EvaluateSortition(
				block.Header().SortitionSeed(), signer,
				td.sandbox.totalPower, vals[i].Power())

			if ok {
				validProof = proof
				validStamp = block.Stamp()
				validVal = vals[i]
			}
		}
	}

	t.Run("invalid proof", func(t *testing.T) {
		invalidProof := td.RandomProof()
		assert.False(t, td.sandbox.VerifyProof(validStamp, invalidProof, validVal))
	})
	t.Run("invalid stamp", func(t *testing.T) {
		invalidStamp := td.RandomStamp()
		assert.False(t, td.sandbox.VerifyProof(invalidStamp, validProof, validVal))
	})

	t.Run("genesis stamp", func(t *testing.T) {
		invalidStamp := hash.UndefHash.Stamp()
		assert.False(t, td.sandbox.VerifyProof(invalidStamp, validProof, validVal))
	})

	t.Run("Ok", func(t *testing.T) {
		assert.True(t, td.sandbox.VerifyProof(validStamp, validProof, validVal))
	})
}
