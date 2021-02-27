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
	_, pub1, _ := crypto.GenerateTestKeyPair()
	_, pub2, _ := crypto.GenerateTestKeyPair()
	_, pub3, _ := crypto.GenerateTestKeyPair()
	_, pub4, _ := crypto.GenerateTestKeyPair()
	_, pub5, _ := crypto.GenerateTestKeyPair()
	_, pub6, _ := crypto.GenerateTestKeyPair()
	_, pub7, _ := crypto.GenerateTestKeyPair()

	val1 := NewValidator(pub1, 1, 0)
	val2 := NewValidator(pub2, 2, 0)
	val3 := NewValidator(pub3, 3, 0)
	val4 := NewValidator(pub4, 4, 0)
	val5 := NewValidator(pub5, 5, 0)
	val6 := NewValidator(pub6, 6, 0)
	val7 := NewValidator(pub7, 7, 0)

	vs, err := NewValidatorSet([]*Validator{val1, val2, val3, val4, val5, val6, val7}, 7, val1.Address())
	assert.NoError(t, err)

	//
	// +=+-+-+-+-+-+-+     +-+=+-+-+-+-+-+     +-+-+-+-+-+=+-+     +=+-+-+-+-+-+-+
	// |1|2|3|4|5|6|7| ==> |1|2|3|4|5|6|7| ==> |1|2|3|4|5|6|7| ==> |1|2|3|4|5|6|7|
	// +=+-+-+-+-+-+-+     +-+=+-+-+-+-+-+     +-+-+-+-+-+=+-+     +=+-+-+-+-+-+-+
	//

	// Height 1
	assert.NoError(t, vs.UpdateTheSet(0, nil))
	assert.Equal(t, vs.Proposer(0).Number(), 2)
	assert.Equal(t, vs.Proposer(1).Number(), 3)
	assert.Equal(t, vs.Validators(), []*Validator{val1, val2, val3, val4, val5, val6, val7})

	// Height 2
	assert.NoError(t, vs.UpdateTheSet(3, nil))
	assert.Equal(t, vs.Proposer(0).Number(), 6)

	// Height 3
	assert.NoError(t, vs.UpdateTheSet(1, nil))
	assert.Equal(t, vs.Proposer(0).Number(), 1)
}

func TestProposerJoin(t *testing.T) {
	_, pub1, _ := crypto.GenerateTestKeyPair()
	_, pub2, _ := crypto.GenerateTestKeyPair()
	_, pub3, _ := crypto.GenerateTestKeyPair()
	_, pub4, _ := crypto.GenerateTestKeyPair()
	_, pub5, _ := crypto.GenerateTestKeyPair()
	_, pub6, _ := crypto.GenerateTestKeyPair()
	_, pub7, _ := crypto.GenerateTestKeyPair()

	val1 := NewValidator(pub1, 1, 0)
	val2 := NewValidator(pub2, 2, 0)
	val3 := NewValidator(pub3, 3, 0)
	val4 := NewValidator(pub4, 4, 0)
	val5 := NewValidator(pub5, 5, 0)
	val6 := NewValidator(pub6, 6, 0)
	val7 := NewValidator(pub7, 7, 0)

	vs, err := NewValidatorSet([]*Validator{val1, val2, val3, val4}, 7, val1.Address())
	assert.NoError(t, err)

	// Val1 is already in set
	assert.Error(t, vs.UpdateTheSet(0, []*Validator{val1}))
	// More than 2/3 of power
	assert.Error(t, vs.UpdateTheSet(0, []*Validator{val5, val6, val7}))

	//
	// +=+-+-+-+     +-+-+=+-+-+     +-+-+-+-+=+     +=+-+-+-+-+-+-+     +-+-+=+-+-+-+-+
	// |1|2|3|4| ==> |5|1|2|3|4| ==> |5|1|2|3|4| ==> |5|1|2|3|6|7|4| ==> |5|1|2|3|6|7|4|
	// +=+-+-+-+     +-+-+=+-+-+     +-+-+-+-+=+     +=+-+-+-+-+-+-+     +-+-+=+-+-+-+-+
	//

	// Height 1
	val5.UpdateLastJoinedHeight(1)
	assert.Equal(t, vs.Proposer(0).Number(), 1)
	assert.NoError(t, vs.UpdateTheSet(0, []*Validator{val5}))
	assert.Equal(t, vs.Proposer(0).Number(), 2)
	assert.Equal(t, vs.Validators(), []*Validator{val5, val1, val2, val3, val4})

	// Height 2
	assert.NoError(t, vs.UpdateTheSet(1, nil))
	assert.Equal(t, vs.Proposer(0).Number(), 4)

	// Height 3
	val6.UpdateLastJoinedHeight(3)
	val7.UpdateLastJoinedHeight(3)
	assert.NoError(t, vs.UpdateTheSet(1, []*Validator{val6, val7}))
	assert.Equal(t, vs.Proposer(0).Number(), 1)
	assert.Equal(t, vs.Validators(), []*Validator{val5, val1, val2, val3, val6, val7, val4})

	//
	assert.NoError(t, vs.UpdateTheSet(0, nil))
	assert.Equal(t, vs.Proposer(0).Number(), 2)
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

	val1 := NewValidator(pub1, 1, 0)
	val2 := NewValidator(pub2, 2, 0)
	val3 := NewValidator(pub3, 3, 0)
	val4 := NewValidator(pub4, 4, 0)
	val5 := NewValidator(pub5, 5, 0)
	val6 := NewValidator(pub6, 6, 0)
	val7 := NewValidator(pub7, 7, 0)
	val8 := NewValidator(pub8, 8, 0)
	val9 := NewValidator(pub9, 9, 0)
	valA := NewValidator(pubA, 10, 0)
	valB := NewValidator(pubB, 11, 0)
	valC := NewValidator(pubC, 12, 0)
	valD := NewValidator(pubD, 13, 0)

	vs, err := NewValidatorSet([]*Validator{val1, val2, val3, val4, val5, val6, val7}, 7, val1.Address())
	assert.NoError(t, err)

	// How validator set moves when new validator(s) join(s)?
	//
	// Example:
	//
	// Imagine validators `1` to `7` are in the set, and `1` is the oldest and also proposer.
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
	// Now validator set should be adjusted and the oldest validator should leave.
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
	assert.NoError(t, vs.UpdateTheSet(0, []*Validator{val8}))
	assert.Equal(t, vs.Proposer(0).Number(), 2)
	assert.Equal(t, vs.Validators(), []*Validator{val8, val2, val3, val4, val5, val6, val7})

	// Height 2
	assert.NoError(t, vs.UpdateTheSet(3, nil))
	assert.Equal(t, vs.Proposer(0).Number(), 6)

	// Height 3
	val9.UpdateLastJoinedHeight(3)
	valA.UpdateLastJoinedHeight(3)
	assert.NoError(t, vs.UpdateTheSet(0, []*Validator{val9, valA}))
	assert.Equal(t, vs.Proposer(0).Number(), 7)
	assert.Equal(t, vs.Validators(), []*Validator{val8, val4, val5, val9, valA, val6, val7})

	// Height 4
	valB.UpdateLastJoinedHeight(4)
	assert.NoError(t, vs.UpdateTheSet(0, []*Validator{valB}))
	assert.Equal(t, vs.Proposer(0).Number(), 8)
	assert.Equal(t, vs.Proposer(1).Number(), 5)
	assert.Equal(t, vs.Proposer(2).Number(), 9)
	assert.Equal(t, vs.Validators(), []*Validator{val8, val5, val9, valA, val6, valB, val7})

	// Height 5
	valC.UpdateLastJoinedHeight(5)
	valD.UpdateLastJoinedHeight(5)
	assert.NoError(t, vs.UpdateTheSet(0, []*Validator{valC, valD}))
	assert.Equal(t, vs.Proposer(0).Number(), 9)
	assert.Equal(t, vs.Proposer(1).Number(), 10)
	assert.Equal(t, vs.Proposer(2).Number(), 11)
	assert.Equal(t, vs.Validators(), []*Validator{valC, valD, val8, val9, valA, valB, val7})

	// Height 6
	val1.UpdateLastJoinedHeight(6)
	assert.NoError(t, vs.UpdateTheSet(2, []*Validator{val1}))
	assert.Equal(t, vs.Proposer(0).Number(), 12)
	assert.Equal(t, vs.Proposer(1).Number(), 13)
	assert.Equal(t, vs.Proposer(2).Number(), 8)
	assert.Equal(t, vs.Proposer(3).Number(), 1)
	assert.Equal(t, vs.Proposer(4).Number(), 9)
	assert.Equal(t, vs.Proposer(5).Number(), 10)
	assert.Equal(t, vs.Proposer(6).Number(), 11)
	assert.Equal(t, vs.Validators(), []*Validator{valC, valD, val8, val1, val9, valA, valB})
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
	///	assert.Equal(t, vs.validators, []*Validator{val1, val2, val3, val4})
	assert.Equal(t, vs.Validator(val2.Address()).Hash(), val2.Hash())
	assert.Nil(t, vs.Validator(val5.Address()))
}

func TestCommittee(t *testing.T) {
	val1, _ := GenerateTestValidator(0)
	val2, _ := GenerateTestValidator(1)
	val3, _ := GenerateTestValidator(2)
	val4, _ := GenerateTestValidator(3)

	vs, err := NewValidatorSet([]*Validator{val1, val2, val3, val4}, 4, val1.Address())
	assert.NoError(t, err)

	expected, _ := crypto.HashFromString("fd36b2597b028652ad4430b34a67094ba93ed84bd3abe5cd27f675bf431add48")
	assert.Equal(t, vs.Committee(), []int{0, 1, 2, 3})
	assert.Equal(t, vs.CommitteeHash(), expected)
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

	assert.Equal(t, vs1.CommitteeHash(), vs2.CommitteeHash())
}

func TestCurrentPower(t *testing.T) {
	val1, _ := GenerateTestValidator(0)
	val2, _ := GenerateTestValidator(1)
	val3, _ := GenerateTestValidator(2)
	val4, _ := GenerateTestValidator(3)

	vs, err := NewValidatorSet([]*Validator{val1, val2, val3, val4}, 4, val1.Address())
	assert.NoError(t, err)

	total := val1.Stake() + val2.Stake() + val3.Stake() + val4.Stake()
	assert.Equal(t, vs.currentPower(), total)
}
