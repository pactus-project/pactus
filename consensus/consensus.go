package consensus

import (
	"fmt"
	"time"

	"github.com/sasha-s/go-deadlock"
	"github.com/zarbchain/zarb-go/consensus/hrs"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/logger"
	"github.com/zarbchain/zarb-go/message"
	"github.com/zarbchain/zarb-go/state"
	"github.com/zarbchain/zarb-go/util"
	"github.com/zarbchain/zarb-go/validator"
	"github.com/zarbchain/zarb-go/vote"
)

//-----------------------------------------------------------------------------

var (
	msgQueueSize = 1000
)

type Consensus struct {
	lk deadlock.RWMutex

	config        *Config
	hrs           hrs.HRS
	votes         *HeightVoteSet
	valset        *validator.ValidatorSet
	privValidator *validator.PrivValidator
	isCommitted   bool
	state         state.State
	broadcastCh   chan message.Message
	logger        *logger.Logger
}

func NewConsensus(
	conf *Config,
	state state.State,
	privValidator *validator.PrivValidator,
	broadcastCh chan message.Message) (*Consensus, error) {
	cs := &Consensus{
		config:        conf,
		state:         state,
		valset:        state.ValidatorSet(),
		broadcastCh:   broadcastCh,
		privValidator: privValidator,
	}

	// Update height later, See enterNewHeight.
	cs.votes = NewHeightVoteSet(-1, cs.valset)
	cs.hrs = hrs.NewHRS(-1, -1, hrs.StepTypeNewHeight)
	cs.logger = logger.NewLogger("_consensus", cs)

	return cs, nil
}

func (cs *Consensus) Fingerprint() string {
	return fmt.Sprintf("{%v}",
		cs.hrs.Fingerprint())
}

func (cs *Consensus) HRS() hrs.HRS {
	cs.lk.RLock()
	defer cs.lk.RUnlock()

	return cs.hrs
}

func (cs *Consensus) updateRoundStep(round int, step hrs.StepType) {
	cs.hrs.UpdateRoundStep(round, step)

	hasProposal := cs.votes.HasRoundProposal(cs.hrs.Round())
	msg := message.NewHeartBeatMessage(cs.state.LastBlockHash(), cs.hrs, hasProposal)
	cs.broadcastCh <- msg
}

func (cs *Consensus) updateHeight(height int) {
	cs.hrs.UpdateHeight(height)
}

func (cs *Consensus) HasProposal() bool {
	cs.lk.RLock()
	defer cs.lk.RUnlock()

	return cs.votes.HasRoundProposal(cs.hrs.Round())
}

func (cs *Consensus) LastProposal() *vote.Proposal {
	cs.lk.RLock()
	defer cs.lk.RUnlock()

	return cs.votes.RoundProposal(cs.hrs.Round())
}

func (cs *Consensus) AllVotes() []*vote.Vote {
	cs.lk.Lock()
	defer cs.lk.Unlock()

	votes := cs.votes.votes
	slice := make([]*vote.Vote, len(votes))
	i := 0
	for _, v := range votes {
		slice[i] = v
		i++
	}
	return slice
}

func (cs *Consensus) AllVotesHashes() []crypto.Hash {
	cs.lk.Lock()
	defer cs.lk.Unlock()

	votes := cs.votes.votes
	slice := make([]crypto.Hash, len(votes))
	i := 0
	for _, v := range votes {
		slice[i] = v.Hash()
		i++
	}
	return slice
}

func (cs *Consensus) Vote(h crypto.Hash) *vote.Vote {
	cs.lk.Lock()
	defer cs.lk.Unlock()

	return cs.votes.votes[h]
}

func (cs *Consensus) scheduleTimeout(duration time.Duration, height int, round int, step hrs.StepType) {
	to := timeout{duration, height, round, step}

	if cs.config.FuzzTesting {
		to.Duration = time.Duration(util.RandInt(8)) * time.Second
	}
	timer := time.NewTimer(duration)
	go func() {
		<-timer.C
		cs.handleTimeout(to)
	}()
	logger.Debug("Scheduled timeout", "dur", duration, "height", height, "round", round, "step", step)
}

func (cs *Consensus) invalidHeight(height int) bool {
	return cs.hrs.Height() != height
}

func (cs *Consensus) invalidHeightRound(height int, round int) bool {
	return cs.hrs.Height() != height || cs.hrs.Round() != round
}

func (cs *Consensus) invalidHeightRoundStep(height int, round int, step hrs.StepType) bool {
	return cs.hrs.Height() != height || cs.hrs.Round() != round || cs.hrs.Step() > step
}

func (cs *Consensus) AddVote(v *vote.Vote) error {
	cs.lk.Lock()
	defer cs.lk.Unlock()

	return cs.addVote(v)
}

func (cs *Consensus) SetProposal(proposal *vote.Proposal) {
	cs.lk.Lock()
	defer cs.lk.Unlock()

	cs.setProposal(proposal)
}

func (cs *Consensus) handleTimeout(ti timeout) {
	cs.lk.Lock()
	defer cs.lk.Unlock()

	cs.logger.Debug("Handle timeout", "timeout", ti)

	// timeouts must be for current height
	if ti.Height != cs.hrs.Height() {
		cs.logger.Debug("Ignoring timeout", "timeout", ti)
		return
	}

	switch ti.Step {
	case hrs.StepTypeNewHeight:
		cs.enterNewHeight(ti.Height + 1)
	case hrs.StepTypeNewRound:
		cs.enterNewRound(ti.Height, ti.Round+1)
	case hrs.StepTypePrevote:
		cs.enterPrevote(ti.Height, ti.Round)
	case hrs.StepTypePrecommit:
		cs.enterPrecommit(ti.Height, ti.Round)
	default:
		panic(fmt.Sprintf("Invalid timeout step: %v", ti.Step))
	}

}

//-----------------------------------------------------------------------------

func (cs *Consensus) addVote(v *vote.Vote) error {
	// Height mismatch is ignored.
	if cs.invalidHeight(v.Height()) {
		return errors.Errorf(errors.ErrInvalidVote, "Vote ignored, height mismatch: %v", v.Height())
	}

	added, err := cs.votes.AddVote(v)
	if err != nil {
		return err
	}
	if !added {
		// we probably had this vote before
		return nil
	}

	height := v.Height()
	round := v.Round()
	switch v.VoteType() {
	case vote.VoteTypePrevote:
		prevotes := cs.votes.Prevotes(round)
		cs.logger.Debug("Vote added to prevote", "vote", v, "voteset", prevotes)
		// current round
		if cs.hrs.Round() == round {
			if ok := prevotes.HasQuorum(); ok {
				blockHash := prevotes.QuorumBlock()
				if blockHash == nil {
					cs.enterPrevoteWait(height, round)
				} else if blockHash.IsUndef() {
					cs.enterPrecommit(height, round)
				} else {
					cs.enterPrecommit(height, round)
				}
			}
		}

	case vote.VoteTypePrecommit:
		precommits := cs.votes.Precommits(round)
		cs.logger.Debug("Vote added to precommit", "vote", v, "voteset", precommits)
		// current round
		if cs.hrs.Round() == round {
			if ok := precommits.HasQuorum(); ok {
				blockHash := precommits.QuorumBlock()
				if blockHash == nil {
					cs.enterPrecommitWait(height, round)
				} else if blockHash.IsUndef() {
					cs.enterNewRound(height, round+1)
				} else {
					cs.enterCommit(height, round)
				}
			}
		}

	default:
		cs.logger.Panic("Unexpected vote type %X", v.VoteType)
	}

	return err
}

func (cs *Consensus) signAddVote(msgType vote.VoteType, hash crypto.Hash) {
	address := cs.privValidator.Address()
	if !cs.valset.Contains(address) {
		cs.logger.Info("This node is not in validator set", "addr", address)
		return
	}

	// Sign the vote
	v := vote.NewVote(msgType, cs.hrs.Height(), cs.hrs.Round(), hash, address)
	cs.privValidator.SignMsg(v)
	err := cs.addVote(v)
	if err != nil {
		cs.logger.Error("Error on adding our vote!", "error", err)
	}

	// Broadcast our vote
	msg := message.NewVoteMessage(v)
	cs.broadcastCh <- msg
}
