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
	d1, _ := hex.DecodeString("a50101020003582212202942316bcfd7550c000000000000000000000000000000000000000000000000040b055877a101a60102021858030004582047088d70ce69076d02609a4dde8ce6014e0b67b0c32e084f1a771dce4c149695055501145f14e9182d002a94bae32e6ca53627c5e17e1d0658303679874b01cdd1eb79470ee5a630e366f2192cb64d008073605b85f56a0e86280aecfc083130f48a2efc5615f7298581")
	// Compressed
	d2, _ := hex.DecodeString("a50101020103582212202942316bcfd7550c000000000000000000000000000000000000000000000000040b0558931f8b08000000000000ff00770088ffa101a60102021858030004582047088d70ce69076d02609a4dde8ce6014e0b67b0c32e084f1a771dce4c149695055501145f14e9182d002a94bae32e6ca53627c5e17e1d0658303679874b01cdd1eb79470ee5a630e366f2192cb64d008073605b85f56a0e86280aecfc083130f48a2efc5615f7298581010000ffff763d7ce777000000")
	m1 := new(Message)
	m2 := new(Message)
	assert.NoError(t, m1.Decode(d1))
	assert.NoError(t, m2.Decode(d2))
	assert.NoError(t, m2.SanityCheck())

	assert.Equal(t, m1.Payload, m2.Payload)
}
