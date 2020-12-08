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

func TestProposerMove(t *testing.T) {
	val1, _ := GenerateTestValidator()
	val2, _ := GenerateTestValidator()
	val3, _ := GenerateTestValidator()
	val4, _ := GenerateTestValidator()
	val5, _ := GenerateTestValidator()
	val6, _ := GenerateTestValidator()
	val7, _ := GenerateTestValidator()

	vs, err := NewValidatorSet([]*Validator{val1, val2, val3, val4, val5, val6, val7}, 7, val1.Address())
	assert.NoError(t, err)

	//
	// +=+-+-+-+-+-+-+       +-+=+-+-+-+-+-+
	// |1|2|3|4|5|6|7|  ==>  |1|2|3|4|5|6|7|
	// +=+-+-+-+-+-+-+       +-+=+-+-+-+-+-+
	//
	vs.proposerIndex = 0
	assert.Equal(t, vs.Proposer(0).Address(), val1.Address())
	vs.MoveToNewHeight(0)
	assert.Equal(t, vs.proposerIndex, 1)
	assert.Equal(t, vs.Proposer(0).Address(), val2.Address())
	assert.Equal(t, vs.Proposer(1).Address(), val3.Address())

	//
	// +-+-+-+=+-+-+-+       +-+-+-+-+=+-+-+
	// |1|2|3|4|5|6|7|  ==>  |1|2|3|4|5|6|7|
	// +-+-+-+=+-+-+-+       +-+-+-+-+=+-+-+
	//
	vs.proposerIndex = 3
	assert.Equal(t, vs.Proposer(0).Address(), val4.Address())
	vs.MoveToNewHeight(0)
	assert.Equal(t, vs.proposerIndex, 4)
	assert.Equal(t, vs.Proposer(0).Address(), val5.Address())

	//
	// +-+-+-+-+-+-+=+       +=+-+-+-+-+-+-+
	// |1|2|3|4|5|6|7|  ==>  |1|2|3|4|5|6|7|
	// +-+-+-+-+-+-+=+       +=+-+-+-+-+-+-+
	//
	vs.proposerIndex = 6
	assert.Equal(t, vs.Proposer(0).Address(), val7.Address())
	vs.MoveToNewHeight(0)
	assert.Equal(t, vs.proposerIndex, 0)
	assert.Equal(t, vs.Proposer(0).Address(), val1.Address())
}

func TestProposerMoveMoreRounds(t *testing.T) {
	val1, _ := GenerateTestValidator()
	val2, _ := GenerateTestValidator()
	val3, _ := GenerateTestValidator()
	val4, _ := GenerateTestValidator()
	val5, _ := GenerateTestValidator()
	val6, _ := GenerateTestValidator()
	val7, _ := GenerateTestValidator()

	vs, err := NewValidatorSet([]*Validator{val1, val2, val3, val4, val5, val6, val7}, 7, val1.Address())
	assert.NoError(t, err)

	//
	// +=+-+-+-+-+-+-+       +-+-+-+=+-+-+-+
	// |1|2|3|4|5|6|7|  ==>  |1|2|3|4|5|6|7|
	// +=+-+-+-+-+-+-+       +-+-+-+=+-+-+-+
	//
	vs.proposerIndex = 0
	assert.Equal(t, vs.Proposer(0).Address(), val1.Address())
	vs.MoveToNewHeight(2)
	assert.Equal(t, vs.proposerIndex, 3)
	assert.Equal(t, vs.Proposer(0).Address(), val4.Address())
	assert.Equal(t, vs.Proposer(1).Address(), val5.Address())

	//
	// +-+-+-+=+-+-+-+       +=+-+-+-+-+-+-+
	// |1|2|3|4|5|6|7|  ==>  |1|2|3|4|5|6|7|
	// +-+-+-+=+-+-+-+       +=+-+-+-+-+-+-+
	//
	vs.proposerIndex = 3
	assert.Equal(t, vs.Proposer(0).Address(), val4.Address())
	vs.MoveToNewHeight(3)
	assert.Equal(t, vs.proposerIndex, 0)
	assert.Equal(t, vs.Proposer(0).Address(), val1.Address())

	//
	// +-+-+-+-+-+-+=+       +-+=+-+-+-+-+-+
	// |1|2|3|4|5|6|7|  ==>  |1|2|3|4|5|6|7|
	// +-+-+-+-+-+-+=+       +-+=+-+-+-+-+-+
	//
	vs.proposerIndex = 6
	assert.Equal(t, vs.Proposer(0).Address(), val7.Address())
	vs.MoveToNewHeight(1)
	assert.Equal(t, vs.proposerIndex, 1)
	assert.Equal(t, vs.Proposer(0).Address(), val2.Address())
}

func TestProposerJoinAndLeave(t *testing.T) {
	val1, _ := GenerateTestValidator()
	val2, _ := GenerateTestValidator()
	val3, _ := GenerateTestValidator()
	val4, _ := GenerateTestValidator()
	val5, _ := GenerateTestValidator()
	val6, _ := GenerateTestValidator()
	val7, _ := GenerateTestValidator()
	val8, _ := GenerateTestValidator()
	val9, _ := GenerateTestValidator()
	valA, _ := GenerateTestValidator()
	valB, _ := GenerateTestValidator()
	valC, _ := GenerateTestValidator()
	valD, _ := GenerateTestValidator()

	vs, err := NewValidatorSet([]*Validator{val1, val2, val3, val4, val5, val6, val7}, 7, val1.Address())
	assert.NoError(t, err)

	// Val1 is already in set
	err = vs.Join(val1)
	assert.Error(t, err)

	//
	// +=+-+-+-+-+-+-+       +=+-+-+-+-+-+-+
	// |1|2|3|4|5|6|7|  ==>  |2|3|4|5|6|7|8|
	// +=+-+-+-+-+-+-+       +=+-+-+-+-+-+-+
	//
	vs.proposerIndex = 0
	assert.NoError(t, vs.Join(val8))
	vs.MoveToNewHeight(0)
	assert.Equal(t, vs.proposerIndex, 0)
	assert.Equal(t, vs.Proposer(0).Address(), val2.Address())
	assert.Equal(t, vs.Proposer(1).Address(), val3.Address())

	//
	// +-+-+=+-+-+-+-+       +-+=+-+-+-+-+-+
	// |2|3|4|5|6|7|8|  ==>  |4|5|6|7|8|9|A|
	// +-+-+=+-+-+-+-+       +-+=+-+-+-+-+-+
	//
	//
	vs.proposerIndex = 2
	assert.NoError(t, vs.Join(val9))
	assert.NoError(t, vs.Join(valA))
	vs.MoveToNewHeight(0)
	assert.Equal(t, vs.proposerIndex, 1)
	assert.Equal(t, vs.Proposer(0).Address(), val5.Address())

	//
	// +-+-+-+-+-+-+=+       +=+-+-+-+-+-+-+
	// |4|5|6|7|8|9|A|  ==>  |5|6|7|8|9|A|B|
	// +-+-+-+-+-+-+=+       +=+-+-+-+-+-+-+
	//
	vs.proposerIndex = 6
	assert.NoError(t, vs.Join(valB))
	vs.MoveToNewHeight(0)
	assert.Equal(t, vs.proposerIndex, 0)
	assert.Equal(t, vs.Proposer(0).Address(), val5.Address())

	//
	// +-+-+-+-+-+-+=+       +=+-+-+-+-+-+-+
	// |5|6|7|8|9|A|B|  ==>  |7|8|9|A|B|C|D|
	// +-+-+-+-+-+-+=+       +=+-+-+-+-+-+-+
	//
	vs.proposerIndex = 6
	assert.NoError(t, vs.Join(valC))
	assert.NoError(t, vs.Join(valD))
	vs.MoveToNewHeight(0)
	assert.Equal(t, vs.proposerIndex, 0)
	assert.Equal(t, vs.Proposer(0).Address(), val7.Address())
}

func TestProposerJoinAndLeaveMoreRound(t *testing.T) {
	val1, _ := GenerateTestValidator()
	val2, _ := GenerateTestValidator()
	val3, _ := GenerateTestValidator()
	val4, _ := GenerateTestValidator()
	val5, _ := GenerateTestValidator()
	val6, _ := GenerateTestValidator()
	val7, _ := GenerateTestValidator()
	val8, _ := GenerateTestValidator()
	val9, _ := GenerateTestValidator()
	valA, _ := GenerateTestValidator()
	valB, _ := GenerateTestValidator()
	valC, _ := GenerateTestValidator()
	valD, _ := GenerateTestValidator()

	vs, err := NewValidatorSet([]*Validator{val1, val2, val3, val4, val5, val6, val7}, 7, val1.Address())
	assert.NoError(t, err)

	// Val1 is already in set
	err = vs.Join(val1)
	assert.Error(t, err)

	//
	// +=+-+-+-+-+-+-+       +-+-+=+-+-+-+-+
	// |1|2|3|4|5|6|7|  ==>  |2|3|4|5|6|7|8|
	// +=+-+-+-+-+-+-+       +-+-+=+-+-+-+-+
	//
	vs.proposerIndex = 0
	assert.NoError(t, vs.Join(val8))
	vs.MoveToNewHeight(2)
	assert.Equal(t, vs.proposerIndex, 2)
	assert.Equal(t, vs.Proposer(0).Address(), val4.Address())
	assert.Equal(t, vs.Proposer(1).Address(), val5.Address())

	//
	// +-+-+=+-+-+-+-+       +-+-+-+-+=+-+-+
	// |2|3|4|5|6|7|8|  ==>  |4|5|6|7|8|9|A|
	// +-+-+=+-+-+-+-+       +-+-+-+-+=+-+-+
	//
	//
	vs.proposerIndex = 2
	assert.NoError(t, vs.Join(val9))
	assert.NoError(t, vs.Join(valA))
	vs.MoveToNewHeight(3)
	assert.Equal(t, vs.proposerIndex, 4)
	assert.Equal(t, vs.Proposer(0).Address(), val8.Address())

	//
	// +-+-+-+-+-+-+=+       +-+=+-+-+-+-+-+
	// |4|5|6|7|8|9|A|  ==>  |5|6|7|8|9|A|B|
	// +-+-+-+-+-+-+=+       +-+=+-+-+-+-+-+
	//
	// 5 is offline
	vs.proposerIndex = 6
	assert.NoError(t, vs.Join(valB))
	vs.MoveToNewHeight(2)
	assert.Equal(t, vs.proposerIndex, 1)
	assert.Equal(t, vs.Proposer(0).Address(), val6.Address())

	//
	// +-+-+-+-+-+-+=+       +-+-+-+-+-+-+-+
	// |5|6|7|8|9|A|B|  ==>  |7|8|9|A|B|C|D|
	// +-+-+-+-+-+-+=+       +-+-+-+-+-+-+-+
	//
	vs.proposerIndex = 5
	assert.NoError(t, vs.Join(valC))
	assert.NoError(t, vs.Join(valD))
	vs.MoveToNewHeight(2)
	assert.Equal(t, vs.proposerIndex, 0)
	assert.Equal(t, vs.Proposer(0).Address(), val7.Address())
}
