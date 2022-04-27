package network

import (
	"context"
	"encoding/hex"
	"fmt"

	lp2p "github.com/libp2p/go-libp2p"
	lp2pcore "github.com/libp2p/go-libp2p-core"
	lp2pcrypto "github.com/libp2p/go-libp2p-core/crypto"
	lp2phost "github.com/libp2p/go-libp2p-core/host"
	lp2peer "github.com/libp2p/go-libp2p-core/peer"
	lp2pps "github.com/libp2p/go-libp2p-pubsub"
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/logger"
	"github.com/zarbchain/zarb-go/util"
	"github.com/zarbchain/zarb-go/version"
)

var _ Network = &network{}

type network struct {
	ctx            context.Context
	cancel         func()
	config         *Config
	host           lp2phost.Host
	mdns           *mdnsService
	dht            *dhtService
	stream         *streamService
	gossip         *gossipService
	generalTopic   *lp2pps.Topic
	consensusTopic *lp2pps.Topic
	eventChannel   chan Event
	logger         *logger.Logger
}

func loadOrCreateKey(path string) (lp2pcrypto.PrivKey, error) {
	if util.PathExists(path) {
		h, err := util.ReadFile(path)
		if err != nil {
			return nil, err
		}
		bs, err := hex.DecodeString(string(h))
		if err != nil {
			return nil, err
		}
		key, err := lp2pcrypto.UnmarshalPrivateKey(bs)
		if err != nil {
			return nil, err
		}
		return key, nil
	}
	key, _, err := lp2pcrypto.GenerateEd25519Key(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to generate private key")
	}
	bs, err := lp2pcrypto.MarshalPrivateKey(key)
	if err != nil {
		return nil, err
	}
	h := hex.EncodeToString(bs)
	err = util.WriteFile(path, []byte(h))
	if err != nil {
		return nil, err
	}
	return key, nil
}

func NewNetwork(conf *Config) (Network, error) {
	nodeKey, err := loadOrCreateKey(conf.NodeKey)
	if err != nil {
		return nil, errors.Errorf(errors.ErrNetwork, err.Error())
	}

	opts := []lp2p.Option{
		lp2p.Identity(nodeKey),
		lp2p.ListenAddrStrings(conf.Listens...),
		lp2p.Ping(conf.EnablePing),
		lp2p.UserAgent(version.Agent()),
	}

	if conf.EnableNAT {
		opts = append(opts,
			lp2p.EnableNATService(),
			lp2p.NATPortMap())
	}

	if conf.EnableRelay {
		opts = append(opts,
			lp2p.EnableRelay())
	} else {
		opts = append(opts,
			lp2p.DisableRelay())
	}
	host, err := lp2p.New(opts...)
	if err != nil {
		return nil, errors.Errorf(errors.ErrNetwork, err.Error())
	}

	ctx, cancel := context.WithCancel(context.Background())

	n := &network{
		ctx:          ctx,
		cancel:       cancel,
		config:       conf,
		host:         host,
		eventChannel: make(chan Event, 100),
	}

	n.logger = logger.NewLogger("_network", n)

	if conf.EnableMdns {
		n.mdns = newMdnsService(ctx, n.host, n.logger)
	}

	if conf.EnableDHT {
		kadProtocolID := lp2pcore.ProtocolID(fmt.Sprintf("/%s/kad/v1", n.config.Name))
		n.dht = newDHTService(n.ctx, n.host, kadProtocolID, conf.Bootstrap, n.logger)
	}

	streamProtocolID := lp2pcore.ProtocolID(fmt.Sprintf("/%s/stream/v1", n.config.Name))
	n.stream = newStreamService(ctx, host, streamProtocolID, n.eventChannel, n.logger)

	n.gossip = newGossipService(ctx, host, n.eventChannel, n.logger)

	n.logger.Debug("network setup", "id", n.host.ID(), "address", conf.Listens)

	return n, nil
}

func (n *network) EventChannel() <-chan Event {
	return n.eventChannel
}

func (n *network) Start() error {
	if n.dht != nil {
		if err := n.dht.Start(); err != nil {
			return err
		}
	}

	if n.mdns != nil {
		if err := n.mdns.Start(); err != nil {
			return err
		}
	}
	n.gossip.Start()
	n.stream.Start()

	n.logger.Info("network started", "addr", n.host.Addrs())
	return nil
}

func (n *network) Stop() {
	n.cancel()

	if n.mdns != nil {
		n.mdns.Stop()
	}
	if n.dht != nil {
		n.dht.Stop()
	}
	n.gossip.Stop()
	n.stream.Stop()

	if err := n.host.Close(); err != nil {
		n.logger.Error("unable to close the network", "err", err)
	}
}

func (n *network) SelfID() lp2peer.ID {
	return n.host.ID()
}

func (n *network) SendTo(msg []byte, pid lp2pcore.PeerID) error {
	return n.stream.SendRequest(msg, pid)
}

func (n *network) Broadcast(msg []byte, topicID TopicID) error {
	n.logger.Debug("publishing new message", "topic", topicID)
	switch topicID {
	case TopicIDGeneral:
		if n.generalTopic == nil {
			return errors.Errorf(errors.ErrNetwork, "not subscribed to general topic")
		}
		return n.gossip.BroadcastMessage(msg, n.generalTopic)

	case TopicIDConsensus:
		if n.consensusTopic == nil {
			return errors.Errorf(errors.ErrNetwork, "not subscribed to consensus topic")
		}
		return n.gossip.BroadcastMessage(msg, n.consensusTopic)

	default:
		return errors.Errorf(errors.ErrNetwork, "invalid topic")
	}
}

func (n *network) JoinGeneralTopic() error {
	if n.generalTopic != nil {
		n.logger.Debug("already subscribed to general topic")
		return nil
	}
	topic, err := n.gossip.JoinTopic(n.TopicName("general"))
	if err != nil {
		return err
	}
	n.generalTopic = topic
	return nil
}

func (n *network) JoinConsensusTopic() error {
	if n.consensusTopic != nil {
		n.logger.Debug("already subscribed to consensus topic")
		return nil
	}
	topic, err := n.gossip.JoinTopic(n.TopicName("consensus"))
	if err != nil {
		return err
	}
	n.consensusTopic = topic
	return nil
}

func (n *network) TopicName(topic string) string {
	return fmt.Sprintf("/%s/topic/%s/v1", n.config.Name, topic)
}

func (n *network) CloseConnection(pid lp2peer.ID) {
	if err := n.host.Network().ClosePeer(pid); err != nil {
		n.logger.Warn("unable to close connection", "peer", pid)
	}
}

func (n *network) Fingerprint() string {
	return fmt.Sprintf("{%d}", n.NumConnectedPeers())
}

func (n *network) NumConnectedPeers() int {
	return len(n.host.Network().Peers())
}
