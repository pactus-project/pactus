package errors

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCode(t *testing.T) {
	err2 := fmt.Errorf("Nope")
	assert.Equal(t, ErrGeneric, Code(err2))

	assert.Equal(t, ErrNone, Code(nil))
}

func TestMessages(t *testing.T) {
	for i := 0; i < ErrCount; i++ {
		assert.NotEmpty(t, messages[i], "Error code %v", i)
	}
}
