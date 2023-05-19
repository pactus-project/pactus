package state

import (
	"fmt"
	"testing"
	"time"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/genesis"
	"github.com/pactus-project/pactus/store"
	"github.com/pactus-project/pactus/txpool"
	"github.com/pactus-project/pactus/types/account"
	"github.com/pactus-project/pactus/types/block"
	"github.com/pactus-project/pactus/types/param"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/tx/payload"
	"github.com/pactus-project/pactus/types/validator"
	"github.com/pactus-project/pactus/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var tState1 *state
var tState2 *state
var tState3 *state
var tState4 *state
var tValSigner1 crypto.Signer
var tValSigner2 crypto.Signer
var tValSigner3 crypto.Signer
var tValSigner4 crypto.Signer
var tCommonTxPool *txpool.MockTxPool

func setup(t *testing.T) {
	pub1, prv1 := bls.GenerateTestKeyPair()
	pub2, prv2 := bls.GenerateTestKeyPair()
	pub3, prv3 := bls.GenerateTestKeyPair()
	pub4, prv4 := bls.GenerateTestKeyPair()

	tValSigner1 = crypto.NewSigner(prv1)
	tValSigner2 = crypto.NewSigner(prv2)
	tValSigner3 = crypto.NewSigner(prv3)
	tValSigner4 = crypto.NewSigner(prv4)

	genTime := util.RoundNow(10)
	tCommonTxPool = txpool.MockingTxPool()

	store1 := store.MockingStore()
	store2 := store.MockingStore()
	store3 := store.MockingStore()
	store4 := store.MockingStore()

	acc := account.NewAccount(0)
	acc.AddToBalance(21 * 1e14) // 2,100,000,000,000,000
	val1 := validator.NewValidator(pub1, 0)
	val2 := validator.NewValidator(pub2, 1)
	val3 := validator.NewValidator(pub3, 2)
	val4 := validator.NewValidator(pub4, 3)
	params := param.DefaultParams()
	params.CommitteeSize = 5
	params.BondInterval = 10

	accs := map[crypto.Address]*account.Account{crypto.TreasuryAddress: acc}
	vals := []*validator.Validator{val1, val2, val3, val4}
	gnDoc := genesis.MakeGenesis(genTime, accs, vals, params)

	st1, err := LoadOrNewState(gnDoc, []crypto.Signer{tValSigner1}, store1, tCommonTxPool, nil)
	require.NoError(t, err)
	st2, err := LoadOrNewState(gnDoc, []crypto.Signer{tValSigner2}, store2, tCommonTxPool, nil)
	require.NoError(t, err)
	st3, err := LoadOrNewState(gnDoc, []crypto.Signer{tValSigner3}, store3, tCommonTxPool, nil)
	require.NoError(t, err)
	st4, err := LoadOrNewState(gnDoc, []crypto.Signer{tValSigner4}, store4, tCommonTxPool, nil)
	require.NoError(t, err)

	tState1, _ = st1.(*state)
	tState2, _ = st2.(*state)
	tState3, _ = st3.(*state)
	tState4, _ = st4.(*state)
}

func makeBlockAndCertificate(t *testing.T, round int16,
	signers ...crypto.Signer) (*block.Block, *block.Certificate) {
	var st *state
	if tState1.committee.IsProposer(tState1.signers[0].Address(), round) {
		st = tState1
	} else if tState1.committee.IsProposer(tState2.signers[0].Address(), round) {
		st = tState2
	} else if tState1.committee.IsProposer(tState3.signers[0].Address(), round) {
		st = tState3
	} else {
		st = tState4
	}

	rewardAddr := st.signers[0].Address()
	b, err := st.ProposeBlock(st.signers[0], rewardAddr, round)
	require.NoError(t, err)
	c := makeCertificateAndSign(t, b.Hash(), round, signers...)

	return b, c
}

func makeCertificateAndSign(t *testing.T, blockHash hash.Hash, round int16,
	signers ...crypto.Signer) *block.Certificate {
	assert.NotZero(t, len(signers))
	sigs := make([]*bls.Signature, len(signers))
	sb := block.CertificateSignBytes(blockHash, round)
	committers := []int32{0, 1, 2, 3}
	signedBy := []int32{}

	for i, s := range signers {
		if s.Address().EqualsTo(tValSigner1.Address()) {
			signedBy = append(signedBy, 0)
		}

		if s.Address().EqualsTo(tValSigner2.Address()) {
			signedBy = append(signedBy, 1)
		}

		if s.Address().EqualsTo(tValSigner3.Address()) {
			signedBy = append(signedBy, 2)
		}

		if s.Address().EqualsTo(tValSigner4.Address()) {
			signedBy = append(signedBy, 3)
		}
		sigs[i] = s.SignData(sb).(*bls.Signature)
	}

	absentees := util.Subtracts(committers, signedBy)
	return block.NewCertificate(round, committers, absentees, bls.Aggregate(sigs))
}

func CommitBlockForAllStates(t *testing.T, b *block.Block, c *block.Certificate) {
	assert.NoError(t, tState1.CommitBlock(tState1.lastInfo.BlockHeight()+1, b, c))
	assert.NoError(t, tState2.CommitBlock(tState2.lastInfo.BlockHeight()+1, b, c))
	assert.NoError(t, tState3.CommitBlock(tState3.lastInfo.BlockHeight()+1, b, c))
	assert.NoError(t, tState4.CommitBlock(tState4.lastInfo.BlockHeight()+1, b, c))
}

func moveToNextHeightForAllStates(t *testing.T) {
	b, c := makeBlockAndCertificate(t, 0, tValSigner1, tValSigner2, tValSigner3, tValSigner4)
	CommitBlockForAllStates(t, b, c)
}

func TestProposeBlockAndValidation(t *testing.T) {
	setup(t)
	moveToNextHeightForAllStates(t)

	b1, err := tState1.ProposeBlock(tState1.signers[0], crypto.GenerateTestAddress(), 0)
	assert.Error(t, err, "Should not propose")
	assert.Nil(t, b1)

	trx := tx.NewTransferTx(tState2.lastInfo.BlockHash().Stamp(), 1, tValSigner1.Address(),
		tValSigner2.Address(), 1000, 1000, "")
	tValSigner1.SignMsg(trx)
	assert.NoError(t, tCommonTxPool.AppendTx(trx))

	b2, err := tState2.ProposeBlock(tState2.signers[0], crypto.GenerateTestAddress(), 0)
	assert.NoError(t, err)
	assert.NotNil(t, b2)
	assert.Equal(t, b2.Transactions().Len(), 2)
	require.NoError(t, tState1.ValidateBlock(b2))

	// Propose and validate again
	b3, err := tState2.ProposeBlock(tState2.signers[0], crypto.GenerateTestAddress(), 0)
	assert.NoError(t, err)
	assert.NotNil(t, b3)
	assert.Equal(t, b3.Transactions().Len(), 2)
	require.NoError(t, tState1.ValidateBlock(b3))
}

func TestBlockSubsidyTx(t *testing.T) {
	setup(t)

	// Without reward address in config
	rewardAddr := crypto.GenerateTestAddress()
	trx := tState1.createSubsidyTx(rewardAddr, 7)
	assert.True(t, trx.IsSubsidyTx())
	assert.Equal(t, trx.Payload().Value(), tState1.params.BlockReward+7)
	assert.Equal(t, trx.Payload().(*payload.TransferPayload).Sender, crypto.TreasuryAddress)
	assert.Equal(t, trx.Payload().(*payload.TransferPayload).Receiver, rewardAddr)
}

func TestCommitBlocks(t *testing.T) {
	setup(t)

	b1, c1 := makeBlockAndCertificate(t, 1, tValSigner1, tValSigner2, tValSigner3)
	invBlock := block.GenerateTestBlock(nil, nil)
	assert.Error(t, tState1.CommitBlock(1, invBlock, c1))
	// No error here but block is ignored, because the height is invalid
	assert.NoError(t, tState1.CommitBlock(2, b1, c1))
	assert.NoError(t, tState1.CommitBlock(1, b1, c1))

	assert.Equal(t, tState1.LastBlockHash(), b1.Hash())
	assert.Equal(t, tState1.LastBlockTime(), b1.Header().Time())
	assert.Equal(t, tState1.LastCertificate().Hash(), c1.Hash())
	assert.Equal(t, tState1.LastBlockHeight(), uint32(1))
	assert.Equal(t, tState1.GenesisHash(), tState2.GenesisHash())
}

func TestCommitSandbox(t *testing.T) {
	t.Run("Add new account", func(t *testing.T) {
		setup(t)

		addr := crypto.GenerateTestAddress()
		sb := tState1.concreteSandbox()
		newAcc := sb.MakeNewAccount(addr)
		newAcc.AddToBalance(1)
		tState1.commitSandbox(sb, 0)

		assert.NotNil(t, tState1.AccountByAddress(addr))
	})

	t.Run("Add new validator", func(t *testing.T) {
		setup(t)

		pub, _ := bls.GenerateTestKeyPair()
		sb := tState1.concreteSandbox()
		newVal := sb.MakeNewValidator(pub)
		newVal.AddToStake(123)
		sb.UpdateValidator(newVal)
		tState1.commitSandbox(sb, 0)

		assert.True(t, tState1.store.HasValidator(pub.Address()))
	})

	t.Run("Modify account", func(t *testing.T) {
		setup(t)

		sb := tState1.concreteSandbox()
		acc := sb.Account(crypto.TreasuryAddress)
		acc.SubtractFromBalance(1)
		sb.UpdateAccount(crypto.TreasuryAddress, acc)
		tState1.commitSandbox(sb, 0)

		acc1 := tState1.AccountByAddress(crypto.TreasuryAddress)
		assert.Equal(t, acc1.Balance(), acc.Balance())
	})

	t.Run("Modify validator", func(t *testing.T) {
		setup(t)

		sb := tState1.concreteSandbox()
		val := sb.Validator(tValSigner2.Address())
		val.AddToStake(2002)
		sb.UpdateValidator(val)
		tState1.commitSandbox(sb, 0)

		val1, _ := tState1.store.Validator(tValSigner2.Address())
		assert.Equal(t, val1.Stake(), val.Stake())
	})

	t.Run("Move committee", func(t *testing.T) {
		setup(t)

		nextProposer := tState1.committee.Proposer(1)

		sb := tState1.concreteSandbox()
		tState1.commitSandbox(sb, 0)

		assert.Equal(t, tState1.committee.Proposer(0).Address(), nextProposer.Address())
	})

	t.Run("Move committee next round", func(t *testing.T) {
		setup(t)

		nextNextProposer := tState1.committee.Proposer(2)

		sb := tState1.concreteSandbox()
		tState1.commitSandbox(sb, 1)

		assert.Equal(t, tState1.committee.Proposer(0).Address(), nextNextProposer.Address())
	})
}

func TestUpdateLastCertificate(t *testing.T) {
	setup(t)
	b1, c1 := makeBlockAndCertificate(t, 0, tValSigner1, tValSigner3, tValSigner4)
	b11, c11 := makeBlockAndCertificate(t, 0, tValSigner1, tValSigner2, tValSigner3, tValSigner4)
	_, c12 := makeBlockAndCertificate(t, 1, tValSigner1, tValSigner2, tValSigner3, tValSigner4)

	CommitBlockForAllStates(t, b1, c1)

	assert.Equal(t, b1.Hash(), b11.Hash())
	assert.Equal(t, tState1.lastInfo.Certificate().Hash(), c1.Hash())
	assert.Error(t, tState1.UpdateLastCertificate(c12))
	assert.NoError(t, tState1.UpdateLastCertificate(c1))
	assert.Equal(t, tState1.lastInfo.Certificate().Hash(), c1.Hash())
	assert.NoError(t, tState1.UpdateLastCertificate(c11))
	assert.Equal(t, tState1.lastInfo.Certificate().Hash(), c11.Hash())
}

func TestInvalidProposerProposeBlock(t *testing.T) {
	setup(t)

	_, err := tState2.ProposeBlock(tState2.signers[0], crypto.GenerateTestAddress(), 0)
	assert.Error(t, err, "Should not propose")
	_, err = tState2.ProposeBlock(tState2.signers[0], crypto.GenerateTestAddress(), 1)
	assert.NoError(t, err, "Should propose")
}

func TestBlockProposal(t *testing.T) {
	setup(t)
	moveToNextHeightForAllStates(t)

	t.Run("validity of proposed block", func(t *testing.T) {
		b, err := tState2.ProposeBlock(tState2.signers[0], crypto.GenerateTestAddress(), 0)
		assert.NoError(t, err)
		assert.NoError(t, tState1.ValidateBlock(b))
	})

	t.Run("Tx pool has two subsidy transactions", func(t *testing.T) {
		trx := tState3.createSubsidyTx(crypto.GenerateTestAddress(), 0)
		assert.NoError(t, tState3.txPool.AppendTx(trx))

		// Moving to the next round
		b, err := tState3.ProposeBlock(tState3.signers[0], crypto.GenerateTestAddress(), 1)
		assert.NoError(t, err)
		assert.NoError(t, tState1.ValidateBlock(b))
		assert.Equal(t, b.Transactions().Len(), 1)
	})
}

func TestInvalidBlock(t *testing.T) {
	setup(t)

	b := block.GenerateTestBlock(nil, nil)
	assert.Error(t, tState1.ValidateBlock(b))
}

func TestForkDetection(t *testing.T) {
	setup(t)

	moveToNextHeightForAllStates(t)

	b5m, c5m := makeBlockAndCertificate(t, 0, tValSigner1, tValSigner2, tValSigner3)
	b5f, c5f := makeBlockAndCertificate(t, 1, tValSigner1, tValSigner2, tValSigner3)
	assert.NoError(t, tState1.CommitBlock(2, b5m, c5m))
	assert.NoError(t, tState2.CommitBlock(2, b5m, c5m))
	assert.NoError(t, tState3.CommitBlock(2, b5m, c5m))
	assert.NoError(t, tState4.CommitBlock(2, b5f, c5f))

	b6, c6 := makeBlockAndCertificate(t, 0, tValSigner1, tValSigner2, tValSigner3)

	assert.NoError(t, tState1.CommitBlock(3, b6, c6))
	assert.NoError(t, tState2.CommitBlock(3, b6, c6))
	assert.NoError(t, tState3.CommitBlock(3, b6, c6))
	t.Run("Fork is detected, Should panic ", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("The code did not panic")
			}
		}()
		assert.Error(t, tState4.CommitBlock(3, b6, c6))
	})
}

func TestSortition(t *testing.T) {
	setup(t)

	pub, prv := bls.GenerateTestKeyPair()
	signer := crypto.NewSigner(prv)
	store := store.MockingStore()
	St1, _ := LoadOrNewState(tState1.genDoc, []crypto.Signer{signer}, store, tCommonTxPool, nil)
	stNew := St1.(*state)

	assert.False(t, stNew.evaluateSortition()) //  not a validator
	assert.Equal(t, tState1.CommitteePower(), int64(4))

	height := uint32(1)
	for ; height <= 11; height++ {
		if height == 2 {
			trx := tx.NewBondTx(tState1.lastInfo.BlockHash().Stamp(), 1, tValSigner1.Address(),
				pub.Address(), pub, 10000000, 1000, "")
			tValSigner1.SignMsg(trx)

			assert.NoError(t, tCommonTxPool.AppendTx(trx))
		}

		b, c := makeBlockAndCertificate(t, 0, tValSigner1, tValSigner2, tValSigner3, tValSigner4)
		CommitBlockForAllStates(t, b, c)
		require.NoError(t, stNew.CommitBlock(height, b, c))
	}

	assert.False(t, stNew.evaluateSortition()) //  bonding period

	// Certificate next block
	b, c := makeBlockAndCertificate(t, 0, tValSigner1, tValSigner2, tValSigner3, tValSigner4)
	CommitBlockForAllStates(t, b, c)
	require.NoError(t, stNew.CommitBlock(height, b, c))
	height++

	assert.True(t, stNew.evaluateSortition())                  //  ok
	assert.False(t, tState1.committee.Contains(pub.Address())) // still not in the committee

	// ---------------------------------------------
	// Certificate next block, new validator should be in the committee now
	b, c = makeBlockAndCertificate(t, 0, tValSigner1, tValSigner2, tValSigner3, tValSigner4)
	CommitBlockForAllStates(t, b, c)
	require.NoError(t, stNew.CommitBlock(height, b, c))

	assert.True(t, stNew.evaluateSortition()) // in the committee
	assert.True(t, tState1.committee.Contains(tValSigner1.Address()))
	assert.True(t, tState1.committee.Contains(pub.Address()))

	// ---------------------------------------------
	// Let's save and load tState1
	tState1.Close()
	St1, _ = LoadOrNewState(tState1.genDoc, []crypto.Signer{tValSigner1}, store, tCommonTxPool, nil)
	st1 := St1.(*state)

	// ---------------------------------------------
	// Let's commit another block with the new committee
	b14, err := stNew.ProposeBlock(stNew.signers[0], crypto.GenerateTestAddress(), 3)
	require.NoError(t, err)
	require.NotNil(t, b14)

	sigs := make([]*bls.Signature, 4)
	sb := block.CertificateSignBytes(b14.Hash(), 3)

	sigs[0] = tValSigner2.SignData(sb).(*bls.Signature)
	sigs[1] = tValSigner3.SignData(sb).(*bls.Signature)
	sigs[2] = tValSigner4.SignData(sb).(*bls.Signature)
	sigs[3] = signer.SignData(sb).(*bls.Signature)
	c14 := block.NewCertificate(3, []int32{4, 0, 1, 2, 3}, []int32{0}, bls.Aggregate(sigs))

	height++
	assert.NoError(t, st1.CommitBlock(height, b14, c14))
	assert.NoError(t, tState1.CommitBlock(height, b14, c14))
	assert.NoError(t, tState2.CommitBlock(height, b14, c14))
	assert.NoError(t, tState3.CommitBlock(height, b14, c14))
	assert.NoError(t, tState4.CommitBlock(height, b14, c14))

	assert.Equal(t, tState1.CommitteePower(), int64(10000004))
	assert.Equal(t, tState1.TotalValidators(), int32(5))
}

func TestValidateBlockTime(t *testing.T) {
	setup(t)
	fmt.Printf("BlockTimeInSecond: %d\n", tState1.params.BlockTimeInSecond)

	// Time not rounded
	roundedNow := util.RoundNow(10)
	assert.Error(t, tState1.validateBlockTime(roundedNow.Add(-15*time.Second)))
	assert.Error(t, tState1.validateBlockTime(roundedNow.Add(-5*time.Second)))
	assert.Error(t, tState1.validateBlockTime(roundedNow.Add(5*time.Second)))
	assert.Error(t, tState1.validateBlockTime(roundedNow.Add(15*time.Second)))

	t.Run("Last block is committed 10 seconds ago", func(t *testing.T) {
		tState1.lastInfo.SetBlockTime(roundedNow.Add(-10 * time.Second))

		// Before or same as the last block time
		assert.Error(t, tState1.validateBlockTime(roundedNow.Add(-20*time.Second)))
		assert.Error(t, tState1.validateBlockTime(roundedNow.Add(-10*time.Second)))

		// Ok
		assert.NoError(t, tState1.validateBlockTime(roundedNow))
		assert.NoError(t, tState1.validateBlockTime(roundedNow.Add(10*time.Second)))
		assert.Equal(t, tState1.proposeNextBlockTime(), roundedNow, "Invalid proposed time for the next block")

		// More than threshold
		assert.Error(t, tState1.validateBlockTime(roundedNow.Add(20*time.Second)))
	})

	t.Run("Last block is committed one minute ago", func(t *testing.T) {
		tState1.lastInfo.SetBlockTime(roundedNow.Add(-1 * time.Minute)) // One minute ago

		// Before or same as the last block time
		assert.Error(t, tState1.validateBlockTime(tState1.lastInfo.BlockTime().Add(-10*time.Second)))
		assert.Error(t, tState1.validateBlockTime(tState1.lastInfo.BlockTime()))

		// Ok
		assert.NoError(t, tState1.validateBlockTime(roundedNow.Add(-10*time.Second)))
		assert.NoError(t, tState1.validateBlockTime(roundedNow))
		assert.NoError(t, tState1.validateBlockTime(roundedNow.Add(10*time.Second)))
		assert.Equal(t, tState1.proposeNextBlockTime(), roundedNow, "Invalid proposed time for the next block")

		// More than threshold
		assert.Error(t, tState1.validateBlockTime(roundedNow.Add(20*time.Second)))
	})

	t.Run("Last block is committed in future", func(t *testing.T) {
		tState1.lastInfo.SetBlockTime(roundedNow.Add(1 * time.Minute)) // One minute later

		assert.Error(t, tState1.validateBlockTime(tState1.lastInfo.BlockTime().Add(+1*time.Minute)))

		// Before the last block time
		assert.Error(t, tState1.validateBlockTime(tState1.lastInfo.BlockTime().Add(-10*time.Second)))
		assert.Error(t, tState1.validateBlockTime(tState1.lastInfo.BlockTime()))

		// Ok
		assert.NoError(t, tState1.validateBlockTime(tState1.lastInfo.BlockTime().Add(10*time.Second)))
		assert.NoError(t, tState1.validateBlockTime(tState1.lastInfo.BlockTime().Add(20*time.Second)))

		// More than threshold
		assert.Error(t, tState1.validateBlockTime(tState1.lastInfo.BlockTime().Add(30*time.Second)))
	})
}

func TestInvalidBlockVersion(t *testing.T) {
	setup(t)

	tState1.params.BlockVersion = 2
	b, _ := tState1.ProposeBlock(tState1.signers[0], crypto.GenerateTestAddress(), 0)
	assert.Error(t, tState2.ValidateBlock(b))
}

func TestInvalidBlockTime(t *testing.T) {
	setup(t)
	moveToNextHeightForAllStates(t)

	validBlock, _ := makeBlockAndCertificate(t, 0, tValSigner1, tValSigner2, tValSigner3, tValSigner4)
	invalidBlock := block.MakeBlock(
		validBlock.Header().Version(),
		validBlock.Header().Time().Add(30*time.Second),
		validBlock.Transactions(),
		validBlock.Header().PrevBlockHash(),
		validBlock.Header().StateRoot(),
		validBlock.PrevCertificate(),
		validBlock.Header().SortitionSeed(),
		validBlock.Header().ProposerAddress())

	assert.NoError(t, tState1.ValidateBlock(validBlock))
	assert.Error(t, tState1.ValidateBlock(invalidBlock))
}

func TestValidatorHelpers(t *testing.T) {
	setup(t)

	t.Run("Should return nil for non-existing Validator Address", func(t *testing.T) {
		_, prv5 := bls.GenerateTestKeyPair()
		nonExistenceValidator := tState1.ValidatorByAddress(prv5.PublicKey().Address())
		assert.Nil(t, nonExistenceValidator, "State 1 returned non-nil For non-existing validator")
		nonExistenceValidator = tState2.ValidatorByAddress(prv5.PublicKey().Address())
		assert.Nil(t, nonExistenceValidator, "State 2 returned non-nil For non-existing validator")
		nonExistenceValidator = tState3.ValidatorByAddress(prv5.PublicKey().Address())
		assert.Nil(t, nonExistenceValidator, "State 3 returned non-nil For non-existing validator")
		nonExistenceValidator = tState4.ValidatorByAddress(prv5.PublicKey().Address())
		assert.Nil(t, nonExistenceValidator, "State 4 returned non-nil For non-existing validator")
	})

	t.Run("Should return validator for valid committee Validator Address", func(t *testing.T) {
		existingValidator := tState4.ValidatorByAddress(tValSigner1.Address())
		assert.NotNil(t, existingValidator)
		assert.Zero(t, existingValidator.Number())
	})

	t.Run("Should return validator for corresponding Validator number", func(t *testing.T) {
		existingValidator := tState4.ValidatorByNumber(1)
		assert.NotNil(t, existingValidator)
		assert.Equal(t, tValSigner2.Address(), existingValidator.Address())
	})

	t.Run("Should return nil for invalid Validator number", func(t *testing.T) {
		nonExistenceValidator := tState4.ValidatorByNumber(10)
		assert.Nil(t, nonExistenceValidator)
	})
}

func TestLoadState(t *testing.T) {
	setup(t)

	// Add a bond transactions to change total power (stake)
	pub, _ := bls.GenerateTestKeyPair()
	tx2 := tx.NewBondTx(tState1.LastBlockHash().Stamp(), 1, tValSigner1.Address(), pub.Address(), pub, 8888000, 8888, "")
	tValSigner1.SignMsg((tx2))

	assert.NoError(t, tCommonTxPool.AppendTx(tx2))

	for i := 0; i < 4; i++ {
		moveToNextHeightForAllStates(t)
	}
	b5, c5 := makeBlockAndCertificate(t, 1, tValSigner1, tValSigner2, tValSigner3, tValSigner4)
	CommitBlockForAllStates(t, b5, c5)

	b6, c6 := makeBlockAndCertificate(t, 0, tValSigner1, tValSigner2, tValSigner3, tValSigner4)

	// Load last state info
	st1Load, err := LoadOrNewState(tState1.genDoc, []crypto.Signer{tValSigner1}, tState1.store, tCommonTxPool, nil)
	require.NoError(t, err)

	assert.Equal(t, tState1.store.TotalAccounts(), st1Load.(*state).store.TotalAccounts())
	assert.Equal(t, tState1.store.TotalValidators(), st1Load.(*state).store.TotalValidators())
	assert.Equal(t, tState1.committee.Committers(), st1Load.(*state).committee.Committers())
	assert.Equal(t, tState1.committee.TotalPower(), st1Load.(*state).committee.TotalPower())
	assert.Equal(t, tState1.TotalPower(), st1Load.(*state).totalPower())
	assert.Equal(t, tState1.store.TotalAccounts(), int32(5))

	require.NoError(t, st1Load.CommitBlock(6, b6, c6))
	require.NoError(t, tState2.CommitBlock(6, b6, c6))
}

func TestLoadStateAfterChangingGenesis(t *testing.T) {
	setup(t)

	// Let's commit some blocks
	i := 0
	for ; i < 10; i++ {
		moveToNextHeightForAllStates(t)
	}

	_, err := LoadOrNewState(tState1.genDoc, []crypto.Signer{tValSigner1}, tState1.store, txpool.MockingTxPool(), nil)
	require.NoError(t, err)

	pub, _ := bls.GenerateTestKeyPair()
	val := validator.NewValidator(pub, 4)
	vals := append(tState1.genDoc.Validators(), val)

	genDoc := genesis.MakeGenesis(
		tState1.genDoc.GenesisTime(),
		tState1.genDoc.Accounts(),
		vals,
		tState1.genDoc.Params())

	// Load last state info after modifying genesis
	_, err = LoadOrNewState(genDoc, []crypto.Signer{tValSigner1}, tState1.store, txpool.MockingTxPool(), nil)
	require.Error(t, err)
}

func TestSetBlockTime(t *testing.T) {
	setup(t)

	assert.Equal(t, tState1.BlockTime(), 10*time.Second)

	t.Run("Last block time is a bit far in past", func(t *testing.T) {
		tState1.lastInfo.SetBlockTime(util.RoundNow(10).Add(-20 * time.Second))
		b, _ := tState1.ProposeBlock(tState1.signers[0], crypto.GenerateTestAddress(), 0)
		fmt.Printf("last block time: %s\nproposed time  : %s\n", tState1.lastInfo.BlockTime(), b.Header().Time().UTC())
		assert.True(t, b.Header().Time().After(tState1.lastInfo.BlockTime()))
		assert.True(t, b.Header().Time().Before(util.Now().Add(10*time.Second)))
		assert.Zero(t, b.Header().Time().Second()%10)
	})

	t.Run("Last block time is almost good", func(t *testing.T) {
		tState1.lastInfo.SetBlockTime(util.RoundNow(10).Add(-10 * time.Second))
		b, _ := tState1.ProposeBlock(tState1.signers[0], crypto.GenerateTestAddress(), 0)
		fmt.Printf("last block time: %s\nproposed time  : %s\n", tState1.lastInfo.BlockTime(), b.Header().Time().UTC())
		assert.True(t, b.Header().Time().After(tState1.lastInfo.BlockTime()))
		assert.True(t, b.Header().Time().Before(util.Now().Add(10*time.Second)))
		assert.Zero(t, b.Header().Time().Second()%10)
	})

	// After our time
	t.Run("Last block time is in near future", func(t *testing.T) {
		tState1.lastInfo.SetBlockTime(util.RoundNow(10).Add(+10 * time.Second))
		b, _ := tState1.ProposeBlock(tState1.signers[0], crypto.GenerateTestAddress(), 0)
		fmt.Printf("last block time: %s\nproposed time  : %s\n", tState1.lastInfo.BlockTime(), b.Header().Time().UTC())
		assert.True(t, b.Header().Time().After(tState1.lastInfo.BlockTime()))
		assert.Zero(t, b.Header().Time().Second()%10)
	})

	t.Run("Last block time is more than a block in future", func(t *testing.T) {
		tState1.lastInfo.SetBlockTime(util.RoundNow(10).Add(+20 * time.Second))
		b, _ := tState1.ProposeBlock(tState1.signers[0], crypto.GenerateTestAddress(), 0)
		fmt.Printf("last block time: %s\nproposed time  : %s\n", tState1.lastInfo.BlockTime(), b.Header().Time().UTC())
		assert.True(t, b.Header().Time().After(tState1.lastInfo.BlockTime()))
		assert.Zero(t, b.Header().Time().Second()%10)
	})
}

func TestIsValidator(t *testing.T) {
	setup(t)

	assert.True(t, tState1.IsInCommittee(tValSigner1.Address()))
	assert.True(t, tState1.IsProposer(tValSigner1.Address(), 0))
	assert.True(t, tState1.IsProposer(tValSigner2.Address(), 1))
	assert.True(t, tState1.IsInCommittee(tValSigner2.Address()))
	assert.True(t, tState1.IsValidator(tValSigner2.Address()))

	addr := crypto.GenerateTestAddress()
	assert.False(t, tState1.IsInCommittee(addr))
	assert.False(t, tState1.IsProposer(addr, 0))
	assert.False(t, tState1.IsInCommittee(addr))
	assert.False(t, tState1.IsValidator(addr))
}

func TestCalculatingGenesisState(t *testing.T) {
	setup(t)

	r := tState1.calculateGenesisStateRootFromGenesisDoc()
	assert.Equal(t, tState1.stateRoot(), r)
}

func TestCommittingInvalidBlock(t *testing.T) {
	setup(t)

	moveToNextHeightForAllStates(t)

	txs := block.NewTxs()
	trx := tState2.createSubsidyTx(crypto.GenerateTestAddress(), 0)
	txs.Append(trx)
	b := block.MakeBlock(2, util.Now(), txs, tState2.lastInfo.BlockHash(), tState2.stateRoot(),
		tState2.lastInfo.Certificate(), tState2.lastInfo.SortitionSeed(), tState2.signers[0].Address())
	c := makeCertificateAndSign(t, b.Hash(), 0, tValSigner1, tValSigner2, tValSigner3, tValSigner4)

	// tState1 receives a block with version 2 and rejects it.
	// It is possible that the same block would be considered valid by tState2.
	assert.Error(t, tState1.CommitBlock(2, b, c))
}
