package message

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/crypto/hash"
	"github.com/zarbchain/zarb-go/errors"
)

func TestHeartBeatType(t *testing.T) {
	m := &HeartBeatMessage{}
	assert.Equal(t, m.Type(), MessageTypeHeartBeat)
}

func TestHeartBeatMessage(t *testing.T) {
	t.Run("Invalid height", func(t *testing.T) {
		m := NewHeartBeatMessage(-1, 0, hash.GenerateTestHash())

		assert.Equal(t, errors.Code(m.SanityCheck()), errors.ErrInvalidHeight)
	})

	t.Run("Invalid round", func(t *testing.T) {
		m := NewHeartBeatMessage(100, -1, hash.GenerateTestHash())

		assert.Equal(t, errors.Code(m.SanityCheck()), errors.ErrInvalidRound)
	})

	t.Run("OK", func(t *testing.T) {
		m := NewHeartBeatMessage(100, 1, hash.GenerateTestHash())

		assert.NoError(t, m.SanityCheck())
		assert.Contains(t, m.Fingerprint(), "100")
	})
}
