package sync

import (
	"bytes"
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
	tConfig      *Config
	tState       *state.MockState
	tConsensus   *consensus.MockConsensus
	tNetwork     *network.MockNetwork
	tSync        *synchronizer
	tBroadcastCh chan payload.Payload
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
	tConfig = TestConfig()
	tConfig.Moniker = "Alice"
}

func setup(t *testing.T) {
	signer := bls.GenerateTestSigner()
	committee, _ := committee.GenerateTestCommittee()
	tState = state.MockingState(committee)
	tConsensus = consensus.MockingConsensus(tState)
	tBroadcastCh = make(chan payload.Payload, 1000)
	tNetwork = network.MockingNetwork(util.RandomPeerID())

	testAddBlocks(t, tState, 21)

	sync1, err := NewSynchronizer(tConfig,
		signer,
		tState,
		tConsensus,
		tNetwork,
		tBroadcastCh,
	)
	assert.NoError(t, err)
	tSync = sync1.(*synchronizer)

	assert.NoError(t, tSync.Start())
	shouldPublishPayloadWithThisType(t, tNetwork, payload.PayloadTypeHello)

	logger.Info("setup finished, running the tests", "name", t.Name())
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
			// Decode message again to check the payload type
			msg := new(message.Message)
			_, err := msg.Decode(bytes.NewReader(b.Data))
			require.NoError(t, err)
			assert.Equal(t, msg.Initiator, net.ID)

			// -----------
			// Check flags
			require.True(t, util.IsFlagSet(msg.Flags, message.MsgFlagNetworkLibP2P), "invalid flag: %v", msg)

			if b.Target == nil {
				require.True(t, util.IsFlagSet(msg.Flags, message.MsgFlagBroadcasted), "invalid flag: %v", msg)
			} else {
				require.False(t, util.IsFlagSet(msg.Flags, message.MsgFlagBroadcasted), "invalid flag: %v", msg)
			}

			if msg.Payload.Type() == payload.PayloadTypeHello {
				require.True(t, util.IsFlagSet(msg.Flags, message.MsgFlagHelloMessage), "invalid flag: %v", msg)
			} else {
				require.False(t, util.IsFlagSet(msg.Flags, message.MsgFlagHelloMessage), "invalid flag: %v", msg)
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
			// Decode message again to check the payload type
			msg := new(message.Message)
			_, err := msg.Decode(bytes.NewReader(b.Data))
			require.NoError(t, err)
			assert.NotEqual(t, msg.Payload.Type(), payloadType)
		}
	}
}

func testReceiveingNewMessage(sync *synchronizer, pld payload.Payload, from peer.ID) error {
	msg := message.NewMessage(from, pld)
	return sync.processIncomingMessage(msg)
}

func testAddBlocks(t *testing.T, state *state.MockState, count int) {
	lastBlockHash := state.LastBlockHash()
	for i := 0; i < count; i++ {
		b, trxs := block.GenerateTestBlock(nil, &lastBlockHash)
		c := block.GenerateTestCertificate(b.Hash())
		lastBlockHash = b.Hash()

		state.AddBlock(state.LastBlockHeight()+1, b, trxs)
		state.LastBlockCertificate = c
	}
	assert.Equal(t, lastBlockHash, state.LastBlockHash())
}

func testAddPeer(t *testing.T, pub crypto.PublicKey, pid peer.ID) *peerset.Peer {
	p := tSync.peerSet.MustGetPeer(pid)
	require.NotNil(t, p)
	p.UpdateMoniker("test")
	p.UpdatePublicKey(pub)
	p.UpdateStatus(peerset.StatusCodeKnown)

	return p
}

func testAddPeerToCommittee(t *testing.T, pid peer.ID, pub crypto.PublicKey) {
	if pub == nil {
		pub, _ = bls.GenerateTestKeyPair()
	}
	p := testAddPeer(t, pub, pid)
	val := validator.NewValidator(pub.(*bls.PublicKey), util.RandInt(0))
	val.UpdateLastJoinedHeight(tState.LastBlockHeight())
	assert.NoError(t, tState.Committee.Update(0, []*validator.Validator{val}))
	require.True(t, tState.Committee.Contains(p.Address()))
}

func checkPeerStatus(t *testing.T, pid peer.ID, code peerset.StatusCode) {
	peer := tSync.peerSet.GetPeer(pid)
	require.Equal(t, peer.Status(), code)
}

func TestStop(t *testing.T) {
	setup(t)
	// Should stop normally
	tSync.Stop()
}

func TestBroadcastInvalidMessage(t *testing.T) {
	setup(t)
	t.Run("Should not publish invalid messages", func(t *testing.T) {
		pld := payload.NewHeartBeatPayload(-1, -1, hash.GenerateTestHash())
		tBroadcastCh <- pld
		shouldNotPublishPayloadWithThisType(t, tNetwork, payload.PayloadTypeHeartBeat)
	})
}
