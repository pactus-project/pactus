package state

import (
	"testing"
	"time"

	"github.com/pactus-project/pactus/types/protocol"
	"github.com/pactus-project/pactus/types/vote"
	"github.com/stretchr/testify/require"
)

// A vote that can't update the last certificate is refused with the read lock,
// so a burst of them doesn't contend with the commit of blocks.
func TestUpdateLastCertificateInvalidVoteDoesNotTakeTheWriteLock(t *testing.T) {
	td := setupWithBlocks(t, protocol.ProtocolVersionLatest, 2)
	signer := td.genValKeys[3]

	invalid := []*vote.Vote{
		vote.NewPrepareVote(td.RandHash(), td.state.LastBlockHeight(), 0, signer.Address()),
		vote.NewPrecommitVote(td.RandHash(), td.state.LastBlockHeight()-1, 0, signer.Address()),
		vote.NewPrecommitVote(td.RandHash(), td.state.LastBlockHeight(),
			td.state.LastCertificate().Round()+1, signer.Address()),
	}

	// A reader holds the lock: the write lock would block.
	td.state.lk.RLock()
	defer td.state.lk.RUnlock()

	for _, vte := range invalid {
		td.HelperSignVote(signer, vte)

		done := make(chan error, 1)
		go func() { done <- td.state.UpdateLastCertificate(vte) }()

		select {
		case err := <-done:
			require.ErrorIs(t, err, InvalidVoteForCertificateError{Vote: vte})
		case <-time.After(3 * time.Second):
			require.Fail(t, "the vote should be refused without the write lock")
		}
	}
}
