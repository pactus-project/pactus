package state

import (
	"testing"

	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/types/protocol"
	"github.com/pactus-project/pactus/types/vote"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
)

func TestErrorMessages(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	rootHash1 := ts.RandHash()
	rootHash2 := ts.RandHash()
	addr1 := ts.RandAccAddress()
	addr2 := ts.RandAccAddress()
	vte := vote.NewPrepareVote(hash.UndefHash, 1, 0, ts.RandValAddress())

	tests := []struct {
		err      error
		contains []string
	}{
		{
			InvalidBlockVersionError{Version: protocol.ProtocolVersion3},
			[]string{"invalid block version", protocol.ProtocolVersion3.String()},
		},
		{
			InvalidSubsidyAmountError{Expected: 1_000_000_000, Got: 2_000_000_000},
			[]string{"invalid subsidy amount", "1", "2"},
		},
		{
			InvalidVoteForCertificateError{Vote: vte},
			[]string{"invalid vote to update the last certificate", vte.Type().String()},
		},
		{
			InvalidStateRootHashError{Expected: rootHash1, Got: rootHash2},
			[]string{"invalid state root hash", rootHash1.String(), rootHash2.String()},
		},
		{
			InvalidProposerError{Expected: addr1, Got: addr2},
			[]string{"invalid block proposer", addr1.String(), addr2.String()},
		},
		{
			InvalidBlockTimeError{Reason: "too far in the future"},
			[]string{"invalid block time", "too far in the future"},
		},
	}

	for _, tt := range tests {
		for _, part := range tt.contains {
			assert.Contains(t, tt.err.Error(), part)
		}
	}
}
