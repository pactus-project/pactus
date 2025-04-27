package pipeline

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetName(t *testing.T) {
	p := New[int](context.TODO(), "test", 10)
	assert.Equal(t, "test", p.Name())
}

func TestClosePipeline(t *testing.T) {
	p := New[int](context.TODO(), "test", 10)

	p.RegisterReceiver(func(a int) {
		time.Sleep(time.Duration(a) * time.Millisecond)
	})

	for i := 0; i < 10; i++ {
		go p.Send(i)
	}

	p.Close()
	p.Close()

	assert.Eventually(t, p.IsClosed, 500*time.Millisecond, 10*time.Millisecond)
}

func TestSendReceive(t *testing.T) {
	p := New[float64](context.TODO(), "test", 10)

	received := make(chan float64, 1)
	receiver := func(data float64) {
		received <- data
	}

	p.RegisterReceiver(receiver)

	// Test multiple messages
	p.Send(3.141592)
	p.Send(2.718281)
	p.Send(1.618033)

	assert.Equal(t, 3.141592, <-received)
	assert.Equal(t, 2.718281, <-received)
	assert.Equal(t, 1.618033, <-received)
}

// TestSendAfterClose verifies error handling
func TestSendAfterClose(t *testing.T) {
	p := New[string](context.TODO(), "test", 10)

	// Close the pipeline first
	p.Close()

	// Send should fail gracefully
	assert.NotPanics(t, func() {
		p.Send("should not panic")
	})
}

func TestContextCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	p := New[string](ctx, "test", 10)

	// Cancel the context
	cancel()

	// Send should fail gracefully
	assert.NotPanics(t, func() {
		p.Send("should not panic after cancel")
	})
}

// TestReceiveError simulates error in receiver
// func TestReceiveError(t *testing.T) {
// 	p := NewPipeline(context.TODO(), "test", 10)

// 	errCh := make(chan error, 1)
// 	p.OnReceive(func(data any) {
// 		if data == "error" {
// 			errCh <- errors.New("mock error")
// 		}
// 	})

// 	// Trigger error case
// 	p.Send("error")

// 	select {
// 	case err := <-errCh:
// 		assert.EqualError(t, err, "mock error")
// 	case <-time.After(100 * time.Millisecond):
// 		t.Fatal("receiver did not process message")
// 	}
// }
