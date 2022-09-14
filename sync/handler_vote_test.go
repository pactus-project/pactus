package sync

import (
	"testing"

	"github.com/pactus-project/pactus/network"
	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/pactus-project/pactus/types/vote"
	"github.com/stretchr/testify/assert"
)

func TestParsingVoteMessages(t *testing.T) {
	setup(t)

	t.Run("Parsing vote message", func(t *testing.T) {
		v, _ := vote.GenerateTestPrecommitVote(1, 0)
		msg := message.NewVoteMessage(v)

		assert.NoError(t, testReceivingNewMessage(tSync, msg, network.TestRandomPeerID()))
		assert.Equal(t, tConsensus.Votes[0].Hash(), v.Hash())
	})
}
