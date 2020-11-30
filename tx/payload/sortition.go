package payload

import (
	"fmt"

	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/errors"
)

type SortitionPayload struct {
	Address crypto.Address `cbor:"1,keyasint"`
	Index   int64          `cbor:"2,keyasint"`
	Proof   []byte         `cbor:"3,keyasint"`
}

func (p *SortitionPayload) Signer() crypto.Address {
	return p.Address
}

func (p *SortitionPayload) Value() int64 {
	return 0
}

func (p *SortitionPayload) SanityCheck() error {
	if p.Index < 0 {
		return errors.Errorf(errors.ErrInvalidTx, "Invalid index")
	}
	if len(p.Proof) != 64 {
		return errors.Errorf(errors.ErrInvalidTx, "Invalid prrof")
	}

	return nil
}

func (p *SortitionPayload) Type() PayloadType {
	return PayloadTypeSortition
}

func (p *SortitionPayload) Fingerprint() string {
	return fmt.Sprintf("{Sortiton: %v",
		p.Address.Fingerprint())
}
