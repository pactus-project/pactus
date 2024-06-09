package ipblocker

import (
	"net"
)

type IPBlocker struct {
	cidrs []*net.IPNet
}

func New(bannedNets []string) (*IPBlocker, error) {
	ipBlocker := &IPBlocker{
		cidrs: make([]*net.IPNet, 0),
	}

	for _, cidr := range bannedNets {
		_, ipNet, err := net.ParseCIDR(cidr)
		if err != nil {
			return nil, err
		}
		ipBlocker.cidrs = append(ipBlocker.cidrs, ipNet)
	}

	return ipBlocker, nil
}

func (i *IPBlocker) IsBanned(ip string) bool {
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return false
	}

	// TODO: if scaled cidrs and ips items we can improve using trie or radix tree
	for _, cidr := range i.cidrs {
		if cidr.Contains(parsedIP) {
			return true
		}
	}

	return false
}
