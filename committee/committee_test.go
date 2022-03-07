package committee

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/crypto/bls"
	"github.com/zarbchain/zarb-go/validator"
)

func TestContains(t *testing.T) {
	committee, signers := GenerateTestCommittee()
	nonExist := crypto.GenerateTestAddress()

	assert.True(t, committee.Contains(signers[0].Address()))
	assert.False(t, committee.Contains(nonExist))
}

func TestProposer(t *testing.T) {
	committee, signers := GenerateTestCommittee()

	assert.Equal(t, committee.Proposer(0).Address(), signers[0].Address())
	assert.Equal(t, committee.Proposer(3).Address(), signers[3].Address())
	assert.Equal(t, committee.Proposer(4).Address(), signers[0].Address())

	committee.Update(0, nil)
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
	pub1, _ := bls.GenerateTestKeyPair()
	pub2, _ := bls.GenerateTestKeyPair()
	pub3, _ := bls.GenerateTestKeyPair()
	pub4, _ := bls.GenerateTestKeyPair()
	pub5, _ := bls.GenerateTestKeyPair()
	pub6, _ := bls.GenerateTestKeyPair()
	pub7, _ := bls.GenerateTestKeyPair()

	val1 := validator.NewValidator(pub1, 1)
	val2 := validator.NewValidator(pub2, 2)
	val3 := validator.NewValidator(pub3, 3)
	val4 := validator.NewValidator(pub4, 4)
	val5 := validator.NewValidator(pub5, 5)
	val6 := validator.NewValidator(pub6, 6)
	val7 := validator.NewValidator(pub7, 7)

	committee, err := NewCommittee([]*validator.Validator{val1, val2, val3, val4, val5, val6, val7}, 7, val1.Address())
	assert.NoError(t, err)

	//
	// +=+-+-+-+-+-+-+     +-+=+-+-+-+-+-+     +-+-+-+-+-+=+-+     +=+-+-+-+-+-+-+
	// |1|2|3|4|5|6|7| ==> |1|2|3|4|5|6|7| ==> |1|2|3|4|5|6|7| ==> |1|2|3|4|5|6|7|
	// +=+-+-+-+-+-+-+     +-+=+-+-+-+-+-+     +-+-+-+-+-+=+-+     +=+-+-+-+-+-+-+
	//

	// Height 1001
	committee.Update(0, nil)
	assert.Equal(t, committee.Proposer(0).Number(), 2)
	assert.Equal(t, committee.Proposer(1).Number(), 3)
	assert.Equal(t, committee.Validators(), []*validator.Validator{val1, val2, val3, val4, val5, val6, val7})

	// Height 1002
	committee.Update(3, nil)
	assert.Equal(t, committee.Proposer(0).Number(), 6)

	// Height 1003
	committee.Update(1, nil)
	assert.Equal(t, committee.Proposer(0).Number(), 1)
}

func TestProposerJoin(t *testing.T) {
	pub1, _ := bls.GenerateTestKeyPair()
	pub2, _ := bls.GenerateTestKeyPair()
	pub3, _ := bls.GenerateTestKeyPair()
	pub4, _ := bls.GenerateTestKeyPair()
	pub5, _ := bls.GenerateTestKeyPair()
	pub6, _ := bls.GenerateTestKeyPair()
	pub7, _ := bls.GenerateTestKeyPair()

	val1 := validator.NewValidator(pub1, 1)
	val2 := validator.NewValidator(pub2, 2)
	val3 := validator.NewValidator(pub3, 3)
	val4 := validator.NewValidator(pub4, 4)
	val5 := validator.NewValidator(pub5, 5)
	val6 := validator.NewValidator(pub6, 6)
	val7 := validator.NewValidator(pub7, 7)

	committee, err := NewCommittee([]*validator.Validator{val1, val2, val3, val4}, 7, val1.Address())
	assert.NoError(t, err)
	assert.Equal(t, committee.Size(), 4)

	//
	// r=0           r=0             r=1             r=1                 r=0
	// +=+-+-+-+     +-+*+=+-+-+     +-+-+-+!+=+     +*+*+!+=+-+-+-+     +-+-+-+-+=+-+-+
	// |1|2|3|4| ==> |1|5|2|3|4| ==> |1|5|2|3|4| ==> |6|7|1|5|2|3|4| ==> |6|7|1|5|2|3|4|
	// +-+-+-+-+     +-+-+-+-+-+     +-+-+-+-+-+     +-+-+-+-+-+-+-+     +-+-+-+-+-+-+-+
	//

	// Height 1000
	// Val1 is in the committee
	val2.UpdateLastJoinedHeight(1000)
	committee.Update(0, []*validator.Validator{val2})
	assert.Equal(t, committee.Proposer(0).Number(), 2)
	assert.Equal(t, committee.Committers(), []int{1, 2, 3, 4})
	assert.Equal(t, committee.Size(), 4)

	// Height 1001
	val5.UpdateLastJoinedHeight(1001)
	committee.Update(0, []*validator.Validator{val5})
	assert.Equal(t, committee.Proposer(0).Number(), 3)
	assert.Equal(t, committee.Proposer(1).Number(), 4)
	assert.Equal(t, committee.Committers(), []int{1, 5, 2, 3, 4})
	assert.Equal(t, committee.Size(), 5)

	// Height 1002
	committee.Update(1, nil)
	assert.Equal(t, committee.Proposer(0).Number(), 1)
	assert.Equal(t, committee.Proposer(1).Number(), 5)

	// Height 1003
	val3.UpdateLastJoinedHeight(1003)
	val6.UpdateLastJoinedHeight(1003)
	val7.UpdateLastJoinedHeight(1003)
	committee.Update(1, []*validator.Validator{val7, val3, val6})
	assert.Equal(t, committee.Proposer(0).Number(), 2)
	assert.Equal(t, committee.Committers(), []int{6, 7, 1, 5, 2, 3, 4})
	assert.Equal(t, committee.Size(), 7)

	// Height 1004
	committee.Update(0, nil)
	assert.Equal(t, committee.Proposer(0).Number(), 3)
}

func TestProposerJoinAndLeave(t *testing.T) {
	pub1, _ := bls.GenerateTestKeyPair()
	pub2, _ := bls.GenerateTestKeyPair()
	pub3, _ := bls.GenerateTestKeyPair()
	pub4, _ := bls.GenerateTestKeyPair()
	pub5, _ := bls.GenerateTestKeyPair()
	pub6, _ := bls.GenerateTestKeyPair()
	pub7, _ := bls.GenerateTestKeyPair()
	pub8, _ := bls.GenerateTestKeyPair()
	pub9, _ := bls.GenerateTestKeyPair()
	pubA, _ := bls.GenerateTestKeyPair()
	pubB, _ := bls.GenerateTestKeyPair()
	pubC, _ := bls.GenerateTestKeyPair()

	val1 := validator.NewValidator(pub1, 1)
	val2 := validator.NewValidator(pub2, 2)
	val3 := validator.NewValidator(pub3, 3)
	val4 := validator.NewValidator(pub4, 4)
	val5 := validator.NewValidator(pub5, 5)
	val6 := validator.NewValidator(pub6, 6)
	val7 := validator.NewValidator(pub7, 7)
	val8 := validator.NewValidator(pub8, 8)
	val9 := validator.NewValidator(pub9, 9)
	valA := validator.NewValidator(pubA, 10)
	valB := validator.NewValidator(pubB, 11)
	valC := validator.NewValidator(pubC, 12)

	committee, err := NewCommittee([]*validator.Validator{val1, val2, val3, val4, val5, val6, val7}, 7, val1.Address())
	assert.NoError(t, err)

	// How committee moves when new validator(s) join(s)?
	//
	// Example:
	//
	// Validators `1` to `7` are in the committee, and `1` is the oldest and also proposer.
	// +=+-+-+-+-+-+-+
	// |1|2|3|4|5|6|7|
	// +-+-+-+-+-+-+-+
	//
	// New validator joins and seats before proposer.
	// In this example `8` seats before `1` (current proposer):
	// +*+=+-+-+-+-+-+-+
	// |8|1|2|3|4|5|6|7|
	// +-+-+-+-+-+-+-+-+
	//
	// Now committee should be adjusted and the oldest validator should leave.
	// In this example `1` also is the oldest validator:
	// +*+-+-+-+-+-+-+
	// |8|2|3|4|5|6|7|
	// +-+-+-+-+-+-+-+
	//
	// Now we move to the next proposer.
	// +-+=+-+-+-+-+-+
	// |8|2|3|4|5|6|7|
	// +-+-+-+-+-+-+-+
	//
	//-------------------------------------
	// In this test we cover these movement:
	//
	// h=1000              h=1001              h=1002              h=1003              h=1004              h=1005              h=1006              h=1007
	// r=0                 r=3                 r=0                 r=0                 r=1                 r=0                 r=2                 r=4
	// +=+-+-+-+-+-+-+     +*+-+-+-+-+-+-+     +-+!+!+!+=+-+-+     +-+-+-+-+*+=+-+     +=+-+-+-+-+*+!+     +-+*+*+!+!+=+-+     +!+=+-+-+-+-+!+     +*+*+-+-+!+!+=+
	// |1|2|3|4|5|6|7| ==> |8|2|3|4|5|6|7| ==> |8|2|3|4|5|6|7| ==> |8|2|3|5|9|6|7| ==> |8|2|3|9|6|A|7| ==> |8|B|C|2|3|9|A| ==> |B|C|2|1|3|9|A| ==> |5|6|B|C|2|1|3|
	// +-+-+-+-+-+-+-+     +-+-+-+-+-+-+-+     +-+-+-+-+-+-+-+     +-+-+-+-+-+-+-+     +-+-+-+-+-+-+-+     +-+-+-+-+-+-+-+     +-+-+-+-+-+-+-+     +-+-+-+-+-+-+-+
	//

	// Height 1001
	val8.UpdateLastJoinedHeight(1001)
	committee.Update(0, []*validator.Validator{val8})
	assert.Equal(t, committee.Proposer(0).Number(), 2)
	assert.Equal(t, committee.Proposer(1).Number(), 3)
	assert.Equal(t, committee.Proposer(2).Number(), 4)
	assert.Equal(t, committee.Committers(), []int{8, 2, 3, 4, 5, 6, 7})

	// Height 1002
	val3.UpdateLastJoinedHeight(1002)
	committee.Update(3, []*validator.Validator{val3})
	assert.Equal(t, committee.Proposer(0).Number(), 6)

	// Height 1003
	val2.UpdateLastJoinedHeight(1003)
	val9.UpdateLastJoinedHeight(1003)
	committee.Update(0, []*validator.Validator{val9, val2})
	assert.Equal(t, committee.Proposer(0).Number(), 7)
	assert.Equal(t, committee.Proposer(1).Number(), 8)
	assert.Equal(t, committee.Committers(), []int{8, 2, 3, 5, 9, 6, 7})

	// Height 1004
	valA.UpdateLastJoinedHeight(1004)
	committee.Update(1, []*validator.Validator{valA})
	assert.Equal(t, committee.Proposer(0).Number(), 2)
	assert.Equal(t, committee.Committers(), []int{8, 2, 3, 9, 6, 10, 7})

	// Height 1005
	valB.UpdateLastJoinedHeight(1005)
	valC.UpdateLastJoinedHeight(1005)
	committee.Update(0, []*validator.Validator{valC, valB})
	assert.Equal(t, committee.Proposer(0).Number(), 3)
	assert.Equal(t, committee.Proposer(1).Number(), 9)
	assert.Equal(t, committee.Proposer(2).Number(), 10)
	assert.Equal(t, committee.Committers(), []int{8, 11, 12, 2, 3, 9, 10})

	// Height 1006
	val1.UpdateLastJoinedHeight(1006)
	committee.Update(2, []*validator.Validator{val1})
	assert.Equal(t, committee.Proposer(0).Number(), 11)
	assert.Equal(t, committee.Committers(), []int{11, 12, 2, 1, 3, 9, 10})

	// Height 1007
	val2.UpdateLastJoinedHeight(1007)
	val3.UpdateLastJoinedHeight(1007)
	val5.UpdateLastJoinedHeight(1007)
	val6.UpdateLastJoinedHeight(1007)
	committee.Update(4, []*validator.Validator{val2, val3, val5, val6})
	assert.Equal(t, committee.Proposer(0).Number(), 5)
	assert.Equal(t, committee.Committers(), []int{5, 6, 11, 12, 2, 1, 3})
}

func TestIsProposer(t *testing.T) {
	val1, _ := validator.GenerateTestValidator(0)
	val2, _ := validator.GenerateTestValidator(1)
	val3, _ := validator.GenerateTestValidator(2)
	val4, _ := validator.GenerateTestValidator(3)
	//val5, _ := validator.GenerateTestValidator(4)

	committee, err := NewCommittee([]*validator.Validator{val1, val2, val3, val4}, 4, val1.Address())
	assert.NoError(t, err)

	assert.Equal(t, committee.Proposer(0).Address(), val1.Address())
	assert.Equal(t, committee.Proposer(1).Address(), val2.Address())
	assert.True(t, committee.IsProposer(val3.Address(), 2))
	assert.False(t, committee.IsProposer(val4.Address(), 2))
	assert.Equal(t, committee.Validators(), []*validator.Validator{val1, val2, val3, val4})
	// assert.Equal(t, committee.Validator(val2.Address()).Hash(), val2.Hash())
	// assert.Nil(t, committee.Validator(val5.Address()))
}

func TestCommittee(t *testing.T) {
	val1, _ := validator.GenerateTestValidator(0)
	val2, _ := validator.GenerateTestValidator(1)
	val3, _ := validator.GenerateTestValidator(2)
	val4, _ := validator.GenerateTestValidator(3)

	committee, err := NewCommittee([]*validator.Validator{val1, val2, val3, val4}, 4, val1.Address())
	assert.NoError(t, err)
	assert.Equal(t, committee.Committers(), []int{0, 1, 2, 3})
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

	committee1.Update(0, []*validator.Validator{val5, val6, val7})
	committee2.Update(0, []*validator.Validator{val7, val5, val6})
}

func TestTotalPower(t *testing.T) {
	pub, _ := bls.GenerateTestKeyPair()
	val0 := validator.NewValidator(pub, 0) // Bootstrap validator
	val1, _ := validator.GenerateTestValidator(0)
	val2, _ := validator.GenerateTestValidator(1)
	val3, _ := validator.GenerateTestValidator(2)
	val4, _ := validator.GenerateTestValidator(3)

	committee, err := NewCommittee([]*validator.Validator{val0, val1, val2, val3, val4}, 4, val1.Address())
	assert.NoError(t, err)

	totalPower := val0.Power() + val1.Power() + val2.Power() + val3.Power() + val4.Power()
	totalStake := val0.Stake() + val1.Stake() + val2.Stake() + val3.Stake() + val4.Stake()
	assert.Equal(t, committee.TotalPower(), totalPower)
	assert.Equal(t, committee.TotalPower(), totalStake+1)
}
