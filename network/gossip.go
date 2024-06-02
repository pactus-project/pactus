package network

import (
	"context"
	"fmt"
	"sync"
	"time"

	lp2pps "github.com/libp2p/go-libp2p-pubsub"
	lp2pcore "github.com/libp2p/go-libp2p/core"
	lp2phost "github.com/libp2p/go-libp2p/core/host"
	"github.com/pactus-project/pactus/util/logger"
)

type gossipService struct {
	ctx                    context.Context
	wg                     sync.WaitGroup
	config                 *Config
	host                   lp2phost.Host
	pubsub                 *lp2pps.PubSub
	topics                 []*lp2pps.Topic
	subs                   []*lp2pps.Subscription
	generalTopicDeprecated *lp2pps.Topic
	blockTopic             *lp2pps.Topic
	transactionTopic       *lp2pps.Topic
	consensusTopic         *lp2pps.Topic
	eventCh                chan Event
	logger                 *logger.SubLogger
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
		config:  conf,
		pubsub:  pubsub,
		wg:      sync.WaitGroup{},
		eventCh: eventCh,
		logger:  log,
	}
}

func (g *gossipService) Broadcast(msg []byte, topicID TopicID) error {
	g.logger.Trace("publishing new message", "topic", topicID)
	switch topicID {
	case TopicIDGeneralDeprecated:
		if g.generalTopicDeprecated == nil {
			return NotSubscribedError{TopicID: topicID}
		}

		return g.BroadcastMessage(msg, g.generalTopicDeprecated)

	case TopicIDBlock:
		if g.blockTopic == nil {
			return NotSubscribedError{TopicID: topicID}
		}

		return g.BroadcastMessage(msg, g.blockTopic)

	case TopicIDTransaction:
		if g.transactionTopic == nil {
			return NotSubscribedError{TopicID: topicID}
		}

		return g.BroadcastMessage(msg, g.transactionTopic)

	case TopicIDConsensus:
		if g.consensusTopic == nil {
			return NotSubscribedError{TopicID: topicID}
		}

		return g.BroadcastMessage(msg, g.consensusTopic)

	default:
		return InvalidTopicError{TopicID: topicID}
	}
}

func (g *gossipService) JoinGeneralTopic(sp ShouldPropagate) error {
	if g.generalTopicDeprecated != nil {
		g.logger.Debug("already subscribed to general topic")

		return nil
	}
	topic, err := g.JoinTopic(g.generalTopicNameDeprecated(), sp)
	if err != nil {
		return err
	}
	g.generalTopicDeprecated = topic

	return nil
}

func (g *gossipService) JoinBlockTopic(sp ShouldPropagate) error {
	if g.blockTopic != nil {
		g.logger.Debug("already subscribed to general topic")

		return nil
	}
	topic, err := g.JoinTopic(g.blockTopicName(), sp)
	if err != nil {
		return err
	}
	g.blockTopic = topic

	return nil
}

func (g *gossipService) JoinTransactionTopic(sp ShouldPropagate) error {
	if g.transactionTopic != nil {
		g.logger.Debug("already subscribed to general topic")

		return nil
	}
	topic, err := g.JoinTopic(g.transactionTopicName(), sp)
	if err != nil {
		return err
	}
	g.transactionTopic = topic

	return nil
}

func (g *gossipService) JoinConsensusTopic(sp ShouldPropagate) error {
	if g.consensusTopic != nil {
		g.logger.Debug("already subscribed to consensus topic")

		return nil
	}
	topic, err := g.JoinTopic(g.consensusTopicName(), sp)
	if err != nil {
		return err
	}
	g.consensusTopic = topic

	return nil
}

func (g *gossipService) generalTopicNameDeprecated() string {
	return g.TopicName("general")
}

func (g *gossipService) blockTopicName() string {
	return g.TopicName("block")
}

func (g *gossipService) transactionTopicName() string {
	return g.TopicName("transaction")
}

func (g *gossipService) consensusTopicName() string {
	return g.TopicName("consensus")
}

func (g *gossipService) TopicName(topic string) string {
	return fmt.Sprintf("/%s/topic/%s/v1", g.config.NetworkName, topic)
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
				g.logger.Warn("message ignored", "from", peerId, "topic", name)

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
