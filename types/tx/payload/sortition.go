package payload

import (
	"fmt"
	"io"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/sortition"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/util/encoding"
)

type SortitionPayload struct {
	Address crypto.Address
	Proof   sortition.Proof
}

func (*SortitionPayload) Type() Type {
	return TypeSortition
}

func (p *SortitionPayload) Signer() crypto.Address {
	return p.Address
}

func (*SortitionPayload) Value() amount.Amount {
	return 0
}

// BasicCheck performs basic checks on the Sortition payload.
func (p *SortitionPayload) BasicCheck() error {
	if !p.Address.IsValidatorAddress() {
		return BasicCheckError{
			Reason: "address is not a validator address: " + p.Address.String(),
		}
	}

	return nil
}

func (*SortitionPayload) SerializeSize() int {
	return 69 // 48+21
}

func (p *SortitionPayload) Encode(w io.Writer) error {
	err := p.Address.Encode(w)
	if err != nil {
		return err
	}

	return encoding.WriteElements(w, &p.Proof)
}

func (p *SortitionPayload) Decode(_ DecodeContext, r io.Reader) error {
	return encoding.ReadElements(r, &p.Address, &p.Proof)
}

// LogString returns a concise string representation intended for use in logs.
func (p *SortitionPayload) LogString() string {
	return fmt.Sprintf("{Sortition 🎯 %s",
		p.Address.LogString())
}
