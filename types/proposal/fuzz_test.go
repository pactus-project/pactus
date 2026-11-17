package proposal_test

import (
	"testing"

	"github.com/pactus-project/pactus/types/proposal"
)

// A proposal comes from any peer: decoding and checking it never panics.
func FuzzProposalDecode(f *testing.F) {
	// The seeds are in testdata/fuzz/FuzzProposalDecode: valid proposals.
	f.Add([]byte{})

	f.Fuzz(func(_ *testing.T, data []byte) {
		prop := new(proposal.Proposal)
		if err := prop.UnmarshalCBOR(data); err != nil {
			return
		}
		if err := prop.BasicCheck(); err != nil {
			return
		}

		_ = prop.LogString()
		_ = prop.Hash()
		_ = prop.SignBytes()
	})
}
