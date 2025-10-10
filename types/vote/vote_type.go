package vote

import "fmt"

type Type int

const (
	VoteTypePrepare    = Type(1) // Deprecated  prepare vote
	VoteTypePrecommit  = Type(2) // precommit vote
	VoteTypeCPPreVote  = Type(3) // change-proposer:pre-vote
	VoteTypeCPMainVote = Type(4) // change-proposer:main-vote
	VoteTypeCPDecided  = Type(5) // change-proposer:decided
)

func (t Type) IsValid() bool {
	switch t {
	case VoteTypePrepare, VoteTypePrecommit,
		VoteTypeCPPreVote, VoteTypeCPMainVote, VoteTypeCPDecided:
		return true
	}

	return false
}

func (t Type) String() string {
	switch t {
	case VoteTypePrepare:
		return "PREPARE"
	case VoteTypePrecommit:
		return "PRECOMMIT"
	case VoteTypeCPPreVote:
		return "PRE-VOTE"
	case VoteTypeCPMainVote:
		return "MAIN-VOTE"
	case VoteTypeCPDecided:
		return "DECIDED"
	default:
		return fmt.Sprintf("%d", t)
	}
}
