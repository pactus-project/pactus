package block

import (
	"encoding/json"

	"github.com/fxamacker/cbor/v2"
	"github.com/zarbchain/zarb-go/crypto"
	simpleMerkle "github.com/zarbchain/zarb-go/libs/merkle"
)

type TxHashes struct {
	data hashesData
}

type hashesData struct {
	Hashes []crypto.Hash `cbor:"1,keyasint"`
}

func NewTxHashes() TxHashes {
	return TxHashes{
		data: hashesData{
			Hashes: make([]crypto.Hash, 0),
		},
	}
}
func (txs *TxHashes) Append(hash crypto.Hash) {
	txs.data.Hashes = append(txs.data.Hashes, hash)
}

func (txs TxHashes) Hash() crypto.Hash {
	merkle := simpleMerkle.NewTreeFromHashes(txs.data.Hashes)
	root := merkle.Root()

	return *root
}

func (txs TxHashes) Hashes() []crypto.Hash {
	return txs.data.Hashes
}

func (txs TxHashes) IsEmpty() bool {
	return txs.Count() == 0
}

func (txs TxHashes) Count() int {
	return len(txs.data.Hashes)
}

func (txs *TxHashes) MarshalCBOR() ([]byte, error) {
	return cbor.Marshal(txs.data)
}

func (txs *TxHashes) UnmarshalCBOR(bs []byte) error {
	return cbor.Unmarshal(bs, &txs.data)
}

func (txs TxHashes) MarshalJSON() ([]byte, error) {
	return json.Marshal(txs.data)
}

func (txs *TxHashes) UnmarshalJSON(bz []byte) error {
	return json.Unmarshal(bz, &txs.data)
}
