package addr

import (
	"fmt"
	"regexp"
)

type P2PAddr struct {
	transport string
	protocol  string
	addr      string
	port      string
	peerID    string
}

const reg = `/(?P<transport>ip4|ip6|dns)/(?P<addr>[^/]+)/(?P<protocol>tcp|udp)/(?P<port>\d+)(/p2p/(?P<peerid>[^/]+))?`

func Parse(addr string) (P2PAddr, error) {
	re := regexp.MustCompile(reg)

	match := re.FindStringSubmatch(addr)
	if match == nil {
		return P2PAddr{}, fmt.Errorf("invalid address format: %s", addr)
	}

	address := P2PAddr{}
	for i, name := range re.SubexpNames() {
		if i > 0 && name != "" {
			switch name {
			case "transport":
				address.transport = match[i]
			case "protocol":
				address.protocol = match[i]
			case "addr":
				address.addr = match[i]
			case "port":
				address.port = match[i]
			case "peerid":
				address.peerID = match[i]
			}
		}
	}

	return address, nil
}

func (p *P2PAddr) Transport() string {
	return p.transport
}

func (p *P2PAddr) Protocol() string {
	return p.protocol
}

func (p *P2PAddr) Address() string {
	return p.addr
}

func (p *P2PAddr) Port() string {
	return p.port
}

func (p *P2PAddr) PeerID() string {
	return p.peerID
}
