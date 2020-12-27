package consensus

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/consensus/hrs"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/vote"
)

func TestNotAcceptingProposal(t *testing.T) {
	setup(t)

	tConsY.MoveToNewHeight()

	addr := tSigners[tIndexX].Address()
	b, _ := block.GenerateTestBlock(&addr, nil)
	invalidProposal := vote.NewProposal(1, 0, *b)
	tSigners[tIndexX].SignMsg(invalidProposal)
	tConsY.SetProposal(invalidProposal)
	assert.Nil(t, tConsY.LastProposal())
}

//Imagine we have four nodes: (Nx, Ny, Nb, Np) which:
// Nb is a byzantine node and Nx, Ny, Np are honest nodes,
// however Np is partitioned and see the network through Nb (Byzantine node).
// In Height H, B sends its pre-votes to all the nodes
// but only sends valid pre-commit to P and nil pre-commit to X,Y.
// For should not hapens
func TestByzantineVote(t *testing.T) {
	setup(t)

	tConsX.enterNewHeight()
	p := tConsX.LastProposal()

	tConsP.enterNewHeight()
	tConsP.SetProposal(p)

	testAddVote(t, tConsP, vote.VoteTypePrepare, 1, 0, p.Block().Hash(), tIndexX, false)
	testAddVote(t, tConsP, vote.VoteTypePrepare, 1, 0, p.Block().Hash(), tIndexB, false)
	checkHRS(t, tConsP, 1, 0, hrs.StepTypePrecommit)

	testAddVote(t, tConsP, vote.VoteTypePrecommit, 1, 0, p.Block().Hash(), tIndexX, false)
	testAddVote(t, tConsP, vote.VoteTypePrecommit, 1, 0, crypto.UndefHash, tIndexB, false) // Byzantine vote

	shouldPublishVote(t, tConsP, vote.VoteTypePrepare, p.Block().Hash())
	shouldPublishVote(t, tConsP, vote.VoteTypePrecommit, p.Block().Hash())

	// Partition heals
	testAddVote(t, tConsP, vote.VoteTypePrecommit, 1, 0, p.Block().Hash(), tIndexY, false)
	checkHRS(t, tConsP, 1, 0, hrs.StepTypeCommit)
}
