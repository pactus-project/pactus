package message

import (
	"fmt"

	"github.com/zarbchain/zarb-go/consensus/hrs"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/errors"
)

type HeartBeatPayload struct {
	LastBlockHash crypto.Hash `cbor:"1,keyasint"`
	HRS           hrs.HRS     `cbor:"2,keyasint"`
	HasProposal   bool        `cbor:"3,keyasint"`
}

func NewHeartBeatMessage(lastBlockHash crypto.Hash, hrs hrs.HRS, hasProposal bool) Message {
	return Message{
		Type: PayloadTypeHeartBeat,
		Payload: &HeartBeatPayload{
			HRS:         hrs,
			HasProposal: hasProposal,
		},
	}
}

func (p *HeartBeatPayload) SanityCheck() error {
	if !p.HRS.IsValid() {
		return errors.Errorf(errors.ErrInvalidMessage, "Invalid step")
	}
	return nil
}

func (p *HeartBeatPayload) Type() PayloadType {
	return PayloadTypeHeartBeat
}

func (p *HeartBeatPayload) Fingerprint() string {
	return fmt.Sprintf("{%s}", p.HRS.Fingerprint())
}
