package ipblocker

import (
	"net"
)

type IPBlocker struct {
	cidrs []*net.IPNet
}

func New(blackListCidr []string) (*IPBlocker, error) {
	ipblocker := &IPBlocker{
		cidrs: make([]*net.IPNet, 0),
	}

	for _, cidr := range blackListCidr {
		_, ipnet, err := net.ParseCIDR(cidr)
		if err != nil {
			return nil, err
		}
		ipblocker.cidrs = append(ipblocker.cidrs, ipnet)
	}

	return ipblocker, nil
}

func (i *IPBlocker) IsBlocked(ip string) bool {
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return false
	}

	// TODO: if scaled cidrs and ips items we can improve using trys or radix tree
	for _, cidr := range i.cidrs {
		if cidr.Contains(parsedIP) {
			return true
		}
	}

	return false
}
