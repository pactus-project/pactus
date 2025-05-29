package payload

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithdrawType(t *testing.T) {
	pld := WithdrawPayload{}
	assert.Equal(t, TypeWithdraw, pld.Type())
}

func TestWithdrawString(t *testing.T) {
	pld := WithdrawPayload{}
	assert.Contains(t, pld.String(), "{Withdraw ")
}
