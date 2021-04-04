package consensus

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/proposal"
	"github.com/zarbchain/zarb-go/vote"
)

func TestProposeBlock(t *testing.T) {
	setup(t)
	testEnterNewHeight(tConsX)
	shouldPublishProposal(t, tConsX)
}

func TestSetProposalInvalidProposer(t *testing.T) {
	setup(t)

	testEnterNewHeight(tConsY)
	assert.Nil(t, tConsY.RoundProposal(0))

	addr := tSigners[tIndexB].Address()
	b, _ := block.GenerateTestBlock(&addr, nil)
	p := proposal.NewProposal(1, 0, *b)

	tConsY.SetProposal(p)
	assert.Nil(t, tConsY.RoundProposal(0))

	tSigners[tIndexB].SignMsg(p) // Invalid signature
	tConsY.SetProposal(p)
	assert.Nil(t, tConsY.RoundProposal(0))
}

func TestSetProposalInvalidBlock(t *testing.T) {
	setup(t)

	a := tSigners[tIndexB].Address()
	invBlock, _ := block.GenerateTestBlock(&a, nil)
	p := proposal.NewProposal(1, 2, *invBlock)
	tSigners[tIndexB].SignMsg(p)

	testEnterNewHeight(tConsP)
	testEnterNewRound(tConsP)

	tConsP.SetProposal(p)
	assert.Nil(t, tConsP.RoundProposal(2))
}

func TestSetProposalInvalidHeight(t *testing.T) {
	setup(t)

	a := tSigners[tIndexB].Address()
	invBlock, _ := block.GenerateTestBlock(&a, nil)
	p := proposal.NewProposal(2, 0, *invBlock)
	tSigners[tIndexB].SignMsg(p)

	testEnterNewHeight(tConsY)
	tConsY.SetProposal(p)
	assert.Nil(t, tConsY.RoundProposal(2))
}

func TestConsensusSetProposalAfterCommit(t *testing.T) {
	setup(t)

	p := makeProposal(t, 1, 0)

	testEnterNewHeight(tConsP)
	commitBlockForAllStates(t)
	tConsP.SetProposal(p)
	assert.Nil(t, tConsP.RoundProposal(0))
}

func TestSecondProposalCommitted(t *testing.T) {
	setup(t)

	commitBlockForAllStates(t)
	commitBlockForAllStates(t)

	testEnterNewHeight(tConsX)

	// Now it's turn for Byzantine node to propose a block
	// Other nodes are going to not accept its proposal, even it is valid
	p1 := makeProposal(t, 3, 0) // valid proposal for round 0, byzantine proposer
	p2 := makeProposal(t, 3, 1) // valid proposal for round 1, partitioned proposer

	// Probably we have blocked Byzantine node
	//tConsX.SetProposal(p1)

	shouldPublishVote(t, tConsX, vote.VoteTypePrepare, crypto.UndefHash)
	testAddVote(t, tConsX, vote.VoteTypePrepare, 3, 0, crypto.UndefHash, tIndexY)
	testAddVote(t, tConsX, vote.VoteTypePrepare, 3, 0, p1.Block().Hash(), tIndexB)
	testAddVote(t, tConsX, vote.VoteTypePrepare, 3, 0, crypto.UndefHash, tIndexP)

	shouldPublishVote(t, tConsX, vote.VoteTypePrecommit, crypto.UndefHash)
	testAddVote(t, tConsX, vote.VoteTypePrecommit, 3, 0, crypto.UndefHash, tIndexY)
	testAddVote(t, tConsX, vote.VoteTypePrecommit, 3, 0, p1.Block().Hash(), tIndexB)
	testAddVote(t, tConsX, vote.VoteTypePrecommit, 3, 0, crypto.UndefHash, tIndexP)

	tConsX.SetProposal(p2)

	shouldPublishVote(t, tConsX, vote.VoteTypePrepare, p2.Block().Hash())
	testAddVote(t, tConsX, vote.VoteTypePrepare, 3, 1, p2.Block().Hash(), tIndexY)
	testAddVote(t, tConsX, vote.VoteTypePrepare, 3, 1, crypto.UndefHash, tIndexB)
	testAddVote(t, tConsX, vote.VoteTypePrepare, 3, 1, p2.Block().Hash(), tIndexP)

	shouldPublishVote(t, tConsX, vote.VoteTypePrecommit, p2.Block().Hash())
	testAddVote(t, tConsX, vote.VoteTypePrecommit, 3, 1, p2.Block().Hash(), tIndexY)
	testAddVote(t, tConsX, vote.VoteTypePrecommit, 3, 1, crypto.UndefHash, tIndexB)
	testAddVote(t, tConsX, vote.VoteTypePrecommit, 3, 1, p2.Block().Hash(), tIndexP)

	shouldPublishBlockAnnounce(t, tConsX, p2.Block().Hash())
}

func TestNetworkLagging(t *testing.T) {
	setup(t)

	testEnterNewHeight(tConsP)

	h := 1
	r := 0
	p := makeProposal(t, h, r)
	// We don't set proposal for second validator here
	// tConsP.SetProposal(p)

	shouldPublishQueryProposal(t, tConsP, h, r)
	shouldPublishVote(t, tConsP, vote.VoteTypePrepare, crypto.UndefHash)

	testAddVote(t, tConsP, vote.VoteTypePrepare, h, r, p.Block().Hash(), tIndexX)
	testAddVote(t, tConsP, vote.VoteTypePrepare, h, r, p.Block().Hash(), tIndexY)
	testAddVote(t, tConsP, vote.VoteTypePrepare, h, r, crypto.UndefHash, tIndexB)

	// Now let's set the proposal
	tConsP.SetProposal(p)
	shouldPublishVote(t, tConsP, vote.VoteTypePrecommit, p.Block().Hash())
}

func TestLateUndefVote(t *testing.T) {
	setup(t)

	commitBlockForAllStates(t)
	commitBlockForAllStates(t)

	testEnterNewHeight(tConsP)

	h := 3
	r := 0
	p := makeProposal(t, h, r) // Other nodes doesn't accept byzantine proposal

	// tConsP is partitioned
	testAddVote(t, tConsP, vote.VoteTypePrepare, h, r, p.Block().Hash(), tIndexB)
	testAddVote(t, tConsP, vote.VoteTypePrecommit, h, r, p.Block().Hash(), tIndexB)

	//  tConsP doesn't accept tConsB's proposal
	shouldPublishVote(t, tConsP, vote.VoteTypePrepare, crypto.UndefHash)

	// Now partition healed
	testAddVote(t, tConsP, vote.VoteTypePrepare, h, r, crypto.UndefHash, tIndexX)

	// Enough prepare votes, tConsP can vote for undef precommit
	shouldPublishVote(t, tConsP, vote.VoteTypePrecommit, crypto.UndefHash)
}

func TestProposalInvalidRound(t *testing.T) {
	setup(t)

	commitBlockForAllStates(t)

	testEnterNewHeight(tConsX)

	// Byzantine node sends proposal for second round (his turn)
	b, err := tConsB.state.ProposeBlock(1)
	assert.NoError(t, err)
	p := proposal.NewProposal(2, 1, *b)
	tSigners[tIndexB].SignMsg(p)

	tConsX.SetProposal(p)

	// tConsX doesn't accept the proposal for next rounds
	assert.Nil(t, tConsX.RoundProposal(1))
}
