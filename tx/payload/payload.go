package payload

import (
	"fmt"

	"github.com/zarbchain/zarb-go/crypto"
)

type PayloadType int

const (
	PayloadTypeSend      = PayloadType(1)
	PayloadTypeBond      = PayloadType(2)
	PayloadTypeUnbond    = PayloadType(3)
	PayloadTypeSortition = PayloadType(4)
)

func (t PayloadType) String() string {
	switch t {
	case PayloadTypeSend:
		return "send"
	case PayloadTypeBond:
		return "bond"
	case PayloadTypeUnbond:
		return "unbond"
	case PayloadTypeSortition:
		return "sortition"
	}
	return fmt.Sprintf("%d", t)
}

type Payload interface {
	Signer() crypto.Address
	Value() int64
	Type() PayloadType
	SanityCheck() error
	Fingerprint() string
}
