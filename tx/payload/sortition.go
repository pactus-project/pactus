package payload

import (
	"fmt"

	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/errors"
)

type SortitionPayload struct {
	Address crypto.Address `cbor:"1,keyasint"`
	Proof   []byte         `cbor:"2,keyasint"`
}

func (p *SortitionPayload) Type() PayloadType {
	return PayloadTypeSortition
}

func (p *SortitionPayload) Signer() crypto.Address {
	return p.Address
}

func (p *SortitionPayload) Value() int64 {
	return 0
}

func (p *SortitionPayload) SanityCheck() error {
	if err := p.Address.SanityCheck(); err != nil {
		return errors.Errorf(errors.ErrInvalidTx, "Invalid address")
	}
	if len(p.Proof) != crypto.SignatureSize {
		return errors.Errorf(errors.ErrInvalidTx, "Invalid proof")
	}

	return nil
}

func (p *SortitionPayload) Fingerprint() string {
	return fmt.Sprintf("{Sortition: %v",
		p.Address.Fingerprint())
}
