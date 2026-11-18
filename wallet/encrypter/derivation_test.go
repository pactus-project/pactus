package encrypter

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// At most maxConcurrentDerivations keys are derived at the same time: the next one waits.
func TestDerivationsAreBounded(t *testing.T) {
	for range maxConcurrentDerivations {
		derivations <- struct{}{}
	}

	derived := make(chan []byte)
	go func() {
		derived <- deriveKey([]byte("password"), []byte("salt-salt-salt"), 1, 8, 1, 32)
	}()

	select {
	case <-derived:
		require.FailNow(t, "a derivation ran past the bound")
	case <-time.After(200 * time.Millisecond):
	}

	<-derivations
	select {
	case key := <-derived:
		assert.Len(t, key, 32)
	case <-time.After(10 * time.Second):
		require.FailNow(t, "the derivation didn't run once a place was free")
	}

	for range maxConcurrentDerivations - 1 {
		<-derivations
	}
}
