package sync

import (
	"testing"

	lp2pnetwork "github.com/libp2p/go-libp2p/core/network"
	"github.com/pactus-project/pactus/network"
	"github.com/pactus-project/pactus/sync/bundle"
	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/pactus-project/pactus/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// GossipMessage.From is libp2p's ReceivedFrom: the last hop (the neighbour that
// forwarded the message), not the author (the pubsub is configured with
// StrictNoSign + NoAuthor). The topic validators return Propagate without
// decoding the payload, so a stream-type bundle (e.g. BlocksRequest) published
// by one peer is relayed by honest nodes, and each receiver then closes the
// connection with the honest relay that delivered it.
func TestGossipOfStreamMessageDoesNotPunishTheForwarder(t *testing.T) {
	td := setup(t, nil)

	forwarder := td.RandPeerID() // honest neighbour relaying a message it did not author
	td.network.OtherNets[forwarder] = network.MockingNetwork(td.TestSuite, forwarder)
	td.connectPeer(forwarder, lp2pnetwork.DirInbound, false)

	// Stream-only message wrapped in a gossip envelope, with a trailing height
	// that passes the "expired message" topic filter.
	bdl := bundle.NewBundle(message.NewBlocksRequestMessage(1, 1, 1))
	bdl.Flags = util.SetFlag(bdl.Flags, bundle.BundleFlagCarrierLibP2P|bundle.BundleFlagNetworkMainnet)
	data, err := bdl.Encode()
	require.NoError(t, err)

	require.False(t, td.network.IsClosed(forwarder))

	td.sync.processGossipMessage(&network.GossipMessage{From: forwarder, Data: data})

	assert.False(t, td.network.IsClosed(forwarder),
		"the connection with the relaying peer must not be closed for a message it did not author")
}
