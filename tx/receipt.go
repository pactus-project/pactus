package tx

import (
	"encoding/json"

	"github.com/fxamacker/cbor/v2"
	"gitlab.com/zarb-chain/zarb-go/crypto"
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
	TxHash    crypto.Hash  `cbor:"1,keyasint"`
	Status    int          `cbor:"2,keyasint"`
	GasUsed   int          `cbor:"3,keyasint"`
	Output    []byte       `cbor:"4,keyasint"`
	BlockHash *crypto.Hash `cbor:"100,keyasint,omitempty"`
}

func (r *Receipt) TxHash() crypto.Hash     { return r.data.TxHash }
func (r *Receipt) Status() int             { return r.data.Status }
func (r *Receipt) GasUsed() int            { return r.data.GasUsed }
func (r *Receipt) Output() []byte          { return r.data.Output }
func (r *Receipt) BlockHash() *crypto.Hash { return r.data.BlockHash }

func (r *Receipt) SetGasUsed(gasUsed int) {
	r.data.GasUsed = gasUsed
}

func (r *Receipt) SetOutput(output []byte) {
	r.data.Output = output
}

func (r *Receipt) SetBlockHash(hash crypto.Hash) {
	r.data.BlockHash = &hash
}

func (r *Receipt) Hash() crypto.Hash {
	// Consensus receipt has no blockhash
	r2 := r
	r2.data.BlockHash = nil
	bz, _ := r2.Encode()
	return crypto.HashH(bz)
}

func (r *Receipt) Encode() ([]byte, error) {
	return cbor.Marshal(r.data)
}

func (r *Receipt) Decode(bs []byte) error {
	return cbor.Unmarshal(bs, &r.data)
}

func (r *Receipt) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.data)
}

func (r *Receipt) UnmarshalJSON(bs []byte) error {
	return json.Unmarshal(bs, &r.data)
}
