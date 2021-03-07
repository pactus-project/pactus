package committee

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/validator"
)

func TestContains(t *testing.T) {
	committee, signers := GenerateTestCommittee()

	assert.True(t, committee.Contains(signers[0].Address()))
	assert.True(t, committee.Contains(committee.Proposer(0).Address()))
}

func TestProposer(t *testing.T) {
	committee, signers := GenerateTestCommittee()

	assert.Equal(t, committee.Proposer(0).Address(), signers[0].Address())
	assert.Equal(t, committee.Proposer(3).Address(), signers[3].Address())
	assert.Equal(t, committee.Proposer(4).Address(), signers[0].Address())

	assert.NoError(t, committee.Update(0, nil))
	assert.Equal(t, committee.Proposer(0).Address(), signers[1].Address())
}

func TestInvalidProposerJoinAndLeave(t *testing.T) {
	val1, _ := validator.GenerateTestValidator(0)
	val2, _ := validator.GenerateTestValidator(1)
	val3, _ := validator.GenerateTestValidator(2)
	val4, _ := validator.GenerateTestValidator(3)
	val5, _ := validator.GenerateTestValidator(4)

	committee, err := NewCommittee([]*validator.Validator{val1, val2, val3, val4}, 4, val5.Address())
	assert.Error(t, err)
	assert.Nil(t, committee)
}

func TestProposerMove(t *testing.T) {
	val1, _ := validator.GenerateTestValidator(1)
	val2, _ := validator.GenerateTestValidator(2)
	val3, _ := validator.GenerateTestValidator(3)
	val4, _ := validator.GenerateTestValidator(4)
	val5, _ := validator.GenerateTestValidator(5)
	val6, _ := validator.GenerateTestValidator(6)
	val7, _ := validator.GenerateTestValidator(7)

	committee, err := NewCommittee([]*validator.Validator{val1, val2, val3, val4, val5, val6, val7}, 7, val1.Address())
	assert.NoError(t, err)

	//
	// +=+-+-+-+-+-+-+       +-+=+-+-+-+-+-+
	// |1|2|3|4|5|6|7|  ==>  |1|2|3|4|5|6|7|
	// +=+-+-+-+-+-+-+       +-+=+-+-+-+-+-+
	//
	committee.proposerIndex = 0
	assert.Equal(t, committee.Proposer(0).Number(), 1)
	assert.NoError(t, committee.Update(0, nil))
	assert.Equal(t, committee.proposerIndex, 1)
	assert.Equal(t, committee.Proposer(0).Number(), 2)
	assert.Equal(t, committee.Proposer(1).Number(), 3)

	//
	// +-+-+-+=+-+-+-+       +-+-+-+-+=+-+-+
	// |1|2|3|4|5|6|7|  ==>  |1|2|3|4|5|6|7|
	// +-+-+-+=+-+-+-+       +-+-+-+-+=+-+-+
	//
	committee.proposerIndex = 3
	assert.Equal(t, committee.Proposer(0).Number(), 4)
	assert.NoError(t, committee.Update(0, nil))
	assert.Equal(t, committee.proposerIndex, 4)
	assert.Equal(t, committee.Proposer(0).Number(), 5)

	//
	// +-+-+-+-+-+-+=+       +=+-+-+-+-+-+-+
	// |1|2|3|4|5|6|7|  ==>  |1|2|3|4|5|6|7|
	// +-+-+-+-+-+-+=+       +=+-+-+-+-+-+-+
	//
	committee.proposerIndex = 6
	assert.Equal(t, committee.Proposer(0).Number(), 7)
	assert.NoError(t, committee.Update(0, nil))
	assert.Equal(t, committee.proposerIndex, 0)
	assert.Equal(t, committee.Proposer(0).Number(), 1)
}

func TestProposerMoveMoreRounds(t *testing.T) {
	val1, _ := validator.GenerateTestValidator(1)
	val2, _ := validator.GenerateTestValidator(2)
	val3, _ := validator.GenerateTestValidator(3)
	val4, _ := validator.GenerateTestValidator(4)
	val5, _ := validator.GenerateTestValidator(5)
	val6, _ := validator.GenerateTestValidator(6)
	val7, _ := validator.GenerateTestValidator(7)

	committee, err := NewCommittee([]*validator.Validator{val1, val2, val3, val4, val5, val6, val7}, 7, val1.Address())
	assert.NoError(t, err)

	//
	// +=+-+-+-+-+-+-+       +-+-+-+=+-+-+-+
	// |1|2|3|4|5|6|7|  ==>  |1|2|3|4|5|6|7|
	// +=+-+-+-+-+-+-+       +-+-+-+=+-+-+-+
	//
	committee.proposerIndex = 0
	assert.Equal(t, committee.Proposer(0).Number(), 1)
	assert.NoError(t, committee.Update(2, nil))
	assert.Equal(t, committee.proposerIndex, 3)
	assert.Equal(t, committee.Proposer(0).Number(), 4)
	assert.Equal(t, committee.Proposer(1).Number(), 5)

	//
	// +-+-+-+=+-+-+-+       +=+-+-+-+-+-+-+
	// |1|2|3|4|5|6|7|  ==>  |1|2|3|4|5|6|7|
	// +-+-+-+=+-+-+-+       +=+-+-+-+-+-+-+
	//
	committee.proposerIndex = 3
	assert.Equal(t, committee.Proposer(0).Number(), 4)
	assert.NoError(t, committee.Update(3, nil))
	assert.Equal(t, committee.proposerIndex, 0)
	assert.Equal(t, committee.Proposer(0).Number(), 1)

	//
	// +-+-+-+-+-+-+=+       +-+=+-+-+-+-+-+
	// |1|2|3|4|5|6|7|  ==>  |1|2|3|4|5|6|7|
	// +-+-+-+-+-+-+=+       +-+=+-+-+-+-+-+
	//
	committee.proposerIndex = 6
	assert.Equal(t, committee.Proposer(0).Number(), 7)
	assert.NoError(t, committee.Update(1, nil))
	assert.Equal(t, committee.proposerIndex, 1)
	assert.Equal(t, committee.Proposer(0).Number(), 2)
}

func TestProposerJoinAndLeave(t *testing.T) {
	val1, _ := validator.GenerateTestValidator(1)
	val2, _ := validator.GenerateTestValidator(2)
	val3, _ := validator.GenerateTestValidator(3)
	val4, _ := validator.GenerateTestValidator(4)
	val5, _ := validator.GenerateTestValidator(5)
	val6, _ := validator.GenerateTestValidator(6)
	val7, _ := validator.GenerateTestValidator(7)
	val8, _ := validator.GenerateTestValidator(8)
	val9, _ := validator.GenerateTestValidator(9)
	valA, _ := validator.GenerateTestValidator(10)
	valB, _ := validator.GenerateTestValidator(11)
	valC, _ := validator.GenerateTestValidator(12)
	valD, _ := validator.GenerateTestValidator(13)

	committee, err := NewCommittee([]*validator.Validator{val1, val2, val3, val4, val5, val6, val7}, 7, val1.Address())
	assert.NoError(t, err)

	// Val1 is already in committee
	assert.Error(t, committee.Update(0, []*validator.Validator{val1}))

	//
	// +=+-+-+-+-+-+-+       +=+-+-+-+-+-+-+
	// |1|2|3|4|5|6|7|  ==>  |2|3|4|5|6|7|8|
	// +=+-+-+-+-+-+-+       +=+-+-+-+-+-+-+
	//
	committee.proposerIndex = 0
	assert.Equal(t, committee.Proposer(0).Number(), 1)
	assert.NoError(t, committee.Update(0, []*validator.Validator{val8}))
	assert.Equal(t, committee.proposerIndex, 0)
	assert.Equal(t, committee.Proposer(0).Number(), 2)

	//
	// +-+-+=+-+-+-+-+       +-+=+-+-+-+-+-+
	// |2|3|4|5|6|7|8|  ==>  |4|5|6|7|8|9|A|
	// +-+-+=+-+-+-+-+       +-+=+-+-+-+-+-+
	//
	//
	committee.proposerIndex = 2
	assert.Equal(t, committee.Proposer(0).Number(), 4)
	assert.NoError(t, committee.Update(0, []*validator.Validator{val9, valA}))
	assert.Equal(t, committee.proposerIndex, 1)
	assert.Equal(t, committee.Proposer(0).Number(), 5)

	//
	// +-+-+-+-+-+-+=+       +=+-+-+-+-+-+-+
	// |4|5|6|7|8|9|A|  ==>  |5|6|7|8|9|A|B|
	// +-+-+-+-+-+-+=+       +=+-+-+-+-+-+-+
	//
	committee.proposerIndex = 6
	assert.Equal(t, committee.Proposer(0).Number(), 10)
	assert.NoError(t, committee.Update(0, []*validator.Validator{valB}))
	assert.Equal(t, committee.proposerIndex, 0)
	assert.Equal(t, committee.Proposer(0).Number(), 5)

	//
	// +-+-+-+-+-+-+=+       +=+-+-+-+-+-+-+
	// |5|6|7|8|9|A|B|  ==>  |7|8|9|A|B|C|D|
	// +-+-+-+-+-+-+=+       +=+-+-+-+-+-+-+
	//
	committee.proposerIndex = 6
	assert.Equal(t, committee.Proposer(0).Number(), 11)
	assert.NoError(t, committee.Update(0, []*validator.Validator{valC, valD}))
	assert.Equal(t, committee.proposerIndex, 0)
	assert.Equal(t, committee.Proposer(0).Number(), 7)
}

func TestProposerJoinAndLeaveMoreRound(t *testing.T) {
	val1, _ := validator.GenerateTestValidator(1)
	val2, _ := validator.GenerateTestValidator(2)
	val3, _ := validator.GenerateTestValidator(3)
	val4, _ := validator.GenerateTestValidator(4)
	val5, _ := validator.GenerateTestValidator(5)
	val6, _ := validator.GenerateTestValidator(6)
	val7, _ := validator.GenerateTestValidator(7)
	val8, _ := validator.GenerateTestValidator(8)
	val9, _ := validator.GenerateTestValidator(9)
	valA, _ := validator.GenerateTestValidator(10)
	valB, _ := validator.GenerateTestValidator(11)
	valC, _ := validator.GenerateTestValidator(12)
	valD, _ := validator.GenerateTestValidator(13)

	committee, err := NewCommittee([]*validator.Validator{val1, val2, val3, val4, val5, val6, val7}, 7, val1.Address())
	assert.NoError(t, err)

	// Val1 is already in committee
	assert.Error(t, committee.Update(0, []*validator.Validator{val1}))

	//
	// +=+-+-+-+-+-+-+       +-+-+=+-+-+-+-+
	// |1|2|3|4|5|6|7|  ==>  |2|3|4|5|6|7|8|
	// +=+-+-+-+-+-+-+       +-+-+=+-+-+-+-+
	//
	committee.proposerIndex = 0
	assert.Equal(t, committee.Proposer(0).Number(), 1)
	assert.NoError(t, committee.Update(2, []*validator.Validator{val8}))
	assert.Equal(t, committee.proposerIndex, 2)
	assert.Equal(t, committee.Proposer(0).Number(), 4)
	assert.Equal(t, committee.Proposer(1).Number(), 5)

	//
	// +-+-+=+-+-+-+-+       +-+-+-+-+=+-+-+
	// |2|3|4|5|6|7|8|  ==>  |4|5|6|7|8|9|A|
	// +-+-+=+-+-+-+-+       +-+-+-+-+=+-+-+
	//
	//
	committee.proposerIndex = 2
	assert.Equal(t, committee.Proposer(0).Number(), 4)
	assert.NoError(t, committee.Update(3, []*validator.Validator{val9, valA}))
	assert.Equal(t, committee.proposerIndex, 4)
	assert.Equal(t, committee.Proposer(0).Number(), 8)

	//
	// +-+-+-+-+-+-+=+       +-+=+-+-+-+-+-+
	// |4|5|6|7|8|9|A|  ==>  |5|6|7|8|9|A|B|
	// +-+-+-+-+-+-+=+       +-+=+-+-+-+-+-+
	//
	// 5 is offline
	committee.proposerIndex = 6
	assert.Equal(t, committee.Proposer(0).Number(), 10)
	assert.NoError(t, committee.Update(2, []*validator.Validator{valB}))
	assert.Equal(t, committee.proposerIndex, 1)
	assert.Equal(t, committee.Proposer(0).Number(), 6)

	//
	// +-+-+-+-+-+=+-+       +=+-+-+-+-+-+-+
	// |5|6|7|8|9|A|B|  ==>  |7|8|9|A|B|C|D|
	// +-+-+-+-+-+=+-+       +=+-+-+-+-+-+-+
	//
	committee.proposerIndex = 5
	assert.Equal(t, committee.Proposer(0).Number(), 10)
	assert.NoError(t, committee.Update(2, []*validator.Validator{valC, valD}))
	assert.Equal(t, committee.proposerIndex, 0)
	assert.Equal(t, committee.Proposer(0).Number(), 7)
}

func TestJoinMoreThatOneThird(t *testing.T) {
	val1, _ := validator.GenerateTestValidator(0)
	val2, _ := validator.GenerateTestValidator(1)
	val3, _ := validator.GenerateTestValidator(2)
	val4, _ := validator.GenerateTestValidator(3)
	val5, _ := validator.GenerateTestValidator(4)
	val6, _ := validator.GenerateTestValidator(6)

	committee, err := NewCommittee([]*validator.Validator{val1, val2, val3, val4}, 4, val1.Address())
	assert.NoError(t, err)

	assert.Error(t, committee.Update(0, []*validator.Validator{val5, val6}))
}

func TestIsProposer(t *testing.T) {
	val1, _ := validator.GenerateTestValidator(0)
	val2, _ := validator.GenerateTestValidator(1)
	val3, _ := validator.GenerateTestValidator(2)
	val4, _ := validator.GenerateTestValidator(3)
	val5, _ := validator.GenerateTestValidator(4)

	committee, err := NewCommittee([]*validator.Validator{val1, val2, val3, val4}, 4, val1.Address())
	assert.NoError(t, err)

	assert.Equal(t, committee.Proposer(0).Address(), val1.Address())
	assert.Equal(t, committee.Proposer(1).Address(), val2.Address())
	assert.True(t, committee.IsProposer(val3.Address(), 2))
	assert.False(t, committee.IsProposer(val4.Address(), 2))
	assert.Equal(t, committee.validators, []*validator.Validator{val1, val2, val3, val4})
	assert.Equal(t, committee.Validator(val2.Address()).Hash(), val2.Hash())
	assert.Nil(t, committee.Validator(val5.Address()))
}

func TestCommittee(t *testing.T) {
	val1, _ := validator.GenerateTestValidator(0)
	val2, _ := validator.GenerateTestValidator(1)
	val3, _ := validator.GenerateTestValidator(2)
	val4, _ := validator.GenerateTestValidator(3)

	committee, err := NewCommittee([]*validator.Validator{val1, val2, val3, val4}, 4, val1.Address())
	assert.NoError(t, err)

	expected, _ := crypto.HashFromString("fd36b2597b028652ad4430b34a67094ba93ed84bd3abe5cd27f675bf431add48")
	assert.Equal(t, committee.CommitteeHash(), expected)
	assert.Equal(t, committee.CommitteeHash(), expected)
}

func TestCopyValidators(t *testing.T) {
	committee, _ := GenerateTestCommittee()
	assert.Equal(t, committee.CopyValidators(), committee.validators)
}

func TestSortJoined(t *testing.T) {
	val1, _ := validator.GenerateTestValidator(0)
	val2, _ := validator.GenerateTestValidator(1)
	val3, _ := validator.GenerateTestValidator(2)
	val4, _ := validator.GenerateTestValidator(3)
	val5, _ := validator.GenerateTestValidator(4)
	val6, _ := validator.GenerateTestValidator(5)
	val7, _ := validator.GenerateTestValidator(6)

	committee1, _ := NewCommittee([]*validator.Validator{val1, val2, val3, val4}, 17, val1.Address())
	committee2, _ := NewCommittee([]*validator.Validator{val1, val2, val3, val4}, 17, val1.Address())

	assert.NoError(t, committee1.Update(0, []*validator.Validator{val5, val6, val7}))
	assert.NoError(t, committee2.Update(0, []*validator.Validator{val7, val5, val6}))

	assert.Equal(t, committee1.CommitteeHash(), committee2.CommitteeHash())
}

func TestCurrentPower(t *testing.T) {
	val1, _ := validator.GenerateTestValidator(0)
	val2, _ := validator.GenerateTestValidator(1)
	val3, _ := validator.GenerateTestValidator(2)
	val4, _ := validator.GenerateTestValidator(3)

	committee, err := NewCommittee([]*validator.Validator{val1, val2, val3, val4}, 4, val1.Address())
	assert.NoError(t, err)

	total := val1.Stake() + val2.Stake() + val3.Stake() + val4.Stake()
	assert.Equal(t, committee.currentPower(), total)
}
