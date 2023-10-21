package network

import (
	"context"
	"time"

	lp2phost "github.com/libp2p/go-libp2p/core/host"
	lp2pnetwork "github.com/libp2p/go-libp2p/core/network"
	lp2ppeer "github.com/libp2p/go-libp2p/core/peer"
	ma "github.com/multiformats/go-multiaddr"
	"github.com/pactus-project/pactus/util/logger"
)

// MakeAddressInfo from Multi-address string.
func MakeAddressInfo(addr string) (*lp2ppeer.AddrInfo, error) {
	maddr, err := ma.NewMultiaddr(addr)
	if err != nil {
		return nil, err
	}
	return lp2ppeer.AddrInfoFromP2pAddr(maddr)
}

func ConnectAsync(ctx context.Context, h lp2phost.Host, addrInfo lp2ppeer.AddrInfo, logger *logger.SubLogger) {
	go func() {
		err := h.Connect(lp2pnetwork.WithDialPeerTimeout(ctx, 10*time.Second), addrInfo)
		if err != nil {
			if logger != nil {
				logger.Warn("connection failed", "addr", addrInfo.Addrs, "err", err)
			}
		} else {
			if logger != nil {
				logger.Debug("connected", "addr", addrInfo.Addrs)
			}
		}
	}()
}
