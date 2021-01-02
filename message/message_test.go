package message

import (
	"encoding/hex"
	"fmt"
	"testing"

	peer "github.com/libp2p/go-libp2p-peer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	assert.Error(t, m2.Decode(d))
}

func TestSalamMessage(t *testing.T) {
	h := crypto.GenerateTestHash()
	_, pub, _ := crypto.GenerateTestKeyPair()
	id, _ := peer.IDB58Decode("12D3KooWDX68JokeBo8wtHtv937vM8Hj6NeNxTEh1LxZ1ragEUgn")
	m := NewSalamMessage("abc", pub, id, h, 112, 0x1)
	assert.NoError(t, m.SanityCheck())
	bs, err := m.Encode()
	fmt.Printf("%x\n", bs)
	assert.NoError(t, err)
	m2 := new(Message)
	assert.NoError(t, m2.Decode(bs))
	assert.Equal(t, m.SignBytes(), m2.SignBytes())
	assert.Equal(t, m.Type, m.Payload.Type())
	assert.Equal(t, m.Version, LastVersion)
}

func TestAleykMessage(t *testing.T) {
	_, pub, _ := crypto.GenerateTestKeyPair()
	id, _ := peer.IDB58Decode("12D3KooWDX68JokeBo8wtHtv937vM8Hj6NeNxTEh1LxZ1ragEUgn")
	m := NewAleykMessage("abc", pub, id, 112, 0x2, payload.SalamResponseCodeRejected, "Invalid genesis")
	assert.NoError(t, m.SanityCheck())
	bs, err := m.Encode()
	assert.NoError(t, err)
	m2 := new(Message)
	assert.NoError(t, m2.Decode(bs))
	assert.Equal(t, m.SignBytes(), m2.SignBytes())
	assert.Equal(t, m.Type, m.Payload.Type())
	assert.Equal(t, m.Version, LastVersion)
}

func TestLatestBlockRequestMessage(t *testing.T) {
	h := crypto.GenerateTestHash()
	invMsg := NewLatestBlocksRequestMessage(-1, h)
	assert.Error(t, invMsg.SanityCheck())
	m := NewLatestBlocksRequestMessage(1, h)
	bs, err := m.Encode()
	assert.NoError(t, err)
	m2 := new(Message)
	assert.NoError(t, m2.Decode(bs))
	assert.Equal(t, m.SignBytes(), m2.SignBytes())
	assert.Equal(t, m.Type, m.Payload.Type())
	assert.Equal(t, m.Version, LastVersion)
}

func TestLatestBlocksMessage(t *testing.T) {
	b, _ := block.GenerateTestBlock(nil, nil)
	invMsg := NewLatestBlocksMessage(4, nil, nil, nil)
	assert.Error(t, invMsg.SanityCheck())
	m := NewLatestBlocksMessage(4, []*block.Block{b}, nil, nil)
	bs, err := m.Encode()
	assert.NoError(t, err)
	m2 := new(Message)
	assert.NoError(t, m2.Decode(bs))
	assert.Equal(t, m.SignBytes(), m2.SignBytes())
	assert.Equal(t, m.Type, m.Payload.Type())
	assert.Equal(t, m.Version, LastVersion)
}

func TestTransactionsRequestMessage(t *testing.T) {
	h := crypto.GenerateTestHash()
	invMsg := NewTransactionsRequestMessage([]crypto.Hash{})
	assert.Error(t, invMsg.SanityCheck())
	m := NewTransactionsRequestMessage([]crypto.Hash{h})
	bs, err := m.Encode()
	assert.NoError(t, err)
	m2 := new(Message)
	assert.NoError(t, m2.Decode(bs))
	assert.Equal(t, m.SignBytes(), m2.SignBytes())
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
	assert.Equal(t, m.SignBytes(), m2.SignBytes())
	assert.Equal(t, m.Type, m.Payload.Type())
	assert.Equal(t, m.Version, LastVersion)
}

func TestProposalRequestMessage(t *testing.T) {
	invMsg := NewProposalRequestMessage(4, -11)
	assert.Error(t, invMsg.SanityCheck())
	m := NewProposalRequestMessage(4, 1)
	bs, err := m.Encode()
	assert.NoError(t, err)
	m2 := new(Message)
	assert.NoError(t, m2.Decode(bs))
	assert.Equal(t, m.SignBytes(), m2.SignBytes())
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
	assert.Equal(t, m.SignBytes(), m2.SignBytes())
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
	assert.Equal(t, m.SignBytes(), m2.SignBytes())
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
	assert.Equal(t, m.SignBytes(), m2.SignBytes())
	assert.Equal(t, m.Type, m.Payload.Type())
	assert.Equal(t, m.Version, LastVersion)
}

func TestHeartbeatMessage(t *testing.T) {
	m := NewHeartBeatMessage(crypto.GenerateTestHash(), hrs.NewHRS(1, 2, 3))
	assert.NoError(t, m.SanityCheck())
	bs, err := m.Encode()
	assert.NoError(t, err)
	m2 := new(Message)
	assert.NoError(t, m2.Decode(bs))
	assert.Equal(t, m.SignBytes(), m2.SignBytes())
	assert.Equal(t, m.Type, m.Payload.Type())
	assert.Equal(t, m.Version, LastVersion)
}

func TestMessageFingerprint(t *testing.T) {
	msg := NewProposalRequestMessage(1, 1)
	assert.Contains(t, msg.Fingerprint(), msg.Payload.Fingerprint())
}

func TestBlocksMessageCompress(t *testing.T) {
	var blocks = []*block.Block{}
	var trxs = []*tx.Tx{}
	for i := 0; i < 100; i++ {
		b, t := block.GenerateTestBlock(nil, nil)
		trxs = append(trxs, t...)
		blocks = append(blocks, b)
	}
	m := NewLatestBlocksMessage(888, blocks, trxs, nil)
	bs0, err := m.Encode()
	assert.NoError(t, err)
	m.CompressIt()
	bs, err := m.Encode()
	assert.NoError(t, err)
	fmt.Printf("Compressed :%v%%\n", 100-len(bs)*100/(len(bs0)))
	m2 := new(Message)
	assert.NoError(t, m2.Decode(bs))
	assert.NoError(t, m2.SanityCheck())
	assert.Equal(t, m2.Flags, 0x1)
}

func TestDecodeVoteMessage(t *testing.T) {
	d1, _ := hex.DecodeString("a401010200030a045875a101a6010102010301045820c16f004da39883f7082d39a959d9444f1cf5fb45ce5d7b0d03b6ab58f6ce5fae0554f04595cf4e14db1b179b31ae05c0656b0d835e7e065830f6510fbc1bfffa661d9562a987c5a9600004084609ded7a3d8ddbf8b09b8dc22cffcf7f19e518c90e13769ee3efe2a95")
	// Compressed
	d2, _ := hex.DecodeString("a401010201030a0458911f8b08000000000000ff0075008affa101a6010102010301045820c16f004da39883f7082d39a959d9444f1cf5fb45ce5d7b0d03b6ab58f6ce5fae0554f04595cf4e14db1b179b31ae05c0656b0d835e7e065830f6510fbc1bfffa661d9562a987c5a9600004084609ded7a3d8ddbf8b09b8dc22cffcf7f19e518c90e13769ee3efe2a95010000ffff0efa80b175000000")
	// With Signature
	d3, _ := hex.DecodeString("a501010202030a045875a101a6010102010301045820c16f004da39883f7082d39a959d9444f1cf5fb45ce5d7b0d03b6ab58f6ce5fae0554f04595cf4e14db1b179b31ae05c0656b0d835e7e065830f6510fbc1bfffa661d9562a987c5a9600004084609ded7a3d8ddbf8b09b8dc22cffcf7f19e518c90e13769ee3efe2a95155830882e40ab14c705196049bd1520e5f256ac87ff5208aacd84f7120429bd15b7407caabb190cfad0a8f19743063517c40a")
	m1 := new(Message)
	m2 := new(Message)
	m3 := new(Message)
	assert.NoError(t, m1.Decode(d1))
	assert.NoError(t, m2.Decode(d2))
	assert.NoError(t, m3.Decode(d3))
	assert.NoError(t, m2.SanityCheck())
	assert.NoError(t, m3.SanityCheck())

	assert.Equal(t, m1.SignBytes(), m2.SignBytes())
	assert.Equal(t, m1.SignBytes(), m3.SignBytes())
}

func TestAddSignature(t *testing.T) {
	m1 := NewProposalRequestMessage(1, 1)
	m2 := new(Message)
	_, _, priv := crypto.GenerateTestKeyPair()
	sig := priv.Sign(m1.SignBytes())
	m1.SetSignature(sig)
	bs, err := m1.Encode()
	require.NoError(t, err)
	err = m2.Decode(bs)
	assert.NoError(t, err)

	assert.NoError(t, m2.SanityCheck())
	assert.NotNil(t, m2.Signature)
}
