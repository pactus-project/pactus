package network

import (
	"testing"

	lp2pspb "github.com/libp2p/go-libp2p-pubsub/pb"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/stretchr/testify/assert"
)

// The identifier is the hash of the topic, a zero byte and the data.
func TestMessageIDIsTheHashOfTopicAndData(t *testing.T) {
	topic := "/pactus/topic/consensus/v1"
	for _, data := range [][]byte{nil, {}, []byte("vote"), make([]byte, 1<<20)} {
		concat := append(append([]byte(topic), 0), data...)
		want := hash.CalcHash(concat)

		assert.Equal(t, string(want[:20]), MessageIDFunc(&lp2pspb.Message{Topic: &topic, Data: data}))
	}
}

func BenchmarkMessageID(b *testing.B) {
	topic := "/pactus/topic/block/v1"
	msg := &lp2pspb.Message{Topic: &topic, Data: make([]byte, 1<<20)}

	b.ReportAllocs()
	for range b.N {
		MessageIDFunc(msg)
	}
}
