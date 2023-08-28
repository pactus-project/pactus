package message

import (
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
		b := ts.GenerateTestBlock(nil, nil)
		c := certificate.NewCertificate(0, 0, nil, nil, nil)
		m := NewBlockAnnounceMessage(100, b, c)
		err := m.BasicCheck()

		assert.ErrorIs(t, err, certificate.BasicCheckError{
			Reason: "height is not positive: 0",
		})
	})

	t.Run("OK", func(t *testing.T) {
		b := ts.GenerateTestBlock(nil, nil)
		c := ts.GenerateTestCertificate()
		m := NewBlockAnnounceMessage(100, b, c)

		assert.NoError(t, m.BasicCheck())
		assert.Contains(t, m.String(), "100")
	})
}
