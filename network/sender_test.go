package network

import (
	"fmt"
	"sync"
	"testing"
	"time"

	lp2ppeer "github.com/libp2p/go-libp2p/core/peer"
	"github.com/pactus-project/gopkg/logger"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type sentMessage struct {
	pid lp2ppeer.ID
	msg string
}

// fakeSender records the sent messages. The send of a message of a peer in
// blockPeers waits until the matching channel is closed.
type fakeSender struct {
	lk         sync.Mutex
	sent       []sentMessage
	blockPeers map[lp2ppeer.ID]chan struct{}
}

func (f *fakeSender) send(msg []byte, pid lp2ppeer.ID) error {
	if gate, ok := f.blockPeers[pid]; ok {
		<-gate
	}

	f.lk.Lock()
	defer f.lk.Unlock()
	f.sent = append(f.sent, sentMessage{pid: pid, msg: string(msg)})

	return nil
}

func (f *fakeSender) messagesTo(pid lp2ppeer.ID) []string {
	f.lk.Lock()
	defer f.lk.Unlock()

	msgs := make([]string, 0, len(f.sent))
	for _, m := range f.sent {
		if m.pid == pid {
			msgs = append(msgs, m.msg)
		}
	}

	return msgs
}

func makeTestSender(t *testing.T, fake *fakeSender) *network {
	t.Helper()

	log := logger.NewSubLogger(t.Context(), "_network", nil)

	return &network{
		ctx:    t.Context(),
		logger: log,
		sender: newOrderedSender(t.Context(), fake.send, log),
	}
}

func TestSendToKeepsOrder(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	pid := ts.RandPeerID()
	gate := make(chan struct{})
	fake := &fakeSender{blockPeers: map[lp2ppeer.ID]chan struct{}{pid: gate}}
	net := makeTestSender(t, fake)

	const count = 50
	expected := make([]string, 0, count)
	for i := 0; i < count; i++ {
		msg := fmt.Sprintf("msg-%d", i)
		expected = append(expected, msg)
		net.SendTo([]byte(msg), pid)
	}

	// The first message is still being sent while the others are queued.
	close(gate)

	require.Eventually(t, func() bool {
		return len(fake.messagesTo(pid)) == count
	}, 10*time.Second, 5*time.Millisecond)
	assert.Equal(t, expected, fake.messagesTo(pid))
}

func TestSendToDoesNotBlockOtherPeers(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	slowPeer := ts.RandPeerID()
	fastPeer := ts.RandPeerID()
	gate := make(chan struct{})
	fake := &fakeSender{blockPeers: map[lp2ppeer.ID]chan struct{}{slowPeer: gate}}
	net := makeTestSender(t, fake)

	net.SendTo([]byte("slow-1"), slowPeer)
	net.SendTo([]byte("slow-2"), slowPeer)
	net.SendTo([]byte("fast-1"), fastPeer)
	net.SendTo([]byte("fast-2"), fastPeer)

	require.Eventually(t, func() bool {
		return len(fake.messagesTo(fastPeer)) == 2
	}, 10*time.Second, 5*time.Millisecond)
	assert.Equal(t, []string{"fast-1", "fast-2"}, fake.messagesTo(fastPeer))
	assert.Empty(t, fake.messagesTo(slowPeer))

	close(gate)
	require.Eventually(t, func() bool {
		return len(fake.messagesTo(slowPeer)) == 2
	}, 10*time.Second, 5*time.Millisecond)
	assert.Equal(t, []string{"slow-1", "slow-2"}, fake.messagesTo(slowPeer))
}

func TestSendToBoundsPendingMessages(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	pid := ts.RandPeerID()
	gate := make(chan struct{})
	fake := &fakeSender{blockPeers: map[lp2ppeer.ID]chan struct{}{pid: gate}}
	net := makeTestSender(t, fake)

	const extra = 10
	for i := 0; i < maxPendingMessages+extra; i++ {
		net.SendTo([]byte(fmt.Sprintf("msg-%d", i)), pid)
	}

	close(gate)
	require.Eventually(t, func() bool {
		return len(fake.messagesTo(pid)) == maxPendingMessages
	}, 10*time.Second, 5*time.Millisecond)

	// The messages over the limit are dropped; the others keep their order.
	msgs := fake.messagesTo(pid)
	assert.Equal(t, "msg-0", msgs[0])
	assert.Equal(t, fmt.Sprintf("msg-%d", maxPendingMessages-1), msgs[len(msgs)-1])
}

func TestSendToAfterQueueDrained(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	pid := ts.RandPeerID()
	fake := &fakeSender{}
	net := makeTestSender(t, fake)

	net.SendTo([]byte("first"), pid)
	require.Eventually(t, func() bool {
		return len(fake.messagesTo(pid)) == 1
	}, 10*time.Second, 5*time.Millisecond)

	// Once the worker is idle it exits; a new message starts a new one.
	require.Eventually(t, func() bool {
		net.sender.lk.Lock()
		defer net.sender.lk.Unlock()

		return len(net.sender.queues) == 0
	}, 10*time.Second, 5*time.Millisecond)

	net.SendTo([]byte("second"), pid)
	require.Eventually(t, func() bool {
		return len(fake.messagesTo(pid)) == 2
	}, 10*time.Second, 5*time.Millisecond)
	assert.Equal(t, []string{"first", "second"}, fake.messagesTo(pid))
}
