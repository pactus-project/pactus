package payload

import (
	"fmt"
	"io"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/sortition"
	"github.com/pactus-project/pactus/util/encoding"
	"github.com/pactus-project/pactus/util/errors"
)

type SortitionPayload struct {
	Address crypto.Address
	Proof   sortition.Proof
}

func (p *SortitionPayload) Type() Type {
	return TypeSortition
}

func (p *SortitionPayload) Signer() crypto.Address {
	return p.Address
}

func (p *SortitionPayload) Value() int64 {
	return 0
}

func (p *SortitionPayload) BasicCheck() error {
	if err := p.Address.BasicCheck(); err != nil {
		return errors.Error(errors.ErrInvalidAddress)
	}

	return nil
}

func (p *SortitionPayload) SerializeSize() int {
	return 69 // 48+21
}

func (p *SortitionPayload) Encode(w io.Writer) error {
	return encoding.WriteElements(w, &p.Address, &p.Proof)
}

func (p *SortitionPayload) Decode(r io.Reader) error {
	return encoding.ReadElements(r, &p.Address, &p.Proof)
}

func (p *SortitionPayload) String() string {
	return fmt.Sprintf("{Sortition 🎯 %v",
		p.Address.ShortString())
}

func (p *SortitionPayload) ReceiverAddr() *crypto.Address {
	return nil
}
