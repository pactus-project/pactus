package hrs

type StepType int

const (
	StepTypeUnknown   = StepType(0)
	StepTypeNewHeight = StepType(1)
	StepTypeNewRound  = StepType(2)
	StepTypePropose   = StepType(3)
	StepTypePrepare   = StepType(4)
	StepTypePrecommit = StepType(5)
	StepTypeCommit    = StepType(6)
)

// IsValid returns true if the step is valid, otherwise false.
func (rs StepType) IsValid() bool {
	return rs >= StepTypeNewHeight && rs <= StepTypeCommit
}

func (rs StepType) String() string {
	switch rs {
	case StepTypeUnknown:
		return "Unknown"
	case StepTypeNewHeight:
		return "NewHeight"
	case StepTypeNewRound:
		return "NewRound"
	case StepTypePropose:
		return "Propose"
	case StepTypePrepare:
		return "Prepare"
	case StepTypePrecommit:
		return "Precommit"
	case StepTypeCommit:
		return "Commit"
	default:
		return "Unknown"
	}
}
