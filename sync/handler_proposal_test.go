package sync

import (
	"testing"

	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/stretchr/testify/assert"
)

func TestParsingProposalMessages(t *testing.T) {
	td := setup(t, nil)

	t.Run("Parsing proposal message", func(t *testing.T) {
		consensusHeight := td.state.LastBlockHeight() + 1
		prop, _ := td.GenerateTestProposal(consensusHeight, 0)
		msg := message.NewProposalMessage(prop)

		assert.NoError(t, td.receivingNewMessage(td.sync, msg, td.RandomPeerID()))
		assert.NotNil(t, td.consMgr.RoundProposal(0))
	})
}
