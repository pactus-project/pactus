package network

import (
	"context"
	"fmt"

	lp2pcore "github.com/libp2p/go-libp2p-core"
	lp2phost "github.com/libp2p/go-libp2p-core/host"
	lp2pdht "github.com/libp2p/go-libp2p-kad-dht"
)

func (n *network) setupKademlia(ctx context.Context, host lp2phost.Host) (*lp2pdht.IpfsDHT, error) {

	opts := []lp2pdht.Option{
		lp2pdht.Mode(lp2pdht.ModeAuto),
		lp2pdht.ProtocolPrefix(lp2pcore.ProtocolID(fmt.Sprintf("/%s/kad/v1", n.config.Name))),
	}

	dht, err := lp2pdht.New(ctx, host, opts...)
	if err != nil {
		return nil, err
	}

	return dht, nil
}
