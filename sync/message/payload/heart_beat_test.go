package payload

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/crypto/hash"
)

func TestHeartBeatType(t *testing.T) {
	p := &HeartBeatPayload{}
	assert.Equal(t, p.Type(), PayloadTypeHeartBeat)
}

func TestHeartBeatPayload(t *testing.T) {
	t.Run("Invalid height", func(t *testing.T) {
		p := NewHeartBeatPayload(-1, 0, hash.GenerateTestHash())

		assert.Error(t, p.SanityCheck())
	})

	t.Run("Invalid round", func(t *testing.T) {
		p := NewHeartBeatPayload(100, -1, hash.GenerateTestHash())

		assert.Error(t, p.SanityCheck())
	})

	t.Run("OK", func(t *testing.T) {
		p := NewHeartBeatPayload(100, 1, hash.GenerateTestHash())

		assert.NoError(t, p.SanityCheck())
		assert.Contains(t, p.Fingerprint(), "100")
	})
}
