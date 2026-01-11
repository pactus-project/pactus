package consensus

import (
	"fmt"
	"sync"
	"time"

	"github.com/pactus-project/pactus/consensus/log"
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/state"
	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/pactus-project/pactus/types/block"
	"github.com/pactus-project/pactus/types/certificate"
	"github.com/pactus-project/pactus/types/proposal"
	"github.com/pactus-project/pactus/types/validator"
	"github.com/pactus-project/pactus/types/vote"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/logger"
	"github.com/pactus-project/pactus/util/pipeline"
)

type broadcaster func(crypto.Address, message.Message)

type consensus struct {
	lk sync.RWMutex

	config          *Config
	logger          *logger.SubLogger
	log             *log.Log
	validators      []*validator.Validator
	cpWeakValidity  *hash.Hash // The change proposer's weak validity that is a prepared block hash
	cpDecided       int
	height          uint32
	round           int16
	cpRound         int16
	valKey          *bls.ValidatorKey
	rewardAddr      crypto.Address
	bcState         state.Facade // Blockchain state
	changeProposer  *changeProposer
	newHeightState  consState
	proposeState    consState
	prepareState    consState
	precommitState  consState
	commitState     consState
	cpPreVoteState  consState
	cpMainVoteState consState
	cpDecideState   consState
	currentState    consState
	broadcaster     broadcaster
	mediator        mediator
	active          bool
}

func NewConsensus(
	conf *Config,
	bcState state.Facade,
	valKey *bls.ValidatorKey,
	rewardAddr crypto.Address,
	broadcastPipe pipeline.Pipeline[message.Message],
	mediator mediator,
) Consensus {
	broadcaster := func(_ crypto.Address, msg message.Message) {
		broadcastPipe.Send(msg)
	}

	return makeConsensus(conf, bcState,
		valKey, rewardAddr, broadcaster, mediator)
}

func makeConsensus(
	conf *Config,
	bcState state.Facade,
	valKey *bls.ValidatorKey,
	rewardAddr crypto.Address,
	broadcaster broadcaster,
	mediator mediator,
) *consensus {
	cons := &consensus{
		config:      conf,
		bcState:     bcState,
		broadcaster: broadcaster,
		valKey:      valKey,
	}

	// Update height later, See enterNewHeight.
	cons.log = log.NewLog()
	cons.logger = logger.NewSubLogger("_consensus", cons)
	cons.rewardAddr = rewardAddr

	cons.changeProposer = &changeProposer{cons}
	cons.newHeightState = &newHeightState{cons}
	cons.proposeState = &proposeState{cons}
	cons.prepareState = &prepareState{cons, false}
	cons.precommitState = &precommitState{cons, false}
	cons.commitState = &commitState{cons}
	cons.cpPreVoteState = &cpPreVoteState{cons.changeProposer}
	cons.cpMainVoteState = &cpMainVoteState{cons.changeProposer}
	cons.cpDecideState = &cpDecideState{cons.changeProposer}
	cons.currentState = cons.newHeightState
	cons.mediator = mediator

	cons.height = 0
	cons.round = 0
	cons.active = false
	cons.mediator = mediator

	mediator.Register(cons)

	logger.Info("consensus instance created (Legacy)",
		"validator address", valKey.Address().String(),
		"reward address", rewardAddr.String())

	return cons
}

// LogString returns a concise string representation intended for use in logs.
func (cs *consensus) LogString() string {
	return fmt.Sprintf("{%s %d/%d/%s/%d}",
		cs.valKey.Address(),
		cs.height, cs.round, cs.currentState.name(), cs.cpRound)
}

func (cs *consensus) ConsensusKey() *bls.PublicKey {
	cs.lk.RLock()
	defer cs.lk.RUnlock()

	return cs.valKey.PublicKey()
}

func (cs *consensus) HeightRound() (uint32, int16) {
	cs.lk.RLock()
	defer cs.lk.RUnlock()

	return cs.height, cs.round
}

func (cs *consensus) HasVote(h hash.Hash) bool {
	cs.lk.RLock()
	defer cs.lk.RUnlock()

	return cs.log.HasVote(h)
}

// AllVotes returns all valid votes inside the consensus log up to and including
// the current consensus round.
// Valid votes from subsequent rounds are not included.
func (cs *consensus) AllVotes() []*vote.Vote {
	cs.lk.RLock()
	defer cs.lk.RUnlock()

	votes := []*vote.Vote{}
	for r := int16(0); r <= cs.round; r++ {
		m := cs.log.RoundMessages(r)
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

	if cs.IsDeprecated() {
		return
	}

	stateHeight := cs.bcState.LastBlockHeight()
	if cs.height != stateHeight+1 {
		cs.enterNewState(cs.newHeightState)
	}
}

func (cs *consensus) scheduleTimeout(duration time.Duration, height uint32, round int16, target tickerTarget) {
	ticker := &ticker{duration, height, round, target}
	timer := time.NewTimer(duration)
	cs.logger.Trace("new timer scheduled ⏱️", "duration", duration, "height", height, "round", round, "target", target)

	go func() {
		<-timer.C
		cs.handleTimeout(ticker)
	}()
}

func (cs *consensus) handleTimeout(ticker *ticker) {
	cs.lk.Lock()
	defer cs.lk.Unlock()

	cs.logger.Trace("handle ticker", "ticker", ticker)

	// Old tickers might be triggered now. Ignore them.
	if cs.height != ticker.Height || cs.round != ticker.Round {
		cs.logger.Trace("stale ticker", "ticker", ticker)

		return
	}

	cs.logger.Debug("timer expired", "ticker", ticker)
	cs.currentState.onTimeout(ticker)
}

func (cs *consensus) SetProposal(prop *proposal.Proposal) {
	cs.lk.Lock()
	defer cs.lk.Unlock()

	if !cs.active {
		return
	}

	if prop.Height() != cs.height {
		return
	}

	if prop.Round() < cs.round {
		cs.logger.Debug("proposal for expired round", "proposal", prop)

		return
	}

	if err := prop.BasicCheck(); err != nil {
		cs.logger.Warn("invalid proposal", "proposal", prop, "error", err)

		return
	}

	roundProposal := cs.log.RoundProposal(prop.Round())
	if roundProposal != nil {
		cs.logger.Trace("this round has proposal", "proposal", prop)

		return
	}

	if prop.Height() == cs.bcState.LastBlockHeight() {
		// A slow node might receive a proposal after committing the proposed block.
		// In this case, we accept the proposal and allow nodes to continue.
		// By doing so, we enable the validator to broadcast its votes and
		// prevent it from being marked as absent in the block certificate.
		cs.logger.Warn("block committed before receiving proposal", "proposal", prop)
		if prop.Block().Hash() != cs.bcState.LastBlockHash() {
			cs.logger.Warn("proposal is not for the committed block", "proposal", prop)

			return
		}
	} else {
		proposer := cs.proposer(prop.Round())
		if err := prop.Verify(proposer.PublicKey()); err != nil {
			cs.logger.Warn("proposal is invalid", "proposal", prop, "error", err)

			return
		}

		if err := cs.bcState.ValidateBlock(prop.Block(), prop.Round()); err != nil {
			cs.logger.Warn("invalid block", "proposal", prop, "error", err)

			return
		}
	}

	cs.logger.Info("proposal set", "proposal", prop)
	cs.log.SetRoundProposal(prop.Round(), prop)

	cs.currentState.onSetProposal(prop)
}

func (cs *consensus) AddVote(vte *vote.Vote) {
	cs.lk.Lock()
	defer cs.lk.Unlock()

	if !cs.active {
		return
	}

	if vte.Height() != cs.height {
		return
	}

	if vte.Round() < cs.round {
		cs.logger.Debug("vote for expired round", "vote", vte)

		return
	}

	if vte.Type() == vote.VoteTypeCPPreVote ||
		vte.Type() == vote.VoteTypeCPMainVote ||
		vte.Type() == vote.VoteTypeCPDecided {
		err := cs.changeProposer.checkJust(vte)
		if err != nil {
			cs.logger.Error("error on adding a cp vote", "vote", vte, "error", err)

			return
		}
	}

	added, err := cs.log.AddVote(vte)
	if err != nil {
		cs.logger.Error("error on adding a vote", "vote", vte, "error", err)
	}
	if added {
		cs.logger.Info("new vote added", "vote", vte)

		cs.currentState.onAddVote(vte)
	}
}

func (cs *consensus) proposer(round int16) *validator.Validator {
	return cs.bcState.Proposer(round)
}

func (cs *consensus) IsProposer() bool {
	cs.lk.RLock()
	defer cs.lk.RUnlock()

	return cs.isProposer()
}

func (cs *consensus) isProposer() bool {
	return cs.proposer(cs.round).Address() == cs.valKey.Address()
}

func (cs *consensus) signAddCPPreVote(h hash.Hash,
	cpRound int16, cpValue vote.CPValue, just vote.Just,
) {
	v := vote.NewCPPreVote(h, cs.height,
		cs.round, cpRound, cpValue, just, cs.valKey.Address())
	cs.signAddVote(v)
}

func (cs *consensus) signAddCPMainVote(h hash.Hash,
	cpRound int16, cpValue vote.CPValue, just vote.Just,
) {
	v := vote.NewCPMainVote(h, cs.height, cs.round,
		cpRound, cpValue, just, cs.valKey.Address())
	cs.signAddVote(v)
}

func (cs *consensus) signAddCPDecidedVote(h hash.Hash,
	cpRound int16, cpValue vote.CPValue, just vote.Just,
) {
	v := vote.NewCPDecidedVote(h, cs.height, cs.round,
		cpRound, cpValue, just, cs.valKey.Address())
	cs.signAddVote(v)
}

func (cs *consensus) signAddPrepareVote(h hash.Hash) {
	v := vote.NewPrepareVote(h, cs.height, cs.round, cs.valKey.Address())
	cs.signAddVote(v)
}

func (cs *consensus) signAddPrecommitVote(h hash.Hash) {
	v := vote.NewPrecommitVote(h, cs.height, cs.round, cs.valKey.Address())
	cs.signAddVote(v)
}

func (cs *consensus) signAddVote(vte *vote.Vote) {
	sig := cs.valKey.Sign(vte.SignBytes())
	vte.SetSignature(sig)
	cs.logger.Info("our vote signed and broadcasted", "vote", vte)

	_, err := cs.log.AddVote(vte)
	if err != nil {
		cs.logger.Warn("error on adding our vote", "error", err, "vote", vte)
	}
	cs.broadcastVote(vte)
}

// queryProposal requests any missing proposal from other validators.
func (cs *consensus) queryProposal() {
	cs.broadcaster(cs.valKey.Address(),
		message.NewQueryProposalMessage(cs.height, cs.round, cs.valKey.Address()))
}

// queryVote requests any missing votes from other validators.
func (cs *consensus) queryVote() {
	cs.broadcaster(cs.valKey.Address(),
		message.NewQueryVoteMessage(cs.height, cs.round, cs.valKey.Address()))
}

func (cs *consensus) broadcastProposal(p *proposal.Proposal) {
	go cs.mediator.OnPublishProposal(cs, p)
	cs.broadcaster(cs.valKey.Address(),
		message.NewProposalMessage(p))
}

func (cs *consensus) broadcastVote(v *vote.Vote) {
	go cs.mediator.OnPublishVote(cs, v)
	cs.broadcaster(cs.valKey.Address(),
		message.NewVoteMessage(v))
}

func (cs *consensus) announceNewBlock(blk *block.Block,
	cert *certificate.Certificate, proof *certificate.Certificate,
) {
	go cs.mediator.OnBlockAnnounce(cs)
	cs.broadcaster(cs.valKey.Address(),
		message.NewBlockAnnounceMessage(blk, cert, proof))
}

func (cs *consensus) makeCertificate(votes map[crypto.Address]*vote.Vote,
) *certificate.Certificate {
	cert := certificate.NewCertificate(cs.height, cs.round)
	cert.SetSignature(cs.signersInfo(votes))

	return cert
}

// signersInfo processes a map of votes from validators and provides these information:
// - A list of all validators' numbers eligible to vote in this step.
// - A list of absentee validators' numbers who did not vote in this step.
// - An aggregated signature generated from the signatures of participating validators.
func (cs *consensus) signersInfo(votes map[crypto.Address]*vote.Vote) (
	committers, absentees []int32, aggSig *bls.Signature,
) {
	vals := cs.validators
	committers = make([]int32, len(vals))
	absentees = make([]int32, 0)
	sigs := make([]*bls.Signature, 0)

	for i, val := range vals {
		vote := votes[val.Address()]
		if vote != nil {
			sigs = append(sigs, vote.Signature())
		} else {
			absentees = append(absentees, val.Number())
		}

		committers[i] = val.Number()
	}

	aggSig, _ = bls.SignatureAggregate(sigs...)

	return committers, absentees, aggSig
}

// IsActive checks if the consensus is in an active state and participating in the consensus algorithm.
func (cs *consensus) IsActive() bool {
	cs.lk.RLock()
	defer cs.lk.RUnlock()

	if cs.IsDeprecated() {
		return false
	}

	return cs.active
}

func (cs *consensus) Proposal() *proposal.Proposal {
	cs.lk.RLock()
	defer cs.lk.RUnlock()

	return cs.log.RoundProposal(cs.round)
}

func (cs *consensus) HandleQueryProposal(height uint32, round int16) *proposal.Proposal {
	cs.lk.RLock()
	defer cs.lk.RUnlock()

	if !cs.active {
		return nil
	}

	if height != cs.height {
		return nil
	}

	if round != cs.round {
		return nil
	}

	if cs.isProposer() {
		return cs.log.RoundProposal(cs.round)
	}

	if cs.cpDecided == 0 {
		// It is decided not to change the proposer and the proposal is locked.
		// Locked proposals can be sent by all validators.
		// This helps prevent a situation where the proposer goes offline after proposing the block.
		return cs.log.RoundProposal(cs.round)
	}

	return nil
}

// TODO: Improve the performance?
func (cs *consensus) HandleQueryVote(height uint32, round int16) *vote.Vote {
	cs.lk.RLock()
	defer cs.lk.RUnlock()

	if !cs.active {
		return nil
	}

	if height != cs.height {
		return nil
	}

	votes := []*vote.Vote{}
	switch {
	case round < cs.round:
		// Past round: Only broadcast cp:decided votes
		vs := cs.log.CPDecidedVoteSet(round)
		votes = append(votes, vs.AllVotes()...)

	case round == cs.round:
		// Current round
		m := cs.log.RoundMessages(round)
		votes = append(votes, m.AllVotes()...)

	case round > cs.round:
		// Future round
	}

	if len(votes) == 0 {
		return nil
	}

	return votes[util.RandInt32(int32(len(votes)))]
}

func (cs *consensus) startChangingProposer() {
	// If it is not decided yet.
	// TODO: can we remove this condition in new consensus model?
	if cs.cpDecided == -1 {
		cs.logger.Info("changing proposer started",
			"cpRound", cs.cpRound, "proposer", cs.proposer(cs.round).Address())
		cs.enterNewState(cs.cpPreVoteState)
	}
}

func (*consensus) IsDeprecated() bool {
	return false
}
