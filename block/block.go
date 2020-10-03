package block

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/fxamacker/cbor/v2"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/logger"
)

type Block struct {
	data blockData
}

type blockData struct {
	Header     Header  `cbor:"1,keyasint"`
	Txs        Txs     `cbor:"2,keyasint"`
	LastCommit *Commit `cbor:"3,keyasint"`
}

func NewBlock(header Header, txs Txs, lastCommit *Commit) (*Block, error) {
	b := &Block{
		data: blockData{
			Header:     header,
			Txs:        txs,
			LastCommit: lastCommit,
		},
	}

	if err := b.SanityCheck(); err != nil {
		return nil, errors.Error(errors.ErrInvalidBlock)
	}
	return b, nil
}

func MakeBlock(timestamp time.Time, txs Txs,
	lastBlockHash, nextValHash, stateHash, lastReceiptsHash crypto.Hash,
	lastCommit *Commit, proposer crypto.Address) Block {

	txsHash := txs.Hash()
	LastCommitHash := lastCommit.Hash()
	header := NewHeader(1, timestamp,
		txsHash, lastBlockHash, nextValHash, stateHash, lastReceiptsHash, LastCommitHash, proposer)

	b := Block{
		data: blockData{
			Header:     header,
			Txs:        txs,
			LastCommit: lastCommit,
		},
	}

	if err := b.SanityCheck(); err != nil {
		logger.Panic("Invalid block information", "err", err)
	}
	return b
}

func (b *Block) Header() *Header     { return &b.data.Header }
func (b *Block) Txs() *Txs           { return &b.data.Txs }
func (b *Block) LastCommit() *Commit { return b.data.LastCommit }

func (b Block) SanityCheck() error {
	if err := b.data.Header.SanityCheck(); err != nil {
		return err
	}

	if b.data.LastCommit != nil {
		if err := b.data.LastCommit.SanityCheck(); err != nil {
			return err
		}
	} else {
		// Check for genesis block
		if !b.data.Header.LastBlockHash().IsUndef() ||
			!b.data.Header.LastCommitHash().IsUndef() ||
			!b.data.Header.LastReceiptsHash().IsUndef() {
			return errors.Errorf(errors.ErrInvalidBlock, "Invalid Genesis block")
		}
	}

	if b.data.Header.LastCommitHash() != b.data.LastCommit.Hash() {
		return errors.Errorf(errors.ErrInvalidBlock, "Invalid LastCommitHash")
	}

	if b.data.Header.TxsHash() != b.data.Txs.Hash() {
		return errors.Errorf(errors.ErrInvalidBlock, "Invalid Txs Hash")
	}

	return nil
}

func (b Block) Hash() crypto.Hash {
	return b.data.Header.Hash()
}

func (b Block) HashesTo(hash crypto.Hash) bool {
	return b.Hash().EqualsTo(hash)
}

func (b *Block) Fingerprint() string {
	return fmt.Sprintf("{%v by:%v Tx:%d App:%v Vals:%v}",
		b.Hash().Fingerprint(),
		b.data.Header.ProposerAddress().Fingerprint(),
		b.data.Txs.Count(),
		b.data.Header.StateHash().Fingerprint(),
		b.data.Header.NextValidatorsHash().Fingerprint(),
	)
}

func (b *Block) Size() int {
	bz, _ := b.Encode()
	return len(bz)
}

func (b Block) Encode() ([]byte, error) {
	bs, err := cbor.Marshal(b.data)
	if err != nil {
		return nil, err
	}
	return bs, nil
}

func (b *Block) Decode(bs []byte) error {
	return cbor.Unmarshal(bs, &b.data)
}


func (b *Block) MarshalCBOR() ([]byte, error) {
	return cbor.Marshal(b.data)
}

func (b *Block) UnmarshalCBOR(bs []byte) error {
	return cbor.Unmarshal(bs, &b.data)
}

func (b Block) MarshalJSON() ([]byte, error) {
	return json.Marshal(b.data)
}

func (b *Block) UnmarshalJSON(bz []byte) error {
	return json.Unmarshal(bz, &b.data)
}
