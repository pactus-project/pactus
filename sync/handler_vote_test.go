package sync

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/sync/message/payload"
	"github.com/zarbchain/zarb-go/util"
	"github.com/zarbchain/zarb-go/vote"
)

func TestParsingVoteMessages(t *testing.T) {
	setup(t)
	disableHeartbeat(t)

	t.Run("Alice receives a vote. she sends it to consensus", func(t *testing.T) {
		consensusHeight := tAliceState.LastBlockHeight() + 1
		v1, _ := vote.GenerateTestPrecommitVote(consensusHeight, 0)
		pld := payload.NewVotePayload(v1)

		tAliceNet.ReceivingMessageFromOtherPeer(util.RandomPeerID(), pld)
		assert.Equal(t, tAliceConsensus.Votes[0].Hash(), v1.Hash())
	})
}
