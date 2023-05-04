package consensus

import (
	"fmt"
	"sync"
	"time"

	"github.com/pactus-project/pactus/consensus/log"
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/state"
	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/pactus-project/pactus/types/block"
	"github.com/pactus-project/pactus/types/proposal"
	"github.com/pactus-project/pactus/types/validator"
	"github.com/pactus-project/pactus/types/vote"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/logger"
)

type consensus struct {
	lk sync.RWMutex

	config              *Config
	log                 *log.Log
	signer              crypto.Signer
	rewardAddr          crypto.Address
	state               state.Facade
	height              uint32
	round               int16
	active              bool
	newHeightState      consState // TODO: rename consState to consPhase to prevent confusion with the blockchain state
	proposeState        consState
	prepareState        consState
	precommitState      consState
	commitState         consState
	currentState        consState
	changeProposerState consState
	broadcastCh         chan message.Message
	mediator            mediator
	logger              *logger.Logger
}

func NewConsensus(
	conf *Config,
	state state.Facade,
	signer crypto.Signer,
	rewardAddr crypto.Address,
	broadcastCh chan message.Message,
	mediator mediator) Consensus {
	cs := &consensus{
		config:      conf,
		state:       state,
		broadcastCh: broadcastCh,
		signer:      signer,
	}

	// Update height later, See enterNewHeight.
	cs.log = log.NewLog()
	cs.logger = logger.NewLogger("_consensus", cs)
	cs.rewardAddr = rewardAddr

	cs.newHeightState = &newHeightState{cs}
	cs.proposeState = &proposeState{cs}
	cs.prepareState = &prepareState{cs, false}
	cs.precommitState = &precommitState{cs, false}
	cs.commitState = &commitState{cs}
	cs.changeProposerState = &changeProposerState{cs}

	cs.height = 0
	cs.round = 0
	cs.active = false
	cs.mediator = mediator

	mediator.Register(cs)

	return cs
}

func (cs *consensus) Fingerprint() string {
	return fmt.Sprintf("{%s %d/%d/%s}",
		cs.signer.Address().Fingerprint(),
		cs.height, cs.round, cs.currentState.name())
}

func (cs *consensus) SignerKey() crypto.PublicKey {
	cs.lk.RLock()
	defer cs.lk.RUnlock()

	return cs.signer.PublicKey()
}

func (cs *consensus) HeightRound() (uint32, int16) {
	cs.lk.RLock()
	defer cs.lk.RUnlock()

	return cs.height, cs.round
}

func (cs *consensus) RoundProposal(round int16) *proposal.Proposal {
	cs.lk.RLock()
	defer cs.lk.RUnlock()

	return cs.log.RoundProposal(round)
}

func (cs *consensus) AllVotes() []*vote.Vote {
	cs.lk.RLock()
	defer cs.lk.RUnlock()

	votes := []*vote.Vote{}
	for r := int16(0); r <= cs.round; r++ {
		m := cs.log.MustGetRoundMessages(r)
		votes = append(votes, m.AllVotes()...)
	}
	return votes
}

func (cs *consensus) enterNewState(s consState) {
	cs.currentState = s
	cs.currentState.enter()
}

func (cs *consensus) MoveToNewHeight() {
	cs.lk.Lock()
	defer cs.lk.Unlock()

	// Move the consensus to a new height only if it is behind the state height.
	if cs.currentState != cs.newHeightState {
		cs.enterNewState(cs.newHeightState)
	}
}

func (cs *consensus) scheduleTimeout(duration time.Duration, height uint32, round int16, target tickerTarget) {
	ti := &ticker{duration, height, round, target}
	timer := time.NewTimer(duration)
	cs.logger.Debug("new timer scheduled ⏱️", "duration", duration, "height", height, "round", round, "target", target)

	go func() {
		<-timer.C
		cs.handleTimeout(ti)
	}()
}

func (cs *consensus) SetProposal(p *proposal.Proposal) {
	cs.lk.Lock()
	defer cs.lk.Unlock()

	if !cs.active {
		cs.logger.Trace("we are not in the committee")
		return
	}

	if p.Height() != cs.height {
		cs.logger.Trace("invalid height", "proposal", p)
		return
	}

	if p.Height() == cs.state.LastBlockHeight() {
		cs.logger.Trace("block is committed", "proposal", p, "block hash", cs.state.LastBlockHash())
		return
	}

	if p.Round() < cs.round {
		cs.logger.Trace("expired round", "proposal", p)
		return
	}

	roundProposal := cs.log.RoundProposal(p.Round())
	if roundProposal != nil {
		cs.logger.Trace("this round has proposal", "proposal", p)
		return
	}

	proposer := cs.proposer(p.Round())
	if err := p.Verify(proposer.PublicKey()); err != nil {
		cs.logger.Warn("proposal has invalid signature", "proposal", p, "err", err)
		return
	}

	if err := cs.state.ValidateBlock(p.Block()); err != nil {
		cs.logger.Warn("invalid block", "proposal", p, "err", err)
		return
	}

	cs.currentState.onSetProposal(p)
}

func (cs *consensus) doSetProposal(p *proposal.Proposal) {
	cs.logger.Info("proposal set", "proposal", p)
	cs.log.SetRoundProposal(p.Round(), p)
}

func (cs *consensus) handleTimeout(t *ticker) {
	cs.lk.Lock()
	defer cs.lk.Unlock()

	cs.logger.Trace("handle ticker", "ticker", t)

	// Old tickers might be triggered now. Ignore them.
	if cs.height != t.Height || cs.round != t.Round {
		cs.logger.Trace("stale ticker", "ticker", t)
		return
	}

	cs.currentState.onTimeout(t)
}

func (cs *consensus) AddVote(v *vote.Vote) {
	cs.lk.Lock()
	defer cs.lk.Unlock()

	if !cs.active {
		cs.logger.Trace("we are not in the committee")
		return
	}

	if v.Height() != cs.height {
		cs.logger.Trace("vote has invalid height", "vote", v)
		return
	}

	if cs.log.HasVote(v.Hash()) {
		cs.logger.Trace("vote exists", "vote", v)
		return
	}

	cs.currentState.onAddVote(v)
}

func (cs *consensus) doAddVote(v *vote.Vote) {
	err := cs.log.AddVote(v)
	if err != nil {
		cs.logger.Error("error on adding a vote", "vote", v, "err", err)
	}

	cs.logger.Debug("new vote added", "vote", v)
}

func (cs *consensus) proposer(round int16) *validator.Validator {
	return cs.state.Proposer(round)
}

func (cs *consensus) signAddVote(msgType vote.Type, hash hash.Hash) {
	address := cs.signer.Address()
	if !cs.log.CanVote(address) {
		cs.logger.Trace("this node is not in committee", "addr", address)
		return
	}

	// Sign the vote
	v := vote.NewVote(msgType, cs.height, cs.round, hash, address)
	cs.signer.SignMsg(v)
	cs.logger.Info("our vote signed and broadcasted", "vote", v)

	err := cs.log.AddVote(v)
	if err != nil {
		cs.logger.Error("error on adding our vote", "err", err, "vote", v)
	} else {
		cs.broadcastVote(v)
	}
}

func (cs *consensus) queryProposal() {
	cs.broadcast(message.NewQueryProposalMessage(cs.height, cs.round))
}

func (cs *consensus) broadcastProposal(p *proposal.Proposal) {
	go cs.mediator.OnPublishProposal(cs, p)
	cs.broadcast(message.NewProposalMessage(p))
}

func (cs *consensus) broadcastVote(v *vote.Vote) {
	go cs.mediator.OnPublishVote(cs, v)
	cs.broadcast(message.NewVoteMessage(v))
}

func (cs *consensus) announceNewBlock(h uint32, b *block.Block, c *block.Certificate) {
	go cs.mediator.OnBlockAnnounce(cs)
	cs.broadcast(message.NewBlockAnnounceMessage(h, b, c))
}

func (cs *consensus) broadcast(msg message.Message) {
	if !cs.active {
		cs.logger.Warn("non-active validators should not publish messages")
		panic("should not happens")
	}

	cs.broadcastCh <- msg
}

// IsActive checks if the consensus is in an active state and participating in the consensus algorithm.
func (cs *consensus) IsActive() bool {
	cs.lk.RLock()
	defer cs.lk.RUnlock()

	return cs.active
}

// TODO: Improve the performance?
func (cs *consensus) PickRandomVote() *vote.Vote {
	cs.lk.RLock()
	defer cs.lk.RUnlock()

	rndRound := util.RandInt16(cs.round + 1)
	votes := []*vote.Vote{}
	if rndRound == cs.round {
		m := cs.log.MustGetRoundMessages(rndRound)
		votes = append(votes, m.AllVotes()...)
	} else {
		// Don't broadcast prepare and precommit votes for previous rounds
		vs := cs.log.ChangeProposerVoteSet(rndRound)
		votes = append(votes, vs.AllVotes()...)
	}
	if len(votes) == 0 {
		return nil
	}
	return votes[util.RandInt32(int32(len(votes)))]
}
