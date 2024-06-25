package sync

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/pactus-project/pactus/consensus"
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
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

func testConfig() *Config {
	return &Config{
		Moniker:             "test",
		SessionTimeout:      time.Second * 1,
		NodeNetwork:         true,
		BlockPerMessage:     11,
		MaxSessions:         8,
		LatestBlockInterval: 23,
		Firewall:            firewall.DefaultConfig(),
		LatestSupportingVer: DefaultConfig().LatestSupportingVer,
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

	consMgr, consMocks := consensus.MockingManager(ts, []*bls.ValidatorKey{valKeys[0], valKeys[1]})
	consMgr.MoveToNewHeight()

	broadcastCh := make(chan message.Message, 1000)
	mockNetwork := network.MockingNetwork(ts, ts.RandPeerID())

	syncInst, err := NewSynchronizer(config,
		valKeys,
		mockState,
		consMgr,
		mockNetwork,
		broadcastCh,
	)
	assert.NoError(t, err)
	sync := syncInst.(*synchronizer)

	td := &testData{
		TestSuite:   ts,
		config:      config,
		state:       mockState,
		consMgr:     consMgr,
		consMocks:   consMocks,
		network:     mockNetwork,
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

func (td *testData) shouldPublishMessageWithThisType(t *testing.T, msgType message.Type,
) *bundle.Bundle {
	t.Helper()

	return shouldPublishMessageWithThisType(t, td.network, msgType)
}

func shouldNotPublishMessageWithThisType(t *testing.T, net *network.MockNetwork, msgType message.Type) {
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

func (td *testData) shouldNotPublishMessageWithThisType(t *testing.T, msgType message.Type) {
	t.Helper()

	shouldNotPublishMessageWithThisType(t, td.network, msgType)
}

func (*testData) receivingNewMessage(sync *synchronizer, msg message.Message, from peer.ID) error {
	bdl := bundle.NewBundle(msg)
	bdl.Flags = util.SetFlag(bdl.Flags, bundle.BundleFlagCarrierLibP2P|bundle.BundleFlagNetworkMainnet)

	return sync.processIncomingBundle(bdl, from)
}

func (td *testData) addPeer(t *testing.T, s status.Status, services service.Services) peer.ID {
	t.Helper()

	pid := td.RandPeerID()
	pub, _ := td.RandBLSKeyPair()

	td.sync.peerSet.UpdateInfo(pid, t.Name(),
		version.NodeAgent.String(), []*bls.PublicKey{pub}, services)
	td.sync.peerSet.UpdateStatus(pid, s)

	return pid
}

func (td *testData) addValidatorToCommittee(t *testing.T, pub crypto.PublicKey) {
	t.Helper()

	if pub == nil {
		pub, _ = td.RandBLSKeyPair()
	}
	val := validator.NewValidator(pub.(*bls.PublicKey), td.RandInt32(1000))
	// Note: This may not be completely accurate, but it has no harm for testing purposes.
	val.UpdateLastSortitionHeight(td.state.TestCommittee.Proposer(0).LastSortitionHeight() + 1)
	td.state.TestStore.UpdateValidator(val)
	td.state.TestCommittee.Update(0, []*validator.Validator{val})
	require.True(t, td.state.TestCommittee.Contains(pub.(*bls.PublicKey).ValidatorAddress()))

	for _, cons := range td.consMocks {
		cons.SetActive(cons.ValKey.PublicKey().EqualsTo(pub))
	}
}

func (td *testData) checkPeerStatus(t *testing.T, pid peer.ID, code status.Status) {
	t.Helper()

	require.Equal(t, td.sync.peerSet.GetPeerStatus(pid), code)
}

func TestStop(t *testing.T) {
	td := setup(t, nil)

	// Should stop gracefully.
	td.sync.Stop()
}

func TestConnectEvent(t *testing.T) {
	conf := testConfig()
	conf.Firewall.BannedNets = []string{
		"84.247.0.0/24",
		"115.193.0.0/16",
		"240e:390:8a1:ae80:7dbc:64b6:e84c:d2bf/64",
	}

	td := setup(t, conf)

	pid := td.RandPeerID()
	ce := &network.ConnectEvent{
		PeerID:        pid,
		RemoteAddress: "/ip4/2.2.2.2/tcp/21888",
	}
	td.network.EventCh <- ce

	assert.Eventually(t, func() bool {
		p := td.sync.peerSet.GetPeer(pid)
		if p == nil {
			return false
		}
		assert.Equal(t, p.Address, "/ip4/2.2.2.2/tcp/21888")

		return p.Status == status.StatusConnected
	}, time.Second, 100*time.Millisecond)

	p1 := td.sync.peerSet.GetPeer(pid)
	assert.Equal(t, status.StatusConnected, p1.Status)

	// Receiving connect event for the banned address
	pid = td.RandPeerID()
	ce = &network.ConnectEvent{
		PeerID:        pid,
		RemoteAddress: "/ip4/115.193.2.1/tcp/21888",
	}
	td.network.EventCh <- ce

	assert.Eventually(t, func() bool {
		p := td.sync.peerSet.GetPeer(pid)
		if p == nil {
			return false
		}

		isBlocked := td.sync.firewall.IsBannedAddress(p.Address)

		if isBlocked {
			p.Status = status.StatusBanned
		}

		return isBlocked
	}, time.Second, 100*time.Millisecond)

	p2 := td.sync.peerSet.GetPeer(pid)
	assert.Equal(t, status.StatusBanned, p2.Status)
}

func TestDisconnectEvent(t *testing.T) {
	td := setup(t, nil)
	pid := td.RandPeerID()
	td.network.EventCh <- &network.DisconnectEvent{
		PeerID: pid,
	}

	assert.Eventually(t, func() bool {
		s := td.sync.peerSet.GetPeerStatus(pid)

		return s.IsDisconnected()
	}, time.Second, 100*time.Millisecond)
}

func TestProtocolsEvent(t *testing.T) {
	td := setup(t, nil)

	pid := td.RandPeerID()
	td.network.EventCh <- &network.ProtocolsEvents{
		PeerID:    pid,
		Protocols: []string{"protocol-1"},
	}
	td.shouldPublishMessageWithThisType(t, message.TypeHello)
}

func TestTestNetFlags(t *testing.T) {
	td := setup(t, nil)

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
		assert.NoError(t, td.receivingNewMessage(td.sync, baMsg, pid))

		td.shouldNotPublishMessageWithThisType(t, message.TypeBlocksRequest)
		td.network.IsClosed(pid)
	})

	t.Run("try to download blocks, but the peer is not a network node", func(t *testing.T) {
		td := setup(t, conf)

		pid := td.addPeer(t, status.StatusKnown, service.New(service.None))
		blk, cert := td.GenerateTestBlock(td.RandHeight())
		baMsg := message.NewBlockAnnounceMessage(blk, cert)
		assert.NoError(t, td.receivingNewMessage(td.sync, baMsg, pid))

		td.shouldNotPublishMessageWithThisType(t, message.TypeBlocksRequest)
		td.network.IsClosed(pid)
	})

	t.Run("try to download blocks and the peer is a network node", func(t *testing.T) {
		td := setup(t, conf)

		pid := td.addPeer(t, status.StatusKnown, service.New(service.Network))
		blk, cert := td.GenerateTestBlock(td.RandHeight())
		baMsg := message.NewBlockAnnounceMessage(blk, cert)
		assert.NoError(t, td.receivingNewMessage(td.sync, baMsg, pid))

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
		assert.NoError(t, td.receivingNewMessage(td.sync, msg, pid))

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

func TestBundleSequenceNo(t *testing.T) {
	td := setup(t, nil)

	msg := message.NewQueryProposalMessage(td.RandHeight(), td.RandRound(), td.RandValAddress())

	td.sync.broadcast(msg)
	bdl1 := td.shouldPublishMessageWithThisType(t, message.TypeQueryProposal)
	assert.Equal(t, 0, bdl1.SequenceNo)

	// Sending the same message again
	td.sync.broadcast(msg)
	bdl2 := td.shouldPublishMessageWithThisType(t, message.TypeQueryProposal)
	assert.Equal(t, 1, bdl2.SequenceNo)
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
