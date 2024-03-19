package network

import (
	"context"
	"sync"
	"time"

	lp2pps "github.com/libp2p/go-libp2p-pubsub"
	lp2pcore "github.com/libp2p/go-libp2p/core"
	lp2phost "github.com/libp2p/go-libp2p/core/host"
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
	conf *Config, log *logger.SubLogger,
) *gossipService {
	opts := []lp2pps.Option{
		lp2pps.WithFloodPublish(true),
		lp2pps.WithMessageSignaturePolicy(lp2pps.StrictNoSign),
		lp2pps.WithNoAuthor(),
		lp2pps.WithMessageIdFn(MessageIDFunc),
		lp2pps.WithPeerOutboundQueueSize(600),
	}

	if conf.IsBootstrapper {
		// enable Peer eXchange on bootstrappers
		opts = append(opts, lp2pps.WithPeerExchange(true))
	}

	gsParams := lp2pps.DefaultGossipSubParams()
	if conf.IsBootstrapper {
		gsParams.Dhi = 12
		gsParams.D = 8
		gsParams.Dlo = 6
		gsParams.HeartbeatInterval = 700 * time.Millisecond
	}
	opts = append(opts, lp2pps.WithGossipSubParams(gsParams))

	pubsub, err := lp2pps.NewGossipSub(ctx, host, opts...)
	if err != nil {
		log.Panic("unable to start Gossip service", "error", err)

		return nil
	}

	return &gossipService{
		ctx:     ctx,
		host:    host,
		pubsub:  pubsub,
		wg:      sync.WaitGroup{},
		eventCh: eventCh,
		logger:  log,
	}
}

// BroadcastMessage broadcasts a message to the specified topic.
func (g *gossipService) BroadcastMessage(msg []byte, topic *lp2pps.Topic) error {
	err := topic.Publish(g.ctx, msg)
	if err != nil {
		return LibP2PError{Err: err}
	}

	return nil
}

// JoinTopic joins a topic with the given name.
// It creates a subscription to the topic and returns the joined topic.
func (g *gossipService) JoinTopic(name string, sp ShouldPropagate) (*lp2pps.Topic, error) {
	topic, err := g.pubsub.Join(name)
	if err != nil {
		return nil, LibP2PError{Err: err}
	}

	sub, err := topic.Subscribe()
	if err != nil {
		return nil, LibP2PError{Err: err}
	}

	err = g.pubsub.RegisterTopicValidator(name,
		func(_ context.Context, peerId lp2pcore.PeerID, m *lp2pps.Message) lp2pps.ValidationResult {
			msg := &GossipMessage{
				From: peerId,
				Data: m.Data,
			}
			if !sp(msg) {
				// Consume the message first
				g.onReceiveMessage(m)

				return lp2pps.ValidationIgnore
			}

			return lp2pps.ValidationAccept
		})
	if err != nil {
		return nil, LibP2PError{Err: err}
	}

	g.topics = append(g.topics, topic)
	g.subs = append(g.subs, sub)
	g.wg.Add(1)

	go func() {
		defer g.wg.Done()

		for {
			m, err := sub.Next(g.ctx)
			if err != nil {
				g.logger.Debug("readLoop error", "error", err)

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

	g.logger.Trace("receiving new gossip message",
		"source", m.GetFrom(), "from", m.ReceivedFrom)
	event := &GossipMessage{
		From: m.ReceivedFrom,
		Data: m.Data,
	}

	g.eventCh <- event
}
