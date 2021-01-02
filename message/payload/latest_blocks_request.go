package payload

import (
	"fmt"

	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/errors"
)

type LatestBlocksRequestPayload struct {
	From          int         `cbor:"1,keyasint"`
	LastBlockHash crypto.Hash `cbor:"2,keyasint"`
}

func (p *LatestBlocksRequestPayload) SanityCheck() error {
	if p.From <= 0 {
		return errors.Errorf(errors.ErrInvalidMessage, "invalid Height")
	}
	return nil
}

func (p *LatestBlocksRequestPayload) Type() PayloadType {
	return PayloadTypeLatestBlocksRequest
}

func (p *LatestBlocksRequestPayload) Fingerprint() string {
	return fmt.Sprintf("{%v ⌘ %v}", p.From, p.LastBlockHash.Fingerprint())
}
