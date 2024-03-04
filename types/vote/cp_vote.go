package vote

import (
	"fmt"

	"github.com/fxamacker/cbor/v2"
	"github.com/pactus-project/pactus/util/errors"
)

type CPValue int8

const (
	CPValueZero    = CPValue(0)
	CPValueOne     = CPValue(1)
	CPValueAbstain = CPValue(2)
)

func (v CPValue) String() string {
	switch v {
	case CPValueZero:
		return "zero"
	case CPValueOne:
		return "one"
	case CPValueAbstain:
		return "abstain"
	default:
		return fmt.Sprintf("unknown: %d", v)
	}
}

type cpVote struct {
	Round int16
	Value CPValue
	Just  Just
}

func (v *cpVote) BasicCheck() error {
	if v.Round < 0 {
		return errors.Error(errors.ErrInvalidRound)
	}
	if v.Value < CPValueZero ||
		v.Value > CPValueAbstain {
		// Invalid values
		return errors.Errorf(errors.ErrInvalidVote, "cp value should be 0, 1 or abstain")
	}

	return v.Just.BasicCheck()
}

type _cpVote struct {
	Round    int16    `cbor:"1,keyasint"`
	Value    CPValue  `cbor:"2,keyasint"`
	JustType JustType `cbor:"3,keyasint"`
	JustData []byte   `cbor:"4,keyasint"`
}

type _JustMainVoteConflict struct {
	Just0Type JustType `cbor:"1,keyasint"`
	Just0Data []byte   `cbor:"2,keyasint"`
	Just1Type JustType `cbor:"3,keyasint"`
	Just1Data []byte   `cbor:"4,keyasint"`
}

// MarshalCBOR marshals the cpVote into CBOR format.
func (v *cpVote) MarshalCBOR() ([]byte, error) {
	justData := []byte{}
	if v.Just.Type() == JustTypeMainVoteConflict {
		conflictJust := v.Just.(*JustMainVoteConflict)
		data0, err := cbor.Marshal(conflictJust.Just0)
		if err != nil {
			return nil, err
		}
		data1, err := cbor.Marshal(conflictJust.Just1)
		if err != nil {
			return nil, err
		}

		_conflictingJust := _JustMainVoteConflict{
			Just0Type: conflictJust.Just0.Type(),
			Just0Data: data0,
			Just1Type: conflictJust.Just1.Type(),
			Just1Data: data1,
		}
		data, err := cbor.Marshal(_conflictingJust)
		if err != nil {
			return nil, err
		}
		justData = append(justData, data...)
	} else {
		data, err := cbor.Marshal(v.Just)
		if err != nil {
			return nil, err
		}
		justData = append(justData, data...)
	}

	msg := &_cpVote{
		Round:    v.Round,
		Value:    v.Value,
		JustType: v.Just.Type(),
		JustData: justData,
	}

	return cbor.Marshal(msg)
}

// UnmarshalCBOR unmarshals the cpVote from CBOR format.
func (v *cpVote) UnmarshalCBOR(bs []byte) error {
	var _cp _cpVote
	err := cbor.Unmarshal(bs, &_cp)
	if err != nil {
		return err
	}

	var just Just
	if _cp.JustType == JustTypeMainVoteConflict {
		_conflictingJust := &_JustMainVoteConflict{}
		err := cbor.Unmarshal(_cp.JustData, _conflictingJust)
		if err != nil {
			return err
		}

		just0 := makeJust(_conflictingJust.Just0Type)
		err = cbor.Unmarshal(_conflictingJust.Just0Data, just0)
		if err != nil {
			return err
		}

		just1 := makeJust(_conflictingJust.Just1Type)
		err = cbor.Unmarshal(_conflictingJust.Just1Data, just1)
		if err != nil {
			return err
		}

		just = &JustMainVoteConflict{
			Just0: just0,
			Just1: just1,
		}
	} else {
		just = makeJust(_cp.JustType)
		err := cbor.Unmarshal(_cp.JustData, just)
		if err != nil {
			return err
		}
	}

	v.Round = _cp.Round
	v.Value = _cp.Value
	v.Just = just

	return nil
}
