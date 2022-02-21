package errors

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCode(t *testing.T) {
	err1 := Error(ErrInsufficientFunds)
	assert.Equal(t, Code(err1), ErrInsufficientFunds)

	err2 := fmt.Errorf("Nope")
	assert.Equal(t, Code(err2), ErrGeneric)

	assert.Equal(t, Code(nil), ErrNone)
}

func TestMessages(t *testing.T) {
	for i := 0; i < ErrCount; i++ {
		assert.NotEmpty(t, messages[i], "Error code %v", i)
	}
}

func TestErrorCode(t *testing.T) {
	err1 := Error(ErrInsufficientFunds)
	err2 := Errorf(ErrInvalidTx, err1.Error())
	err3 := Errorf(ErrInvalidBlock, err1.Error())

	assert.Equal(t, Code(err1), ErrInsufficientFunds)
	assert.Equal(t, Code(err2), ErrInvalidTx)
	assert.Equal(t, Code(err3), ErrInvalidBlock)
	assert.Equal(t, "Insufficient funds", err1.Error())
	assert.Equal(t, "Invalid transaction: Insufficient funds", err2.Error())
	assert.Equal(t, "Invalid block: Insufficient funds", err3.Error())
}
