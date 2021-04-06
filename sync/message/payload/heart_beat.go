package payload

import (
	"fmt"

	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/errors"
)

type HeartBeatPayload struct {
	Height        int         `cbor:"1,keyasint"`
	Round         int         `cbor:"2,keyasint"`
	LastBlockHash crypto.Hash `cbor:"3,keyasint"`
}

func NewHeartBeatPayload(h, r int, hash crypto.Hash) Payload {
	return &HeartBeatPayload{
		Height:        h,
		Round:         r,
		LastBlockHash: hash,
	}
}

func (p *HeartBeatPayload) SanityCheck() error {
	if p.Height <= 0 {
		return errors.Errorf(errors.ErrInvalidMessage, "Invalid height")
	}
	if p.Round < 0 {
		return errors.Errorf(errors.ErrInvalidMessage, "Invalid round")
	}
	return nil
}

func (p *HeartBeatPayload) Type() PayloadType {
	return PayloadTypeHeartBeat
}

func (p *HeartBeatPayload) Fingerprint() string {
	return fmt.Sprintf("{%d/%d}", p.Height, p.Round)
}
