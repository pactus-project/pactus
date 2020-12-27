package payload

import (
	"fmt"

	"github.com/zarbchain/zarb-go/consensus/hrs"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/errors"
)

type HeartBeatPayload struct {
	Pulse         hrs.HRS     `cbor:"1,keyasint"`
	LastBlockHash crypto.Hash `cbor:"2,keyasint"`
}

func (p *HeartBeatPayload) SanityCheck() error {
	if !p.Pulse.IsValid() {
		return errors.Errorf(errors.ErrInvalidMessage, "Invalid step")
	}
	return nil
}

func (p *HeartBeatPayload) Type() PayloadType {
	return PayloadTypeHeartBeat
}

func (p *HeartBeatPayload) Fingerprint() string {
	return fmt.Sprintf("{%s}", p.Pulse.String())
}
