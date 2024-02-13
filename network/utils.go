package network

import (
	"context"
	"fmt"
	"net"
	"slices"
	"time"

	lp2pspb "github.com/libp2p/go-libp2p-pubsub/pb"
	lp2phost "github.com/libp2p/go-libp2p/core/host"
	lp2pnetwork "github.com/libp2p/go-libp2p/core/network"
	lp2ppeer "github.com/libp2p/go-libp2p/core/peer"
	lp2prcmgr "github.com/libp2p/go-libp2p/p2p/host/resource-manager"
	lp2pswarm "github.com/libp2p/go-libp2p/p2p/net/swarm"
	"github.com/multiformats/go-multiaddr"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/util/logger"
)

// MakeMultiAddrs converts a slice of string peer addresses to MultiAddress.
func MakeMultiAddrs(addrs []string) ([]multiaddr.Multiaddr, error) {
	mas := make([]multiaddr.Multiaddr, 0, len(addrs))
	for _, addr := range addrs {
		ma, err := multiaddr.NewMultiaddr(addr)
		if err != nil {
			return nil, err
		}
		mas = append(mas, ma)
	}

	return mas, nil
}

// MakeAddrInfos converts a slice of string peer addresses to AddrInfo.
func MakeAddrInfos(addrs []string) ([]lp2ppeer.AddrInfo, error) {
	pis := make([]lp2ppeer.AddrInfo, 0, len(addrs))
	for _, addr := range addrs {
		pinfo, err := lp2ppeer.AddrInfoFromString(addr)
		if err != nil {
			return nil, err
		}
		pis = append(pis, *pinfo)
	}

	return pis, nil
}

func IPToMultiAddr(ip string, port int) (multiaddr.Multiaddr, error) {
	ipParsed := net.ParseIP(ip)
	if ipParsed == nil {
		return nil, fmt.Errorf("invalid IP address: %s", ip)
	}

	var addr string
	if ipParsed.To4() != nil {
		addr = fmt.Sprintf("/ip4/%s/tcp/%d", ip, port)
	} else {
		addr = fmt.Sprintf("/ip6/%s/tcp/%d", ip, port)
	}
	ma, err := multiaddr.NewMultiaddr(addr)
	if err != nil {
		return nil, err
	}

	return ma, nil
}

// HasPID checks if a peer ID exists in a list of peer IDs.
func HasPID(pids []lp2ppeer.ID, pid lp2ppeer.ID) bool {
	return slices.Contains(pids, pid)
}

func ConnectAsync(ctx context.Context, h lp2phost.Host, addrInfo lp2ppeer.AddrInfo, log *logger.SubLogger) {
	go func() {
		err := ConnectSync(ctx, h, addrInfo)
		if log != nil {
			if err != nil {
				log.Warn("connection failed", "addr", addrInfo.Addrs, "err", err)
			} else {
				log.Debug("connection successful", "addr", addrInfo.Addrs)
			}
		}
	}()
}

func ConnectSync(ctx context.Context, h lp2phost.Host, addrInfo lp2ppeer.AddrInfo) error {
	if swarm, ok := h.Network().(*lp2pswarm.Swarm); ok {
		swarm.Backoff().Clear(addrInfo.ID)
	}

	return h.Connect(lp2pnetwork.WithDialPeerTimeout(ctx, 30*time.Second), addrInfo)
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

func BuildConcreteLimitConfig(maxConns int) lp2prcmgr.ConcreteLimitConfig {
	changes := lp2prcmgr.PartialLimitConfig{}

	updateResourceLimits := func(limit *lp2prcmgr.ResourceLimits, maxConns int, coefficient float32) {
		maxConnVal := lp2prcmgr.LimitVal(int(float32(maxConns) * coefficient))

		limit.ConnsInbound = maxConnVal
		limit.ConnsOutbound = maxConnVal
		limit.Conns = maxConnVal
		limit.StreamsInbound = maxConnVal * 8
		limit.StreamsOutbound = maxConnVal * 8
		limit.Streams = maxConnVal * 8
	}

	updateResourceLimits(&changes.System, maxConns, 1)
	updateResourceLimits(&changes.ServiceDefault, maxConns, 1)
	updateResourceLimits(&changes.ProtocolDefault, maxConns, 1)
	updateResourceLimits(&changes.ProtocolPeerDefault, maxConns, 1)
	updateResourceLimits(&changes.Transient, maxConns, 1)

	defaultLimitConfig := lp2prcmgr.DefaultLimits.AutoScale()
	changedLimitConfig := changes.Build(defaultLimitConfig)

	return changedLimitConfig
}

func MessageIDFunc(m *lp2pspb.Message) string {
	h := hash.CalcHash(m.Data)

	return string(h[:20])
}
