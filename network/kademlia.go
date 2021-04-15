package network

import (
	"context"
	"fmt"

	host "github.com/libp2p/go-libp2p-core/host"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	protocol "github.com/libp2p/go-libp2p-protocol"
)

func (n *network) setupKademlia(ctx context.Context, host host.Host) (*dht.IpfsDHT, error) {

	opts := []dht.Option{
		dht.Mode(dht.ModeAuto),
		dht.ProtocolPrefix(protocol.ID(fmt.Sprintf("/zarb/kad/%s", n.config.Name))),
	}

	dht, err := dht.New(ctx, host, opts...)
	if err != nil {
		return nil, err
	}

	return dht, nil
}
