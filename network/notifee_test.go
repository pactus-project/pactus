package network

import (
	"errors"
	"sync/atomic"
	"testing"
	"time"

	lp2pevent "github.com/libp2p/go-libp2p/core/event"
	lp2phost "github.com/libp2p/go-libp2p/core/host"
	lp2pnetwork "github.com/libp2p/go-libp2p/core/network"
	"github.com/pactus-project/gopkg/logger"
	"github.com/pactus-project/gopkg/pipeline"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// fakeEventSub is an event subscription that counts how many times its
// channel is requested.
type fakeEventSub struct {
	out      chan any
	outCalls atomic.Int64
}

func newFakeEventSub() *fakeEventSub {
	return &fakeEventSub{out: make(chan any, 1)}
}

func (s *fakeEventSub) Out() <-chan any {
	s.outCalls.Add(1)

	return s.out
}

func (s *fakeEventSub) Close() error {
	close(s.out)

	return nil
}

func (*fakeEventSub) Name() string {
	return "fake"
}

type failingEventBus struct {
	lp2pevent.Bus
}

func (*failingEventBus) Subscribe(any, ...lp2pevent.SubscriptionOpt) (lp2pevent.Subscription, error) {
	return nil, errors.New("subscribe failed")
}

type failingEventBusHost struct {
	lp2phost.Host
}

func (*failingEventBusHost) EventBus() lp2pevent.Bus {
	return &failingEventBus{}
}

func makeFakeNotifee(t *testing.T) (*NotifeeService, *fakeEventSub) {
	t.Helper()

	sub := newFakeEventSub()
	notifee := &NotifeeService{
		ctx:          t.Context(),
		lp2pEventSub: sub,
		networkPipe:  pipeline.New[Event](t.Context()),
		logger:       logger.NewSubLogger(t.Context(), "_network", nil),
	}

	return notifee, sub
}

func TestNotifeeStopEndsEventLoop(t *testing.T) {
	notifee, sub := makeFakeNotifee(t)

	notifee.Start()
	require.Eventually(t, func() bool {
		return sub.outCalls.Load() > 0
	}, 10*time.Second, 5*time.Millisecond)

	// Closing the subscription ends the loop, even though the context is alive.
	notifee.Stop()

	require.Eventually(t, func() bool {
		calls := sub.outCalls.Load()
		time.Sleep(20 * time.Millisecond)

		return sub.outCalls.Load() == calls
	}, 10*time.Second, 5*time.Millisecond)
}

func TestNotifeeReachability(t *testing.T) {
	notifee, sub := makeFakeNotifee(t)

	assert.Equal(t, lp2pnetwork.ReachabilityUnknown, notifee.Reachability())

	notifee.Start()
	sub.out <- lp2pevent.EvtLocalReachabilityChanged{Reachability: lp2pnetwork.ReachabilityPublic}

	require.Eventually(t, func() bool {
		return notifee.Reachability() == lp2pnetwork.ReachabilityPublic
	}, 10*time.Second, 5*time.Millisecond)

	notifee.Stop()
}

func TestNewNotifeeServiceSubscribeError(t *testing.T) {
	notifee, err := newNotifeeService(t.Context(), &failingEventBusHost{}, nil, nil, "", nil)
	require.Error(t, err)
	assert.Nil(t, notifee)
}
