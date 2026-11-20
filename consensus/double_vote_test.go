package consensus

import (
	"testing"

	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/pactus-project/pactus/types/vote"
	"github.com/stretchr/testify/assert"
)

// After a double proposal an honest validator must not sign a second, different prepare vote.
func TestNoSecondPrepareVoteInTheSameRound(t *testing.T) {
	td := setup(t)

	td.enterNewHeight(td.consP)

	firstHash, secondHash := td.RandHash(), td.RandHash()
	td.consP.lk.Lock()
	td.consP.signAddPrepareVote(firstHash)
	td.consP.signAddPrepareVote(secondHash)
	td.consP.signAddPrepareVote(firstHash) // The same vote can be sent again.
	td.consP.lk.Unlock()

	votes := []hash.Hash{}
	for _, consMsg := range td.consMessages {
		if consMsg.sender != td.consP.valKey.Address() || consMsg.message.Type() != message.TypeVote {
			continue
		}
		v := consMsg.message.(*message.VoteMessage).Vote
		if v.Type() == vote.VoteTypePrepare {
			votes = append(votes, v.BlockHash())
		}
	}

	assert.Equal(t, []hash.Hash{firstHash, firstHash}, votes)
}

// After a double proposal, when the change-proposer phase decides to keep the proposer,
// the validator prepares again in the same round, for the block that the phase decided.
// Without its vote, the other honest validators may not reach a quorum.
// Any other second prepare vote is still refused.
func TestSecondPrepareVoteForTheDecidedBlock(t *testing.T) {
	td := setup(t)

	td.enterNewHeight(td.consP)

	droppedHash, decidedHash, otherHash := td.RandHash(), td.RandHash(), td.RandHash()
	td.consP.lk.Lock()
	td.consP.signAddPrepareVote(droppedHash)

	td.consP.cpDecided = 0
	td.consP.cpWeakValidity = &decidedHash
	td.consP.signAddPrepareVote(otherHash)
	td.consP.signAddPrepareVote(decidedHash)
	td.consP.lk.Unlock()

	votes := []hash.Hash{}
	for _, consMsg := range td.consMessages {
		if consMsg.sender != td.consP.valKey.Address() || consMsg.message.Type() != message.TypeVote {
			continue
		}
		v := consMsg.message.(*message.VoteMessage).Vote
		if v.Type() == vote.VoteTypePrepare {
			votes = append(votes, v.BlockHash())
		}
	}

	assert.Equal(t, []hash.Hash{droppedHash, decidedHash}, votes)
}

// The second vote that the change-proposer phase allows is a prepare vote only, after a decision
// to keep the proposer, in the current round.
func TestSecondVoteExceptionIsNarrow(t *testing.T) {
	signedVotes := func(td *testData, voteType vote.Type) []hash.Hash {
		hashes := []hash.Hash{}
		for _, consMsg := range td.consMessages {
			if consMsg.sender != td.consP.valKey.Address() || consMsg.message.Type() != message.TypeVote {
				continue
			}
			v := consMsg.message.(*message.VoteMessage).Vote
			if v.Type() == voteType {
				hashes = append(hashes, v.BlockHash())
			}
		}

		return hashes
	}

	t.Run("not for a second precommit vote", func(t *testing.T) {
		td := setup(t)
		td.enterNewHeight(td.consP)

		firstHash, decidedHash := td.RandHash(), td.RandHash()
		td.consP.lk.Lock()
		td.consP.signAddPrecommitVote(firstHash)
		td.consP.cpDecided = 0
		td.consP.cpWeakValidity = &decidedHash
		td.consP.signAddPrecommitVote(decidedHash)
		td.consP.lk.Unlock()

		assert.Equal(t, []hash.Hash{firstHash}, signedVotes(td, vote.VoteTypePrecommit))
	})

	t.Run("not after a decision to change the proposer", func(t *testing.T) {
		td := setup(t)
		td.enterNewHeight(td.consP)

		firstHash, decidedHash := td.RandHash(), td.RandHash()
		td.consP.lk.Lock()
		td.consP.signAddPrepareVote(firstHash)
		td.consP.cpDecided = 1
		td.consP.cpWeakValidity = &decidedHash
		td.consP.signAddPrepareVote(decidedHash)
		td.consP.lk.Unlock()

		assert.Equal(t, []hash.Hash{firstHash}, signedVotes(td, vote.VoteTypePrepare))
	})

	t.Run("not for another round", func(t *testing.T) {
		td := setup(t)
		td.enterNewHeight(td.consP)

		firstHash, decidedHash := td.RandHash(), td.RandHash()
		td.consP.lk.Lock()
		td.consP.signAddPrepareVote(firstHash)
		oldRound := td.consP.round
		td.consP.round++
		td.consP.cpDecided = 0
		td.consP.cpWeakValidity = &decidedHash
		td.consP.signAddVote(vote.NewPrepareVote(decidedHash, td.consP.height, oldRound, td.consP.valKey.Address()))
		td.consP.lk.Unlock()

		assert.Equal(t, []hash.Hash{firstHash}, signedVotes(td, vote.VoteTypePrepare))
	})
}
