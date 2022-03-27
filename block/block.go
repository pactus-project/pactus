package block

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/crypto/hash"
	"github.com/zarbchain/zarb-go/encoding"
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/sortition"
	"github.com/zarbchain/zarb-go/tx"
	"github.com/zarbchain/zarb-go/util"
)

type Block struct {
	memorizedHash *hash.Hash
	memorizedData []byte
	data          blockData
}

type blockData struct {
	Header   Header
	PrevCert *Certificate
	Txs      Txs
}

func NewBlock(header Header, prevCert *Certificate, txs Txs) *Block {
	return &Block{
		data: blockData{
			Header:   header,
			PrevCert: prevCert,
			Txs:      txs,
		},
	}
}

func MakeBlock(version uint8, timestamp time.Time, txs Txs,
	prevBlockHash, stateRoot hash.Hash,
	prevCert *Certificate, sortitionSeed sortition.VerifiableSeed, proposer crypto.Address) *Block {
	header := NewHeader(version, timestamp,
		stateRoot, prevBlockHash, sortitionSeed, proposer)

	b := NewBlock(header, prevCert, txs)
	if err := b.SanityCheck(); err != nil {
		panic(err)
	}
	return b
}

func (b *Block) Header() *Header               { return &b.data.Header }
func (b *Block) PrevCertificate() *Certificate { return b.data.PrevCert }
func (b *Block) Transactions() Txs             { return b.data.Txs }

func (b *Block) SanityCheck() error {
	if err := b.Header().SanityCheck(); err != nil {
		return err
	}
	if b.Transactions().Len() == 0 {
		return errors.Errorf(errors.ErrInvalidBlock, "block at least should have one transaction")
	}
	if b.Transactions().Len() > 1000 {
		return errors.Errorf(errors.ErrInvalidBlock, "block is full")
	}
	if b.PrevCertificate() != nil {
		if err := b.PrevCertificate().SanityCheck(); err != nil {
			return err
		}
		if err := b.Header().PrevBlockHash().SanityCheck(); err != nil {
			return err
		}
	} else {
		// Genesis block checks
		if !b.Header().PrevBlockHash().IsUndef() {
			return errors.Errorf(errors.ErrInvalidBlock, "invalid previous block hash")
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
	if b.memorizedHash != nil {
		return *b.memorizedHash
	}

	w := &bytes.Buffer{}
	b.data.Header.Encode(w)
	w.Write(b.data.PrevCert.Hash().RawBytes())
	w.Write(b.data.Txs.Root().RawBytes())
	w.Write(util.Int32ToSlice(int32(b.data.Txs.Len())))

	h := hash.CalcHash(w.Bytes())
	b.memorizedHash = &h
	return h
}

func (b *Block) Stamp() hash.Stamp {
	return b.Hash().Stamp()
}

func (b *Block) Fingerprint() string {
	return fmt.Sprintf("{⌘ %v 👤 %v 💻 %v 📨 %d}",
		b.Hash().Fingerprint(),
		b.data.Header.ProposerAddress().Fingerprint(),
		b.data.Header.StateRoot().Fingerprint(),
		b.data.Txs.Len(),
	)
}

// Encode encodes the receiver to w.
func (b *Block) Encode(w io.Writer) error {
	if err := b.data.Header.Encode(w); err != nil {
		return err
	}
	if err := b.data.PrevCert.Encode(w); err != nil {
		return err
	}
	encoding.WriteVarInt(w, uint64(b.data.Txs.Len()))
	for _, tx := range b.Transactions() {
		if err := tx.Encode(w); err != nil {
			return err
		}
	}
	return nil
}

func (b *Block) Decode(r io.Reader) error {
	if err := b.data.Header.Decode(r); err != nil {
		return err
	}
	if err := b.data.PrevCert.Decode(r); err != nil {
		return err
	}
	len, err := encoding.ReadVarInt(r)
	if err != nil {
		return err
	}
	b.data.Txs = make([]*tx.Tx, len)
	for _, tx := range b.Transactions() {
		if err := tx.Decode(r); err != nil {
			return err
		}
	}
	return nil
}

// SerializeSize returns the number of bytes it would take to serialize the block
func (b *Block) SerializeSize() int {
	n := b.Header().SerializeSize() +
		b.PrevCertificate().SerializeSize()

	for _, tx := range b.Transactions() {
		n += tx.SerializeSize()
	}
	return n
}

// Bytes returns the serialized bytes for the Block. It caches the
// result so subsequent calls are more efficient.
func (b *Block) Bytes() ([]byte, error) {
	// Return the cached serialized bytes if it has already been generated.
	if len(b.memorizedData) != 0 {
		return b.memorizedData, nil
	}

	w := bytes.NewBuffer(make([]byte, 0, b.SerializeSize()))
	err := b.Encode(w)
	if err != nil {
		return nil, err
	}

	// Cache the serialized bytes and return them.
	b.memorizedData = w.Bytes()
	return b.memorizedData, nil
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
		hash.GenerateTestHash(),
		*prevBlockHash,
		sortitionSeed,
		*proposer)

	return NewBlock(header, cert, txs)
}
