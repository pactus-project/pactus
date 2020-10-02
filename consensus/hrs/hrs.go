package hrs

import "fmt"

type HRS struct {
	height int
	round  int
	step   StepType
}

func NewHRS(height int, round int, step StepType) HRS {
	return HRS{height, round, step}
}

func (hrs HRS) Height() int {
	return hrs.height
}

func (hrs HRS) Round() int {
	return hrs.round
}

func (hrs HRS) Step() StepType {
	return hrs.step
}

func (hrs *HRS) UpdateHeight(height int) {
	hrs.height = height
}

func (hrs *HRS) UpdateRoundStep(round int, step StepType) {
	hrs.round = round
	hrs.step = step
}

func (hrs *HRS) UpdateHeightRoundStep(height int, round int, step StepType) {
	hrs.height = height
	hrs.round = round
	hrs.step = step
}

func (hrs HRS) InvalidHeight(height int) bool {
	return hrs.height != height
}

func (hrs HRS) InvalidHeightRound(height int, round int) bool {
	return hrs.height != height || hrs.round != round
}

func (hrs HRS) InvalidHeightRoundStep(height int, round int, step StepType) bool {
	return hrs.height != height || hrs.round != round || hrs.step > step
}

func (hrs *HRS) Fingerprint() string {
	return fmt.Sprintf("%v/%v/%s",
		hrs.height, hrs.round, hrs.step)
}
