package network

import (
	"context"

	lp2pdht "github.com/libp2p/go-libp2p-kad-dht"
	lp2pcore "github.com/libp2p/go-libp2p/core"
	lp2phost "github.com/libp2p/go-libp2p/core/host"
	"github.com/pactus-project/pactus/util/logger"
)

type dhtService struct {
	ctx      context.Context
	host     lp2phost.Host
	kademlia *lp2pdht.IpfsDHT
	peerMgr  *peerMgr
	logger   *logger.SubLogger
}

func newDHTService(ctx context.Context, host lp2phost.Host, protocolID lp2pcore.ProtocolID,
	conf *Config, logger *logger.SubLogger,
) *dhtService {
	mode := lp2pdht.ModeAuto
	if conf.Bootstrapper {
		mode = lp2pdht.ModeServer
	}

	bootsrapAddrs := PeerAddrsToAddrInfo(conf.BootstrapAddrs)

	opts := []lp2pdht.Option{
		lp2pdht.Mode(mode),
		lp2pdht.ProtocolPrefix(protocolID),
		lp2pdht.BootstrapPeers(bootsrapAddrs...),
	}

	kademlia, err := lp2pdht.New(ctx, host, opts...)
	if err != nil {
		logger.Panic("unable to start DHT service", "error", err)
		return nil
	}

	err = kademlia.Bootstrap(ctx)
	if err != nil {
		panic(err.Error())
	}

	bootstrap := newPeerMgr(ctx, host, host.Network(), kademlia,
		bootsrapAddrs, conf.MinConns, conf.MaxConns, logger)

	return &dhtService{
		ctx:      ctx,
		host:     host,
		kademlia: kademlia,
		peerMgr:  bootstrap,
		logger:   logger,
	}
}

func (dht *dhtService) Start() error {
	dht.peerMgr.Start()
	return nil
}

func (dht *dhtService) Stop() {
	if err := dht.kademlia.Close(); err != nil {
		dht.logger.Error("unable to close Kademlia", "error", err)
	}

	dht.peerMgr.Stop()
}
