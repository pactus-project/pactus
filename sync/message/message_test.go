package message

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/consensus/vote"
	"github.com/zarbchain/zarb-go/sync/message/payload"
	"github.com/zarbchain/zarb-go/tx"
	"github.com/zarbchain/zarb-go/util"
)

var tPeerID1 peer.ID
var tPeerID2 peer.ID

func init() {
	tPeerID1 = util.RandomPeerID()
	tPeerID2 = util.RandomPeerID()

}

func TestInvalidCBOR(t *testing.T) {
	d, _ := hex.DecodeString("a40100030004010aa301")
	m2 := new(Message)
	assert.Error(t, m2.Decode(d))
}
func TestMessageCompress(t *testing.T) {
	var blocks = []*block.Block{}
	var trxs = []*tx.Tx{}
	for i := 0; i < 10; i++ {
		b, t := block.GenerateTestBlock(nil, nil)
		trxs = append(trxs, t...)
		blocks = append(blocks, b)
	}
	pld := payload.NewLatestBlocksResponsePayload(payload.ResponseCodeBusy, 1234, tPeerID2, 888, blocks, trxs, nil)
	msg := NewMessage(tPeerID1, pld)
	bs0, err := msg.Encode()
	assert.NoError(t, err)
	msg.CompressIt()
	bs1, err := msg.Encode()
	assert.NoError(t, err)
	fmt.Printf("Compressed :%v%%\n", 100-len(bs1)*100/(len(bs0)))
	fmt.Printf("Uncompressed len :%v\n", len(bs0))
	fmt.Printf("Compressed len :%v\n", len(bs1))
	msg2 := new(Message)
	msg3 := new(Message)
	assert.NoError(t, msg2.Decode(bs0))
	assert.NoError(t, msg3.Decode(bs1))
	assert.NoError(t, msg2.SanityCheck())
	assert.NoError(t, msg3.SanityCheck())
}

func TestDecodeVoteMessage(t *testing.T) {
	v, _ := vote.GenerateTestPrecommitVote(88, 0)
	pld := payload.NewVotePayload(v)
	msg := NewMessage(tPeerID1, pld)
	bs0, err := msg.Encode()
	assert.NoError(t, err)
	msg.CompressIt()
	bs1, err := msg.Encode()
	assert.NoError(t, err)
	fmt.Printf("Compressed :%v%%\n", 100-len(bs1)*100/(len(bs0)))
	fmt.Printf("Uncompressed len :%v\n", len(bs0))
	fmt.Printf("Compressed len :%v\n", len(bs1))
}

func TestDecodeVoteCBOR(t *testing.T) {
	d1, _ := hex.DecodeString("a50101020003582212202202894105afb730000000000000000000000000000000000000000000000000040b055876a101a601020218580300045820d9c5f68143672b2ad64da4be8390224ee61566c389ba844040e16a8c136a816d055495ee67563b33004f893fd5616ddc3f910ffc9be1065830fa5ecd896f0eba5ee6e1115a253184780c45b17a01bd2ff2adc0197c96edd50c4be33cdac2f31321734d3103b7a00c84")
	// Compressed
	d2, _ := hex.DecodeString("a50101020103582212202202894105afb730000000000000000000000000000000000000000000000000040b0558921f8b08000000000000ff00760089ffa101a601020218580300045820d9c5f68143672b2ad64da4be8390224ee61566c389ba844040e16a8c136a816d055495ee67563b33004f893fd5616ddc3f910ffc9be1065830fa5ecd896f0eba5ee6e1115a253184780c45b17a01bd2ff2adc0197c96edd50c4be33cdac2f31321734d3103b7a00c84010000fffff37a941a76000000")
	m1 := new(Message)
	m2 := new(Message)
	assert.NoError(t, m1.Decode(d1))
	assert.NoError(t, m2.Decode(d2))
	assert.NoError(t, m2.SanityCheck())

	assert.Equal(t, m1.Payload, m2.Payload)
}
