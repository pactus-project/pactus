package payload

import (
	"fmt"
	"io"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/types/amount"
)

type Type uint8

const (
	TypeTransfer  = Type(1)
	TypeBond      = Type(2)
	TypeSortition = Type(3)
	TypeUnbond    = Type(4)
	TypeWithdraw  = Type(5)
)

func (t Type) String() string {
	switch t {
	case TypeTransfer:
		return "transfer"
	case TypeBond:
		return "bond"
	case TypeUnbond:
		return "unbond"
	case TypeWithdraw:
		return "withdraw"
	case TypeSortition:
		return "sortition"
	}

	return fmt.Sprintf("%d", t)
}

type Payload interface {
	Signer() crypto.Address
	Value() amount.Amount
	Type() Type
	SerializeSize() int
	Encode(io.Writer) error
	Decode(io.Reader) error
	BasicCheck() error
	String() string
	Receiver() *crypto.Address
}
