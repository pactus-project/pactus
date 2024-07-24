package message

import (
	"testing"

	"github.com/pactus-project/pactus/util/errors"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
)

func TestQueryProposalType(t *testing.T) {
	m := &QueryProposalMessage{}
	assert.Equal(t, TypeQueryProposal, m.Type())
}

func TestQueryProposalMessage(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	t.Run("Invalid round", func(t *testing.T) {
		m := NewQueryProposalMessage(0, -1, ts.RandValAddress())

		assert.Equal(t, errors.ErrInvalidRound, errors.Code(m.BasicCheck()))
	})

	t.Run("OK", func(t *testing.T) {
		m := NewQueryProposalMessage(100, 0, ts.RandValAddress())

		assert.NoError(t, m.BasicCheck())
		assert.Contains(t, m.String(), "100")
	})
}
