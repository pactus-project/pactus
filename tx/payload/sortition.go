package payload

import (
	"fmt"
	"io"

	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/encoding"
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/sortition"
)

type SortitionPayload struct {
	Address crypto.Address
	Proof   sortition.Proof
}

func (p *SortitionPayload) Type() Type {
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
		return errors.Errorf(errors.ErrInvalidTx, "invalid address")
	}

	return nil
}

func (p *SortitionPayload) SerializeSize() int {
	return 69 //48+21
}

func (p *SortitionPayload) Encode(w io.Writer) error {
	return encoding.WriteElements(w, &p.Address, &p.Proof)
}

func (p *SortitionPayload) Decode(r io.Reader) error {
	return encoding.ReadElements(r, &p.Address, &p.Proof)
}

func (p *SortitionPayload) Fingerprint() string {
	return fmt.Sprintf("{Sortition 🎯 %v",
		p.Address.Fingerprint())
}
