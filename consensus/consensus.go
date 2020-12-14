package consensus

import (
	"fmt"
	"time"

	"github.com/sasha-s/go-deadlock"
	"github.com/zarbchain/zarb-go/consensus/hrs"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/logger"
	"github.com/zarbchain/zarb-go/message"
	"github.com/zarbchain/zarb-go/state"
	"github.com/zarbchain/zarb-go/util"
	"github.com/zarbchain/zarb-go/validator"
	"github.com/zarbchain/zarb-go/vote"
)

type Consensus struct {
	lk deadlock.RWMutex

	config      *Config
	hrs         hrs.HRS
	votes       *HeightVoteSet
	valset      *validator.ValidatorSet
	signer      crypto.Signer
	isCommitted bool
	state       state.State
	broadcastCh chan *message.Message
	logger      *logger.Logger
}

func NewConsensus(
	conf *Config,
	state state.State,
	signer crypto.Signer,
	broadcastCh chan *message.Message) (*Consensus, error) {
	cs := &Consensus{
		config:      conf,
		state:       state,
		valset:      state.ValidatorSet(),
		broadcastCh: broadcastCh,
		signer:      signer,
	}

	// Update height later, See enterNewHeight.
	cs.votes = NewHeightVoteSet(0, cs.valset)
	cs.hrs = hrs.NewHRS(0, 0, hrs.StepTypeNewHeight)
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
}

func (cs *Consensus) updateHeight(height int) {
	if cs.hrs.Height() != height {
		cs.votes.Reset(height)
		cs.hrs.UpdateHeight(height)
	}
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
	if cs.isCommitted {
		return true
	}
	return cs.hrs.Height() != height
}

func (cs *Consensus) invalidHeightRound(height int, round int) bool {
	if cs.isCommitted {
		return true
	}
	return cs.hrs.Height() != height || cs.hrs.Round() != round
}

func (cs *Consensus) invalidHeightRoundStep(height int, round int, step hrs.StepType) bool {
	if cs.isCommitted {
		return true
	}
	return cs.hrs.Height() != height || cs.hrs.Round() != round || cs.hrs.Step() > step
}

func (cs *Consensus) AddVote(v *vote.Vote) {
	cs.lk.Lock()
	defer cs.lk.Unlock()

	if err := cs.addVote(v); err != nil {
		cs.logger.Error("Error on adding a vote", "vote", v, "error", err)
	}
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
	if cs.hrs.Height() != v.Height() {
		return nil
	}

	added, err := cs.votes.AddVote(v)
	if err != nil {
		if v.Signer().EqualsTo(cs.signer.Address()) {
			cs.logger.Error("Detecting a duplicated vote from ourself. Did you restart the node?")
		}

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

	case vote.VoteTypePrecommit:
		precommits := cs.votes.Precommits(round)
		cs.logger.Debug("Vote added to precommit", "vote", v, "voteset", precommits)

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

	default:
		cs.logger.Panic("Unexpected vote type %X", v.VoteType)
	}

	return err
}

func (cs *Consensus) signAddVote(msgType vote.VoteType, hash crypto.Hash) {
	address := cs.signer.Address()
	if !cs.valset.Contains(address) {
		cs.logger.Trace("This node is not in validator set", "addr", address)
		return
	}

	// Sign the vote
	v := vote.NewVote(msgType, cs.hrs.Height(), cs.hrs.Round(), hash, address)
	cs.signer.SignMsg(v)
	cs.logger.Info("Our vote signed and broadcasted", "vote", v)

	err := cs.addVote(v)
	if err != nil {
		cs.logger.Error("Error on adding our vote!", "error", err, "vote", v)
		return
	}

	// Broadcast our vote
	msg := message.NewVoteMessage(v)
	cs.broadcastCh <- msg
}

func (cs *Consensus) requestForProposal() {
	msg := message.NewProposalReqMessage(cs.hrs.Height(), cs.hrs.Round())
	cs.broadcastCh <- msg
}
