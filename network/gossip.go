package network

import (
	"context"
	"sync"

	lp2pps "github.com/libp2p/go-libp2p-pubsub"
	lp2phost "github.com/libp2p/go-libp2p/core/host"
	"github.com/pactus-project/pactus/util/errors"
	"github.com/pactus-project/pactus/util/logger"
)

type gossipService struct {
	ctx     context.Context
	wg      sync.WaitGroup
	host    lp2phost.Host
	pubsub  *lp2pps.PubSub
	topics  []*lp2pps.Topic
	subs    []*lp2pps.Subscription
	eventCh chan Event
	logger  *logger.SubLogger
}

func newGossipService(ctx context.Context, host lp2phost.Host, eventCh chan Event,
	logger *logger.SubLogger) *gossipService {
	pubsub, err := lp2pps.NewGossipSub(ctx, host)
	if err != nil {
		logger.Panic("unable to start Gossip service", "err", err)
		return nil
	}

	return &gossipService{
		ctx:     ctx,
		host:    host,
		pubsub:  pubsub,
		wg:      sync.WaitGroup{},
		eventCh: eventCh,
		logger:  logger,
	}
}

// BroadcastMessage broadcasts a message to the specified topic.
func (g *gossipService) BroadcastMessage(msg []byte, topic *lp2pps.Topic) error {
	return topic.Publish(g.ctx, msg)
}

// JoinTopic joins a topic with the given name.
// It creates a subscription to the topic and returns the joined topic.
func (g *gossipService) JoinTopic(name string) (*lp2pps.Topic, error) {
	topic, err := g.pubsub.Join(name)
	if err != nil {
		return nil, errors.Errorf(errors.ErrNetwork, err.Error())
	}
	sub, err := topic.Subscribe()
	if err != nil {
		return nil, errors.Errorf(errors.ErrNetwork, err.Error())
	}

	g.topics = append(g.topics, topic)
	g.subs = append(g.subs, sub)
	g.wg.Add(1)

	go func() {
		defer g.wg.Done()

		for {
			m, err := sub.Next(g.ctx)
			if err != nil {
				g.logger.Debug("readLoop error", "err", err)
				return
			}

			g.onReceiveMessage(m)
		}
	}()

	return topic, nil
}

// Start starts the gossip service.
func (g *gossipService) Start() {
}

// Stop stops the gossip service.
// It closes all the joined topics and cancels all the subscriptions.
func (g *gossipService) Stop() {
	for _, t := range g.topics {
		t.Close()
	}
	for _, s := range g.subs {
		s.Cancel()
	}

	g.wg.Wait()
}

func (g *gossipService) onReceiveMessage(m *lp2pps.Message) {
	// only forward messages delivered by others
	if m.ReceivedFrom == g.host.ID() {
		return
	}

	g.logger.Debug("receiving new gossip message", "from", m.GetFrom(), "received from", m.ReceivedFrom)
	event := &GossipMessage{
		Source: m.GetFrom(),
		From:   m.ReceivedFrom,
		Data:   m.Data,
	}

	g.eventCh <- event
}
