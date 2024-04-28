package message

import (
	"fmt"
	"testing"

	"github.com/pactus-project/pactus/types/certificate"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
)

func TestBlockAnnounceType(t *testing.T) {
	m := &BlockAnnounceMessage{}
	assert.Equal(t, m.Type(), TypeBlockAnnounce)
}

func TestBlockAnnounceMessage(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	t.Run("Invalid certificate", func(t *testing.T) {
		blk, _ := ts.GenerateTestBlock(ts.RandHeight())
		cert := certificate.NewBlockCertificate(0, 0, false)
		m := NewBlockAnnounceMessage(blk, cert)
		err := m.BasicCheck()

		assert.ErrorIs(t, err, certificate.BasicCheckError{
			Reason: "height is not positive: 0",
		})
	})

	t.Run("OK", func(t *testing.T) {
		height := ts.RandHeight()
		blk, cert := ts.GenerateTestBlock(height)
		m := NewBlockAnnounceMessage(blk, cert)

		assert.NoError(t, m.BasicCheck())
		assert.Equal(t, height, m.Height())
		assert.Contains(t, m.String(), fmt.Sprintf("%d", height))
	})
}
