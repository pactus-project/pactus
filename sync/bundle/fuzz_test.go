package bundle

import (
	"bytes"
	"testing"
)

// A bundle comes from any peer: decoding it, and checking its message, never panics.
func FuzzBundleDecode(f *testing.F) {
	// The seeds are in testdata/fuzz/FuzzBundleDecode: a vote, a proposal, a block announcement
	// and a transactions message.
	f.Add([]byte{})

	f.Fuzz(func(_ *testing.T, data []byte) {
		bdl := new(Bundle)
		if _, err := bdl.Decode(bytes.NewReader(data)); err != nil {
			return
		}

		_ = bdl.BasicCheck()
		_ = bdl.LogString()
	})
}
