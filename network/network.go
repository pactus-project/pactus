package network

import (
	"context"
	"encoding/hex"
	"fmt"
	"time"

	lp2p "github.com/libp2p/go-libp2p"
	lp2pps "github.com/libp2p/go-libp2p-pubsub"
	lp2pcore "github.com/libp2p/go-libp2p/core"
	lp2pcrypto "github.com/libp2p/go-libp2p/core/crypto"
	lp2phost "github.com/libp2p/go-libp2p/core/host"
	lp2pmetrics "github.com/libp2p/go-libp2p/core/metrics"
	lp2ppeer "github.com/libp2p/go-libp2p/core/peer"
	lp2prcmgr "github.com/libp2p/go-libp2p/p2p/host/resource-manager"
	lp2pconnmgr "github.com/libp2p/go-libp2p/p2p/net/connmgr"
	"github.com/multiformats/go-multiaddr"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/logger"
	"github.com/pactus-project/pactus/version"
	"github.com/prometheus/client_golang/prometheus"
)

var _ Network = &network{}

type network struct {
	// TODO: Keeping context inside struct is bad practice:
	// Read more here: https://go.dev/blog/context-and-structs
	// We should remove it from here and pass it as first argument of functions
	// Adding these linter later:  contextcheck and containedctx
	ctx            context.Context
	cancel         func()
	config         *Config
	host           lp2phost.Host
	mdns           *mdnsService
	dht            *dhtService
	peerMgr        *peerMgr
	connGater      *ConnectionGater
	stream         *streamService
	gossip         *gossipService
	notifee        *NotifeeService
	generalTopic   *lp2pps.Topic
	consensusTopic *lp2pps.Topic
	eventChannel   chan Event
	logger         *logger.SubLogger
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
	log := logger.NewSubLogger("_network", nil)
	return newNetwork(conf, log, []lp2p.Option{})
}

func newNetwork(conf *Config, log *logger.SubLogger, opts []lp2p.Option) (*network, error) {
	networkKey, err := loadOrCreateKey(conf.NetworkKey)
	if err != nil {
		return nil, LibP2PError{Err: err}
	}

	rcMgrOpt := []lp2prcmgr.Option{}
	if conf.EnableMetrics {
		log.Info("metric enabled")
		str, err := lp2prcmgr.NewStatsTraceReporter()
		if err != nil {
			return nil, LibP2PError{Err: err}
		}

		// metrics for rcMgr
		lp2prcmgr.MustRegisterWith(prometheus.DefaultRegisterer)
		rcMgrOpt = append(rcMgrOpt, lp2prcmgr.WithTraceReporter(str))

		// metrics for lp2p
		bandwidthCounter := lp2pmetrics.NewBandwidthCounter()
		opts = append(opts, lp2p.BandwidthReporter(bandwidthCounter))
	} else {
		rcMgrOpt = append(rcMgrOpt, lp2prcmgr.WithMetricsDisabled())
		opts = append(opts, lp2p.DisableMetrics())
	}

	limit := MakeScalingLimitConfig(conf.MinConns, conf.MaxConns)
	resMgr, err := lp2prcmgr.NewResourceManager(
		lp2prcmgr.NewFixedLimiter(limit.AutoScale()),
		rcMgrOpt...,
	)
	if err != nil {
		return nil, LibP2PError{Err: err}
	}

	connMgr, err := lp2pconnmgr.NewConnManager(
		conf.MinConns, // Low Watermark
		conf.MaxConns, // High Watermark
		lp2pconnmgr.WithGracePeriod(time.Minute),
	)
	if err != nil {
		return nil, LibP2PError{Err: err}
	}

	opts = append(opts,
		lp2p.Identity(networkKey),
		lp2p.ListenAddrs(conf.ListenAddrs()...),
		lp2p.UserAgent(version.Agent()),
		lp2p.ResourceManager(resMgr),
		lp2p.ConnectionManager(connMgr),
	)

	if conf.EnableNATService {
		log.Info("Nat service enabled")
		opts = append(opts,
			lp2p.EnableNATService(),
		)
	}

	if conf.EnableUPnP {
		log.Info("UPnP enabled")
		opts = append(opts,
			lp2p.NATPortMap(),
		)
	}

	if conf.EnableRelay {
		log.Info("relay enabled", "relay addrs", conf.RelayAddrStrings)
		opts = append(opts,
			lp2p.EnableRelay(),
			lp2p.EnableAutoRelayWithStaticRelays(conf.RelayAddrInfos()),
			lp2p.EnableHolePunching(),
		)
	} else {
		log.Info("relay disabled")
		opts = append(opts,
			lp2p.DisableRelay(),
		)
	}

	// TODO: should include relay addresses
	privateSubnets := PrivateSubnets()
	privateFilters := SubnetsToFilters(privateSubnets, multiaddr.ActionDeny)
	publicAddrs := conf.PublicAddr()

	addrFactory := lp2p.AddrsFactory(func(mas []multiaddr.Multiaddr) []multiaddr.Multiaddr {
		addrs := []multiaddr.Multiaddr{}
		for _, addr := range mas {
			if conf.ForcePrivateNetwork || !privateFilters.AddrBlocked(addr) {
				addrs = append(addrs, addr)
			}
		}
		if publicAddrs != nil {
			addrs = append(addrs, publicAddrs)
		}
		// if len(addrs) == 0 {
		// for _, addr := range conf.RelayAddrStrings {
		// 	// To connect a peer over relay, we need a relay address.
		// 	// The format for the relay address is defined here:
		// 	// https://docs.libp2p.io/concepts/nat/circuit-relay/#relay-addresses
		// 	pid, _ := peer.IDFromPrivateKey(networkKey)
		// 	circuitAddr, _ := multiaddr.NewMultiaddr(fmt.Sprintf("%s/p2p-circuit/p2p/%s", addr, pid))
		// 	addrs = append(addrs, circuitAddr)
		// }

		// }
		return addrs
	})
	opts = append(opts, addrFactory)

	connGater, err := NewConnectionGater(conf, log)
	if err != nil {
		return nil, LibP2PError{Err: err}
	}
	opts = append(opts, lp2p.ConnectionGater(connGater))

	host, err := lp2p.New(opts...)
	if err != nil {
		return nil, LibP2PError{Err: err}
	}

	ctx, cancel := context.WithCancel(context.Background())

	n := &network{
		ctx:          ctx,
		cancel:       cancel,
		config:       conf,
		logger:       log,
		host:         host,
		connGater:    connGater,
		eventChannel: make(chan Event, 100),
	}

	log.SetObj(n)

	isBootstrapper := conf.IsBootstrapper(host.ID())
	kadProtocolID := lp2pcore.ProtocolID(fmt.Sprintf("/%s/gossip/v1", conf.NetworkName)) // TODO: better name?
	streamProtocolID := lp2pcore.ProtocolID(fmt.Sprintf("/%s/stream/v1", conf.NetworkName))

	if conf.EnableMdns {
		n.mdns = newMdnsService(ctx, n.host, n.logger)
	}
	n.dht = newDHTService(n.ctx, n.host, kadProtocolID, isBootstrapper, conf, n.logger)
	n.peerMgr = newPeerMgr(ctx, host, n.dht.kademlia, conf, n.logger)
	n.stream = newStreamService(ctx, n.host, streamProtocolID, n.eventChannel, n.logger)
	n.gossip = newGossipService(ctx, n.host, n.eventChannel, isBootstrapper, n.logger)
	n.notifee = newNotifeeService(ctx, n.host, n.eventChannel, n.peerMgr, streamProtocolID, isBootstrapper, n.logger)

	n.host.Network().Notify(n.notifee)
	n.connGater.SetPeerManager(n.peerMgr)

	n.logger.Info("network setup", "id", n.host.ID(),
		"name", conf.NetworkName,
		"address", conf.ListenAddrStrings,
		"bootstrapper", isBootstrapper)

	return n, nil
}

func (n *network) EventChannel() <-chan Event {
	return n.eventChannel
}

func (n *network) Start() error {
	if err := n.dht.Start(); err != nil {
		return LibP2PError{Err: err}
	}
	if n.mdns != nil {
		if err := n.mdns.Start(); err != nil {
			return LibP2PError{Err: err}
		}
	}
	n.gossip.Start()
	n.stream.Start()
	n.peerMgr.Start()
	n.notifee.Start()

	n.logger.Info("network started", "addr", n.host.Addrs(), "id", n.host.ID())
	return nil
}

func (n *network) Stop() {
	n.cancel()

	if n.mdns != nil {
		n.mdns.Stop()
	}

	n.dht.Stop()
	n.gossip.Stop()
	n.stream.Stop()
	n.peerMgr.Stop()
	n.notifee.Stop()

	if err := n.host.Close(); err != nil {
		n.logger.Error("unable to close the network", "error", err)
	}
}

func (n *network) SelfID() lp2ppeer.ID {
	return n.host.ID()
}

func (n *network) Protect(pid lp2pcore.PeerID, tag string) {
	n.host.ConnManager().Protect(pid, tag)
}

func (n *network) SendTo(msg []byte, pid lp2pcore.PeerID) error {
	n.logger.Trace("Sending new message", "to", pid)
	return n.stream.SendRequest(msg, pid)
}

func (n *network) Broadcast(msg []byte, topicID TopicID) error {
	n.logger.Trace("publishing new message", "topic", topicID)
	switch topicID {
	case TopicIDGeneral:
		if n.generalTopic == nil {
			return NotSubscribedError{TopicID: topicID}
		}
		return n.gossip.BroadcastMessage(msg, n.generalTopic)

	case TopicIDConsensus:
		if n.consensusTopic == nil {
			return NotSubscribedError{TopicID: topicID}
		}
		return n.gossip.BroadcastMessage(msg, n.consensusTopic)

	default:
		return InvalidTopicError{TopicID: topicID}
	}
}

func (n *network) JoinGeneralTopic(sp ShouldPropagate) error {
	if n.generalTopic != nil {
		n.logger.Debug("already subscribed to general topic")
		return nil
	}
	topic, err := n.gossip.JoinTopic(n.generalTopicName(), sp)
	if err != nil {
		return err
	}
	n.generalTopic = topic
	return nil
}

func (n *network) JoinConsensusTopic(sp ShouldPropagate) error {
	if n.consensusTopic != nil {
		n.logger.Debug("already subscribed to consensus topic")
		return nil
	}
	topic, err := n.gossip.JoinTopic(n.consensusTopicName(), sp)
	if err != nil {
		return err
	}
	n.consensusTopic = topic
	return nil
}

func (n *network) generalTopicName() string {
	return n.TopicName("general")
}

func (n *network) consensusTopicName() string {
	return n.TopicName("consensus")
}

func (n *network) TopicName(topic string) string {
	return fmt.Sprintf("/%s/topic/%s/v1", n.config.NetworkName, topic)
}

func (n *network) CloseConnection(pid lp2ppeer.ID) {
	if err := n.host.Network().ClosePeer(pid); err != nil {
		n.logger.Warn("unable to close connection", "peer", pid)
	}
}

func (n *network) String() string {
	return fmt.Sprintf("{%d}", n.NumConnectedPeers())
}

func (n *network) NumConnectedPeers() int {
	return len(n.host.Network().Peers())
}
