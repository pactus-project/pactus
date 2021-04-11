package sync

import (
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/zarbchain/zarb-go/consensus"
	"github.com/zarbchain/zarb-go/logger"
	"github.com/zarbchain/zarb-go/proposal"
	"github.com/zarbchain/zarb-go/sync/message"
	"github.com/zarbchain/zarb-go/sync/message/payload"
	"github.com/zarbchain/zarb-go/tx"
	"github.com/zarbchain/zarb-go/vote"
)

type ConsensusSync struct {
	config    *Config
	selfID    peer.ID
	consensus consensus.Consensus
	logger    *logger.Logger
	publishFn publishMessageFn
}

func NewConsensusSync(
	conf *Config,
	selfID peer.ID,
	consensus consensus.Consensus,
	logger *logger.Logger,
	publishFn publishMessageFn) *ConsensusSync {
	return &ConsensusSync{
		config:    conf,
		selfID:    selfID,
		consensus: consensus,
		logger:    logger,
		publishFn: publishFn,
	}
}

func (cs *ConsensusSync) BroadcastQueryProposal(height, round int) {
	msg := message.NewQueryProposalMessage(cs.selfID, height, round)
	cs.publishFn(msg)
}

func (cs *ConsensusSync) BroadcastProposal(p *proposal.Proposal) {
	msg := message.NewProposalMessage(p)
	cs.publishFn(msg)
}

func (cs *ConsensusSync) BroadcastVote(v *vote.Vote) {
	msg := message.NewVoteMessage(v)
	cs.publishFn(msg)
}

func (cs *ConsensusSync) BroadcastQueryVotes(height, round int) {
	msg := message.NewQueryVoteMessage(cs.selfID, height, round)
	cs.publishFn(msg)
}

func (cs *ConsensusSync) BroadcastQueryTransaction(ids []tx.ID) {
	msg := message.NewQueryTransactionsMessage(cs.selfID, ids)
	cs.publishFn(msg)
}

func (cs *ConsensusSync) ProcessVotePayload(pld *payload.VotePayload) {
	cs.logger.Trace("Process vote payload", "pld", pld)

	cs.consensus.AddVote(pld.Vote)
}

func (cs *ConsensusSync) ProcessQueryVotesPayload(pld *payload.QueryVotesPayload) {
	cs.logger.Trace("Process vote-set payload", "pld", pld)

	hrs := cs.consensus.HRS()
	if pld.Height == hrs.Height() {
		v := cs.consensus.PickRandomVote(pld.Round)
		if v != nil {
			cs.BroadcastVote(v)
		}
	}
}
func (cs *ConsensusSync) ProcessQueryProposalPayload(pld *payload.QueryProposalPayload) {
	cs.logger.Trace("Process proposal request payload", "pld", pld)

	hrs := cs.consensus.HRS()
	if pld.Height == hrs.Height() {
		p := cs.consensus.RoundProposal(pld.Round)
		if p != nil {
			cs.BroadcastProposal(p)
		}
	}
}

func (cs *ConsensusSync) ProcessProposalPayload(pld *payload.ProposalPayload) {
	cs.logger.Trace("Process proposal payload", "pld", pld)

	cs.consensus.SetProposal(pld.Proposal)
}
