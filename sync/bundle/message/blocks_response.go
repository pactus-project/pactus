package message

import (
	"fmt"

	"github.com/pactus-project/pactus/network"
	"github.com/pactus-project/pactus/types/certificate"
)

type BlocksResponseMessage struct {
	ResponseCode    ResponseCode             `cbor:"1,keyasint"`
	SessionID       int                      `cbor:"2,keyasint"`
	From            uint32                   `cbor:"3,keyasint"`
	BlocksData      [][]byte                 `cbor:"4,keyasint"`
	LastCertificate *certificate.Certificate `cbor:"5,keyasint"`
	Reason          string                   `cbor:"6,keyasint"`
}

func NewBlocksResponseMessage(code ResponseCode, reason string, sid int, from uint32,
	blocksData [][]byte, lastCert *certificate.Certificate,
) *BlocksResponseMessage {
	return &BlocksResponseMessage{
		ResponseCode:    code,
		SessionID:       sid,
		From:            from,
		BlocksData:      blocksData,
		LastCertificate: lastCert,
		Reason:          reason,
	}
}

func (m *BlocksResponseMessage) BasicCheck() error {
	if m.LastCertificate != nil {
		if err := m.LastCertificate.BasicCheck(); err != nil {
			return err
		}
	}

	return nil
}

func (*BlocksResponseMessage) Type() Type {
	return TypeBlocksResponse
}

func (*BlocksResponseMessage) TopicID() network.TopicID {
	return network.TopicIDUnspecified
}

func (*BlocksResponseMessage) ShouldBroadcast() bool {
	return false
}

func (*BlocksResponseMessage) ConsensusHeight() uint32 {
	return 0
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

// LogString returns a concise string representation intended for use in logs.
func (m *BlocksResponseMessage) LogString() string {
	return fmt.Sprintf("{⚓ %d %s %v-%v}", m.SessionID, m.ResponseCode, m.From, m.To())
}

func (m *BlocksResponseMessage) IsRequestRejected() bool {
	return m.ResponseCode == ResponseCodeRejected
}
