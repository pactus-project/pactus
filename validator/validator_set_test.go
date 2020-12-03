package validator

import (
	"testing"

	"github.com/zarbchain/zarb-go/crypto"

	"github.com/stretchr/testify/assert"
)

func TestContains(t *testing.T) {
	vs, keys := GenerateTestValidatorSet()
	a, _, _ := crypto.GenerateTestKeyPair()

	assert.True(t, vs.Contains(keys[0].PublicKey().Address()))
	assert.True(t, vs.Contains(vs.Proposer(0).Address()))
	assert.False(t, vs.Contains(a))
}

func TestProposerMoves(t *testing.T) {
	vs, keys := GenerateTestValidatorSet()

	assert.Equal(t, vs.Proposer(0).Address(), keys[0].PublicKey().Address())
	assert.Equal(t, vs.Proposer(3).Address(), keys[3].PublicKey().Address())
	assert.Equal(t, vs.Proposer(4).Address(), keys[0].PublicKey().Address())

	vs.MoveToNewHeight(0)
	assert.Equal(t, vs.Proposer(0).Address(), keys[1].PublicKey().Address())
}

func TestInvalidProposerJoinAndLeave(t *testing.T) {
	val1, _ := GenerateTestValidator()
	val2, _ := GenerateTestValidator()
	val3, _ := GenerateTestValidator()
	val4, _ := GenerateTestValidator()
	val5, _ := GenerateTestValidator()

	vs, err := NewValidatorSet([]*Validator{val1, val2, val3, val4}, 4, val5.Address())
	assert.Error(t, err)
	assert.Nil(t, vs)
}

func TestProposerJoinAndLeave(t *testing.T) {
	val1, _ := GenerateTestValidator()
	val2, _ := GenerateTestValidator()
	val3, _ := GenerateTestValidator()
	val4, _ := GenerateTestValidator()
	val5, _ := GenerateTestValidator()
	val6, _ := GenerateTestValidator()
	val7, _ := GenerateTestValidator()

	vs, err := NewValidatorSet([]*Validator{val1, val2, val3, val4}, 4, val1.Address())
	assert.NoError(t, err)

	// Val1 is already in set
	err = vs.Join(val1)
	assert.Error(t, err)

	// -------------------
	// Validator 1 is first proposer
	assert.Equal(t, vs.validators[0].Address(), val1.Address())
	assert.Equal(t, vs.validators[1].Address(), val2.Address())
	assert.Equal(t, vs.validators[2].Address(), val3.Address())
	assert.Equal(t, vs.validators[3].Address(), val4.Address())

	assert.Equal(t, vs.Proposer(0).Address(), val1.Address())
	vs.MoveToNewHeight(0)

	// -------------------
	// Validator 2 is offline, proposer will be validator 3 for the second round
	assert.Equal(t, vs.validators[0].Address(), val1.Address())
	assert.Equal(t, vs.validators[1].Address(), val2.Address())
	assert.Equal(t, vs.validators[2].Address(), val3.Address())
	assert.Equal(t, vs.validators[3].Address(), val4.Address())

	assert.Equal(t, vs.Proposer(0).Address(), val2.Address())
	assert.Equal(t, vs.Proposer(1).Address(), val3.Address())
	vs.MoveToNewHeight(1)

	// -------------------
	// Next psoposer is validator 4
	// In this height, validator 5 joins the set and validator 1 should leave the set
	assert.Equal(t, vs.validators[0].Address(), val1.Address())
	assert.Equal(t, vs.validators[1].Address(), val2.Address())
	assert.Equal(t, vs.validators[2].Address(), val3.Address())
	assert.Equal(t, vs.validators[3].Address(), val4.Address())

	assert.Equal(t, vs.Proposer(0).Address(), val4.Address())
	err = vs.Join(val5)
	assert.NoError(t, err)

	// Can't enter for second time
	err = vs.Join(val6)
	assert.Error(t, err)

	vs.MoveToNewHeight(0)

	// -------------------------
	// Now validator 1 has left the set and validator 5 is proposer
	assert.Equal(t, vs.validators[0].Address(), val2.Address())
	assert.Equal(t, vs.validators[1].Address(), val3.Address())
	assert.Equal(t, vs.validators[2].Address(), val4.Address())
	assert.Equal(t, vs.validators[3].Address(), val5.Address())

	assert.Equal(t, vs.Proposer(0).Address(), val5.Address())
	assert.False(t, vs.Contains(val1.Address()))
	vs.MoveToNewHeight(0)

	// -------------------------
	// Now validator 2 is the proposer but he is offline
	// At this time validator 6 also joins the set
	assert.Equal(t, vs.validators[0].Address(), val2.Address())
	assert.Equal(t, vs.validators[1].Address(), val3.Address())
	assert.Equal(t, vs.validators[2].Address(), val4.Address())
	assert.Equal(t, vs.validators[3].Address(), val5.Address())

	assert.Equal(t, vs.Proposer(0).Address(), val2.Address())
	assert.Equal(t, vs.Proposer(1).Address(), val3.Address())
	vs.Join(val6)
	vs.MoveToNewHeight(1)

	// -------------------------
	// Now validator 4 is proposer now
	assert.Equal(t, vs.validators[0].Address(), val3.Address())
	assert.Equal(t, vs.validators[1].Address(), val4.Address())
	assert.Equal(t, vs.validators[2].Address(), val5.Address())
	assert.Equal(t, vs.validators[3].Address(), val6.Address())

	assert.Equal(t, vs.Proposer(0).Address(), val4.Address())
	vs.MoveToNewHeight(0)

	// -------------------------
	// Now validator 5 is proposer now
	// Validator 7 joins the set and validator 3 should leave
	assert.Equal(t, vs.validators[0].Address(), val3.Address())
	assert.Equal(t, vs.validators[1].Address(), val4.Address())
	assert.Equal(t, vs.validators[2].Address(), val5.Address())
	assert.Equal(t, vs.validators[3].Address(), val6.Address())

	assert.Equal(t, vs.Proposer(0).Address(), val5.Address())
	vs.Join(val7)
	vs.MoveToNewHeight(0)

	// -------------------------
	// Now validator 3 has left the set and validator 6 is proposer
	assert.Equal(t, vs.validators[0].Address(), val4.Address())
	assert.Equal(t, vs.validators[1].Address(), val5.Address())
	assert.Equal(t, vs.validators[2].Address(), val6.Address())
	assert.Equal(t, vs.validators[3].Address(), val7.Address())

	assert.Equal(t, vs.Proposer(0).Address(), val6.Address())
	assert.False(t, vs.Contains(val3.Address()))
	vs.MoveToNewHeight(0)

	// -------------------------
	// Validator 7 is proposer (validators in set : 4,5,6,7)
	assert.Equal(t, vs.validators[0].Address(), val4.Address())
	assert.Equal(t, vs.validators[1].Address(), val5.Address())
	assert.Equal(t, vs.validators[2].Address(), val6.Address())
	assert.Equal(t, vs.validators[3].Address(), val7.Address())

	assert.Equal(t, vs.Proposer(0).Address(), val7.Address())

	vs.MoveToNewHeight(0)

	// -------------------------
	// Validator 4 is proposer
	// Validator 1 joins again
	assert.Equal(t, vs.validators[0].Address(), val4.Address())
	assert.Equal(t, vs.validators[1].Address(), val5.Address())
	assert.Equal(t, vs.validators[2].Address(), val6.Address())
	assert.Equal(t, vs.validators[3].Address(), val7.Address())

	assert.Equal(t, vs.Proposer(0).Address(), val4.Address())

	vs.Join(val1)
	vs.MoveToNewHeight(0)

	// -------------------------
	// validator 1 joined, and 4 left
	assert.Equal(t, vs.validators[0].Address(), val5.Address())
	assert.Equal(t, vs.validators[1].Address(), val6.Address())
	assert.Equal(t, vs.validators[2].Address(), val7.Address())
	assert.Equal(t, vs.validators[3].Address(), val1.Address())

	assert.Equal(t, vs.Proposer(0).Address(), val5.Address())

	vs.Join(val3)
	vs.MoveToNewHeight(0)

	// -------------------------
	// Validator 3 joins again
	assert.Equal(t, vs.validators[0].Address(), val6.Address())
	assert.Equal(t, vs.validators[1].Address(), val7.Address())
	assert.Equal(t, vs.validators[2].Address(), val1.Address())
	assert.Equal(t, vs.validators[3].Address(), val3.Address())

	assert.Equal(t, vs.Proposer(0).Address(), val6.Address())
}
