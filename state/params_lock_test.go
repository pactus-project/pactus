package state

import (
	"testing"

	"github.com/pactus-project/pactus/types"
	"github.com/pactus-project/pactus/types/protocol"
	"github.com/stretchr/testify/assert"
)

// A Hello message reads the parameters while a block is being committed.
func TestParamsDuringCommit(t *testing.T) {
	td := setupWithBlocks(t, protocol.ProtocolVersionLatest, 0)

	stop := make(chan struct{})
	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			select {
			case <-stop:
				return
			default:
				_ = td.state.Params().BlockVersion
			}
		}
	}()

	td.commitBlocks(t, 5)
	close(stop)
	<-done

	assert.Equal(t, types.Height(5), td.state.LastBlockHeight())
}
