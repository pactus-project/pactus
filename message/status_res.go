package message

import (
	"fmt"

	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/version"
)

type StatusResPayload struct {
	Version         version.Version `cbor:"1,keyasint"`
	LastBlockHeight int             `cbor:"2,keyasint"`
	LastBlockHash   crypto.Hash     `cbor:"3,keyasint"`
}

func NewStatusResMessage(height int, hash crypto.Hash) Message {
	return Message{
		Type: PayloadTypeStatusRes,
		Payload: &StatusResPayload{
			Version:         version.NodeVersion,
			LastBlockHeight: height,
			LastBlockHash:   hash,
		},
	}

}
func (p *StatusResPayload) SanityCheck() error {
	if p.LastBlockHeight < 0 {
		return errors.Errorf(errors.ErrInvalidMessage, "invalid Height")
	}
	if err := p.LastBlockHash.SanityCheck(); err != nil {
		return errors.Errorf(errors.ErrInvalidMessage, "Invalid hash: %v", err)
	}
	return nil
}

func (p *StatusResPayload) Type() PayloadType {
	return PayloadTypeStatusRes
}

func (p *StatusResPayload) Fingerprint() string {
	return fmt.Sprintf("{%v}", p.LastBlockHeight)
}
