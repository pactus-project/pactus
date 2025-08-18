package sync

import (
	"testing"

	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/pactus-project/pactus/types/protocol"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
)

func TestHandlerProposalParsingMessages(t *testing.T) {
	td := setup(t, nil)

	t.Run("Set proposal for consensus", func(t *testing.T) {
		consensusHeight := td.state.LastBlockHeight() + 1
		prop := td.GenerateTestProposal(consensusHeight, 0)
		msg := message.NewProposalMessage(prop)
		pid := td.RandPeerID()

		td.receivingNewMessage(td.sync, msg, pid)
		assert.Equal(t, prop, td.sync.getConsMgr().Proposal())
	})

	t.Run("Update protocol version of the proposer", func(t *testing.T) {
		consensusHeight := td.state.LastBlockHeight() + 1
		valKey := td.RandValKey()
		val := td.GenerateTestValidator(testsuite.ValidatorWithPublicKey(valKey.PublicKey()))
		td.state.TestStore.UpdateValidator(val)
		prop := td.GenerateTestProposal(consensusHeight, 0, testsuite.ProposalWithKey(valKey))
		msg := message.NewProposalMessage(prop)
		pid := td.RandPeerID()

		td.receivingNewMessage(td.sync, msg, pid)

		assert.Equal(t, protocol.ProtocolVersionLatest,
			td.state.ValidatorByAddress(valKey.Address()).ProtocolVersion())
	})
}
