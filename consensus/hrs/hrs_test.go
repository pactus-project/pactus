package hrs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOperator(t *testing.T) {
	hrs1 := NewHRS(100, 1, 1)
	hrs2 := NewHRS(100, 1, 2) // bigger than hrs1
	hrs3 := NewHRS(100, 2, 1) // bigger than hrs2
	hrs4 := NewHRS(101, 0, 0) // bigger than hrs3
	hrs5 := new(HRS)
	hrs6 := new(HRS)

	assert.True(t, hrs1.LessThan(hrs2))
	assert.True(t, hrs1.LessThan(hrs3))
	assert.True(t, hrs1.LessThan(hrs4))

	assert.True(t, hrs4.GreaterThan(hrs1))
	assert.True(t, hrs4.GreaterThan(hrs2))
	assert.True(t, hrs4.GreaterThan(hrs3))

	hrs5.UpdateHeight(101)
	hrs5.UpdateRoundStep(0, 0)

	assert.True(t, hrs5.EqualsTo(hrs4))
	assert.False(t, hrs5.EqualsTo(hrs1))
	assert.False(t, hrs5.EqualsTo(hrs2))
	assert.False(t, hrs5.EqualsTo(hrs3))

	hrs6.UpdateHeightRoundStep(hrs5.Height(), hrs5.Round(), hrs5.Step())
	assert.True(t, hrs6.EqualsTo(*hrs5))
}

func TestMarshaling(t *testing.T) {
	hrs1 := NewHRS(100, 1, 1)
	hrs2 := new(HRS)

	bs, err := hrs1.MarshalCBOR()
	assert.NoError(t, err)
	assert.NoError(t, hrs2.UnmarshalCBOR(bs))
	assert.Equal(t, hrs1, *hrs2)
}
