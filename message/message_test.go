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
	"github.com/zarbchain/zarb-go/message/payload"
	"github.com/zarbchain/zarb-go/tx"
	"github.com/zarbchain/zarb-go/vote"
)

func TestInvalidCBOR(t *testing.T) {
	d, _ := hex.DecodeString("a40100030004010aa301a3654d616a6f7201654d696e6f720065506174636800025820a723d8cb4fa6a4a67a4a6a8984f81cb737defb3a8daaacafcd3159f1dc545c03031870")
	m2 := new(Message)
	assert.Error(t, m2.UnmarshalCBOR(d))
}

func TestSalamMessage(t *testing.T) {
	h := crypto.GenerateTestHash()
	_, pub, _ := crypto.GenerateTestKeyPair()
	id, _ := peer.IDB58Decode("12D3KooWDX68JokeBo8wtHtv937vM8Hj6NeNxTEh1LxZ1ragEUgn")
	m := NewSalamMessage("abc", pub, id, h, 112)
	assert.NoError(t, m.SanityCheck())
	bs, err := m.MarshalCBOR()
	fmt.Printf("%x\n", bs)
	assert.NoError(t, err)
	m2 := new(Message)
	assert.NoError(t, m2.UnmarshalCBOR(bs))
	assert.Equal(t, m.SignBytes(), m2.SignBytes())
	assert.Equal(t, m.Type, m.Payload.Type())
	assert.Equal(t, m.Version, LastVersion)
}

func TestAleykMessage(t *testing.T) {
	h := crypto.GenerateTestHash()
	_, pub, _ := crypto.GenerateTestKeyPair()
	id, _ := peer.IDB58Decode("12D3KooWDX68JokeBo8wtHtv937vM8Hj6NeNxTEh1LxZ1ragEUgn")
	m := NewAleykMessage("abc", pub, id, h, 112, payload.SalamResponseCodeRejected, "Invalid genesis")
	assert.NoError(t, m.SanityCheck())
	bs, err := m.MarshalCBOR()
	assert.NoError(t, err)
	m2 := new(Message)
	assert.NoError(t, m2.UnmarshalCBOR(bs))
	assert.Equal(t, m.SignBytes(), m2.SignBytes())
	assert.Equal(t, m.Type, m.Payload.Type())
	assert.Equal(t, m.Version, LastVersion)
}

func TestBlockReqMessage(t *testing.T) {
	h := crypto.GenerateTestHash()
	invMsg := NewBlocksReqMessage(4, 1, h)
	assert.Error(t, invMsg.SanityCheck())
	m := NewBlocksReqMessage(1, 4, h)
	bs, err := m.MarshalCBOR()
	assert.NoError(t, err)
	m2 := new(Message)
	assert.NoError(t, m2.UnmarshalCBOR(bs))
	assert.Equal(t, m.SignBytes(), m2.SignBytes())
	assert.Equal(t, m.Type, m.Payload.Type())
	assert.Equal(t, m.Version, LastVersion)
}

func TestBlocksMessage(t *testing.T) {
	b, _ := block.GenerateTestBlock(nil, nil)
	invMsg := NewBlocksMessage(4, nil, nil)
	assert.Error(t, invMsg.SanityCheck())
	m := NewBlocksMessage(4, []*block.Block{b}, nil)
	bs, err := m.MarshalCBOR()
	assert.NoError(t, err)
	m2 := new(Message)
	assert.NoError(t, m2.UnmarshalCBOR(bs))
	assert.Equal(t, m.SignBytes(), m2.SignBytes())
	assert.Equal(t, m.Type, m.Payload.Type())
	assert.Equal(t, m.Version, LastVersion)
}

func TestTxReqMessage(t *testing.T) {
	h := crypto.GenerateTestHash()
	invMsg := NewTxsReqMessage([]crypto.Hash{})
	assert.Error(t, invMsg.SanityCheck())
	m := NewTxsReqMessage([]crypto.Hash{h})
	bs, err := m.MarshalCBOR()
	assert.NoError(t, err)
	m2 := new(Message)
	assert.NoError(t, m2.UnmarshalCBOR(bs))
	assert.Equal(t, m.SignBytes(), m2.SignBytes())
	assert.Equal(t, m.Type, m.Payload.Type())
	assert.Equal(t, m.Version, LastVersion)
}

func TestTxsMessage(t *testing.T) {
	trx, _ := tx.GenerateTestSendTx()
	invMsg := NewTxsMessage([]*tx.Tx{})
	assert.Error(t, invMsg.SanityCheck())
	m := NewTxsMessage([]*tx.Tx{trx})
	bs, err := m.MarshalCBOR()
	assert.NoError(t, err)
	m2 := new(Message)
	assert.NoError(t, m2.UnmarshalCBOR(bs))
	assert.Equal(t, m.SignBytes(), m2.SignBytes())
	assert.Equal(t, m.Type, m.Payload.Type())
	assert.Equal(t, m.Version, LastVersion)
}

func TestProposalReqMessage(t *testing.T) {
	invMsg := NewProposalReqMessage(4, -11)
	assert.Error(t, invMsg.SanityCheck())
	m := NewProposalReqMessage(4, 1)
	bs, err := m.MarshalCBOR()
	assert.NoError(t, err)
	m2 := new(Message)
	assert.NoError(t, m2.UnmarshalCBOR(bs))
	assert.Equal(t, m.SignBytes(), m2.SignBytes())
	assert.Equal(t, m.Type, m.Payload.Type())
	assert.Equal(t, m.Version, LastVersion)
}

func TestProposalsMessage(t *testing.T) {
	p, _ := vote.GenerateTestProposal(5, 1)
	invMsg := NewProposalMessage(nil)
	assert.Error(t, invMsg.SanityCheck())
	m := NewProposalMessage(p)
	bs, err := m.MarshalCBOR()
	assert.NoError(t, err)
	m2 := new(Message)
	assert.NoError(t, m2.UnmarshalCBOR(bs))
	assert.Equal(t, m.SignBytes(), m2.SignBytes())
	assert.Equal(t, m.Type, m.Payload.Type())
	assert.Equal(t, m.Version, LastVersion)
}

func TestVoteSetMessage(t *testing.T) {
	m := NewVoteSetMessage(4, 1, []crypto.Hash{})
	assert.NoError(t, m.SanityCheck())
	bs, err := m.MarshalCBOR()
	assert.NoError(t, err)
	m2 := new(Message)
	assert.NoError(t, m2.UnmarshalCBOR(bs))
	assert.Equal(t, m.SignBytes(), m2.SignBytes())
	assert.Equal(t, m.Type, m.Payload.Type())
	assert.Equal(t, m.Version, LastVersion)
}

func TestVoteMessage(t *testing.T) {
	v, _ := vote.GenerateTestPrepareVote(1, 1)
	m := NewVoteMessage(v)
	assert.NoError(t, m.SanityCheck())
	bs, err := m.MarshalCBOR()
	assert.NoError(t, err)
	m2 := new(Message)
	assert.NoError(t, m2.UnmarshalCBOR(bs))
	assert.Equal(t, m.SignBytes(), m2.SignBytes())
	assert.Equal(t, m.Type, m.Payload.Type())
	assert.Equal(t, m.Version, LastVersion)
}

func TestHeartbeatMessage(t *testing.T) {
	m := NewHeartBeatMessage(crypto.GenerateTestHash(), hrs.NewHRS(1, 2, 3))
	assert.NoError(t, m.SanityCheck())
	bs, err := m.MarshalCBOR()
	assert.NoError(t, err)
	m2 := new(Message)
	assert.NoError(t, m2.UnmarshalCBOR(bs))
	assert.Equal(t, m.SignBytes(), m2.SignBytes())
	assert.Equal(t, m.Type, m.Payload.Type())
	assert.Equal(t, m.Version, LastVersion)
}

func TestMessageFingerprint(t *testing.T) {
	msg := NewProposalReqMessage(1, 1)
	assert.Contains(t, msg.Fingerprint(), msg.Payload.Fingerprint())
}

func TestBlocksMessageCompress(t *testing.T) {
	var blocks = []*block.Block{}
	for i := 0; i < 100; i++ {
		b, _ := block.GenerateTestBlock(nil, nil)
		blocks = append(blocks, b)
	}
	m := NewBlocksMessage(888, blocks, nil)
	bs, err := m.MarshalCBOR()
	assert.NoError(t, err)
	m2 := new(Message)
	assert.NoError(t, m2.UnmarshalCBOR(bs))
	assert.NoError(t, m2.SanityCheck())
	assert.Equal(t, m2.Flags, 0x1)
}

func TestDecodeVoteMessage(t *testing.T) {
	d1, _ := hex.DecodeString("a401010200030a045875a101a6010102010301045820c16f004da39883f7082d39a959d9444f1cf5fb45ce5d7b0d03b6ab58f6ce5fae0554f04595cf4e14db1b179b31ae05c0656b0d835e7e065830f6510fbc1bfffa661d9562a987c5a9600004084609ded7a3d8ddbf8b09b8dc22cffcf7f19e518c90e13769ee3efe2a95")
	// Compressed
	d2, _ := hex.DecodeString("a401010201030a0458911f8b08000000000000ff0075008affa101a6010102010301045820c16f004da39883f7082d39a959d9444f1cf5fb45ce5d7b0d03b6ab58f6ce5fae0554f04595cf4e14db1b179b31ae05c0656b0d835e7e065830f6510fbc1bfffa661d9562a987c5a9600004084609ded7a3d8ddbf8b09b8dc22cffcf7f19e518c90e13769ee3efe2a95010000ffff0efa80b175000000")
	m1 := new(Message)
	m2 := new(Message)
	assert.NoError(t, m1.UnmarshalCBOR(d1))
	assert.NoError(t, m2.UnmarshalCBOR(d2))
	assert.NoError(t, m2.SanityCheck())

	assert.Equal(t, m1.SignBytes(), m2.SignBytes())
}
