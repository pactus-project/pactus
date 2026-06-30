package sandbox

import (
	"testing"

	"github.com/pactus-project/pactus/committee"
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/genesis"
	"github.com/pactus-project/pactus/state/param"
	"github.com/pactus-project/pactus/store"
	"github.com/pactus-project/pactus/types/account"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/validator"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type testData struct {
	*testsuite.TestSuite

	fakeCommittee *committee.FakeCommittee
	fakeStore     *store.FakeStore
	fakeParams    *param.Params
	sbx           *sandbox
}

func setup(t *testing.T) *testData {
	t.Helper()

	ts := testsuite.NewTestSuite(t)
	fakeStore := store.NewFakeStore(ts)
	fakeCommittee := committee.NewFakeCommittee(ts)

	genDoc := genesis.MainnetGenesis()
	fakeParams := param.FromGenesis(genDoc)
	totalPower := ts.RandInt64()
	fakeHeight := ts.RandHeight()

	sbx := NewSandbox(fakeHeight, fakeStore, fakeParams, fakeCommittee, totalPower).(*sandbox)
	assert.Equal(t, fakeHeight+1, sbx.CurrentHeight())
	assert.Equal(t, fakeParams, sbx.Params())

	return &testData{
		TestSuite:  ts,
		fakeStore:  fakeStore,
		fakeParams: fakeParams,
		sbx:        sbx,
	}
}

func TestAccountChange(t *testing.T) {
	td := setup(t)

	t.Run("Should returns nil for invalid address", func(t *testing.T) {
		invAddr := td.RandAccAddress()
		td.fakeStore.EXPECT().Account(invAddr).Return(nil, store.ErrNotFound).Times(1)

		assert.Nil(t, td.sbx.Account(invAddr))
		td.sbx.IterateAccounts(func(_ crypto.Address, _ *account.Account, _ bool) {
			panic("should be empty")
		})
	})

	t.Run("Retrieve an account from store and update it", func(t *testing.T) {
		acc, addr := td.GenerateTestAccount()
		bal := acc.Balance()
		td.fakeStore.EXPECT().Account(addr).Return(acc, nil).Times(1)

		sbAcc1 := td.sbx.Account(addr)
		assert.Equal(t, acc, sbAcc1)

		sbAcc1.AddToBalance(1)

		assert.False(t, td.sbx.accounts[addr].updated)
		assert.Equal(t, bal, td.sbx.Account(addr).Balance())
		td.sbx.UpdateAccount(addr, sbAcc1)
		assert.True(t, td.sbx.accounts[addr].updated)
		assert.Equal(t, bal+1, td.sbx.Account(addr).Balance())

		t.Run("Update the same account again", func(t *testing.T) {
			sbAcc2 := td.sbx.Account(addr)
			sbAcc2.AddToBalance(1)

			assert.True(t, td.sbx.accounts[addr].updated, "it is updated before")
			assert.Equal(t, bal+1, td.sbx.Account(addr).Balance())
			td.sbx.UpdateAccount(addr, sbAcc2)
			assert.True(t, td.sbx.accounts[addr].updated)
			assert.Equal(t, bal+2, td.sbx.Account(addr).Balance())
		})

		t.Run("Should be iterated", func(t *testing.T) {
			td.sbx.IterateAccounts(func(a crypto.Address, acc *account.Account, updated bool) {
				assert.Equal(t, addr, a)
				assert.True(t, updated)
				assert.Equal(t, bal+2, acc.Balance())
			})
		})
	})

	t.Run("Make new account", func(t *testing.T) {
		addr := td.RandAccAddress()
		td.fakeStore.EXPECT().HasAccount(addr).Return(false).Times(1)

		acc := td.sbx.MakeNewAccount(addr)

		acc.AddToBalance(1)

		td.sbx.UpdateAccount(addr, acc)
		sbAcc := td.sbx.Account(addr)
		assert.Equal(t, acc, sbAcc)

		t.Run("Should be iterated", func(t *testing.T) {
			td.sbx.IterateAccounts(func(a crypto.Address, acc *account.Account, updated bool) {
				if a == addr {
					assert.True(t, updated)
					assert.Equal(t, amount.Amount(1), acc.Balance())
				}
			})
		})
	})
}

func TestRecentTransaction(t *testing.T) {
	td := setup(t)

	randTx1 := td.GenerateTestTransferTx()
	randTx2 := td.GenerateTestTransferTx()
	td.sbx.CommitTransaction(randTx1)
	td.sbx.CommitTransaction(randTx2)

	assert.True(t, td.sbx.RecentTransaction(randTx1.ID()))
	assert.True(t, td.sbx.RecentTransaction(randTx2.ID()))

	totalTxFees := randTx1.Fee() + randTx2.Fee()
	assert.Equal(t, totalTxFees, td.sbx.AccumulatedFee())
}

func TestValidatorChange(t *testing.T) {
	td := setup(t)

	t.Run("Should returns nil for invalid address", func(t *testing.T) {
		invAddr := td.RandAccAddress()
		td.fakeStore.EXPECT().Validator(invAddr).Return(nil, store.ErrNotFound).Times(1)

		assert.Nil(t, td.sbx.Validator(invAddr))

		td.sbx.IterateValidators(func(_ *validator.Validator, _ bool, _ bool) {
			panic("should be empty")
		})
	})

	t.Run("Retrieve an validator from store and update it", func(t *testing.T) {
		val := td.GenerateTestValidator()
		addr := val.Address()
		stk := val.Stake()
		td.fakeStore.EXPECT().Validator(addr).Return(val, nil).Times(1)

		sbVal1 := td.sbx.Validator(addr)
		assert.Equal(t, val.Hash(), sbVal1.Hash())

		sbVal1.AddToStake(1)

		assert.False(t, td.sbx.validators[addr].updated)
		assert.Equal(t, stk, td.sbx.Validator(addr).Stake())
		td.sbx.UpdateValidator(sbVal1)
		assert.True(t, td.sbx.validators[sbVal1.Address()].updated)
		assert.Equal(t, stk+1, td.sbx.Validator(addr).Stake())

		t.Run("Update the same validator again", func(t *testing.T) {
			sbVal2 := td.sbx.Validator(addr)
			sbVal2.AddToStake(1)

			assert.True(t, td.sbx.validators[addr].updated, "it is updated before")
			assert.Equal(t, stk+1, td.sbx.Validator(addr).Stake())
			td.sbx.UpdateValidator(sbVal2)
			assert.True(t, td.sbx.validators[sbVal1.Address()].updated)
			assert.Equal(t, stk+2, td.sbx.Validator(addr).Stake())
		})

		t.Run("Should be iterated", func(t *testing.T) {
			td.sbx.IterateValidators(func(val *validator.Validator, updated bool, joined bool) {
				assert.True(t, updated)
				assert.False(t, joined)
				assert.Equal(t, stk+2, val.Stake())
			})
		})
	})

	t.Run("Make new validator", func(t *testing.T) {
		pub, _ := td.RandBLSKeyPair()
		td.fakeStore.EXPECT().HasValidator(pub.ValidatorAddress()).Return(false).Times(1)

		val := td.sbx.MakeNewValidator(pub)

		val.AddToStake(1)

		td.sbx.UpdateValidator(val)
		sbVal := td.sbx.Validator(val.Address())
		assert.Equal(t, val, sbVal)

		t.Run("Should be iterated", func(t *testing.T) {
			td.sbx.IterateValidators(func(val *validator.Validator, updated bool, joined bool) {
				if val.PublicKey() == pub {
					assert.True(t, updated)
					assert.False(t, joined)
					assert.Equal(t, amount.Amount(1), val.Stake())
				}
			})
		})
	})
}

func TestTotalAccountCounter(t *testing.T) {
	td := setup(t)

	t.Run("Should update total account counter", func(t *testing.T) {
		totalAccs := td.fakeStore.TotalAccounts()

		td.fakeStore.EXPECT().HasAccount(gomock.Any()).Return(false).Times(2)

		acc1 := td.sbx.MakeNewAccount(td.RandAccAddress())
		assert.Equal(t, totalAccs, acc1.Number())

		acc2 := td.sbx.MakeNewAccount(td.RandAccAddress())
		assert.Equal(t, totalAccs+1, acc2.Number())
	})
}

func TestTotalValidatorCounter(t *testing.T) {
	td := setup(t)

	t.Run("Should update total validator counter", func(t *testing.T) {
		totalVals := td.fakeStore.TotalValidators()

		td.fakeStore.EXPECT().HasValidator(gomock.Any()).Return(false).Times(2)

		pub, _ := td.RandBLSKeyPair()
		pub2, _ := td.RandBLSKeyPair()
		val1 := td.sbx.MakeNewValidator(pub)
		assert.Equal(t, totalVals, val1.Number())

		val2 := td.sbx.MakeNewValidator(pub2)
		assert.Equal(t, totalVals+1, val2.Number())
	})
}

func TestCreateDuplicated(t *testing.T) {
	td := setup(t)

	t.Run("Try creating duplicated account, Should panic", func(t *testing.T) {
		assert.Panics(t, func() {
			addr := td.RandAccAddress()
			td.fakeStore.EXPECT().HasAccount(addr).Return(true).Times(1)

			td.sbx.MakeNewAccount(addr)
		})
	})

	t.Run("Try creating duplicated validator, Should panic", func(t *testing.T) {
		assert.Panics(t, func() {
			pub, _ := td.RandBLSKeyPair()
			td.fakeStore.EXPECT().HasValidator(pub.ValidatorAddress()).Return(true).Times(1)

			td.sbx.MakeNewValidator(pub)
		})
	})
}

func TestUpdateFromOutsideTheSandbox(t *testing.T) {
	td := setup(t)

	t.Run("Try update an account from outside the sandbox, Should panic", func(t *testing.T) {
		assert.Panics(t, func() {
			acc, addr := td.GenerateTestAccount()
			td.sbx.UpdateAccount(addr, acc)
		})
	})

	t.Run("Try update a validator from outside the sandbox, Should panic", func(t *testing.T) {
		assert.Panics(t, func() {
			val := td.GenerateTestValidator()
			td.sbx.UpdateValidator(val)
		})
	})
}

func TestAccountDeepCopy(t *testing.T) {
	td := setup(t)

	t.Run("non existing account", func(t *testing.T) {
		addr := td.RandAccAddress()
		td.fakeStore.EXPECT().HasAccount(addr).Return(false).Times(1)

		acc1 := td.sbx.MakeNewAccount(addr)

		acc2 := td.sbx.Account(addr)
		acc2.AddToBalance(1)

		assert.NotEqual(t, acc1, acc2)
	})

	t.Run("existing account", func(t *testing.T) {
		acc1, addr := td.GenerateTestAccount()
		td.fakeStore.EXPECT().Account(addr).Return(acc1, nil).Times(1)

		acc2 := td.sbx.Account(addr)
		acc2.AddToBalance(2)

		acc3 := td.sbx.Account(addr)
		acc3.AddToBalance(1)

		assert.NotEqual(t, acc1, acc2)
		assert.NotEqual(t, acc2, acc3)
	})
}

func TestValidatorDeepCopy(t *testing.T) {
	td := setup(t)

	t.Run("non existing validator", func(t *testing.T) {
		pub, _ := td.RandBLSKeyPair()
		td.fakeStore.EXPECT().HasValidator(pub.ValidatorAddress()).Return(false).Times(1)

		val1 := td.sbx.MakeNewValidator(pub)

		val2 := td.sbx.Validator(pub.ValidatorAddress())
		val2.AddToStake(1)

		assert.NotEqual(t, val1, val2)
	})

	t.Run("existing validator", func(t *testing.T) {
		val1 := td.GenerateTestValidator()
		addr := val1.Address()
		td.fakeStore.EXPECT().Validator(addr).Return(val1, nil).Times(1)

		val2 := td.sbx.Validator(addr)
		val2.AddToStake(2)

		val3 := td.sbx.Validator(addr)
		val3.AddToStake(1)

		assert.NotEqual(t, val1, val2)
		assert.NotEqual(t, val2, val3)
	})
}

func TestPowerDelta(t *testing.T) {
	td := setup(t)

	assert.Zero(t, td.sbx.PowerDelta())
	td.sbx.UpdatePowerDelta(1)
	assert.Equal(t, int64(1), td.sbx.PowerDelta())
	td.sbx.UpdatePowerDelta(-1)
	assert.Zero(t, td.sbx.PowerDelta())
}

// func TestVerifyProof(t *testing.T) {
// 	td := setup(t)

// 	lastCert := td.store.LastCertificate()
// 	lastHeight := lastCert.Height()
// 	vals := td.sbx.committee.Validators()

// 	// Try to evaluate a valid sortition
// 	var validProof sortition.Proof
// 	var validLockTime types.Height
// 	var validVal *validator.Validator
// 	for height := lastHeight; height > 0; height-- {
// 		block := td.store.Blocks[height]
// 		for index, valKey := range td.valKeys {
// 			ok, proof := sortition.EvaluateSortition(
// 				block.Header().SortitionSeed(), valKey.PrivateKey(),
// 				td.sbx.totalPower, vals[index].Power(),
// 			)

// 			if ok {
// 				validProof = proof
// 				validLockTime = height
// 				validVal = vals[index]
// 			}
// 		}
// 	}

// 	t.Run("invalid proof", func(t *testing.T) {
// 		invalidProof := td.RandProof()
// 		assert.False(t, td.sbx.VerifyProof(validLockTime, invalidProof, validVal))
// 	})
// 	t.Run("invalid height", func(t *testing.T) {
// 		assert.False(t, td.sbx.VerifyProof(td.RandHeight(), validProof, validVal))
// 	})

// 	t.Run("genesis block height", func(t *testing.T) {
// 		assert.False(t, td.sbx.VerifyProof(0, validProof, validVal))
// 	})

// 	t.Run("Ok", func(t *testing.T) {
// 		assert.True(t, td.sbx.VerifyProof(validLockTime, validProof, validVal))
// 	})
// }

// func TestJoinedToCommittee(t *testing.T) {
// 	td := setup(t)

// 	pub, _ := td.RandBLSKeyPair()
// 	addr := pub.ValidatorAddress()
// 	td.sbx.MakeNewValidator(pub)

// 	// assert.False(t, td.sbx.IsJoinedCommittee(td.RandValAddress()))
// 	assert.False(t, td.sbx.IsJoinedCommittee(addr))

// 	td.sbx.JoinToCommittee(addr)
// 	assert.True(t, td.sbx.IsJoinedCommittee(addr))
// }
