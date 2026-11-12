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

// A block is committed between the verification of a vote for the last certificate and its addition:
// the vote is for the previous block now, and it is not added to the certificate of the new block.
func TestVoteForLastCertificateAfterNewBlock(t *testing.T) {
	td := setupWithBlocks(t, protocol.ProtocolVersionLatest, 2)

	// The last validator is absent from the certificates that the tests make.
	absentee := td.genValKeys[3]
	lastCert := td.state.LastCertificate()
	vte := vote.NewPrecommitVote(td.state.LastBlockHash(), lastCert.Height(), lastCert.Round(), absentee.Address())
	td.HelperSignVote(absentee, vte)
	require.NoError(t, td.state.verifyVoteForLastCertificate(vte))

	td.commit(t, td.proposeEmpty(t))

	require.ErrorIs(t, td.state.addVoteToLastCertificate(vte), InvalidVoteForCertificateError{Vote: vte})
	require.NoError(t, td.state.LastCertificate().ValidatePrecommit(td.state.lastInfo.Validators(),
		td.state.LastBlockHash()))
}
