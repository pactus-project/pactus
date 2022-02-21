package network

import (
	"context"
	syncer "sync"

	lp2phost "github.com/libp2p/go-libp2p-core/host"
	lp2pps "github.com/libp2p/go-libp2p-pubsub"
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/logger"
	"github.com/zarbchain/zarb-go/util"
)

type gossipService struct {
	ctx     context.Context
	wg      syncer.WaitGroup
	host    lp2phost.Host
	pubsub  *lp2pps.PubSub
	topics  []*lp2pps.Topic
	subs    []*lp2pps.Subscription
	eventCh chan Event
	logger  *logger.Logger
}

func newGossipService(ctx context.Context, host lp2phost.Host, eventCh chan Event, logger *logger.Logger) *gossipService {
	pubsub, err := lp2pps.NewGossipSub(ctx, host)
	if err != nil {
		logger.Panic("Unable to start Gossip service: %v", err)
		return nil
	}

	return &gossipService{
		ctx:     ctx,
		host:    host,
		pubsub:  pubsub,
		wg:      syncer.WaitGroup{},
		eventCh: eventCh,
		logger:  logger,
	}
}

func (g *gossipService) BroadcastMessage(msg []byte, topic *lp2pps.Topic) error {
	return topic.Publish(g.ctx, msg)
}

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
func (g *gossipService) Start() {
}
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

	g.logger.Debug("Receiving new gossip message", "from", util.FingerprintPeerID(m.GetFrom()), "received from", util.FingerprintPeerID(m.ReceivedFrom))
	event := &GossipMessage{
		Source: m.GetFrom(),
		From:   m.ReceivedFrom,
		Data:   m.Data,
	}

	g.eventCh <- event
}
