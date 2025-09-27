package network

import (
	"context"

	lp2pdht "github.com/libp2p/go-libp2p-kad-dht"
	lp2pcore "github.com/libp2p/go-libp2p/core"
	lp2phost "github.com/libp2p/go-libp2p/core/host"
	"github.com/multiformats/go-multiaddr"
	"github.com/pactus-project/pactus/util/logger"
)

type dhtService struct {
	ctx      context.Context
	host     lp2phost.Host
	kademlia *lp2pdht.IpfsDHT
	logger   *logger.SubLogger
}

func newDHTService(ctx context.Context, host lp2phost.Host, conf *Config,
	protocolID lp2pcore.ProtocolID, log *logger.SubLogger,
) *dhtService {
	// A dirty code in LibP2P!!!
	// prevent apply default bootstrap node of libp2p
	lp2pdht.DefaultBootstrapPeers = []multiaddr.Multiaddr{}

	mode := lp2pdht.ModeAuto
	if conf.IsBootstrapper {
		mode = lp2pdht.ModeServer
	}
	opts := []lp2pdht.Option{
		lp2pdht.Mode(mode),
		lp2pdht.ProtocolPrefix(protocolID),
		lp2pdht.BootstrapPeers(conf.BootstrapAddrInfos()...),
	}

	kademlia, err := lp2pdht.New(ctx, host, opts...)
	if err != nil {
		panic(err)
	}

	return &dhtService{
		ctx:      ctx,
		host:     host,
		kademlia: kademlia,
		logger:   log,
	}
}

func (s *dhtService) Start() error {
	return s.kademlia.Bootstrap(s.ctx)
}

func (s *dhtService) Stop() {
	if err := s.kademlia.Close(); err != nil {
		s.logger.Error("unable to close Kademlia", "error", err)
	}
}
