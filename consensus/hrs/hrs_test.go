package hrs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOperator(t *testing.T) {
	hrs1 := NewHRS(100, 10, 5)
	hrs2 := NewHRS(100, 10, 6)
	hrs3 := NewHRS(100, 11, 5)
	hrs4 := NewHRS(101, 10, 5)
	hrs5 := NewHRS(100, 10, 5)

	assert.True(t, hrs1.LessThan(hrs2))
	assert.True(t, hrs1.LessThan(hrs3))
	assert.True(t, hrs1.LessThan(hrs4))
	assert.True(t, hrs1.EqualsTo(hrs5))
	assert.False(t, hrs1.LessThan(hrs5))
}
