package bundle

import (
	"bytes"
	"testing"

	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Bundle.Decode calls a CBOR decoder that stops after reading
// one item; it never checks that the reader is exhausted afterward. A bundle
// with arbitrary bytes appended after a valid CBOR item still decodes to the
// exact same logical message.
//
// This matters for gossip: network.MessageIDFunc (network/utils.go) computes
// the pubsub message ID as hash.CalcHash(rawBytes)[:20] — a hash of the RAW
// WIRE BYTES, not of the decoded message. Two wire messages that differ only
// by trailing junk decode to the identical Bundle, but get different gossip
// IDs. An attacker can therefore make libp2p-pubsub's own deduplication (keyed
// by that ID) treat the same logical message as "new" on every relay by
// varying the junk suffix, bypassing the de-duplication that would otherwise
// make the rest of the mesh ignore a message it already processed.
//
// This test expresses the desired behavior (trailing bytes must be rejected)
// and fails today.
func TestDecodeAcceptsTrailingBytesBreakingGossipDedup(t *testing.T) {
	msg := message.NewHelloAckMessage(message.ResponseCodeOK, "ok", 1)
	bdl := NewBundle(msg)

	clean, err := bdl.Encode()
	require.NoError(t, err)

	junked := append(bytes.Clone(clean), 0xDE, 0xAD, 0xBE, 0xEF, 0x00, 0x01, 0x02, 0x03)

	decodedClean := new(Bundle)
	_, err = decodedClean.Decode(bytes.NewReader(clean))
	require.NoError(t, err)

	decodedJunked := new(Bundle)
	bytesRead, err := decodedJunked.Decode(bytes.NewReader(junked))

	// Desired behavior: the extra bytes must not be silently accepted.
	require.Error(t, err, "trailing bytes after a valid bundle must be rejected, not silently ignored")

	if err == nil {
		// Demonstrate the concrete gossip-dedup consequence of accepting it.
		require.Equal(t, len(clean), bytesRead, "the decoder only consumed the valid prefix")
		require.Equal(t, decodedClean.LogString(), decodedJunked.LogString(),
			"both wire messages decode to the exact same logical bundle")

		cleanID := hash.CalcHash(clean)
		junkedID := hash.CalcHash(junked)
		assert.NotEqual(t, cleanID, junkedID,
			"the two wire-identical-content messages get different gossip IDs, defeating pubsub deduplication")
	}
}
