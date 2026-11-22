package network

import (
	"context"
	"fmt"
	"testing"
	"time"

	lp2p "github.com/libp2p/go-libp2p"
	lp2pps "github.com/libp2p/go-libp2p-pubsub"
	lp2pspb "github.com/libp2p/go-libp2p-pubsub/pb"
	lp2phost "github.com/libp2p/go-libp2p/core/host"
	lp2ppeer "github.com/libp2p/go-libp2p/core/peer"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/stretchr/testify/require"
)

// previousMessageID is the gossip identifier of the releases before the topic was added to it.
func previousMessageID(m *lp2pspb.Message) string {
	h := hash.CalcHash(m.Data)

	return string(h[:20])
}

type gossipNode struct {
	name  string
	host  lp2phost.Host
	topic *lp2pps.Topic
	sub   *lp2pps.Subscription
}

// newGossipNode creates a node with the gossip options of the network, and the given identifier function.
func newGossipNode(t *testing.T, name string, idFn func(*lp2pspb.Message) string) *gossipNode {
	t.Helper()

	host, err := lp2p.New(lp2p.ListenAddrStrings("/ip4/127.0.0.1/tcp/0"), lp2p.DisableRelay())
	require.NoError(t, err)
	t.Cleanup(func() { _ = host.Close() })

	pubsub, err := lp2pps.NewGossipSub(t.Context(), host,
		lp2pps.WithFloodPublish(true),
		lp2pps.WithMessageSignaturePolicy(lp2pps.StrictNoSign),
		lp2pps.WithNoAuthor(),
		lp2pps.WithMessageIdFn(idFn),
		lp2pps.WithSeenMessagesTTL(60*time.Second),
	)
	require.NoError(t, err)

	topic, err := pubsub.Join("/test/topic/consensus/v1")
	require.NoError(t, err)
	sub, err := topic.Subscribe()
	require.NoError(t, err)

	return &gossipNode{name: name, host: host, topic: topic, sub: sub}
}

func connectGossipNodes(t *testing.T, left, right *gossipNode) {
	t.Helper()

	info := lp2ppeer.AddrInfo{ID: right.host.ID(), Addrs: right.host.Addrs()}
	require.NoError(t, left.host.Connect(t.Context(), info))
}

// receive returns the next message of the node that it didn't publish itself.
func (n *gossipNode) receive(t *testing.T) []byte {
	t.Helper()

	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	defer cancel()

	for {
		msg, err := n.sub.Next(ctx)
		require.NoError(t, err, "%s received nothing", n.name)
		if msg.ReceivedFrom != n.host.ID() {
			return msg.Data
		}
	}
}

// receivesNothing checks that the node gets no other message: a message is delivered once.
func (n *gossipNode) receivesNothing(t *testing.T) {
	t.Helper()

	ctx, cancel := context.WithTimeout(t.Context(), 2*time.Second)
	defer cancel()

	for {
		msg, err := n.sub.Next(ctx)
		if err != nil {
			return
		}
		require.Equal(t, n.host.ID(), msg.ReceivedFrom, "%s received %q twice", n.name, msg.Data)
	}
}

// During an upgrade, the nodes of the previous release and the new nodes are on the same
// network, and they compute different identifiers for the same message. The nodes are in a
// line, so that a message reaches the last node only if the middle node relays it: a node
// of one release relays the messages of the other, in both directions, and each node
// delivers a message once.
func TestGossipBetweenTheTwoMessageIdentifiers(t *testing.T) {
	tests := []struct {
		name  string
		idFns [3]func(*lp2pspb.Message) string
	}{
		{"previous - new - previous", [3]func(*lp2pspb.Message) string{previousMessageID, MessageIDFunc, previousMessageID}},
		{"new - previous - new", [3]func(*lp2pspb.Message) string{MessageIDFunc, previousMessageID, MessageIDFunc}},
		{"previous - previous - new", [3]func(*lp2pspb.Message) string{previousMessageID, previousMessageID, MessageIDFunc}},
		{"new - new - previous", [3]func(*lp2pspb.Message) string{MessageIDFunc, MessageIDFunc, previousMessageID}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nodes := make([]*gossipNode, 3)
			for i := range nodes {
				nodes[i] = newGossipNode(t, fmt.Sprintf("node %d", i), tt.idFns[i])
			}
			connectGossipNodes(t, nodes[0], nodes[1])
			connectGossipNodes(t, nodes[1], nodes[2])

			// The meshes are built at the heartbeats.
			time.Sleep(2 * time.Second)

			for round := range 20 {
				for sender := range nodes {
					data := []byte(fmt.Sprintf("message %d of node %d", round, sender))
					require.NoError(t, nodes[sender].topic.Publish(t.Context(), data))

					for receiver := range nodes {
						if receiver != sender {
							require.Equal(t, string(data), string(nodes[receiver].receive(t)),
								"round %d: from node %d to node %d", round, sender, receiver)
						}
					}
				}
			}

			// The gossip of the identifiers (IHAVE and IWANT) runs at the heartbeats:
			// it delivers no message a second time.
			for _, node := range nodes {
				node.receivesNothing(t)
			}
		})
	}
}
