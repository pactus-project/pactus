package message

import (
	"testing"

	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
)

func TestQueryProposalType(t *testing.T) {
	m := &QueryProposalMessage{}
	assert.Equal(t, m.Type(), TypeQueryProposal)
}

func TestQueryProposalMessage(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	m := NewQueryProposalMessage(0, ts.RandValAddress())
	assert.NoError(t, m.BasicCheck())
}
