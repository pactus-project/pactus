package block

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/fxamacker/cbor/v2"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/sortition"
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

func MakeBlock(version int, timestamp time.Time, txIDs TxIDs,
	lastBlockHash, committeeHash, stateHash, lastReceiptsHash crypto.Hash,
	lastCommit *Commit, sortitionSeed sortition.Seed, proposer crypto.Address) Block {

	txIDsHash := txIDs.Hash()
	header := NewHeader(version, timestamp,
		txIDsHash, lastBlockHash, committeeHash, stateHash, lastReceiptsHash, lastCommit.Hash(), sortitionSeed, proposer)

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

func (b Block) Header() Header      { return b.data.Header }
func (b Block) LastCommit() *Commit { return b.data.LastCommit }
func (b Block) TxIDs() TxIDs        { return b.data.TxIDs }

func (b Block) SanityCheck() error {
	if err := b.data.Header.SanityCheck(); err != nil {
		return err
	}
	if b.data.TxIDs.Len() == 0 {
		return errors.Errorf(errors.ErrInvalidBlock, "Block at least should have one transaction")
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
		b.data.Header.CommitteeHash().Fingerprint(),
		b.data.TxIDs.Len(),
	)
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
func GenerateTestBlock(proposer *crypto.Address, lastBlockHash *crypto.Hash) (*Block, []*tx.Tx) {
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
	if lastBlockHash == nil {
		h := crypto.GenerateTestHash()
		lastBlockHash = &h
	}
	lastReceiptsHash := crypto.GenerateTestHash()
	commit := GenerateTestCommit(*lastBlockHash)
	if lastBlockHash.IsUndef() {
		commit = nil
		lastReceiptsHash = crypto.UndefHash
	}
	sortitionSeed := sortition.GenerateRandomSeed()
	block := MakeBlock(1, util.Now(), ids,
		*lastBlockHash,
		crypto.GenerateTestHash(),
		crypto.GenerateTestHash(),
		lastReceiptsHash,
		commit,
		sortitionSeed,
		*proposer)

	return &block, txs
}

func GenerateTestCommit(blockhash crypto.Hash) *Commit {
	_, _, priv2 := crypto.GenerateTestKeyPair()
	_, _, priv3 := crypto.GenerateTestKeyPair()
	_, _, priv4 := crypto.GenerateTestKeyPair()

	sigs := []crypto.Signature{
		priv2.Sign(blockhash.RawBytes()),
		priv3.Sign(blockhash.RawBytes()),
		priv4.Sign(blockhash.RawBytes()),
	}
	sig := crypto.Aggregate(sigs)

	return NewCommit(blockhash, util.RandInt(10),
		[]Committer{
			{Status: CommitNotSigned, Number: 0},
			{Status: CommitSigned, Number: 1},
			{Status: CommitSigned, Number: 2},
			{Status: CommitSigned, Number: 3},
		},
		sig)
}
