package firewall

import (
	"testing"

	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/network"
	"github.com/pactus-project/pactus/sync/bundle"
	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/pactus-project/pactus/types"
	"github.com/pactus-project/pactus/types/block"
	"github.com/pactus-project/pactus/types/proposal"
	"github.com/pactus-project/pactus/types/validator"
	"github.com/pactus-project/pactus/types/vote"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

// consensusGossip encodes a consensus message as a gossip message of the mainnet.
func consensusGossip(t *testing.T, td *testData, msg message.Message) *network.GossipMessage {
	t.Helper()

	bdl := bundle.NewBundle(msg)
	bdl.Flags = util.SetFlag(bdl.Flags, bundle.BundleFlagNetworkMainnet)
	data, err := bdl.Encode()
	require.NoError(t, err)

	return &network.GossipMessage{From: td.goodPeerID, Data: data}
}

func setupCommitteeCheck(t *testing.T) (*testData, []*bls.ValidatorKey, types.Height) {
	t.Helper()

	td := setup(t, nil)
	lastHeight := td.RandHeight()
	keys := make([]*bls.ValidatorKey, 0, 4)
	committee := make([]*validator.Validator, 0, 4)
	for i := 0; i < 4; i++ {
		key := td.RandValKey()
		keys = append(keys, key)
		committee = append(committee, td.GenerateTestValidator(
			testsuite.ValidatorWithPublicKey(key.PublicKey())))
	}

	td.state.EXPECT().LastBlockHeight().Return(lastHeight).AnyTimes()
	td.state.EXPECT().CommitteeValidators().Return(committee).AnyTimes()
	// The proposer of every round is the first member, as a simple rule.
	td.state.EXPECT().Proposer(gomock.Any()).Return(committee[0]).AnyTimes()

	return td, keys, lastHeight + 1
}

func TestConsensusRequestVotes(t *testing.T) {
	td, keys, height := setupCommitteeCheck(t)

	t.Run("a vote signed by its committee signer is relayed", func(t *testing.T) {
		vte := vote.NewPrecommitVote(td.RandHash(), height, 0, keys[1].Address())
		td.HelperSignVote(keys[1], vte)

		assert.Equal(t, network.Propagate,
			td.firewall.AllowConsensusRequest(consensusGossip(t, td, message.NewVoteMessage(vte))))
	})

	t.Run("a vote with a forged signature is dropped", func(t *testing.T) {
		vte := vote.NewPrecommitVote(td.RandHash(), height, 0, keys[1].Address())
		td.HelperSignVote(keys[2], vte)

		assert.Equal(t, network.Drop,
			td.firewall.AllowConsensusRequest(consensusGossip(t, td, message.NewVoteMessage(vte))))
	})

	t.Run("a vote of a signer outside the committee is dropped", func(t *testing.T) {
		outsider := td.RandValKey()
		vte := vote.NewPrecommitVote(td.RandHash(), height, 0, outsider.Address())
		td.HelperSignVote(outsider, vte)

		assert.Equal(t, network.Drop,
			td.firewall.AllowConsensusRequest(consensusGossip(t, td, message.NewVoteMessage(vte))))
	})

	t.Run("a vote of the next height is left to the consensus", func(t *testing.T) {
		// The committee of the next height may differ, so it is not checked here.
		outsider := td.RandValKey()
		vte := vote.NewPrecommitVote(td.RandHash(), height+1, 0, outsider.Address())
		td.HelperSignVote(td.RandValKey(), vte)

		assert.Equal(t, network.Propagate,
			td.firewall.AllowConsensusRequest(consensusGossip(t, td, message.NewVoteMessage(vte))))
	})

	t.Run("a query is not checked", func(t *testing.T) {
		query := message.NewQueryVoteMessage(height, 0, keys[0].Address())

		assert.Equal(t, network.Propagate,
			td.firewall.AllowConsensusRequest(consensusGossip(t, td, query)))
	})
}

func TestConsensusRequestProposals(t *testing.T) {
	td, keys, height := setupCommitteeCheck(t)

	t.Run("a proposal signed by the round's proposer is relayed", func(t *testing.T) {
		prop := td.GenerateTestProposal(height, 0, testsuite.ProposalWithKey(keys[0]))

		assert.Equal(t, network.Propagate,
			td.firewall.AllowConsensusRequest(consensusGossip(t, td, message.NewProposalMessage(prop))))
	})

	t.Run("a proposal with a forged signature is dropped", func(t *testing.T) {
		prop := td.GenerateTestProposal(height, 0, testsuite.ProposalWithKey(keys[0]))
		prop.SetSignature(keys[1].Sign(prop.SignBytes()))

		assert.Equal(t, network.Drop,
			td.firewall.AllowConsensusRequest(consensusGossip(t, td, message.NewProposalMessage(prop))))
	})
}

func TestCommitteeKeysRefreshPerHeight(t *testing.T) {
	td := setup(t, nil)
	first := td.GenerateTestValidator()
	second := td.GenerateTestValidator()

	gomockCall := td.state.EXPECT().CommitteeValidators().Return([]*validator.Validator{first}).Times(1)
	assert.NotNil(t, td.firewall.committeePublicKey(10, first.Address()))
	assert.NotNil(t, td.firewall.committeePublicKey(10, first.Address()), "read once per height")
	assert.Nil(t, td.firewall.committeePublicKey(10, second.Address()))

	td.state.EXPECT().CommitteeValidators().Return([]*validator.Validator{second}).Times(1).After(gomockCall)
	assert.NotNil(t, td.firewall.committeePublicKey(11, second.Address()))
	assert.Nil(t, td.firewall.committeePublicKey(11, first.Address()))
}

func TestConsensusRequestSignatureBudget(t *testing.T) {
	td, keys, height := setupCommitteeCheck(t)

	// A peer flooding forged votes gets them checked up to its budget, then
	// they are no longer relayed, and no longer checked.
	results := map[network.PropagationPolicy]int{}
	for i := 0; i < maxSignatureChecksPerPeer+10; i++ {
		vte := vote.NewPrecommitVote(td.RandHash(), height, 0, keys[1].Address())
		td.HelperSignVote(keys[2], vte)
		results[td.firewall.AllowConsensusRequest(consensusGossip(t, td, message.NewVoteMessage(vte)))]++
	}

	assert.Equal(t, maxSignatureChecksPerPeer, results[network.Drop])
	assert.Equal(t, 10, results[network.DropButConsume])
	assert.Zero(t, results[network.Propagate])
}

func TestConsensusRequestBudgetCountsForgedOnly(t *testing.T) {
	td, keys, height := setupCommitteeCheck(t)

	// The peer is the relay of the messages, not their author.
	// Messages that can't be checked, and valid messages, must not use up its budget:
	// its authentic votes would no longer be relayed.
	for range 2 * maxSignatureChecksPerPeer {
		query := message.NewQueryVoteMessage(height, 0, keys[0].Address())
		require.Equal(t, network.Propagate,
			td.firewall.AllowConsensusRequest(consensusGossip(t, td, query)))

		next := vote.NewPrecommitVote(td.RandHash(), height+1, 0, keys[1].Address())
		td.HelperSignVote(keys[2], next)
		require.Equal(t, network.Propagate,
			td.firewall.AllowConsensusRequest(consensusGossip(t, td, message.NewVoteMessage(next))))
	}

	for range maxSignatureChecksPerPeer + 10 {
		vte := vote.NewPrecommitVote(td.RandHash(), height, 0, keys[1].Address())
		td.HelperSignVote(keys[1], vte)
		require.Equal(t, network.Propagate,
			td.firewall.AllowConsensusRequest(consensusGossip(t, td, message.NewVoteMessage(vte))))
	}
}

func TestConsensusRequestForgedProposalSkipsBlockCheck(t *testing.T) {
	td, keys, height := setupCommitteeCheck(t)

	// The block holds a transaction with an invalid signature.
	// A forged proposal is dropped for its own signature, and a proposal of the proposer
	// is relayed: the block is checked by the consensus, not by the gossip validator.
	trx := td.GenerateTestTransferTx()
	trx.SetSignature(td.RandBLSSignature())
	blk, _ := td.GenerateTestBlock(height,
		testsuite.BlockWithTransactions(block.Txs{trx}),
		testsuite.BlockWithProposer(keys[0].Address()))
	require.Error(t, blk.BasicCheck())

	forged := proposal.NewProposal(height, 0, blk)
	forged.SetSignature(keys[1].Sign(forged.SignBytes()))
	assert.Equal(t, network.Drop,
		td.firewall.AllowConsensusRequest(consensusGossip(t, td, message.NewProposalMessage(forged))))

	signed := proposal.NewProposal(height, 0, blk)
	signed.SetSignature(keys[0].Sign(signed.SignBytes()))
	assert.Equal(t, network.Propagate,
		td.firewall.AllowConsensusRequest(consensusGossip(t, td, message.NewProposalMessage(signed))))
}
