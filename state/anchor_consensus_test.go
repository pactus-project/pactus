package state

import (
	"bytes"
	"slices"
	"testing"

	"github.com/pactus-project/gopkg/pipeline"
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/execution/executor"
	"github.com/pactus-project/pactus/store"
	"github.com/pactus-project/pactus/txpool"
	"github.com/pactus-project/pactus/types"
	"github.com/pactus-project/pactus/types/account"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/block"
	"github.com/pactus-project/pactus/types/protocol"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/tx/payload"
	"github.com/pactus-project/pactus/types/validator"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

// newPeer starts a second, independent node from the same genesis and replays
// every block the main node has committed so far.
func (td *testData) newPeer(t *testing.T) *state {
	t.Helper()

	peerStore := store.NewFakeStore(td.TestSuite)
	peerStore.EXPECT().IsBanned(gomock.Any()).Return(false).AnyTimes()
	peerStore.EXPECT().RecentTransaction(gomock.Any()).DoAndReturn(func(id tx.ID) bool {
		_, ok := peerStore.RecentTxs[id]

		return ok
	}).AnyTimes()

	peerPool := txpool.NewFakeTxPool(td.TestSuite)
	peerPool.EXPECT().SetNewSandboxAndRecheck(gomock.Any()).Return().AnyTimes()
	peerPool.EXPECT().HandleCommittedBlock(gomock.Any()).Return().AnyTimes()

	// No validator key, so the peer never runs sortition.
	loaded, err := LoadOrNewState(t.Context(), td.state.genDoc, nil,
		peerStore, peerPool, pipeline.New[any](t.Context()))
	require.NoError(t, err)
	peer := loaded.(*state)

	mainStore := td.state.store.(*store.FakeStore)
	for height := types.Height(1); height <= td.state.LastBlockHeight(); height++ {
		require.NoError(t, peer.CommitBlock(mainStore.FakeBlocks[height], mainStore.FakeCertificates[height]))
	}
	require.Equal(t, td.state.stateRoot(), peer.stateRoot())

	return peer
}

// commitOnBoth commits a block proposed by the main node on both nodes.
// The peer validates it first, as a follower would.
func (td *testData) commitOnBoth(t *testing.T, peer *state, blk *block.Block) {
	t.Helper()

	require.NoError(t, peer.ValidateBlock(blk, 0))
	cert := td.makeCertificateAndSign(t, blk.Hash(), 0)
	require.NoError(t, td.state.CommitBlock(blk, cert))
	require.NoError(t, peer.CommitBlock(blk, cert))
	require.Equal(t, td.state.stateRoot(), peer.stateRoot())
}

// requireSameAnchors checks that both nodes store the same bytes for the given
// accounts and list the same anchors.
func requireSameAnchors(t *testing.T, left, right *state, addrs ...crypto.Address) {
	t.Helper()

	for _, addr := range addrs {
		leftAcc, err := left.AccountByAddress(addr)
		require.NoError(t, err)
		rightAcc, err := right.AccountByAddress(addr)
		require.NoError(t, err)
		leftBytes, err := leftAcc.Bytes()
		require.NoError(t, err)
		rightBytes, err := rightAcc.Bytes()
		require.NoError(t, err)
		require.Equal(t, leftBytes, rightBytes, addr.String())
	}

	leftItems, leftTotal := left.ListAnchors(0, 100)
	rightItems, rightTotal := right.ListAnchors(0, 100)
	require.Equal(t, leftTotal, rightTotal)
	require.Len(t, rightItems, len(leftItems))
	for i := range leftItems {
		require.Equal(t, leftItems[i].Address, rightItems[i].Address)
		require.Equal(t, leftItems[i].Account.Hash(), rightItems[i].Account.Hash())
	}
}

// totalCoinsOf sums spendable balances, locked anchor deposits and stakes.
func totalCoinsOf(chain *state) amount.Amount {
	total := amount.Amount(0)
	chain.store.IterateAccounts(func(_ crypto.Address, acc *account.Account) bool {
		total += acc.Balance() + acc.LockedDeposit()

		return false
	})
	chain.store.IterateValidators(func(val *validator.Validator) bool {
		total += val.Stake()

		return false
	})

	return total
}

func TestAnchorPeerReachesSameState(t *testing.T) {
	td := setupWithVersion(t, protocol.ProtocolVersion5)
	td.fakeTxPool.EXPECT().AppendTxAndBroadcast(gomock.Any()).Return(nil).AnyTimes()
	peer := td.newPeer(t)
	sender := td.sender()
	receiver := td.RandAccAddress()

	steps := []struct {
		name      string
		txs       func() block.Txs
		hasAnchor bool
	}{
		{
			name: "create",
			txs: func() block.Txs {
				return block.Txs{td.anchor(td.nextHeight(), payload.AnchorActionSet,
					bytes.Repeat([]byte{0x61}, 32), "create", executor.MinAnchorDeposit, 1)}
			},
			hasAnchor: true,
		},
		{
			name: "update with a deposit, next to a transfer",
			txs: func() block.Txs {
				transfer := tx.NewTransferTx(td.nextHeight(), sender, receiver, 7, 1)
				td.HelperSignTransaction(td.genAccKey, transfer)

				return block.Txs{
					td.anchor(td.nextHeight(), payload.AnchorActionSet,
						bytes.Repeat([]byte{0x62}, 64), "update", 5, 1),
					transfer,
				}
			},
			hasAnchor: true,
		},
		{
			name: "delete",
			txs: func() block.Txs {
				return block.Txs{td.anchor(td.nextHeight(), payload.AnchorActionDelete, nil, "", 0, 1)}
			},
			hasAnchor: false,
		},
		{
			name: "recreate",
			txs: func() block.Txs {
				return block.Txs{td.anchor(td.nextHeight(), payload.AnchorActionSet,
					bytes.Repeat([]byte{0x63}, 48), "again", executor.MinAnchorDeposit+3, 2)}
			},
			hasAnchor: true,
		},
	}

	for _, step := range steps {
		txs := step.txs()
		blk := td.propose(t, txs)
		require.Equal(t, len(txs)+1, blk.Transactions().Len(), step.name)
		td.commitOnBoth(t, peer, blk)

		peerAcc, err := peer.AccountByAddress(sender)
		require.NoError(t, err)
		require.Equal(t, step.hasAnchor, peerAcc.HasAnchor(), step.name)
		requireSameAnchors(t, td.state, peer, sender)
	}
}

type chainSender struct {
	addr crypto.Address
	prv  crypto.PrivateKey
}

// randomChainTx returns a random anchor or transfer transaction.
// Many of them are invalid on purpose; the proposer must leave those out.
func (td *testData) randomChainTx(senders []chainSender) *tx.Tx {
	sender := senders[td.RandIntMax(len(senders))]
	fee := amount.Amount(td.RandIntMax(4))
	lockTime := td.nextHeight()

	var trx *tx.Tx
	switch td.RandIntMax(4) {
	case 0, 1:
		deposits := []amount.Amount{
			0,
			executor.MinAnchorDeposit - 1,
			executor.MinAnchorDeposit,
			executor.MinAnchorDeposit + td.RandAmount(),
			1000 * executor.MinAnchorDeposit,
		}
		trx = tx.NewAnchorTx(lockTime, sender.addr, payload.AnchorActionSet,
			td.RandBytes(32+td.RandIntMax(33)), td.RandString(td.RandIntMax(20)),
			uint8(td.RandIntMax(256)), deposits[td.RandIntMax(len(deposits))], fee)
	case 2:
		trx = tx.NewAnchorTx(lockTime, sender.addr, payload.AnchorActionDelete, nil, "", 0, 0, fee)
	default:
		receiver := senders[td.RandIntMax(len(senders))].addr
		trx = tx.NewTransferTx(lockTime, sender.addr, receiver, td.RandAmount(20*executor.MinAnchorDeposit), fee)
	}
	td.HelperSignTransaction(sender.prv, trx)

	return trx
}

// Two nodes follow a chain of random anchor and transfer transactions.
// They must agree on every state root, and no coin is created or destroyed.
func TestAnchorRandomChainTwoNodes(t *testing.T) {
	td := setupWithVersion(t, protocol.ProtocolVersion5)
	td.fakeTxPool.EXPECT().AppendTxAndBroadcast(gomock.Any()).Return(nil).AnyTimes()
	peer := td.newPeer(t)
	supply := totalCoinsOf(td.state)
	require.Equal(t, supply, totalCoinsOf(peer))

	blsPub, blsPrv := td.RandBLSKeyPair()
	edPub, edPrv := td.RandEd25519KeyPair()
	secpPub, secpPrv := td.RandSecp256k1KeyPair()
	senders := []chainSender{
		{addr: blsPub.AccountAddress(), prv: blsPrv},
		{addr: edPub.AccountAddress(), prv: edPrv},
		{addr: secpPub.AccountAddress(), prv: secpPrv},
	}

	funding := make(block.Txs, 0, len(senders))
	for _, sender := range senders {
		trx := tx.NewTransferTx(td.nextHeight(), td.sender(), sender.addr, 10*executor.MinAnchorDeposit, 0)
		td.HelperSignTransaction(td.genAccKey, trx)
		funding = append(funding, trx)
	}
	td.commitOnBoth(t, peer, td.propose(t, funding))

	blocks := 25
	if testing.Short() {
		blocks = 8
	}
	submitted, included := 0, 0
	for i := 0; i < blocks; i++ {
		txs := make(block.Txs, 0, 6)
		if i == 0 {
			// Make sure at least one anchor exists early in the chain.
			first := tx.NewAnchorTx(td.nextHeight(), senders[0].addr, payload.AnchorActionSet,
				bytes.Repeat([]byte{0x71}, 32), "first", 0, executor.MinAnchorDeposit, 1)
			td.HelperSignTransaction(senders[0].prv, first)
			txs = append(txs, first)
		}
		for n := td.RandIntNonZero(6); n > 0; n-- {
			txs = append(txs, td.randomChainTx(senders))
		}

		// ProposeBlock removes invalid transactions in place, so give it a copy.
		blk := td.propose(t, slices.Clone(txs))
		for _, trx := range txs {
			if trx.Payload().Type() == payload.TypeAnchor {
				submitted++
				if blockHasTx(blk, trx.ID()) {
					included++
				}
			}
		}
		td.commitOnBoth(t, peer, blk)
		require.Equal(t, supply, totalCoinsOf(td.state), "block %d", i)
		require.Equal(t, supply, totalCoinsOf(peer), "block %d", i)
	}

	addrs := []crypto.Address{senders[0].addr, senders[1].addr, senders[2].addr}
	requireSameAnchors(t, td.state, peer, addrs...)
	require.Positive(t, included, "no anchor reached the chain")
	t.Logf("anchor transactions: %d submitted, %d included", submitted, included)
}

// A node restarted from its store keeps the anchors, the anchor index and the
// state root, and can keep committing anchor blocks.
func TestAnchorStateReload(t *testing.T) {
	td := setupWithVersion(t, protocol.ProtocolVersion5)
	td.fakeTxPool.EXPECT().AppendTxAndBroadcast(gomock.Any()).Return(nil).AnyTimes()
	sender := td.sender()

	td.commit(t, td.propose(t, block.Txs{
		td.anchor(td.nextHeight(), payload.AnchorActionSet, bytes.Repeat([]byte{0x81}, 32), "kept",
			executor.MinAnchorDeposit, 1),
	}))
	td.commit(t, td.propose(t, block.Txs{
		td.anchor(td.nextHeight(), payload.AnchorActionSet, bytes.Repeat([]byte{0x82}, 32), "updated", 0, 1),
	}))

	loaded, err := LoadOrNewState(t.Context(), td.state.genDoc, td.state.valKeys,
		td.state.store, td.fakeTxPool, pipeline.New[any](t.Context()))
	require.NoError(t, err)
	reloaded := loaded.(*state)
	require.Equal(t, td.state.stateRoot(), reloaded.stateRoot())
	require.Equal(t, td.state.params.BlockVersion, reloaded.params.BlockVersion)
	requireSameAnchors(t, td.state, reloaded, sender)

	blk := td.propose(t, block.Txs{td.anchor(td.nextHeight(), payload.AnchorActionDelete, nil, "", 0, 1)})
	cert := td.makeCertificateAndSign(t, blk.Hash(), 0)
	require.NoError(t, reloaded.CommitBlock(blk, cert))

	acc, err := reloaded.AccountByAddress(sender)
	require.NoError(t, err)
	require.False(t, acc.HasAnchor())
	_, total := reloaded.ListAnchors(0, 100)
	require.Zero(t, total)
}

// Before version 5 the pool rejects anchors. The first version-5 block
// gives the pool a version-5 sandbox, and the same transaction is accepted.
func TestAnchorPoolAdmissionFollowsActivation(t *testing.T) {
	td := setupWithVersion(t, protocol.ProtocolVersion4)
	td.fakeTxPool.EXPECT().AppendTxAndBroadcast(gomock.Any()).Return(nil).AnyTimes()

	anchorTx := td.anchor(td.nextHeight(), payload.AnchorActionSet, bytes.Repeat([]byte{0x91}, 32), "early",
		executor.MinAnchorDeposit, 1)
	require.ErrorIs(t, td.state.CheckTransaction(anchorTx),
		executor.InvalidPayloadTypeError{PayloadType: payload.TypeAnchor})

	applyCommittee(t, td.state, []protocol.Version{
		protocol.ProtocolVersion5, protocol.ProtocolVersion5,
		protocol.ProtocolVersion5, protocol.ProtocolVersion5,
	})
	upgrade := td.proposeEmpty(t)
	require.Equal(t, protocol.ProtocolVersion5, upgrade.Header().Version())
	td.commit(t, upgrade)
	require.Equal(t, protocol.ProtocolVersion5, td.state.params.BlockVersion)

	require.NoError(t, td.state.CheckTransaction(anchorTx))
}
