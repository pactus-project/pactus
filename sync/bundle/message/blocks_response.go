package message

import (
	"fmt"

	"github.com/pactus-project/pactus/types/certificate"
	"github.com/pactus-project/pactus/util/errors"
)

type BlocksResponseMessage struct {
	ResponseCode        ResponseCode             `cbor:"1,keyasint"`
	SessionID           int                      `cbor:"2,keyasint"`
	From                uint32                   `cbor:"3,keyasint"`
	CommittedBlocksData [][]byte                 `cbor:"4,keyasint"`
	LastCertificate     *certificate.Certificate `cbor:"5,keyasint"`
	Reason              string                   `cbor:"6,keyasint"`
}

func NewBlocksResponseMessage(code ResponseCode, reason string, sid int, from uint32,
	blocksData [][]byte, lastCert *certificate.Certificate,
) *BlocksResponseMessage {
	return &BlocksResponseMessage{
		ResponseCode:        code,
		SessionID:           sid,
		From:                from,
		CommittedBlocksData: blocksData,
		LastCertificate:     lastCert,
		Reason:              reason,
	}
}

func (m *BlocksResponseMessage) BasicCheck() error {
	if m.From == 0 && len(m.CommittedBlocksData) != 0 {
		return errors.Errorf(errors.ErrInvalidHeight, "unexpected block for height zero")
	}
	if m.LastCertificate != nil {
		if err := m.LastCertificate.BasicCheck(); err != nil {
			return err
		}
	}

	return nil
}

func (m *BlocksResponseMessage) Type() Type {
	return TypeBlocksResponse
}

func (m *BlocksResponseMessage) Count() uint32 {
	return uint32(len(m.CommittedBlocksData))
}

func (m *BlocksResponseMessage) To() uint32 {
	// response message without any block
	if len(m.CommittedBlocksData) == 0 {
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

func (m *BlocksResponseMessage) String() string {
	return fmt.Sprintf("{⚓ %d %s %v-%v}", m.SessionID, m.ResponseCode, m.From, m.To())
}

func (m *BlocksResponseMessage) IsRequestRejected() bool {
	return m.ResponseCode == ResponseCodeRejected
}
