package payload

import (
	"fmt"

	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/zarbchain/zarb-go/consensus/hrs"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/errors"
)

type HeartBeatPayload struct {
	PeerID        peer.ID     `cbor:"1,keyasint"`
	Pulse         hrs.HRS     `cbor:"2,keyasint"`
	LastBlockHash crypto.Hash `cbor:"3,keyasint"`
}

func (p *HeartBeatPayload) SanityCheck() error {
	if err := p.PeerID.Validate(); err != nil {
		return errors.Errorf(errors.ErrInvalidMessage, "Invalid peer id: %v", err)
	}
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
