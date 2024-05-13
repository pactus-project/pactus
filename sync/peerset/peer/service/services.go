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
	Foo     Service = 0x02 // For future use
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
	flags := s
	if util.IsFlagSet(flags, Services(Network)) {
		services += "NETWORK | "
		flags = util.UnsetFlag(flags, Services(Network))
	}

	if util.IsFlagSet(flags, Services(Foo)) {
		services += "FOO | "
		flags = util.UnsetFlag(flags, Services(Foo))
	}

	if flags != 0 {
		services += fmt.Sprintf("%d", flags)
	} else if services != "" {
		services = services[:len(services)-3]
	}

	return services
}

func (s Services) IsNetwork() bool {
	return util.IsFlagSet(s, Services(Network))
}

func (s Services) IsFoo() bool {
	return util.IsFlagSet(s, Services(Foo))
}
