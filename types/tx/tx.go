package tx

import (
	"bytes"
	"fmt"
	"io"

	"github.com/fxamacker/cbor/v2"
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/sortition"
	"github.com/pactus-project/pactus/types/tx/payload"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/encoding"
	"github.com/pactus-project/pactus/util/errors"
)

const versionLatest = 0x01
const flagLockTime = 0x80

const maxMemoLength = 64

type ID = hash.Hash

type Tx struct {
	memorizedID   *ID
	sanityChecked bool

	data txData
}

type txData struct {
	// Version format
	//  7             0
	// +-+-+-+-+-+-+-+-+
	// |L|R|R|R|VERSION|
	// +-+-+-+-+-+-+-+-+
	// L: Lock Time transacion
	// R: Reserved bit
	//
	Version   uint8
	Stamp     hash.Stamp
	LockTime  uint32
	Sequence  int32
	Fee       int64
	Payload   payload.Payload
	Memo      string
	PublicKey crypto.PublicKey
	Signature crypto.Signature
}

func NewTx(stamp hash.Stamp, seq int32, pld payload.Payload, fee int64,
	memo string) *Tx {
	trx := &Tx{
		data: txData{
			Stamp:    stamp,
			Sequence: seq,
			Version:  versionLatest,
			Payload:  pld,
			Fee:      fee,
			Memo:     memo,
		},
	}

	return trx
}

func NewLockTimeTx(lockTime uint32, seq int32, pld payload.Payload, fee int64,
	memo string) *Tx {
	trx := &Tx{
		data: txData{
			LockTime: lockTime,
			Sequence: seq,
			Version:  versionLatest | flagLockTime,
			Payload:  pld,
			Fee:      fee,
			Memo:     memo,
		},
	}

	return trx
}

// FromBytes constructs a new transaction from byte array.
func FromBytes(bs []byte) (*Tx, error) {
	tx := new(Tx)
	r := bytes.NewReader(bs)
	if err := tx.Decode(r); err != nil {
		return nil, err
	}
	return tx, nil
}

func (tx *Tx) Version() uint8 {
	return tx.data.Version & 0x0f
}

func (tx *Tx) Stamp() hash.Stamp {
	return tx.data.Stamp
}

func (tx *Tx) Sequence() int32 {
	return tx.data.Sequence
}

func (tx *Tx) Payload() payload.Payload {
	return tx.data.Payload
}

func (tx *Tx) Fee() int64 {
	return tx.data.Fee
}

func (tx *Tx) Memo() string {
	return tx.data.Memo
}

func (tx *Tx) PublicKey() crypto.PublicKey {
	return tx.data.PublicKey
}

func (tx *Tx) Signature() crypto.Signature {
	return tx.data.Signature
}

func (tx *Tx) LockTime() uint32 {
	return tx.data.LockTime
}

func (tx *Tx) IsStamped() bool {
	return tx.data.Version&flagLockTime == 0x00
}

func (tx *Tx) IsLockTime() bool {
	return tx.data.Version&flagLockTime == flagLockTime
}

func (tx *Tx) SetSignature(sig crypto.Signature) {
	tx.sanityChecked = false
	tx.data.Signature = sig
}

func (tx *Tx) SetPublicKey(pub crypto.PublicKey) {
	tx.sanityChecked = false
	tx.data.PublicKey = pub
}

func (tx *Tx) SanityCheck() error {
	if tx.sanityChecked {
		return nil
	}
	if tx.Version() != versionLatest {
		return errors.Errorf(errors.ErrInvalidTx, "invalid version")
	}
	if tx.Sequence() < 0 {
		return errors.Error(errors.ErrInvalidSequence)
	}
	if tx.Payload().Value() < 0 || tx.Payload().Value() > 21*1e14 {
		return errors.Error(errors.ErrInvalidAmount)
	}
	if err := tx.checkFee(); err != nil {
		return err
	}
	if err := tx.Payload().SanityCheck(); err != nil {
		return err
	}
	if len(tx.Memo()) > maxMemoLength {
		return errors.Error(errors.ErrInvalidMemo)
	}
	if err := tx.checkSignature(); err != nil {
		return err
	}

	tx.sanityChecked = true

	return nil
}

func (tx *Tx) checkFee() error {
	if tx.IsFreeTx() {
		if tx.Fee() != 0 {
			return errors.Errorf(errors.ErrInvalidFee, "fee should set to zero")
		}
	} else {
		if tx.Fee() <= 0 {
			return errors.Errorf(errors.ErrInvalidFee, "fee should not be negative")
		}
	}

	return nil
}

func (tx *Tx) checkSignature() error {
	if tx.IsSubsidyTx() {
		if tx.PublicKey() != nil {
			return errors.Errorf(errors.ErrInvalidPublicKey, "subsidy transaction should not have public key")
		}
		if tx.Signature() != nil {
			return errors.Errorf(errors.ErrInvalidSignature, "subsidy transaction should not have signature")
		}
	} else {
		if tx.PublicKey() == nil {
			return errors.Errorf(errors.ErrInvalidPublicKey, "no public key")
		}
		if tx.Signature() == nil {
			return errors.Errorf(errors.ErrInvalidSignature, "no signature")
		}
		if err := tx.PublicKey().VerifyAddress(tx.Payload().Signer()); err != nil {
			return err
		}
		bs := tx.SignBytes()
		if err := tx.PublicKey().Verify(bs, tx.Signature()); err != nil {
			return errors.Error(errors.ErrInvalidSignature)
		}
	}
	return nil
}

// Bytes returns the serialized bytes for the Transaction.
func (tx *Tx) Bytes() ([]byte, error) {
	w := bytes.NewBuffer(make([]byte, 0, tx.SerializeSize()))
	err := tx.Encode(w)
	if err != nil {
		return nil, err
	}

	return w.Bytes(), nil
}

func (tx *Tx) MarshalCBOR() ([]byte, error) {
	buf := bytes.NewBuffer(make([]byte, 0, tx.SerializeSize()))
	if err := tx.Encode(buf); err != nil {
		return nil, err
	}
	return cbor.Marshal(buf.Bytes())
}

func (tx *Tx) UnmarshalCBOR(bs []byte) error {
	data := make([]byte, 0, tx.SerializeSize())
	err := cbor.Unmarshal(bs, &data)
	if err != nil {
		return err
	}
	buf := bytes.NewBuffer(data)
	return tx.Decode(buf)
}

// SerializeSize returns the number of bytes it would take to serialize the transaction.
func (tx *Tx) SerializeSize() int {
	n := 150 +
		encoding.VarIntSerializeSize(uint64(tx.Sequence())) +
		encoding.VarIntSerializeSize(uint64(tx.Fee())) +
		encoding.VarStringSerializeSize(tx.Memo())
	if tx.Payload() != nil {
		n += tx.Payload().SerializeSize()
	}

	return n
}

func (tx *Tx) Encode(w io.Writer) error {
	err := tx.EncodeWithNoSignatory(w)
	if err != nil {
		return err
	}
	if tx.data.Signature != nil {
		err = tx.data.Signature.Encode(w)
		if err != nil {
			return err
		}
	}
	if tx.data.PublicKey != nil {
		err = tx.data.PublicKey.Encode(w)
		if err != nil {
			return err
		}
	}

	return nil
}

func (tx *Tx) EncodeWithNoSignatory(w io.Writer) error {
	err := encoding.WriteElements(w, tx.data.Version, tx.data.Stamp)
	if err != nil {
		return err
	}
	err = encoding.WriteVarInt(w, uint64(tx.data.Sequence))
	if err != nil {
		return err
	}
	err = encoding.WriteVarInt(w, uint64(tx.data.Fee))
	if err != nil {
		return err
	}
	err = encoding.WriteElement(w, uint8(tx.data.Payload.Type()))
	if err != nil {
		return err
	}
	err = tx.data.Payload.Encode(w)
	if err != nil {
		return err
	}
	err = encoding.WriteVarString(w, tx.data.Memo)
	if err != nil {
		return err
	}
	return nil
}

func (tx *Tx) DecodeWithNoSignatory(r io.Reader) error {
	err := encoding.ReadElements(r, &tx.data.Version, &tx.data.Stamp)
	if err != nil {
		return err
	}

	seq, err := encoding.ReadVarInt(r)
	if err != nil {
		return err
	}
	tx.data.Sequence = int32(seq)

	fee, err := encoding.ReadVarInt(r)
	if err != nil {
		return err
	}
	tx.data.Fee = int64(fee)

	payloadType := uint8(0)
	err = encoding.ReadElement(r, &payloadType)
	if err != nil {
		return err
	}

	switch payload.Type(payloadType) {
	case payload.PayloadTypeTransfer:
		tx.data.Payload = &payload.TransferPayload{}
	case payload.PayloadTypeBond:
		tx.data.Payload = &payload.BondPayload{}
	case payload.PayloadTypeUnbond:
		tx.data.Payload = &payload.UnbondPayload{}
	case payload.PayloadTypeWithdraw:
		tx.data.Payload = &payload.WithdrawPayload{}
	case payload.PayloadTypeSortition:
		tx.data.Payload = &payload.SortitionPayload{}

	default:
		return errors.Errorf(errors.ErrInvalidTx, "invalid payload")
	}

	err = tx.data.Payload.Decode(r)
	if err != nil {
		return err
	}
	tx.data.Memo, err = encoding.ReadVarString(r)
	if err != nil {
		return err
	}
	return nil
}
func (tx *Tx) Decode(r io.Reader) error {
	err := tx.DecodeWithNoSignatory(r)
	if err != nil {
		return err
	}

	if !tx.IsSubsidyTx() {
		sig := new(bls.Signature)
		err = sig.Decode(r)
		if err != nil {
			return err
		}
		tx.data.Signature = sig

		pub := new(bls.PublicKey)
		err = pub.Decode(r)
		if err != nil {
			return err
		}
		tx.data.PublicKey = pub
	}

	return nil
}

func (tx *Tx) Fingerprint() string {
	return fmt.Sprintf("{⌘ %v 🏵 %v %v}",
		tx.ID().Fingerprint(),
		tx.data.Stamp.String(),
		tx.data.Payload.Fingerprint())
}

func (tx *Tx) SignBytes() []byte {
	buf := bytes.Buffer{}
	err := tx.EncodeWithNoSignatory(&buf)
	if err != nil {
		return nil
	}
	return buf.Bytes()
}

func (tx *Tx) ID() ID {
	if tx.memorizedID != nil {
		return *tx.memorizedID
	}
	id := hash.CalcHash(tx.SignBytes())
	tx.memorizedID = &id
	return id
}

func (tx *Tx) IsSendTx() bool {
	return tx.Payload().Type() == payload.PayloadTypeTransfer &&
		!tx.data.Payload.Signer().EqualsTo(crypto.TreasuryAddress)
}

func (tx *Tx) IsBondTx() bool {
	return tx.Payload().Type() == payload.PayloadTypeBond
}

func (tx *Tx) IsSubsidyTx() bool {
	return tx.Payload().Type() == payload.PayloadTypeTransfer &&
		tx.data.Payload.Signer().EqualsTo(crypto.TreasuryAddress)
}

func (tx *Tx) IsSortitionTx() bool {
	return tx.Payload().Type() == payload.PayloadTypeSortition
}

func (tx *Tx) IsUnbondTx() bool {
	return tx.Payload().Type() == payload.PayloadTypeUnbond
}

func (tx *Tx) IsWithdrawTx() bool {
	return tx.Payload().Type() == payload.PayloadTypeWithdraw
}

// IsFreeTx will checks if transaction fee is 0.
func (tx *Tx) IsFreeTx() bool {
	return tx.IsSubsidyTx() || tx.IsSortitionTx() || tx.IsUnbondTx()
}

// GenerateTestSendTx generates a send transaction for testing.
func GenerateTestSendTx() (*Tx, crypto.Signer) {
	stamp := hash.GenerateTestStamp()
	s := bls.GenerateTestSigner()
	pub, _ := bls.GenerateTestKeyPair()
	tx := NewTransferTx(stamp, util.RandInt32(1000), s.Address(), pub.Address(),
		util.RandInt64(1000*1e10), util.RandInt64(1*1e10), "test send-tx")
	s.SignMsg(tx)
	return tx, s
}

// GenerateTestSendTx generates a bond transaction for testing.
func GenerateTestBondTx() (*Tx, crypto.Signer) {
	stamp := hash.GenerateTestStamp()
	s := bls.GenerateTestSigner()
	pub, _ := bls.GenerateTestKeyPair()
	tx := NewBondTx(stamp, util.RandInt32(1000), s.Address(), pub.Address(),
		pub, util.RandInt64(1000*1e10), util.RandInt64(1*1e10), "test bond-tx")
	s.SignMsg(tx)
	return tx, s
}

// GenerateTestSendTx generates a sortition transaction for testing.
func GenerateTestSortitionTx() (*Tx, crypto.Signer) {
	stamp := hash.GenerateTestStamp()
	s := bls.GenerateTestSigner()
	proof := sortition.GenerateRandomProof()
	tx := NewSortitionTx(stamp, util.RandInt32(1000), s.Address(), proof)
	s.SignMsg(tx)
	return tx, s
}

// GenerateTestSendTx generates an unbond transaction for testing.
func GenerateTestUnbondTx() (*Tx, crypto.Signer) {
	stamp := hash.GenerateTestStamp()
	s := bls.GenerateTestSigner()
	tx := NewUnbondTx(stamp, util.RandInt32(1000), s.Address(), "test unbond-tx")
	s.SignMsg(tx)
	return tx, s
}

// GenerateTestSendTx generates a withdraw transaction for testing.
func GenerateTestWithdrawTx() (*Tx, crypto.Signer) {
	stamp := hash.GenerateTestStamp()
	s := bls.GenerateTestSigner()
	tx := NewWithdrawTx(stamp, util.RandInt32(1000), s.Address(), crypto.GenerateTestAddress(),
		util.RandInt64(1000*1e10), util.RandInt64(1*1e10), "test withdraw-tx")
	s.SignMsg(tx)
	return tx, s
}
