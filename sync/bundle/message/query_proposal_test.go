package message

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQueryProposalType(t *testing.T) {
	m := &QueryProposalMessage{}
	assert.Equal(t, m.Type(), TypeQueryProposal)
}

func TestQueryProposalMessage(t *testing.T) {
	m := NewQueryProposalMessage(0)
	assert.NoError(t, m.BasicCheck())
}
