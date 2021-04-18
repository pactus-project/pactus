package network

import (
	"context"
	"encoding/hex"
	"fmt"
	"sync"

	"github.com/libp2p/go-libp2p"
	circuit "github.com/libp2p/go-libp2p-circuit"
	acrypto "github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	libp2pdht "github.com/libp2p/go-libp2p-kad-dht"
	libp2pps "github.com/libp2p/go-libp2p-pubsub"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/p2p/discovery"
	"github.com/sasha-s/go-deadlock"
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/logger"
	"github.com/zarbchain/zarb-go/util"
	"github.com/zarbchain/zarb-go/version"
)

type network struct {
	lk deadlock.RWMutex

	ctx            context.Context
	config         *Config
	host           host.Host
	wg             sync.WaitGroup
	mdns           discovery.Service
	kademlia       *libp2pdht.IpfsDHT
	pubsub         *libp2pps.PubSub
	generalTopic   *pubsub.Topic
	downloadTopic  *pubsub.Topic
	dataTopic      *pubsub.Topic
	consensusTopic *pubsub.Topic
	generalSub     *pubsub.Subscription
	downloadSub    *pubsub.Subscription
	dataSub        *pubsub.Subscription
	consensusSub   *pubsub.Subscription
	callback       CallbackFn
	bootstrapper   *Bootstrapper
	logger         *logger.Logger
}

func loadOrCreateKey(path string) (acrypto.PrivKey, error) {
	if util.PathExists(path) {
		h, err := util.ReadFile(path)
		if err != nil {
			return nil, err
		}
		bs, err := hex.DecodeString(string(h))
		if err != nil {
			return nil, err
		}
		key, err := acrypto.UnmarshalPrivateKey(bs)
		if err != nil {
			return nil, err
		}
		return key, nil
	}
	key, _, err := acrypto.GenerateEd25519Key(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to generate private key")
	}
	bs, err := acrypto.MarshalPrivateKey(key)
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

	opts := []libp2p.Option{
		libp2p.Identity(nodeKey),
		libp2p.ListenAddrStrings(conf.ListenAddress...),
		libp2p.Ping(true),
		libp2p.UserAgent("zarb-" + version.NodeVersion.String()),
	}
	if conf.EnableNATService {
		opts = append(opts,
			libp2p.EnableNATService(),
			libp2p.NATPortMap())
	}
	if conf.EnableRelay {
		opts = append(opts,
			libp2p.EnableRelay(circuit.OptHop))
	}
	host, err := libp2p.New(ctx, opts...)
	if err != nil {
		return nil, errors.Errorf(errors.ErrNetwork, err.Error())
	}

	pubsub, err := libp2pps.NewGossipSub(ctx, host)
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

func (n *network) SelfID() peer.ID {
	return n.host.ID()
}

func (n *network) Fingerprint() string {
	return fmt.Sprintf("{%d}", len(n.host.Network().Peers()))
}

func (n *network) joinTopic(name string) (*pubsub.Topic, error) {
	topic := fmt.Sprintf("/zarb/pubsub/%s/v1/%s", n.config.Name, name)
	return n.pubsub.Join(topic)
}
