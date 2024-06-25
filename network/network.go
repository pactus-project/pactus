package network

import (
	"context"
	"encoding/hex"
	"fmt"
	"time"

	lp2p "github.com/libp2p/go-libp2p"
	lp2pcore "github.com/libp2p/go-libp2p/core"
	lp2pcrypto "github.com/libp2p/go-libp2p/core/crypto"
	lp2phost "github.com/libp2p/go-libp2p/core/host"
	lp2pmetrics "github.com/libp2p/go-libp2p/core/metrics"
	lp2ppeer "github.com/libp2p/go-libp2p/core/peer"
	lp2pautorelay "github.com/libp2p/go-libp2p/p2p/host/autorelay"
	lp2prcmgr "github.com/libp2p/go-libp2p/p2p/host/resource-manager"
	lp2pconnmgr "github.com/libp2p/go-libp2p/p2p/net/connmgr"
	lp2pproto "github.com/libp2p/go-libp2p/p2p/protocol/circuitv2/proto"
	lp2quic "github.com/libp2p/go-libp2p/p2p/transport/quic"
	lp2ptcp "github.com/libp2p/go-libp2p/p2p/transport/tcp"
	"github.com/multiformats/go-multiaddr"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/logger"
	"github.com/pactus-project/pactus/version"
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/exp/slices"
)

var _ Network = &network{}

type network struct {
	ctx          context.Context
	cancel       context.CancelFunc
	config       *Config
	host         lp2phost.Host
	mdns         *mdnsService
	dht          *dhtService
	peerMgr      *peerMgr
	connGater    *ConnectionGater
	stream       *streamService
	gossip       *gossipService
	notifee      *NotifeeService
	eventChannel chan Event
	logger       *logger.SubLogger
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

	return makeNetwork(conf, log, []lp2p.Option{})
}

func makeNetwork(conf *Config, log *logger.SubLogger, opts []lp2p.Option) (*network, error) {
	self := new(network)

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

	resMgr, err := lp2prcmgr.NewResourceManager(
		lp2prcmgr.NewFixedLimiter(lp2prcmgr.InfiniteLimits),
		rcMgrOpt...,
	)
	if err != nil {
		return nil, LibP2PError{Err: err}
	}

	// https://docs.libp2p.io/concepts/security/dos-mitigation/#limit-the-number-of-connections-your-application-needs
	// The ConnManager is in charge of pruning connections to stay below the defined high watermark,
	// in contrast, the Resource Manager represents a hard limit where connections will fail to
	// be created in the first place once we’ve reached our limits.
	//
	lowWM := conf.MaxConns                        // Low  Watermark, ex: 64 (if max_conn = 64)
	highWM := conf.MaxConns + (conf.MaxConns / 4) // High Watermark, ex: 80 (if max_conn = 64)
	connMgr, err := lp2pconnmgr.NewConnManager(
		lowWM, highWM,
		lp2pconnmgr.WithGracePeriod(time.Minute),
	)
	if err != nil {
		return nil, LibP2PError{Err: err}
	}
	log.Info("connection manager created", "lowWM", lowWM, "highWM", highWM)

	opts = append(opts,
		lp2p.Identity(networkKey),
		lp2p.ListenAddrs(conf.ListenAddrs()...),
		lp2p.UserAgent(version.NodeAgent.String()),
		lp2p.ResourceManager(resMgr),
		lp2p.ConnectionManager(connMgr),
		lp2p.Ping(true),
		lp2p.Transport(lp2ptcp.NewTCPTransport),
	)

	if conf.EnableUDP {
		log.Info("UDP is enabled")
		opts = append(opts,
			lp2p.Transport(lp2quic.NewTransport))
	}

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
	// networkReady is a channel used to wait until the network is ready.
	// This is primarily to avoid reading while writing.
	networkReady := make(chan struct{})
	defer close(networkReady)

	networkGetter := func() *network {
		<-networkReady              // Closed when newNetwork is finished
		time.Sleep(1 * time.Second) // Adding a safety delay

		return self
	}

	if conf.EnableRelay {
		log.Info("relay enabled")

		autoRelayOpt := []lp2pautorelay.Option{
			lp2pautorelay.WithMinCandidates(1),
			lp2pautorelay.WithMaxCandidates(4),
			lp2pautorelay.WithMinInterval(1 * time.Minute),
		}

		opts = append(opts,
			lp2p.EnableRelay(),
			lp2p.EnableAutoRelayWithPeerSource(findRelayPeers(networkGetter), autoRelayOpt...),
			lp2p.EnableHolePunching(),
		)
	} else {
		log.Info("relay disabled")
		opts = append(opts,
			lp2p.DisableRelay(),
		)
	}

	if conf.EnableRelayService {
		log.Info("relay service enabled")
		opts = append(opts, lp2p.EnableRelayService())
	}

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

	self.ctx = ctx
	self.cancel = cancel
	self.config = conf
	self.logger = log
	self.host = host
	self.connGater = connGater
	self.eventChannel = make(chan Event, 100)

	log.SetObj(self)

	conf.CheckIsBootstrapper(host.ID())

	kadProtocolID := lp2pcore.ProtocolID(fmt.Sprintf("/%s/gossip/v1", conf.NetworkName)) // TODO: better name?
	streamProtocolID := lp2pcore.ProtocolID(fmt.Sprintf("/%s/stream/v1", conf.NetworkName))

	if conf.EnableMdns {
		self.mdns = newMdnsService(ctx, self.host, self.logger)
	}

	self.dht = newDHTService(self.ctx, self.host, kadProtocolID, conf, self.logger)
	self.peerMgr = newPeerMgr(ctx, host, conf, self.logger)
	self.stream = newStreamService(ctx, self.host, streamProtocolID, self.eventChannel, self.logger)
	self.gossip = newGossipService(ctx, self.host, self.eventChannel, conf, self.logger)
	self.notifee = newNotifeeService(ctx, self.host, self.eventChannel, self.peerMgr, streamProtocolID, self.logger)

	self.logger.Info("network setup", "id", self.host.ID(),
		"name", conf.NetworkName,
		"address", conf.ListenAddrs(),
		"bootstrapper", conf.IsBootstrapper,
		"maxConns", conf.MaxConns)

	return self, nil
}

func findRelayPeers(networkGetter func() *network) func(ctx context.Context,
	num int) <-chan lp2ppeer.AddrInfo {
	return func(ctx context.Context, num int) <-chan lp2ppeer.AddrInfo {
		// make a channel to return, and put items from numPeers on
		// that channel up to numPeers. Then close it.
		r := make(chan lp2ppeer.AddrInfo, num)
		defer close(r)

		// Because the network is initialized after relay, we need to
		// obtain them indirectly this way.
		n := networkGetter()
		if n == nil { // context canceled etc.
			return r
		}

		n.logger.Debug("try to find relay peers", "num", num)

		peerStore := n.host.Peerstore()
		for _, id := range peerStore.Peers() {
			protos, err := peerStore.GetProtocols(id)
			if err != nil {
				continue
			}

			if !slices.Contains(protos, lp2pproto.ProtoIDv2Hop) {
				continue
			}

			addr := peerStore.Addrs(id)
			n.logger.Debug("found relay peer", "addr", addr)
			dhtPeer := lp2ppeer.AddrInfo{ID: id, Addrs: addr}
			// Attempt to put peers on r if we have space,
			// otherwise return (we reached numPeers)
			select {
			case r <- dhtPeer:
			case <-ctx.Done():
				return r
			default:
				return r
			}
		}

		return r
	}
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

	n.host.Network().Notify(n.notifee)
	n.connGater.SetPeerManager(n.peerMgr)

	n.logger.Info("network started", "addr", n.host.Addrs(), "id", n.host.ID())

	return nil
}

func (n *network) Stop() {
	n.cancel()
	n.logger.Debug("context closed", "reason", n.ctx.Err())

	if n.mdns != nil {
		n.mdns.Stop()
	}

	n.gossip.Stop()
	n.stream.Stop()
	n.peerMgr.Stop()
	n.notifee.Stop()
	n.dht.Stop()

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

// SendTo sends a message to a specific peer identified by pid asynchronously.
// It uses a goroutine to ensure that if sending is blocked, receiving messages won't be blocked.
func (n *network) SendTo(msg []byte, pid lp2pcore.PeerID) {
	go func() {
		err := n.stream.SendRequest(msg, pid)
		if err != nil {
			n.logger.Warn("error on sending msg", "pid", pid, "error", err)
		}
	}()
}

// Broadcast sends a message to all peers subscribed to a specific topic asynchronously.
// It uses a goroutine to ensure that if broadcasting is blocked, receiving messages won't be blocked.
func (n *network) Broadcast(msg []byte, topicID TopicID) {
	go func() {
		err := n.gossip.Broadcast(msg, topicID)
		if err != nil {
			n.logger.Warn("error on broadcasting msg", "error", err)
		}
	}()
}

func (n *network) JoinTopic(topicID TopicID, sp ShouldPropagate) error {
	return n.gossip.JoinTopic(topicID, sp)
}

func (n *network) CloseConnection(pid lp2ppeer.ID) {
	n.logger.Debug("closing connection", "pid", pid)

	if err := n.host.Network().ClosePeer(pid); err != nil {
		n.logger.Warn("unable to close connection", "peer", pid)
	}

	n.logger.Debug("connection closed", "pid", pid)
}

func (n *network) String() string {
	return fmt.Sprintf("{%d}", n.NumConnectedPeers())
}

func (n *network) NumConnectedPeers() int {
	return len(n.host.Network().Peers())
}

func (n *network) ReachabilityStatus() string {
	return n.notifee.Reachability().String()
}

func (n *network) HostAddrs() []string {
	addrs := make([]string, 0, len(n.host.Addrs()))
	for _, addr := range n.host.Addrs() {
		addrs = append(addrs, addr.String())
	}

	return addrs
}

func (n *network) Name() string {
	return n.config.NetworkName
}

func (n *network) Protocols() []string {
	protocols := []string{}
	for _, p := range n.host.Mux().Protocols() {
		protocols = append(protocols, string(p))
	}

	return protocols
}

func (n *network) NumInbound() int {
	return n.peerMgr.NumInbound()
}

func (n *network) NumOutbound() int {
	return n.peerMgr.NumOutbound()
}
