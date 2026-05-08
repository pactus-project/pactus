package message

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHelloAckType(t *testing.T) {
	smg := &HelloAckMessage{}
	assert.Equal(t, TypeHelloAck, smg.Type())
}

func TestHelloAckMessage(t *testing.T) {
	msg := NewHelloAckMessage(ResponseCodeRejected, "rejected", 100)

	require.NoError(t, msg.BasicCheck())
	assert.Zero(t, msg.ConsensusHeight())
	assert.Contains(t, msg.LogString(), "100")
}
