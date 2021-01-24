package consensus

import (
	"testing"

	"github.com/zarbchain/zarb-go/consensus/hrs"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/vote"
)

// Imagine we have four nodes: (Nx, Ny, Nb, Np) which:
// Nb is a byzantine node and Nx, Ny, Np are honest nodes,
// however Np is partitioned and see the network through Nb (Byzantine node).
// In Height H, B sends its pre-votes to all the nodes
// but only sends valid pre-commit to P and nil pre-commit to X,Y.
// For should not hapens
func TestByzantineVote(t *testing.T) {
	setup(t)

	h := 1
	r := 0
	p := makeProposal(t, h, r)

	tConsP.enterNewHeight()
	tConsP.SetProposal(p)

	testAddVote(t, tConsP, vote.VoteTypePrepare, h, r, p.Block().Hash(), tIndexX, false)
	testAddVote(t, tConsP, vote.VoteTypePrepare, h, r, p.Block().Hash(), tIndexB, false)
	checkHRS(t, tConsP, h, r, hrs.StepTypePrecommit)

	testAddVote(t, tConsP, vote.VoteTypePrecommit, h, r, p.Block().Hash(), tIndexX, false)
	testAddVote(t, tConsP, vote.VoteTypePrecommit, h, r, crypto.UndefHash, tIndexB, false) // Byzantine vote

	shouldPublishVote(t, tConsP, vote.VoteTypePrepare, p.Block().Hash())
	shouldPublishVote(t, tConsP, vote.VoteTypePrecommit, p.Block().Hash())

	// Partition heals
	testAddVote(t, tConsP, vote.VoteTypePrecommit, h, r, p.Block().Hash(), tIndexY, false)
	checkHRS(t, tConsP, h, r, hrs.StepTypeCommit)
}

func TestPrepareTimeout(t *testing.T) {
	setup(t)

	commitBlockForAllStates(t)
	commitBlockForAllStates(t)

	tConsY.enterNewHeight()

	shouldPublishVote(t, tConsY, vote.VoteTypePrepare, crypto.UndefHash)
}
