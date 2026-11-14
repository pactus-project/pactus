package message

import (
	"testing"

	"github.com/fxamacker/cbor/v2"
)

// A message comes from any peer: decoding each type of message from any bytes,
// and checking it, never panics.
func FuzzMessageDecode(f *testing.F) {
	// The seeds are in testdata/fuzz/FuzzMessageDecode: one valid message of each type.
	for msgType := TypeHello; msgType <= TypeBlocksResponse; msgType++ {
		f.Add(int32(msgType), []byte{})
	}

	f.Fuzz(func(_ *testing.T, msgType int32, data []byte) {
		msg, err := MakeMessage(Type(msgType))
		if err != nil {
			return
		}
		if err := cbor.Unmarshal(data, msg); err != nil {
			return
		}

		_ = msg.BasicCheck()
		_ = msg.LogString()
		_ = msg.ConsensusHeight()
		_ = msg.ShouldBroadcast()
		_ = msg.TopicID()
	})
}
