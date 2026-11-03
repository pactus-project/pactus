package state

import (
	"testing"

	"github.com/pactus-project/pactus/types/block"
	"github.com/pactus-project/pactus/types/vote"
	"github.com/stretchr/testify/require"
)

// A late precommit vote of an absent validator is merged into the certificate
// of the last block. The vote must be a precommit for the last block; a precommit
// for another block hash (equivocation or garbage from a validator that is
// absent in the certificate) must not corrupt the stored certificate: this
// certificate is embedded in the next block, which would become invalid.
func TestUpdateLastCertificateIgnoresOtherBlockHash(t *testing.T) {
	td := setup(t)

	td.fakeTxPool.EXPECT().PrepareBlockTransactions().Return(block.Txs{}).Times(1)
	blk, cert := td.makeBlockAndCertificate(t, 0)
	require.NoError(t, td.state.CommitBlock(blk, cert))

	absentee := td.genValKeys[len(td.genValKeys)-1]
	vte := vote.NewPrecommitVote(td.RandHash(), cert.Height(), cert.Round(), absentee.Address())
	td.HelperSignVote(absentee, vte)

	_ = td.state.UpdateLastCertificate(vte)

	lastCert := td.state.LastCertificate()
	err := lastCert.ValidatePrecommit(td.state.lastInfo.Validators(), td.state.LastBlockHash())
	require.NoError(t, err, "the certificate of the last block is corrupted by a vote for another block")
}

// Consequence: the next block of every proposer embeds the corrupted certificate.
func TestProposeBlockAfterGarbageLateVote(t *testing.T) {
	td := setup(t)

	td.fakeTxPool.EXPECT().PrepareBlockTransactions().Return(block.Txs{}).AnyTimes()
	blk, cert := td.makeBlockAndCertificate(t, 0)
	require.NoError(t, td.state.CommitBlock(blk, cert))

	absentee := td.genValKeys[len(td.genValKeys)-1]
	vte := vote.NewPrecommitVote(td.RandHash(), cert.Height(), cert.Round(), absentee.Address())
	td.HelperSignVote(absentee, vte)
	_ = td.state.UpdateLastCertificate(vte)

	next, err := td.state.ProposeBlock(td.proposerKey(t, 0), td.RandAccAddress())
	require.NoError(t, err)
	require.NoError(t, td.state.ValidateBlock(next, 0))
}
