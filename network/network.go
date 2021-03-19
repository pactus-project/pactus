package network

import (
	"context"
	"encoding/hex"
	"fmt"

	"github.com/libp2p/go-libp2p"
	circuit "github.com/libp2p/go-libp2p-circuit"
	acrypto "github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/peerstore"
	libp2pdht "github.com/libp2p/go-libp2p-kad-dht"
	libp2pps "github.com/libp2p/go-libp2p-pubsub"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/p2p/discovery"
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/logger"
	"github.com/zarbchain/zarb-go/util"
	"github.com/zarbchain/zarb-go/version"
)

type Network struct {
	ctx          context.Context
	config       *Config
	host         host.Host
	pubsub       *libp2pps.PubSub
	mdns         discovery.Service
	kademlia     *libp2pdht.IpfsDHT
	bootstrapper *Bootstrapper
	logger       *logger.Logger
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

func NewNetwork(conf *Config) (*Network, error) {
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

	n := &Network{
		ctx:    ctx,
		config: conf,
		host:   host,
		pubsub: pubsub,
	}
	n.logger = logger.NewLogger("_network", n)
	n.logger.Info("Network started", "id", n.host.ID(), "address", conf.ListenAddress)

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

func (n *Network) Start() {
	if n.bootstrapper != nil {
		n.bootstrapper.Start()
	}
}

func (n *Network) Stop() {
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

func (n *Network) ID() peer.ID {
	return n.host.ID()
}

func (n *Network) Peerstore() peerstore.Peerstore {
	return n.host.Peerstore()
}

func (n *Network) Fingerprint() string {
	return fmt.Sprintf("{%d}", len(n.host.Network().Peers()))
}

func (n *Network) JoinTopic(name string) (*pubsub.Topic, error) {
	topic := fmt.Sprintf("/zarb/pubsub/%s/v1/%s", n.config.Name, name)
	return n.pubsub.Join(topic)
}
