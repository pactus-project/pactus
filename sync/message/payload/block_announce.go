package payload

import (
	"fmt"

	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/errors"
)

type BlockAnnouncePayload struct {
	PeerID      peer.ID            `cbor:"1,keyasint"`
	Height      int                `cbor:"2,keyasint"`
	Block       *block.Block       `cbor:"3,keyasint"`
	Certificate *block.Certificate `cbor:"4,keyasint"`
}

func (p *BlockAnnouncePayload) SanityCheck() error {
	if p.Height < 0 {
		return errors.Errorf(errors.ErrInvalidMessage, "Invalid height")
	}
	if err := p.Block.SanityCheck(); err != nil {
		return errors.Errorf(errors.ErrInvalidMessage, "Invalid block: %v", err)
	}
	if err := p.Certificate.SanityCheck(); err != nil {
		return errors.Errorf(errors.ErrInvalidMessage, "Invalid commit: %v", err)
	}
	if err := p.PeerID.Validate(); err != nil {
		return errors.Errorf(errors.ErrInvalidMessage, "Invalid querier peer id: %v", err)
	}

	return nil
}

func (p *BlockAnnouncePayload) Type() PayloadType {
	return PayloadTypeBlockAnnounce
}

func (p *BlockAnnouncePayload) Fingerprint() string {
	return fmt.Sprintf("{⌘ %v}", p.Block.Hash().Fingerprint())
}
