package consensus

import (
	"testing"
	"time"

	"github.com/zarbchain/zarb-go/consensus/vote"
	"github.com/zarbchain/zarb-go/types/crypto/hash"
)

func TestPrepareQueryProposal(t *testing.T) {
	setup(t)

	commitBlockForAllStates(t)

	testEnterNewHeight(tConsP)

	// After receiving one vote, it should query for proposal (if don't have it yet)
	testAddVote(tConsP, vote.VoteTypePrepare, 2, 0, hash.GenerateTestHash(), tIndexX)

	shouldPublishQueryProposal(t, tConsP, 2, 0)
}

func TestGoToChangeProposerFromPrepare(t *testing.T) {
	setup(t)

	commitBlockForAllStates(t)

	testEnterNewHeight(tConsP)
	p := makeProposal(t, 2, 0)

	testAddVote(tConsP, vote.VoteTypeChangeProposer, 2, 0, hash.UndefHash, tIndexX)
	testAddVote(tConsP, vote.VoteTypeChangeProposer, 2, 0, hash.UndefHash, tIndexY)

	tConsP.SetProposal(p)
	shouldPublishVote(t, tConsP, vote.VoteTypeChangeProposer, hash.UndefHash)
}

// We have four nodes: Nx, Ny, Nb, Np, which:
// Nb is a byzantine node and Nx, Ny, Np are honest nodes,
// however Np is partitioned and see the network through Nb (Byzantine node).
//
// In Height H, B sends prepare votes to Nx, Ny and change-proposer vote to Np.
// Np should not move to change-proposer stage unless it has 1/3+ votes from other replicas.
func TestByzantineVote1(t *testing.T) {
	setup(t)

	commitBlockForAllStates(t)

	h := int32(2)
	r := int16(0)
	p := makeProposal(t, h, r)

	testEnterNewHeight(tConsX)
	testEnterNewHeight(tConsP)

	// =================================
	// Nx votes
	testAddVote(tConsX, vote.VoteTypePrepare, h, r, p.Block().Hash(), tIndexX)
	testAddVote(tConsX, vote.VoteTypePrepare, h, r, p.Block().Hash(), tIndexY)
	testAddVote(tConsX, vote.VoteTypePrepare, h, r, p.Block().Hash(), tIndexB)

	testAddVote(tConsX, vote.VoteTypePrecommit, h, r, p.Block().Hash(), tIndexX)
	testAddVote(tConsX, vote.VoteTypePrecommit, h, r, p.Block().Hash(), tIndexY)
	// Byzantine node doesn't broadcast its precommit vote

	// Nx, Ny are unable to progress

	// =================================
	// Np votes
	testAddVote(tConsP, vote.VoteTypeChangeProposer, h, r, hash.UndefHash, tIndexB) // Byzantine vote
	shouldPublishVote(t, tConsP, vote.VoteTypeChangeProposer, hash.UndefHash)
	// Np is unable to progress

	// =================================
	time.Sleep(1 * time.Second)
	checkHeightRound(t, tConsX, h, r)
	checkHeightRound(t, tConsP, h, r)

	// =================================
	// Now, Partition heals

	for _, v := range tConsP.AllVotes() {
		tConsX.AddVote(v)
	}
	shouldPublishVote(t, tConsX, vote.VoteTypeChangeProposer, hash.UndefHash)
	checkHeightRoundWait(t, tConsX, h, r+1)

	tConsP.SetProposal(p)
	for _, v := range tConsX.AllVotes() {
		tConsP.AddVote(v)
	}
	checkHeightRoundWait(t, tConsP, h, r+1)
}

// Np should propose a block. Np is partitioned and Nb doesn't send proposal to Nx, Ny.
func TestByzantineVote2(t *testing.T) {
	setup(t)

	commitBlockForAllStates(t)
	commitBlockForAllStates(t)
	commitBlockForAllStates(t)

	h := int32(4)
	r := int16(0)
	p1 := makeProposal(t, h, r)
	p2 := makeProposal(t, h, r+1)

	testEnterNewHeight(tConsX)
	testEnterNewHeight(tConsP)

	// =================================
	// Np votes
	testAddVote(tConsP, vote.VoteTypePrepare, h, r, p1.Block().Hash(), tIndexB) // Byzantine vote
	shouldPublishVote(t, tConsP, vote.VoteTypePrepare, p1.Block().Hash())

	// Partitioned node is unable to progress

	// =================================
	// Nx votes
	testAddVote(tConsX, vote.VoteTypeChangeProposer, h, r, hash.UndefHash, tIndexY)
	testAddVote(tConsX, vote.VoteTypeChangeProposer, h, r, hash.UndefHash, tIndexB) // Nb sends change proposer vote to Nx, Ny

	shouldPublishVote(t, tConsX, vote.VoteTypeChangeProposer, hash.UndefHash)
	// Nx goes to the next round

	testAddVote(tConsX, vote.VoteTypePrepare, h, r+1, p2.Block().Hash(), tIndexY)
	shouldPublishVote(t, tConsX, vote.VoteTypePrepare, p2.Block().Hash())

	// Nx, Ny are unable to progress

	// =================================
	time.Sleep(1 * time.Second)
	checkHeightRound(t, tConsX, h, r+1)
	checkHeightRound(t, tConsP, h, r)

	// =================================
	// Now, Partition heals
	tConsP.SetProposal(p2)
	for _, v := range tConsX.AllVotes() {
		tConsP.AddVote(v)
	}
	shouldPublishVote(t, tConsP, vote.VoteTypeChangeProposer, hash.UndefHash)
	checkHeightRoundWait(t, tConsP, h, r+1)

	for _, v := range tConsP.AllVotes() {
		tConsX.AddVote(v)
	}
	checkHeightRoundWait(t, tConsX, h, r+1)
}

func TestQueryProposal(t *testing.T) {
	setup(t)

	commitBlockForAllStates(t)
	testEnterNewHeight(tConsX)
	testEnterNextRound(tConsX)
	shouldPublishQueryProposal(t, tConsX, 2, 1)
}
