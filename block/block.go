package block

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/fxamacker/cbor/v2"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/crypto/bls"
	"github.com/zarbchain/zarb-go/crypto/hash"
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/sortition"
	"github.com/zarbchain/zarb-go/tx"
	"github.com/zarbchain/zarb-go/util"
)

type Block struct {
	lk            sync.RWMutex
	memorizedHash *hash.Hash

	data blockData
}

type blockData struct {
	Header          Header       `cbor:"1,keyasint"`
	LastCertificate *Certificate `cbor:"2,keyasint"`
	TxIDs           TxIDs        `cbor:"3,keyasint"`
}

func MakeBlock(version int, timestamp time.Time, txIDs TxIDs,
	lastBlockHash, stateHash hash.Hash,
	lastCertificate *Certificate, sortitionSeed sortition.Seed, proposer crypto.Address) *Block {
	txIDsHash := txIDs.Hash()
	lastCertHash := hash.UndefHash
	if lastCertificate != nil {
		lastCertHash = lastCertificate.Hash()
	}
	header := NewHeader(version, timestamp,
		txIDsHash, lastBlockHash, stateHash, lastCertHash, sortitionSeed, proposer)

	b := &Block{
		data: blockData{
			Header:          header,
			LastCertificate: lastCertificate,
			TxIDs:           txIDs,
		},
	}

	if err := b.SanityCheck(); err != nil {
		panic(err)
	}
	return b
}

func (b *Block) Header() *Header               { return &b.data.Header }
func (b *Block) LastCertificate() *Certificate { return b.data.LastCertificate }
func (b *Block) TxIDs() TxIDs                  { return b.data.TxIDs }

func (b *Block) SanityCheck() error {
	if err := b.Header().SanityCheck(); err != nil {
		return err
	}
	if b.TxIDs().Len() == 0 {
		return errors.Errorf(errors.ErrInvalidBlock, "block at least should have one transaction")
	}
	if !b.Header().TxIDsHash().EqualsTo(b.data.TxIDs.Hash()) {
		return errors.Errorf(errors.ErrInvalidBlock, "invalid Txs Hash")
	}
	if b.LastCertificate() != nil {
		if err := b.LastCertificate().SanityCheck(); err != nil {
			return err
		}
		if !b.Header().LastCertificateHash().EqualsTo(b.LastCertificate().Hash()) {
			return errors.Errorf(errors.ErrInvalidBlock, "invalid Last Certificate hash")
		}
	} else {
		// Check for genesis block
		if !b.Header().LastCertificateHash().IsUndef() {
			return errors.Errorf(errors.ErrInvalidBlock, "invalid Last Certificate hash")
		}
	}

	return nil
}

func (b *Block) Hash() hash.Hash {
	b.lk.Lock()
	defer b.lk.Unlock()

	if b.memorizedHash == nil {
		h := b.data.Header.Hash()
		b.memorizedHash = &h
	}

	return *b.memorizedHash
}

func (b *Block) HashesTo(hash hash.Hash) bool {
	return b.Hash().EqualsTo(hash)
}

func (b *Block) Fingerprint() string {
	return fmt.Sprintf("{⌘ %v 👤 %v 💻 %v 📨 %d}",
		b.Hash().Fingerprint(),
		b.data.Header.ProposerAddress().Fingerprint(),
		b.data.Header.StateHash().Fingerprint(),
		b.data.TxIDs.Len(),
	)
}

func (b *Block) Encode() ([]byte, error) {
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

func (b *Block) MarshalJSON() ([]byte, error) {
	return json.Marshal(b.data)
}

// ---------
// For tests
func GenerateTestBlock(proposer *crypto.Address, lastBlockHash *hash.Hash) (*Block, []*tx.Tx) {
	if proposer == nil {
		addr, _, _ := bls.GenerateTestKeyPair()
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
		h := hash.GenerateTestHash()
		lastBlockHash = &h
	}
	cert := GenerateTestCertificate(*lastBlockHash)
	if lastBlockHash.IsUndef() {
		cert = nil
	}
	sortitionSeed := sortition.GenerateRandomSeed()
	block := MakeBlock(1, util.Now(), ids,
		*lastBlockHash,
		hash.GenerateTestHash(),
		cert,
		sortitionSeed,
		*proposer)

	return block, txs
}
