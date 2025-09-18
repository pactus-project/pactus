package certificate_test

import (
	"testing"

	"github.com/pactus-project/pactus/types/certificate"
	"github.com/stretchr/testify/assert"
)

func TestPowerFunctions(t *testing.T) {
	tests := []struct {
		totalPower   int64
		signedPower  int64
		has1FP1Power bool
		has2FP1Power bool
		has3FP1Power bool
	}{
		{totalPower: 7, signedPower: 2, has1FP1Power: false, has2FP1Power: false, has3FP1Power: false},
		{totalPower: 7, signedPower: 3, has1FP1Power: true, has2FP1Power: false, has3FP1Power: false},
		{totalPower: 7, signedPower: 5, has1FP1Power: true, has2FP1Power: true, has3FP1Power: false},
		{totalPower: 7, signedPower: 7, has1FP1Power: true, has2FP1Power: true, has3FP1Power: true},

		{totalPower: 8, signedPower: 2, has1FP1Power: false, has2FP1Power: false, has3FP1Power: false},
		{totalPower: 8, signedPower: 3, has1FP1Power: true, has2FP1Power: false, has3FP1Power: false},
		{totalPower: 8, signedPower: 5, has1FP1Power: true, has2FP1Power: true, has3FP1Power: false},
		{totalPower: 8, signedPower: 7, has1FP1Power: true, has2FP1Power: true, has3FP1Power: true},

		{totalPower: 9, signedPower: 2, has1FP1Power: false, has2FP1Power: false, has3FP1Power: false},
		{totalPower: 9, signedPower: 3, has1FP1Power: true, has2FP1Power: false, has3FP1Power: false},
		{totalPower: 9, signedPower: 5, has1FP1Power: true, has2FP1Power: true, has3FP1Power: false},
		{totalPower: 9, signedPower: 7, has1FP1Power: true, has2FP1Power: true, has3FP1Power: true},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.has1FP1Power, certificate.Has1FP1Power(tt.totalPower, tt.signedPower),
			"Has1FP1Power(%d, %d)", tt.totalPower, tt.signedPower)

		assert.Equal(t, tt.has2FP1Power, certificate.Has2FP1Power(tt.totalPower, tt.signedPower),
			"Has2FP1Power(%d, %d)", tt.totalPower, tt.signedPower)

		assert.Equal(t, tt.has3FP1Power, certificate.Has3FP1Power(tt.totalPower, tt.signedPower),
			"Has3FP1Power(%d, %d)", tt.totalPower, tt.signedPower)
	}
}
