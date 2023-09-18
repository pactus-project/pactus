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
	"github.com/pactus-project/pactus/network"
	"github.com/pactus-project/pactus/state"
	"github.com/pactus-project/pactus/sync/bundle"
	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/pactus-project/pactus/sync/firewall"
	"github.com/pactus-project/pactus/sync/peerset"
	"github.com/pactus-project/pactus/sync/services"
	"github.com/pactus-project/pactus/types/validator"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/logger"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/pactus-project/pactus/version"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testData struct {
	*testsuite.TestSuite

	config      *Config
	state       *state.MockState
	consMgr     consensus.Manager
	consMocks   []*consensus.MockConsensus
	network     *network.MockNetwork
	sync        *synchronizer
	broadcastCh chan message.Message
}

type OverrideStringer struct {
	sync *synchronizer
	name string
}

func init() {
	LatestBlockInterval = 23
}

func (o *OverrideStringer) String() string {
	return o.name + o.sync.String()
}

func testConfig() *Config {
	return &Config{
		Moniker:         "test",
		SessionTimeout:  time.Second * 1,
		NodeNetwork:     true,
		BlockPerMessage: 11,
		CacheSize:       1000,
		Firewall:        firewall.DefaultConfig(),
	}
}

func setup(t *testing.T, config *Config) *testData {
	t.Helper()

	ts := testsuite.NewTestSuite(t)

	if config == nil {
		config = testConfig()
		config.Moniker = "Alice"
	}
	signers := []crypto.Signer{ts.RandSigner(), ts.RandSigner()}
	state := state.MockingState(ts)
	consMgr, consMocks := consensus.MockingManager(ts, signers)
	broadcastCh := make(chan message.Message, 1000)
	network := network.MockingNetwork(ts, ts.RandPeerID())

	Sync, err := NewSynchronizer(config,
		signers,
		state,
		consMgr,
		network,
		broadcastCh,
	)
	assert.NoError(t, err)
	sync := Sync.(*synchronizer)

	td := &testData{
		TestSuite:   ts,
		config:      config,
		state:       state,
		consMgr:     consMgr,
		consMocks:   consMocks,
		network:     network,
		sync:        sync,
		broadcastCh: broadcastCh,
	}

	assert.NoError(t, td.sync.Start())
	assert.Equal(t, td.sync.Moniker(), config.Moniker)

	logger.Info("setup finished, running the tests", "name", t.Name())

	return td
}

func shouldPublishMessageWithThisType(t *testing.T, net *network.MockNetwork, msgType message.Type) *bundle.Bundle {
	t.Helper()

	timeout := time.NewTimer(3 * time.Second)

	for {
		select {
		case <-timeout.C:
			require.NoError(t, fmt.Errorf("shouldPublishMessageWithThisType %v: Timeout, test: %v", msgType, t.Name()))
			return nil
		case b := <-net.PublishCh:
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

			if bdl.Message.Type() == message.TypeHello ||
				bdl.Message.Type() == message.TypeHelloAck {
				require.True(t, util.IsFlagSet(bdl.Flags, bundle.BundleFlagHandshaking), "invalid flag: %v", bdl)
			} else {
				require.False(t, util.IsFlagSet(bdl.Flags, bundle.BundleFlagHandshaking), "invalid flag: %v", bdl)
			}
			// -----------

			if bdl.Message.Type() == msgType {
				logger.Info("shouldPublishMessageWithThisType",
					"bundle", bdl, "type", msgType.String())
				return bdl
			}
		}
	}
}

func (td *testData) shouldPublishMessageWithThisType(t *testing.T, net *network.MockNetwork,
	msgType message.Type,
) *bundle.Bundle {
	t.Helper()

	return shouldPublishMessageWithThisType(t, net, msgType)
}

func (td *testData) shouldNotPublishMessageWithThisType(t *testing.T, net *network.MockNetwork, msgType message.Type) {
	t.Helper()

	timeout := time.NewTimer(3 * time.Millisecond)

	for {
		select {
		case <-timeout.C:
			return
		case b := <-net.PublishCh:
			// Decode message again to check the message type
			bdl := new(bundle.Bundle)
			_, err := bdl.Decode(bytes.NewReader(b.Data))
			require.NoError(t, err)
			assert.NotEqual(t, msgType, bdl.Message.Type(),
				"not expected %s", msgType)
		}
	}
}

func (td *testData) receivingNewMessage(sync *synchronizer, msg message.Message, from peer.ID) error {
	bdl := bundle.NewBundle(from, msg)
	bdl.Flags = util.SetFlag(bdl.Flags, bundle.BundleFlagCarrierLibP2P|bundle.BundleFlagNetworkMainnet)
	return sync.processIncomingBundle(bdl)
}

func (td *testData) addPeer(t *testing.T, pub crypto.PublicKey, pid peer.ID, services services.Services) {
	t.Helper()

	td.sync.peerSet.UpdateInfo(pid, t.Name(),
		version.Agent(), []*bls.PublicKey{pub.(*bls.PublicKey)}, services)
	td.sync.peerSet.UpdateStatus(pid, peerset.StatusCodeKnown)
}

func (td *testData) addPeerToCommittee(t *testing.T, pid peer.ID, pub crypto.PublicKey) {
	t.Helper()

	if pub == nil {
		pub, _ = td.RandBLSKeyPair()
	}
	td.addPeer(t, pub, pid, services.New(services.Network))
	val := validator.NewValidator(pub.(*bls.PublicKey), td.RandInt32(1000))
	// Note: This may not be completely accurate, but it poses no harm for testing purposes.
	val.UpdateLastSortitionHeight(td.state.TestCommittee.Proposer(0).LastSortitionHeight() + 1)
	td.state.TestStore.UpdateValidator(val)
	td.state.TestCommittee.Update(0, []*validator.Validator{val})
	require.True(t, td.state.TestCommittee.Contains(pub.Address()))

	for _, cons := range td.consMocks {
		cons.SetActive(cons.Signer.PublicKey().EqualsTo(pub))
	}
}

func (td *testData) checkPeerStatus(t *testing.T, pid peer.ID, code peerset.StatusCode) {
	t.Helper()

	require.Equal(t, td.sync.peerSet.GetPeer(pid).Status, code)
}

func TestStop(t *testing.T) {
	td := setup(t, nil)

	// Should stop gracefully.
	td.sync.Stop()
}

func TestConnectEvents(t *testing.T) {
	td := setup(t, nil)

	pid := td.RandPeerID()
	td.network.EventCh <- &network.ConnectEvent{
		PeerID: pid,
	}
	td.shouldPublishMessageWithThisType(t, td.network, message.TypeHello)
	assert.Equal(t, td.sync.peerSet.GetPeer(pid).Status, peerset.StatusCodeConnected)
}

func TestTestNetFlags(t *testing.T) {
	td := setup(t, nil)

	td.addPeerToCommittee(t, td.sync.SelfID(), td.sync.signers[0].PublicKey())
	td.state.TestParams.BlockVersion = 0x3f
	bdl := td.sync.prepareBundle(message.NewQueryProposalMessage(td.RandHeight(), td.RandRound()))
	require.False(t, util.IsFlagSet(bdl.Flags, bundle.BundleFlagNetworkMainnet), "invalid flag: %v", bdl)
	require.True(t, util.IsFlagSet(bdl.Flags, bundle.BundleFlagNetworkTestnet), "invalid flag: %v", bdl)
}

func TestDownload(t *testing.T) {
	td := setup(t, nil)

	ourBlockHeight := td.state.LastBlockHeight()
	// To make sure the peer is not synced,
	// we add a block in past (more than 2 hours)

	blk := td.GenerateTestBlock(nil)
	cert := td.GenerateTestCertificate()
	pid := td.RandPeerID()
	msg := message.NewBlockAnnounceMessage(ourBlockHeight+LatestBlockInterval+1, blk, cert)

	t.Run("try to download blocks, but the peer is not known", func(t *testing.T) {
		assert.NoError(t, td.receivingNewMessage(td.sync, msg, pid))

		td.shouldNotPublishMessageWithThisType(t, td.network, message.TypeBlocksRequest)
	})

	t.Run("try to download blocks, but the peer is not a network node", func(t *testing.T) {
		pub, _ := td.RandBLSKeyPair()
		td.addPeer(t, pub, pid, services.New(services.None))

		assert.NoError(t, td.receivingNewMessage(td.sync, msg, pid))

		td.shouldNotPublishMessageWithThisType(t, td.network, message.TypeBlocksRequest)
	})

	t.Run("try to download blocks and the peer is a network node", func(t *testing.T) {
		pub, _ := td.RandBLSKeyPair()
		td.addPeer(t, pub, pid, services.New(services.Network))

		assert.NoError(t, td.receivingNewMessage(td.sync, msg, pid))

		bdl := td.shouldPublishMessageWithThisType(t, td.network, message.TypeBlocksRequest)

		// Let's close the opened session
		td.sync.peerSet.CloseSession(bdl.Message.(*message.BlocksRequestMessage).SessionID)
	})

	t.Run("download request is rejected", func(t *testing.T) {
		session := td.sync.peerSet.OpenSession(pid)
		msg2 := message.NewBlocksResponseMessage(message.ResponseCodeRejected, t.Name(),
			session.SessionID(), 1, nil, nil)
		assert.NoError(t, td.receivingNewMessage(td.sync, msg2, pid))
		bdl := td.shouldPublishMessageWithThisType(t, td.network, message.TypeBlocksRequest)

		// Let's close the opened session
		td.sync.peerSet.CloseSession(bdl.Message.(*message.BlocksRequestMessage).SessionID)
	})

	t.Run("testing send failure", func(t *testing.T) {
		td.network.SendError = fmt.Errorf("send error")

		assert.NoError(t, td.receivingNewMessage(td.sync, msg, pid))

		td.shouldNotPublishMessageWithThisType(t, td.network, message.TypeBlocksRequest)
	})
}
