package consensus

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/types/vote"
)

func TestPrepareQueryProposal(t *testing.T) {
	td := setup(t)

	td.commitBlockForAllStates(t)

	td.enterNewHeight(td.consP)

	// After receiving one vote, it should query for proposal (if don't have it yet)
	td.addVote(td.consP, vote.VoteTypePrepare, 2, 0, td.RandomHash(), tIndexX)

	td.shouldPublishQueryProposal(t, td.consP, 2, 0)
}

func TestGoToChangeProposerFromPrepare(t *testing.T) {
	td := setup(t)

	td.commitBlockForAllStates(t)

	td.enterNewHeight(td.consP)
	p := td.makeProposal(t, 2, 0)

	td.addVote(td.consP, vote.VoteTypeChangeProposer, 2, 0, hash.UndefHash, tIndexX)
	td.addVote(td.consP, vote.VoteTypeChangeProposer, 2, 0, hash.UndefHash, tIndexY)

	td.consP.SetProposal(p)
	td.shouldPublishVote(t, td.consP, vote.VoteTypeChangeProposer, hash.UndefHash)
}

// We have four nodes: Nx, Ny, Nb, and Np, which:
// - Nb is a Byzantine node
// - Nx, Ny, and Np are honest nodes
// - However, Np is partitioned and sees the network through Nb (Byzantine node).
//
// In Height H, Nb sends prepare votes to Nx, Ny, and a change-proposer vote to Np.
// Np should not move to the change-proposer state.
// After the partition heals, honest nodes move to the next round.
func TestByzantineVote1(t *testing.T) {
	td := setup(t)

	td.commitBlockForAllStates(t)
	td.commitBlockForAllStates(t)

	h := uint32(3)
	r := int16(0)
	p := td.makeProposal(t, h, r)

	wg := sync.WaitGroup{}
	wg.Add(2)

	// =================================
	// Nx votes
	go func() {
		td.enterNewHeight(td.consX)
		td.consX.SetProposal(p)

		td.shouldPublishVote(t, td.consX, vote.VoteTypePrepare, p.Block().Hash())
		td.addVote(td.consX, vote.VoteTypePrepare, h, r, p.Block().Hash(), tIndexY)
		td.addVote(td.consX, vote.VoteTypePrepare, h, r, p.Block().Hash(), tIndexB)

		td.shouldPublishVote(t, td.consX, vote.VoteTypePrecommit, p.Block().Hash())
		td.addVote(td.consX, vote.VoteTypePrecommit, h, r, p.Block().Hash(), tIndexY)
		// Byzantine node doesn't broadcast its precommit vote

		// Nx and Ny are unable to progress

		wg.Done()
	}()

	// =================================
	// Np votes
	go func() {
		td.enterNewHeight(td.consP)

		td.shouldPublishVote(t, td.consP, vote.VoteTypeChangeProposer, hash.UndefHash)
		td.addVote(td.consP, vote.VoteTypeChangeProposer, h, r, hash.UndefHash, tIndexB) // Byzantine vote

		wg.Done()
	}()
	// Np is unable to progress

	// =================================
	wg.Wait()
	td.checkHeightRound(t, td.consX, h, r)
	td.checkHeightRound(t, td.consP, h, r)

	// =================================
	// Now, Partition heals
	time.Sleep(1 * time.Second)
	fmt.Println("=== Partition healed")

	for _, v := range td.consP.AllVotes() {
		td.consX.AddVote(v)
	}
	td.shouldPublishVote(t, td.consX, vote.VoteTypeChangeProposer, hash.UndefHash)
	td.checkHeightRoundWait(t, td.consX, h, r+1)

	td.consP.SetProposal(p)
	for _, v := range td.consX.AllVotes() {
		td.consP.AddVote(v)
	}
	td.checkHeightRoundWait(t, td.consP, h, r+1)
}

// We have four nodes: Nx, Ny, Nb, and Np, which:
// - Nb is a Byzantine node
// - Nx, Ny, and Np are honest nodes
// - However, Np is partitioned and sees the network through Nb (Byzantine node).
//
// In Height H, Nb sends change-proposer votes to Nx, Ny, and a prepare vote to Np.
// Np moves to the precommit state.
// After the partition heals, honest nodes move to the next round.
func TestByzantineVote2(t *testing.T) {
	td := setup(t)

	td.commitBlockForAllStates(t)
	td.commitBlockForAllStates(t)

	h := uint32(3)
	r := int16(0)
	p := td.makeProposal(t, h, r)

	wg := sync.WaitGroup{}
	wg.Add(2)

	// =================================
	// Nx votes
	go func() {
		td.enterNewHeight(td.consX)

		td.consX.SetProposal(p)
		td.shouldPublishVote(t, td.consX, vote.VoteTypePrepare, p.Block().Hash())
		td.addVote(td.consX, vote.VoteTypePrepare, h, r, p.Block().Hash(), tIndexY)
		td.addVote(td.consX, vote.VoteTypePrepare, h, r, hash.UndefHash, tIndexB)

		td.shouldPublishVote(t, td.consX, vote.VoteTypeChangeProposer, hash.UndefHash)
		td.addVote(td.consX, vote.VoteTypeChangeProposer, h, r, hash.UndefHash, tIndexY)
		td.addVote(td.consX, vote.VoteTypeChangeProposer, h, r, hash.UndefHash, tIndexB)

		// Nx and Ny move to next round

		wg.Done()
	}()
	// =================================
	// Np votes
	go func() {
		td.enterNewHeight(td.consP)

		td.consP.SetProposal(p)
		td.shouldPublishVote(t, td.consP, vote.VoteTypePrepare, p.Block().Hash())
		td.addVote(td.consP, vote.VoteTypePrepare, h, r, p.Block().Hash(), tIndexX)
		td.addVote(td.consP, vote.VoteTypePrepare, h, r, p.Block().Hash(), tIndexB)

		td.shouldPublishVote(t, td.consP, vote.VoteTypePrecommit, p.Block().Hash())
		td.addVote(td.consP, vote.VoteTypePrecommit, h, r, p.Block().Hash(), tIndexB)
		// Np is unable to progress

		wg.Done()
	}()

	// =================================
	wg.Wait()
	td.checkHeightRound(t, td.consX, h, r+1)
	td.checkHeightRound(t, td.consP, h, r)

	// =================================
	// Now, Partition heals

	time.Sleep(1 * time.Second)
	fmt.Println("=== Partition healed")

	for _, v := range td.consP.AllVotes() {
		td.consX.AddVote(v)
	}
	td.checkHeightRoundWait(t, td.consX, h, r+1)

	td.consP.SetProposal(p)
	for _, v := range td.consX.AllVotes() {
		td.consP.AddVote(v)
	}
	td.shouldPublishVote(t, td.consP, vote.VoteTypeChangeProposer, hash.UndefHash)
	td.checkHeightRoundWait(t, td.consP, h, r+1)
}

func TestQueryProposal(t *testing.T) {
	td := setup(t)

	td.commitBlockForAllStates(t)
	td.enterNewHeight(td.consX)
	td.enterNextRound(td.consX)
	td.shouldPublishQueryProposal(t, td.consX, 2, 1)
}
