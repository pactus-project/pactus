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
	"github.com/zarbchain/zarb-go/sync/bundle"
	"github.com/zarbchain/zarb-go/sync/bundle/message"
	"github.com/zarbchain/zarb-go/sync/peerset"
	"github.com/zarbchain/zarb-go/util"
	"github.com/zarbchain/zarb-go/validator"
	"github.com/zarbchain/zarb-go/version"
)

var (
	tConfig      *Config
	tState       *state.MockState
	tConsensus   *consensus.MockConsensus
	tNetwork     *network.MockNetwork
	tSync        *synchronizer
	tBroadcastCh chan message.Message
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
	tBroadcastCh = make(chan message.Message, 1000)
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
	shouldPublishMessageWithThisType(t, tNetwork, message.MessageTypeHello)

	logger.Info("setup finished, running the tests", "name", t.Name())
}

func shouldPublishMessageWithThisType(t *testing.T, net *network.MockNetwork, msgType message.Type) *bundle.Bundle {
	timeout := time.NewTimer(2 * time.Second)

	for {
		select {
		case <-timeout.C:
			require.NoError(t, fmt.Errorf("shouldPublishMessageWithThisType %v: Timeout, test: %v", msgType, t.Name()))
			return nil
		case b := <-net.BroadcastCh:
			net.SendToOthers(b.Data, b.Target)
			// Decode message again to check the message type
			bdl := new(bundle.Bundle)
			_, err := bdl.Decode(bytes.NewReader(b.Data))
			require.NoError(t, err)
			assert.Equal(t, bdl.Initiator, net.ID)

			// -----------
			// Check flags
			require.True(t, util.IsFlagSet(bdl.Flags, bundle.BundleFlagNetworkLibP2P), "invalid flag: %v", bdl)

			if b.Target == nil {
				require.True(t, util.IsFlagSet(bdl.Flags, bundle.BundleFlagBroadcasted), "invalid flag: %v", bdl)
			} else {
				require.False(t, util.IsFlagSet(bdl.Flags, bundle.BundleFlagBroadcasted), "invalid flag: %v", bdl)
			}

			if bdl.Message.Type() == message.MessageTypeHello {
				require.True(t, util.IsFlagSet(bdl.Flags, bundle.BundleFlagHelloMessage), "invalid flag: %v", bdl)
			} else {
				require.False(t, util.IsFlagSet(bdl.Flags, bundle.BundleFlagHelloMessage), "invalid flag: %v", bdl)
			}
			// -----------

			if bdl.Message.Type() == msgType {
				logger.Info("shouldPublishMessageWithThisType", "bundle", bdl, "type", msgType.String())
				return bdl
			}
		}
	}
}

func shouldNotPublishMessageWithThisType(t *testing.T, net *network.MockNetwork, msgType message.Type) {
	timeout := time.NewTimer(300 * time.Millisecond)

	for {
		select {
		case <-timeout.C:
			return
		case b := <-net.BroadcastCh:
			// Decode message again to check the message type
			bdl := new(bundle.Bundle)
			_, err := bdl.Decode(bytes.NewReader(b.Data))
			require.NoError(t, err)
			assert.NotEqual(t, bdl.Message.Type(), msgType)
		}
	}
}

func testReceiveingNewMessage(sync *synchronizer, msg message.Message, from peer.ID) error {
	bdl := bundle.NewBundle(from, msg)
	return sync.processIncomingBundle(bdl)
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

func testAddPeer(t *testing.T, pub crypto.PublicKey, pid peer.ID) {
	tSync.peerSet.UpdatePeer(pid, peerset.StatusCodeKnown, t.Name(), version.Agent(), 0, pub.(*bls.PublicKey), false)
}

func testAddPeerToCommittee(t *testing.T, pid peer.ID, pub crypto.PublicKey) {
	if pub == nil {
		pub, _ = bls.GenerateTestKeyPair()
	}
	testAddPeer(t, pub, pid)
	val := validator.NewValidator(pub.(*bls.PublicKey), util.RandInt(0))
	val.UpdateLastJoinedHeight(tState.LastBlockHeight())
	assert.NoError(t, tState.Committee.Update(0, []*validator.Validator{val}))
	require.True(t, tState.Committee.Contains(pub.Address()))
}

func checkPeerStatus(t *testing.T, pid peer.ID, code peerset.StatusCode) {
	require.Equal(t, tSync.peerSet.GetPeer(pid).Status, code)
}

func TestStop(t *testing.T) {
	setup(t)
	// Should stop normally
	tSync.Stop()
}

func TestBroadcastInvalidMessage(t *testing.T) {
	setup(t)
	t.Run("Should not publish invalid messages", func(t *testing.T) {
		tBroadcastCh <- message.NewHeartBeatMessage(-1, -1, hash.GenerateTestHash())
		shouldNotPublishMessageWithThisType(t, tNetwork, message.MessageTypeHeartBeat)
	})
}
