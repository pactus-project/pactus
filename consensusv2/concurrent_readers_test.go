package consensusv2

import (
	"sync"
	"testing"
)

// TestConcurrentQueries calls the query accessors from several goroutines while
// the log has no round yet, right after a new height starts. Run it with -race.
func TestConcurrentQueries(t *testing.T) {
	td := setup(t)
	td.enterNewHeight(td.consX)

	td.consX.lk.Lock()
	td.consX.log.MoveToNewHeight(td.consX.validators)
	td.consX.lk.Unlock()

	var wg sync.WaitGroup
	for i := 0; i < 4; i++ {
		wg.Go(func() {
			_ = td.consX.Proposal()
			_ = td.consX.AllVotes()
			_ = td.consX.HandleQueryProposal(td.consX.height, td.consX.round)
			_ = td.consX.HandleQueryVote(td.consX.height, td.consX.round)
		})
	}
	wg.Wait()
}
