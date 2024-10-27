package score

import (
	"testing"

	"github.com/pactus-project/pactus/types/certificate"
	"github.com/stretchr/testify/assert"
)

func TestScoreManager(t *testing.T) {
	maxCert := uint32(3)
	scoreMgr := NewScoreManager(maxCert)

	cert1 := certificate.NewBlockCertificate(1, 0)
	cert1.SetSignature([]int32{0, 1, 2, 3}, []int32{0}, nil)

	cert2 := certificate.NewBlockCertificate(2, 0)
	cert2.SetSignature([]int32{0, 1, 2, 3}, []int32{3}, nil)

	cert3 := certificate.NewBlockCertificate(3, 0)
	cert3.SetSignature([]int32{1, 2, 3, 4}, []int32{2}, nil)

	cert4 := certificate.NewBlockCertificate(4, 0)
	cert4.SetSignature([]int32{1, 2, 3, 4}, []int32{2}, nil)

	cert5 := certificate.NewBlockCertificate(5, 0)
	cert5.SetSignature([]int32{1, 2, 3, 4}, []int32{2}, nil)

	tests := []struct {
		cert   *certificate.BlockCertificate
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

	for no, tt := range tests {
		scoreMgr.SetCertificate(tt.cert)

		score0 := scoreMgr.AvailabilityScore(0)
		assert.Equal(t, tt.score0, score0, "#%v: invalid score0, expected %v, got %v",
			no, tt.score0, score0)

		score1 := scoreMgr.AvailabilityScore(1)
		assert.Equal(t, tt.score1, score1, "#%v: invalid score1, expected %v, got %v",
			no, tt.score1, score1)

		score2 := scoreMgr.AvailabilityScore(2)
		assert.Equal(t, tt.score2, score2, "#%v: invalid score2, expected %v, got %v",
			no, tt.score2, score2)

		score3 := scoreMgr.AvailabilityScore(3)
		assert.Equal(t, tt.score3, score3, "#%v: invalid score3, expected %v, got %v",
			no, tt.score3, score3)

		score4 := scoreMgr.AvailabilityScore(4)
		assert.Equal(t, tt.score4, score4, "#%v: invalid score4, expected %v, got %v",
			no, tt.score4, score4)
	}
}
