package sync

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/consensus/vote"
	"github.com/zarbchain/zarb-go/sync/bundle/message"
	"github.com/zarbchain/zarb-go/util"
)

func TestParsingVoteMessages(t *testing.T) {
	setup(t)

	t.Run("Parsing vote message", func(t *testing.T) {
		v, _ := vote.GenerateTestPrecommitVote(1, 0)
		msg := message.NewVoteMessage(v)

		assert.NoError(t, testReceiveingNewMessage(tSync, msg, util.RandomPeerID()))
		assert.Equal(t, tConsensus.Votes[0].Hash(), v.Hash())
	})
}
