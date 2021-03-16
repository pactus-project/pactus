package hrs

import (
	"fmt"

	"github.com/fxamacker/cbor/v2"
	"github.com/sasha-s/go-deadlock"
)

type HRS struct {
	lk deadlock.RWMutex

	data hrsData
}

type hrsData struct {
	Height int
	Round  int
	Step   StepType
}

func NewHRS(height int, round int, step StepType) *HRS {
	return &HRS{
		data: hrsData{
			Height: height,
			Round:  round,
			Step:   step,
		},
	}
}

func (hrs *HRS) IsValid() bool {
	hrs.lk.Lock()
	defer hrs.lk.Unlock()

	if hrs.data.Height <= 0 || hrs.data.Round < 0 {
		return false
	}
	return hrs.data.Step.IsValid()
}

func (hrs *HRS) Height() int {
	hrs.lk.Lock()
	defer hrs.lk.Unlock()

	return hrs.data.Height
}

func (hrs *HRS) Round() int {
	hrs.lk.Lock()
	defer hrs.lk.Unlock()

	return hrs.data.Round
}

func (hrs *HRS) Step() StepType {
	hrs.lk.Lock()
	defer hrs.lk.Unlock()

	return hrs.data.Step
}

func (hrs *HRS) UpdateHeight(height int) {
	hrs.lk.Lock()
	defer hrs.lk.Unlock()

	hrs.data.Height = height
}

func (hrs *HRS) UpdateRound(round int) {
	hrs.lk.Lock()
	defer hrs.lk.Unlock()

	hrs.data.Round = round
}

func (hrs *HRS) UpdateStep(step StepType) {
	hrs.lk.Lock()
	defer hrs.lk.Unlock()

	hrs.data.Step = step
}

func (hrs *HRS) LessThan(r *HRS) bool {
	hrs.lk.Lock()
	defer hrs.lk.Unlock()

	if hrs.data.Height < r.data.Height ||
		(hrs.data.Height == r.data.Height && hrs.data.Round < r.data.Round) ||
		(hrs.data.Height == r.data.Height && hrs.data.Round == r.data.Round && hrs.data.Step < r.data.Step) {
		return true
	}
	return false
}

func (hrs *HRS) EqualsTo(r *HRS) bool {
	hrs.lk.Lock()
	defer hrs.lk.Unlock()

	if hrs.data.Height == r.data.Height && hrs.data.Round == r.data.Round && hrs.data.Step == r.data.Step {
		return true
	}
	return false
}

func (hrs *HRS) GreaterThan(r *HRS) bool {
	if hrs.LessThan(r) {
		return false
	}
	if hrs.EqualsTo(r) {
		return false
	}
	return true
}

func (hrs *HRS) String() string {
	hrs.lk.Lock()
	defer hrs.lk.Unlock()

	return fmt.Sprintf("%v/%v/%s",
		hrs.data.Height, hrs.data.Round, hrs.data.Step)
}

func (hrs *HRS) MarshalCBOR() ([]byte, error) {
	hrs.lk.Lock()
	defer hrs.lk.Unlock()

	return cbor.Marshal(hrs.data)
}

func (hrs *HRS) UnmarshalCBOR(bs []byte) error {
	hrs.lk.Lock()
	defer hrs.lk.Unlock()

	return cbor.Unmarshal(bs, &hrs.data)
}
