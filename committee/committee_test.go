package committee_test

import (
	"fmt"
	"testing"

	"github.com/pactus-project/pactus/committee"
	"github.com/pactus-project/pactus/types/validator"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
)

func TestContains(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	cmt, valKeys := ts.GenerateTestCommittee(21)
	nonExist := ts.RandAccAddress()

	assert.True(t, cmt.Contains(valKeys[0].Address()))
	assert.False(t, cmt.Contains(nonExist))
}

func TestProposer(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	cmt, _ := ts.GenerateTestCommittee(4)

	assert.Equal(t, cmt.Proposer(0).Number(), int32(0))
	assert.Equal(t, cmt.Proposer(3).Number(), int32(3))
	assert.Equal(t, cmt.Proposer(4).Number(), int32(0))

	cmt.Update(0, nil)
	assert.Equal(t, cmt.Proposer(0).Number(), int32(1))
}

func TestInvalidProposerJoinAndLeave(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	val1, _ := ts.GenerateTestValidator(0)
	val2, _ := ts.GenerateTestValidator(1)
	val3, _ := ts.GenerateTestValidator(2)
	val4, _ := ts.GenerateTestValidator(3)
	val5, _ := ts.GenerateTestValidator(4)

	cmt, err := committee.NewCommittee([]*validator.Validator{val1, val2, val3, val4}, 4, val5.Address())
	assert.Error(t, err)
	assert.Nil(t, cmt)
}

func TestProposerMove(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	val1, _ := ts.GenerateTestValidator(1)
	val2, _ := ts.GenerateTestValidator(2)
	val3, _ := ts.GenerateTestValidator(3)
	val4, _ := ts.GenerateTestValidator(4)
	val5, _ := ts.GenerateTestValidator(5)
	val6, _ := ts.GenerateTestValidator(6)
	val7, _ := ts.GenerateTestValidator(7)

	cmt, err := committee.NewCommittee(
		[]*validator.Validator{val1, val2, val3, val4, val5, val6, val7}, 7, val1.Address())
	assert.NoError(t, err)

	//
	// +*+-+-+-+-+-+-+    +-+*+-+-+-+-+-+    +-+-+-+-+-+*+-+    +*+-+-+-+-+-+-+
	// |1|2|3|4|5|6|7| => |1|2|3|4|5|6|7| => |1|2|3|4|5|6|7| => |1|2|3|4|5|6|7|
	// +-+-+-+-+-+-+-+    +-+-+-+-+-+-+-+    +-+-+-+-+-+-+-+    +-+-+-+-+-+-+-+
	//

	assert.Equal(t, cmt.Proposer(0).Number(), int32(1))
	assert.Equal(t, cmt.Proposer(7).Number(), int32(1))
	assert.Equal(t, cmt.Validators(), []*validator.Validator{val1, val2, val3, val4, val5, val6, val7})
	fmt.Println(cmt.String())

	// Height 1001
	cmt.Update(0, nil)
	assert.Equal(t, cmt.Proposer(0).Number(), int32(2))
	assert.Equal(t, cmt.Proposer(1).Number(), int32(3))
	assert.Equal(t, cmt.Proposer(7).Number(), int32(2))
	assert.Equal(t, cmt.Validators(), []*validator.Validator{val1, val2, val3, val4, val5, val6, val7})
	fmt.Println(cmt.String())

	// Height 1002
	cmt.Update(3, nil)
	assert.Equal(t, cmt.Proposer(0).Number(), int32(6))
	assert.Equal(t, cmt.Validators(), []*validator.Validator{val1, val2, val3, val4, val5, val6, val7})
	fmt.Println(cmt.String())

	// Height 1003
	cmt.Update(1, nil)
	assert.Equal(t, cmt.Proposer(0).Number(), int32(1))
	assert.Equal(t, cmt.Validators(), []*validator.Validator{val1, val2, val3, val4, val5, val6, val7})
	fmt.Println(cmt.String())
}

func TestValidatorConsistency(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	val1, _ := ts.GenerateTestValidator(1)
	val2, _ := ts.GenerateTestValidator(2)
	val3, _ := ts.GenerateTestValidator(3)
	val4, _ := ts.GenerateTestValidator(4)

	cmt, _ := committee.NewCommittee([]*validator.Validator{val1, val2, val3, val4}, 4, val1.Address())

	t.Run("Updating validators' stake, Should panic", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("The code did not panic")
			}
		}()

		val1.AddToStake(1)
		cmt.Update(0, []*validator.Validator{val1})
	})
}

func TestProposerJoin(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	val1, _ := ts.GenerateTestValidator(1)
	val2, _ := ts.GenerateTestValidator(2)
	val3, _ := ts.GenerateTestValidator(3)
	val4, _ := ts.GenerateTestValidator(4)
	val5, _ := ts.GenerateTestValidator(5)
	val6, _ := ts.GenerateTestValidator(6)
	val7, _ := ts.GenerateTestValidator(7)

	cmt, err := committee.NewCommittee([]*validator.Validator{val1, val2, val3, val4}, 7, val1.Address())
	assert.NoError(t, err)
	assert.Equal(t, cmt.Size(), 4)

	//
	// h=1000, r=0  h=1001, r=0    h=1002, r=1    h=1003, r=1        h=1004, r=0
	// +-+*+-+-+    +-+$+-+*+-+    +*+-+-+!+-+    +$+$+-+!+*+-+-+    +-+-+-+-+-+*+-+
	// |1|2|3|4| => |1|5|2|3|4| => |1|5|2|3|4| => |6|7|1|5|2|3|4| => |6|7|1|5|2|3|4|
	// +-+-+-+-+    +-+-+-+-+-+    +-+-+-+-+-+    +-+-+-+-+-+-+-+    +-+-+-+-+-+-+-+
	//

	// Height 1000
	// Val2 is in the committee
	val2.UpdateLastSortitionHeight(1000)
	cmt.Update(0, []*validator.Validator{val2})
	assert.Equal(t, cmt.Proposer(0).Number(), int32(2))
	assert.Equal(t, cmt.Committers(), []int32{1, 2, 3, 4})
	assert.Equal(t, cmt.Size(), 4)
	fmt.Println(cmt.String())

	// Height 1001
	val5.UpdateLastSortitionHeight(1001)
	cmt.Update(0, []*validator.Validator{val5})
	assert.Equal(t, cmt.Proposer(0).Number(), int32(3))
	assert.Equal(t, cmt.Proposer(1).Number(), int32(4))
	assert.Equal(t, cmt.Committers(), []int32{1, 5, 2, 3, 4})
	assert.Equal(t, cmt.Size(), 5)
	fmt.Println(cmt.String())

	// Height 1002
	cmt.Update(1, nil)
	assert.Equal(t, cmt.Proposer(0).Number(), int32(1))
	assert.Equal(t, cmt.Proposer(1).Number(), int32(5))
	fmt.Println(cmt.String())

	// Height 1003
	val3.UpdateLastSortitionHeight(1003)
	val6.UpdateLastSortitionHeight(1003)
	val7.UpdateLastSortitionHeight(1003)
	cmt.Update(1, []*validator.Validator{val7, val3, val6})
	assert.Equal(t, cmt.Proposer(0).Number(), int32(2))
	assert.Equal(t, cmt.Committers(), []int32{6, 7, 1, 5, 2, 3, 4})
	assert.Equal(t, cmt.Size(), 7)
	fmt.Println(cmt.String())

	// Height 1004
	cmt.Update(0, nil)
	assert.Equal(t, cmt.Proposer(0).Number(), int32(3))
	fmt.Println(cmt.String())
}

func TestProposerJoinAndLeave(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	pub1, _ := ts.RandBLSKeyPair()
	pub2, _ := ts.RandBLSKeyPair()
	pub3, _ := ts.RandBLSKeyPair()
	pub4, _ := ts.RandBLSKeyPair()
	pub5, _ := ts.RandBLSKeyPair()
	pub6, _ := ts.RandBLSKeyPair()
	pub7, _ := ts.RandBLSKeyPair()
	pub8, _ := ts.RandBLSKeyPair()
	pub9, _ := ts.RandBLSKeyPair()
	pubA, _ := ts.RandBLSKeyPair()
	pubB, _ := ts.RandBLSKeyPair()
	pubC, _ := ts.RandBLSKeyPair()

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

	cmt, err := committee.NewCommittee(
		[]*validator.Validator{val1, val2, val3, val4, val5, val6, val7}, 7, val1.Address())
	assert.NoError(t, err)
	fmt.Println(cmt.String())

	// This code comment explains how the committee changes when new validators join.
	//
	// The symbols used in the explanation are as follows:
	// * represents the current proposer
	// ! represents a failed proposer
	// $ represents a joined validator
	// h is height and r is round number
	//
	// Initially, the committee consists of validators numbered from 1 to 7.
	// Validator 1 is the oldest and the current proposer.
	// The committee configuration is represented as:
	// +*+-+-+-+-+-+-+
	// |1|2|3|4|5|6|7|
	// +-+-+-+-+-+-+-+
	//
	// When new validators join, they are inserted before the current proposer.
	// For example, validator 8 joins the committee:
	// +$+*+-+-+-+-+-+-+
	// |8|1|2|3|4|5|6|7|
	// +-+-+-+-+-+-+-+-+
	//
	// After the addition of a new validator, the committee needs to be adjusted,
	// and the oldest validator should leave:
	// +$+-+-+-+-+-+-+
	// |8|2|3|4|5|6|7|
	// +-+-+-+-+-+-+-+
	//
	// Next, we move to the next proposer.
	// +-+*+-+-+-+-+-+
	// |8|2|3|4|5|6|7|
	// +-+-+-+-+-+-+-+
	//
	//-------------------------------------
	// In this test, we cover the following movements:
	//
	// h=1000, r=0        h=1001, r=0        h=1002, r=3        h=1003, r=0
	// +*+-+-+-+-+-+-+    +$+*+-+-+-+-+-+    +-+-+!+!+!+*+-+    +-+$+-+-+$+-+*+
	// |1|2|3|4|5|6|7| => |8|2|3|4|5|6|7| => |8|2|3|4|5|6|7| => |8|2|3|5|9|6|7| =>
	// +-+-+-+-+-+-+-+    +-+-+-+-+-+-+-+    +-+-+$+-+-+-+-+    +-+-+-+-+-+-+-+
	//
	// h=1004, r=1        h=1005, r=0        h=1006, r=2        h=1007, r=4
	// +!+*+-+-+-+$+-+    +-+$+$+-+*+-+-+    +*+-+-+$+-+!+!+    +*+$+-+!+!+!+!+
	// |8|2|3|9|6|A|7| => |8|B|C|2|3|9|A| => |B|C|2|1|3|9|A| => |5|6|B|C|2|1|3|
	// +-+-+-+-+-+-+-+    +-+-+-+-+-+-+-+    +-+-+-+-+-+-+-+    +$+-+-+-+$+-+$+

	// Height 1001
	val8.UpdateLastSortitionHeight(1001)
	cmt.Update(0, []*validator.Validator{val8})
	assert.Equal(t, cmt.Proposer(0).Number(), int32(2))
	assert.Equal(t, cmt.Proposer(1).Number(), int32(3))
	assert.Equal(t, cmt.Proposer(2).Number(), int32(4))
	assert.Equal(t, cmt.Committers(), []int32{8, 2, 3, 4, 5, 6, 7})
	fmt.Println(cmt.String())

	// Height 1002
	val3.UpdateLastSortitionHeight(1002)
	cmt.Update(3, []*validator.Validator{val3})
	assert.Equal(t, cmt.Proposer(0).Number(), int32(6))
	fmt.Println(cmt.String())

	// Height 1003
	val2.UpdateLastSortitionHeight(1003)
	val9.UpdateLastSortitionHeight(1003)
	cmt.Update(0, []*validator.Validator{val9, val2})
	assert.Equal(t, cmt.Proposer(0).Number(), int32(7))
	assert.Equal(t, cmt.Proposer(1).Number(), int32(8))
	assert.Equal(t, cmt.Committers(), []int32{8, 2, 3, 5, 9, 6, 7})
	fmt.Println(cmt.String())

	// Height 1004
	valA.UpdateLastSortitionHeight(1004)
	cmt.Update(1, []*validator.Validator{valA})
	assert.Equal(t, cmt.Proposer(0).Number(), int32(2))
	assert.Equal(t, cmt.Committers(), []int32{8, 2, 3, 9, 6, 10, 7})
	fmt.Println(cmt.String())

	// Height 1005
	valB.UpdateLastSortitionHeight(1005)
	valC.UpdateLastSortitionHeight(1005)
	cmt.Update(0, []*validator.Validator{valC, valB})
	assert.Equal(t, cmt.Proposer(0).Number(), int32(3))
	assert.Equal(t, cmt.Proposer(1).Number(), int32(9))
	assert.Equal(t, cmt.Proposer(2).Number(), int32(10))
	assert.Equal(t, cmt.Committers(), []int32{8, 11, 12, 2, 3, 9, 10})
	fmt.Println(cmt.String())

	// Height 1006
	val1.UpdateLastSortitionHeight(1006)
	cmt.Update(2, []*validator.Validator{val1})
	assert.Equal(t, cmt.Proposer(0).Number(), int32(11))
	assert.Equal(t, cmt.Committers(), []int32{11, 12, 2, 1, 3, 9, 10})
	fmt.Println(cmt.String())

	// Height 1007
	val2.UpdateLastSortitionHeight(1007)
	val3.UpdateLastSortitionHeight(1007)
	val5.UpdateLastSortitionHeight(1007)
	val6.UpdateLastSortitionHeight(1007)
	cmt.Update(4, []*validator.Validator{val2, val3, val5, val6})
	assert.Equal(t, cmt.Proposer(0).Number(), int32(5))
	assert.Equal(t, cmt.Committers(), []int32{5, 6, 11, 12, 2, 1, 3})
	fmt.Println(cmt.String())
}

func TestIsProposer(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	val1, _ := ts.GenerateTestValidator(0)
	val2, _ := ts.GenerateTestValidator(1)
	val3, _ := ts.GenerateTestValidator(2)
	val4, _ := ts.GenerateTestValidator(3)

	cmt, err := committee.NewCommittee([]*validator.Validator{val1, val2, val3, val4}, 4, val1.Address())
	assert.NoError(t, err)

	assert.Equal(t, cmt.Proposer(0).Number(), int32(0))
	assert.Equal(t, cmt.Proposer(1).Number(), int32(1))
	assert.True(t, cmt.IsProposer(val3.Address(), 2))
	assert.False(t, cmt.IsProposer(val4.Address(), 2))
	assert.False(t, cmt.IsProposer(ts.RandAccAddress(), 2))
	assert.Equal(t, cmt.Validators(), []*validator.Validator{val1, val2, val3, val4})
}

func TestCommitters(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	val1, _ := ts.GenerateTestValidator(0)
	val2, _ := ts.GenerateTestValidator(1)
	val3, _ := ts.GenerateTestValidator(2)
	val4, _ := ts.GenerateTestValidator(3)

	cmt, err := committee.NewCommittee([]*validator.Validator{val1, val2, val3, val4}, 4, val1.Address())
	assert.NoError(t, err)
	assert.Equal(t, cmt.Committers(), []int32{0, 1, 2, 3})
}

func TestSortJoined(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	val1, _ := ts.GenerateTestValidator(0)
	val2, _ := ts.GenerateTestValidator(1)
	val3, _ := ts.GenerateTestValidator(2)
	val4, _ := ts.GenerateTestValidator(3)
	val5, _ := ts.GenerateTestValidator(4)
	val6, _ := ts.GenerateTestValidator(5)
	val7, _ := ts.GenerateTestValidator(6)

	cmt, err := committee.NewCommittee([]*validator.Validator{val1, val2, val3, val4}, 17, val1.Address())
	assert.NoError(t, err)
	committee2, err := committee.NewCommittee([]*validator.Validator{val1, val2, val3, val4}, 17, val1.Address())
	assert.NoError(t, err)

	cmt.Update(0, []*validator.Validator{val5, val6, val7})
	committee2.Update(0, []*validator.Validator{val7, val5, val6})
}

func TestTotalPower(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	pub, _ := ts.RandBLSKeyPair()
	val0 := validator.NewValidator(pub, 0) // Bootstrap validator
	val1, _ := ts.GenerateTestValidator(0)
	val2, _ := ts.GenerateTestValidator(1)
	val3, _ := ts.GenerateTestValidator(2)
	val4, _ := ts.GenerateTestValidator(3)

	cmt, err := committee.NewCommittee([]*validator.Validator{val0, val1, val2, val3, val4}, 4, val1.Address())
	assert.NoError(t, err)

	totalPower := val0.Power() + val1.Power() + val2.Power() + val3.Power() + val4.Power()
	totalStake := val0.Stake() + val1.Stake() + val2.Stake() + val3.Stake() + val4.Stake()
	assert.Equal(t, cmt.TotalPower(), totalPower)
	assert.Equal(t, cmt.TotalPower(), int64(totalStake+1))
}
