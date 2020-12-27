package vote

type VoteType int

const (
	VoteTypePrepare   = VoteType(1)
	VoteTypePrecommit = VoteType(2)
)

func (t VoteType) IsValid() bool {
	switch t {
	case VoteTypePrepare, VoteTypePrecommit:
		return true
	}

	return false
}

func (t VoteType) String() string {
	switch t {
	case VoteTypePrepare:
		return "Prepare"
	case VoteTypePrecommit:
		return "Precommit"
	default:
		return ("Invalid vote type")
	}
}
