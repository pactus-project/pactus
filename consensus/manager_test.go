package consensus

import (
	"testing"
	"time"

	"github.com/pactus-project/pactus/committee"
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/crypto/hash"
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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestManager(t *testing.T) {
	_, committeeSigners := committee.GenerateTestCommittee(5)
	acc := account.NewAccount(0)
	acc.AddToBalance(21 * 1e14)
	params := param.DefaultParams()
	params.BlockTimeInSecond = 1
	vals := make([]*validator.Validator, 5)
	for i, s := range committeeSigners {
		val := validator.NewValidator(s.PublicKey().(*bls.PublicKey), int32(i))
		vals[i] = val
	}
	accs := map[crypto.Address]*account.Account{crypto.TreasuryAddress: acc}
	// to prevent triggering timers before starting the tests to avoid double entries for new heights in some tests.
	getTime := util.RoundNow(params.BlockTimeInSecond).Add(time.Duration(params.BlockTimeInSecond) * time.Second)
	genDoc := genesis.MakeGenesis(getTime, accs, vals, params)

	rewardAddrs := []crypto.Address{
		crypto.GenerateTestAddress(), crypto.GenerateTestAddress(),
		crypto.GenerateTestAddress(), crypto.GenerateTestAddress(),
		crypto.GenerateTestAddress(),
	}
	signers := make([]crypto.Signer, 5)
	signers[0] = committeeSigners[0]
	signers[1] = committeeSigners[1]
	signers[2] = bls.GenerateTestSigner()
	signers[3] = bls.GenerateTestSigner()
	signers[4] = bls.GenerateTestSigner()
	broadcastCh := make(chan message.Message, 500)
	txPool := txpool.MockingTxPool()

	state, err := state.LoadOrNewState(genDoc, signers, store.MockingStore(), txPool, nil)
	require.NoError(t, err)

	Mgr := NewManager(testConfig(), state, signers, rewardAddrs, broadcastCh)
	mgr := Mgr.(*manager)

	consA := mgr.instances[0].(*consensus) // active
	consB := mgr.instances[1].(*consensus) // active
	consC := mgr.instances[2].(*consensus) // inactive
	consD := mgr.instances[3].(*consensus) // inactive
	consE := mgr.instances[4].(*consensus) // inactive

	assert.False(t, mgr.HasActiveInstance())

	t.Run("Check if keys are assigned properly", func(t *testing.T) {
		instances := mgr.Instances()

		assert.Equal(t, signers[0].PublicKey(), instances[0].SignerKey())
		assert.Equal(t, signers[1].PublicKey(), instances[1].SignerKey())
		assert.Equal(t, signers[2].PublicKey(), instances[2].SignerKey())
		assert.Equal(t, signers[3].PublicKey(), instances[3].SignerKey())
		assert.Equal(t, signers[4].PublicKey(), instances[4].SignerKey())
	})

	t.Run("Check if all instances move to new height", func(t *testing.T) {
		mgr.MoveToNewHeight()

		checkHeightRoundWait(t, consA, 1, 0)
		checkHeightRoundWait(t, consB, 1, 0)
		checkHeightRoundWait(t, consC, 1, 0)
		checkHeightRoundWait(t, consD, 1, 0)
		checkHeightRoundWait(t, consE, 1, 0)

		assert.True(t, mgr.HasActiveInstance())
	})

	t.Run("Testing add vote", func(t *testing.T) {
		v := vote.NewVote(vote.VoteTypeChangeProposer, 1, 0, hash.UndefHash, committeeSigners[2].Address())
		committeeSigners[2].SignMsg(v)

		mgr.AddVote(v)

		assert.True(t, consA.HasVote(v.Hash()))
		assert.True(t, consB.HasVote(v.Hash()))
		assert.False(t, consC.HasVote(v.Hash()))
		assert.False(t, consD.HasVote(v.Hash()))
	})

	t.Run("Testing set proposal", func(t *testing.T) {
		b, _ := state.ProposeBlock(committeeSigners[2], committeeSigners[2].Address(), 2)
		p := proposal.NewProposal(1, 2, b)
		committeeSigners[2].SignMsg(p)

		mgr.SetProposal(p)

		assert.Equal(t, consA.RoundProposal(2), p)
		assert.Equal(t, consB.RoundProposal(2), p)
		assert.Nil(t, consC.RoundProposal(2))
		assert.Nil(t, consD.RoundProposal(2))
	})

	t.Run("Check if one instance publishes a proposal, the other instances receive it", func(t *testing.T) {
		p := shouldPublishProposal(t, consA, 1, 0)

		assert.Equal(t, consA.RoundProposal(0), p)
		assert.Equal(t, consB.RoundProposal(0), p)
		assert.Nil(t, consC.RoundProposal(0))
		assert.Nil(t, consD.RoundProposal(0))
	})

	t.Run("Testing moving to the next round proposal", func(t *testing.T) {
		v3 := vote.NewVote(vote.VoteTypeChangeProposer, 1, 0, hash.UndefHash, committeeSigners[2].Address())
		committeeSigners[2].SignMsg(v3)

		v4 := vote.NewVote(vote.VoteTypeChangeProposer, 1, 0, hash.UndefHash, committeeSigners[3].Address())
		committeeSigners[3].SignMsg(v4)

		mgr.AddVote(v3)
		mgr.AddVote(v4)

		checkHeightRoundWait(t, consA, 1, 1)

		h, r := mgr.HeightRound()
		assert.Equal(t, h, uint32(1))
		assert.Equal(t, r, int16(1))
		assert.Nil(t, mgr.RoundProposal(0))
	})
}
