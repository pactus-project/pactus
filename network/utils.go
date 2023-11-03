package network

import (
	"context"
	"math/bits"
	"net"
	"time"

	lp2phost "github.com/libp2p/go-libp2p/core/host"
	lp2pnetwork "github.com/libp2p/go-libp2p/core/network"
	lp2ppeer "github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
	"github.com/pactus-project/pactus/util/logger"
)

// PeerAddrsToAddrInfo converts a slice of string peer addresses
// to AddrInfo.
func PeerAddrsToAddrInfo(addrs []string) []lp2ppeer.AddrInfo {
	pis := make([]lp2ppeer.AddrInfo, 0, len(addrs))
	for _, addr := range addrs {
		pinfo, _ := MakeAddressInfo(addr)
		if pinfo != nil {
			pis = append(pis, *pinfo)
		}
	}
	return pis
}

// MakeAddressInfo from Multi-address string.
func MakeAddressInfo(addr string) (*lp2ppeer.AddrInfo, error) {
	maddr, err := multiaddr.NewMultiaddr(addr)
	if err != nil {
		return nil, err
	}
	return lp2ppeer.AddrInfoFromP2pAddr(maddr)
}

// HasPID checks if a peer ID exists in a list of peer IDs.
func HasPID(pids []lp2ppeer.ID, pid lp2ppeer.ID) bool {
	for _, p := range pids {
		if p == pid {
			return true
		}
	}
	return false
}

func ConnectAsync(ctx context.Context, h lp2phost.Host, addrInfo lp2ppeer.AddrInfo, log *logger.SubLogger) {
	go func() {
		err := h.Connect(lp2pnetwork.WithDialPeerTimeout(ctx, 30*time.Second), addrInfo)
		if err != nil {
			if log != nil {
				log.Warn("connection failed", "addr", addrInfo.Addrs, "err", err)
			}
		} else {
			if log != nil {
				log.Debug("connected", "addr", addrInfo.Addrs)
			}
		}
	}()
}

func LogScale(val int) int {
	bitlen := bits.Len(uint(val))
	return 1 << bitlen
}

func PrivateSubnets() []*net.IPNet {
	privateCIDRs := []string{
		// -- Ipv4 --
		// localhost
		"127.0.0.0/8",
		// private networks
		"10.0.0.0/8",
		"100.64.0.0/10",
		"172.16.0.0/12",
		"192.168.0.0/16",
		// link local
		"169.254.0.0/16",

		// -- Ipv6 --
		// localhost
		"::1/128",
		// ULA reserved
		"fc00::/7",
		// link local
		"fe80::/10",
	}

	subnets := []*net.IPNet{}
	for _, cidr := range privateCIDRs {
		_, sn, err := net.ParseCIDR(cidr)
		if err != nil {
			panic(err)
		}
		subnets = append(subnets, sn)
	}

	return subnets
}

func SubnetsToFilters(subnets []*net.IPNet, action multiaddr.Action) *multiaddr.Filters {
	filters := multiaddr.NewFilters()
	for _, sn := range subnets {
		filters.AddFilter(*sn, action)
	}

	return filters
}
