package block

import (
	"encoding/json"

	"github.com/fxamacker/cbor/v2"
	"gitlab.com/zarb-chain/zarb-go/crypto"
	simpleMerkle "gitlab.com/zarb-chain/zarb-go/libs/merkle"
)

type Txs struct {
	data txsData
}

type txsData struct {
	Hashes []crypto.Hash `cbor:"1,keyasint"`
}

func NewTxs() Txs {
	return Txs{
		data: txsData{
			Hashes: make([]crypto.Hash, 0),
		},
	}
}
func (txs *Txs) Append(hash crypto.Hash) {
	txs.data.Hashes = append(txs.data.Hashes, hash)
}

func (txs Txs) Hash() crypto.Hash {
	merkle := simpleMerkle.NewTreeFromHashes(txs.data.Hashes)
	root := merkle.Root()

	return *root
}

func (txs Txs) TxHashes() []crypto.Hash {
	return txs.data.Hashes
}

func (txs Txs) IsEmpty() bool {
	return txs.Count() == 0
}

func (txs Txs) Count() int {
	return len(txs.data.Hashes)
}

func (txs *Txs) MarshalCBOR() ([]byte, error) {
	return cbor.Marshal(txs.data)
}

func (txs *Txs) UnmarshalCBOR(bs []byte) error {
	return cbor.Unmarshal(bs, &txs.data)
}

func (txs Txs) MarshalJSON() ([]byte, error) {
	return json.Marshal(txs.data)
}

func (txs *Txs) UnmarshalJSON(bz []byte) error {
	return json.Unmarshal(bz, &txs.data)
}
