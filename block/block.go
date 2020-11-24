package block

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/fxamacker/cbor/v2"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/logger"
	"github.com/zarbchain/zarb-go/tx"
	"github.com/zarbchain/zarb-go/util"
)

type Block struct {
	data blockData
}

type blockData struct {
	Header   Header   `cbor:"1,keyasint"`
	TxHashes TxHashes `cbor:"2,keyasint"`
}

func NewBlock(header Header, txHashes TxHashes) (*Block, error) {
	b := &Block{
		data: blockData{
			Header:   header,
			TxHashes: txHashes,
		},
	}

	if err := b.SanityCheck(); err != nil {
		return nil, errors.Error(errors.ErrInvalidBlock)
	}
	return b, nil
}

func MakeBlock(timestamp time.Time, txHashes TxHashes,
	lastBlockHash, CommitersHash, stateHash, lastReceiptsHash crypto.Hash,
	lastCommit *Commit, proposer crypto.Address) Block {

	txsHash := txHashes.Hash()
	header := NewHeader(1, timestamp,
		txsHash, lastBlockHash, CommitersHash, stateHash, lastReceiptsHash, proposer, lastCommit)

	b := Block{
		data: blockData{
			Header:   header,
			TxHashes: txHashes,
		},
	}

	if err := b.SanityCheck(); err != nil {
		logger.Panic("Invalid block information", "err", err)
	}
	return b
}

func (b Block) Header() *Header     { return &b.data.Header }
func (b Block) TxHashes() *TxHashes { return &b.data.TxHashes }

func (b Block) SanityCheck() error {
	if err := b.data.Header.SanityCheck(); err != nil {
		return err
	}
	if !b.data.Header.TxsHash().EqualsTo(b.data.TxHashes.Hash()) {
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

func (b Block) Fingerprint() string {
	return fmt.Sprintf("{⌘ %v 👤 %v 💻 %v 👥 %v 📨 %d}",
		b.Hash().Fingerprint(),
		b.data.Header.ProposerAddress().Fingerprint(),
		b.data.Header.StateHash().Fingerprint(),
		b.data.Header.CommitersHash().Fingerprint(),
		b.data.TxHashes.Count(),
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

// ---------
// For tests
func GenerateTestBlock(proposer *crypto.Address) (Block, []*tx.Tx) {
	if proposer == nil {
		addr, _, _ := crypto.GenerateTestKeyPair()
		proposer = &addr
	}
	txs := make([]*tx.Tx, 0)
	txs = append(txs, tx.GenerateTestSendTx())
	txs = append(txs, tx.GenerateTestSendTx())
	txs = append(txs, tx.GenerateTestSendTx())
	txs = append(txs, tx.GenerateTestSendTx())

	txHashes := NewTxHashes()
	for _, tx := range txs {
		txHashes.Append(tx.Hash())
	}
	lastBlockHash := crypto.GenerateTestHash()
	addr1, _, pv1 := crypto.GenerateTestKeyPair()
	addr2, _, pv2 := crypto.GenerateTestKeyPair()
	addr3, _, pv3 := crypto.GenerateTestKeyPair()
	addr4, _, _ := crypto.GenerateTestKeyPair()

	sigs := []crypto.Signature{
		*pv1.Sign(lastBlockHash.RawBytes()),
		*pv2.Sign(lastBlockHash.RawBytes()),
		*pv3.Sign(lastBlockHash.RawBytes()),
	}
	sig := crypto.Aggregate(sigs)

	commit := NewCommit(util.RandInt(10),
		[]Commiter{
			Commiter{Signed: true, Address: addr1},
			Commiter{Signed: true, Address: addr2},
			Commiter{Signed: true, Address: addr3},
			Commiter{Signed: false, Address: addr4},
		},
		sig)

	block := MakeBlock(time.Now(), txHashes,
		lastBlockHash,
		crypto.GenerateTestHash(),
		crypto.GenerateTestHash(),
		crypto.GenerateTestHash(),
		commit,
		*proposer)

	return block, txs
}
