package payload

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnbondType(t *testing.T) {
	pld := UnbondPayload{}
	assert.Equal(t, TypeUnbond, pld.Type())
}

func TestUnbondString(t *testing.T) {
	pld := UnbondPayload{}
	assert.Contains(t, pld.String(), "{Unbond ")
}
