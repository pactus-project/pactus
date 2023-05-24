package sync

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/pactus-project/pactus/consensus"
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/network"
	"github.com/pactus-project/pactus/state"
	"github.com/pactus-project/pactus/sync/bundle"
	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/pactus-project/pactus/sync/firewall"
	"github.com/pactus-project/pactus/sync/peerset"
	"github.com/pactus-project/pactus/types/validator"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/logger"
	"github.com/pactus-project/pactus/version"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	tConfig      *Config
	tState       *state.MockState
	tConsMgr     consensus.Manager
	tConsMocks   []*consensus.MockConsensus
	tNetwork     *network.MockNetwork
	tSync        *synchronizer
	tBroadcastCh chan message.Message
)

type OverrideFingerprint struct {
	sync Synchronizer
	name string
}

func init() {
	LatestBlockInterval = 20
	tConfig = testConfig()
	tConfig.Moniker = "Alice"
}

func (o *OverrideFingerprint) Fingerprint() string {
	return o.name + o.sync.Fingerprint()
}

func testConfig() *Config {
	return &Config{
		Moniker:         "test",
		HeartBeatTimer:  0, // Disabling heartbeat
		SessionTimeout:  time.Second * 1,
		NodeNetwork:     true,
		BlockPerMessage: 10,
		MaxOpenSessions: 4,
		CacheSize:       1000,
		Firewall:        firewall.DefaultConfig(),
	}
}

func setup(t *testing.T) {
	signers := []crypto.Signer{bls.GenerateTestSigner(), bls.GenerateTestSigner()}
	tState = state.MockingState()
	tConsMgr, tConsMocks = consensus.MockingManager(signers)
	tBroadcastCh = make(chan message.Message, 1000)
	tNetwork = network.MockingNetwork(network.TestRandomPeerID())

	testAddBlocks(t, tState, 21)

	sync1, err := NewSynchronizer(tConfig,
		signers,
		tState,
		tConsMgr,
		tNetwork,
		tBroadcastCh,
	)
	assert.NoError(t, err)
	tSync = sync1.(*synchronizer)

	assert.NoError(t, tSync.Start())
	shouldPublishMessageWithThisType(t, tNetwork, message.MessageTypeHello) // Alice key 1
	shouldPublishMessageWithThisType(t, tNetwork, message.MessageTypeHello) // Alice key 2

	logger.Info("setup finished, running the tests", "name", t.Name())
}

func shouldPublishMessageWithThisType(t *testing.T, net *network.MockNetwork, msgType message.Type) *bundle.Bundle {
	timeout := time.NewTimer(3 * time.Second)

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
			require.True(t, util.IsFlagSet(bdl.Flags, bundle.BundleFlagCarrierLibP2P), "invalid flag: %v", bdl)
			require.True(t, util.IsFlagSet(bdl.Flags, bundle.BundleFlagNetworkTestnet), "invalid flag: %v", bdl)

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

func testReceivingNewMessage(sync *synchronizer, msg message.Message, from peer.ID) error {
	bdl := bundle.NewBundle(from, msg)
	bdl.Flags = util.SetFlag(bdl.Flags, bundle.BundleFlagCarrierLibP2P|bundle.BundleFlagNetworkMainnet)
	return sync.processIncomingBundle(bdl)
}

func testAddBlocks(t *testing.T, state *state.MockState, count int) {
	h := state.LastBlockHeight()
	state.CommitTestBlocks(count)
	assert.Equal(t, h+uint32(count), state.LastBlockHeight())
}

func testAddPeer(t *testing.T, pub crypto.PublicKey, pid peer.ID) {
	tSync.peerSet.UpdatePeerInfo(pid, peerset.StatusCodeKnown, t.Name(), version.Agent(), pub.(*bls.PublicKey), false)
}

func testAddPeerToCommittee(t *testing.T, pid peer.ID, pub crypto.PublicKey) {
	if pub == nil {
		pub, _ = bls.GenerateTestKeyPair()
	}
	testAddPeer(t, pub, pid)
	val := validator.NewValidator(pub.(*bls.PublicKey), util.RandInt32(0))
	// Note: This may not be completely accurate, but it poses no harm for testing purposes.
	val.UpdateLastJoinedHeight(tState.TestCommittee.Proposer(0).LastJoinedHeight() + 1)
	tState.TestStore.UpdateValidator(val)
	tState.TestCommittee.Update(0, []*validator.Validator{val})
	require.True(t, tState.TestCommittee.Contains(pub.Address()))

	for _, cons := range tConsMocks {
		cons.SetActive(cons.Signer.PublicKey().EqualsTo(pub))
	}
}

func checkPeerStatus(t *testing.T, pid peer.ID, code peerset.StatusCode) {
	require.Equal(t, tSync.peerSet.GetPeer(pid).Status, code)
}

func TestStop(t *testing.T) {
	setup(t)
	// Should stop gracefully.
	tSync.Stop()
}

func TestTestNetFlags(t *testing.T) {
	setup(t)

	tState.TestParams.BlockVersion = 0x3f
	bdl := tSync.prepareBundle(message.NewHeartBeatMessage(1, 0, hash.GenerateTestHash()))
	require.False(t, util.IsFlagSet(bdl.Flags, bundle.BundleFlagNetworkMainnet), "invalid flag: %v", bdl)
	require.True(t, util.IsFlagSet(bdl.Flags, bundle.BundleFlagNetworkTestnet), "invalid flag: %v", bdl)
}
