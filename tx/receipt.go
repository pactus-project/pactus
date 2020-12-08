package tx

import (
	"encoding/json"

	"github.com/fxamacker/cbor/v2"
	"github.com/zarbchain/zarb-go/crypto"
)

const (
	Ok     = 0
	Failed = 1
)

// Transaction receipt
type Receipt struct {
	data receiptData
}
type receiptData struct {
	TxHash    crypto.Hash `cbor:"1,keyasint"`
	BlockHash crypto.Hash `cbor:"2,keyasint"`
	Status    int         `cbor:"3,keyasint"`
}

func (r *Receipt) TxHash() crypto.Hash    { return r.data.TxHash }
func (r *Receipt) BlockHash() crypto.Hash { return r.data.BlockHash }
func (r *Receipt) Status() int            { return r.data.Status }

func (r *Receipt) Hash() crypto.Hash {
	bz, _ := r.MarshalCBOR()
	return crypto.HashH(bz)
}

func (r *Receipt) MarshalCBOR() ([]byte, error) {
	return cbor.Marshal(r.data)
}

func (r *Receipt) UnmarshalCBOR(bs []byte) error {
	return cbor.Unmarshal(bs, &r.data)
}

func (r *Receipt) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.data)
}

func (r *Receipt) UnmarshalJSON(bs []byte) error {
	return json.Unmarshal(bs, &r.data)
}
