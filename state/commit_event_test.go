package state

import (
	"testing"
	"time"

	"github.com/pactus-project/gopkg/pipeline"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

// TestCommitBlockWithSlowStateReader makes sure an event receiver that reads the state,
// like the wallet does, cannot deadlock the commit when the event pipe is full.
func TestCommitBlockWithSlowStateReader(t *testing.T) {
	td := setup(t)

	eventPipe := pipeline.New[any](t.Context(), pipeline.WithBufferSize(1))
	eventPipe.RegisterReceiver(func(any) {
		time.Sleep(10 * time.Millisecond)
		_ = td.state.LastBlockHeight()
	})
	td.state.eventPipe = eventPipe
	td.fakeTxPool.EXPECT().AppendTxAndBroadcast(gomock.Any()).Return(nil).AnyTimes()

	errCh := make(chan error, 1)
	go func() {
		for i := 0; i < 10; i++ {
			td.fakeTxPool.EXPECT().PrepareBlockTransactions().Return(nil).Times(1)
			blk, cert := td.makeBlockAndCertificate(t, 0)
			if err := td.state.CommitBlock(blk, cert); err != nil {
				errCh <- err

				return
			}
		}
		errCh <- nil
	}()

	select {
	case err := <-errCh:
		require.NoError(t, err)
	case <-time.After(10 * time.Second):
		require.FailNow(t, "committing blocks is deadlocked")
	}
}
