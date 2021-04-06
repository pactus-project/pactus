package consensus

import (
	"fmt"
	"time"

	"github.com/sasha-s/go-deadlock"
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/consensus/pending_votes"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/logger"
	"github.com/zarbchain/zarb-go/proposal"
	"github.com/zarbchain/zarb-go/state"
	"github.com/zarbchain/zarb-go/sync/message/payload"
	"github.com/zarbchain/zarb-go/util"
	"github.com/zarbchain/zarb-go/validator"
	"github.com/zarbchain/zarb-go/vote"
)

type consensus struct {
	lk deadlock.RWMutex

	config              *Config
	pendingVotes        *pending_votes.PendingVotes
	signer              crypto.Signer
	state               state.StateFacade
	height              int
	round               int
	newHeightState      consState
	newRoundState       consState
	proposeState        consState
	prepareState        consState
	precommitState      consState
	commitState         consState
	currentState        consState
	changeProposerState consState
	broadcastCh         chan payload.Payload
	logger              *logger.Logger
}

func NewConsensus(
	conf *Config,
	state state.StateFacade,
	signer crypto.Signer,
	broadcastCh chan payload.Payload) (Consensus, error) {
	cs := &consensus{
		config:      conf,
		state:       state,
		broadcastCh: broadcastCh,
		signer:      signer,
	}

	// Update height later, See enterNewHeight.
	cs.pendingVotes = pending_votes.NewPendingVotes()
	cs.logger = logger.NewLogger("_consensus", cs)

	cs.newHeightState = &newHeightState{cs}
	cs.newRoundState = &newRoundState{cs}
	cs.proposeState = &proposeState{cs}
	cs.prepareState = &prepareState{cs, false}
	cs.precommitState = &precommitState{cs, false}
	cs.commitState = &commitState{cs}
	cs.changeProposerState = &changeProposerState{cs}

	cs.height = -1
	cs.round = -1
	cs.MoveToNewHeight()

	return cs, nil
}

func (cs *consensus) Stop() {
}

func (cs *consensus) Fingerprint() string {
	return fmt.Sprintf("{%d/%d/%s}", cs.height, cs.round, cs.currentState.name())
}

func (cs *consensus) HeightRound() (int, int) {
	cs.lk.RLock()
	defer cs.lk.RUnlock()

	return cs.height, cs.round
}

func (cs *consensus) Height() int {
	cs.lk.RLock()
	defer cs.lk.RUnlock()

	return cs.height
}

func (cs *consensus) Round() int {
	cs.lk.RLock()
	defer cs.lk.RUnlock()

	return cs.round
}

func (cs *consensus) RoundProposal(round int) *proposal.Proposal {
	cs.lk.RLock()
	defer cs.lk.RUnlock()

	return cs.pendingVotes.RoundProposal(round)
}

func (cs *consensus) RoundVotes(round int) []*vote.Vote {
	cs.lk.Lock()
	defer cs.lk.Unlock()

	rv := cs.pendingVotes.MustGetRoundVotes(round)
	return rv.AllVotes()
}

func (cs *consensus) HasVote(hash crypto.Hash) bool {
	cs.lk.Lock()
	defer cs.lk.Unlock()

	return cs.pendingVotes.HasVote(hash)
}

func (cs *consensus) enterNewState(s consState) {
	cs.currentState = s
	cs.currentState.enter()
}

func (cs *consensus) MoveToNewHeight() {
	cs.lk.Lock()
	defer cs.lk.Unlock()

	if cs.state.LastBlockHeight()+1 > cs.height {
		cs.enterNewState(cs.newHeightState)
	}
}

func (cs *consensus) scheduleTimeout(duration time.Duration, height int, round int, target tickerTarget) {
	ti := &ticker{duration, height, round, target}
	timer := time.NewTimer(duration)
	go func() {
		<-timer.C
		cs.handleTimeout(ti)
	}()
	logger.Trace("Scheduled timeout", "duration", duration, "height", height, "round", round, "target", target)
}

func (cs *consensus) SetProposal(p *proposal.Proposal) {
	cs.lk.Lock()
	defer cs.lk.Unlock()

	if p.Height() != cs.height {
		cs.logger.Trace("Invalid height", "proposal", p)
		return
	}

	roundProposal := cs.pendingVotes.RoundProposal(p.Round())
	if roundProposal != nil {
		cs.logger.Trace("This round has proposal", "proposal", p)
		return
	}

	proposer := cs.proposer(p.Round())
	if err := p.Verify(proposer.PublicKey()); err != nil {
		cs.logger.Error("Proposal has invalid signature", "proposal", p, "err", err)
		return
	}

	if err := cs.state.ValidateBlock(p.Block()); err != nil {
		cs.logger.Warn("Invalid block", "proposal", p, "err", err)
		return
	}

	cs.currentState.onSetProposal(p)
}

func (cs *consensus) doSetProposal(p *proposal.Proposal) {
	cs.logger.Info("Proposal set", "proposal", p)
	cs.pendingVotes.SetRoundProposal(p.Round(), p)
}

func (cs *consensus) handleTimeout(t *ticker) {
	cs.lk.Lock()
	defer cs.lk.Unlock()

	cs.logger.Trace("Handle ticker", "ticker", t)

	// Old tickers might trigged now. Ignore them
	if cs.height != t.Height || cs.round != t.Round {
		cs.logger.Trace("Stale ticker", "ticker", t)
		return
	}

	cs.currentState.onTimedout(t)
}

func (cs *consensus) AddVote(v *vote.Vote) {
	cs.lk.Lock()
	defer cs.lk.Unlock()

	if v.Height() != cs.height {
		cs.logger.Trace("Vote has invalid height", "vote", v)
		return
	}

	if cs.pendingVotes.HasVote(v.Hash()) {
		cs.logger.Trace("Vote exists", "vote", v)
		return
	}

	cs.currentState.onAddVote(v)
}

func (cs *consensus) doAddVote(v *vote.Vote) {
	err := cs.pendingVotes.AddVote(v)
	if err != nil {
		cs.logger.Error("Error on adding a vote", "vote", v, "err", err)
	}

	cs.logger.Debug("New vote added", "vote", v)
}

func (cs *consensus) proposer(round int) *validator.Validator {
	return cs.state.Proposer(round)
}

func (cs *consensus) signAddVote(msgType vote.VoteType, hash crypto.Hash) {
	address := cs.signer.Address()
	if !cs.pendingVotes.CanVote(address) {
		cs.logger.Trace("This node is not in committee", "addr", address)
		return
	}

	// Sign the vote
	v := vote.NewVote(msgType, cs.height, cs.round, hash, address)
	cs.signer.SignMsg(v)
	cs.logger.Info("Our vote signed and broadcasted", "vote", v)

	err := cs.pendingVotes.AddVote(v)
	if err != nil {
		cs.logger.Error("Error on adding our vote!", "err", err, "vote", v)
	} else {
		cs.broadcastVote(v)
	}
}

func (cs *consensus) queryProposal() {
	pld := payload.NewQueryProposalPayload(cs.height, cs.round)
	cs.broadcastCh <- pld
}

func (cs *consensus) broadcastProposal(p *proposal.Proposal) {
	pld := payload.NewProposalPayload(*p)
	cs.broadcastCh <- pld
}

func (cs *consensus) broadcastVote(v *vote.Vote) {
	pld := payload.NewVotePayload(*v)
	cs.broadcastCh <- pld
}

func (cs *consensus) announceNewBlock(h int, b *block.Block, c *block.Certificate) {
	pld := payload.NewBlockAnnouncePayload(h, b, c)
	cs.broadcastCh <- pld
}

func (cs *consensus) PickRandomVote() *vote.Vote {
	cs.lk.RLock()
	defer cs.lk.RUnlock()

	round := util.RandInt(cs.round + 1)
	rv := cs.pendingVotes.MustGetRoundVotes(round)
	votes := rv.AllVotes()
	if len(votes) == 0 {
		return nil
	}
	r := util.RandInt(len(votes))
	return votes[r]
}
