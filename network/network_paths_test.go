package network

import (
	"errors"
	"io"
	"testing"
	"time"

	lp2pnetwork "github.com/libp2p/go-libp2p/core/network"
	"github.com/multiformats/go-multiaddr"
	"github.com/pactus-project/gopkg/logger"
	"github.com/pactus-project/gopkg/pipeline"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestErrorMessages(t *testing.T) {
	assert.Equal(t, "bad config", ConfigError{Reason: "bad config"}.Error())
	assert.Equal(t, "not subscribed to the 'block' topic",
		NotSubscribedError{TopicID: TopicIDBlock}.Error())
	assert.Equal(t, "invalid topic: 'unspecified'",
		InvalidTopicError{TopicID: TopicIDUnspecified}.Error())
	assert.Equal(t, "libp2p error: closed",
		LibP2PError{Err: errors.New("closed")}.Error())
}

func TestEventTypeString(t *testing.T) {
	tests := []struct {
		eventType EventType
		expected  string
	}{
		{EventTypeConnect, "connect"},
		{EventTypeDisconnect, "disconnect"},
		{EventTypeProtocols, "protocols"},
		{EventTypeGossip, "gossip-msg"},
		{EventTypeStream, "stream-msg"},
		{EventType(0), "invalid"},
		{EventType(99), "invalid"},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.expected, tt.eventType.String())
	}
}

// The public constructor creates a network that starts and stops.
func TestNewNetworkStartsAndStops(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	pipe := pipeline.New[Event](t.Context())
	net, err := NewNetwork(t.Context(), testConfig(), pipe)
	require.NoError(t, err)

	require.NoError(t, net.Start())
	assert.NotEmpty(t, net.SelfID())
	assert.NotEmpty(t, net.HostAddrs())
	assert.Zero(t, net.NumConnectedPeers())

	// Protecting and closing the connection of a peer that is not connected is harmless.
	pid := ts.RandPeerID()
	net.Protect(pid, "test")
	net.CloseConnection(pid)

	net.Stop()
}

// A message can be broadcast on a topic only once the topic is joined.
func TestBroadcastNeedsAJoinedTopic(t *testing.T) {
	net := makeTestNetwork(t, testConfig(), nil)
	msg := []byte("test-msg")

	err := net.gossip.Broadcast(msg, TopicIDUnspecified)
	require.ErrorIs(t, err, InvalidTopicError{TopicID: TopicIDUnspecified})

	err = net.gossip.Broadcast(msg, TopicID(99))
	require.ErrorIs(t, err, InvalidTopicError{TopicID: TopicID(99)})

	for _, topicID := range []TopicID{TopicIDBlock, TopicIDTransaction, TopicIDConsensus} {
		err = net.gossip.Broadcast(msg, topicID)
		require.ErrorIs(t, err, NotSubscribedError{TopicID: topicID})

		require.NoError(t, net.JoinTopic(topicID, alwaysPropagate))
		require.NoError(t, net.gossip.Broadcast(msg, topicID))
	}

	require.ErrorIs(t, net.JoinTopic(TopicIDUnspecified, alwaysPropagate),
		InvalidTopicError{TopicID: TopicIDUnspecified})
}

// Publishing on a topic that is closed returns an error, it doesn't panic.
func TestPublishOnAClosedTopic(t *testing.T) {
	net := makeTestNetwork(t, testConfig(), nil)
	require.NoError(t, net.JoinTopic(TopicIDBlock, alwaysPropagate))

	topic := net.gossip.topicBlock
	net.gossip.Stop()

	err := net.gossip.publish([]byte("test-msg"), topic)
	var libp2pErr LibP2PError
	require.ErrorAs(t, err, &libp2pErr)

	// The asynchronous broadcast only logs the error.
	net.Broadcast([]byte("test-msg"), TopicIDBlock)
	net.Broadcast([]byte("test-msg"), TopicIDUnspecified)
}

// Without a peer manager the gater has no count of the connections: it doesn't limit them.
func TestGaterWithoutPeerManager(t *testing.T) {
	ts := testsuite.NewTestSuite(t)
	log := logger.NewSubLogger(t.Context(), "_gater", nil)

	gater, err := NewConnectionGater(testConfig(), log)
	require.NoError(t, err)

	ma := multiaddr.StringCast("/ip4/8.8.8.8/tcp/1234")
	pid := ts.RandPeerID()

	assert.True(t, gater.InterceptPeerDial(pid))
	assert.True(t, gater.InterceptAddrDial(pid, ma))
	assert.True(t, gater.InterceptAccept(&mockConnMultiaddrs{remote: ma}))
	assert.True(t, gater.InterceptSecured(lp2pnetwork.DirInbound, pid, &mockConnMultiaddrs{remote: ma}))

	allowed, reason := gater.InterceptUpgraded(nil)
	assert.True(t, allowed)
	assert.Zero(t, reason)
}

// The mock network is what the other packages test with: two linked mocks deliver
// a direct message to its target only, and a broadcast to all the others.
func TestMockNetwork(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	netA := MockingNetwork(ts, ts.RandPeerID())
	netB := MockingNetwork(ts, ts.RandPeerID())
	netC := MockingNetwork(ts, ts.RandPeerID())

	require.NoError(t, netA.Start())
	require.NoError(t, netA.JoinTopic(TopicIDBlock, alwaysPropagate))
	netA.Protect(netB.SelfID(), "test")

	assert.True(t, netA.IsClosed(netB.SelfID()))

	netA.AddAnotherNetwork(netB, lp2pnetwork.DirOutbound)
	netA.AddAnotherNetwork(netC, lp2pnetwork.DirInbound)

	assert.False(t, netA.IsClosed(netB.SelfID()))
	assert.Equal(t, 2, netA.NumConnectedPeers())
	assert.Equal(t, 2, netA.NumInbound())
	assert.Equal(t, 2, netA.NumOutbound())

	// Linking a network sends the connect and the protocols events.
	eventsA := netA.EventPipe.UnsafeGetChannel()
	connect := receiveMockEvent(t, eventsA).(*ConnectEvent)
	assert.Equal(t, netB.SelfID(), connect.PeerID)
	assert.Equal(t, lp2pnetwork.DirOutbound, connect.Direction)
	protocols := receiveMockEvent(t, eventsA).(*ProtocolsEvents)
	assert.Equal(t, netB.SelfID(), protocols.PeerID)

	// A direct message reaches its target only.
	netA.SendTo([]byte("direct"), netB.SelfID())
	stream := receiveMockEvent(t, netB.EventPipe.UnsafeGetChannel()).(*StreamMessage)
	assert.Equal(t, netA.SelfID(), stream.From)
	data, err := io.ReadAll(stream.Reader)
	require.NoError(t, err)
	assert.Equal(t, "direct", string(data))
	noMockEvent(t, netC.EventPipe.UnsafeGetChannel())

	published := <-netA.PublishCh
	assert.Equal(t, "direct", string(published.Data))
	assert.Equal(t, netB.SelfID(), *published.Target)

	// A broadcast reaches all the others.
	netA.Broadcast([]byte("gossip"), TopicIDBlock)
	for _, other := range []*MockNetwork{netB, netC} {
		gossip := receiveMockEvent(t, other.EventPipe.UnsafeGetChannel()).(*GossipMessage)
		assert.Equal(t, netA.SelfID(), gossip.From)
		assert.Equal(t, "gossip", string(gossip.Data))
	}
	published = <-netA.PublishCh
	assert.Nil(t, published.Target)

	netA.CloseConnection(netB.SelfID())
	assert.True(t, netA.IsClosed(netB.SelfID()))
	assert.Equal(t, 1, netA.NumConnectedPeers())

	assert.Equal(t, "Unknown", netA.ReachabilityStatus())
	assert.Equal(t, []string{"localhost"}, netA.HostAddrs())
	assert.Equal(t, "pactus", netA.Name())
	assert.Equal(t, []string{"gossip"}, netA.Protocols())

	netA.Stop()
}

func receiveMockEvent(t *testing.T, events <-chan Event) Event {
	t.Helper()

	select {
	case evt := <-events:
		return evt
	case <-time.After(5 * time.Second):
		require.FailNow(t, "no event")

		return nil
	}
}

func noMockEvent(t *testing.T, events <-chan Event) {
	t.Helper()

	select {
	case evt := <-events:
		require.FailNow(t, "unexpected event", "type %v", evt.Type())
	case <-time.After(100 * time.Millisecond):
	}
}
