package service

import (
	"fmt"

	"github.com/pactus-project/pactus/util"
)

type (
	Services int
	Service  int
)

const (
	None    Service = 0x00
	Network Service = 0x01
	Gossip  Service = 0x02
)

func New(flags ...Service) Services {
	s := None
	for _, f := range flags {
		s = util.SetFlag(s, f)
	}

	return Services(s)
}

func (s *Services) Append(flag Service) {
	*s = util.SetFlag(*s, Services(flag))
}

func (s Services) String() string {
	services := ""
	if util.IsFlagSet(s, Services(Network)) {
		services += "NETWORK | "
		s = util.UnsetFlag(s, Services(Network))
	}

	if util.IsFlagSet(s, Services(Gossip)) {
		services += "GOSSIP | "
		s = util.UnsetFlag(s, Services(Gossip))
	}

	if s != 0 {
		services += fmt.Sprintf("%d", s)
	} else if len(services) > 0 {
		services = services[:len(services)-3]
	}

	return services
}

func (s Services) IsNetwork() bool {
	return util.IsFlagSet(s, Services(Network))
}

func (s Services) IsGossip() bool {
	return util.IsFlagSet(s, Services(Gossip))
}
