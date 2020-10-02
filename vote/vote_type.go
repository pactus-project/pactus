package vote

type VoteType int

const (
	VoteTypePrevote   = VoteType(1)
	VoteTypePrecommit = VoteType(2)
)

func (t VoteType) IsValid() bool {
	switch t {
	case VoteTypePrevote, VoteTypePrecommit:
		return true
	}

	return false
}

func (t VoteType) String() string {
	switch t {
	case VoteTypePrevote:
		return "Prevote"
	case VoteTypePrecommit:
		return "Precommit"
	default:
		return ("Invalid vote type")
	}
}
