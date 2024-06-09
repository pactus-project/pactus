package ntp

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewNtpChecker(t *testing.T) {
	checker := NewNtpChecker(1*time.Second, 500*time.Millisecond)
	assert.NotNil(t, checker)
	assert.Equal(t, 1*time.Second, checker.interval)
	assert.Equal(t, 500*time.Millisecond, checker.threshold)
}

func TestOutOfSync(t *testing.T) {
	checker := NewNtpChecker(1*time.Second, 500*time.Millisecond)
	assert.False(t, checker.OutOfSync(300*time.Millisecond))
	assert.True(t, checker.OutOfSync(600*time.Millisecond))
}

func TestGetClockOffset(t *testing.T) {
	checker := NewNtpChecker(1*time.Second, 500*time.Millisecond)
	checker.lock.Lock()
	checker.offset = 300 * time.Millisecond
	checker.lock.Unlock()

	offset, err := checker.GetClockOffset()
	assert.NoError(t, err)
	assert.Equal(t, 300*time.Millisecond, offset)

	checker.lock.Lock()
	checker.offset = maxClockOffset
	checker.lock.Unlock()

	offset, err = checker.GetClockOffset()
	assert.Error(t, err)
	assert.Equal(t, 0*time.Millisecond, offset)
}

func TestStartAndStop(t *testing.T) {
	checker := NewNtpChecker(1*time.Second, 500*time.Millisecond)
	assert.NotNil(t, checker)

	// Using a WaitGroup to wait for the goroutine to start
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		wg.Done()
		checker.Start()
	}()

	// Wait for the goroutine to start
	wg.Wait()

	// Let the checker run for a short period
	time.Sleep(2 * time.Second)

	checker.Stop()
}
