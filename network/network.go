package network

import (
	"context"
	"encoding/hex"
	"fmt"
	"sync"

	lp2p "github.com/libp2p/go-libp2p"
	lp2pcircuit "github.com/libp2p/go-libp2p-circuit"
	lp2pcrypto "github.com/libp2p/go-libp2p-core/crypto"
	lp2phost "github.com/libp2p/go-libp2p-core/host"
	lp2peer "github.com/libp2p/go-libp2p-core/peer"
	lp2pdht "github.com/libp2p/go-libp2p-kad-dht"
	lp2pps "github.com/libp2p/go-libp2p-pubsub"
	lp2pdiscovery "github.com/libp2p/go-libp2p/p2p/discovery"
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/logger"
	"github.com/zarbchain/zarb-go/util"
	"github.com/zarbchain/zarb-go/version"
)

type network struct {
	lk sync.RWMutex

	ctx            context.Context
	config         *Config
	host           lp2phost.Host
	wg             sync.WaitGroup
	mdns           lp2pdiscovery.Service
	kademlia       *lp2pdht.IpfsDHT
	pubsub         *lp2pps.PubSub
	generalTopic   *lp2pps.Topic
	downloadTopic  *lp2pps.Topic
	dataTopic      *lp2pps.Topic
	consensusTopic *lp2pps.Topic
	generalSub     *lp2pps.Subscription
	downloadSub    *lp2pps.Subscription
	dataSub        *lp2pps.Subscription
	consensusSub   *lp2pps.Subscription
	callback       CallbackFn
	bootstrapper   *Bootstrapper
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
	ctx := context.Background()

	nodeKey, err := loadOrCreateKey(conf.NodeKeyFile)
	if err != nil {
		return nil, errors.Errorf(errors.ErrNetwork, err.Error())
	}

	opts := []lp2p.Option{
		lp2p.Identity(nodeKey),
		lp2p.ListenAddrStrings(conf.ListenAddress...),
		lp2p.Ping(true),
		lp2p.UserAgent("zarb-" + version.Version()),
	}
	if conf.EnableNATService {
		opts = append(opts,
			lp2p.EnableNATService(),
			lp2p.NATPortMap())
	}
	if conf.EnableRelay {
		opts = append(opts,
			lp2p.EnableRelay(lp2pcircuit.OptHop))
	}
	host, err := lp2p.New(ctx, opts...)
	if err != nil {
		return nil, errors.Errorf(errors.ErrNetwork, err.Error())
	}

	pubsub, err := lp2pps.NewGossipSub(ctx, host)
	if err != nil {
		return nil, errors.Errorf(errors.ErrNetwork, err.Error())
	}

	n := &network{
		ctx:    ctx,
		config: conf,
		host:   host,
		pubsub: pubsub,
	}
	n.logger = logger.NewLogger("_network", n)
	n.logger.Info("network started", "id", n.host.ID(), "address", conf.ListenAddress)

	if conf.EnableMDNS {
		mdns, err := n.setupMNSDiscovery(n.ctx, n.host)
		if err != nil {
			n.logger.Error("Unable to setup mDNS discovery", "err", err)
		}
		n.mdns = mdns
	}

	if conf.EnableKademlia {
		kademlia, err := n.setupKademlia(n.ctx, n.host)
		if err != nil {
			n.logger.Error("Unable to setup Kademlia DHT", "err", err)
		}
		n.kademlia = kademlia
		n.bootstrapper = NewBootstrapper(ctx,
			host, host.Network(), kademlia,
			conf.Bootstrap, n.logger)
	}

	return n, nil
}

func (n *network) Start() error {
	if n.bootstrapper != nil {
		n.bootstrapper.Start()
	}
	return nil
}

func (n *network) Stop() {
	n.ctx.Done()
	n.closeTopics()

	if n.mdns != nil {
		if err := n.mdns.Close(); err != nil {
			n.logger.Error("Unable to close mDNS", "err", err)
		}
	}
	if n.kademlia != nil {
		if err := n.kademlia.Close(); err != nil {
			n.logger.Error("Unable to close Kademlia", "err", err)
		}
	}
	if n.bootstrapper != nil {
		n.bootstrapper.Stop()
	}
	if err := n.host.Close(); err != nil {
		n.logger.Error("Unable to close the network", "err", err)
	}
}

func (n *network) SelfID() lp2peer.ID {
	return n.host.ID()
}

func (n *network) CloseConnection(pid lp2peer.ID) {
	if err := n.host.Network().ClosePeer(pid); err != nil {
		n.logger.Warn("Unable to close connection", "peer", pid)
	}
}

func (n *network) Fingerprint() string {
	return fmt.Sprintf("{%d}", len(n.host.Network().Peers()))
}

func (n *network) joinTopic(name string) (*lp2pps.Topic, error) {
	topic := fmt.Sprintf("/zarb/pubsub/%s/v1/%s", n.config.Name, name)
	return n.pubsub.Join(topic)
}
