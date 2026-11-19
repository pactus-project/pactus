package network

import (
	"encoding/hex"
	"testing"

	lp2pspb "github.com/libp2p/go-libp2p-pubsub/pb"
	"github.com/stretchr/testify/assert"
)

// The identifier is the same on every node: a fixed vector pins it (the first 20 bytes of
// BLAKE2b-256 of "topic", a zero byte and "data", computed outside of Go). The topic and the data
// are separated, so the bytes of the topic can't be moved into the data.
func TestMessageIDVector(t *testing.T) {
	msgID := func(topic, data string) string {
		return MessageIDFunc(&lp2pspb.Message{Topic: &topic, Data: []byte(data)})
	}

	assert.Equal(t, "e751d7edb9c2c1c8b9a216e9495b99f003880e9e", hex.EncodeToString([]byte(msgID("topic", "data"))))
	assert.NotEqual(t, msgID("ab", "c"), msgID("a", "bc"))
}
