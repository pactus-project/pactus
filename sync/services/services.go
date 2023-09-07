package services

import (
	"fmt"

	"github.com/pactus-project/pactus/util"
)

type Services int

const (
	None    = 0x00
	Network = 0x01
)

func New(flags ...int) Services {
	s := 0
	for _, f := range flags {
		s = util.SetFlag(s, f)
	}
	return Services(s)
}

func (s Services) String() string {
	services := ""
	if util.IsFlagSet(s, Network) {
		services += "NETWORK | "
		s = util.UnsetFlag(s, Network)
	}

	if s != 0 {
		services += fmt.Sprintf("%d", s)
	} else if len(services) > 0 {
		services = services[:len(services)-3]
	}

	return services
}

func (s Services) IsNetwork() bool {
	return util.IsFlagSet(s, Network)
}
