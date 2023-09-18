package vote

type Type int

const (
	VoteTypePrepare    = Type(1) // prepare vote
	VoteTypePrecommit  = Type(2) // precommit vote
	VoteTypeCPPreVote  = Type(3) // change-proposer:pre-vote
	VoteTypeCPMainVote = Type(4) // change-proposer:main-vote
)

func (t Type) IsValid() bool {
	switch t {
	case VoteTypePrepare, VoteTypePrecommit,
		VoteTypeCPPreVote, VoteTypeCPMainVote:
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
	default:
		return ("invalid vote type")
	}
}
