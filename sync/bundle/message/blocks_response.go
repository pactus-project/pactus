package message

import (
	"fmt"

	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/errors"
)

const LatestBlocksResponseCodeOK = 0
const LatestBlocksResponseCodeNoMoreBlock = 1

type BlocksResponseMessage struct {
	ResponseCode    ResponseCode       `cbor:"1,keyasint"`
	SessionID       int                `cbor:"2,keyasint"`
	From            int                `cbor:"3,keyasint"`
	Blocks          []*block.Block     `cbor:"4,keyasint"`
	LastCertificate *block.Certificate `cbor:"6,keyasint"`
}

func NewBlocksResponseMessage(code ResponseCode, sid int, from int,
	blocks []*block.Block, cert *block.Certificate) *BlocksResponseMessage {
	return &BlocksResponseMessage{
		ResponseCode:    code,
		SessionID:       sid,
		From:            from,
		Blocks:          blocks,
		LastCertificate: cert,
	}
}
func (m *BlocksResponseMessage) SanityCheck() error {
	if m.From < 0 {
		return errors.Errorf(errors.ErrInvalidMessage, "invalid Height")
	}
	for _, b := range m.Blocks {
		if err := b.SanityCheck(); err != nil {
			return errors.Errorf(errors.ErrInvalidMessage, err.Error())
		}
	}
	if m.LastCertificate != nil {
		if err := m.LastCertificate.SanityCheck(); err != nil {
			return errors.Errorf(errors.ErrInvalidMessage, err.Error())
		}
	}

	return nil
}

func (m *BlocksResponseMessage) Type() Type {
	return MessageTypeBlocksResponse
}

func (m *BlocksResponseMessage) To() int {
	if len(m.Blocks) == 0 {
		return m.From
	}
	return m.From + len(m.Blocks) - 1
}

func (m *BlocksResponseMessage) Fingerprint() string {
	return fmt.Sprintf("{⚓ %d %s %v-%v}", m.SessionID, m.ResponseCode, m.From, m.To())
}

func (m *BlocksResponseMessage) IsRequestRejected() bool {
	if m.ResponseCode == ResponseCodeBusy ||
		m.ResponseCode == ResponseCodeRejected {
		return true
	}

	return false
}
