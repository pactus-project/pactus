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

func (hrs *HRS) IsValid() bool {
	if hrs.data.Height < 0 || hrs.data.Round < 0 {
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

func (hrs *HRS) UpdateRoundStep(round int, step StepType) {
	hrs.data.Round = round
	hrs.data.Step = step
}

func (hrs *HRS) UpdateHeightRoundStep(height int, round int, step StepType) {
	hrs.data.Height = height
	hrs.data.Round = round
	hrs.data.Step = step
}

func (hrs HRS) InvalidHeight(height int) bool {
	return hrs.data.Height != height
}

func (hrs HRS) InvalidHeightRound(height int, round int) bool {
	return hrs.data.Height != height || hrs.data.Round != round
}

func (hrs HRS) InvalidHeightRoundStep(height int, round int, step StepType) bool {
	return hrs.data.Height != height || hrs.data.Round != round || hrs.data.Step > step
}

func (hrs *HRS) Fingerprint() string {
	return fmt.Sprintf("%v/%v/%s",
		hrs.data.Height, hrs.data.Round, hrs.data.Step)
}

func (hrs *HRS) MarshalCBOR() ([]byte, error) {
	return cbor.Marshal(hrs.data)
}

func (hrs *HRS) UnmarshalCBOR(bs []byte) error {
	return cbor.Unmarshal(bs, &hrs.data)
}
