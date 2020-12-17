package block

import (
	"encoding/json"

	"github.com/fxamacker/cbor/v2"
	"github.com/zarbchain/zarb-go/crypto"
	simpleMerkle "github.com/zarbchain/zarb-go/libs/merkle"
)

type TxIDs struct {
	data txIDsData
}

type txIDsData struct {
	IDs []crypto.Hash `cbor:"1,keyasint"`
}

func NewTxIDs() TxIDs {
	return TxIDs{
		data: txIDsData{
			IDs: make([]crypto.Hash, 0),
		},
	}
}
func (txs *TxIDs) Append(hash crypto.Hash) {
	txs.data.IDs = append(txs.data.IDs, hash)
}

func (txs TxIDs) Hash() crypto.Hash {
	merkle := simpleMerkle.NewTreeFromHashes(txs.data.IDs)
	return merkle.Root()
}

func (txs TxIDs) IDs() []crypto.Hash {
	return txs.data.IDs
}

func (txs TxIDs) IsEmpty() bool {
	return txs.Count() == 0
}

func (txs TxIDs) Count() int {
	return len(txs.data.IDs)
}

func (txs *TxIDs) MarshalCBOR() ([]byte, error) {
	return cbor.Marshal(txs.data)
}

func (txs *TxIDs) UnmarshalCBOR(bs []byte) error {
	return cbor.Unmarshal(bs, &txs.data)
}

func (txs TxIDs) MarshalJSON() ([]byte, error) {
	return json.Marshal(txs.data)
}

func (txs *TxIDs) UnmarshalJSON(bz []byte) error {
	return json.Unmarshal(bz, &txs.data)
}
