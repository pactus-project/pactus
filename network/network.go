package network

import (
	"context"
	"encoding/hex"
	"fmt"
	"time"

	exchange "github.com/ipfs/go-ipfs-exchange-interface"
	"github.com/libp2p/go-libp2p"
	acrypto "github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/routing"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	libp2pps "github.com/libp2p/go-libp2p-pubsub"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/logger"
	"github.com/zarbchain/zarb-go/utils"
)

// DiscoveryInterval is how often we re-publish our mDNS records.
const DiscoveryInterval = time.Hour

// DiscoveryServiceTag is used in our mDNS advertisements to discover other chat peers.
const DiscoveryServiceTag = "pubsub-chat-example"

type Network struct {
	ctx         context.Context
	config      *Config
	networkName string
	host        host.Host
	router      routing.Routing // Router is a router from IPFS
	pubsub      *libp2pps.PubSub
	bitswap     exchange.Interface
	logger      *logger.Logger

	// Network *net.Network
}

func loadOrCreateKey(path string) (acrypto.PrivKey, error) {
	if utils.PathExists(path) {
		h, err := utils.ReadFile(path)
		if err != nil {
			return nil, errors.Errorf(errors.ErrNetwork, err.Error())
		}
		bs, err := hex.DecodeString(string(h))
		if err != nil {
			return nil, errors.Errorf(errors.ErrNetwork, err.Error())
		}
		key, err := acrypto.UnmarshalPrivateKey(bs)
		if err != nil {
			return nil, errors.Errorf(errors.ErrNetwork, err.Error())
		}
		return key, nil
	}
	key, _, err := acrypto.GenerateEd25519Key(nil)
	if err != nil {
		return nil, errors.Errorf(errors.ErrNetwork, "failed to create peer key")
	}
	bs, err := acrypto.MarshalPrivateKey(key)
	if err != nil {
		return nil, errors.Errorf(errors.ErrNetwork, err.Error())
	}
	h := hex.EncodeToString(bs)
	err = utils.WriteFile(path, []byte(h))
	if err != nil {
		return nil, errors.Errorf(errors.ErrNetwork, err.Error())
	}
	return key, nil
}

func NewNetwork(conf *Config) (*Network, error) {
	ctx := context.Background()

	var router routing.Routing
	makeDHT := func(h host.Host) (routing.Routing, error) {
		r, err := dht.New(
			ctx,
			h,
			// dhtopts.Protocols(protocol.ID(fmt.Sprintf("/zarb/kad/%s", network))),
		)
		if err != nil {
			return nil, errors.Errorf(errors.ErrNetwork, err.Error())
		}
		router = r
		return r, err
	}
	host, err := buildHost(ctx, conf, makeDHT)
	if err != nil {
		return nil, err
	}

	libp2pps.GossipSubHeartbeatInterval = 100 * time.Millisecond
	pubsub, err := libp2pps.NewGossipSub(ctx, host)
	if err != nil {
		return nil, errors.Errorf(errors.ErrNetwork, err.Error())
	}

	n := &Network{
		ctx:         ctx,
		config:      conf,
		networkName: conf.Name,
		host:        host,
		router:      router,
		pubsub:      pubsub,
	}

	n.logger = logger.NewLogger("_Network", n)
	n.logger.Info("Network started", "id", n.host.ID(), "address", conf.Address)

	return n, nil
}

func buildHost(ctx context.Context, conf *Config, makeDHT func(host host.Host) (routing.Routing, error)) (host.Host, error) {
	// Node must build a host acting as a libp2p relay.  Additionally it
	// runs the autoNAT service which allows other nodes to check for their
	// own dialability by having this node attempt to dial them.
	makeDHTRightType := func(h host.Host) (routing.PeerRouting, error) {
		return makeDHT(h)
	}

	nodeKey, err := loadOrCreateKey(conf.NodeKey)
	if err != nil {
		return nil, err
	}

	libP2pOpts := []libp2p.Option{
		libp2p.ListenAddrStrings(conf.Address),
		libp2p.Identity(nodeKey),
	}

	return libp2p.New(
		ctx,
		libp2p.EnableAutoRelay(),
		libp2p.Routing(makeDHTRightType),
		libp2p.ChainOptions(libP2pOpts...),
	)
}

func (n *Network) Start() {
	n.setupMNSDiscovery(n.ctx, n.host)
	//NewBootstrapper()

}

func (n *Network) Stop() {
	if err := n.host.Close(); err != nil {
		n.logger.Panic("Unable to close the network", "err", err)
	}
}

func (n *Network) ID() peer.ID {
	return n.host.ID()
}

func (n *Network) Fingerprint() string {
	return fmt.Sprintf("{%d}", len(n.host.Network().Peers()))
}

func (n *Network) JoinTopic(name string) (*pubsub.Topic, error) {
	// TODO : add topic validator
	topic := fmt.Sprintf("/zarb/%s/%s", n.networkName, name)
	return n.pubsub.Join(topic)
}
