package http

import (
	"time"

	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/crypto"
)

func bytesToHash(bs []byte, err error) crypto.Hash {
	h, _ := crypto.HashFromRawBytes(bs)
	return h
}

func bytesToAddress(bs []byte, err error) crypto.Address {
	a, _ := crypto.AddressFromRawBytes(bs)
	return a
}

func bytesToSignature(bs []byte, err error) crypto.Signature {
	sig, _ := crypto.SignatureFromRawBytes(bs)
	return sig
}

type BlockInfo struct {
	Hash   crypto.Hash
	Height int
	Data   string
	Time   time.Time
	Block  block.Block
}
