package message

import (
	"testing"

	"github.com/pactus-project/pactus/types/block"
	"github.com/pactus-project/pactus/util/errors"
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
		b := ts.GenerateTestBlock(nil, nil)
		c := block.NewCertificate(-1, nil, nil, nil)
		m := NewBlockAnnounceMessage(100, b, c)

		assert.Equal(t, errors.Code(m.SanityCheck()), errors.ErrInvalidRound)
	})

	t.Run("OK", func(t *testing.T) {
		b := ts.GenerateTestBlock(nil, nil)
		c := ts.GenerateTestCertificate(b.Hash())
		m := NewBlockAnnounceMessage(100, b, c)

		assert.NoError(t, m.SanityCheck())
		assert.Contains(t, m.String(), "100")
	})
}
