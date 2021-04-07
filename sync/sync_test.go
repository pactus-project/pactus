package sync

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/committee"
	"github.com/zarbchain/zarb-go/consensus"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/logger"
	"github.com/zarbchain/zarb-go/network"
	"github.com/zarbchain/zarb-go/state"
	"github.com/zarbchain/zarb-go/sync/message"
	"github.com/zarbchain/zarb-go/sync/message/payload"
	"github.com/zarbchain/zarb-go/util"
	"github.com/zarbchain/zarb-go/validator"
)

var (
	tAliceConfig      *Config
	tBobConfig        *Config
	tAliceState       *state.MockState
	tBobState         *state.MockState
	tAliceConsensus   *consensus.MockConsensus
	tBobConsensus     *consensus.MockConsensus
	tAliceNet         *network.MockNetwork
	tBobNet           *network.MockNetwork
	tAliceSync        *synchronizer
	tBobSync          *synchronizer
	tAliceBroadcastCh chan payload.Payload
	tBobBroadcastCh   chan payload.Payload
	tAlicePeerID      peer.ID
	tBobPeerID        peer.ID
	tAnotherPeerID    peer.ID
)

type OverrideFingerprint struct {
	sync Synchronizer
	name string
}

func (o *OverrideFingerprint) Fingerprint() string {
	return o.name + o.sync.Fingerprint()
}

func init() {
	logger.InitLogger(logger.TestConfig())
	tAliceConfig = TestConfig()
	tBobConfig = TestConfig()

	tAliceConfig.Moniker = "alice"
	tBobConfig.Moniker = "bob"
}

func setup(t *testing.T) {
	_, _, priv1 := crypto.GenerateTestKeyPair()
	_, _, priv2 := crypto.GenerateTestKeyPair()
	aliceSigner := crypto.NewSigner(priv1)
	bobSigner := crypto.NewSigner(priv2)

	committee, _ := committee.GenerateTestCommittee()
	tAlicePeerID = util.RandomPeerID()
	tBobPeerID = util.RandomPeerID()
	tAnotherPeerID = util.RandomPeerID()
	tAliceState = state.MockingState(committee)
	tBobState = state.MockingState(committee)
	tAliceConsensus = consensus.MockingConsensus(tAliceState)
	tBobConsensus = consensus.MockingConsensus(tBobState)
	tAliceBroadcastCh = make(chan payload.Payload, 1000)
	tBobBroadcastCh = make(chan payload.Payload, 1000)
	tAliceNet = network.MockingNetwork(tAlicePeerID)
	tBobNet = network.MockingNetwork(tBobPeerID)

	tBobState.GenHash = tAliceState.GenHash

	// Apply 20 blocks for both Alice and Bob
	lastBlockHash := crypto.Hash{}
	for i := 0; i < 21; i++ {
		b, trxs := block.GenerateTestBlock(nil, &lastBlockHash)
		c := block.GenerateTestCertificate(b.Hash())
		lastBlockHash = b.Hash()

		tAliceState.AddBlock(i+1, b, trxs)
		tAliceState.LastBlockCertificate = c

		tBobState.AddBlock(i+1, b, trxs)
		tBobState.LastBlockCertificate = c
	}

	tAliceSync = &synchronizer{ctx: context.Background()}
	aliceSync, err := NewSynchronizer(tAliceConfig,
		aliceSigner,
		tAliceState,
		tAliceConsensus,
		tAliceNet,
		tAliceBroadcastCh,
	)
	assert.NoError(t, err)
	tAliceSync = aliceSync.(*synchronizer)
	tAliceSync.logger = logger.NewLogger("_sync", &OverrideFingerprint{name: "alice: ", sync: tAliceSync})

	tBobSync = &synchronizer{ctx: context.Background()}
	bobSync, err := NewSynchronizer(tBobConfig,
		bobSigner,
		tBobState,
		tBobConsensus,
		tBobNet,
		tBobBroadcastCh,
	)
	assert.NoError(t, err)
	tBobSync = bobSync.(*synchronizer)
	tBobSync.logger = logger.NewLogger("_sync", &OverrideFingerprint{name: "bob: ", sync: tBobSync})

	tAliceNet.OtherNet = tBobNet
	tBobNet.OtherNet = tAliceNet

	assert.NoError(t, tAliceSync.Start())
	assert.NoError(t, tBobSync.Start())

	shouldPublishPayloadWithThisType(t, tAliceNet, payload.PayloadTypeSalam)
	shouldPublishPayloadWithThisType(t, tBobNet, payload.PayloadTypeSalam)

	shouldPublishPayloadWithThisType(t, tAliceNet, payload.PayloadTypeAleyk)
	shouldPublishPayloadWithThisType(t, tBobNet, payload.PayloadTypeAleyk)

	assert.Equal(t, tAliceState.LastBlockHeight(), tBobState.LastBlockHeight())
}

func shouldPublishPayloadWithThisType(t *testing.T, net *network.MockNetwork, payloadType payload.PayloadType) *message.Message {
	timeout := time.NewTimer(2 * time.Second)

	for {
		select {
		case <-timeout.C:
			require.NoError(t, fmt.Errorf("ShouldPublishPayloadWithThisType %v: Timeout", payloadType))
			return nil
		case msg := <-net.BroadcastCh:
			logger.Info("shouldPublishPayloadWithThisType", "msg", msg, "type", payloadType.String())
			net.OtherNet.ReceivedMessage(msg, net.SelfID())
			logger.Info("Nessage sent to other peer", "msg", msg)

			if msg.Payload.Type() == payloadType {
				return msg
			}
		}
	}
}

func shouldNotPublishMessageWithThisType(t *testing.T, net *network.MockNetwork, payloadType payload.PayloadType) {
	timeout := time.NewTimer(300 * time.Millisecond)

	for {
		select {
		case <-timeout.C:
			return
		case msg := <-net.BroadcastCh:
			logger.Info("shouldNotPublishMessageWithThisType", "msg", msg, "type", payloadType.String())
			net.OtherNet.ReceivedMessage(msg, net.SelfID())
			logger.Info("Nessage sent to other peer", "msg", msg)

			assert.NotEqual(t, msg.Payload.Type(), payloadType)
		}
	}
}

func addMoreBlocksForBob(t *testing.T, count int) {
	lastBlockHash := tBobState.LastBlockHash()
	for i := 0; i < count; i++ {
		b, trxs := block.GenerateTestBlock(nil, &lastBlockHash)
		c := block.GenerateTestCertificate(b.Hash())
		lastBlockHash = b.Hash()

		tBobState.AddBlock(tBobState.LastBlockHeight()+1, b, trxs)
		tBobState.LastBlockCertificate = c
	}
}

func addMoreBlocksForBobAndAnnounceLastBlock(t *testing.T, count int) {
	addMoreBlocksForBob(t, count)

	pld := payload.NewBlockAnnouncePayload(
		tBobState.LastBlockHeight(),
		tBobState.Block(tBobState.LastBlockHeight()),
		tBobState.LastCertificate())

	tBobBroadcastCh <- pld
}

func disableHeartbeat(t *testing.T) {
	tAliceSync.heartBeatTicker.Stop()
	tBobSync.heartBeatTicker.Stop()
}

func joinAliceToTheSet(t *testing.T) {
	val := validator.NewValidator(tAliceSync.signer.PublicKey(), 4, tAliceState.LastBlockHeight())
	val.UpdateLastJoinedHeight(tAliceState.LastBlockHeight())

	assert.NoError(t, tAliceState.Committee.Update(0, []*validator.Validator{val}))
}

func joinBobToTheSet(t *testing.T) {
	val := validator.NewValidator(tBobSync.signer.PublicKey(), 5, tBobState.LastBlockHeight())
	val.UpdateLastJoinedHeight(tBobState.LastBlockHeight())

	assert.NoError(t, tAliceState.Committee.Update(0, []*validator.Validator{val}))
}

func TestAccessors(t *testing.T) {
	setup(t)

	assert.Equal(t, tAliceSync.PeerID(), tAlicePeerID)
	assert.Equal(t, len(tAliceSync.Peers()), 1)
}

func TestStop(t *testing.T) {
	setup(t)
	// Should stop normally
	tAliceSync.Stop()
	tBobSync.Stop()
}

/*

func TestSendSalamBadGenesisHash(t *testing.T) {
	setup(t)

	invGenHash := crypto.GenerateTestHash()
	_, pub, _ := crypto.GenerateTestKeyPair()

	pld := payload.NewSalamPayload("bad-genesis", pub, invGenHash, 0, 0)
	tAliceNet.CheckAndParsPayload(pld, tAnotherPeerID)

	pld2 := shouldPublishPayloadWithThisType(t, tAliceNet, payload.PayloadTypeAleyk)
	pldAleyk2 := pld2.(*payload.AleykPayload)

	assert.Equal(t, pldAleyk2.ResponseCode, payload.ResponseCodeRejected)
}

func TestSendSalamPeerBehind(t *testing.T) {
	setup(t)
	_, pub, _ := crypto.GenerateTestKeyPair()

	pld := payload.NewSalamPayload("kitty", pub, tAliceState.GenHash, 3, 0x1)
	tAliceNet.CheckAndParsPayload(pld, tAnotherPeerID)

	pld2 := shouldPublishPayloadWithThisType(t, tAliceNet, payload.PayloadTypeAleyk)
	pldAleyk2 := pld2.(*payload.AleykPayload)

	assert.Equal(t, pldAleyk2.ResponseCode, payload.ResponseCodeOK)
	assert.Equal(t, tBobSync.peerSet.MaxClaimedHeight(), tAliceState.LastBlockHeight())

	p := tAliceSync.peerSet.GetPeer(tAnotherPeerID)
	assert.Equal(t, p.NodeVersion(), version.NodeVersion)
	assert.Equal(t, p.Moniker(), "kitty")
	assert.True(t, pub.EqualsTo(p.PublicKey()))
	assert.Equal(t, p.PeerID(), tAnotherPeerID)
	assert.Equal(t, p.Height(), 3)
	assert.Equal(t, p.InitialBlockDownload(), true)
}

func TestSendSalamPeerAhead(t *testing.T) {
	setup(t)

	_, pub, _ := crypto.GenerateTestKeyPair()
	claimedHeight := tAliceState.LastBlockHeight() + 5
	pld := payload.NewSalamPayload("kitty", pub, tAliceState.GenHash, claimedHeight, 0)
	tAliceNet.CheckAndParsPayload(pld, tAnotherPeerID)

	shouldPublishPayloadWithThisType(t, tAliceNet, payload.PayloadTypeAleyk)
	shouldPublishPayloadWithThisType(t, tAliceNet, payload.PayloadTypeLatestBlocksRequest)
	assert.Equal(t, tAliceSync.peerSet.MaxClaimedHeight(), claimedHeight)
}

func TestSendAleykPeerBehind(t *testing.T) {
	setup(t)

	_, pub, _ := crypto.GenerateTestKeyPair()
	pld := payload.NewAleykPayload(payload.ResponseCodeOK, "Welcome!", "kitty", pub, 1, 0)
	tAliceNet.CheckAndParsPayload(pld, tAnotherPeerID)

	shouldNotPublishMessageWithThisType(t, tAliceSync, payload.PayloadTypeLatestBlocksRequest)
}

func TestSendAleykPeerAhead(t *testing.T) {
	setup(t)

	_, pub, _ := crypto.GenerateTestKeyPair()
	claimedHeight := tAliceState.LastBlockHeight() + 5
	pld := payload.NewAleykPayload(payload.ResponseCodeOK, "Welcome!", "kitty", pub, claimedHeight, 0)
	tAliceNet.CheckAndParsPayload(pld, tAnotherPeerID)

	shouldPublishPayloadWithThisType(t, tAliceNet, payload.PayloadTypeLatestBlocksRequest)
	assert.Equal(t, tAliceSync.peerSet.MaxClaimedHeight(), claimedHeight)
}

func TestSendAleykPeerSameHeight(t *testing.T) {
	setup(t)

	_, pub, _ := crypto.GenerateTestKeyPair()
	claimedHeight := tAliceState.LastBlockHeight()
	pld := payload.NewAleykPayload(payload.ResponseCodeOK, "Welcome!", "kitty", pub, claimedHeight, 0)
	tAliceNet.CheckAndParsPayload(pld, tAnotherPeerID)

	shouldNotPublishMessageWithThisType(t, tAliceSync, payload.PayloadTypeLatestBlocksRequest)
}

// func TestIncreaseHeight(t *testing.T) {
// 	setup(t)

// 	_, pub, _ := crypto.GenerateTestKeyPair()
// 	msg1 := message.NewMessage(tAnotherPeerID, payload.NewSalamPayload("kitty", pub, tAliceState.GenesisHash(), 103, 0))
// 	tAliceSync.ParsMessage(msg1, tAnotherPeerID)
// 	assert.Equal(t, tAliceSync.peerSet.MaxClaimedHeight(), 103)

// 	msg2 := message.NewMessage(tAnotherPeerID, payload.NewAleykPayload(payload.ResponseCodeOK, "Welcome!", "kitty-2", pub, 104, 0))
// 	tAliceSync.ParsMessage(msg2, tAnotherPeerID)
// 	assert.Equal(t, tAliceSync.peerSet.MaxClaimedHeight(), 104)

// 	msg3 := message.NewMessage(tAnotherPeerID, payload.NewHeartBeatPayload(106, 0, crypto.GenerateTestHash()))
// 	tAliceSync.ParsMessage(msg3, tAnotherPeerID)
// 	assert.Equal(t, tAliceSync.peerSet.MaxClaimedHeight(), 105)
// }

func TestQueryTransaction(t *testing.T) {
	setup(t)

	trx1, _ := tx.GenerateTestBondTx()
	trx2, _ := tx.GenerateTestSendTx()
	trx3, _ := tx.GenerateTestSendTx()
	trx4, _ := tx.GenerateTestSendTx()

	// Alice has trx1 in his cache
	tAliceSync.cache.AddTransaction(trx1)
	tAliceSync.cache.AddTransaction(trx3)
	tAliceSync.cache.AddTransaction(trx4)
	tBobSync.cache.AddTransaction(trx2)
	pld := payload.NewQueryTransactionsPayload([]crypto.Hash{trx2.ID(), trx3.ID(), trx4.ID()})

	t.Run("Alice should not send query transaction message because she is not an active validator", func(t *testing.T) {
		tAliceBroadcastCh <- pld
		shouldNotPublishMessageWithThisType(t, tAliceSync, payload.PayloadTypeQueryTransactions)
	})

	t.Run("Bob should not process alice message because he is not an active validator", func(t *testing.T) {
		tBobNet.CheckAndParsPayload(pld, tAnotherPeerID)
		shouldNotPublishMessageWithThisType(t, tBobSync, payload.PayloadTypeTransactions)
	})

	joinBobToTheSet(t)

	t.Run("Bob should not process alice message because she is not an active validator", func(t *testing.T) {
		tBobNet.CheckAndParsPayload(pld, tAnotherPeerID)
		shouldNotPublishMessageWithThisType(t, tBobSync, payload.PayloadTypeTransactions)
	})

	joinAliceToTheSet(t)

	t.Run("Alice sends query transaction message", func(t *testing.T) {
		tAliceBroadcastCh <- pld
		shouldPublishPayloadWithThisType(t, tAliceNet, payload.PayloadTypeQueryTransactions)
	})

	t.Run("Alice sends query transaction message, but she has it in the cache", func(t *testing.T) {
		tAliceBroadcastCh <- payload.NewQueryTransactionsPayload([]crypto.Hash{trx1.ID()})
		shouldNotPublishMessageWithThisType(t, tAliceSync, payload.PayloadTypeQueryTransactions)
	})

	t.Run("Bob processes alice message", func(t *testing.T) {
		tBobNet.CheckAndParsPayload(pld, tAnotherPeerID)
		shouldPublishPayloadWithThisType(t, tBobNet, payload.PayloadTypeTransactions)
	})
}

func TestQueryProposal(t *testing.T) {
	setup(t)

	consensusHeight := tAliceState.LastBlockHeight() + 1
	p1, _ := proposal.GenerateTestProposal(consensusHeight, 0)
	p2, _ := proposal.GenerateTestProposal(consensusHeight, 1)

	tAliceSync.cache.AddProposal(p1)
	tBobConsensus.SetProposal(p2)
	pld := payload.NewQueryProposalPayload(consensusHeight, 1)

	t.Run("Alice should not send query proposal message because she is not an active validator", func(t *testing.T) {
		tAliceBroadcastCh <- pld
		shouldNotPublishMessageWithThisType(t, tAliceSync, payload.PayloadTypeQueryProposal)
	})

	t.Run("Bob should not process alice message because he is not an active validator", func(t *testing.T) {
		tBobNet.CheckAndParsPayload(pld, tAnotherPeerID)
		shouldNotPublishMessageWithThisType(t, tBobSync, payload.PayloadTypeTransactions)
	})

	joinBobToTheSet(t)

	t.Run("Bob should not process alice message because she is not an active validator", func(t *testing.T) {
		tBobNet.CheckAndParsPayload(pld, tAnotherPeerID)
		shouldNotPublishMessageWithThisType(t, tBobSync, payload.PayloadTypeTransactions)
	})

	joinAliceToTheSet(t)

	t.Run("Alice sends query transaction message", func(t *testing.T) {
		tAliceBroadcastCh <- pld
		shouldPublishPayloadWithThisType(t, tAliceNet, payload.PayloadTypeQueryProposal)
	})

	t.Run("Alice sends query transaction message, but she has it in her cache", func(t *testing.T) {
		tAliceBroadcastCh <- payload.NewQueryProposalPayload(consensusHeight, 0)
		shouldNotPublishMessageWithThisType(t, tAliceSync, payload.PayloadTypeQueryProposal)
	})

	t.Run("Bob processes alice message", func(t *testing.T) {
		tBobNet.CheckAndParsPayload(pld, tAnotherPeerID)
		shouldPublishPayloadWithThisType(t, tBobNet, payload.PayloadTypeProposal)
	})
}
func TestHeartbeatNotInSet(t *testing.T) {
	setup(t)

	// Alice is not in committee
	tAliceSync.broadcastHeartBeat()
	shouldPublishPayloadWithThisType(t, tAliceNet, payload.PayloadTypeHeartBeat)
	shouldNotPublishMessageWithThisType(t, tAliceSync, payload.PayloadTypeVote)

	joinAliceToTheSet(t)
	aliceH, _ := tAliceConsensus.HeightRound()
	v1, _ := vote.GenerateTestPrepareVote(aliceH, 0)
	tAliceConsensus.Votes = []*vote.Vote{v1}

	// Alice is in committee
	tAliceSync.broadcastHeartBeat()
	shouldPublishPayloadWithThisType(t, tAliceNet, payload.PayloadTypeHeartBeat)
	shouldPublishPayloadWithThisType(t, tAliceNet, payload.PayloadTypeVote)
}

func TestBlockAnnounceMessage(t *testing.T) {
	setup(t)

	joinBobToTheSet(t)

	t.Run("Bob should broadcast block announce message because he is an active validator", func(t *testing.T) {
		addMoreBlocksForBobAndAnnounceLastBlock(t, 1)

		shouldPublishPayloadWithThisType(t, tBobNet, payload.PayloadTypeBlockAnnounce)
		shouldNotPublishMessageWithThisType(t, tAliceSync, payload.PayloadTypeLatestBlocksRequest)

		aliceH, aliceR := tAliceConsensus.HeightRound()
		bobH, bobR := tAliceConsensus.HeightRound()
		assert.Equal(t, aliceH, bobH)
		assert.Equal(t, aliceR, bobR)
		assert.Equal(t, tAliceState.LastBlockHeight(), tBobState.LastBlockHeight())
	})
}

func TestNotActiveValidator(t *testing.T) {
	setup(t)

	t.Run("Alice is not an active validator, She should not send query proposal message", func(t *testing.T) {
		tAliceSync.queryProposal(1, 1)
		shouldNotPublishMessageWithThisType(t, tAliceSync, payload.PayloadTypeQueryProposal)
	})

	t.Run("Alice is not an active validator, She should not send query votes message", func(t *testing.T) {
		tAliceSync.queryVotes(1, 1)
		shouldNotPublishMessageWithThisType(t, tAliceSync, payload.PayloadTypeQueryVotes)
	})

	joinAliceToTheSet(t)

	t.Run("Alice is an active validator, She can send query proposal message", func(t *testing.T) {
		tAliceSync.queryProposal(1, 1)
		shouldPublishPayloadWithThisType(t, tAliceNet, payload.PayloadTypeQueryProposal)
	})

	t.Run("Alice is an active validator, She can send query votes message", func(t *testing.T) {
		tAliceSync.queryVotes(1, 1)
		shouldPublishPayloadWithThisType(t, tAliceNet, payload.PayloadTypeQueryVotes)
	})
}

func TestProposalToCache(t *testing.T) {
	setup(t)

	p, _ := proposal.GenerateTestProposal(106, 0)

	tAliceSync.BroadcastProposal(p)
	shouldPublishPayloadWithThisType(t, tAliceNet, payload.PayloadTypeProposal)
	assert.NotNil(t, tBobSync.cache.GetProposal(p.Height(), p.Round()))
}

func TestRequestForProposal(t *testing.T) {
	setup(t)

	joinAliceToTheSet(t)
	joinBobToTheSet(t)

	aliceH, _ := tAliceConsensus.HeightRound()
	assert.Equal(t, aliceH, tAliceState.LastBlockHeight()+1)

	t.Run("Alice and bob are in same height. Alice doesn't have have proposal. Bob ask for the proposal", func(t *testing.T) {
		tBobBroadcastCh <- payload.NewQueryProposalPayload(aliceH, 0)
		shouldPublishPayloadWithThisType(t, tBobNet, payload.PayloadTypeQueryProposal)

		// Alice doesn't respond
		shouldNotPublishMessageWithThisType(t, tAliceSync, payload.PayloadTypeProposal)
	})

	p1, _ := proposal.GenerateTestProposal(aliceH, 0)
	tAliceConsensus.SetProposal(p1)

	t.Run("Alice and bob are in same height. Alice has proposal. Bob ask for the proposal", func(t *testing.T) {
		tBobBroadcastCh <- payload.NewQueryProposalPayload(aliceH, 0)

		shouldPublishPayloadWithThisType(t, tBobNet, payload.PayloadTypeQueryProposal)
		shouldPublishPayloadWithThisType(t, tAliceNet, payload.PayloadTypeProposal)

		assert.Equal(t, tBobConsensus.Proposal.Hash(), p1.Hash())
	})

	t.Run("Alice and bob are in same height. Bob is in next round. Bob ask for the proposal", func(t *testing.T) {
		tBobBroadcastCh <- payload.NewQueryProposalPayload(aliceH, 1)
		shouldPublishPayloadWithThisType(t, tBobNet, payload.PayloadTypeQueryProposal)

		// Alice doesn't respond
		shouldNotPublishMessageWithThisType(t, tAliceSync, payload.PayloadTypeProposal)
	})

	p2, _ := proposal.GenerateTestProposal(aliceH, 1)
	tAliceConsensus.Proposal = p2
	tAliceConsensus.Round = 1

	t.Run("Alice and bob are in same height. Alice is in next round. Alice has proposal. Bob ask for the proposal", func(t *testing.T) {
		tBobBroadcastCh <- payload.NewQueryProposalPayload(aliceH, 1)

		shouldPublishPayloadWithThisType(t, tBobNet, payload.PayloadTypeQueryProposal)
		shouldPublishPayloadWithThisType(t, tAliceNet, payload.PayloadTypeProposal)

		assert.Equal(t, tBobConsensus.Proposal.Hash(), p2.Hash())
	})
}

func TestUpdateConsensus(t *testing.T) {
	setup(t)

	v, _ := vote.GenerateTestPrecommitVote(1, 1)
	p, _ := proposal.GenerateTestProposal(1, 1)

	tAliceSync.BroadcastVote(v)
	shouldPublishPayloadWithThisType(t, tAliceNet, payload.PayloadTypeVote)

	tAliceSync.BroadcastProposal(p)
	shouldPublishPayloadWithThisType(t, tAliceNet, payload.PayloadTypeProposal)

	assert.Equal(t, tBobConsensus.Votes[0].Hash(), v.Hash())
	assert.Equal(t, tBobConsensus.Proposal.Hash(), p.Hash())
}

func TestProcessQueryVote(t *testing.T) {
	setup(t)

	disableHeartbeat(t)
	joinAliceToTheSet(t)
	joinBobToTheSet(t)

	aliceH, _ := tAliceConsensus.HeightRound()
	v1, _ := vote.GenerateTestPrepareVote(aliceH, 0)
	v2, _ := vote.GenerateTestPrepareVote(aliceH, 1)
	tAliceConsensus.Votes = []*vote.Vote{v1, v2}

	t.Run("Alice and bob are in same height. Bob queries for votes, alice sends a random vote", func(t *testing.T) {
		tBobSync.BroadcastQueryVotes(aliceH, 1)
		shouldPublishPayloadWithThisType(t, tBobNet, payload.PayloadTypeQueryVotes)
		shouldPublishPayloadWithThisType(t, tAliceNet, payload.PayloadTypeVote)
	})
}

// func TestProcessHeartbeatForQueryProposal(t *testing.T) {
// 	setup(t)

// 	joinAliceToTheSet(t)
// 	joinBobToTheSet(t)

// 	aliceH, aliceR := tAliceConsensus.HeightRound()
// 	p, _ := proposal.GenerateTestProposal(aliceH, aliceR)
// 	tBobConsensus.SetProposal(p)

// 	t.Run("Alice Doesn't have proposal. She should query it.", func(t *testing.T) {
// 		tAliceSync.broadcastHeartBeat()
// 		shouldPublishPayloadWithThisType(t, tAliceNet,payload.PayloadTypeQueryProposal)
// 		shouldPublishPayloadWithThisType(t, tBobNet,payload.PayloadTypeProposal)
// 	})

// 	v1, _ := vote.GenerateTestPrepareVote(aliceH, 0)
// 	v2, _ := vote.GenerateTestPrepareVote(aliceH, 1)
// 	tAliceConsensus.Votes = []*vote.Vote{v1, v2}

// 	t.Run("Alice and bob are in same height.", func(t *testing.T) {
// 		tAliceSync.broadcastHeartBeat()
// 		shouldPublishPayloadWithThisType(t, tAliceNet,payload.PayloadTypeVote)
// 		shouldPublishPayloadWithThisType(t, tAliceNet,payload.PayloadTypeHeartBeat)
// 	})

// 	tAliceConsensus.Round = 1
// 	t.Run("Alice is in the next round. Bob isn't", func(t *testing.T) {
// 		tAliceSync.broadcastHeartBeat()
// 		shouldPublishPayloadWithThisType(t, tAliceNet,payload.PayloadTypeVote)
// 		shouldPublishPayloadWithThisType(t, tAliceNet,payload.PayloadTypeHeartBeat)

// 		shouldPublishPayloadWithThisType(t, tBobNet,payload.PayloadTypeQueryVotes)
// 	})
// }

// TODO:

// func TestAddBlockToCache(t *testing.T) {
// 	setup(t)
// 	disableHeartbeat(t)

// 	b1, trxs1 := block.GenerateTestBlock(nil, nil)
// 	b2, trxs2 := block.GenerateTestBlock(nil, nil)

// 	// Alice send a block to another peer, bob should cache it
// 	tAliceSync.BroadcastLatestBlocksResponse(payload.ResponseCodeMoreBlocks, tAnotherPeerID, 123, 1001, []*block.Block{b1}, trxs1, nil)
// 	shouldPublishPayloadWithThisType(t, tAliceNet,payload.PayloadTypeLatestBlocksResponse)
// 	assert.Equal(t, tBobSync.cache.GetBlock(1001).Hash(), b1.Hash())

// 	// Alice send a block to bob, bob should cache it
// 	tAliceSync.BroadcastLatestBlocksResponse(payload.ResponseCodeMoreBlocks, tBobPeerID, 123, 1002, []*block.Block{b2}, trxs2, nil)
// 	shouldPublishPayloadWithThisType(t, tAliceNet,payload.PayloadTypeLatestBlocksResponse)
// 	assert.Equal(t, tBobSync.cache.GetBlock(1002).Hash(), b2.Hash())
// }

// func TestAddTxToCache(t *testing.T) {
// 	setup(t)

// 	trx1, _ := tx.GenerateTestBondTx()

// 	// Alice send transaction to bob, bob should cache it
// 	tAliceSync.BroadcastTransactions([]*tx.Tx{trx1})
// 	shouldPublishPayloadWithThisType(t, tAliceNet,payload.PayloadTypeTransactions)
// 	assert.NotNil(t, tBobSync.cache.GetTransaction(trx1.ID()))
// 	assert.NotNil(t, tBobSync.state.PendingTx(trx1.ID()))
// }

// func TestRequestForBlocksNotVeryFar(t *testing.T) {
// 	setup(t)

// 	addMoreBlocksForBobAndAnnounceLastBlock(t, 15)

// 	tAliceSync.BroadcastLatestBlocksRequest(tBobPeerID)
// 	shouldPublishPayloadWithThisType(t, tAliceNet,payload.PayloadTypeLatestBlocksRequest)
// 	shouldPublishPayloadWithThisType(t, tBobNet,payload.PayloadTypeLatestBlocksResponse) // blocks 21-30
// 	shouldPublishPayloadWithThisType(t, tBobNet,payload.PayloadTypeLatestBlocksResponse) // blocks 31-35
// 	shouldPublishPayloadWithThisType(t, tBobNet,payload.PayloadTypeLatestBlocksResponse) // last commit + sync response
// }

// TODO:
// func TestPrepareLastBlock(t *testing.T) {
// 	setup(t)

// 	h := tAliceState.LastBlockHeight()
// 	b, _ := tAliceSync.prepareBlocksAndTransactions(h, 10)
// 	assert.Equal(t, len(b), 1)
// }

// func TestInvalidRangeForDownload(t *testing.T) {
// 	setup(t)

// 	t.Run("Bob is not target", func(t *testing.T) {
// 		pld := &payload.DownloadRequestPayload{
// 			SessionID: 1,
// 			Initiator: tAnotherPeerID,
// 			Target:    util.RandomPeerID(),
// 			From:      1000,
// 			To:        1001,
// 		}
// 		tBobSync.ProcessDownloadRequestPayload(pld)
// 		shouldNotPublishMessageWithThisType(t, tBobSync,payload.PayloadTypeDownloadResponse)
// 	})

// 	t.Run("Ask Bob to send big range of blocks", func(t *testing.T) {
// 		pld := &payload.DownloadRequestPayload{
// 			SessionID: 1,
// 			Initiator: tAnotherPeerID,
// 			Target:    tBobPeerID,
// 			From:      1000,
// 			To:        2000,
// 		}
// 		tBobSync.ProcessDownloadRequestPayload(pld)
// 		shouldNotPublishMessageWithThisType(t, tBobSync,payload.PayloadTypeDownloadResponse)
// 	})

// 	t.Run("Ask bob for the blocks that he doesn't have", func(t *testing.T) {
// 		pld := &payload.DownloadRequestPayload{
// 			SessionID: 1,
// 			Initiator: tAnotherPeerID,
// 			Target:    tBobPeerID,
// 			From:      1000,
// 			To:        1010,
// 		}
// 		tBobSync.ProcessDownloadRequestPayload(pld)
// 	})

// }

func TestDownloadBlocks(t *testing.T) {
	tAliceConfig.RequestBlockInterval = 30
	tBobConfig.RequestBlockInterval = 30
	setup(t)

	disableHeartbeat(t)

	// Clear alice store
	tAliceSync.cache.Clear()
	tAliceState.Store.Blocks = make(map[int]*block.Block)
	tAliceConsensus.Scheduled = false

	joinBobToTheSet(t)
	addMoreBlocksForBobAndAnnounceLastBlock(t, 80)
	shouldPublishPayloadWithThisType(t, tBobNet,payload.PayloadTypeBlockAnnounce)

	shouldPublishPayloadWithThisType(t, tAliceNet,payload.PayloadTypeDownloadRequest)
	shouldPublishPayloadWithThisType(t, tBobNet,payload.PayloadTypeDownloadResponse) // 1-10
	shouldPublishPayloadWithThisType(t, tBobNet,payload.PayloadTypeDownloadResponse) // 11-20
	shouldPublishPayloadWithThisType(t, tBobNet,payload.PayloadTypeDownloadResponse) // 21-31 (one extra block)
	shouldPublishPayloadWithThisType(t, tBobNet,payload.PayloadTypeDownloadResponse) // NoMoreBlock

	shouldPublishPayloadWithThisType(t, tAliceNet,payload.PayloadTypeDownloadRequest)
	shouldPublishPayloadWithThisType(t, tBobNet,payload.PayloadTypeDownloadResponse) // 40-49
	shouldPublishPayloadWithThisType(t, tBobNet,payload.PayloadTypeDownloadResponse) // 50-59
	shouldPublishPayloadWithThisType(t, tBobNet,payload.PayloadTypeDownloadResponse) // 60-70 (one extra block)
	shouldPublishPayloadWithThisType(t, tBobNet,payload.PayloadTypeDownloadResponse) // NoMoreBlock

	shouldPublishPayloadWithThisType(t, tAliceNet,payload.PayloadTypeDownloadRequest)
	shouldPublishPayloadWithThisType(t, tBobNet,payload.PayloadTypeDownloadResponse) // 61-70
	shouldPublishPayloadWithThisType(t, tBobNet,payload.PayloadTypeDownloadResponse) // 71-80
	shouldPublishPayloadWithThisType(t, tBobNet,payload.PayloadTypeDownloadResponse) // 81-91 (one extra block)
	shouldPublishPayloadWithThisType(t, tBobNet,payload.PayloadTypeDownloadResponse) // NoMoreBlock

	// Latest block requests
	shouldPublishPayloadWithThisType(t, tAliceNet,payload.PayloadTypeLatestBlocksRequest)
	shouldPublishPayloadWithThisType(t, tBobNet,payload.PayloadTypeLatestBlocksResponse) // 91-100
	shouldPublishPayloadWithThisType(t, tBobNet,payload.PayloadTypeLatestBlocksResponse) // 101-101
	shouldPublishPayloadWithThisType(t, tBobNet,payload.PayloadTypeLatestBlocksResponse) // Synced

	assert.Equal(t, tAliceState.LastBlockHash(), tBobState.LastBlockHash())
	assert.Equal(t, tAliceState.LastBlockHeight(), tBobState.LastBlockHeight())
	assert.False(t, tAliceSync.peerSet.HasAnyValidSession())
	assert.False(t, tBobSync.peerSet.HasAnyValidSession())
}

func TestSessionTimeout(t *testing.T) {
	tAliceConfig.SessionTimeout = 200 * time.Millisecond
	setup(t)

	p := tAliceSync.peerSet.MustGetPeer(tAnotherPeerID)
	p.UpdateInitialBlockDownload(true)
	p.UpdateHeight(1000)
	tAliceSync.peerSet.UpdateMaxClaimedHeight(1000)
	tAliceSync.sendBlocksRequestIfWeAreBehind()
	assert.True(t, tAliceSync.peerSet.HasAnyValidSession())
	time.Sleep(2 * tAliceConfig.SessionTimeout)
	assert.False(t, tAliceSync.peerSet.HasAnyValidSession())
}

func TestOneBlockBehind(t *testing.T) {
	setup(t)

	t.Run("Bob is not in the committee. Bob commits one block. Bob should broadcasts heartbeat. Alice should ask for the last block.", func(t *testing.T) {
		addMoreBlocksForBob(t, 1)

		shouldPublishPayloadWithThisType(t, tBobNet, payload.PayloadTypeHeartBeat)
		shouldPublishPayloadWithThisType(t, tAliceNet, payload.PayloadTypeLatestBlocksRequest)
		shouldPublishPayloadWithThisType(t, tBobNet, payload.PayloadTypeLatestBlocksResponse) // 22
		shouldPublishPayloadWithThisType(t, tBobNet, payload.PayloadTypeLatestBlocksResponse) // No more block

		assert.Equal(t, tAliceState.LastBlockHash(), tBobState.LastBlockHash())
		assert.Equal(t, tAliceState.LastBlockHeight(), tBobState.LastBlockHeight())
	})
}
*/
