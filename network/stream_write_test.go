package network

import (
	"errors"
	"testing"
	"time"

	lp2pnetwork "github.com/libp2p/go-libp2p/core/network"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// failingStream is a stream whose write or close for writing fails.
type failingStream struct {
	lp2pnetwork.Stream

	writeErr      error
	closeWriteErr error
	written       []byte
	closedWrite   bool
	reset         bool
}

func (*failingStream) SetDeadline(time.Time) error { return nil }

func (s *failingStream) Write(data []byte) (int, error) {
	if s.writeErr != nil {
		return 0, s.writeErr
	}
	s.written = append(s.written, data...)

	return len(data), nil
}

func (s *failingStream) CloseWrite() error {
	s.closedWrite = true

	return s.closeWriteErr
}

func (s *failingStream) Reset() error {
	s.reset = true

	return nil
}

func TestWriteMessage(t *testing.T) {
	service := &streamService{timeout: time.Second}
	errTest := errors.New("test")

	t.Run("the message is written and the stream is closed for writing", func(t *testing.T) {
		stream := &failingStream{}
		require.NoError(t, service.writeMessage(stream, []byte("message")))
		assert.Equal(t, []byte("message"), stream.written)
		assert.True(t, stream.closedWrite)
		assert.False(t, stream.reset)
	})

	t.Run("a stream that can't be written is reset", func(t *testing.T) {
		stream := &failingStream{writeErr: errTest}
		require.ErrorAs(t, service.writeMessage(stream, []byte("message")), &LibP2PError{})
		assert.True(t, stream.reset)
		assert.False(t, stream.closedWrite)
	})

	t.Run("a stream that can't be half-closed is reset", func(t *testing.T) {
		stream := &failingStream{closeWriteErr: errTest}
		require.ErrorAs(t, service.writeMessage(stream, []byte("message")), &LibP2PError{})
		assert.True(t, stream.reset)
	})
}
