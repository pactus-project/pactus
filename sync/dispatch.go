package sync

import (
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/message"
	"github.com/zarbchain/zarb-go/tx"
	"github.com/zarbchain/zarb-go/vote"
)

func (syncer *Synchronizer) broadcastSalam() {
	height := syncer.state.LastBlockHeight()
	msg := message.NewSalamMessage(height)
	syncer.networkApi.PublishMessage(msg)
}

func (syncer *Synchronizer) broadcastBlocksReq(from, to int) {
	msg := message.NewBlocksReqMessage(from, to, syncer.state.LastBlockHash())
	syncer.networkApi.PublishMessage(msg)
}

func (syncer *Synchronizer) broadcastBlocks(from int, blocks []block.Block, lastCommit *block.Commit) {
	msg := message.NewBlocksMessage(from, blocks, lastCommit)
	syncer.networkApi.PublishMessage(msg)
}

func (syncer *Synchronizer) broadcastTxs(txs []tx.Tx) {
	msg := message.NewTxsMessage(txs)
	syncer.networkApi.PublishMessage(msg)
}

func (syncer *Synchronizer) broadcastTxsReq(hashes []crypto.Hash) {
	msg := message.NewTxsReqMessage(hashes)
	syncer.networkApi.PublishMessage(msg)
}

func (syncer *Synchronizer) broadcastHeartBeat() {
	hasProposal := syncer.consensus.HasProposal()
	hrs := syncer.consensus.HRS()

	msg := message.NewHeartBeatMessage(syncer.state.LastBlockHash(), hrs, hasProposal)
	syncer.networkApi.PublishMessage(msg)
}

func (syncer *Synchronizer) broadcastProposal(proposal vote.Proposal) {
	msg := message.NewProposalMessage(proposal)
	syncer.networkApi.PublishMessage(msg)
}

func (syncer *Synchronizer) broadcastVote(v *vote.Vote) {
	msg := message.NewVoteMessage(v)
	syncer.networkApi.PublishMessage(msg)
}

func (syncer *Synchronizer) broadcastVoteSet(height int, hashes []crypto.Hash) {
	msg := message.NewVoteSetMessage(height, hashes)
	syncer.networkApi.PublishMessage(msg)
}
