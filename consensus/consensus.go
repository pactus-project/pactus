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
	state           state.Facade // TODO: rename `state` to `bcState` (blockchain state)
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
	state state.Facade,
	valKey *bls.ValidatorKey,
	rewardAddr crypto.Address,
	broadcastCh chan message.Message,
	mediator mediator,
) Consensus {
	broadcaster := func(_ crypto.Address, msg message.Message) {
		broadcastCh <- msg
	}
	return newConsensus(conf, state,
		valKey, rewardAddr, broadcaster, mediator)
}

func newConsensus(
	conf *Config,
	state state.Facade,
	valKey *bls.ValidatorKey,
	rewardAddr crypto.Address,
	broadcaster broadcaster,
	mediator mediator,
) *consensus {
	cs := &consensus{
		config:      conf,
		state:       state,
		broadcaster: broadcaster,
		valKey:      valKey,
	}

	// Update height later, See enterNewHeight.
	cs.log = log.NewLog()
	cs.logger = logger.NewSubLogger("_consensus", cs)
	cs.rewardAddr = rewardAddr

	cs.newHeightState = &newHeightState{cs}
	cs.proposeState = &proposeState{cs}
	cs.prepareState = &prepareState{cs, false}
	cs.precommitState = &precommitState{cs, false}
	cs.commitState = &commitState{cs}
	cs.cpPreVoteState = &cpPreVoteState{cs}
	cs.cpMainVoteState = &cpMainVoteState{cs}
	cs.cpDecideState = &cpDecideState{cs}
	cs.currentState = cs.newHeightState
	cs.mediator = mediator

	cs.height = 0
	cs.round = 0
	cs.active = false
	cs.mediator = mediator

	mediator.Register(cs)

	logger.Info("consensus instance created",
		"validator address", valKey.Address().String(),
		"reward address", rewardAddr.String())

	return cs
}

func (cs *consensus) String() string {
	return fmt.Sprintf("{%s %d/%d/%s/%d}",
		cs.valKey.Address().ShortString(),
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

func (cs *consensus) RoundProposal(round int16) *proposal.Proposal {
	cs.lk.RLock()
	defer cs.lk.RUnlock()

	return cs.log.RoundProposal(round)
}

func (cs *consensus) HasVote(hash hash.Hash) bool {
	cs.lk.RLock()
	defer cs.lk.RUnlock()

	return cs.log.HasVote(hash)
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

	cs.enterNewState(cs.newHeightState)
}

func (cs *consensus) scheduleTimeout(duration time.Duration, height uint32, round int16, target tickerTarget) {
	ti := &ticker{duration, height, round, target}
	timer := time.NewTimer(duration)
	cs.logger.Trace("new timer scheduled ⏱️", "duration", duration, "height", height, "round", round, "target", target)

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

	if p.Round() > cs.round+2 {
		cs.logger.Trace("proposal round number exceeding the round limit", "proposal", p)
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

	if p.Height() == cs.state.LastBlockHeight() {
		// A slow node might receive a proposal after committing the proposed block.
		// In this case, we accept the proposal and allow nodes to continue.
		// By doing so, we enable the validator to broadcast its votes and
		// prevent it from being marked as absent in the block certificate.
		cs.logger.Trace("block is committed for this height", "proposal", p)
		if p.Block().Hash() != cs.state.LastBlockHash() {
			cs.logger.Warn("proposal is not for the committed block", "proposal", p)
			return
		}
	} else {
		proposer := cs.proposer(p.Round())
		if err := p.Verify(proposer.PublicKey()); err != nil {
			cs.logger.Warn("proposal has invalid signature", "proposal", p, "error", err)
			return
		}

		if err := cs.state.ValidateBlock(p.Block()); err != nil {
			cs.logger.Warn("invalid block", "proposal", p, "error", err)
			return
		}
	}

	cs.logger.Info("proposal set", "proposal", p)
	cs.log.SetRoundProposal(p.Round(), p)

	cs.currentState.onSetProposal(p)
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

	if v.Round() > cs.round+2 {
		cs.logger.Trace("vote round number exceeding the round limit", "vote", v)
		return
	}

	if v.Type() == vote.VoteTypeCPPreVote ||
		v.Type() == vote.VoteTypeCPMainVote {
		err := cs.checkJust(v)
		if err != nil {
			cs.logger.Error("error on adding a cp vote", "vote", v, "error", err)
			return
		}

		if v.Round() == cs.round && cs.cpWeakValidity == nil {
			if v.CPValue() == vote.CPValueZero ||
				v.CPValue() == vote.CPValueAbstain {
				bh := v.BlockHash()
				cs.cpWeakValidity = &bh

				roundProposal := cs.log.RoundProposal(cs.round)

				if roundProposal != nil &&
					roundProposal.Block().Hash() != bh {
					cs.logger.Warn("double proposal detected",
						"prepared", bh.ShortString(),
						"roundProposal", roundProposal.Block().Hash().ShortString())

					cs.log.SetRoundProposal(cs.round, nil)
					cs.queryProposal()
				}
			}
		}
	}

	added, err := cs.log.AddVote(v)
	if err != nil {
		cs.logger.Error("error on adding a vote", "vote", v, "error", err)
	}
	if added {
		cs.logger.Debug("new vote added", "vote", v)

		cs.currentState.onAddVote(v)
	}
}

func (cs *consensus) checkCPValue(vote *vote.Vote, allowedValues ...vote.CPValue) error {
	for _, v := range allowedValues {
		if vote.CPValue() == v {
			return nil
		}
	}

	return invalidJustificationError{
		JustType: vote.CPJust().Type(),
		Reason:   fmt.Sprintf("invalid value: %v", vote.CPValue()),
	}
}

func (cs *consensus) checkJustInitZero(just vote.Just, blockHash hash.Hash) error {
	j, ok := just.(*vote.JustInitZero)
	if !ok {
		return invalidJustificationError{
			JustType: just.Type(),
			Reason:   "invalid just data",
		}
	}

	sb := certificate.BlockCertificateSignBytes(blockHash,
		j.QCert.Height(),
		j.QCert.Round())
	sb = append(sb, util.StringToBytes(vote.VoteTypePrepare.String())...)

	err := j.QCert.Validate(cs.height, cs.validators, sb)
	if err != nil {
		return invalidJustificationError{
			JustType: j.Type(),
			Reason:   err.Error(),
		}
	}
	return nil
}

func (cs *consensus) checkJustInitOne(just vote.Just) error {
	_, ok := just.(*vote.JustInitOne)
	if !ok {
		return invalidJustificationError{
			JustType: just.Type(),
			Reason:   "invalid just data",
		}
	}

	return nil
}

func (cs *consensus) checkJustPreVoteHard(just vote.Just,
	blockHash hash.Hash, cpRound int16, cpValue vote.CPValue,
) error {
	j, ok := just.(*vote.JustPreVoteHard)
	if !ok {
		return invalidJustificationError{
			JustType: just.Type(),
			Reason:   "invalid just data",
		}
	}

	sb := certificate.BlockCertificateSignBytes(blockHash,
		j.QCert.Height(),
		j.QCert.Round())
	sb = append(sb, util.StringToBytes(vote.VoteTypeCPPreVote.String())...)
	sb = append(sb, util.Int16ToSlice(cpRound-1)...)
	sb = append(sb, byte(cpValue))

	err := j.QCert.Validate(cs.height, cs.validators, sb)
	if err != nil {
		return invalidJustificationError{
			JustType: just.Type(),
			Reason:   err.Error(),
		}
	}
	return nil
}

func (cs *consensus) checkJustPreVoteSoft(just vote.Just,
	blockHash hash.Hash, cpRound int16,
) error {
	j, ok := just.(*vote.JustPreVoteSoft)
	if !ok {
		return invalidJustificationError{
			JustType: just.Type(),
			Reason:   "invalid just data",
		}
	}

	sb := certificate.BlockCertificateSignBytes(blockHash,
		j.QCert.Height(),
		j.QCert.Round())
	sb = append(sb, util.StringToBytes(vote.VoteTypeCPMainVote.String())...)
	sb = append(sb, util.Int16ToSlice(cpRound-1)...)
	sb = append(sb, byte(vote.CPValueAbstain))

	err := j.QCert.Validate(cs.height, cs.validators, sb)
	if err != nil {
		return invalidJustificationError{
			JustType: just.Type(),
			Reason:   err.Error(),
		}
	}
	return nil
}

func (cs *consensus) checkJustMainVoteNoConflict(just vote.Just,
	blockHash hash.Hash, cpRound int16, cpValue vote.CPValue,
) error {
	j, ok := just.(*vote.JustMainVoteNoConflict)
	if !ok {
		return invalidJustificationError{
			JustType: just.Type(),
			Reason:   "invalid just data",
		}
	}

	sb := certificate.BlockCertificateSignBytes(blockHash,
		j.QCert.Height(),
		j.QCert.Round())
	sb = append(sb, util.StringToBytes(vote.VoteTypeCPPreVote.String())...)
	sb = append(sb, util.Int16ToSlice(cpRound)...)
	sb = append(sb, byte(cpValue))

	err := j.QCert.Validate(cs.height, cs.validators, sb)
	if err != nil {
		return invalidJustificationError{
			JustType: j.Type(),
			Reason:   err.Error(),
		}
	}
	return nil
}

func (cs *consensus) checkJustMainVoteConflict(just vote.Just,
	blockHash hash.Hash, cpRound int16,
) error {
	j, ok := just.(*vote.JustMainVoteConflict)
	if !ok {
		return invalidJustificationError{
			JustType: just.Type(),
			Reason:   "invalid just data",
		}
	}

	if cpRound == 0 {
		switch j.Just0.Type() {
		case vote.JustTypeInitZero:
			err := cs.checkJustInitZero(j.Just0, blockHash)
			if err != nil {
				return err
			}

		default:
			return invalidJustificationError{
				JustType: just.Type(),
				Reason:   fmt.Sprintf("unexpected justification: %s", j.Just0.Type()),
			}
		}

		switch j.Just1.Type() {
		case vote.JustTypeInitOne:
			err := cs.checkJustInitOne(j.Just1)
			if err != nil {
				return err
			}

		default:
			return invalidJustificationError{
				JustType: just.Type(),
				Reason:   fmt.Sprintf("unexpected justification: %s", j.Just1.Type()),
			}
		}
	} else {
		switch j.Just0.Type() {
		case vote.JustTypePreVoteSoft:
			err := cs.checkJustPreVoteSoft(j.Just0, blockHash, cpRound)
			if err != nil {
				return err
			}
		case vote.JustTypePreVoteHard:
			err := cs.checkJustPreVoteHard(j.Just0, blockHash, cpRound, vote.CPValueZero)
			if err != nil {
				return err
			}
		default:
			return invalidJustificationError{
				JustType: just.Type(),
				Reason:   fmt.Sprintf("unexpected justification: %s", j.Just0.Type()),
			}
		}

		switch j.Just1.Type() {
		case vote.JustTypePreVoteHard:
			err := cs.checkJustPreVoteHard(j.Just1, hash.UndefHash, cpRound, vote.CPValueOne)
			if err != nil {
				return err
			}

		default:
			return invalidJustificationError{
				JustType: just.Type(),
				Reason:   fmt.Sprintf("unexpected justification: %s", j.Just1.Type()),
			}
		}
	}

	return nil
}

func (cs *consensus) checkJustPreVote(v *vote.Vote) error {
	just := v.CPJust()
	if v.CPRound() == 0 {
		switch just.Type() {
		case vote.JustTypeInitZero:
			err := cs.checkCPValue(v, vote.CPValueZero)
			if err != nil {
				return err
			}
			return cs.checkJustInitZero(just, v.BlockHash())

		case vote.JustTypeInitOne:
			err := cs.checkCPValue(v, vote.CPValueOne)
			if err != nil {
				return err
			}
			if v.BlockHash() != hash.UndefHash {
				return invalidJustificationError{
					JustType: just.Type(),
					Reason:   "invalid block hash",
				}
			}
			return cs.checkJustInitOne(just)
		default:
			return invalidJustificationError{
				JustType: just.Type(),
				Reason:   "invalid pre-vote justification",
			}
		}
	} else {
		switch just.Type() {
		case vote.JustTypePreVoteSoft:
			err := cs.checkCPValue(v, vote.CPValueZero, vote.CPValueOne)
			if err != nil {
				return err
			}
			return cs.checkJustPreVoteSoft(just, v.BlockHash(), v.CPRound())

		case vote.JustTypePreVoteHard:
			err := cs.checkCPValue(v, vote.CPValueZero, vote.CPValueOne)
			if err != nil {
				return err
			}
			return cs.checkJustPreVoteHard(just, v.BlockHash(), v.CPRound(), v.CPValue())

		default:
			return invalidJustificationError{
				JustType: just.Type(),
				Reason:   "invalid pre-vote justification",
			}
		}
	}
}

func (cs *consensus) checkJustMainVote(v *vote.Vote) error {
	just := v.CPJust()
	switch just.Type() {
	case vote.JustTypeMainVoteNoConflict:
		err := cs.checkCPValue(v, vote.CPValueZero, vote.CPValueOne)
		if err != nil {
			return err
		}
		return cs.checkJustMainVoteNoConflict(just, v.BlockHash(), v.CPRound(), v.CPValue())

	case vote.JustTypeMainVoteConflict:
		err := cs.checkCPValue(v, vote.CPValueAbstain)
		if err != nil {
			return err
		}
		return cs.checkJustMainVoteConflict(just, v.BlockHash(), v.CPRound())

	default:
		return invalidJustificationError{
			JustType: just.Type(),
			Reason:   "invalid main-vote justification",
		}
	}
}

func (cs *consensus) checkJust(v *vote.Vote) error {
	if v.Type() == vote.VoteTypeCPPreVote {
		return cs.checkJustPreVote(v)
	} else if v.Type() == vote.VoteTypeCPMainVote {
		return cs.checkJustMainVote(v)
	} else {
		panic("unreachable")
	}
}

func (cs *consensus) proposer(round int16) *validator.Validator {
	return cs.state.Proposer(round)
}

func (cs *consensus) signAddCPPreVote(hash hash.Hash,
	cpRound int16, cpValue vote.CPValue, just vote.Just,
) {
	v := vote.NewCPPreVote(hash, cs.height,
		cs.round, cpRound, cpValue, just, cs.valKey.Address())
	cs.signAddVote(v)
}

func (cs *consensus) signAddCPMainVote(hash hash.Hash,
	cpRound int16, cpValue vote.CPValue, just vote.Just,
) {
	v := vote.NewCPMainVote(hash, cs.height, cs.round,
		cpRound, cpValue, just, cs.valKey.Address())
	cs.signAddVote(v)
}

func (cs *consensus) signAddPrepareVote(hash hash.Hash) {
	v := vote.NewPrepareVote(hash, cs.height, cs.round, cs.valKey.Address())
	cs.signAddVote(v)
}

func (cs *consensus) signAddPrecommitVote(hash hash.Hash) {
	v := vote.NewPrecommitVote(hash, cs.height, cs.round, cs.valKey.Address())
	cs.signAddVote(v)
}

func (cs *consensus) signAddVote(v *vote.Vote) {
	sig := cs.valKey.Sign(v.SignBytes())
	v.SetSignature(sig)
	cs.logger.Info("our vote signed and broadcasted", "vote", v)

	_, err := cs.log.AddVote(v)
	if err != nil {
		cs.logger.Error("error on adding our vote", "error", err, "vote", v)
	}
	cs.broadcastVote(v)
}

func (cs *consensus) queryProposal() {
	cs.broadcaster(cs.valKey.Address(),
		message.NewQueryProposalMessage(cs.height, cs.round))
}

func (cs *consensus) queryVotes() {
	cs.broadcaster(cs.valKey.Address(),
		message.NewQueryVotesMessage(cs.height, cs.round))
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

func (cs *consensus) announceNewBlock(h uint32, b *block.Block, c *certificate.Certificate) {
	go cs.mediator.OnBlockAnnounce(cs)
	cs.broadcaster(cs.valKey.Address(),
		message.NewBlockAnnounceMessage(h, b, c))
}

func (cs *consensus) makeCertificate(votes map[crypto.Address]*vote.Vote) *certificate.Certificate {
	vals := cs.state.CommitteeValidators()
	committers := make([]int32, len(vals))
	absentees := make([]int32, 0)
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

	aggSig := bls.SignatureAggregate(sigs...)

	return certificate.NewCertificate(cs.height, cs.round, committers, absentees, aggSig)
}

// IsActive checks if the consensus is in an active state and participating in the consensus algorithm.
func (cs *consensus) IsActive() bool {
	cs.lk.RLock()
	defer cs.lk.RUnlock()

	return cs.active
}

// TODO: Improve the performance?
func (cs *consensus) PickRandomVote(round int16) *vote.Vote {
	cs.lk.RLock()
	defer cs.lk.RUnlock()

	votes := []*vote.Vote{}
	if round == cs.round {
		m := cs.log.RoundMessages(round)
		votes = append(votes, m.AllVotes()...)
	} else {
		// Don't broadcast prepare and precommit votes for previous rounds
		vs0 := cs.log.CPPreVoteVoteSet(round)
		vs1 := cs.log.CPMainVoteVoteSet(round)
		votes = append(votes, vs0.AllVotes()...)
		votes = append(votes, vs1.AllVotes()...)
	}
	if len(votes) == 0 {
		return nil
	}
	return votes[util.RandInt32(int32(len(votes)))]
}

func (cs *consensus) startChangingProposer() {
	// It is not timeout before
	if cs.cpDecided == -1 {
		cs.logger.Debug("changing proposer started")
		cs.enterNewState(cs.cpPreVoteState)
	}
}
