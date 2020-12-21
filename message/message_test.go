package message

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/consensus/hrs"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/tx"
	"github.com/zarbchain/zarb-go/vote"
)

func TestSalamMessage(t *testing.T) {
	h := crypto.GenerateTestHash()
	m := NewSalamMessage(h, 112)
	assert.NoError(t, m.SanityCheck())
	bs, err := m.MarshalCBOR()
	assert.NoError(t, err)
	m2 := new(Message)
	assert.NoError(t, m2.UnmarshalCBOR(bs))
	assert.Equal(t, m, m2)
	assert.Equal(t, m.Type, m.Payload.Type())
}

func TestAleykMessage(t *testing.T) {
	h := crypto.GenerateTestHash()
	m := NewAleykMessage(h, 112)
	assert.NoError(t, m.SanityCheck())
	bs, err := m.MarshalCBOR()
	assert.NoError(t, err)
	m2 := new(Message)
	assert.NoError(t, m2.UnmarshalCBOR(bs))
	assert.Equal(t, m, m2)
	assert.Equal(t, m.Type, m.Payload.Type())
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
	assert.Equal(t, m, m2)
	assert.Equal(t, m.Type, m.Payload.Type())
}

func TestBlocksMessage(t *testing.T) {
	b, _ := block.GenerateTestBlock(nil)
	invMsg := NewBlocksMessage(4, nil, nil)
	assert.Error(t, invMsg.SanityCheck())
	m := NewBlocksMessage(4, []*block.Block{b}, nil)
	bs, err := m.MarshalCBOR()
	assert.NoError(t, err)
	m2 := new(Message)
	assert.NoError(t, m2.UnmarshalCBOR(bs))
	bs2, err := m.MarshalCBOR()
	assert.NoError(t, err)
	assert.EqualValues(t, bs, bs2)
	assert.Equal(t, m.Type, m.Payload.Type())
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
	assert.Equal(t, m, m2)
	assert.Equal(t, m.Type, m.Payload.Type())
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
	bs2, err := m.MarshalCBOR()
	assert.NoError(t, err)
	assert.EqualValues(t, bs, bs2)
	assert.Equal(t, m.Type, m.Payload.Type())
}

func TestProposalReqMessage(t *testing.T) {
	invMsg := NewProposalReqMessage(4, -11)
	assert.Error(t, invMsg.SanityCheck())
	m := NewProposalReqMessage(4, 1)
	bs, err := m.MarshalCBOR()
	assert.NoError(t, err)
	m2 := new(Message)
	assert.NoError(t, m2.UnmarshalCBOR(bs))
	assert.Equal(t, m, m2)
	assert.Equal(t, m.Type, m.Payload.Type())
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
	bs2, err := m.MarshalCBOR()
	assert.NoError(t, err)
	assert.EqualValues(t, bs, bs2)
	assert.Equal(t, m.Type, m.Payload.Type())
}

func TestVoteSetMessage(t *testing.T) {
	m := NewVoteSetMessage(4, []crypto.Hash{})
	assert.NoError(t, m.SanityCheck())
	bs, err := m.MarshalCBOR()
	assert.NoError(t, err)
	m2 := new(Message)
	assert.NoError(t, m2.UnmarshalCBOR(bs))
	assert.Equal(t, m, m2)
	assert.Equal(t, m.Type, m.Payload.Type())
}

func TestVoteMessage(t *testing.T) {
	v, _ := vote.GenerateTestPrecommitVote(1, 1)
	m := NewVoteMessage(v)
	assert.NoError(t, m.SanityCheck())
	bs, err := m.MarshalCBOR()
	assert.NoError(t, err)
	m2 := new(Message)
	assert.NoError(t, m2.UnmarshalCBOR(bs))
	bs2, err := m.MarshalCBOR()
	assert.NoError(t, err)
	assert.EqualValues(t, bs, bs2)
	assert.Equal(t, m.Type, m.Payload.Type())
}

func TestHeartbeatMessage(t *testing.T) {
	m := NewHeartBeatMessage(crypto.GenerateTestHash(), hrs.NewHRS(1, 2, 3))
	assert.NoError(t, m.SanityCheck())
	bs, err := m.MarshalCBOR()
	assert.NoError(t, err)
	m2 := new(Message)
	assert.NoError(t, m2.UnmarshalCBOR(bs))
	assert.Equal(t, m, m2)
	assert.Equal(t, m.Type, m.Payload.Type())
}
