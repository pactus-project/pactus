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
	"github.com/pactus-project/pactus/types/certificate"
	"github.com/pactus-project/pactus/types/param"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/tx/payload"
	"github.com/pactus-project/pactus/types/validator"
	"github.com/pactus-project/pactus/types/vote"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/errors"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TODO: clean me

type testData struct {
	*testsuite.TestSuite

	state1       *state
	state2       *state
	state3       *state
	state4       *state
	valKey1      *bls.ValidatorKey
	valKey2      *bls.ValidatorKey
	valKey3      *bls.ValidatorKey
	valKey4      *bls.ValidatorKey
	commonTxPool *txpool.MockTxPool
}

func setup(t *testing.T) *testData {
	t.Helper()

	ts := testsuite.NewTestSuite(t)

	pub1, prv1 := ts.RandBLSKeyPair()
	pub2, prv2 := ts.RandBLSKeyPair()
	pub3, prv3 := ts.RandBLSKeyPair()
	pub4, prv4 := ts.RandBLSKeyPair()

	valKey1 := bls.NewValidatorKey(prv1)
	valKey2 := bls.NewValidatorKey(prv2)
	valKey3 := bls.NewValidatorKey(prv3)
	valKey4 := bls.NewValidatorKey(prv4)

	genTime := util.RoundNow(10)
	commonTxPool := txpool.MockingTxPool()

	store1 := store.MockingStore(ts)
	store2 := store.MockingStore(ts)
	store3 := store.MockingStore(ts)
	store4 := store.MockingStore(ts)

	val1 := validator.NewValidator(pub1, 0)
	val2 := validator.NewValidator(pub2, 1)
	val3 := validator.NewValidator(pub3, 2)
	val4 := validator.NewValidator(pub4, 3)
	params := param.DefaultParams()
	params.CommitteeSize = 5
	params.BondInterval = 10

	acc1 := account.NewAccount(0)
	acc1.AddToBalance(21 * 1e15) // 21,000,000,000,000,000
	acc2 := account.NewAccount(1)
	acc2.AddToBalance(21 * 1e15) // 21,000,000,000,000,000

	accs := map[crypto.Address]*account.Account{
		crypto.TreasuryAddress: acc1,
		ts.RandAccAddress():    acc2,
	}
	vals := []*validator.Validator{val1, val2, val3, val4}
	gnDoc := genesis.MakeGenesis(genTime, accs, vals, params)

	st1, err := LoadOrNewState(gnDoc, []*bls.ValidatorKey{valKey1}, store1, commonTxPool, nil)
	require.NoError(t, err)
	st2, err := LoadOrNewState(gnDoc, []*bls.ValidatorKey{valKey2}, store2, commonTxPool, nil)
	require.NoError(t, err)
	st3, err := LoadOrNewState(gnDoc, []*bls.ValidatorKey{valKey3}, store3, commonTxPool, nil)
	require.NoError(t, err)
	st4, err := LoadOrNewState(gnDoc, []*bls.ValidatorKey{valKey4}, store4, commonTxPool, nil)
	require.NoError(t, err)

	state1, _ := st1.(*state)
	state2, _ := st2.(*state)
	state3, _ := st3.(*state)
	state4, _ := st4.(*state)

	return &testData{
		TestSuite:    ts,
		state1:       state1,
		state2:       state2,
		state3:       state3,
		state4:       state4,
		valKey1:      valKey1,
		valKey2:      valKey2,
		valKey3:      valKey3,
		valKey4:      valKey4,
		commonTxPool: commonTxPool,
	}
}

func (td *testData) makeBlockAndCertificate(t *testing.T, round int16,
	valKeys ...*bls.ValidatorKey,
) (*block.Block, *certificate.Certificate) {
	t.Helper()

	var st *state
	if td.state1.committee.IsProposer(td.state1.valKeys[0].Address(), round) {
		st = td.state1
	} else if td.state1.committee.IsProposer(td.state2.valKeys[0].Address(), round) {
		st = td.state2
	} else if td.state1.committee.IsProposer(td.state3.valKeys[0].Address(), round) {
		st = td.state3
	} else {
		st = td.state4
	}

	rewardAddr := st.valKeys[0].Address()
	b, err := st.ProposeBlock(st.valKeys[0], rewardAddr, round)
	require.NoError(t, err)
	c := td.makeCertificateAndSign(t, b.Hash(), round, valKeys...)

	return b, c
}

func (td *testData) makeCertificateAndSign(t *testing.T, blockHash hash.Hash, round int16,
	valKeys ...*bls.ValidatorKey,
) *certificate.Certificate {
	t.Helper()

	assert.NotZero(t, len(valKeys))
	sigs := make([]*bls.Signature, len(valKeys))
	height := td.state1.LastBlockHeight()
	sb := certificate.BlockCertificateSignBytes(blockHash, height+1, round)
	committers := []int32{0, 1, 2, 3}
	var signedBy []int32

	for i, s := range valKeys {
		if s.Address() == td.valKey1.Address() {
			signedBy = append(signedBy, 0)
		}

		if s.Address() == td.valKey2.Address() {
			signedBy = append(signedBy, 1)
		}

		if s.Address() == td.valKey3.Address() {
			signedBy = append(signedBy, 2)
		}

		if s.Address() == td.valKey4.Address() {
			signedBy = append(signedBy, 3)
		}
		sigs[i] = s.Sign(sb)
	}

	absentees := util.Subtracts(committers, signedBy)
	return certificate.NewCertificate(height+1, round, committers, absentees, bls.SignatureAggregate(sigs...))
}

func (td *testData) commitBlockForAllStates(t *testing.T,
	blk *block.Block, cert *certificate.Certificate,
) {
	t.Helper()

	assert.NoError(t, td.state1.CommitBlock(blk, cert))
	assert.NoError(t, td.state2.CommitBlock(blk, cert))
	assert.NoError(t, td.state3.CommitBlock(blk, cert))
	assert.NoError(t, td.state4.CommitBlock(blk, cert))
}

func (td *testData) moveToNextHeightForAllStates(t *testing.T) {
	t.Helper()

	b, c := td.makeBlockAndCertificate(t, 0, td.valKey1, td.valKey2, td.valKey3, td.valKey4)
	td.commitBlockForAllStates(t, b, c)
}

func TestProposeBlockAndValidation(t *testing.T) {
	td := setup(t)

	td.moveToNextHeightForAllStates(t)

	b1, err := td.state1.ProposeBlock(td.state1.valKeys[0], td.RandAccAddress(), 0)
	assert.Error(t, err, "Should not propose")
	assert.Nil(t, b1)

	trx := tx.NewTransferTx(td.state1.lastInfo.BlockHeight()+1, td.valKey1.Address(),
		td.valKey2.Address(), 1000, 1000, "")
	td.HelperSignTransaction(td.valKey1.PrivateKey(), trx)
	assert.NoError(t, td.commonTxPool.AppendTx(trx))

	b2, err := td.state2.ProposeBlock(td.state2.valKeys[0], td.RandAccAddress(), 0)
	assert.NoError(t, err)
	assert.NotNil(t, b2)
	assert.Equal(t, b2.Transactions().Len(), 2)
	require.NoError(t, td.state1.ValidateBlock(b2))

	// Propose and validate again
	b3, err := td.state2.ProposeBlock(td.state2.valKeys[0], td.RandAccAddress(), 0)
	assert.NoError(t, err)
	assert.NotNil(t, b3)
	assert.Equal(t, b3.Transactions().Len(), 2)
	require.NoError(t, td.state1.ValidateBlock(b3))
}

func TestBlockSubsidyTx(t *testing.T) {
	td := setup(t)

	// Without reward address in config
	rewardAddr := td.RandAccAddress()
	trx := td.state1.createSubsidyTx(rewardAddr, 7)
	assert.True(t, trx.IsSubsidyTx())
	assert.Equal(t, trx.Payload().Value(), td.state1.params.BlockReward+7)
	assert.Equal(t, trx.Payload().(*payload.TransferPayload).From, crypto.TreasuryAddress)
	assert.Equal(t, trx.Payload().(*payload.TransferPayload).To, rewardAddr)
}

func TestBlockTime(t *testing.T) {
	td := setup(t)

	t.Run("No blocks: LastBlockTime is the genesis time", func(t *testing.T) {
		assert.Equal(t, td.state1.LastBlockTime(), td.state1.Genesis().GenesisTime())
	})

	t.Run("Commit one block: LastBlockTime is the time of the first block", func(t *testing.T) {
		blk, cert := td.makeBlockAndCertificate(t, 1, td.valKey1, td.valKey2, td.valKey3)
		assert.NoError(t, td.state1.CommitBlock(blk, cert))

		assert.NotEqual(t, td.state1.LastBlockTime(), td.state1.Genesis().GenesisTime())
		assert.Equal(t, td.state1.LastBlockTime(), blk.Header().Time())
	})
}

func TestCommitBlocks(t *testing.T) {
	td := setup(t)

	b1, c1 := td.makeBlockAndCertificate(t, 1, td.valKey1, td.valKey2, td.valKey3)
	invBlock, invCert := td.GenerateTestBlock(1)
	assert.Error(t, td.state1.CommitBlock(invBlock, c1))
	assert.Error(t, td.state1.CommitBlock(b1, invCert))
	// No error here but block is ignored, because the height is invalid
	assert.NoError(t, td.state1.CommitBlock(b1, c1))
	assert.NoError(t, td.state1.CommitBlock(b1, c1))

	assert.Equal(t, td.state1.LastBlockHash(), b1.Hash())
	assert.Equal(t, td.state1.LastBlockTime(), b1.Header().Time())
	assert.Equal(t, td.state1.LastCertificate().Hash(), c1.Hash())
	assert.Equal(t, td.state1.LastBlockHeight(), uint32(1))
	assert.Equal(t, td.state1.Genesis().Hash(), td.state2.Genesis().Hash())
}

func TestCommitSandbox(t *testing.T) {
	t.Run("Add new account", func(t *testing.T) {
		td := setup(t)

		addr := td.RandAccAddress()
		sb := td.state1.concreteSandbox()
		newAcc := sb.MakeNewAccount(addr)
		newAcc.AddToBalance(1)
		td.state1.commitSandbox(sb, 0)

		assert.NotNil(t, td.state1.AccountByAddress(addr))
	})

	t.Run("Add new validator", func(t *testing.T) {
		td := setup(t)

		pub, _ := td.RandBLSKeyPair()
		sb := td.state1.concreteSandbox()
		newVal := sb.MakeNewValidator(pub)
		newVal.AddToStake(123)
		sb.UpdateValidator(newVal)
		td.state1.commitSandbox(sb, 0)

		assert.True(t, td.state1.store.HasValidator(pub.ValidatorAddress()))
	})

	t.Run("Modify account", func(t *testing.T) {
		td := setup(t)

		sb := td.state1.concreteSandbox()
		acc := sb.Account(crypto.TreasuryAddress)
		acc.SubtractFromBalance(1)
		sb.UpdateAccount(crypto.TreasuryAddress, acc)
		td.state1.commitSandbox(sb, 0)

		acc1 := td.state1.AccountByAddress(crypto.TreasuryAddress)
		assert.Equal(t, acc1.Balance(), acc.Balance())
	})

	t.Run("Modify validator", func(t *testing.T) {
		td := setup(t)

		sb := td.state1.concreteSandbox()
		val := sb.Validator(td.valKey2.Address())
		val.AddToStake(2002)
		sb.UpdateValidator(val)
		td.state1.commitSandbox(sb, 0)

		val1, _ := td.state1.store.Validator(td.valKey2.Address())
		assert.Equal(t, val1.Stake(), val.Stake())
	})

	t.Run("Move committee", func(t *testing.T) {
		td := setup(t)

		nextProposer := td.state1.committee.Proposer(1)

		sb := td.state1.concreteSandbox()
		td.state1.commitSandbox(sb, 0)

		assert.Equal(t, td.state1.committee.Proposer(0).Address(), nextProposer.Address())
	})

	t.Run("Move committee next round", func(t *testing.T) {
		td := setup(t)

		nextNextProposer := td.state1.committee.Proposer(2)

		sb := td.state1.concreteSandbox()
		td.state1.commitSandbox(sb, 1)

		assert.Equal(t, td.state1.committee.Proposer(0).Address(), nextNextProposer.Address())
	})
}

func TestUpdateLastCertificate(t *testing.T) {
	td := setup(t)

	blk, cert := td.makeBlockAndCertificate(t, 1, td.valKey1, td.valKey2, td.valKey3)
	td.commitBlockForAllStates(t, blk, cert)

	invValKey := td.RandValKey()
	notActiveValKey := td.RandValKey()
	valNum := int32(4) // [0..3] are in the committee now
	val := validator.NewValidator(notActiveValKey.PublicKey(), valNum)
	td.state1.store.UpdateValidator(val)

	v1 := vote.NewPrepareVote(blk.Hash(), cert.Height(), cert.Round(), td.valKey3.Address())
	v2 := vote.NewPrecommitVote(blk.Hash(), cert.Height()+1, cert.Round(), td.valKey3.Address())
	v3 := vote.NewPrecommitVote(blk.Hash(), cert.Height(), cert.Round()-1, td.valKey3.Address())
	v4 := vote.NewPrecommitVote(blk.Hash(), cert.Height(), cert.Round(), td.valKey4.Address())
	v5 := vote.NewPrecommitVote(blk.Hash(), cert.Height(), cert.Round(), invValKey.Address())
	v6 := vote.NewPrecommitVote(blk.Hash(), cert.Height(), cert.Round(), notActiveValKey.Address())
	v7 := vote.NewPrecommitVote(blk.Hash(), cert.Height(), cert.Round(), td.valKey4.Address())

	td.HelperSignVote(td.valKey3, v1)
	td.HelperSignVote(td.valKey3, v2)
	td.HelperSignVote(td.valKey3, v3)
	td.HelperSignVote(invValKey, v4)
	td.HelperSignVote(invValKey, v5)
	td.HelperSignVote(notActiveValKey, v6)
	td.HelperSignVote(td.valKey4, v7)

	tests := []struct {
		vote *vote.Vote
		err  error
	}{
		{v1, InvalidVoteForCertificateError{Vote: v1}},
		{v2, InvalidVoteForCertificateError{Vote: v2}},
		{v3, InvalidVoteForCertificateError{Vote: v3}},
		{v4, crypto.ErrInvalidSignature},
		{v5, store.ErrNotFound},
		{v6, InvalidVoteForCertificateError{Vote: v6}},
		{v7, nil},
	}

	for i, test := range tests {
		err := td.state1.UpdateLastCertificate(test.vote)
		assert.ErrorIs(t, test.err, err, "error not matched for test %v", i)
	}
}

func TestInvalidProposerProposeBlock(t *testing.T) {
	td := setup(t)

	_, err := td.state2.ProposeBlock(td.state2.valKeys[0], td.RandAccAddress(), 0)
	assert.Error(t, err, "Should not propose")
	_, err = td.state2.ProposeBlock(td.state2.valKeys[0], td.RandAccAddress(), 1)
	assert.NoError(t, err, "Should propose")
}

func TestBlockProposal(t *testing.T) {
	td := setup(t)

	td.moveToNextHeightForAllStates(t)

	t.Run("validity of proposed block", func(t *testing.T) {
		b, err := td.state2.ProposeBlock(td.state2.valKeys[0], td.RandAccAddress(), 0)
		assert.NoError(t, err)
		assert.NoError(t, td.state1.ValidateBlock(b))
	})

	t.Run("Tx pool has two subsidy transactions", func(t *testing.T) {
		trx := td.state3.createSubsidyTx(td.RandAccAddress(), 0)
		assert.NoError(t, td.state3.txPool.AppendTx(trx))

		// Moving to the next round
		b, err := td.state3.ProposeBlock(td.state3.valKeys[0], td.RandAccAddress(), 1)
		assert.NoError(t, err)
		assert.NoError(t, td.state1.ValidateBlock(b))
		assert.Equal(t, b.Transactions().Len(), 1)
	})
}

func TestInvalidBlock(t *testing.T) {
	td := setup(t)

	invBlk, _ := td.GenerateTestBlock(td.RandHeight())
	assert.Error(t, td.state1.ValidateBlock(invBlk))
}

func TestForkDetection(t *testing.T) {
	td := setup(t)

	td.moveToNextHeightForAllStates(t)

	b2m, c2m := td.makeBlockAndCertificate(t, 0, td.valKey1, td.valKey2, td.valKey3)
	b2f, c2f := td.makeBlockAndCertificate(t, 1, td.valKey1, td.valKey2, td.valKey3)
	assert.NoError(t, td.state1.CommitBlock(b2m, c2m))
	assert.NoError(t, td.state2.CommitBlock(b2m, c2m))
	assert.NoError(t, td.state3.CommitBlock(b2m, c2m))
	assert.NoError(t, td.state4.CommitBlock(b2f, c2f))

	b3, c3 := td.makeBlockAndCertificate(t, 0, td.valKey1, td.valKey2, td.valKey3)

	assert.NoError(t, td.state1.CommitBlock(b3, c3))
	assert.NoError(t, td.state2.CommitBlock(b3, c3))
	assert.NoError(t, td.state3.CommitBlock(b3, c3))
	t.Run("Fork is detected, Should panic ", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("The code did not panic")
			}
		}()
		assert.Error(t, td.state4.CommitBlock(b3, c3))
	})
}

func TestSortition(t *testing.T) {
	td := setup(t)

	pub, prv := td.RandBLSKeyPair()
	valKey := bls.NewValidatorKey(prv)
	mockStore := store.MockingStore(td.TestSuite)
	St1, _ := LoadOrNewState(td.state1.genDoc, []*bls.ValidatorKey{valKey}, mockStore, td.commonTxPool, nil)
	stNew := St1.(*state)

	assert.False(t, stNew.evaluateSortition()) //  not a validator
	assert.Equal(t, td.state1.CommitteePower(), int64(4))

	height := uint32(1)
	for ; height <= 15; height++ {
		if height == 6 {
			trx := tx.NewBondTx(1, td.valKey1.Address(),
				pub.ValidatorAddress(), pub, 1000000000, 100000, "")
			td.HelperSignTransaction(td.valKey1.PrivateKey(), trx)

			assert.NoError(t, td.commonTxPool.AppendTx(trx))
		}

		b, c := td.makeBlockAndCertificate(t, 0, td.valKey1, td.valKey2, td.valKey3, td.valKey4)
		td.commitBlockForAllStates(t, b, c)
		require.NoError(t, stNew.CommitBlock(b, c))
	}

	assert.False(t, stNew.evaluateSortition()) //  bonding period

	// Certificate next block
	b, c := td.makeBlockAndCertificate(t, 0, td.valKey1, td.valKey2, td.valKey3, td.valKey4)
	td.commitBlockForAllStates(t, b, c)
	require.NoError(t, stNew.CommitBlock(b, c))
	height++

	assert.True(t, stNew.evaluateSortition())                             //  ok
	assert.False(t, td.state1.committee.Contains(pub.ValidatorAddress())) // still not in the committee

	// ---------------------------------------------
	// Certificate next block, new validator should be in the committee now
	b, c = td.makeBlockAndCertificate(t, 0, td.valKey1, td.valKey2, td.valKey3, td.valKey4)
	td.commitBlockForAllStates(t, b, c)
	require.NoError(t, stNew.CommitBlock(b, c))

	assert.True(t, stNew.evaluateSortition()) // in the committee
	assert.True(t, td.state1.committee.Contains(td.valKey1.Address()))
	assert.True(t, td.state1.committee.Contains(pub.ValidatorAddress()))

	// ---------------------------------------------
	// Let's save and load td.state1
	_ = td.state1.Close()
	St1, _ = LoadOrNewState(td.state1.genDoc, []*bls.ValidatorKey{td.valKey1}, mockStore, td.commonTxPool, nil)
	st1 := St1.(*state)

	// ---------------------------------------------
	// Let's commit another block with the new committee
	height++

	b14, err := stNew.ProposeBlock(stNew.valKeys[0], td.RandAccAddress(), 3)
	require.NoError(t, err)
	require.NotNil(t, b14)

	sigs := make([]*bls.Signature, 4)
	sb := certificate.BlockCertificateSignBytes(b14.Hash(), height, 3)

	sigs[0] = td.valKey2.Sign(sb)
	sigs[1] = td.valKey3.Sign(sb)
	sigs[2] = td.valKey4.Sign(sb)
	sigs[3] = valKey.Sign(sb)

	c14 := certificate.NewCertificate(height, 3, []int32{4, 0, 1, 2, 3}, []int32{0}, bls.SignatureAggregate(sigs...))

	assert.NoError(t, st1.CommitBlock(b14, c14))
	assert.NoError(t, td.state1.CommitBlock(b14, c14))
	assert.NoError(t, td.state2.CommitBlock(b14, c14))
	assert.NoError(t, td.state3.CommitBlock(b14, c14))
	assert.NoError(t, td.state4.CommitBlock(b14, c14))

	assert.Equal(t, td.state1.CommitteePower(), int64(1000000004))
	assert.Equal(t, td.state1.TotalValidators(), int32(5))
}

func TestValidateBlockTime(t *testing.T) {
	td := setup(t)

	fmt.Printf("BlockTimeInSecond: %d\n", td.state1.params.BlockIntervalInSecond)

	// Time is not rounded
	roundedNow := util.RoundNow(10)
	assert.Error(t, td.state1.validateBlockTime(roundedNow.Add(-15*time.Second)))
	assert.Error(t, td.state1.validateBlockTime(roundedNow.Add(-5*time.Second)))
	assert.Error(t, td.state1.validateBlockTime(roundedNow.Add(5*time.Second)))
	assert.Error(t, td.state1.validateBlockTime(roundedNow.Add(15*time.Second)))

	t.Run("Last block is committed 10 seconds ago", func(t *testing.T) {
		td.state1.lastInfo.UpdateBlockTime(roundedNow.Add(-10 * time.Second))

		// Before or same as the last block time
		assert.Error(t, td.state1.validateBlockTime(roundedNow.Add(-20*time.Second)))
		assert.Error(t, td.state1.validateBlockTime(roundedNow.Add(-10*time.Second)))

		// Ok
		assert.NoError(t, td.state1.validateBlockTime(roundedNow))
		assert.NoError(t, td.state1.validateBlockTime(roundedNow.Add(10*time.Second)))
		assert.Equal(t, td.state1.proposeNextBlockTime(), roundedNow, "Invalid proposed time for the next block")

		// More than the threshold
		assert.Error(t, td.state1.validateBlockTime(roundedNow.Add(20*time.Second)))
	})

	t.Run("Last block is committed one minute ago", func(t *testing.T) {
		td.state1.lastInfo.UpdateBlockTime(roundedNow.Add(-1 * time.Minute)) // One minute ago

		// Before or same as the last block time
		assert.Error(t, td.state1.validateBlockTime(td.state1.lastInfo.BlockTime().Add(-10*time.Second)))
		assert.Error(t, td.state1.validateBlockTime(td.state1.lastInfo.BlockTime()))

		// Ok
		assert.NoError(t, td.state1.validateBlockTime(roundedNow.Add(-10*time.Second)))
		assert.NoError(t, td.state1.validateBlockTime(roundedNow))
		assert.NoError(t, td.state1.validateBlockTime(roundedNow.Add(10*time.Second)))
		assert.Equal(t, td.state1.proposeNextBlockTime(), roundedNow, "Invalid proposed time for the next block")

		// More than the threshold
		assert.Error(t, td.state1.validateBlockTime(roundedNow.Add(20*time.Second)))
	})

	t.Run("Last block is committed in future", func(t *testing.T) {
		td.state1.lastInfo.UpdateBlockTime(roundedNow.Add(1 * time.Minute)) // One minute later

		assert.Error(t, td.state1.validateBlockTime(td.state1.lastInfo.BlockTime().Add(+1*time.Minute)))

		// Before the last block time
		assert.Error(t, td.state1.validateBlockTime(td.state1.lastInfo.BlockTime().Add(-10*time.Second)))
		assert.Error(t, td.state1.validateBlockTime(td.state1.lastInfo.BlockTime()))

		// Ok
		assert.NoError(t, td.state1.validateBlockTime(td.state1.lastInfo.BlockTime().Add(10*time.Second)))
		assert.NoError(t, td.state1.validateBlockTime(td.state1.lastInfo.BlockTime().Add(20*time.Second)))

		// More than the threshold
		assert.Error(t, td.state1.validateBlockTime(td.state1.lastInfo.BlockTime().Add(30*time.Second)))
	})
}

func TestInvalidBlockVersion(t *testing.T) {
	td := setup(t)

	td.state1.params.BlockVersion = 2
	b, _ := td.state1.ProposeBlock(td.state1.valKeys[0], td.RandAccAddress(), 0)
	assert.Error(t, td.state2.ValidateBlock(b))
}

func TestInvalidBlockTime(t *testing.T) {
	td := setup(t)

	td.moveToNextHeightForAllStates(t)

	validBlock, _ := td.makeBlockAndCertificate(t, 0, td.valKey1, td.valKey2, td.valKey3, td.valKey4)
	invalidBlock := block.MakeBlock(
		validBlock.Header().Version(),
		validBlock.Header().Time().Add(30*time.Second),
		validBlock.Transactions(),
		validBlock.Header().PrevBlockHash(),
		validBlock.Header().StateRoot(),
		validBlock.PrevCertificate(),
		validBlock.Header().SortitionSeed(),
		validBlock.Header().ProposerAddress())

	assert.NoError(t, td.state1.ValidateBlock(validBlock))
	assert.Error(t, td.state1.ValidateBlock(invalidBlock))
}

func TestValidatorHelpers(t *testing.T) {
	td := setup(t)

	t.Run("Should return nil for non-existing Validator Address", func(t *testing.T) {
		_, prv5 := td.RandBLSKeyPair()
		valKey := bls.NewValidatorKey(prv5)
		nonExistenceValidator := td.state1.ValidatorByAddress(valKey.PublicKey().ValidatorAddress())
		assert.Nil(t, nonExistenceValidator, "State 1 returned non-nil For non-existing validator")
		nonExistenceValidator = td.state2.ValidatorByAddress(valKey.PublicKey().ValidatorAddress())
		assert.Nil(t, nonExistenceValidator, "State 2 returned non-nil For non-existing validator")
		nonExistenceValidator = td.state3.ValidatorByAddress(valKey.PublicKey().ValidatorAddress())
		assert.Nil(t, nonExistenceValidator, "State 3 returned non-nil For non-existing validator")
		nonExistenceValidator = td.state4.ValidatorByAddress(valKey.PublicKey().ValidatorAddress())
		assert.Nil(t, nonExistenceValidator, "State 4 returned non-nil For non-existing validator")
	})

	t.Run("Should return validator for valid committee Validator Address", func(t *testing.T) {
		existingValidator := td.state4.ValidatorByAddress(td.valKey1.Address())
		assert.NotNil(t, existingValidator)
		assert.Zero(t, existingValidator.Number())
	})

	t.Run("Should return validator for corresponding Validator number", func(t *testing.T) {
		existingValidator := td.state4.ValidatorByNumber(1)
		assert.NotNil(t, existingValidator)
		assert.Equal(t, td.valKey2.Address(), existingValidator.Address())
	})

	t.Run("Should return nil for invalid Validator number", func(t *testing.T) {
		nonExistenceValidator := td.state4.ValidatorByNumber(10)
		assert.Nil(t, nonExistenceValidator)
	})
}

func TestLoadState(t *testing.T) {
	td := setup(t)

	// Add a bond transactions to change total power (stake)
	pub, _ := td.RandBLSKeyPair()
	tx2 := tx.NewBondTx(1, td.valKey1.Address(),
		pub.ValidatorAddress(), pub, 8888000, 8888, "")
	td.HelperSignTransaction(td.valKey1.PrivateKey(), tx2)

	assert.NoError(t, td.commonTxPool.AppendTx(tx2))

	for i := 0; i < 4; i++ {
		td.moveToNextHeightForAllStates(t)
	}
	b5, c5 := td.makeBlockAndCertificate(t, 1, td.valKey1, td.valKey2, td.valKey3, td.valKey4)
	td.commitBlockForAllStates(t, b5, c5)

	b6, c6 := td.makeBlockAndCertificate(t, 0, td.valKey1, td.valKey2, td.valKey3, td.valKey4)

	// Load last state info
	st1Load, err := LoadOrNewState(td.state1.genDoc, []*bls.ValidatorKey{td.valKey1},
		td.state1.store, td.commonTxPool, nil)
	require.NoError(t, err)

	assert.Equal(t, td.state1.store.TotalAccounts(), st1Load.(*state).store.TotalAccounts())
	assert.Equal(t, td.state1.store.TotalValidators(), st1Load.(*state).store.TotalValidators())
	assert.Equal(t, td.state1.committee.Committers(), st1Load.(*state).committee.Committers())
	assert.Equal(t, td.state1.committee.TotalPower(), st1Load.(*state).committee.TotalPower())
	assert.Equal(t, td.state1.TotalPower(), st1Load.(*state).TotalPower())
	assert.Equal(t, td.state1.store.TotalAccounts(), int32(6))

	require.NoError(t, st1Load.CommitBlock(b6, c6))
	require.NoError(t, td.state2.CommitBlock(b6, c6))
}

func TestLoadStateAfterChangingGenesis(t *testing.T) {
	td := setup(t)

	// Let's commit some blocks
	i := 0
	for ; i < 10; i++ {
		td.moveToNextHeightForAllStates(t)
	}

	_, err := LoadOrNewState(td.state1.genDoc, []*bls.ValidatorKey{td.valKey1},
		td.state1.store, txpool.MockingTxPool(), nil)
	require.NoError(t, err)

	pub, _ := td.RandBLSKeyPair()
	val := validator.NewValidator(pub, 4)
	vals := append(td.state1.genDoc.Validators(), val)

	genDoc := genesis.MakeGenesis(
		td.state1.genDoc.GenesisTime(),
		td.state1.genDoc.Accounts(),
		vals,
		td.state1.genDoc.Params())

	// Load last state info after modifying genesis
	_, err = LoadOrNewState(genDoc, []*bls.ValidatorKey{td.valKey1}, td.state1.store, txpool.MockingTxPool(), nil)
	require.Error(t, err)
}

func TestSetBlockTime(t *testing.T) {
	td := setup(t)

	t.Run("Last block time is a bit far in past", func(t *testing.T) {
		td.state1.lastInfo.UpdateBlockTime(util.RoundNow(10).Add(-20 * time.Second))
		b, _ := td.state1.ProposeBlock(td.state1.valKeys[0], td.RandAccAddress(), 0)
		fmt.Printf("last block time: %s\nproposed time  : %s\n", td.state1.lastInfo.BlockTime(), b.Header().Time().UTC())
		assert.True(t, b.Header().Time().After(td.state1.lastInfo.BlockTime()))
		assert.True(t, b.Header().Time().Before(util.Now().Add(10*time.Second)))
		assert.Zero(t, b.Header().Time().Second()%10)
	})

	t.Run("Last block time is almost good", func(t *testing.T) {
		td.state1.lastInfo.UpdateBlockTime(util.RoundNow(10).Add(-10 * time.Second))
		b, _ := td.state1.ProposeBlock(td.state1.valKeys[0], td.RandAccAddress(), 0)
		fmt.Printf("last block time: %s\nproposed time  : %s\n", td.state1.lastInfo.BlockTime(), b.Header().Time().UTC())
		assert.True(t, b.Header().Time().After(td.state1.lastInfo.BlockTime()))
		assert.True(t, b.Header().Time().Before(util.Now().Add(10*time.Second)))
		assert.Zero(t, b.Header().Time().Second()%10)
	})

	// After our time
	t.Run("Last block time is in near future", func(t *testing.T) {
		td.state1.lastInfo.UpdateBlockTime(util.RoundNow(10).Add(+10 * time.Second))
		b, _ := td.state1.ProposeBlock(td.state1.valKeys[0], td.RandAccAddress(), 0)
		fmt.Printf("last block time: %s\nproposed time  : %s\n", td.state1.lastInfo.BlockTime(), b.Header().Time().UTC())
		assert.True(t, b.Header().Time().After(td.state1.lastInfo.BlockTime()))
		assert.Zero(t, b.Header().Time().Second()%10)
	})

	t.Run("Last block time is more than a block in future", func(t *testing.T) {
		td.state1.lastInfo.UpdateBlockTime(util.RoundNow(10).Add(+20 * time.Second))
		b, _ := td.state1.ProposeBlock(td.state1.valKeys[0], td.RandAccAddress(), 0)
		fmt.Printf("last block time: %s\nproposed time  : %s\n", td.state1.lastInfo.BlockTime(), b.Header().Time().UTC())
		assert.True(t, b.Header().Time().After(td.state1.lastInfo.BlockTime()))
		assert.Zero(t, b.Header().Time().Second()%10)
	})
}

func TestIsValidator(t *testing.T) {
	td := setup(t)

	assert.True(t, td.state1.IsInCommittee(td.valKey1.Address()))
	assert.True(t, td.state1.IsProposer(td.valKey1.Address(), 0))
	assert.True(t, td.state1.IsProposer(td.valKey2.Address(), 1))
	assert.True(t, td.state1.IsInCommittee(td.valKey2.Address()))
	assert.True(t, td.state1.IsValidator(td.valKey2.Address()))

	addr := td.RandAccAddress()
	assert.False(t, td.state1.IsInCommittee(addr))
	assert.False(t, td.state1.IsProposer(addr, 0))
	assert.False(t, td.state1.IsInCommittee(addr))
	assert.False(t, td.state1.IsValidator(addr))
}

func TestCalculatingGenesisState(t *testing.T) {
	td := setup(t)

	r := td.state1.calculateGenesisStateRootFromGenesisDoc()
	assert.Equal(t, td.state1.stateRoot(), r)
}

func TestCommittingInvalidBlock(t *testing.T) {
	td := setup(t)

	td.moveToNextHeightForAllStates(t)

	txs := block.NewTxs()
	trx := td.state2.createSubsidyTx(td.RandAccAddress(), 0)
	txs.Append(trx)
	b := block.MakeBlock(2, util.Now(), txs, td.state2.lastInfo.BlockHash(), td.state2.stateRoot(),
		td.state2.lastInfo.Certificate(), td.state2.lastInfo.SortitionSeed(), td.state2.valKeys[0].Address())
	c := td.makeCertificateAndSign(t, b.Hash(), 0, td.valKey1, td.valKey2, td.valKey3, td.valKey4)

	// td.state1 receives a block with version 2 and rejects it.
	// It is possible that the same block would be considered valid by td.state2.
	assert.Error(t, td.state1.CommitBlock(b, c))
}

func TestCalcFee(t *testing.T) {
	td := setup(t)
	tests := []struct {
		amount          int64
		pldType         payload.Type
		fee             int64
		expectedFee     int64
		expectedErrCode int
	}{
		{1, payload.TypeTransfer, 1, td.state1.params.MinimumFee, errors.ErrInvalidFee},
		{1, payload.TypeWithdraw, 1001, td.state1.params.MinimumFee, errors.ErrInvalidFee},
		{1, payload.TypeBond, 1000, td.state1.params.MinimumFee, errors.ErrNone},

		{1 * 1e9, payload.TypeTransfer, 1, 100000, errors.ErrInvalidFee},
		{1 * 1e9, payload.TypeWithdraw, 100001, 100000, errors.ErrInvalidFee},
		{1 * 1e9, payload.TypeBond, 100000, 100000, errors.ErrNone},

		{1 * 1e12, payload.TypeTransfer, 1, 1000000, errors.ErrInvalidFee},
		{1 * 1e12, payload.TypeWithdraw, 1000001, 1000000, errors.ErrInvalidFee},
		{1 * 1e12, payload.TypeBond, 1000000, 1000000, errors.ErrNone},

		{1 * 1e12, payload.TypeSortition, 0, 0, errors.ErrInvalidFee},
		{1 * 1e12, payload.TypeUnbond, 0, 0, errors.ErrNone},
	}
	for _, test := range tests {
		fee, err := td.state2.CalculateFee(test.amount, test.pldType)
		assert.NoError(t, err)
		assert.Equal(t, test.expectedFee, fee)

		_, err = td.state2.CalculateFee(test.amount, 6)
		assert.Error(t, err)
	}
}
