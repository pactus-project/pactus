package tx

import (
	"encoding/json"

	"github.com/fxamacker/cbor/v2"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/errors"
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
	Status    int         `cbor:"1,keyasint"`
	TxID      ID          `cbor:"2,keyasint"`
	BlockHash crypto.Hash `cbor:"3,keyasint"`
}

func (r *Receipt) Status() int            { return r.data.Status }
func (r *Receipt) TxID() ID               { return r.data.TxID }
func (r *Receipt) BlockHash() crypto.Hash { return r.data.BlockHash }

func (r *Receipt) Hash() crypto.Hash {
	bz, _ := r.MarshalCBOR()
	return crypto.HashH(bz)
}

func (r *Receipt) SanityCheck() error {
	if r.data.Status != Ok {
		return errors.Errorf(errors.ErrInvalidTx, "invalid status")
	}
	if err := r.data.BlockHash.SanityCheck(); err != nil {
		return errors.Errorf(errors.ErrInvalidTx, "invalid block hash")
	}
	if err := r.data.TxID.SanityCheck(); err != nil {
		return errors.Errorf(errors.ErrInvalidTx, "invalid block hash")
	}
	return nil
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

func (r *Receipt) Encode() ([]byte, error) {
	return cbor.Marshal(r.data)
}

func (r *Receipt) Decode(bs []byte) error {
	return cbor.Unmarshal(bs, &r.data)
}
