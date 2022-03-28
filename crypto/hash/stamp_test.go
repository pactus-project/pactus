package hash

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStampFromString(t *testing.T) {
	stamp1 := Stamp{1, 2, 3, 4}
	stamp2, err := StampFromString(stamp1.String())
	assert.NoError(t, err)
	assert.True(t, stamp1.EqualsTo(stamp2))
	assert.Equal(t, stamp1.Bytes(), []byte{1, 2, 3, 4})

	js, _ := stamp1.MarshalJSON()
	assert.Contains(t, string(js), stamp1.String())

	_, err = StampFromString("")
	assert.Error(t, err)

	_, err = StampFromString("inv")
	assert.Error(t, err)

	_, err = StampFromString("00")
	assert.Error(t, err)
}
