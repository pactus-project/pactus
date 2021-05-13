package committee

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/validator"
)

func TestContains(t *testing.T) {
	committee, signers := GenerateTestCommittee()
	nonExist, _, _ := crypto.GenerateTestKeyPair()

	assert.True(t, committee.Contains(signers[0].Address()))
	assert.False(t, committee.Contains(nonExist))
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
	_, pub1, _ := crypto.GenerateTestKeyPair()
	_, pub2, _ := crypto.GenerateTestKeyPair()
	_, pub3, _ := crypto.GenerateTestKeyPair()
	_, pub4, _ := crypto.GenerateTestKeyPair()
	_, pub5, _ := crypto.GenerateTestKeyPair()
	_, pub6, _ := crypto.GenerateTestKeyPair()
	_, pub7, _ := crypto.GenerateTestKeyPair()

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

	// Height 1
	assert.NoError(t, committee.Update(0, nil))
	assert.Equal(t, committee.Proposer(0).Number(), 2)
	assert.Equal(t, committee.Proposer(1).Number(), 3)
	assert.Equal(t, committee.Validators(), []*validator.Validator{val1, val2, val3, val4, val5, val6, val7})

	// Height 2
	assert.NoError(t, committee.Update(3, nil))
	assert.Equal(t, committee.Proposer(0).Number(), 6)

	// Height 3
	assert.NoError(t, committee.Update(1, nil))
	assert.Equal(t, committee.Proposer(0).Number(), 1)
}

func TestProposerJoin(t *testing.T) {
	_, pub1, _ := crypto.GenerateTestKeyPair()
	_, pub2, _ := crypto.GenerateTestKeyPair()
	_, pub3, _ := crypto.GenerateTestKeyPair()
	_, pub4, _ := crypto.GenerateTestKeyPair()
	_, pub5, _ := crypto.GenerateTestKeyPair()
	_, pub6, _ := crypto.GenerateTestKeyPair()
	_, pub7, _ := crypto.GenerateTestKeyPair()

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

	// Val1 is already in committee
	assert.Error(t, committee.Update(0, []*validator.Validator{val1}))

	//
	// +=+-+-+-+     +-+-+=+-+-+     +-+-+-+-+=+     +=+-+-+-+-+-+-+     +-+-+=+-+-+-+-+
	// |1|2|3|4| ==> |5|1|2|3|4| ==> |5|1|2|3|4| ==> |5|1|2|3|6|7|4| ==> |5|1|2|3|6|7|4|
	// +=+-+-+-+     +-+-+=+-+-+     +-+-+-+-+=+     +=+-+-+-+-+-+-+     +-+-+=+-+-+-+-+
	//

	// Height 1
	val5.UpdateLastJoinedHeight(1)
	assert.Equal(t, committee.Proposer(0).Number(), 1)
	assert.NoError(t, committee.Update(0, []*validator.Validator{val5}))
	assert.Equal(t, committee.Proposer(0).Number(), 2)
	assert.Equal(t, committee.Validators(), []*validator.Validator{val5, val1, val2, val3, val4})
	assert.Equal(t, committee.Size(), 5)

	// Height 2
	assert.NoError(t, committee.Update(1, nil))
	assert.Equal(t, committee.Proposer(0).Number(), 4)

	// Height 3
	val6.UpdateLastJoinedHeight(3)
	val7.UpdateLastJoinedHeight(3)
	assert.NoError(t, committee.Update(1, []*validator.Validator{val6, val7}))
	assert.Equal(t, committee.Proposer(0).Number(), 1)
	assert.Equal(t, committee.Validators(), []*validator.Validator{val5, val1, val2, val3, val6, val7, val4})
	assert.Equal(t, committee.Size(), 7)

	//
	assert.NoError(t, committee.Update(0, nil))
	assert.Equal(t, committee.Proposer(0).Number(), 2)
}

func TestProposerJoinAndLeave(t *testing.T) {
	_, pub1, _ := crypto.GenerateTestKeyPair()
	_, pub2, _ := crypto.GenerateTestKeyPair()
	_, pub3, _ := crypto.GenerateTestKeyPair()
	_, pub4, _ := crypto.GenerateTestKeyPair()
	_, pub5, _ := crypto.GenerateTestKeyPair()
	_, pub6, _ := crypto.GenerateTestKeyPair()
	_, pub7, _ := crypto.GenerateTestKeyPair()
	_, pub8, _ := crypto.GenerateTestKeyPair()
	_, pub9, _ := crypto.GenerateTestKeyPair()
	_, pubA, _ := crypto.GenerateTestKeyPair()
	_, pubB, _ := crypto.GenerateTestKeyPair()
	_, pubC, _ := crypto.GenerateTestKeyPair()
	_, pubD, _ := crypto.GenerateTestKeyPair()

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
	valD := validator.NewValidator(pubD, 13)

	committee, err := NewCommittee([]*validator.Validator{val1, val2, val3, val4, val5, val6, val7}, 7, val1.Address())
	assert.NoError(t, err)

	// How committee moves when new validator(s) join(s)?
	//
	// Example:
	//
	// Imagine validators `1` to `7` are in the committee, and `1` is the oldest and also proposer.
	// +=+-+-+-+-+-+-+
	// |1|2|3|4|5|6|7|
	// +=+-+-+-+-+-+-+
	//
	// New validator joins and sits before proposer.
	// In this example `8` sits before `1` (current proposer):
	// +*+=+-+-+-+-+-+-+
	// |8|1|2|3|4|5|6|7|
	// +*+=+-+-+-+-+-+-+
	//
	// Now committee should be adjusted and the oldest validator should leave.
	// In this example `1` is the oldest validator:
	// +-+-+-+-+-+-+-+
	// |8|2|3|4|5|6|7|
	// +-+-+-+-+-+-+-+
	//
	// Now we move to the next proposer.
	// In this example next proposer is `2`:
	// +-+=+-+-+-+-+-+
	// |8|2|3|4|5|6|7|
	// +-+=+-+-+-+-+-+
	//
	//
	// In this test we are covering these cases:
	//
	// +=+-+-+-+-+-+-+     +-+=+-+-+-+-+-+     +-+-+-+-+-+=+-+     +-+-+-+-+-+-+=+     +=+-+-+-+-+-+-+     +-+-+-+=+-+-+-+     +=+-+-+-+-+-+-+
	// |1|2|3|4|5|6|7| ==> |8|2|3|4|5|6|7| ==> |8|2|3|4|5|6|7| ==> |8|4|5|9|A|6|7| ==> |8|5|9|A|6|B|7| ==> |C|D|8|9|A|B|7| ==> |C|D|8|1|9|A|B|
	// +=+-+-+-+-+-+-+     +-+=+-+-+-+-+-+     +-+-+-+-+-+=+-+     +-+-+-+-+-+-+=+     +=+-+-+-+-+-+-+     +-+-+-+=+-+-+-+     +=+-+-+-+-+-+-+
	//

	// Height 1
	val8.UpdateLastJoinedHeight(1)
	assert.NoError(t, committee.Update(0, []*validator.Validator{val8}))
	assert.Equal(t, committee.Proposer(0).Number(), 2)
	assert.Equal(t, committee.Validators(), []*validator.Validator{val8, val2, val3, val4, val5, val6, val7})

	// Height 2
	assert.NoError(t, committee.Update(3, nil))
	assert.Equal(t, committee.Proposer(0).Number(), 6)

	// Height 3
	val9.UpdateLastJoinedHeight(3)
	valA.UpdateLastJoinedHeight(3)
	assert.NoError(t, committee.Update(0, []*validator.Validator{val9, valA}))
	assert.Equal(t, committee.Proposer(0).Number(), 7)
	assert.Equal(t, committee.Validators(), []*validator.Validator{val8, val4, val5, val9, valA, val6, val7})

	// Height 4
	valB.UpdateLastJoinedHeight(4)
	assert.NoError(t, committee.Update(0, []*validator.Validator{valB}))
	assert.Equal(t, committee.Proposer(0).Number(), 8)
	assert.Equal(t, committee.Proposer(1).Number(), 5)
	assert.Equal(t, committee.Proposer(2).Number(), 9)
	assert.Equal(t, committee.Validators(), []*validator.Validator{val8, val5, val9, valA, val6, valB, val7})

	// Height 5
	valC.UpdateLastJoinedHeight(5)
	valD.UpdateLastJoinedHeight(5)
	assert.NoError(t, committee.Update(0, []*validator.Validator{valC, valD}))
	assert.Equal(t, committee.Proposer(0).Number(), 9)
	assert.Equal(t, committee.Proposer(1).Number(), 10)
	assert.Equal(t, committee.Proposer(2).Number(), 11)
	assert.Equal(t, committee.Validators(), []*validator.Validator{valC, valD, val8, val9, valA, valB, val7})

	// Height 6
	val1.UpdateLastJoinedHeight(6)
	assert.NoError(t, committee.Update(2, []*validator.Validator{val1}))
	assert.Equal(t, committee.Proposer(0).Number(), 12)
	assert.Equal(t, committee.Proposer(1).Number(), 13)
	assert.Equal(t, committee.Proposer(2).Number(), 8)
	assert.Equal(t, committee.Proposer(3).Number(), 1)
	assert.Equal(t, committee.Proposer(4).Number(), 9)
	assert.Equal(t, committee.Proposer(5).Number(), 10)
	assert.Equal(t, committee.Proposer(6).Number(), 11)
	assert.Equal(t, committee.Validators(), []*validator.Validator{valC, valD, val8, val1, val9, valA, valB})
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
	assert.Equal(t, committee.Validators(), []*validator.Validator{val1, val2, val3, val4})
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

	vs1, _ := NewCommittee([]*validator.Validator{val1, val2, val3, val4}, 17, val1.Address())
	vs2, _ := NewCommittee([]*validator.Validator{val1, val2, val3, val4}, 17, val1.Address())

	assert.NoError(t, vs1.Update(0, []*validator.Validator{val5, val6, val7}))
	assert.NoError(t, vs2.Update(0, []*validator.Validator{val7, val5, val6}))
}

func TestTotalPower(t *testing.T) {
	_, pub, _ := crypto.GenerateTestKeyPair()
	val0 := validator.NewValidator(pub, 0) // Bootstrap validator
	val1, _ := validator.GenerateTestValidator(0)
	val2, _ := validator.GenerateTestValidator(1)
	val3, _ := validator.GenerateTestValidator(2)
	val4, _ := validator.GenerateTestValidator(3)

	committee, err := NewCommittee([]*validator.Validator{val0, val1, val2, val3, val4}, 4, val1.Address())
	assert.NoError(t, err)

	totalPower := val0.Power() + val1.Power() + val2.Power() + val3.Power() + val4.Power()
	totalStake := val0.Stake() + val1.Stake() + val2.Stake() + val3.Stake() + val4.Stake()
	assert.Equal(t, committee.TotalStake(), totalStake)
	assert.Equal(t, committee.TotalPower(), totalPower)
	assert.Equal(t, committee.TotalPower(), totalStake+1)
}
