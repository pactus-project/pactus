package network

import (
	"context"
	"fmt"
	"math/bits"
	"net"
	"time"

	lp2phost "github.com/libp2p/go-libp2p/core/host"
	lp2pnetwork "github.com/libp2p/go-libp2p/core/network"
	lp2ppeer "github.com/libp2p/go-libp2p/core/peer"
	lp2prcmgr "github.com/libp2p/go-libp2p/p2p/host/resource-manager"
	"github.com/multiformats/go-multiaddr"
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
		pinfo, err := MakeAddressInfo(addr)
		if err != nil {
			return nil, err
		}
		pis = append(pis, *pinfo)
	}
	return pis, nil
}

// MakeAddressInfo from Multi-address string.
func MakeAddressInfo(addr string) (*lp2ppeer.AddrInfo, error) {
	ma, err := multiaddr.NewMultiaddr(addr)
	if err != nil {
		return nil, err
	}
	return lp2ppeer.AddrInfoFromP2pAddr(ma)
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

// LogScale computes 2^⌈log₂(val)⌉, where ⌈x⌉ represents the ceiling of x.
// For more information, refer to: https://en.wikipedia.org/wiki/Logarithmic_scale
func LogScale(val int) int {
	bitlen := bits.Len(uint(val - 1))
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

func MakeScalingLimitConfig(minConns, maxConns int) lp2prcmgr.ScalingLimitConfig {
	limit := lp2prcmgr.DefaultLimits

	limit.SystemBaseLimit.ConnsInbound = LogScale(maxConns)
	limit.SystemBaseLimit.Conns = LogScale(2 * maxConns)
	limit.SystemBaseLimit.StreamsInbound = LogScale(maxConns)
	limit.SystemBaseLimit.Streams = LogScale(2 * maxConns)

	limit.ServiceLimitIncrease.ConnsInbound = LogScale(minConns)
	limit.ServiceLimitIncrease.Conns = LogScale(2 * minConns)
	limit.ServiceLimitIncrease.StreamsInbound = LogScale(minConns)
	limit.ServiceLimitIncrease.Streams = LogScale(2 * minConns)

	limit.TransientBaseLimit.ConnsInbound = LogScale(maxConns / 2)
	limit.TransientBaseLimit.Conns = LogScale(2 * maxConns / 2)
	limit.TransientBaseLimit.StreamsInbound = LogScale(maxConns / 2)
	limit.TransientBaseLimit.Streams = LogScale(2 * maxConns / 2)

	limit.TransientLimitIncrease.ConnsInbound = LogScale(minConns / 2)
	limit.TransientLimitIncrease.Conns = LogScale(2 * minConns / 2)
	limit.TransientLimitIncrease.StreamsInbound = LogScale(minConns / 2)
	limit.TransientLimitIncrease.Streams = LogScale(2 * minConns / 2)

	return limit
}

// DetectPublicIPv4 returns the public IPv4 address of the local machine by dialing a list of specific IP addresses.
func DetectPublicIPv4() (string, bool) {
	ips := []string{
		"8.8.8.8",        // Google
		"1.1.1.1",        // Cloudflare
		"9.9.9.9",        // Quad9
		"208.67.222.222", // OpenDNS
	}

	for _, ip := range ips {
		conn, err := net.Dial("udp", ip+":53") // Use a UDP connection to avoid blocking
		if err == nil {
			conn.Close()

			localAddr := conn.LocalAddr().(*net.UDPAddr)
			return localAddr.IP.String(), true
		}
	}

	// If all attempts failed, return false
	return "", false
}

// DetectPublicIPv6 returns the public IPv6 address of the local machine by dialing a list of specific IPv6 addresses.
func DetectPublicIPv6() (string, bool) {
	ips := []string{
		"2001:4860:4860::8888", // Google
		"2620:0:ccc::2",        // OpenDNS
		"2620:fe::fe",          // Quad9
		"2606:4700:4700::1111", // Cloudflare
	}

	for _, ip := range ips {
		conn, err := net.Dial("udp", "["+ip+"]:53") // Use a UDP connection to avoid blocking
		if err == nil {
			conn.Close()

			localAddr := conn.LocalAddr().(*net.UDPAddr)
			return localAddr.IP.String(), true
		}
	}

	// If all attempts failed, return false
	return "", false
}
