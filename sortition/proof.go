package sortition

import (
	"crypto/rand"

	"github.com/zarbchain/zarb-go/util"
)

type Proof struct {
	Base []byte `cbor:"1,keyasint"`
	Coin int    `cbor:"2,keyasint"`
}

func NewProof(base []byte, coin int) Proof {
	return Proof{
		Base: base,
		Coin: coin,
	}
}

func GenerateRandomProof() Proof {
	p := Proof{
		Base: util.RandomSlice(48),
	}
	_, err := rand.Read(p.Base[:])
	if err != nil {
		panic(err)
	}
	p.Coin = util.RandInt(21 * 1e6) // 21 milion coin
	return p
}
