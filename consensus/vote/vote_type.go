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
		return "Prepare"
	case VoteTypePrecommit:
		return "Precommit"
	case VoteTypeChangeProposer:
		return "ChangeProposer"
	default:
		return ("Invalid vote type")
	}
}
