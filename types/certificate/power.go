package certificate

// faultyPower calculates the maximum faulty power based on the total voting power.
// The formula used is: `f = (n - 1) / 3`, where `n` is the total voting power.
func faultyPower(totalPower int64) int64 {
	return (totalPower - 1) / 3
}

type RequiredPowerFn func(int64) int64

var Required3FP1Power = func(power int64) int64 {
	f := faultyPower(power)
	p := (3 * f) + 1

	return p
}

var Required2FP1Power = func(power int64) int64 {
	f := faultyPower(power)
	p := (2 * f) + 1

	return p
}

var Required1FP1Power = func(totalPower int64) int64 {
	f := faultyPower(totalPower)
	p := (1 * f) + 1

	return p
}

// Has1FP1Power checks whether the signed power is greater than or equal to f+1,
// where f is the maximum faulty power.
func Has1FP1Power(totalPower, signedPower int64) bool {
	return signedPower >= Required1FP1Power(totalPower)
}

// Has2FP1Power checks whether the signed power is greater than or equal to 2f+1,
// where f is the maximum faulty power.
func Has2FP1Power(totalPower, signedPower int64) bool {
	return signedPower >= Required2FP1Power(totalPower)
}

// Has3FP1Power checks whether the signed power is greater than or equal to 3f+1,
// where f is the maximum faulty power.
func Has3FP1Power(totalPower, signedPower int64) bool {
	return signedPower >= Required3FP1Power(totalPower)
}
