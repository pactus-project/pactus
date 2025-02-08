package sandbox

import (
	"testing"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/genesis"
	"github.com/pactus-project/pactus/sortition"
	"github.com/pactus-project/pactus/state/param"
	"github.com/pactus-project/pactus/store"
	"github.com/pactus-project/pactus/types/account"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/validator"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
)

type testData struct {
	*testsuite.TestSuite

	valKeys []*bls.ValidatorKey
	store   *store.MockStore
	sbx     *sandbox
}

func setup(t *testing.T) *testData {
	t.Helper()

	ts := testsuite.NewTestSuite(t)
	mockStore := store.MockingStore(ts)
	params := genesis.DefaultGenesisParams()
	params.TransactionToLiveInterval = 64

	cmt, valKeys := ts.GenerateTestCommittee(21)
	acc := account.NewAccount(0)
	acc.AddToBalance(21 * 1e14)
	mockStore.UpdateAccount(crypto.TreasuryAddress, acc)

	totalPower := int64(0)
	for _, val := range cmt.Validators() {
		// For testing purpose, we create some test accounts first.
		// Account number is the validator number plus one,
		// since account #0 is the Treasury account.
		newAcc := account.NewAccount(val.Number() + 1)
		mockStore.UpdateValidator(val)
		mockStore.UpdateAccount(val.Address(), newAcc)

		totalPower += val.Power()
	}

	lastHeight := uint32(21)
	for height := uint32(1); height < lastHeight; height++ {
		blk, cert := ts.GenerateTestBlock(height)
		mockStore.SaveBlock(blk, cert)
	}
	sbx := NewSandbox(mockStore.LastHeight,
		mockStore, param.FromGenesis(params), cmt, totalPower).(*sandbox)
	assert.Equal(t, lastHeight, sbx.CurrentHeight())
	assert.Equal(t, param.FromGenesis(params), sbx.Params())

	return &testData{
		TestSuite: ts,
		valKeys:   valKeys,
		store:     mockStore,
		sbx:       sbx,
	}
}

func TestAccountChange(t *testing.T) {
	td := setup(t)

	t.Run("Should returns nil for invalid address", func(t *testing.T) {
		invAddr := td.RandAccAddress()
		assert.Nil(t, td.sbx.Account(invAddr))

		td.sbx.IterateAccounts(func(_ crypto.Address, _ *account.Account, _ bool) {
			panic("should be empty")
		})
	})

	t.Run("Retrieve an account from store and update it", func(t *testing.T) {
		acc, addr := td.GenerateTestAccount()
		bal := acc.Balance()
		td.store.UpdateAccount(addr, acc)

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
		assert.Nil(t, td.sbx.Validator(invAddr))

		td.sbx.IterateValidators(func(_ *validator.Validator, _ bool, _ bool) {
			panic("should be empty")
		})
	})

	t.Run("Retrieve an validator from store and update it", func(t *testing.T) {
		val := td.GenerateTestValidator()
		addr := val.Address()
		stk := val.Stake()
		td.store.UpdateValidator(val)

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
		totalAccs := td.store.TotalAccounts()

		acc1 := td.sbx.MakeNewAccount(td.RandAccAddress())
		assert.Equal(t, totalAccs, acc1.Number())

		acc2 := td.sbx.MakeNewAccount(td.RandAccAddress())
		assert.Equal(t, totalAccs+1, acc2.Number())
	})
}

func TestTotalValidatorCounter(t *testing.T) {
	td := setup(t)

	t.Run("Should update total validator counter", func(t *testing.T) {
		totalVals := td.store.TotalValidators()

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
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("The code did not panic")
			}
		}()
		addr := crypto.TreasuryAddress
		td.sbx.MakeNewAccount(addr)
	})

	t.Run("Try creating duplicated validator, Should panic", func(t *testing.T) {
		assert.Panics(t, func() {
			pub := td.valKeys[3].PublicKey()
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
		acc := td.sbx.MakeNewAccount(addr)
		acc.AddToBalance(1)

		assert.NotEqual(t, acc, td.sbx.Account(addr))
	})

	t.Run("existing account", func(t *testing.T) {
		addr := crypto.TreasuryAddress
		acc := td.sbx.Account(addr)
		acc.AddToBalance(1)

		assert.NotEqual(t, acc, td.sbx.Account(addr))
	})

	t.Run("sandbox account", func(t *testing.T) {
		addr := crypto.TreasuryAddress
		acc := td.sbx.Account(addr)
		acc.AddToBalance(1)

		assert.NotEqual(t, acc, td.sbx.Account(addr))
	})
}

func TestValidatorDeepCopy(t *testing.T) {
	td := setup(t)

	t.Run("non existing validator", func(t *testing.T) {
		pub, _ := td.RandBLSKeyPair()
		val := td.sbx.MakeNewValidator(pub)
		val.AddToStake(1)

		assert.NotEqual(t, val, td.sbx.Validator(pub.ValidatorAddress()))
	})

	val0, _ := td.store.ValidatorByNumber(0)
	addr := val0.Address()
	t.Run("existing validator", func(t *testing.T) {
		val := td.sbx.Validator(addr)
		val.AddToStake(1)

		assert.NotEqual(t, val, td.sbx.Validator(addr))
	})

	t.Run("sandbox validator", func(t *testing.T) {
		val := td.sbx.Validator(addr)
		val.AddToStake(1)

		assert.NotEqual(t, val, td.sbx.Validator(addr))
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

func TestVerifyProof(t *testing.T) {
	td := setup(t)

	lastCert := td.store.LastCertificate()
	lastHeight := lastCert.Height()
	vals := td.sbx.committee.Validators()

	// Try to evaluate a valid sortition
	var validProof sortition.Proof
	var validLockTime uint32
	var validVal *validator.Validator
	for height := lastHeight; height > 0; height-- {
		block := td.store.Blocks[height]
		for index, valKey := range td.valKeys {
			ok, proof := sortition.EvaluateSortition(
				block.Header().SortitionSeed(), valKey.PrivateKey(),
				td.sbx.totalPower, vals[index].Power())

			if ok {
				validProof = proof
				validLockTime = height
				validVal = vals[index]
			}
		}
	}

	t.Run("invalid proof", func(t *testing.T) {
		invalidProof := td.RandProof()
		assert.False(t, td.sbx.VerifyProof(validLockTime, invalidProof, validVal))
	})
	t.Run("invalid height", func(t *testing.T) {
		assert.False(t, td.sbx.VerifyProof(td.RandHeight(), validProof, validVal))
	})

	t.Run("genesis block height", func(t *testing.T) {
		assert.False(t, td.sbx.VerifyProof(0, validProof, validVal))
	})

	t.Run("Ok", func(t *testing.T) {
		assert.True(t, td.sbx.VerifyProof(validLockTime, validProof, validVal))
	})
}

func TestIsBanned(t *testing.T) {
	td := setup(t)

	t.Run("Validator is in the banned list", func(t *testing.T) {
		pub, prv := td.RandBLSKeyPair()
		td.store.TestConfig.BannedAddrs[pub.ValidatorAddress()] = true
		trx := td.GenerateTestSortitionTx(
			testsuite.TransactionWithBLSSigner(prv))

		assert.True(t, td.sbx.IsBanned(trx))
	})

	t.Run("Validator is not in the banned list", func(t *testing.T) {
		trx := td.GenerateTestSortitionTx()

		assert.False(t, td.sbx.IsBanned(trx))
	})

	t.Run("Xeggex account is in a frozen state", func(t *testing.T) {
		t.Run("Attempt to transfer assets to another account should be rejected", func(t *testing.T) {
			pub, prv := td.RandBLSKeyPair()
			td.store.XeggexAccount().DepositAddrs = pub.AccountAddress()
			trx := td.GenerateTestTransferTx(testsuite.TransactionWithBLSSigner(prv))

			assert.True(t, td.sbx.IsBanned(trx))
		})

		t.Run("Attempt to transfer assets to the Watcher account but not the full balance", func(t *testing.T) {
			pub, prv := td.RandBLSKeyPair()
			td.store.XeggexAccount().DepositAddrs = pub.AccountAddress()
			trx := td.GenerateTestTransferTx(
				testsuite.TransactionWithBLSSigner(prv),
				testsuite.TransactionWithReceiver(td.store.XeggexAccount().WatcherAddrs),
			)

			assert.True(t, td.sbx.IsBanned(trx))
		})

		t.Run("Attempt to transfer assets to the Watcher account with the full balance", func(t *testing.T) {
			pub, prv := td.RandBLSKeyPair()
			td.store.XeggexAccount().DepositAddrs = pub.AccountAddress()
			trx := td.GenerateTestTransferTx(
				testsuite.TransactionWithBLSSigner(prv),
				testsuite.TransactionWithReceiver(td.store.XeggexAccount().WatcherAddrs),
				testsuite.TransactionWithAmount(amount.Amount(500_000e9)),
			)

			assert.False(t, td.sbx.IsBanned(trx))
		})
	})

	t.Run("Xeggex account is in an unfrozen state", func(t *testing.T) {
		watcherPub, watcherPrv := td.RandBLSKeyPair()
		xeggexPub, xeggexPrv := td.RandBLSKeyPair()

		trx1 := td.GenerateTestTransferTx(testsuite.TransactionWithBLSSigner(watcherPrv))
		blk, cert := td.GenerateTestBlock(td.RandHeight(), testsuite.BlockWithTransactions([]*tx.Tx{trx1}))
		td.store.SaveBlock(blk, cert)

		td.store.TestConfig.XeggexAccount.DepositAddrs = xeggexPub.AccountAddress()
		td.store.TestConfig.XeggexAccount.WatcherAddrs = watcherPub.AccountAddress()

		trx2 := td.GenerateTestTransferTx(testsuite.TransactionWithBLSSigner(xeggexPrv))
		assert.False(t, td.sbx.IsBanned(trx2))
	})
}
