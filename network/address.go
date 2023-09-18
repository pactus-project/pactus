package network

import (
	lp2ppeer "github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
)

// PeerAddrsToAddrInfo converts a slice of string peer addresses
// to AddrInfo.
func PeerAddrsToAddrInfo(addrs []string) ([]lp2ppeer.AddrInfo, error) {
	pis := make([]lp2ppeer.AddrInfo, 0, len(addrs))
	for _, addr := range addrs {
		a, err := multiaddr.NewMultiaddr(addr)
		if err != nil {
			return nil, err
		}

		pinfo, err := lp2ppeer.AddrInfoFromP2pAddr(a)
		if err != nil {
			return nil, err
		}
		pis = append(pis, *pinfo)
	}
	return pis, nil
}
