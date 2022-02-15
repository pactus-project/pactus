package sync

import (
	"bytes"
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
	"github.com/zarbchain/zarb-go/crypto/bls"
	"github.com/zarbchain/zarb-go/crypto/hash"
	"github.com/zarbchain/zarb-go/logger"
	"github.com/zarbchain/zarb-go/network"
	"github.com/zarbchain/zarb-go/state"
	"github.com/zarbchain/zarb-go/sync/message"
	"github.com/zarbchain/zarb-go/sync/message/payload"
	"github.com/zarbchain/zarb-go/sync/peerset"
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
)

type OverrideFingerprint struct {
	sync Synchronizer
	name string
}

func (o *OverrideFingerprint) Fingerprint() string {
	return o.name + o.sync.Fingerprint()
}

func init() {
	LatestBlockInterval = 20
	logger.InitLogger(logger.TestConfig())
	tAliceConfig = TestConfig()
	tBobConfig = TestConfig()

	tAliceConfig.Moniker = "Alice"
	tBobConfig.Moniker = "Bob"
}

func setup(t *testing.T) {
	_, prv1 := bls.GenerateTestKeyPair()
	_, prv2 := bls.GenerateTestKeyPair()
	signer1 := crypto.NewSigner(prv1)
	signer2 := crypto.NewSigner(prv2)

	committee, _ := committee.GenerateTestCommittee()
	tAlicePeerID = util.RandomPeerID()
	tBobPeerID = util.RandomPeerID()
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
	prevBlockHash := hash.Hash{}
	for i := 0; i < 21; i++ {
		b, trxs := block.GenerateTestBlock(nil, &prevBlockHash)
		c := block.GenerateTestCertificate(b.Hash())
		prevBlockHash = b.Hash()

		tAliceState.AddBlock(i+1, b, trxs)
		tAliceState.LastBlockCertificate = c

		tBobState.AddBlock(i+1, b, trxs)
		tBobState.LastBlockCertificate = c
	}

	tAliceSync = &synchronizer{ctx: context.Background()}
	sync1, err := NewSynchronizer(tAliceConfig,
		signer1,
		tAliceState,
		tAliceConsensus,
		tAliceNet,
		tAliceBroadcastCh,
	)
	assert.NoError(t, err)
	tAliceSync = sync1.(*synchronizer)

	tBobSync = &synchronizer{ctx: context.Background()}
	sync2, err := NewSynchronizer(tBobConfig,
		signer2,
		tBobState,
		tBobConsensus,
		tBobNet,
		tBobBroadcastCh,
	)
	assert.NoError(t, err)
	tBobSync = sync2.(*synchronizer)

	// -------------------------------
	// For better logging when testing
	overrideLogger := func(sync *synchronizer, name string) {
		sync.logger = logger.NewLogger("_sync", &OverrideFingerprint{name: fmt.Sprintf("%s - %s: ", name, t.Name()), sync: sync})
	}

	overrideLogger(tAliceSync, "Alice")
	overrideLogger(tBobSync, "Bob")

	tAliceNet.AddAnotherNetwork(tBobNet)
	tBobNet.AddAnotherNetwork(tAliceNet)
	/// -------------------------------

	assert.NoError(t, tAliceSync.Start())
	assert.NoError(t, tBobSync.Start())

	shouldPublishPayloadWithThisType(t, tAliceNet, payload.PayloadTypeHello)
	shouldPublishPayloadWithThisType(t, tBobNet, payload.PayloadTypeHello)

	// Hello acknowledgments
	shouldPublishPayloadWithThisType(t, tAliceNet, payload.PayloadTypeHello)
	shouldPublishPayloadWithThisType(t, tBobNet, payload.PayloadTypeHello)

	assert.Equal(t, tAliceState.LastBlockHeight(), tBobState.LastBlockHeight())

	logger.Info("setup finished, start running the test", "name", t.Name())
}

func shouldPublishPayloadWithThisType(t *testing.T, net *network.MockNetwork, payloadType payload.Type) *message.Message {
	timeout := time.NewTimer(2 * time.Second)

	for {
		select {
		case <-timeout.C:
			require.NoError(t, fmt.Errorf("ShouldPublishPayloadWithThisType %v: Timeout, test: %v", payloadType, t.Name()))
			return nil
		case b := <-net.BroadcastCh:
			net.SendToOthers(b.Data, b.Target)
			// Re-Decode again to check the payload type
			msg := new(message.Message)
			_, err := msg.Decode(bytes.NewReader(b.Data))
			require.NoError(t, err)
			assert.Equal(t, msg.Initiator, net.ID)

			// -----------
			// Check flags
			require.True(t, util.IsFlagSet(msg.Flags, message.FlagNetworkLibP2P), "invalid flag: %v", msg)

			if b.Target == nil {
				require.True(t, util.IsFlagSet(msg.Flags, message.FlagBroadcasted), "invalid flag: %v", msg)
			} else {
				require.False(t, util.IsFlagSet(msg.Flags, message.FlagBroadcasted), "invalid flag: %v", msg)
			}

			if msg.Payload.Type() == payload.PayloadTypeHello {
				require.True(t, util.IsFlagSet(msg.Flags, message.FlagHelloMessage), "invalid flag: %v", msg)
			} else {
				require.False(t, util.IsFlagSet(msg.Flags, message.FlagHelloMessage), "invalid flag: %v", msg)
			}
			// -----------

			if msg.Payload.Type() == payloadType {
				logger.Info("shouldPublishPayloadWithThisType", "msg", msg, "type", payloadType.String())
				return msg
			}
		}
	}
}

func shouldNotPublishPayloadWithThisType(t *testing.T, net *network.MockNetwork, payloadType payload.Type) {
	timeout := time.NewTimer(300 * time.Millisecond)

	for {
		select {
		case <-timeout.C:
			return
		case b := <-net.BroadcastCh:
			net.SendToOthers(b.Data, b.Target)
			// Re-Decode again to check the payload type
			msg := new(message.Message)
			_, err := msg.Decode(bytes.NewReader(b.Data))
			require.NoError(t, err)
			assert.NotEqual(t, msg.Payload.Type(), payloadType)
		}
	}
}

func simulatingReceiveingNewMessage(t *testing.T, sync *synchronizer, pld payload.Payload, from peer.ID) error {
	msg := message.NewMessage(from, pld)
	return sync.processIncomingMessage(msg)
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
	assert.Equal(t, lastBlockHash, tBobState.LastBlockHash())
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
	require.NotNil(t, tAliceSync)
	require.NotNil(t, tBobSync)

	tAliceSync.heartBeatTicker.Stop()
	tBobSync.heartBeatTicker.Stop()
}

func joinAliceToCommittee(t *testing.T) {
	val := validator.NewValidator(tAliceSync.signer.PublicKey().(*bls.PublicKey), 4)
	val.UpdateLastJoinedHeight(tAliceState.LastBlockHeight())

	assert.NoError(t, tAliceState.Committee.Update(0, []*validator.Validator{val}))
}

func joinBobToCommittee(t *testing.T) {
	val := validator.NewValidator(tBobSync.signer.PublicKey().(*bls.PublicKey), 5)
	val.UpdateLastJoinedHeight(tBobState.LastBlockHeight())

	assert.NoError(t, tAliceState.Committee.Update(0, []*validator.Validator{val}))
}

func checkPeerStatus(t *testing.T, pid peer.ID, code peerset.StatusCode) {
	peer := tAliceSync.peerSet.GetPeer(pid)
	require.Equal(t, peer.Status(), code)
}

func TestAccessors(t *testing.T) {
	setup(t)

	assert.Equal(t, tAliceSync.SelfID(), tAlicePeerID)
	assert.Equal(t, len(tAliceSync.Peers()), 1)
}

func TestStop(t *testing.T) {
	setup(t)
	// Should stop normally
	tAliceSync.Stop()
	tBobSync.Stop()
}

func TestBroadcastInvalidMessage(t *testing.T) {
	setup(t)
	t.Run("Should not publish invalid messages", func(t *testing.T) {
		pld := payload.NewHeartBeatPayload(-1, -1, hash.GenerateTestHash())
		tAliceBroadcastCh <- pld
		shouldNotPublishPayloadWithThisType(t, tAliceNet, payload.PayloadTypeHeartBeat)
	})
}
