package block

import (
	"bytes"
	"fmt"
	"io"
	"time"

	"github.com/fxamacker/cbor/v2"
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/sortition"
	"github.com/pactus-project/pactus/types/certificate"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/encoding"
)

type Block struct {
	memorizedHash *hash.Hash
	memorizedData []byte
	data          blockData
}

type blockData struct {
	Header   *Header
	PrevCert *certificate.BlockCertificate
	Txs      Txs
}

func NewBlock(header *Header, prevCert *certificate.BlockCertificate, txs Txs) *Block {
	return &Block{
		data: blockData{
			Header:   header,
			PrevCert: prevCert,
			Txs:      txs,
		},
	}
}

// FromBytes constructs a new block from byte array.
func FromBytes(data []byte) (*Block, error) {
	b := new(Block)
	r := bytes.NewReader(data)
	if err := b.Decode(r); err != nil {
		return nil, err
	}

	return b, nil
}

func MakeBlock(version uint8, timestamp time.Time, txs Txs,
	prevBlockHash, stateRoot hash.Hash,
	prevCert *certificate.BlockCertificate, sortitionSeed sortition.VerifiableSeed, proposer crypto.Address,
) *Block {
	header := NewHeader(version, timestamp,
		stateRoot, prevBlockHash, sortitionSeed, proposer)

	return NewBlock(header, prevCert, txs)
}

func (b *Block) Header() *Header {
	return b.data.Header
}

func (b *Block) PrevCertificate() *certificate.BlockCertificate {
	return b.data.PrevCert
}

func (b *Block) Transactions() Txs {
	return b.data.Txs
}

func (b *Block) BasicCheck() error {
	if err := b.Header().BasicCheck(); err != nil {
		return err
	}
	if b.Transactions().Len() == 0 {
		// block at least should have one transaction
		return BasicCheckError{
			Reason: "no subsidy transaction",
		}
	}
	if b.Transactions().Len() > 1000 {
		return BasicCheckError{
			Reason: "block is full",
		}
	}
	if b.PrevCertificate() != nil {
		if err := b.PrevCertificate().BasicCheck(); err != nil {
			return BasicCheckError{
				Reason: fmt.Sprintf("invalid certificate: %s", err.Error()),
			}
		}
	} else {
		// Genesis block checks
		if !b.Header().PrevBlockHash().IsUndef() {
			return BasicCheckError{
				Reason: "invalid genesis block hash",
			}
		}
	}

	for _, trx := range b.Transactions() {
		if err := trx.BasicCheck(); err != nil {
			return BasicCheckError{
				Reason: fmt.Sprintf("invalid transaction: %s", err.Error()),
			}
		}
	}

	return nil
}

func (b *Block) Hash() hash.Hash {
	if b.memorizedHash != nil {
		return *b.memorizedHash
	}

	w := &bytes.Buffer{}
	if err := b.data.Header.Encode(w); err != nil {
		return hash.UndefHash
	}
	// Genesis block has no certificate
	if b.data.PrevCert != nil {
		w.Write(b.data.PrevCert.Hash().Bytes())
	}
	w.Write(b.data.Txs.Root().Bytes())
	w.Write(util.Int32ToSlice(int32(b.data.Txs.Len())))

	h := hash.CalcHash(w.Bytes())
	b.memorizedHash = &h

	return h
}

func (b *Block) Height() uint32 {
	if b.data.PrevCert == nil {
		return 1
	}

	return b.PrevCertificate().Height() + 1
}

func (b *Block) String() string {
	return fmt.Sprintf("{⌘ %v 👤 %v 💻 %v 📨 %d}",
		b.Hash().ShortString(),
		b.data.Header.ProposerAddress().ShortString(),
		b.data.Header.StateRoot().ShortString(),
		b.data.Txs.Len(),
	)
}

func (b *Block) MarshalCBOR() ([]byte, error) {
	buf := bytes.NewBuffer(make([]byte, 0, b.SerializeSize()))
	if err := b.Encode(buf); err != nil {
		return nil, err
	}

	return cbor.Marshal(buf.Bytes())
}

func (b *Block) UnmarshalCBOR(bs []byte) error {
	data := make([]byte, 0, b.SerializeSize())
	err := cbor.Unmarshal(bs, &data)
	if err != nil {
		return err
	}
	buf := bytes.NewBuffer(data)

	return b.Decode(buf)
}

func (b *Block) Encode(w io.Writer) error {
	if err := b.data.Header.Encode(w); err != nil {
		return err
	}
	if b.data.PrevCert != nil {
		if err := b.data.PrevCert.Encode(w); err != nil {
			return err
		}
	}
	if err := encoding.WriteVarInt(w, uint64(b.data.Txs.Len())); err != nil {
		return err
	}
	for _, trx := range b.Transactions() {
		if err := trx.Encode(w); err != nil {
			return err
		}
	}

	return nil
}

func (b *Block) Decode(r io.Reader) error {
	b.data.Header = new(Header)
	if err := b.data.Header.Decode(r); err != nil {
		return err
	}
	if !b.data.Header.PrevBlockHash().IsUndef() {
		b.data.PrevCert = new(certificate.BlockCertificate)
		if err := b.data.PrevCert.Decode(r); err != nil {
			return err
		}
	}
	length, err := encoding.ReadVarInt(r)
	if err != nil {
		return err
	}
	b.data.Txs = make([]*tx.Tx, length)
	for i := 0; i < int(length); i++ {
		trx := new(tx.Tx)
		if err := trx.Decode(r); err != nil {
			return err
		}
		b.data.Txs[i] = trx
	}

	return nil
}

// SerializeSize returns the number of bytes it would take to serialize the block.
func (b *Block) SerializeSize() int {
	n := b.Header().SerializeSize()

	if b.PrevCertificate() != nil {
		n += b.PrevCertificate().SerializeSize()
	}

	n += encoding.VarIntSerializeSize(uint64(b.Transactions().Len()))
	for _, trx := range b.Transactions() {
		n += trx.SerializeSize()
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
