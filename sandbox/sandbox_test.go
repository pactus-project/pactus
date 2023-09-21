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
	t.Helper()

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

	lastHeight := uint32(21)
	for i := uint32(1); i < lastHeight; i++ {
		b := ts.GenerateTestBlock()
		c := ts.GenerateTestCertificate()
		store.SaveBlock(i, b, c)
	}
	sandbox := NewSandbox(store.LastHeight,
		store, params, committee, totalPower).(*sandbox)
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
		invAddr := td.RandAccAddress()
		assert.Nil(t, td.sandbox.Account(invAddr))

		td.sandbox.IterateAccounts(func(_ crypto.Address, _ *account.Account, _ bool) {
			panic("should be empty")
		})
	})

	t.Run("Retrieve an account from store and update it", func(t *testing.T) {
		acc, signer := td.GenerateTestAccount(td.RandInt32(10000))
		addr := signer.Address()
		bal := acc.Balance()
		td.store.UpdateAccount(addr, acc)

		sbAcc1 := td.sandbox.Account(addr)
		assert.Equal(t, acc, sbAcc1)

		sbAcc1.AddToBalance(1)

		assert.False(t, td.sandbox.accounts[addr].updated)
		assert.Equal(t, td.sandbox.Account(addr).Balance(), bal)
		td.sandbox.UpdateAccount(addr, sbAcc1)
		assert.True(t, td.sandbox.accounts[addr].updated)
		assert.Equal(t, td.sandbox.Account(addr).Balance(), bal+1)

		t.Run("Update the same account again", func(t *testing.T) {
			sbAcc2 := td.sandbox.Account(addr)
			sbAcc2.AddToBalance(1)

			assert.True(t, td.sandbox.accounts[addr].updated, "it is updated before")
			assert.Equal(t, td.sandbox.Account(addr).Balance(), bal+1)
			td.sandbox.UpdateAccount(addr, sbAcc2)
			assert.True(t, td.sandbox.accounts[addr].updated)
			assert.Equal(t, td.sandbox.Account(addr).Balance(), bal+2)
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
		addr := td.RandAccAddress()
		acc := td.sandbox.MakeNewAccount(addr)

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

func TestAnyRecentTransaction(t *testing.T) {
	td := setup(t)

	randTx1, _ := td.GenerateTestTransferTx()
	randTx2, _ := td.GenerateTestTransferTx()
	td.sandbox.CommitTransaction(randTx1)
	td.sandbox.CommitTransaction(randTx2)

	assert.True(t, td.sandbox.AnyRecentTransaction(randTx1.ID()))
	assert.True(t, td.sandbox.AnyRecentTransaction(randTx2.ID()))

	totalTxFees := randTx1.Fee() + randTx2.Fee()
	assert.Equal(t, td.sandbox.AccumulatedFee(), totalTxFees)
}

func TestValidatorChange(t *testing.T) {
	td := setup(t)

	t.Run("Should returns nil for invalid address", func(t *testing.T) {
		invAddr := td.RandAccAddress()
		assert.Nil(t, td.sandbox.Validator(invAddr))

		td.sandbox.IterateValidators(func(_ *validator.Validator, _ bool, _ bool) {
			panic("should be empty")
		})
	})

	t.Run("Retrieve an validator from store and update it", func(t *testing.T) {
		val, _ := td.GenerateTestValidator(td.RandInt32(10000))
		addr := val.Address()
		stk := val.Stake()
		td.store.UpdateValidator(val)

		sbVal1 := td.sandbox.Validator(addr)
		assert.Equal(t, val.Hash(), sbVal1.Hash())

		sbVal1.AddToStake(1)

		assert.False(t, td.sandbox.validators[addr].updated)
		assert.Equal(t, td.sandbox.Validator(addr).Stake(), stk)
		td.sandbox.UpdateValidator(sbVal1)
		assert.True(t, td.sandbox.validators[sbVal1.Address()].updated)
		assert.Equal(t, td.sandbox.Validator(addr).Stake(), stk+1)

		t.Run("Update the same validator again", func(t *testing.T) {
			sbVal2 := td.sandbox.Validator(addr)
			sbVal2.AddToStake(1)

			assert.True(t, td.sandbox.validators[addr].updated, "it is updated before")
			assert.Equal(t, td.sandbox.Validator(addr).Stake(), stk+1)
			td.sandbox.UpdateValidator(sbVal2)
			assert.True(t, td.sandbox.validators[sbVal1.Address()].updated)
			assert.Equal(t, td.sandbox.Validator(addr).Stake(), stk+2)
		})

		t.Run("Should be iterated", func(t *testing.T) {
			td.sandbox.IterateValidators(func(val *validator.Validator, updated bool, joined bool) {
				assert.True(t, updated)
				assert.False(t, joined)
				assert.Equal(t, val.Stake(), stk+2)
			})
		})
	})

	t.Run("Make new validator", func(t *testing.T) {
		pub, _ := td.RandBLSKeyPair()
		val := td.sandbox.MakeNewValidator(pub)

		val.AddToStake(1)

		td.sandbox.UpdateValidator(val)
		sbVal := td.sandbox.Validator(val.Address())
		assert.Equal(t, val, sbVal)

		t.Run("Should be iterated", func(t *testing.T) {
			td.sandbox.IterateValidators(func(val *validator.Validator, updated bool, joined bool) {
				if val.PublicKey() == pub {
					assert.True(t, updated)
					assert.False(t, joined)
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

		addr1 := td.RandAccAddress()
		addr2 := td.RandAccAddress()
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

		pub, _ := td.RandBLSKeyPair()
		pub2, _ := td.RandBLSKeyPair()
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
		addr := td.RandAccAddress()
		acc := td.sandbox.MakeNewAccount(addr)
		acc.AddToBalance(1)

		assert.NotEqual(t, td.sandbox.Account(addr), acc)
	})

	t.Run("existing account", func(t *testing.T) {
		addr := crypto.TreasuryAddress
		acc := td.sandbox.Account(addr)
		acc.AddToBalance(1)

		assert.NotEqual(t, td.sandbox.Account(addr), acc)
	})

	t.Run("sandbox account", func(t *testing.T) {
		addr := crypto.TreasuryAddress
		acc := td.sandbox.Account(addr)
		acc.AddToBalance(1)

		assert.NotEqual(t, td.sandbox.Account(addr), acc)
	})
}

func TestValidatorDeepCopy(t *testing.T) {
	td := setup(t)

	t.Run("non existing validator", func(t *testing.T) {
		pub, _ := td.RandBLSKeyPair()
		val := td.sandbox.MakeNewValidator(pub)
		val.AddToStake(1)

		assert.NotEqual(t, td.sandbox.Validator(pub.Address()), val)
	})

	val0, _ := td.store.ValidatorByNumber(0)
	addr := val0.Address()
	t.Run("existing validator", func(t *testing.T) {
		val := td.sandbox.Validator(addr)
		val.AddToStake(1)

		assert.NotEqual(t, td.sandbox.Validator(addr), val)
	})

	t.Run("sandbox validator", func(t *testing.T) {
		val := td.sandbox.Validator(addr)
		val.AddToStake(1)

		assert.NotEqual(t, td.sandbox.Validator(addr), val)
	})
}

func TestRecentBlockByStamp(t *testing.T) {
	td := setup(t)

	h, b := td.sandbox.RecentBlockByStamp(td.RandStamp())
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
		invalidProof := td.RandProof()
		assert.False(t, td.sandbox.VerifyProof(validStamp, invalidProof, validVal))
	})
	t.Run("invalid stamp", func(t *testing.T) {
		invalidStamp := td.RandStamp()
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
