package message

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHelloAckType(t *testing.T) {
	m := &HelloAckMessage{}
	assert.Equal(t, m.Type(), TypeHelloAck)
}

func TestHelloAckMessage(t *testing.T) {
	m := NewHelloAckMessage(ResponseCodeRejected, "rejected", 0)
	assert.NoError(t, m.BasicCheck())
}
