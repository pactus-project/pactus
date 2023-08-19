package message

import (
	"fmt"

	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/util/errors"
)

type HeartBeatMessage struct {
	Height        uint32    `cbor:"1,keyasint"`
	Round         int16     `cbor:"2,keyasint"`
	PrevBlockHash hash.Hash `cbor:"3,keyasint"`
}

func NewHeartBeatMessage(h uint32, r int16, hash hash.Hash) *HeartBeatMessage {
	return &HeartBeatMessage{
		Height:        h,
		Round:         r,
		PrevBlockHash: hash,
	}
}

func (m *HeartBeatMessage) BasicCheck() error {
	if m.Height == 0 {
		return errors.Errorf(errors.ErrInvalidHeight, "invalid height")
	}
	return nil
}

func (m *HeartBeatMessage) Type() Type {
	return TypeHeartBeat
}

func (m *HeartBeatMessage) String() string {
	return fmt.Sprintf("{%d/%d}", m.Height, m.Round)
}
