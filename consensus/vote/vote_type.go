package vote

type Type int

const (
	VoteTypePrepare        = Type(1)
	VoteTypePrecommit      = Type(2)
	VoteTypeChangeProposer = Type(3)
)

func (t Type) IsValid() bool {
	switch t {
	case VoteTypePrepare, VoteTypePrecommit, VoteTypeChangeProposer:
		return true
	}

	return false
}

func (t Type) String() string {
	switch t {
	case VoteTypePrepare:
		return "prepare"
	case VoteTypePrecommit:
		return "precommit"
	case VoteTypeChangeProposer:
		return "changeProposer"
	default:
		return ("invalid vote type")
	}
}
