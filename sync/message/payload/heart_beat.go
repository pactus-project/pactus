package payload

import (
	"fmt"

	"github.com/zarbchain/zarb-go/crypto/hash"
	"github.com/zarbchain/zarb-go/errors"
)

type HeartBeatPayload struct {
	Height        int       `cbor:"1,keyasint"`
	Round         int       `cbor:"2,keyasint"`
	PrevBlockHash hash.Hash `cbor:"3,keyasint"`
}

func NewHeartBeatPayload(h, r int, hash hash.Hash) *HeartBeatPayload {
	return &HeartBeatPayload{
		Height:        h,
		Round:         r,
		PrevBlockHash: hash,
	}
}

func (p *HeartBeatPayload) SanityCheck() error {
	if p.Height <= 0 {
		return errors.Errorf(errors.ErrInvalidMessage, "invalid height")
	}
	if p.Round < 0 {
		return errors.Errorf(errors.ErrInvalidMessage, "invalid round")
	}
	return nil
}

func (p *HeartBeatPayload) Type() Type {
	return PayloadTypeHeartBeat
}

func (p *HeartBeatPayload) Fingerprint() string {
	return fmt.Sprintf("{%d/%d}", p.Height, p.Round)
}
