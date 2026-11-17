package state

import (
	"testing"

	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
)

func TestChainErrorMessages(t *testing.T) {
	ts := testsuite.NewTestSuite(t)
	expected, got := ts.RandHash(), ts.RandHash()

	assert.Equal(t,
		"a possible fork is detected, the last block hash is "+expected.String()+", but the block follows "+got.String(),
		ForkDetectedError{Expected: expected, Got: got}.Error())
	assert.Equal(t, "block 12 is ahead of the chain, the last block is 10",
		FutureBlockError{Height: 12, LastHeight: 10}.Error())
}
