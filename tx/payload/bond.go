package payload

import (
	"fmt"

	"github.com/zarbchain/zarb-go/crypto"
)

type BondPayload struct {
}

func (p *BondPayload) Signer() crypto.Address {
	return crypto.MintbaseAddress
}

func (p *BondPayload) SanityCheck() error {

	return nil
}

func (p *BondPayload) Type() PayloadType {
	return PayloadTypeBond
}

func (p *BondPayload) Fingerprint() string {
	return fmt.Sprint("{}")
}
