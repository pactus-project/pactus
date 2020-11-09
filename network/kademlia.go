package network

import (
	"context"

	host "github.com/libp2p/go-libp2p-core/host"
	libp2pdht "github.com/libp2p/go-libp2p-kad-dht"
)

func (n *Network) setupKademlia(ctx context.Context, h host.Host) (*libp2pdht.IpfsDHT, error) {
	kademliaDHT, err := libp2pdht.New(ctx, h)
	if err != nil {
		return nil, err
	}

	return kademliaDHT, nil
}
