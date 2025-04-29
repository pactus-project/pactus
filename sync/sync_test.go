package sync

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/pactus-project/pactus/consensus"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/genesis"
	"github.com/pactus-project/pactus/network"
	"github.com/pactus-project/pactus/state"
	"github.com/pactus-project/pactus/sync/bundle"
	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/pactus-project/pactus/sync/firewall"
	"github.com/pactus-project/pactus/sync/peerset/peer"
	"github.com/pactus-project/pactus/sync/peerset/peer/service"
	"github.com/pactus-project/pactus/sync/peerset/peer/status"
	"github.com/pactus-project/pactus/types/validator"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/logger"
	"github.com/pactus-project/pactus/util/pipeline"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/pactus-project/pactus/version"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testData struct {
	*testsuite.TestSuite

	config    *Config
	state     *state.MockState
	consMgr   consensus.Manager
	consMocks []*consensus.MockConsensus
	network   *network.MockNetwork
	sync      *synchronizer
}

func testConfig() *Config {
	return &Config{
		Moniker:             "test",
		SessionTimeoutStr:   "1s",
		BlockPerMessage:     11,
		MaxSessions:         4,
		BlockPerSession:     23,
		PruneWindow:         13,
		Firewall:            firewall.DefaultConfig(),
		LatestSupportingVer: DefaultConfig().LatestSupportingVer,
		Services:            service.New(service.FullNode, service.PrunedNode),
	}
}

func setup(t *testing.T, config *Config) *testData {
	t.Helper()

	ts := testsuite.NewTestSuite(t)

	if config == nil {
		config = testConfig()
		config.Moniker = "Alice"
	}
	valKeys := []*bls.ValidatorKey{ts.RandValKey(), ts.RandValKey()}
	mockState := state.MockingState(ts)

	consMgr, consMocks := consensus.MockingManager(ts, mockState, []*bls.ValidatorKey{valKeys[0], valKeys[1]})
	consMgr.MoveToNewHeight()

	mockNetwork := network.MockingNetwork(ts, ts.RandPeerID())
	broadcastPipe := pipeline.MockingPipeline[message.Message]()

	syncInst, err := NewSynchronizer(config, valKeys,
		mockState, consMgr, mockNetwork, broadcastPipe, mockNetwork.EventPipe)
	assert.NoError(t, err)
	sync := syncInst.(*synchronizer)

	td := &testData{
		TestSuite: ts,
		config:    config,
		state:     mockState,
		consMgr:   consMgr,
		consMocks: consMocks,
		network:   mockNetwork,
		sync:      sync,
	}

	assert.NoError(t, td.sync.Start())
	assert.Equal(t, config.Moniker, td.sync.Moniker())
	assert.Equal(t, config.Services, td.sync.Services())

	logger.Info("setup finished, running the tests", "name", t.Name())

	return td
}

func shouldPublishMessageWithThisType(t *testing.T, net *network.MockNetwork, msgType message.Type) *bundle.Bundle {
	t.Helper()

	timer := time.NewTimer(3 * time.Second)

	for {
		select {
		case <-timer.C:
			require.NoError(t, fmt.Errorf("shouldPublishMessageWithThisType %v: Timeout, test: %v", msgType, t.Name()))

			return nil
		case data := <-net.PublishCh:
			// Decode message again to check the message type
			bdl := new(bundle.Bundle)
			_, err := bdl.Decode(bytes.NewReader(data.Data))
			require.NoError(t, err)

			// -----------
			// Check flags
			require.True(t, util.IsFlagSet(bdl.Flags, bundle.BundleFlagCarrierLibP2P), "invalid flag: %v", bdl)
			require.True(t, util.IsFlagSet(bdl.Flags, bundle.BundleFlagNetworkMainnet), "invalid flag: %v", bdl)

			if data.Target == nil {
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

			require.Equal(t, bdl.Message.Type(), msgType, "not expected %s", msgType)
			logger.Info("shouldPublishMessageWithThisType", "bundle", bdl, "type", msgType.String())

			return bdl
		}
	}
}

func (td *testData) shouldPublishMessageWithThisType(t *testing.T, msgType message.Type,
) *bundle.Bundle {
	t.Helper()

	return shouldPublishMessageWithThisType(t, td.network, msgType)
}

func shouldNotPublishMessageWithThisType(t *testing.T, net *network.MockNetwork, msgType message.Type) {
	t.Helper()

	timer := time.NewTimer(3 * time.Millisecond)

	for {
		select {
		case <-timer.C:
			return

		case b := <-net.PublishCh:
			// Decode message again to check the message type
			bdl := new(bundle.Bundle)
			_, err := bdl.Decode(bytes.NewReader(b.Data))
			require.NoError(t, err)
			assert.NotEqual(t, msgType, bdl.Message.Type(), "not expected %s", msgType)
		}
	}
}

func (td *testData) shouldNotPublishMessageWithThisType(t *testing.T, msgType message.Type) {
	t.Helper()

	shouldNotPublishMessageWithThisType(t, td.network, msgType)
}

func (*testData) receivingNewMessage(sync *synchronizer, msg message.Message, from peer.ID) {
	bdl := bundle.NewBundle(msg)
	bdl.Flags = util.SetFlag(bdl.Flags, bundle.BundleFlagCarrierLibP2P|bundle.BundleFlagNetworkMainnet)

	sync.processIncomingBundle(bdl, from)
}

func (td *testData) addPeer(t *testing.T, status status.Status, services service.Services) peer.ID {
	t.Helper()

	pid := td.RandPeerID()
	pub, _ := td.RandBLSKeyPair()

	td.sync.peerSet.UpdateInfo(pid, t.Name(),
		version.NodeAgent.String(), []*bls.PublicKey{pub}, services)
	td.sync.peerSet.UpdateStatus(pid, status)

	return pid
}

func (td *testData) addValidatorToCommittee(t *testing.T, pub *bls.PublicKey) {
	t.Helper()

	if pub == nil {
		pub, _ = td.RandBLSKeyPair()
	}
	val := td.GenerateTestValidator(testsuite.ValidatorWithPublicKey(pub))
	// Note: This may not be completely accurate, but it has no harm for testing purposes.
	val.UpdateLastSortitionHeight(td.state.TestCommittee.Proposer(0).LastSortitionHeight() + 1)
	td.state.TestStore.UpdateValidator(val)
	td.state.TestCommittee.Update(0, []*validator.Validator{val})
	require.True(t, td.state.TestCommittee.Contains(pub.ValidatorAddress()))

	for _, cons := range td.consMocks {
		cons.SetActive(cons.ValKey.PublicKey().EqualsTo(pub))
	}
}

func (td *testData) checkPeerStatus(t *testing.T, pid peer.ID, code status.Status) {
	t.Helper()

	require.Equal(t, code, td.sync.peerSet.GetPeerStatus(pid))
}

func TestStop(t *testing.T) {
	td := setup(t, nil)

	// Should stop gracefully.
	td.sync.Stop()
}

func TestConnectEvent(t *testing.T) {
	td := setup(t, nil)

	pid := td.RandPeerID()
	remoteAddr := "/ip4/2.2.2.2/tcp/21888"
	ce := &network.ConnectEvent{
		PeerID:        pid,
		RemoteAddress: remoteAddr,
		Direction:     "Inbound",
	}
	td.network.EventPipe.Send(ce)

	assert.Eventually(t, func() bool {
		return td.sync.peerSet.HasPeer(pid)
	}, time.Second, 100*time.Millisecond)

	p1 := td.sync.peerSet.GetPeer(pid)
	assert.Equal(t, status.StatusConnected, p1.Status)
	assert.Equal(t, remoteAddr, p1.Address)
	assert.Equal(t, "Inbound", p1.Direction)
}

func TestDisconnectEvent(t *testing.T) {
	td := setup(t, nil)
	pid := td.RandPeerID()
	de := &network.DisconnectEvent{
		PeerID: pid,
	}
	td.network.EventPipe.Send(de)

	assert.Eventually(t, func() bool {
		s := td.sync.peerSet.GetPeerStatus(pid)

		return s.IsDisconnected()
	}, time.Second, 100*time.Millisecond)
}

func TestProtocolsEvent(t *testing.T) {
	td := setup(t, nil)

	pid := td.RandPeerID()
	pe := &network.ProtocolsEvents{
		PeerID:    pid,
		Protocols: []string{"protocol-1"},
	}
	td.network.EventPipe.Send(pe)
	td.shouldPublishMessageWithThisType(t, message.TypeHello)
}

func TestTestNetFlags(t *testing.T) {
	td := setup(t, nil)

	td.state.TestGenesis = genesis.TestnetGenesis()
	td.addValidatorToCommittee(t, td.sync.valKeys[0].PublicKey())
	bdl := td.sync.prepareBundle(message.NewQueryProposalMessage(td.RandHeight(), td.RandRound(), td.RandValAddress()))

	require.False(t, util.IsFlagSet(bdl.Flags, bundle.BundleFlagNetworkMainnet), "invalid flag: %v", bdl)
	require.True(t, util.IsFlagSet(bdl.Flags, bundle.BundleFlagNetworkTestnet), "invalid flag: %v", bdl)
}

func TestDownload(t *testing.T) {
	conf := testConfig()
	// Let's not allow `GetRandomPeer` to disappoint us!
	conf.MaxSessions = 32

	t.Run("try to download blocks, but the peer is not known", func(t *testing.T) {
		td := setup(t, conf)

		pid := td.addPeer(t, status.StatusConnected, service.New(service.None))
		blk, cert := td.GenerateTestBlock(td.RandHeight())
		baMsg := message.NewBlockAnnounceMessage(blk, cert)
		td.receivingNewMessage(td.sync, baMsg, pid)

		td.shouldNotPublishMessageWithThisType(t, message.TypeBlocksRequest)
		td.network.IsClosed(pid)
	})

	t.Run("try to download blocks, but the peer is not a network node", func(t *testing.T) {
		td := setup(t, conf)

		pid := td.addPeer(t, status.StatusKnown, service.New(service.None))
		blk, cert := td.GenerateTestBlock(td.RandHeight())
		baMsg := message.NewBlockAnnounceMessage(blk, cert)
		td.receivingNewMessage(td.sync, baMsg, pid)

		td.shouldNotPublishMessageWithThisType(t, message.TypeBlocksRequest)
		td.network.IsClosed(pid)
	})

	t.Run("try to download blocks and the peer is a network node", func(t *testing.T) {
		td := setup(t, conf)

		pid := td.addPeer(t, status.StatusKnown, service.New(service.FullNode))
		blk, cert := td.GenerateTestBlock(td.RandHeight())
		baMsg := message.NewBlockAnnounceMessage(blk, cert)
		td.receivingNewMessage(td.sync, baMsg, pid)

		td.shouldPublishMessageWithThisType(t, message.TypeBlocksRequest)
	})

	t.Run("download request is rejected", func(t *testing.T) {
		td := setup(t, conf)

		pid := td.addPeer(t, status.StatusKnown, service.New(service.None))
		from := td.sync.stateHeight() + 1
		count := uint32(123)
		sid := td.sync.peerSet.OpenSession(pid, from, count)
		msg := message.NewBlocksResponseMessage(message.ResponseCodeRejected, t.Name(),
			sid, 1, nil, nil)
		td.receivingNewMessage(td.sync, msg, pid)

		assert.False(t, td.sync.peerSet.HasOpenSession(pid))
	})
}

func TestBroadcastBlockAnnounce(t *testing.T) {
	td := setup(t, nil)

	t.Run("Should announce the block", func(t *testing.T) {
		blk, cert := td.GenerateTestBlock(td.RandHeight())
		msg := message.NewBlockAnnounceMessage(blk, cert)

		td.sync.broadcast(msg)

		td.shouldPublishMessageWithThisType(t, message.TypeBlockAnnounce)
	})

	t.Run("Should NOT announce the block", func(t *testing.T) {
		blk, cert := td.GenerateTestBlock(td.RandHeight())
		msg := message.NewBlockAnnounceMessage(blk, cert)

		td.sync.cache.AddBlock(blk)
		td.sync.broadcast(msg)

		td.shouldNotPublishMessageWithThisType(t, message.TypeBlockAnnounce)
	})
}

func TestAllBlocksInCache(t *testing.T) {
	td := setup(t, nil)

	blk100, _ := td.GenerateTestBlock(100)
	blk101, _ := td.GenerateTestBlock(101)
	blk102, _ := td.GenerateTestBlock(102)

	td.sync.cache.AddBlock(blk100)
	td.sync.cache.AddBlock(blk101)
	td.sync.cache.AddBlock(blk102)

	res := td.sync.sendBlockRequestToRandomPeer(100, 3, true)
	assert.True(t, res)
}
