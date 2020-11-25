package payload

import (
	"fmt"

	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/errors"
)

type SendPayload struct {
	SenderAddr   crypto.Address `cbor:"1,keyasint"`
	ReceiverAddr crypto.Address `cbor:"2,keyasint"`
	Amount       int64          `cbor:"3,keyasint"`
}

func (p *SendPayload) Signer() crypto.Address {
	return p.SenderAddr
}

func (p *SendPayload) SanityCheck() error {
	if p.Amount < 0 {
		return errors.Errorf(errors.ErrInvalidTx, "Invalid amount")
	}
	if err := p.SenderAddr.SanityCheck(); err != nil {
		return errors.Errorf(errors.ErrInvalidTx, "Invalid sender address")
	}
	if err := p.ReceiverAddr.SanityCheck(); err != nil {
		return errors.Errorf(errors.ErrInvalidTx, "Invalid receiver address")
	}

	return nil
}

func (p *SendPayload) Type() PayloadType {
	return PayloadTypeSend
}

func (p *SendPayload) Fingerprint() string {
	return fmt.Sprintf("{Send: %v->%v 🪙 %v",
		p.SenderAddr.Fingerprint(),
		p.ReceiverAddr.Fingerprint(),
		p.Amount)
}
