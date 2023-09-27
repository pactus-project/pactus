package consensus

import (
	"fmt"
	"testing"
	"time"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/genesis"
	"github.com/pactus-project/pactus/state"
	"github.com/pactus-project/pactus/store"
	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/pactus-project/pactus/txpool"
	"github.com/pactus-project/pactus/types/account"
	"github.com/pactus-project/pactus/types/param"
	"github.com/pactus-project/pactus/types/proposal"
	"github.com/pactus-project/pactus/types/validator"
	"github.com/pactus-project/pactus/types/vote"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/logger"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func shouldPublishProposal(t *testing.T, broadcastCh chan message.Message,
	height uint32, round int16,
) *proposal.Proposal {
	t.Helper()

	timeout := time.NewTimer(1 * time.Second)

	for {
		select {
		case <-timeout.C:
			require.NoError(t, fmt.Errorf("Timeout"))
			return nil
		case msg := <-broadcastCh:
			logger.Info("shouldPublishProposal", "message", msg)

			if msg.Type() == message.TypeProposal {
				m := msg.(*message.ProposalMessage)
				require.Equal(t, m.Proposal.Height(), height)
				require.Equal(t, m.Proposal.Round(), round)
				return m.Proposal
			}
		}
	}
}

func TestManager(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	_, committeeValKeys := ts.GenerateTestCommittee(5)
	acc := account.NewAccount(0)
	acc.AddToBalance(21 * 1e14)
	params := param.DefaultParams()
	params.BlockIntervalInSecond = 1
	vals := make([]*validator.Validator, 5)
	for i, s := range committeeValKeys {
		val := validator.NewValidator(s.PublicKey(), int32(i))
		vals[i] = val
	}
	accs := map[crypto.Address]*account.Account{crypto.TreasuryAddress: acc}
	// to prevent triggering timers before starting the tests to avoid double entries for new heights in some tests.
	getTime := util.RoundNow(params.BlockIntervalInSecond).Add(time.Duration(params.BlockIntervalInSecond) * time.Second)
	genDoc := genesis.MakeGenesis(getTime, accs, vals, params)

	rewardAddrs := []crypto.Address{
		ts.RandAccAddress(), ts.RandAccAddress(),
		ts.RandAccAddress(), ts.RandAccAddress(),
		ts.RandAccAddress(),
	}
	valKeys := make([]*bls.ValidatorKey, 5)
	valKeys[0] = committeeValKeys[0]
	valKeys[1] = ts.RandValKey()
	valKeys[2] = committeeValKeys[1]
	valKeys[3] = ts.RandValKey()
	valKeys[4] = ts.RandValKey()
	broadcastCh := make(chan message.Message, 500)
	txPool := txpool.MockingTxPool()

	state, err := state.LoadOrNewState(genDoc, valKeys, store.MockingStore(ts), txPool, nil)
	require.NoError(t, err)

	Mgr := NewManager(testConfig(), state, valKeys, rewardAddrs, broadcastCh)
	mgr := Mgr.(*manager)

	consA := mgr.instances[0].(*consensus) // active
	consB := mgr.instances[1].(*consensus) // inactive
	consC := mgr.instances[2].(*consensus) // active
	consD := mgr.instances[3].(*consensus) // inactive
	consE := mgr.instances[4].(*consensus) // inactive

	assert.False(t, mgr.HasActiveInstance())
	mgr.MoveToNewHeight()
	newHeightTimeout(consA)
	newHeightTimeout(consB)
	newHeightTimeout(consC)
	newHeightTimeout(consD)
	newHeightTimeout(consE)

	t.Run("Check if keys are assigned properly", func(t *testing.T) {
		assert.Equal(t, valKeys[0].PublicKey(), consA.ConsensusKey())
		assert.Equal(t, valKeys[1].PublicKey(), consB.ConsensusKey())
		assert.Equal(t, valKeys[2].PublicKey(), consC.ConsensusKey())
		assert.Equal(t, valKeys[3].PublicKey(), consD.ConsensusKey())
		assert.Equal(t, valKeys[4].PublicKey(), consE.ConsensusKey())
	})

	t.Run("Check if all instances move to new height", func(t *testing.T) {
		assert.True(t, mgr.HasActiveInstance())
	})

	t.Run("Check if all instances move to new height", func(t *testing.T) {
		h, r := mgr.HeightRound()
		assert.Equal(t, h, uint32(1))
		assert.Equal(t, r, int16(0))
		assert.True(t, mgr.HasActiveInstance())
	})

	t.Run("Testing add vote", func(t *testing.T) {
		v := vote.NewPrepareVote(ts.RandHash(), 1, 0, committeeValKeys[2].Address())
		ts.HelperSignVote(committeeValKeys[2], v)

		mgr.AddVote(v)

		assert.True(t, consA.HasVote(v.Hash()))
		assert.False(t, consB.HasVote(v.Hash()))
		assert.True(t, consC.HasVote(v.Hash()))
		assert.False(t, consD.HasVote(v.Hash()))
	})

	t.Run("Testing set proposal", func(t *testing.T) {
		b, _ := state.ProposeBlock(committeeValKeys[1], committeeValKeys[1].Address(), 1)
		p := proposal.NewProposal(1, 1, b)
		ts.HelperSignProposal(committeeValKeys[1], p)

		mgr.SetProposal(p)

		assert.Equal(t, consA.RoundProposal(1), p)
		assert.Nil(t, consB.RoundProposal(1))
		assert.Equal(t, consC.RoundProposal(1), p)
		assert.Nil(t, consD.RoundProposal(1))
	})

	t.Run("Check if one instance publishes a proposal, the other instances receive it", func(t *testing.T) {
		p := shouldPublishProposal(t, broadcastCh, 1, 0)

		assert.Equal(t, mgr.RoundProposal(0), p)
		assert.Equal(t, consA.RoundProposal(0), p)
		assert.Nil(t, consB.RoundProposal(0))
	})
}
