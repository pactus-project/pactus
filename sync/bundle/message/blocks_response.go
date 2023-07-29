package message

import (
	"fmt"

	"github.com/pactus-project/pactus/types/block"
	"github.com/pactus-project/pactus/util/errors"
)

const LatestBlocksResponseCodeOK = 0
const LatestBlocksResponseCodeNoMoreBlock = 1

type BlocksResponseMessage struct {
	ResponseCode    ResponseCode       `cbor:"1,keyasint"`
	SessionID       int                `cbor:"2,keyasint"`
	From            uint32             `cbor:"3,keyasint"`
	BlocksData      [][]byte           `cbor:"4,keyasint"`
	LastCertificate *block.Certificate `cbor:"6,keyasint"`
	Reason          string             `cbor:"7,keyasint"`
}

func NewBlocksResponseMessage(code ResponseCode, reason string, sid int, from uint32,
	blocksData [][]byte, lastCert *block.Certificate) *BlocksResponseMessage {
	return &BlocksResponseMessage{
		ResponseCode:    code,
		SessionID:       sid,
		From:            from,
		BlocksData:      blocksData,
		LastCertificate: lastCert,
		Reason:          reason,
	}
}
func (m *BlocksResponseMessage) SanityCheck() error {
	if m.From == 0 && len(m.BlocksData) != 0 {
		return errors.Errorf(errors.ErrInvalidHeight, "unexpected block for height zero")
	}
	if m.LastCertificate != nil {
		if err := m.LastCertificate.SanityCheck(); err != nil {
			return err
		}
	}

	return nil
}

func (m *BlocksResponseMessage) Type() Type {
	return TypeBlocksResponse
}

func (m *BlocksResponseMessage) Count() uint32 {
	return uint32(len(m.BlocksData))
}

func (m *BlocksResponseMessage) To() uint32 {
	// response message without any block
	if len(m.BlocksData) == 0 {
		return 0
	}
	return m.From + m.Count() - 1
}

func (m *BlocksResponseMessage) LastCertificateHeight() uint32 {
	if m.LastCertificate != nil {
		return m.From
	}
	return 0
}

func (m *BlocksResponseMessage) Fingerprint() string {
	return fmt.Sprintf("{⚓ %d %s %v-%v}", m.SessionID, m.ResponseCode, m.From, m.To())
}

func (m *BlocksResponseMessage) IsRequestRejected() bool {
	return m.ResponseCode == ResponseCodeRejected
}
