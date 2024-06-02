package ipblocker

import (
	"net"
	"sync"
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

	var wg sync.WaitGroup
	result := make(chan bool, 1)

	for _, cidr := range i.cidrs {
		wg.Add(1)
		go func(cidr *net.IPNet) {
			defer wg.Done()
			if cidr.Contains(parsedIP) {
				select {
				case result <- true:
				default:
				}
			}
		}(cidr)
	}

	go func() {
		wg.Wait()
		close(result)
	}()

	return <-result
}
