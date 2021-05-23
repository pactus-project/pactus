package payload

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/crypto"
)

func TestHeartBeatType(t *testing.T) {
	p := &HeartBeatPayload{}
	assert.Equal(t, p.Type(), PayloadTypeHeartBeat)
}

func TestHeartBeatPayload(t *testing.T) {
	t.Run("Invalid height", func(t *testing.T) {
		p1 := NewHeartBeatPayload(-1, 0, crypto.GenerateTestHash())
		assert.Error(t, p1.SanityCheck())
	})

	t.Run("Invalid round", func(t *testing.T) {
		p1 := NewHeartBeatPayload(100, -1, crypto.GenerateTestHash())
		assert.Error(t, p1.SanityCheck())
	})

	t.Run("OK", func(t *testing.T) {
		p1 := NewHeartBeatPayload(100, 1, crypto.GenerateTestHash())
		assert.NoError(t, p1.SanityCheck())
	})
}
