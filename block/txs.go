package block

import (
	"encoding/json"

	"github.com/fxamacker/cbor/v2"
	"github.com/zarbchain/zarb-go/crypto"
	simpleMerkle "github.com/zarbchain/zarb-go/libs/merkle"
	"github.com/zarbchain/zarb-go/tx"
)

type TxIDs struct {
	data txIDsData
}

type txIDsData struct {
	IDs []tx.ID `cbor:"1,keyasint"`
}

func NewTxIDs() TxIDs {
	return TxIDs{
		data: txIDsData{
			IDs: make([]tx.ID, 0),
		},
	}
}
func (txs *TxIDs) Append(id tx.ID) {
	txs.data.IDs = append(txs.data.IDs, id)
}

func (txs *TxIDs) Prepend(id tx.ID) {
	ids := make([]tx.ID, len(txs.data.IDs)+1)
	ids[0] = id
	copy(ids[1:], txs.data.IDs)
	txs.data.IDs = ids
}

func (txs TxIDs) Hash() crypto.Hash {
	merkle := simpleMerkle.NewTreeFromHashes(txs.data.IDs)
	return merkle.Root()
}

func (txs TxIDs) IDs() []tx.ID {
	return txs.data.IDs
}

func (txs TxIDs) IsEmpty() bool {
	return txs.Len() == 0
}

func (txs TxIDs) Len() int {
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
