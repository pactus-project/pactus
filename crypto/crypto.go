package crypto

import (
	"github.com/herumi/bls-go-binary/bls"
)

var MintbaseAddress = Address{data: addressData{Address: [AddressSize]byte{0}}}

func init() {
	bls.Init(bls.BLS12_381)
}
