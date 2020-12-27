package hrs

import (
	"fmt"

	"github.com/fxamacker/cbor/v2"
)

type HRS struct {
	data hrsData
}

type hrsData struct {
	Height int
	Round  int
	Step   StepType
}

func NewHRS(height int, round int, step StepType) HRS {
	return HRS{
		data: hrsData{
			Height: height,
			Round:  round,
			Step:   step,
		},
	}
}

func (hrs HRS) IsValid() bool {
	if hrs.data.Height <= 0 || hrs.data.Round < 0 {
		return false
	}
	return hrs.data.Step.IsValid()
}

func (hrs HRS) Height() int {
	return hrs.data.Height
}

func (hrs HRS) Round() int {
	return hrs.data.Round
}

func (hrs HRS) Step() StepType {
	return hrs.data.Step
}

func (hrs *HRS) UpdateHeight(height int) {
	hrs.data.Height = height
}

func (hrs *HRS) UpdateRound(round int) {
	hrs.data.Round = round
}

func (hrs *HRS) UpdateStep(step StepType) {
	hrs.data.Step = step
}

func (hrs HRS) LessThan(r HRS) bool {
	if hrs.Height() < r.Height() ||
		(hrs.Height() == r.Height() && hrs.Round() < r.Round()) ||
		(hrs.Height() == r.Height() && hrs.Round() == r.Round() && hrs.Step() < r.Step()) {
		return true
	}
	return false
}

func (hrs HRS) EqualsTo(r HRS) bool {
	if hrs.Height() == r.Height() && hrs.Round() == r.Round() && hrs.Step() == r.Step() {
		return true
	}
	return false
}

func (hrs HRS) GreaterThan(r HRS) bool {
	if hrs.LessThan(r) {
		return false
	}
	if hrs.EqualsTo(r) {
		return false
	}
	return true
}

func (hrs HRS) String() string {
	return fmt.Sprintf("%v/%v/%s",
		hrs.data.Height, hrs.data.Round, hrs.data.Step)
}

func (hrs *HRS) MarshalCBOR() ([]byte, error) {
	return cbor.Marshal(hrs.data)
}

func (hrs *HRS) UnmarshalCBOR(bs []byte) error {
	return cbor.Unmarshal(bs, &hrs.data)
}
