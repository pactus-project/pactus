package state

import (
	"errors"
	"testing"

	"github.com/pactus-project/pactus/types/block"
	"github.com/pactus-project/pactus/types/protocol"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	td := setup(t)

	assert.Same(t, td.state.genDoc, td.state.Genesis())
	assert.Equal(t, td.state.genDoc.Hash(), td.state.Genesis().Hash())
}

func TestCommitteeValidators(t *testing.T) {
	td := setup(t)

	vals := td.state.CommitteeValidators()
	require.Len(t, vals, 4)
	for _, val := range vals {
		assert.True(t, td.state.IsInCommittee(val.Address()))
	}

	t.Run("returned validators are copies of the committee members", func(t *testing.T) {
		powerBefore := td.state.CommitteeInfo().CommitteePower

		vals[0].AddToStake(1_000_000)

		assert.Equal(t, powerBefore, td.state.CommitteeInfo().CommitteePower)
		assert.Equal(t, powerBefore, td.state.CommitteeValidators()[0].Power()+
			td.state.CommitteeValidators()[1].Power()+
			td.state.CommitteeValidators()[2].Power()+
			td.state.CommitteeValidators()[3].Power())
	})
}

func TestCommitteeInfo(t *testing.T) {
	td := setup(t)

	info := td.state.CommitteeInfo()

	assert.Len(t, info.Validators, 4)
	assert.Equal(t, int64(4), info.CommitteePower)
	assert.Equal(t, int64(4), info.TotalPower)

	// The protocol versions are percentages of the committee power.
	sum := 0.0
	for _, percent := range info.ProtocolVersions {
		sum += percent
	}
	assert.InDelta(t, 100.0, sum, 1e-9)

	// It should be consistent with the chain info.
	chainInfo := td.state.ChainInfo()
	assert.Equal(t, chainInfo.CommitteePower, info.CommitteePower)
	assert.Equal(t, chainInfo.TotalPower, info.TotalPower)
	assert.Len(t, info.Validators, chainInfo.CommitteeSize)
}

func TestPendingTransactions(t *testing.T) {
	td := setup(t)

	trx1 := td.GenerateTestTransferTx()
	trx2 := td.GenerateTestTransferTx()
	unknownID := td.GenerateTestTransferTx().ID()

	t.Run("pending transaction is looked up in the transaction pool", func(t *testing.T) {
		td.fakeTxPool.EXPECT().PendingTx(trx1.ID()).Return(trx1).Times(1)
		td.fakeTxPool.EXPECT().PendingTx(unknownID).Return(nil).Times(1)

		assert.Equal(t, trx1, td.state.PendingTx(trx1.ID()))
		assert.Nil(t, td.state.PendingTx(unknownID))
	})

	t.Run("all pending transactions", func(t *testing.T) {
		td.fakeTxPool.EXPECT().AllPendingTxs().Return(block.Txs{trx1, trx2}).Times(1)

		assert.Equal(t, block.Txs{trx1, trx2}, block.Txs(td.state.AllPendingTxs()))
	})

	t.Run("adding a transaction returns the error of the pool", func(t *testing.T) {
		errRejected := errors.New("rejected")
		td.fakeTxPool.EXPECT().AppendTx(trx1).Return(nil).Times(1)
		td.fakeTxPool.EXPECT().AppendTx(trx2).Return(errRejected).Times(1)

		require.NoError(t, td.state.AddPendingTx(trx1))
		require.ErrorIs(t, td.state.AddPendingTx(trx2), errRejected)
	})

	t.Run("adding and broadcasting a transaction returns the error of the pool", func(t *testing.T) {
		errRejected := errors.New("rejected")
		td.fakeTxPool.EXPECT().AppendTxAndBroadcast(trx1).Return(nil).Times(1)
		td.fakeTxPool.EXPECT().AppendTxAndBroadcast(trx2).Return(errRejected).Times(1)

		require.NoError(t, td.state.AddPendingTxAndBroadcast(trx1))
		require.ErrorIs(t, td.state.AddPendingTxAndBroadcast(trx2), errRejected)
	})
}

func TestCommittedTx(t *testing.T) {
	td := setup(t)

	// A valid transaction, otherwise the block proposer drops it.
	trx := td.GenerateTestTransferTx(
		testsuite.TransactionWithLockTime(td.state.LastBlockHeight()),
		testsuite.TransactionWithSigner(td.genAccKey),
	)
	td.commitBlocks(t, 1, block.Txs{trx})

	committedTx, err := td.state.CommittedTx(trx.ID())
	require.NoError(t, err)
	assert.Equal(t, trx.ID(), committedTx.TxID)
	assert.Equal(t, td.state.LastBlockHeight(), committedTx.Height)

	_, err = td.state.CommittedTx(td.GenerateTestTransferTx().ID())
	require.Error(t, err)
}

func TestPublicKey(t *testing.T) {
	td := setup(t)

	pub, _ := td.RandBLSKeyPair()
	unknownAddr := td.RandAccAddress()
	td.fakeStore.EXPECT().PublicKey(pub.AccountAddress()).Return(pub, nil).Times(1)
	td.fakeStore.EXPECT().PublicKey(unknownAddr).Return(nil, errors.New("not found")).Times(1)

	got, err := td.state.PublicKey(pub.AccountAddress())
	require.NoError(t, err)
	assert.Equal(t, pub, got)

	_, err = td.state.PublicKey(unknownAddr)
	require.Error(t, err)
}

func TestAvailabilityScore(t *testing.T) {
	td := setup(t)

	// The certificates made by the test helper have validator number 3 as absentee.
	// So, validators 0, 1, and 2 never missed a block.
	assert.InDelta(t, 1.0, td.state.AvailabilityScore(0), 1e-9)
	assert.InDelta(t, 1.0, td.state.AvailabilityScore(1), 1e-9)
	assert.InDelta(t, 1.0, td.state.AvailabilityScore(2), 1e-9)
	assert.Less(t, td.state.AvailabilityScore(3), 1.0)
	assert.GreaterOrEqual(t, td.state.AvailabilityScore(3), 0.0)

	// A validator that has never been in a committee has a perfect score.
	assert.InDelta(t, 1.0, td.state.AvailabilityScore(1000), 1e-9)
}

func TestUpdateValidatorProtocolVersion(t *testing.T) {
	td := setup(t)

	addr := td.genValKeys[1].Address()
	val, err := td.state.ValidatorByAddress(addr)
	require.NoError(t, err)
	require.NotEqual(t, protocol.ProtocolVersion3, val.ProtocolVersion())

	td.state.UpdateValidatorProtocolVersion(addr, protocol.ProtocolVersion3)

	val, err = td.state.ValidatorByAddress(addr)
	require.NoError(t, err)
	assert.Equal(t, protocol.ProtocolVersion3, val.ProtocolVersion())

	// An unknown validator is ignored.
	require.NotPanics(t, func() {
		td.state.UpdateValidatorProtocolVersion(td.RandValAddress(), protocol.ProtocolVersion3)
	})
}
