package message

import (
	"testing"

	"github.com/pactus-project/pactus/util/errors"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
)

func TestHeartBeatType(t *testing.T) {
	m := &HeartBeatMessage{}
	assert.Equal(t, m.Type(), TypeHeartBeat)
}

func TestHeartBeatMessage(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	t.Run("Invalid height", func(t *testing.T) {
		m := NewHeartBeatMessage(0, 0, ts.RandHash())

		assert.Equal(t, errors.Code(m.BasicCheck()), errors.ErrInvalidHeight)
	})

	t.Run("OK", func(t *testing.T) {
		m := NewHeartBeatMessage(100, 1, ts.RandHash())

		assert.NoError(t, m.BasicCheck())
		assert.Contains(t, m.String(), "100")
	})
}
