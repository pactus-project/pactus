package payload

import (
	"fmt"

	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/errors"
)

type SendPayload struct {
	Sender   crypto.Address `cbor:"1,keyasint"`
	Receiver crypto.Address `cbor:"2,keyasint"`
	Amount   int64          `cbor:"3,keyasint"`
}

func (p *SendPayload) Type() PayloadType {
	return PayloadTypeSend
}

func (p *SendPayload) Signer() crypto.Address {
	return p.Sender
}

func (p *SendPayload) Value() int64 {
	return p.Amount
}

func (p *SendPayload) SanityCheck() error {
	if p.Amount < 0 {
		return errors.Errorf(errors.ErrInvalidTx, "invalid amount")
	}
	if err := p.Receiver.SanityCheck(); err != nil {
		return errors.Errorf(errors.ErrInvalidTx, "invalid receiver address")
	}

	return nil
}

func (p *SendPayload) Fingerprint() string {
	return fmt.Sprintf("{Send: %v->%v 💸 %v",
		p.Sender.Fingerprint(),
		p.Receiver.Fingerprint(),
		p.Amount)
}
