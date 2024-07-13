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
	None Service = 0x00

	// FullNode indicates that the node has a full blockchain history.
	FullNode Service = 0x01

	// PrunedNode indicates that the node has a pruned blockchain history.
	PrunedNode Service = 0x02
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
	if util.IsFlagSet(flags, Services(FullNode)) {
		services += "FULL | "
		flags = util.UnsetFlag(flags, Services(FullNode))
	}

	if util.IsFlagSet(flags, Services(PrunedNode)) {
		services += "PRUNED | "
		flags = util.UnsetFlag(flags, Services(PrunedNode))
	}

	if flags != 0 {
		services += fmt.Sprintf("%d", flags)
	} else if services != "" {
		services = services[:len(services)-3]
	}

	return services
}

func (s Services) IsFullNode() bool {
	return util.IsFlagSet(s, Services(FullNode))
}

func (s Services) IsPrunedNode() bool {
	return util.IsFlagSet(s, Services(PrunedNode))
}
