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
	nodeKey, err := loadOrCreateKey(conf.NodeKeyFile)
	if err != nil {
		return nil, errors.Errorf(errors.ErrNetwork, err.Error())
	}

	opts := []lp2p.Option{
		lp2p.Identity(nodeKey),
		lp2p.ListenAddrStrings(conf.ListenAddress...),
		lp2p.Ping(conf.EnablePing),
		lp2p.UserAgent("zarb-go-" + version.Version()),
	}

	if conf.EnableRelay {
		opts = append(opts,
			lp2p.EnableRelay())
	}
	host, err := lp2p.New(opts...)
	if err != nil {
		return nil, errors.Errorf(errors.ErrNetwork, err.Error())
	}

	ctx, cancel := context.WithCancel(context.Background())

	n := &network{
		ctx:    ctx,
		cancel: cancel,
		config: conf,
		host:   host,
	}

	n.logger = logger.NewLogger("_network", n)

	if conf.EnableMDNS {
		n.mdns = newMDNSService(ctx, n.host, n.logger)
	}

	if conf.EnableKademlia {
		kadProtocolID := lp2pcore.ProtocolID(fmt.Sprintf("/%s/kad/v1", n.config.Name))
		n.dht = newDHTService(n.ctx, n.host, kadProtocolID, conf.Bootstrap, n.logger)
	}

	streamProtocolID := lp2pcore.ProtocolID(fmt.Sprintf("/%s/stream/v1", n.config.Name))
	n.stream = newStreamService(ctx, host, streamProtocolID, n.logger)

	n.gossip = newGossipService(ctx, host, n.logger)

	n.logger.Debug("Network setup", "id", n.host.ID(), "address", conf.ListenAddress)

	return n, nil
}

func (n *network) SetCallback(callbackFn CallbackFn) {
	n.logger.Debug("Callback set")
	n.gossip.SetCallback(callbackFn)
	n.stream.SetCallback(callbackFn)
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

	n.logger.Info("Network started", "addr", n.host.Addrs())
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
		n.logger.Error("Unable to close the network", "err", err)
	}
}

func (n *network) SelfID() lp2peer.ID {
	return n.host.ID()
}

func (n *network) SendMessage(msg []byte, pid lp2pcore.PeerID) error {
	return n.stream.SendRequest(msg, pid)
}

func (n *network) BroadcastMessage(msg []byte, topicID TopicID) error {
	n.logger.Debug("Publishing new message", "topic", topicID)
	switch topicID {
	case GeneralTopic:
		if n.generalTopic != nil {
			return n.gossip.BroadcastMessage(msg, n.generalTopic)
		} else {
			n.logger.Warn("Not subscribed to general topic")
		}

	case ConsensusTopic:
		if n.consensusTopic != nil {
			return n.gossip.BroadcastMessage(msg, n.consensusTopic)
		} else {
			n.logger.Warn("Not subscribed to consensus topic")
		}

	default:
		n.logger.Warn("Invalid topic")
	}

	return nil
}

func (n *network) JoinGeneralTopic() error {
	if n.generalTopic != nil {
		n.logger.Debug("Already subscribed to general topic")
		return nil
	}
	name := fmt.Sprintf("/%s/topic/general/v1", n.config.Name)
	topic, err := n.gossip.JoinTopic(name)
	if err != nil {
		return err
	}
	n.generalTopic = topic
	return nil
}

func (n *network) JoinConsensusTopic() error {
	if n.consensusTopic != nil {
		n.logger.Debug("Already subscribed to consensus topic")
		return nil
	}
	name := fmt.Sprintf("/%s/topic/consensus/v1", n.config.Name)
	topic, err := n.gossip.JoinTopic(name)
	if err != nil {
		return err
	}
	n.generalTopic = topic
	return nil
}

func (n *network) CloseConnection(pid lp2peer.ID) {
	if err := n.host.Network().ClosePeer(pid); err != nil {
		n.logger.Warn("Unable to close connection", "peer", util.FingerprintPeerID(pid))
	}
}

func (n *network) Fingerprint() string {
	return fmt.Sprintf("{%d}", n.NumConnectedPeers())
}

func (n *network) NumConnectedPeers() int {
	return len(n.host.Network().Peers())
}
