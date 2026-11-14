package vote

import (
	"testing"

	"github.com/fxamacker/cbor/v2"
	"github.com/stretchr/testify/require"
)

// TestCPVoteUnmarshalCBORRejectsTags checks that the change-proposer part of a vote,
// and its justifications, are decoded with the decoder of the data from the network.
func TestCPVoteUnmarshalCBORRejectsTags(t *testing.T) {
	tag := []byte{0xd9, 0xd9, 0xf7} // the self-describe tag 55799, which prefixes any item

	justData, err := cbor.Marshal(&JustInitYes{})
	require.NoError(t, err)

	encode := func(cpv *_cpVote) []byte {
		t.Helper()

		data, err := cbor.Marshal(cpv)
		require.NoError(t, err)

		return data
	}
	encodeConflict := func(just0, just1 []byte) []byte {
		t.Helper()

		data, err := cbor.Marshal(&_JustMainVoteConflict{
			Just0Type: JustTypeInitYes,
			Just0Data: just0,
			Just1Type: JustTypeInitYes,
			Just1Data: just1,
		})
		require.NoError(t, err)

		return data
	}

	valid := encode(&_cpVote{Round: 1, Value: CPValueNo, JustType: JustTypeInitYes, JustData: justData})
	require.NoError(t, new(cpVote).UnmarshalCBOR(valid))
	require.Error(t, new(cpVote).UnmarshalCBOR(append(tag, valid...)), "change-proposer vote")

	taggedJust := encode(&_cpVote{
		Round: 1, Value: CPValueNo, JustType: JustTypeInitYes,
		JustData: append(tag, justData...),
	})
	require.Error(t, new(cpVote).UnmarshalCBOR(taggedJust), "justification")

	validConflict := encode(&_cpVote{
		Round: 1, Value: CPValueAbstain, JustType: JustTypeMainVoteConflict,
		JustData: encodeConflict(justData, justData),
	})
	require.NoError(t, new(cpVote).UnmarshalCBOR(validConflict))

	for _, conflict := range [][]byte{
		append(tag, encodeConflict(justData, justData)...),
		encodeConflict(append(tag, justData...), justData),
		encodeConflict(justData, append(tag, justData...)),
	} {
		data := encode(&_cpVote{
			Round: 1, Value: CPValueAbstain, JustType: JustTypeMainVoteConflict,
			JustData: conflict,
		})
		require.Error(t, new(cpVote).UnmarshalCBOR(data), "conflicting justifications")
	}
}
