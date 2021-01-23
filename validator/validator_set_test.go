package validator

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/crypto"
)

func TestContains(t *testing.T) {
	vs, signers := GenerateTestValidatorSet()

	assert.True(t, vs.Contains(signers[0].Address()))
	assert.True(t, vs.Contains(vs.Proposer(0).Address()))
}

func TestProposer(t *testing.T) {
	vs, signers := GenerateTestValidatorSet()

	assert.Equal(t, vs.Proposer(0).Address(), signers[0].Address())
	assert.Equal(t, vs.Proposer(3).Address(), signers[3].Address())
	assert.Equal(t, vs.Proposer(4).Address(), signers[0].Address())

	assert.NoError(t, vs.UpdateTheSet(0, nil))
	assert.Equal(t, vs.Proposer(0).Address(), signers[1].Address())
}

func TestInvalidProposerJoinAndLeave(t *testing.T) {
	val1, _ := GenerateTestValidator(0)
	val2, _ := GenerateTestValidator(1)
	val3, _ := GenerateTestValidator(2)
	val4, _ := GenerateTestValidator(3)
	val5, _ := GenerateTestValidator(4)

	vs, err := NewValidatorSet([]*Validator{val1, val2, val3, val4}, 4, val5.Address())
	assert.Error(t, err)
	assert.Nil(t, vs)
}

func TestProposerMove(t *testing.T) {
	val1, _ := GenerateTestValidator(0)
	val2, _ := GenerateTestValidator(1)
	val3, _ := GenerateTestValidator(2)
	val4, _ := GenerateTestValidator(3)
	val5, _ := GenerateTestValidator(4)
	val6, _ := GenerateTestValidator(5)
	val7, _ := GenerateTestValidator(6)

	vs, err := NewValidatorSet([]*Validator{val1, val2, val3, val4, val5, val6, val7}, 7, val1.Address())
	assert.NoError(t, err)

	//
	// +=+-+-+-+-+-+-+       +-+=+-+-+-+-+-+
	// |1|2|3|4|5|6|7|  ==>  |1|2|3|4|5|6|7|
	// +=+-+-+-+-+-+-+       +-+=+-+-+-+-+-+
	//
	vs.proposerIndex = 0
	assert.Equal(t, vs.Proposer(0).Address(), val1.Address())
	assert.NoError(t, vs.UpdateTheSet(0, nil))
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
	assert.NoError(t, vs.UpdateTheSet(0, nil))
	assert.Equal(t, vs.proposerIndex, 4)
	assert.Equal(t, vs.Proposer(0).Address(), val5.Address())

	//
	// +-+-+-+-+-+-+=+       +=+-+-+-+-+-+-+
	// |1|2|3|4|5|6|7|  ==>  |1|2|3|4|5|6|7|
	// +-+-+-+-+-+-+=+       +=+-+-+-+-+-+-+
	//
	vs.proposerIndex = 6
	assert.Equal(t, vs.Proposer(0).Address(), val7.Address())
	assert.NoError(t, vs.UpdateTheSet(0, nil))
	assert.Equal(t, vs.proposerIndex, 0)
	assert.Equal(t, vs.Proposer(0).Address(), val1.Address())
}

func TestProposerMoveMoreRounds(t *testing.T) {
	val1, _ := GenerateTestValidator(0)
	val2, _ := GenerateTestValidator(1)
	val3, _ := GenerateTestValidator(2)
	val4, _ := GenerateTestValidator(3)
	val5, _ := GenerateTestValidator(4)
	val6, _ := GenerateTestValidator(5)
	val7, _ := GenerateTestValidator(6)

	vs, err := NewValidatorSet([]*Validator{val1, val2, val3, val4, val5, val6, val7}, 7, val1.Address())
	assert.NoError(t, err)

	//
	// +=+-+-+-+-+-+-+       +-+-+-+=+-+-+-+
	// |1|2|3|4|5|6|7|  ==>  |1|2|3|4|5|6|7|
	// +=+-+-+-+-+-+-+       +-+-+-+=+-+-+-+
	//
	vs.proposerIndex = 0
	assert.Equal(t, vs.Proposer(0).Address(), val1.Address())
	assert.NoError(t, vs.UpdateTheSet(2, nil))
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
	assert.NoError(t, vs.UpdateTheSet(3, nil))
	assert.Equal(t, vs.proposerIndex, 0)
	assert.Equal(t, vs.Proposer(0).Address(), val1.Address())

	//
	// +-+-+-+-+-+-+=+       +-+=+-+-+-+-+-+
	// |1|2|3|4|5|6|7|  ==>  |1|2|3|4|5|6|7|
	// +-+-+-+-+-+-+=+       +-+=+-+-+-+-+-+
	//
	vs.proposerIndex = 6
	assert.Equal(t, vs.Proposer(0).Address(), val7.Address())
	assert.NoError(t, vs.UpdateTheSet(1, nil))
	assert.Equal(t, vs.proposerIndex, 1)
	assert.Equal(t, vs.Proposer(0).Address(), val2.Address())
}

func TestProposerJoinAndLeave(t *testing.T) {
	val1, _ := GenerateTestValidator(0)
	val2, _ := GenerateTestValidator(1)
	val3, _ := GenerateTestValidator(2)
	val4, _ := GenerateTestValidator(3)
	val5, _ := GenerateTestValidator(4)
	val6, _ := GenerateTestValidator(5)
	val7, _ := GenerateTestValidator(6)
	val8, _ := GenerateTestValidator(7)
	val9, _ := GenerateTestValidator(8)
	valA, _ := GenerateTestValidator(9)
	valB, _ := GenerateTestValidator(10)
	valC, _ := GenerateTestValidator(11)
	valD, _ := GenerateTestValidator(12)

	vs, err := NewValidatorSet([]*Validator{val1, val2, val3, val4, val5, val6, val7}, 7, val1.Address())
	assert.NoError(t, err)

	// Val1 is already in set
	assert.Error(t, vs.UpdateTheSet(0, []*Validator{val1}))

	//
	// +=+-+-+-+-+-+-+       +=+-+-+-+-+-+-+
	// |1|2|3|4|5|6|7|  ==>  |2|3|4|5|6|7|8|
	// +=+-+-+-+-+-+-+       +=+-+-+-+-+-+-+
	//
	vs.proposerIndex = 0
	assert.NoError(t, vs.UpdateTheSet(0, []*Validator{val8}))
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
	assert.NoError(t, vs.UpdateTheSet(0, []*Validator{val9, valA}))
	assert.Equal(t, vs.proposerIndex, 1)
	assert.Equal(t, vs.Proposer(0).Address(), val5.Address())

	//
	// +-+-+-+-+-+-+=+       +=+-+-+-+-+-+-+
	// |4|5|6|7|8|9|A|  ==>  |5|6|7|8|9|A|B|
	// +-+-+-+-+-+-+=+       +=+-+-+-+-+-+-+
	//
	vs.proposerIndex = 6
	assert.NoError(t, vs.UpdateTheSet(0, []*Validator{valB}))
	assert.Equal(t, vs.proposerIndex, 0)
	assert.Equal(t, vs.Proposer(0).Address(), val5.Address())

	//
	// +-+-+-+-+-+-+=+       +=+-+-+-+-+-+-+
	// |5|6|7|8|9|A|B|  ==>  |7|8|9|A|B|C|D|
	// +-+-+-+-+-+-+=+       +=+-+-+-+-+-+-+
	//
	vs.proposerIndex = 6
	assert.NoError(t, vs.UpdateTheSet(0, []*Validator{valC, valD}))
	assert.Equal(t, vs.proposerIndex, 0)
	assert.Equal(t, vs.Proposer(0).Address(), val7.Address())
}

func TestProposerJoinAndLeaveMoreRound(t *testing.T) {
	val1, _ := GenerateTestValidator(0)
	val2, _ := GenerateTestValidator(1)
	val3, _ := GenerateTestValidator(2)
	val4, _ := GenerateTestValidator(3)
	val5, _ := GenerateTestValidator(4)
	val6, _ := GenerateTestValidator(5)
	val7, _ := GenerateTestValidator(6)
	val8, _ := GenerateTestValidator(7)
	val9, _ := GenerateTestValidator(8)
	valA, _ := GenerateTestValidator(9)
	valB, _ := GenerateTestValidator(10)
	valC, _ := GenerateTestValidator(11)
	valD, _ := GenerateTestValidator(12)

	vs, err := NewValidatorSet([]*Validator{val1, val2, val3, val4, val5, val6, val7}, 7, val1.Address())
	assert.NoError(t, err)

	// Val1 is already in set
	assert.Error(t, vs.UpdateTheSet(0, []*Validator{val1}))

	//
	// +=+-+-+-+-+-+-+       +-+-+=+-+-+-+-+
	// |1|2|3|4|5|6|7|  ==>  |2|3|4|5|6|7|8|
	// +=+-+-+-+-+-+-+       +-+-+=+-+-+-+-+
	//
	vs.proposerIndex = 0
	assert.NoError(t, vs.UpdateTheSet(2, []*Validator{val8}))
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
	assert.NoError(t, vs.UpdateTheSet(3, []*Validator{val9, valA}))
	assert.Equal(t, vs.proposerIndex, 4)
	assert.Equal(t, vs.Proposer(0).Address(), val8.Address())

	//
	// +-+-+-+-+-+-+=+       +-+=+-+-+-+-+-+
	// |4|5|6|7|8|9|A|  ==>  |5|6|7|8|9|A|B|
	// +-+-+-+-+-+-+=+       +-+=+-+-+-+-+-+
	//
	// 5 is offline
	vs.proposerIndex = 6
	assert.NoError(t, vs.UpdateTheSet(2, []*Validator{valB}))

	assert.Equal(t, vs.proposerIndex, 1)
	assert.Equal(t, vs.Proposer(0).Address(), val6.Address())

	//
	// +-+-+-+-+-+-+=+       +=+-+-+-+-+-+-+
	// |5|6|7|8|9|A|B|  ==>  |7|8|9|A|B|C|D|
	// +-+-+-+-+-+-+=+       +=+-+-+-+-+-+-+
	//
	vs.proposerIndex = 5
	assert.NoError(t, vs.UpdateTheSet(2, []*Validator{valC, valD}))

	assert.Equal(t, vs.proposerIndex, 0)
	assert.Equal(t, vs.Proposer(0).Address(), val7.Address())
}

func TestJoinMoreThatOneThird(t *testing.T) {
	val1, _ := GenerateTestValidator(0)
	val2, _ := GenerateTestValidator(1)
	val3, _ := GenerateTestValidator(2)
	val4, _ := GenerateTestValidator(3)
	val5, _ := GenerateTestValidator(4)
	val6, _ := GenerateTestValidator(6)

	vs, err := NewValidatorSet([]*Validator{val1, val2, val3, val4}, 4, val1.Address())
	assert.NoError(t, err)

	assert.Error(t, vs.UpdateTheSet(0, []*Validator{val5, val6}))
}

func TestIsProposer(t *testing.T) {
	val1, _ := GenerateTestValidator(0)
	val2, _ := GenerateTestValidator(1)
	val3, _ := GenerateTestValidator(2)
	val4, _ := GenerateTestValidator(3)
	val5, _ := GenerateTestValidator(4)

	vs, err := NewValidatorSet([]*Validator{val1, val2, val3, val4}, 4, val1.Address())
	assert.NoError(t, err)

	assert.Equal(t, vs.Proposer(0).Address(), val1.Address())
	assert.Equal(t, vs.Proposer(1).Address(), val2.Address())
	assert.True(t, vs.IsProposer(val3.Address(), 2))
	assert.False(t, vs.IsProposer(val4.Address(), 2))
	assert.Equal(t, vs.validators, []*Validator{val1, val2, val3, val4})
	assert.Equal(t, vs.Validator(val2.Address()).Hash(), val2.Hash())
	assert.Nil(t, vs.Validator(val5.Address()))
}

func TestCommittersHash(t *testing.T) {
	val1, _ := GenerateTestValidator(0)
	val2, _ := GenerateTestValidator(1)
	val3, _ := GenerateTestValidator(2)
	val4, _ := GenerateTestValidator(3)

	vs, err := NewValidatorSet([]*Validator{val1, val2, val3, val4}, 4, val1.Address())
	assert.NoError(t, err)

	expected, _ := crypto.HashFromString("fd36b2597b028652ad4430b34a67094ba93ed84bd3abe5cd27f675bf431add48")
	assert.Equal(t, vs.CommittersHash(), expected)
}

func TestPower(t *testing.T) {
	vs, _ := GenerateTestValidatorSet()
	assert.Equal(t, vs.currentPower(), len(vs.validators))
}

func TestCopyValidators(t *testing.T) {
	vs, _ := GenerateTestValidatorSet()
	assert.Equal(t, vs.CopyValidators(), vs.validators)
}

func TestSortJoined(t *testing.T) {
	val1, _ := GenerateTestValidator(0)
	val2, _ := GenerateTestValidator(1)
	val3, _ := GenerateTestValidator(2)
	val4, _ := GenerateTestValidator(3)
	val5, _ := GenerateTestValidator(4)
	val6, _ := GenerateTestValidator(5)
	val7, _ := GenerateTestValidator(6)

	vs1, _ := NewValidatorSet([]*Validator{val1, val2, val3, val4}, 17, val1.Address())
	vs2, _ := NewValidatorSet([]*Validator{val1, val2, val3, val4}, 17, val1.Address())

	assert.NoError(t, vs1.UpdateTheSet(0, []*Validator{val5, val6, val7}))
	assert.NoError(t, vs2.UpdateTheSet(0, []*Validator{val7, val5, val6}))

	assert.Equal(t, vs1.CommittersHash(), vs2.CommittersHash())
}
