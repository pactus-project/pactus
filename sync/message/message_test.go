package message

import (
	"encoding/hex"
	"fmt"
	"testing"

	peer "github.com/libp2p/go-libp2p-peer"
	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/consensus/hrs"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/sync/message/payload"
	"github.com/zarbchain/zarb-go/tx"
	"github.com/zarbchain/zarb-go/vote"
)

var tPeerID1 peer.ID
var tPeerID2 peer.ID

func init() {
	tPeerID1, _ = peer.IDB58Decode("12D3KooWDX68JokeBo8wtHtv937vM8Hj6NeNxTEh1LxZ1ragEUgn")
	tPeerID2, _ = peer.IDB58Decode("12D3KooWBAoXit47Mxbax1FdbE91aphCkR9V2nUZ1yiXDCe3nh18")

}
func TestInvalidPayloadType(t *testing.T) {
	assert.Nil(t, makePayload(0))
	m := &Message{Type: 1, Payload: makePayload(2)}
	assert.Error(t, m.SanityCheck())
}

func TestInvalidCBOR(t *testing.T) {
	d, _ := hex.DecodeString("a40100030004010aa301")
	m2 := new(Message)
	assert.Error(t, m2.Decode(d))
}

func TestSalamMessage(t *testing.T) {
	h := crypto.GenerateTestHash()
	_, pub, _ := crypto.GenerateTestKeyPair()
	m := NewSalamMessage("abc", pub, tPeerID1, h, 112, 0x1)
	assert.NoError(t, m.SanityCheck())
	bs, err := m.Encode()
	fmt.Printf("%x\n", bs)
	assert.NoError(t, err)
	m2 := new(Message)
	assert.NoError(t, m2.Decode(bs))
	assert.Equal(t, m.Type, m.Payload.Type())
	assert.Equal(t, m.Version, LastVersion)
}

func TestAleykMessage(t *testing.T) {
	_, pub, _ := crypto.GenerateTestKeyPair()
	m := NewAleykMessage(payload.ResponseCodeRejected, "Invalid genesis", "cute-kitty", pub, tPeerID1, 112, 0x2)
	assert.NoError(t, m.SanityCheck())
	bs, err := m.Encode()
	assert.NoError(t, err)
	m2 := new(Message)
	assert.NoError(t, m2.Decode(bs))
	assert.Equal(t, m.Type, m.Payload.Type())
	assert.Equal(t, m.Version, LastVersion)
}

func TestLatestBlockRequestMessage(t *testing.T) {
	invMsg := NewLatestBlocksRequestMessage(tPeerID1, tPeerID2, 1234, -1)
	assert.Error(t, invMsg.SanityCheck())
	m := NewLatestBlocksRequestMessage(tPeerID1, tPeerID2, 1234, 100)
	bs, err := m.Encode()
	assert.NoError(t, err)
	m2 := new(Message)
	assert.NoError(t, m2.Decode(bs))
	assert.Equal(t, m.Type, m.Payload.Type())
	assert.Equal(t, m.Version, LastVersion)
}

func TestLatestBlocksResponseMessage(t *testing.T) {
	b, trxs := block.GenerateTestBlock(nil, nil)
	invMsg := NewLatestBlocksResponseMessage(payload.ResponseCodeBusy, tPeerID1, tPeerID2, 1234, -1, nil, nil, nil)
	assert.Error(t, invMsg.SanityCheck())
	m := NewLatestBlocksResponseMessage(payload.ResponseCodeBusy, tPeerID1, tPeerID2, 1234, 4, []*block.Block{b}, trxs, nil)
	bs, err := m.Encode()
	assert.NoError(t, err)
	m2 := new(Message)
	assert.NoError(t, m2.Decode(bs))
	assert.Equal(t, m.Type, m.Payload.Type())
	assert.Equal(t, m.Version, LastVersion)
}

func TestQueryTransactionsMessage(t *testing.T) {
	h := crypto.GenerateTestHash()
	invMsg := NewQueryTransactionsMessage([]crypto.Hash{})
	assert.Error(t, invMsg.SanityCheck())
	m := NewQueryTransactionsMessage([]crypto.Hash{h})
	bs, err := m.Encode()
	assert.NoError(t, err)
	m2 := new(Message)
	assert.NoError(t, m2.Decode(bs))
	assert.Equal(t, m.Type, m.Payload.Type())
	assert.Equal(t, m.Version, LastVersion)
}

func TestTransactionsMessage(t *testing.T) {
	trx, _ := tx.GenerateTestSendTx()
	invMsg := NewTransactionsMessage([]*tx.Tx{})
	assert.Error(t, invMsg.SanityCheck())
	m := NewTransactionsMessage([]*tx.Tx{trx})
	bs, err := m.Encode()
	assert.NoError(t, err)
	m2 := new(Message)
	assert.NoError(t, m2.Decode(bs))
	assert.Equal(t, m.Type, m.Payload.Type())
	assert.Equal(t, m.Version, LastVersion)
}

func TestQueryProposalMessage(t *testing.T) {
	invMsg := NewQueryProposalMessage(4, -11)
	assert.Error(t, invMsg.SanityCheck())
	m := NewQueryProposalMessage(4, 1)
	bs, err := m.Encode()
	assert.NoError(t, err)
	m2 := new(Message)
	assert.NoError(t, m2.Decode(bs))
	assert.Equal(t, m.Type, m.Payload.Type())
	assert.Equal(t, m.Version, LastVersion)
}

func TestProposalsMessage(t *testing.T) {
	p, _ := vote.GenerateTestProposal(5, 1)
	invMsg := NewProposalMessage(nil)
	assert.Error(t, invMsg.SanityCheck())
	m := NewProposalMessage(p)
	bs, err := m.Encode()
	assert.NoError(t, err)
	m2 := new(Message)
	assert.NoError(t, m2.Decode(bs))
	assert.Equal(t, m.Type, m.Payload.Type())
	assert.Equal(t, m.Version, LastVersion)
}

func TestVoteSetMessage(t *testing.T) {
	m := NewVoteSetMessage(4, 1, []crypto.Hash{})
	assert.NoError(t, m.SanityCheck())
	bs, err := m.Encode()
	assert.NoError(t, err)
	m2 := new(Message)
	assert.NoError(t, m2.Decode(bs))
	assert.Equal(t, m.Type, m.Payload.Type())
	assert.Equal(t, m.Version, LastVersion)
}

func TestVoteMessage(t *testing.T) {
	v, _ := vote.GenerateTestPrepareVote(1, 1)
	m := NewVoteMessage(v)
	assert.NoError(t, m.SanityCheck())
	bs, err := m.Encode()
	assert.NoError(t, err)
	m2 := new(Message)
	assert.NoError(t, m2.Decode(bs))
	assert.Equal(t, m.Type, m.Payload.Type())
	assert.Equal(t, m.Version, LastVersion)
}

func TestBlockAnnounceMessage(t *testing.T) {
	b, _ := block.GenerateTestBlock(nil, nil)
	c := block.GenerateTestCommit(b.Hash())
	m := NewBlockAnnounceMessage(1001, b, c)
	assert.NoError(t, m.SanityCheck())
	bs, err := m.Encode()
	assert.NoError(t, err)
	m2 := new(Message)
	assert.NoError(t, m2.Decode(bs))
	assert.Equal(t, m.Type, m.Payload.Type())
	assert.Equal(t, m.Version, LastVersion)
}

func TestHeartbeatMessage(t *testing.T) {
	m := NewHeartBeatMessage(tPeerID1, crypto.GenerateTestHash(), hrs.NewHRS(1, 2, 3))
	assert.NoError(t, m.SanityCheck())
	bs, err := m.Encode()
	assert.NoError(t, err)
	m2 := new(Message)
	assert.NoError(t, m2.Decode(bs))
	assert.Equal(t, m.Type, m.Payload.Type())
	assert.Equal(t, m.Version, LastVersion)
}

func TestDownloadRequest(t *testing.T) {
	m := NewDownloadRequestMessage(tPeerID1, tPeerID2, 6789, 1234, 2234)
	assert.NoError(t, m.SanityCheck())
	bs, err := m.Encode()
	assert.NoError(t, err)
	m2 := new(Message)
	assert.NoError(t, m2.Decode(bs))
	assert.Equal(t, m.Type, m.Payload.Type())
	assert.Equal(t, m.Version, LastVersion)
}

func TestDownloadResponse(t *testing.T) {
	b, trxs := block.GenerateTestBlock(nil, nil)
	m := NewDownloadResponseMessage(payload.ResponseCodeBusy, tPeerID1, tPeerID2, 6789, 1234, []*block.Block{b}, trxs)
	assert.NoError(t, m.SanityCheck())
	bs, err := m.Encode()
	assert.NoError(t, err)
	m2 := new(Message)
	assert.NoError(t, m2.Decode(bs))
	assert.Equal(t, m.Type, m.Payload.Type())
	assert.Equal(t, m.Version, LastVersion)
}

func TestBlocksMessageCompress(t *testing.T) {
	var blocks = []*block.Block{}
	var trxs = []*tx.Tx{}
	for i := 0; i < 10; i++ {
		b, t := block.GenerateTestBlock(nil, nil)
		trxs = append(trxs, t...)
		blocks = append(blocks, b)
	}
	m := NewLatestBlocksResponseMessage(payload.ResponseCodeBusy, tPeerID1, tPeerID2, 1234, 888, blocks, trxs, nil)
	bs0, err := m.Encode()
	assert.NoError(t, err)
	m.CompressIt()
	bs, err := m.Encode()
	assert.NoError(t, err)
	fmt.Printf("Compressed :%v%%\n", 100-len(bs)*100/(len(bs0)))
	fmt.Printf("Compressed len :%v\n", len(bs))
	m2 := new(Message)
	assert.NoError(t, m2.Decode(bs))
	assert.NoError(t, m2.SanityCheck())
	assert.Equal(t, m2.Flags, 0x1)
}

func TestDecodeVoteMessage(t *testing.T) {
	d1, _ := hex.DecodeString("a401010200030a045875a101a601010201030104582070f4338d6ed218ba9e0884352bcee8c19cfe5e110d4aa8b248881a90a96b7d980554d94f1e7acfdc43db98e3cbd51f56e832be4d6d73065830f4dbe13687e792bdf51ddcdf746e38001c60a9b05e94ae86e5ffc116a835e930e32b5da0e6163a0f785ed294f4c69e18")
	// Compressed
	d2, _ := hex.DecodeString("a401010201030a0458911f8b08000000000000ff0075008affa101a601010201030104582070f4338d6ed218ba9e0884352bcee8c19cfe5e110d4aa8b248881a90a96b7d980554d94f1e7acfdc43db98e3cbd51f56e832be4d6d73065830f4dbe13687e792bdf51ddcdf746e38001c60a9b05e94ae86e5ffc116a835e930e32b5da0e6163a0f785ed294f4c69e18010000ffffc303af1e75000000")
	m1 := new(Message)
	m2 := new(Message)
	assert.NoError(t, m1.Decode(d1))
	assert.NoError(t, m2.Decode(d2))
	assert.NoError(t, m2.SanityCheck())

	assert.Equal(t, m1.Payload, m2.Payload)
}
