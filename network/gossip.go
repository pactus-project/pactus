package network

import (
	"context"
	"fmt"
	"sync"
	"time"

	lp2pps "github.com/libp2p/go-libp2p-pubsub"
	lp2pcore "github.com/libp2p/go-libp2p/core"
	lp2phost "github.com/libp2p/go-libp2p/core/host"
	"github.com/pactus-project/pactus/util/flume"
	"github.com/pactus-project/pactus/util/logger"
)

type gossipService struct {
	ctx              context.Context
	wg               sync.WaitGroup
	host             lp2phost.Host
	pubsub           *lp2pps.PubSub
	topics           []*lp2pps.Topic
	subs             []*lp2pps.Subscription
	topicBlock       *lp2pps.Topic
	topicTransaction *lp2pps.Topic
	topicConsensus   *lp2pps.Topic
	networkName      string
	networkPipe      flume.Pipeline[Event]
	logger           *logger.SubLogger
}

func newGossipService(ctx context.Context, host lp2phost.Host, conf *Config,
	networkPipe flume.Pipeline[Event], log *logger.SubLogger,
) *gossipService {
	opts := []lp2pps.Option{
		lp2pps.WithFloodPublish(true),
		lp2pps.WithMessageSignaturePolicy(lp2pps.StrictNoSign),
		lp2pps.WithNoAuthor(),
		lp2pps.WithMessageIdFn(MessageIDFunc),
		lp2pps.WithSeenMessagesTTL(60 * time.Second),
	}

	if conf.IsBootstrapper {
		// enable Peer exchange on bootstrappers
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
		ctx:         ctx,
		networkName: conf.NetworkName,
		host:        host,
		pubsub:      pubsub,
		wg:          sync.WaitGroup{},
		networkPipe: networkPipe,
		logger:      log,
	}
}

// Broadcast broadcasts a message with the specified topic ID to the network.
func (g *gossipService) Broadcast(msg []byte, topicID TopicID) error {
	g.logger.Debug("publishing new message", "topic", topicID)

	switch topicID {
	case TopicIDUnspecified:
		return InvalidTopicError{TopicID: topicID}

	case TopicIDBlock:
		if g.topicBlock == nil {
			return NotSubscribedError{TopicID: topicID}
		}

		return g.publish(msg, g.topicBlock)

	case TopicIDTransaction:
		if g.topicTransaction == nil {
			return NotSubscribedError{TopicID: topicID}
		}

		return g.publish(msg, g.topicTransaction)

	case TopicIDConsensus:
		if g.topicConsensus == nil {
			return NotSubscribedError{TopicID: topicID}
		}

		return g.publish(msg, g.topicConsensus)

	default:
		return InvalidTopicError{TopicID: topicID}
	}
}

// publish publishes a message with the specified topic to the network.
func (g *gossipService) publish(msg []byte, topic *lp2pps.Topic) error {
	err := topic.Publish(g.ctx, msg)
	if err != nil {
		return LibP2PError{Err: err}
	}

	return nil
}

// JoinTopic joins to the topic with the given name and subscribes to receive topic messages.
func (g *gossipService) JoinTopic(topicID TopicID, evaluator PropagationEvaluator) error {
	switch topicID {
	case TopicIDUnspecified:
		return InvalidTopicError{TopicID: topicID}

	case TopicIDBlock:
		topic, err := g.joinTopic(topicID, evaluator)
		if err != nil {
			return err
		}
		g.topicBlock = topic

		return nil

	case TopicIDTransaction:
		topic, err := g.joinTopic(topicID, evaluator)
		if err != nil {
			return err
		}
		g.topicTransaction = topic

		return nil

	case TopicIDConsensus:
		topic, err := g.joinTopic(topicID, evaluator)
		if err != nil {
			return err
		}
		g.topicConsensus = topic

		return nil

	default:
		return InvalidTopicError{TopicID: topicID}
	}
}

func (g *gossipService) TopicName(topicID TopicID) string {
	return fmt.Sprintf("/%s/topic/%s/v1", g.networkName, topicID.String())
}

// joinTopic joins a given topic and registers a validator for it.
// If successful, it returns the topic and subscribes to it.
func (g *gossipService) joinTopic(topicID TopicID, evaluator PropagationEvaluator) (*lp2pps.Topic, error) {
	topicName := g.TopicName(topicID)
	topic, err := g.pubsub.Join(topicName)
	if err != nil {
		return nil, LibP2PError{Err: err}
	}

	sub, err := topic.Subscribe()
	if err != nil {
		return nil, LibP2PError{Err: err}
	}

	err = g.pubsub.RegisterTopicValidator(
		topicName, g.createValidator(topicID, evaluator))
	if err != nil {
		return nil, LibP2PError{Err: err}
	}

	g.topics = append(g.topics, topic)
	g.subs = append(g.subs, sub)
	g.wg.Add(1)

	go func() {
		defer g.wg.Done()

		for {
			lp2pMsg, err := sub.Next(g.ctx)
			if err != nil {
				g.logger.Debug("readLoop error", "error", err)

				return
			}

			msg := &GossipMessage{
				From:    lp2pMsg.ReceivedFrom,
				Data:    lp2pMsg.Data,
				TopicID: topicID,
			}

			g.onReceiveMessage(msg)
		}
	}()

	return topic, nil
}

func (g *gossipService) createValidator(topicID TopicID, evaluator PropagationEvaluator,
) func(context.Context, lp2pcore.PeerID, *lp2pps.Message) lp2pps.ValidationResult {
	return func(_ context.Context, peerId lp2pcore.PeerID, lp2pMsg *lp2pps.Message) lp2pps.ValidationResult {
		msg := &GossipMessage{
			From:    peerId,
			Data:    lp2pMsg.Data,
			TopicID: topicID,
		}

		// Automatically accept and broadcast messages originating from this node
		if msg.From == g.host.ID() {
			return lp2pps.ValidationAccept
		}

		switch evaluator(msg) {
		case Drop:
			g.logger.Debug("message dropped", "from", peerId, "topic", topicID)

			return lp2pps.ValidationIgnore

		case DropButConsume:
			g.logger.Debug("message dropped but consumed", "from", peerId, "topic", topicID)

			// Consume the message first
			g.onReceiveMessage(msg)

			return lp2pps.ValidationIgnore

		case Propagate:
			return lp2pps.ValidationAccept

		default:
			panic("unreachable")
		}
	}
}

// Start starts the gossip service.
func (*gossipService) Start() {
}

// Stop stops the gossip service.
// It closes all the joined topics and cancels all the subscriptions.
func (g *gossipService) Stop() {
	for _, t := range g.topics {
		_ = t.Close()
	}
	for _, s := range g.subs {
		s.Cancel()
	}

	g.wg.Wait()
}

func (g *gossipService) onReceiveMessage(msg *GossipMessage) {
	// only forward messages delivered by others
	if msg.From == g.host.ID() {
		return
	}

	g.logger.Debug("receiving new gossip message", "from", msg.From)
	g.networkPipe.Send(msg)
}
