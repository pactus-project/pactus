package sync

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/consensus/vote"
	"github.com/zarbchain/zarb-go/sync/message/payload"
	"github.com/zarbchain/zarb-go/util"
)

func TestParsingVoteMessages(t *testing.T) {
	setup(t)

	t.Run("Parsing vote message", func(t *testing.T) {
		v, _ := vote.GenerateTestPrecommitVote(1, 0)
		pld := payload.NewVotePayload(v)

		assert.NoError(t, testReceiveingNewMessage(tSync, pld, util.RandomPeerID()))
		assert.Equal(t, tConsensus.Votes[0].Hash(), v.Hash())
	})
}
