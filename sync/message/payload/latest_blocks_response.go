package payload

import (
	"fmt"

	peer "github.com/libp2p/go-libp2p-peer"
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/tx"
	"github.com/zarbchain/zarb-go/util"
)

const LatestBlocksResponseCodeOK = 0
const LatestBlocksResponseCodeNoMoreBlock = 1

type LatestBlocksResponsePayload struct {
	ResponseCode ResponseCode   `cbor:"1,keyasint"`
	SessionID    int            `cbor:"2,keyasint"`
	Initiator    peer.ID        `cbor:"3,keyasint"`
	Target       peer.ID        `cbor:"4,keyasint"`
	From         int            `cbor:"5,keyasint"`
	Blocks       []*block.Block `cbor:"6,keyasint"`
	Transactions []*tx.Tx       `cbor:"7,keyasint"`
	LastCommit   *block.Commit  `cbor:"8,keyasint"`
}

func (p *LatestBlocksResponsePayload) SanityCheck() error {
	if p.From < 0 {
		return errors.Errorf(errors.ErrInvalidMessage, "Invalid Height")
	}
	if err := p.Initiator.Validate(); err != nil {
		return errors.Errorf(errors.ErrInvalidMessage, "Invalid initiator peer is: %v", err)
	}
	if err := p.Target.Validate(); err != nil {
		return errors.Errorf(errors.ErrInvalidMessage, "Invalid target peer is: %v", err)
	}
	for _, b := range p.Blocks {
		if err := b.SanityCheck(); err != nil {
			return errors.Errorf(errors.ErrInvalidMessage, "Invalid block: %v", err)
		}
	}
	if p.LastCommit != nil {
		if err := p.LastCommit.SanityCheck(); err != nil {
			return errors.Errorf(errors.ErrInvalidMessage, "Invalid commit: %v", err)
		}
	}
	for _, trx := range p.Transactions {
		if err := trx.SanityCheck(); err != nil {
			return err
		}
	}
	return nil
}

func (p *LatestBlocksResponsePayload) Type() PayloadType {
	return PayloadTypeLatestBlocksResponse
}

func (p *LatestBlocksResponsePayload) To() int {
	if len(p.Blocks) == 0 {
		return p.From
	}
	return p.From + len(p.Blocks) - 1
}

func (p *LatestBlocksResponsePayload) Fingerprint() string {
	return fmt.Sprintf("{%s %s ⚓ %d %v-%v}", util.FingerprintPeerID(p.Initiator), p.ResponseCode, p.SessionID, p.From, p.To())
}
