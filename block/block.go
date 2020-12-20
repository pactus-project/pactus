package block

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/fxamacker/cbor/v2"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/tx"
	"github.com/zarbchain/zarb-go/util"
)

type Block struct {
	data blockData

	memorizedHash *crypto.Hash
}

type blockData struct {
	Header     Header  `cbor:"1,keyasint"`
	LastCommit *Commit `cbor:"2,keyasint"`
	TxIDs      TxIDs   `cbor:"3,keyasint"`
}

func NewBlock(header Header, lastCommit *Commit, TxIDs TxIDs) (*Block, error) {
	b := &Block{
		data: blockData{
			Header:     header,
			LastCommit: lastCommit,
			TxIDs:      TxIDs,
		},
	}

	if err := b.SanityCheck(); err != nil {
		return nil, errors.Error(errors.ErrInvalidBlock)
	}
	return b, nil
}

func MakeBlock(timestamp time.Time, txIDs TxIDs,
	lastBlockHash, CommittersHash, stateHash, lastReceiptsHash crypto.Hash,
	lastCommit *Commit, proposer crypto.Address) Block {

	txIDsHash := txIDs.Hash()
	header := NewHeader(1, timestamp,
		txIDsHash, lastBlockHash, CommittersHash, stateHash, lastReceiptsHash, lastCommit.Hash(), proposer)

	b := Block{
		data: blockData{
			Header:     header,
			LastCommit: lastCommit,
			TxIDs:      txIDs,
		},
	}

	if err := b.SanityCheck(); err != nil {
		panic(err)
	}
	return b
}

func (b Block) Header() *Header     { return &b.data.Header }
func (b Block) LastCommit() *Commit { return b.data.LastCommit }
func (b Block) TxIDs() *TxIDs       { return &b.data.TxIDs }

func (b Block) SanityCheck() error {
	if err := b.data.Header.SanityCheck(); err != nil {
		return err
	}
	if !b.data.Header.TxIDsHash().EqualsTo(b.data.TxIDs.Hash()) {
		return errors.Errorf(errors.ErrInvalidBlock, "Invalid Txs Hash")
	}
	if b.data.LastCommit != nil {
		if err := b.data.LastCommit.SanityCheck(); err != nil {
			return err
		}
		if !b.data.Header.LastCommitHash().EqualsTo(b.data.LastCommit.Hash()) {
			return errors.Errorf(errors.ErrInvalidBlock, "Invalid Last Commit Hash")
		}
	} else {
		// Check for genesis block
		if !b.data.Header.LastCommitHash().IsUndef() {
			return errors.Errorf(errors.ErrInvalidBlock, "Invalid genesis block hash")
		}
	}

	return nil
}

func (b Block) Hash() crypto.Hash {
	if b.memorizedHash == nil {
		h := b.data.Header.Hash()
		b.memorizedHash = &h
	}

	return *b.memorizedHash
}

func (b Block) HashesTo(hash crypto.Hash) bool {
	return b.Hash().EqualsTo(hash)
}

func (b Block) Fingerprint() string {
	return fmt.Sprintf("{⌘ %v 👤 %v 💻 %v 👥 %v 📨 %d}",
		b.Hash().Fingerprint(),
		b.data.Header.ProposerAddress().Fingerprint(),
		b.data.Header.StateHash().Fingerprint(),
		b.data.Header.CommittersHash().Fingerprint(),
		b.data.TxIDs.Len(),
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
func GenerateTestBlock(proposer *crypto.Address) (*Block, []*tx.Tx) {
	if proposer == nil {
		addr, _, _ := crypto.GenerateTestKeyPair()
		proposer = &addr
	}
	txs := make([]*tx.Tx, 0)
	tx1, _ := tx.GenerateTestSendTx()
	tx2, _ := tx.GenerateTestSendTx()
	tx3, _ := tx.GenerateTestSendTx()
	tx4, _ := tx.GenerateTestSendTx()

	txs = append(txs, tx1)
	txs = append(txs, tx2)
	txs = append(txs, tx3)
	txs = append(txs, tx4)

	ids := NewTxIDs()
	for _, tx := range txs {
		ids.Append(tx.ID())
	}
	lastBlockHash := crypto.GenerateTestHash()
	commit := GenerateTestCommit(lastBlockHash)

	block := MakeBlock(time.Now(), ids,
		lastBlockHash,
		crypto.GenerateTestHash(),
		crypto.GenerateTestHash(),
		crypto.GenerateTestHash(),
		commit,
		*proposer)

	return &block, txs
}

func GenerateTestCommit(blockhash crypto.Hash) *Commit {

	addr1, _, pv1 := crypto.GenerateTestKeyPair()
	addr2, _, pv2 := crypto.GenerateTestKeyPair()
	addr3, _, pv3 := crypto.GenerateTestKeyPair()
	addr4, _, _ := crypto.GenerateTestKeyPair()

	sigs := []*crypto.Signature{
		pv1.Sign(blockhash.RawBytes()),
		pv2.Sign(blockhash.RawBytes()),
		pv3.Sign(blockhash.RawBytes()),
	}
	sig := crypto.Aggregate(sigs)

	return NewCommit(util.RandInt(10),
		[]Committer{
			{Status: CommitSigned, Address: addr1},
			{Status: CommitSigned, Address: addr2},
			{Status: CommitSigned, Address: addr3},
			{Status: CommitNotSigned, Address: addr4},
		},
		sig)
}
