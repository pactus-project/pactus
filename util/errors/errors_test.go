package errors

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCode(t *testing.T) {
	err1 := Error(ErrInvalidAmount)
	assert.Equal(t, ErrInvalidAmount, Code(err1))

	err2 := fmt.Errorf("Nope")
	assert.Equal(t, ErrGeneric, Code(err2))

	assert.Equal(t, ErrNone, Code(nil))
}

func TestMessages(t *testing.T) {
	for i := 0; i < ErrCount; i++ {
		assert.NotEmpty(t, messages[i], "Error code %v", i)
	}
}

func TestErrorCode(t *testing.T) {
	err1 := Error(ErrInvalidAmount)
	err2 := Errorf(ErrInvalidTx, err1.Error())
	err3 := Errorf(ErrInvalidBlock, err1.Error())

	assert.Equal(t, ErrInvalidTx, Code(err2))
	assert.Equal(t, ErrInvalidBlock, Code(err3))
	assert.Equal(t, "invalid amount", err1.Error())
	assert.Equal(t, "invalid transaction: invalid amount", err2.Error())
	assert.Equal(t, "invalid block: invalid amount", err3.Error())
}
