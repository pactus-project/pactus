package hrs

type StepType int

const (
	StepTypeUnknown       = StepType(0)
	StepTypeNewHeight     = StepType(1)
	StepTypeNewRound      = StepType(2)
	StepTypePropose       = StepType(3)
	StepTypePrepare       = StepType(4)
	StepTypePrepareWait   = StepType(5)
	StepTypePrecommit     = StepType(6)
	StepTypePrecommitWait = StepType(7)
	StepTypeCommit        = StepType(8)
)

// IsValid returns true if the step is valid, false if unknown/undefined.
func (rs StepType) IsValid() bool {
	return uint8(rs) >= 0x01 && uint8(rs) <= 0x08
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
	case StepTypePrepareWait:
		return "PrepareWait"
	case StepTypePrecommit:
		return "Precommit"
	case StepTypePrecommitWait:
		return "PrecommitWait"
	case StepTypeCommit:
		return "Commit"
	default:
		return "Unknown"
	}
}
