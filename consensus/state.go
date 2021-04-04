package consensus

import "github.com/zarbchain/zarb-go/vote"

type consState interface {
	enter()
	execute()
	voteAdded(v *vote.Vote)
	timedout(t *ticker)
	name() string
}

type initState struct {
}

func (s *initState) enter()                 {}
func (s *initState) execute()               {}
func (s *initState) timedout(t *ticker)     {}
func (s *initState) voteAdded(v *vote.Vote) {}
func (s *initState) name() string {
	return "initializing"
}
