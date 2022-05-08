package consensus

import (
	"sync"
	"testing"
	"time"

	"github.com/zarbchain/zarb-go/types/crypto/hash"
	"github.com/zarbchain/zarb-go/types/vote"
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
// Np should not move to change-proposer state.
// After partition heals, they move to next round.
func TestByzantineVote1(t *testing.T) {
	setup(t)

	commitBlockForAllStates(t)
	commitBlockForAllStates(t)

	h := int32(3)
	r := int16(0)
	p := makeProposal(t, h, r)

	wg := sync.WaitGroup{}
	wg.Add(2)

	// =================================
	// Nx votes
	go func() {
		testEnterNewHeight(tConsX)
		tConsX.SetProposal(p)

		shouldPublishVote(t, tConsX, vote.VoteTypePrepare, p.Block().Hash())
		testAddVote(tConsX, vote.VoteTypePrepare, h, r, p.Block().Hash(), tIndexY)
		testAddVote(tConsX, vote.VoteTypePrepare, h, r, p.Block().Hash(), tIndexB)

		shouldPublishVote(t, tConsX, vote.VoteTypePrecommit, p.Block().Hash())
		testAddVote(tConsX, vote.VoteTypePrecommit, h, r, p.Block().Hash(), tIndexY)
		// Byzantine node doesn't broadcast its precommit vote

		// Nx and Ny are unable to progress

		wg.Done()
	}()

	// =================================
	// Np votes
	go func() {
		testEnterNewHeight(tConsP)

		shouldPublishVote(t, tConsP, vote.VoteTypeChangeProposer, hash.UndefHash)
		testAddVote(tConsP, vote.VoteTypeChangeProposer, h, r, hash.UndefHash, tIndexB) // Byzantine vote

		wg.Done()
	}()
	// Np is unable to progress

	// =================================
	wg.Wait()
	checkHeightRound(t, tConsX, h, r)
	checkHeightRound(t, tConsP, h, r)

	// =================================
	// Now, Partition heals
	time.Sleep(1 * time.Second)

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

// We have four nodes: Nx, Ny, Nb, Np, which:
// Nb is a byzantine node and Nx, Ny, Np are honest nodes,
// however Np is partitioned and see the network through Nb (Byzantine node).
//
// In Height H, B sends change-proposer votes to Nx, Ny and prepare vote to Np.
// Np moves to precommit state.
// After partition heals, they move to next round.
func TestByzantineVote2(t *testing.T) {
	setup(t)

	commitBlockForAllStates(t)
	commitBlockForAllStates(t)

	h := int32(3)
	r := int16(0)
	p := makeProposal(t, h, r)

	wg := sync.WaitGroup{}
	wg.Add(2)

	// =================================
	// Nx votes
	go func() {
		testEnterNewHeight(tConsX)

		tConsX.SetProposal(p)
		shouldPublishVote(t, tConsX, vote.VoteTypePrepare, p.Block().Hash())
		testAddVote(tConsX, vote.VoteTypePrepare, h, r, p.Block().Hash(), tIndexY)
		testAddVote(tConsX, vote.VoteTypePrepare, h, r, hash.UndefHash, tIndexB)

		shouldPublishVote(t, tConsX, vote.VoteTypeChangeProposer, hash.UndefHash)
		testAddVote(tConsX, vote.VoteTypeChangeProposer, h, r, hash.UndefHash, tIndexY)
		testAddVote(tConsX, vote.VoteTypeChangeProposer, h, r, hash.UndefHash, tIndexB)

		// Nx and Ny move to next round

		wg.Done()
	}()
	// =================================
	// Np votes
	go func() {
		testEnterNewHeight(tConsP)

		tConsP.SetProposal(p)
		shouldPublishVote(t, tConsP, vote.VoteTypePrepare, p.Block().Hash())
		testAddVote(tConsP, vote.VoteTypePrepare, h, r, p.Block().Hash(), tIndexX)
		testAddVote(tConsP, vote.VoteTypePrepare, h, r, p.Block().Hash(), tIndexB)

		shouldPublishVote(t, tConsP, vote.VoteTypePrecommit, p.Block().Hash())
		testAddVote(tConsP, vote.VoteTypePrecommit, h, r, p.Block().Hash(), tIndexB)
		// Np is unable to progress

		wg.Done()
	}()

	// =================================
	wg.Wait()
	checkHeightRound(t, tConsX, h, r+1)
	checkHeightRound(t, tConsP, h, r)

	// =================================
	// Now, Partition heals

	time.Sleep(1 * time.Second)

	for _, v := range tConsP.AllVotes() {
		tConsX.AddVote(v)
	}
	checkHeightRoundWait(t, tConsX, h, r+1)

	tConsP.SetProposal(p)
	for _, v := range tConsX.AllVotes() {
		tConsP.AddVote(v)
	}
	shouldPublishVote(t, tConsP, vote.VoteTypeChangeProposer, hash.UndefHash)
	checkHeightRoundWait(t, tConsP, h, r+1)
}

func TestQueryProposal(t *testing.T) {
	setup(t)

	commitBlockForAllStates(t)
	testEnterNewHeight(tConsX)
	testEnterNextRound(tConsX)
	shouldPublishQueryProposal(t, tConsX, 2, 1)
}
