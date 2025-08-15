package committee_test

import (
	"fmt"
	"testing"

	"github.com/pactus-project/pactus/committee"
	"github.com/pactus-project/pactus/types/protocol"
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

	assert.Equal(t, int32(0), cmt.Proposer(0).Number())
	assert.Equal(t, int32(3), cmt.Proposer(3).Number())
	assert.Equal(t, int32(0), cmt.Proposer(4).Number())

	cmt.Update(0, nil)
	assert.Equal(t, int32(1), cmt.Proposer(0).Number())
}

func TestInvalidProposerJoinAndLeave(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	val1 := ts.GenerateTestValidator(testsuite.ValidatorWithNumber(0))
	val2 := ts.GenerateTestValidator(testsuite.ValidatorWithNumber(1))
	val3 := ts.GenerateTestValidator(testsuite.ValidatorWithNumber(2))
	val4 := ts.GenerateTestValidator(testsuite.ValidatorWithNumber(3))
	val5 := ts.GenerateTestValidator(testsuite.ValidatorWithNumber(4))

	cmt, err := committee.NewCommittee([]*validator.Validator{val1, val2, val3, val4}, 4, val5.Address())
	assert.Error(t, err)
	assert.Nil(t, cmt)
}

func TestProposerMove(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	val1 := ts.GenerateTestValidator(testsuite.ValidatorWithNumber(1))
	val2 := ts.GenerateTestValidator(testsuite.ValidatorWithNumber(2))
	val3 := ts.GenerateTestValidator(testsuite.ValidatorWithNumber(3))
	val4 := ts.GenerateTestValidator(testsuite.ValidatorWithNumber(4))
	val5 := ts.GenerateTestValidator(testsuite.ValidatorWithNumber(5))
	val6 := ts.GenerateTestValidator(testsuite.ValidatorWithNumber(6))
	val7 := ts.GenerateTestValidator(testsuite.ValidatorWithNumber(7))

	cmt, err := committee.NewCommittee(
		[]*validator.Validator{val1, val2, val3, val4, val5, val6, val7}, 7, val1.Address())
	assert.NoError(t, err)

	//
	// +*+-+-+-+-+-+-+    +-+*+-+-+-+-+-+    +-+-+-+-+-+*+-+    +*+-+-+-+-+-+-+
	// |1|2|3|4|5|6|7| => |1|2|3|4|5|6|7| => |1|2|3|4|5|6|7| => |1|2|3|4|5|6|7|
	// +-+-+-+-+-+-+-+    +-+-+-+-+-+-+-+    +-+-+-+-+-+-+-+    +-+-+-+-+-+-+-+
	//

	assert.Equal(t, int32(1), cmt.Proposer(0).Number())
	assert.Equal(t, int32(1), cmt.Proposer(7).Number())
	assert.Equal(t, []*validator.Validator{val1, val2, val3, val4, val5, val6, val7}, cmt.Validators())
	fmt.Println(cmt.String())

	// Height 1001
	cmt.Update(0, nil)
	assert.Equal(t, int32(2), cmt.Proposer(0).Number())
	assert.Equal(t, int32(3), cmt.Proposer(1).Number())
	assert.Equal(t, int32(2), cmt.Proposer(7).Number())
	assert.Equal(t, []*validator.Validator{val1, val2, val3, val4, val5, val6, val7}, cmt.Validators())
	fmt.Println(cmt.String())

	// Height 1002
	cmt.Update(3, nil)
	assert.Equal(t, int32(6), cmt.Proposer(0).Number())
	assert.Equal(t, []*validator.Validator{val1, val2, val3, val4, val5, val6, val7}, cmt.Validators())
	fmt.Println(cmt.String())

	// Height 1003
	cmt.Update(1, nil)
	assert.Equal(t, int32(1), cmt.Proposer(0).Number())
	assert.Equal(t, []*validator.Validator{val1, val2, val3, val4, val5, val6, val7}, cmt.Validators())
	fmt.Println(cmt.String())
}

func TestValidatorConsistency(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	val1 := ts.GenerateTestValidator(testsuite.ValidatorWithNumber(1))
	val2 := ts.GenerateTestValidator(testsuite.ValidatorWithNumber(2))
	val3 := ts.GenerateTestValidator(testsuite.ValidatorWithNumber(3))
	val4 := ts.GenerateTestValidator(testsuite.ValidatorWithNumber(4))

	cmt, _ := committee.NewCommittee([]*validator.Validator{val1, val2, val3, val4}, 4, val1.Address())

	t.Run("Updating validators' stake, Should panic", func(t *testing.T) {
		assert.Panics(t, func() {
			val1.AddToStake(1)
			cmt.Update(0, []*validator.Validator{val1})
		}, "The code did not panic")
	})
}

func TestProposerJoin(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	val1 := ts.GenerateTestValidator(testsuite.ValidatorWithNumber(1))
	val2 := ts.GenerateTestValidator(testsuite.ValidatorWithNumber(2))
	val3 := ts.GenerateTestValidator(testsuite.ValidatorWithNumber(3))
	val4 := ts.GenerateTestValidator(testsuite.ValidatorWithNumber(4))
	val5 := ts.GenerateTestValidator(testsuite.ValidatorWithNumber(5))
	val6 := ts.GenerateTestValidator(testsuite.ValidatorWithNumber(6))
	val7 := ts.GenerateTestValidator(testsuite.ValidatorWithNumber(7))

	cmt, err := committee.NewCommittee([]*validator.Validator{val1, val2, val3, val4}, 7, val1.Address())
	assert.NoError(t, err)
	assert.Equal(t, 4, cmt.Size())

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
	assert.Equal(t, int32(2), cmt.Proposer(0).Number())
	assert.Equal(t, []int32{1, 2, 3, 4}, cmt.Committers())
	assert.Equal(t, 4, cmt.Size())
	fmt.Println(cmt.String())

	// Height 1001
	val5.UpdateLastSortitionHeight(1001)
	cmt.Update(0, []*validator.Validator{val5})
	assert.Equal(t, int32(3), cmt.Proposer(0).Number())
	assert.Equal(t, int32(4), cmt.Proposer(1).Number())
	assert.Equal(t, []int32{1, 5, 2, 3, 4}, cmt.Committers())
	assert.Equal(t, 5, cmt.Size())
	fmt.Println(cmt.String())

	// Height 1002
	cmt.Update(1, nil)
	assert.Equal(t, int32(1), cmt.Proposer(0).Number())
	assert.Equal(t, int32(5), cmt.Proposer(1).Number())
	fmt.Println(cmt.String())

	// Height 1003
	val3.UpdateLastSortitionHeight(1003)
	val6.UpdateLastSortitionHeight(1003)
	val7.UpdateLastSortitionHeight(1003)
	cmt.Update(1, []*validator.Validator{val7, val3, val6})
	assert.Equal(t, int32(2), cmt.Proposer(0).Number())
	assert.Equal(t, []int32{6, 7, 1, 5, 2, 3, 4}, cmt.Committers())
	assert.Equal(t, 7, cmt.Size())
	fmt.Println(cmt.String())

	// Height 1004
	cmt.Update(0, nil)
	assert.Equal(t, int32(3), cmt.Proposer(0).Number())
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
	// -------------------------------------
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
	assert.Equal(t, int32(2), cmt.Proposer(0).Number())
	assert.Equal(t, int32(3), cmt.Proposer(1).Number())
	assert.Equal(t, int32(4), cmt.Proposer(2).Number())
	assert.Equal(t, []int32{8, 2, 3, 4, 5, 6, 7}, cmt.Committers())
	fmt.Println(cmt.String())

	// Height 1002
	val3.UpdateLastSortitionHeight(1002)
	cmt.Update(3, []*validator.Validator{val3})
	assert.Equal(t, int32(6), cmt.Proposer(0).Number())
	fmt.Println(cmt.String())

	// Height 1003
	val2.UpdateLastSortitionHeight(1003)
	val9.UpdateLastSortitionHeight(1003)
	cmt.Update(0, []*validator.Validator{val9, val2})
	assert.Equal(t, int32(7), cmt.Proposer(0).Number())
	assert.Equal(t, int32(8), cmt.Proposer(1).Number())
	assert.Equal(t, []int32{8, 2, 3, 5, 9, 6, 7}, cmt.Committers())
	fmt.Println(cmt.String())

	// Height 1004
	valA.UpdateLastSortitionHeight(1004)
	cmt.Update(1, []*validator.Validator{valA})
	assert.Equal(t, int32(2), cmt.Proposer(0).Number())
	assert.Equal(t, []int32{8, 2, 3, 9, 6, 10, 7}, cmt.Committers())
	fmt.Println(cmt.String())

	// Height 1005
	valB.UpdateLastSortitionHeight(1005)
	valC.UpdateLastSortitionHeight(1005)
	cmt.Update(0, []*validator.Validator{valC, valB})
	assert.Equal(t, int32(3), cmt.Proposer(0).Number())
	assert.Equal(t, int32(9), cmt.Proposer(1).Number())
	assert.Equal(t, int32(10), cmt.Proposer(2).Number())
	assert.Equal(t, []int32{8, 11, 12, 2, 3, 9, 10}, cmt.Committers())
	fmt.Println(cmt.String())

	// Height 1006
	val1.UpdateLastSortitionHeight(1006)
	cmt.Update(2, []*validator.Validator{val1})
	assert.Equal(t, int32(11), cmt.Proposer(0).Number())
	assert.Equal(t, []int32{11, 12, 2, 1, 3, 9, 10}, cmt.Committers())
	fmt.Println(cmt.String())

	// Height 1007
	val2.UpdateLastSortitionHeight(1007)
	val3.UpdateLastSortitionHeight(1007)
	val5.UpdateLastSortitionHeight(1007)
	val6.UpdateLastSortitionHeight(1007)
	cmt.Update(4, []*validator.Validator{val2, val3, val5, val6})
	assert.Equal(t, int32(5), cmt.Proposer(0).Number())
	assert.Equal(t, []int32{5, 6, 11, 12, 2, 1, 3}, cmt.Committers())
	fmt.Println(cmt.String())
}

func TestIsProposer(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	val1 := ts.GenerateTestValidator(testsuite.ValidatorWithNumber(0))
	val2 := ts.GenerateTestValidator(testsuite.ValidatorWithNumber(1))
	val3 := ts.GenerateTestValidator(testsuite.ValidatorWithNumber(2))
	val4 := ts.GenerateTestValidator(testsuite.ValidatorWithNumber(3))

	cmt, err := committee.NewCommittee([]*validator.Validator{val1, val2, val3, val4}, 4, val1.Address())
	assert.NoError(t, err)

	assert.Equal(t, int32(0), cmt.Proposer(0).Number())
	assert.Equal(t, int32(1), cmt.Proposer(1).Number())
	assert.True(t, cmt.IsProposer(val3.Address(), 2))
	assert.False(t, cmt.IsProposer(val4.Address(), 2))
	assert.False(t, cmt.IsProposer(ts.RandAccAddress(), 2))
	assert.Equal(t, []*validator.Validator{val1, val2, val3, val4}, cmt.Validators())
}

func TestCommitters(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	val1 := ts.GenerateTestValidator(testsuite.ValidatorWithNumber(0))
	val2 := ts.GenerateTestValidator(testsuite.ValidatorWithNumber(1))
	val3 := ts.GenerateTestValidator(testsuite.ValidatorWithNumber(2))
	val4 := ts.GenerateTestValidator(testsuite.ValidatorWithNumber(3))

	cmt, err := committee.NewCommittee([]*validator.Validator{val1, val2, val3, val4}, 4, val1.Address())
	assert.NoError(t, err)
	assert.Equal(t, []int32{0, 1, 2, 3}, cmt.Committers())
}

func TestSortJoined(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	val1 := ts.GenerateTestValidator(testsuite.ValidatorWithNumber(0))
	val2 := ts.GenerateTestValidator(testsuite.ValidatorWithNumber(1))
	val3 := ts.GenerateTestValidator(testsuite.ValidatorWithNumber(2))
	val4 := ts.GenerateTestValidator(testsuite.ValidatorWithNumber(3))
	val5 := ts.GenerateTestValidator(testsuite.ValidatorWithNumber(4))
	val6 := ts.GenerateTestValidator(testsuite.ValidatorWithNumber(5))
	val7 := ts.GenerateTestValidator(testsuite.ValidatorWithNumber(6))

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
	val1 := ts.GenerateTestValidator(testsuite.ValidatorWithNumber(0))
	val2 := ts.GenerateTestValidator(testsuite.ValidatorWithNumber(1))
	val3 := ts.GenerateTestValidator(testsuite.ValidatorWithNumber(2))
	val4 := ts.GenerateTestValidator(testsuite.ValidatorWithNumber(3))

	cmt, err := committee.NewCommittee([]*validator.Validator{val0, val1, val2, val3, val4}, 4, val1.Address())
	assert.NoError(t, err)

	totalPower := val0.Power() + val1.Power() + val2.Power() + val3.Power() + val4.Power()
	totalStake := val0.Stake() + val1.Stake() + val2.Stake() + val3.Stake() + val4.Stake()
	assert.Equal(t, totalPower, cmt.TotalPower())
	assert.Equal(t, int64(totalStake+1), cmt.TotalPower())
}

func TestProtocolVersionPercentages(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	val1 := ts.GenerateTestValidator()
	val2 := ts.GenerateTestValidator()
	val3 := ts.GenerateTestValidator()
	val4 := ts.GenerateTestValidator()
	val5 := ts.GenerateTestValidator()

	val2.UpdateProtocolVersion(protocol.ProtocolVersion1)
	val3.UpdateProtocolVersion(protocol.ProtocolVersion2)
	val4.UpdateProtocolVersion(protocol.ProtocolVersion2)
	val5.UpdateProtocolVersion(protocol.ProtocolVersion2)

	cmt, err := committee.NewCommittee([]*validator.Validator{val1, val2, val3, val4, val5}, 5, val1.Address())
	assert.NoError(t, err)

	percentages := cmt.ProtocolVersions()
	assert.Equal(t, float64(20), percentages[protocol.ProtocolVersionUnknown])
	assert.Equal(t, float64(20), percentages[protocol.ProtocolVersion1])
	assert.Equal(t, float64(60), percentages[protocol.ProtocolVersion2])
}
