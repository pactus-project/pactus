package vote_test

import (
	"testing"

	"github.com/pactus-project/pactus/types/vote"
	"github.com/stretchr/testify/require"
)

// A vote comes from any peer: decoding and checking it never panics,
// and a vote that is decoded keeps its hash once it is encoded and decoded again.
func FuzzVoteDecode(f *testing.F) {
	// The seeds are in testdata/fuzz/FuzzVoteDecode: one valid vote of each type.
	f.Add([]byte{})

	f.Fuzz(func(t *testing.T, data []byte) {
		vte := new(vote.Vote)
		if err := vte.UnmarshalCBOR(data); err != nil {
			return
		}

		// A vote is used only once it is checked.
		if err := vte.BasicCheck(); err != nil {
			return
		}

		_ = vte.LogString()
		_ = vte.IsCPVote()
		_ = vte.SignBytes()

		encoded, err := vte.MarshalCBOR()
		require.NoError(t, err)

		again := new(vote.Vote)
		require.NoError(t, again.UnmarshalCBOR(encoded))
		require.Equal(t, vte.Hash(), again.Hash())
	})
}
