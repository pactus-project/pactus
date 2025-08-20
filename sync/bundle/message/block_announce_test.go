package message

import (
	"fmt"
	"testing"

	"github.com/pactus-project/pactus/types/certificate"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
)

func TestBlockAnnounceType(t *testing.T) {
	smg := &BlockAnnounceMessage{}
	assert.Equal(t, TypeBlockAnnounce, smg.Type())
}

func TestBlockAnnounceMessage(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	t.Run("Invalid block certificate", func(t *testing.T) {
		blk, _ := ts.GenerateTestBlock(ts.RandHeight())
		cert := certificate.NewBlockCertificate(0, 0)
		msg := NewBlockAnnounceMessage(blk, cert, nil)
		err := msg.BasicCheck()

		assert.ErrorIs(t, err, certificate.BasicCheckError{
			Reason: "height is not positive: 0",
		})
	})

	t.Run("Invalid proof", func(t *testing.T) {
		height := ts.RandHeight()
		blk, cert := ts.GenerateTestBlock(height)
		proof := ts.GenerateTestVoteCertificate(0)
		msg := NewBlockAnnounceMessage(blk, cert, proof)
		err := msg.BasicCheck()

		assert.ErrorIs(t, err, certificate.BasicCheckError{
			Reason: "height is not positive: 0",
		})
	})

	t.Run("OK", func(t *testing.T) {
		height := ts.RandHeight()
		blk, cert := ts.GenerateTestBlock(height)
		msg := NewBlockAnnounceMessage(blk, cert, nil)

		assert.NoError(t, msg.BasicCheck())
		assert.Equal(t, height, msg.ConsensusHeight())
		assert.Contains(t, msg.String(), fmt.Sprintf("%d", height))
	})
}
