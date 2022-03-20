package hash

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStampFromString(t *testing.T) {
	stamp1 := GenerateTestStamp()
	stamp2, err := StampFromString(stamp1.String())
	assert.NoError(t, err)
	assert.True(t, stamp1.EqualsTo(stamp2))

	js, _ := stamp1.MarshalJSON()
	assert.Contains(t, string(js), stamp1.String())

	_, err = StampFromString("")
	assert.Error(t, err)

	_, err = StampFromString("inv")
	assert.Error(t, err)

	_, err = StampFromString("00")
	assert.Error(t, err)
}
