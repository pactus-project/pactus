package block

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/fxamacker/cbor/v2"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/crypto/hash"
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/sortition"
	"github.com/zarbchain/zarb-go/tx"
	"github.com/zarbchain/zarb-go/util"
)

type Block struct {
	data blockData
}

type blockData struct {
	Header          Header       `cbor:"1,keyasint"`
	PrevCertificate *Certificate `cbor:"2,keyasint"`
	Transactions    Txs          `cbor:"3,keyasint"`
}

func NewBlock(header Header, prevCert *Certificate, txs Txs) *Block {
	return &Block{
		data: blockData{
			Header:          header,
			PrevCertificate: prevCert,
			Transactions:    txs,
		},
	}
}

func MakeBlock(version int, timestamp time.Time, txs Txs,
	prevBlockHash, stateRoot hash.Hash,
	prevCert *Certificate, sortitionSeed sortition.VerifiableSeed, proposer crypto.Address) *Block {
	txsRoot := txs.Root()
	prevCertHash := hash.UndefHash
	if prevCert != nil {
		prevCertHash = prevCert.Hash()
	}
	header := NewHeader(version, timestamp,
		txsRoot, stateRoot, prevBlockHash, prevCertHash, sortitionSeed, proposer)

	b := NewBlock(header, prevCert, txs)
	if err := b.SanityCheck(); err != nil {
		panic(err)
	}
	return b
}

func (b *Block) Header() *Header               { return &b.data.Header }
func (b *Block) PrevCertificate() *Certificate { return b.data.PrevCertificate }
func (b *Block) Transactions() Txs             { return b.data.Transactions }

func (b *Block) SanityCheck() error {
	if err := b.Header().SanityCheck(); err != nil {
		return err
	}
	if b.Transactions().Len() == 0 {
		return errors.Errorf(errors.ErrInvalidBlock, "block at least should have one transaction")
	}
	if !b.Header().TxsRoot().EqualsTo(b.data.Transactions.Root()) {
		return errors.Errorf(errors.ErrInvalidBlock, "transactions root is not matched")
	}
	if b.PrevCertificate() != nil {
		if err := b.PrevCertificate().SanityCheck(); err != nil {
			return err
		}
		if !b.Header().PrevCertificateHash().EqualsTo(b.PrevCertificate().Hash()) {
			return errors.Errorf(errors.ErrInvalidBlock, "invalid previous certificate hash")
		}
	} else {
		// Genesis block checks
		if !b.Header().PrevCertificateHash().IsUndef() {
			return errors.Errorf(errors.ErrInvalidBlock, "invalid genesis certificate hash")
		}
	}

	for _, trx := range b.Transactions() {
		if err := trx.SanityCheck(); err != nil {
			return errors.Errorf(errors.ErrInvalidBlock, err.Error())
		}
	}

	return nil
}

func (b *Block) Hash() hash.Hash {
	return b.Header().Hash()
}

func (b *Block) Stamp() hash.Stamp {
	return b.Hash().Stamp()
}

func (b *Block) Fingerprint() string {
	return fmt.Sprintf("{⌘ %v 👤 %v 💻 %v 📨 %d}",
		b.Hash().Fingerprint(),
		b.data.Header.ProposerAddress().Fingerprint(),
		b.data.Header.StateRoot().Fingerprint(),
		b.data.Transactions.Len(),
	)
}
func (b *Block) MarshalCBOR() ([]byte, error) {
	return b.Encode()
}

func (b *Block) UnmarshalCBOR(bs []byte) error {
	return b.Decode(bs)
}

func (b *Block) Encode() ([]byte, error) {
	return cbor.Marshal(b.data)
}

func (b *Block) Decode(bs []byte) error {
	return cbor.Unmarshal(bs, &b.data)
}

func (b *Block) MarshalJSON() ([]byte, error) {
	return json.Marshal(b.data)
}

// ---------
// For tests
func GenerateTestBlock(proposer *crypto.Address, prevBlockHash *hash.Hash) *Block {
	if proposer == nil {
		addr := crypto.GenerateTestAddress()
		proposer = &addr
	}
	txs := NewTxs()
	tx1, _ := tx.GenerateTestSendTx()
	tx2, _ := tx.GenerateTestSendTx()
	tx3, _ := tx.GenerateTestSendTx()
	tx4, _ := tx.GenerateTestSendTx()

	txs.Append(tx1)
	txs.Append(tx2)
	txs.Append(tx3)
	txs.Append(tx4)

	if prevBlockHash == nil {
		h := hash.GenerateTestHash()
		prevBlockHash = &h
	}
	cert := GenerateTestCertificate(*prevBlockHash)
	if prevBlockHash.IsUndef() {
		cert = nil
	}
	sortitionSeed := sortition.GenerateRandomSeed()
	header := NewHeader(1, util.Now(),
		txs.Root(),
		hash.GenerateTestHash(),
		hash.GenerateTestHash(),
		cert.Hash(),
		sortitionSeed,
		*proposer)

	return NewBlock(header, cert, txs)
}
