package score

import (
	"testing"

	"github.com/pactus-project/pactus/types/certificate"
	"github.com/stretchr/testify/assert"
)

func TestScoreManager(t *testing.T) {
	maxCert := uint32(3)
	sm := NewScoreManager(maxCert)

	cert1 := certificate.NewCertificate(1, 0, []int32{0, 1, 2, 3}, []int32{0}, nil)
	cert2 := certificate.NewCertificate(2, 0, []int32{0, 1, 2, 3}, []int32{3}, nil)
	cert3 := certificate.NewCertificate(3, 0, []int32{1, 2, 3, 4}, []int32{2}, nil)
	cert4 := certificate.NewCertificate(4, 0, []int32{1, 2, 3, 4}, []int32{2}, nil)
	cert5 := certificate.NewCertificate(5, 0, []int32{1, 2, 3, 4}, []int32{2}, nil)

	tests := []struct {
		cert   *certificate.Certificate
		score0 float64
		score1 float64
		score2 float64
		score3 float64
		score4 float64
	}{
		{cert1, 0, 1, 1, 1, 1},
		{cert2, 0.5, 1, 1, 0.5, 1},
		{cert3, 0.5, 1, 1 - (float64(1) / float64(3)), 1 - (float64(1) / float64(3)), 1},
		{cert4, 1, 1, 1 - (float64(2) / float64(3)), 1 - (float64(1) / float64(3)), 1},
		{cert5, 1, 1, 0, 1, 1},
	}

	for i, test := range tests {
		sm.SetCertificate(test.cert)

		score0 := sm.AvailabilityScore(0)
		assert.Equal(t, test.score0, score0, "#%v: invalid score0, expected %v, got %v",
			i, test.score0, score0)

		score1 := sm.AvailabilityScore(1)
		assert.Equal(t, test.score1, score1, "#%v: invalid score1, expected %v, got %v",
			i, test.score1, score1)

		score2 := sm.AvailabilityScore(2)
		assert.Equal(t, test.score2, score2, "#%v: invalid score2, expected %v, got %v",
			i, test.score2, score2)

		score3 := sm.AvailabilityScore(3)
		assert.Equal(t, test.score3, score3, "#%v: invalid score3, expected %v, got %v",
			i, test.score3, score3)

		score4 := sm.AvailabilityScore(4)
		assert.Equal(t, test.score4, score4, "#%v: invalid score4, expected %v, got %v",
			i, test.score4, score4)
	}
}
