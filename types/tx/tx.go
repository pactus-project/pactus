package tx

import (
	"bytes"
	"fmt"
	"io"

	"github.com/fxamacker/cbor/v2"
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/tx/payload"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/encoding"
)

const (
	versionLatest        = 0x01
	flagStripedPublicKey = 0x01
	flagNotSigned        = 0x02
	maxMemoLength        = 64
)

type ID = hash.Hash

type Tx struct {
	memorizedID  *ID
	basicChecked bool

	data txData
}

type txData struct {
	Flags     uint8
	Version   uint8
	LockTime  uint32
	Fee       amount.Amount
	Memo      string
	Payload   payload.Payload
	Signature crypto.Signature
	PublicKey crypto.PublicKey
}

func newTx(lockTime uint32, pld payload.Payload, fee amount.Amount,
	memo string,
) *Tx {
	trx := &Tx{
		data: txData{
			Flags:    flagNotSigned,
			LockTime: lockTime,
			Version:  versionLatest,
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

func (tx *Tx) LockTime() uint32 {
	return tx.data.LockTime
}

func (tx *Tx) Payload() payload.Payload {
	return tx.data.Payload
}

func (tx *Tx) Fee() amount.Amount {
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

func (tx *Tx) SetSignature(sig crypto.Signature) {
	tx.basicChecked = false
	tx.data.Signature = sig

	if sig == nil {
		tx.data.Flags = util.SetFlag(tx.data.Flags, flagNotSigned)
	} else {
		tx.data.Flags = util.UnsetFlag(tx.data.Flags, flagNotSigned)
	}
}

func (tx *Tx) SetPublicKey(pub crypto.PublicKey) {
	tx.basicChecked = false
	tx.data.PublicKey = pub
	if pub == nil {
		if !tx.IsSubsidyTx() {
			tx.data.Flags = util.SetFlag(tx.data.Flags, flagStripedPublicKey)
		}
	} else {
		tx.data.Flags = util.UnsetFlag(tx.data.Flags, flagStripedPublicKey)
	}
}

func (tx *Tx) BasicCheck() error {
	if tx.basicChecked {
		return nil
	}
	if tx.Version() != versionLatest {
		return BasicCheckError{
			Reason: fmt.Sprintf("invalid version: %d", tx.Version()),
		}
	}
	if tx.LockTime() == 0 {
		return BasicCheckError{
			Reason: "lock time is not defined",
		}
	}
	if tx.Payload().Value() < 0 || tx.Payload().Value() > amount.MaxNanoPAC {
		return BasicCheckError{
			Reason: fmt.Sprintf("invalid amount: %s", tx.Payload().Value()),
		}
	}
	if tx.Fee() < 0 || tx.Fee() > amount.MaxNanoPAC {
		return BasicCheckError{
			Reason: fmt.Sprintf("invalid fee: %s", tx.Fee()),
		}
	}
	if len(tx.Memo()) > maxMemoLength {
		return BasicCheckError{
			Reason: fmt.Sprintf("memo length exceeded: %d", len(tx.Memo())),
		}
	}
	if err := tx.Payload().BasicCheck(); err != nil {
		return BasicCheckError{
			Reason: fmt.Sprintf("invalid payload: %s", err.Error()),
		}
	}
	if err := tx.checkSignature(); err != nil {
		return err
	}

	tx.basicChecked = true

	return nil
}

func (tx *Tx) checkSignature() error {
	if tx.IsSubsidyTx() {
		// Ensure no signatory is set for subsidy transactions.
		if tx.PublicKey() != nil || tx.Signature() != nil {
			return BasicCheckError{
				Reason: "subsidy transaction with signatory",
			}
		}

		return nil
	}

	// Non-subsidy transactions should have a valid signatory.
	if tx.PublicKey() == nil {
		return BasicCheckError{
			Reason: "no public key",
		}
	}

	if tx.Signature() == nil {
		return BasicCheckError{
			Reason: "no signature",
		}
	}

	if err := tx.PublicKey().VerifyAddress(tx.Payload().Signer()); err != nil {
		return BasicCheckError{
			Reason: err.Error(),
		}
	}

	bs := tx.SignBytes()
	if err := tx.PublicKey().Verify(bs, tx.Signature()); err != nil {
		return BasicCheckError{
			Reason: "invalid signature",
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
	n := 3 + // one byte version, flag, payload type
		4 + // for tx.LockTime
		encoding.VarIntSerializeSize(uint64(tx.Fee())) +
		encoding.VarStringSerializeSize(tx.Memo())
	if tx.Payload() != nil {
		n += tx.Payload().SerializeSize()
	}
	if tx.data.Signature != nil {
		n += bls.SignatureSize
	}
	if tx.data.PublicKey != nil {
		n += bls.PublicKeySize
	}

	return n
}

func (tx *Tx) Encode(w io.Writer) error {
	err := tx.encodeWithNoSignatory(w)
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

func (tx *Tx) encodeWithNoSignatory(w io.Writer) error {
	err := encoding.WriteElements(w, tx.data.Flags, tx.data.Version, tx.data.LockTime)
	if err != nil {
		return err
	}
	err = encoding.WriteVarInt(w, uint64(tx.data.Fee))
	if err != nil {
		return err
	}
	err = encoding.WriteVarString(w, tx.data.Memo)
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

	return nil
}

func (tx *Tx) Decode(r io.Reader) error {
	err := encoding.ReadElements(r, &tx.data.Flags, &tx.data.Version, &tx.data.LockTime)
	if err != nil {
		return err
	}

	fee, err := encoding.ReadVarInt(r)
	if err != nil {
		return err
	}
	tx.data.Fee = amount.Amount(fee)

	tx.data.Memo, err = encoding.ReadVarString(r)
	if err != nil {
		return err
	}

	payloadType := uint8(0)
	err = encoding.ReadElement(r, &payloadType)
	if err != nil {
		return err
	}

	switch t := payload.Type(payloadType); t {
	case payload.TypeTransfer:
		tx.data.Payload = new(payload.TransferPayload)
	case payload.TypeBond:
		tx.data.Payload = new(payload.BondPayload)
	case payload.TypeUnbond:
		tx.data.Payload = new(payload.UnbondPayload)
	case payload.TypeWithdraw:
		tx.data.Payload = new(payload.WithdrawPayload)
	case payload.TypeSortition:
		tx.data.Payload = new(payload.SortitionPayload)

	default:
		return InvalidPayloadTypeError{
			PayloadType: t,
		}
	}

	err = tx.data.Payload.Decode(r)
	if err != nil {
		return err
	}

	if !util.IsFlagSet(tx.data.Flags, flagNotSigned) {
		sig := new(bls.Signature)
		err = sig.Decode(r)
		if err != nil {
			return err
		}
		tx.data.Signature = sig

		if !tx.IsPublicKeyStriped() {
			pub := new(bls.PublicKey)
			err = pub.Decode(r)
			if err != nil {
				return err
			}
			tx.data.PublicKey = pub
		}
	}

	return nil
}

func (tx *Tx) String() string {
	return fmt.Sprintf("{⌘ %v - %v 🏵 %v}",
		tx.ID().ShortString(),
		tx.LockTime(),
		tx.data.Payload.String())
}

func (tx *Tx) SignBytes() []byte {
	buf := bytes.Buffer{}
	err := tx.encodeWithNoSignatory(&buf)
	if err != nil {
		return nil
	}

	return buf.Bytes()[1:] // Exclude flags
}

func (tx *Tx) ID() ID {
	if tx.memorizedID != nil {
		return *tx.memorizedID
	}
	id := hash.CalcHash(tx.SignBytes())
	tx.memorizedID = &id

	return id
}

func (tx *Tx) IsTransferTx() bool {
	return tx.Payload().Type() == payload.TypeTransfer &&
		tx.Payload().Signer() != crypto.TreasuryAddress
}

func (tx *Tx) IsBondTx() bool {
	return tx.Payload().Type() == payload.TypeBond
}

func (tx *Tx) IsSubsidyTx() bool {
	return tx.Payload().Type() == payload.TypeTransfer &&
		tx.Payload().Signer() == crypto.TreasuryAddress
}

func (tx *Tx) IsSortitionTx() bool {
	return tx.Payload().Type() == payload.TypeSortition
}

func (tx *Tx) IsUnbondTx() bool {
	return tx.Payload().Type() == payload.TypeUnbond
}

func (tx *Tx) IsWithdrawTx() bool {
	return tx.Payload().Type() == payload.TypeWithdraw
}

// StripPublicKey removes the public key from the transaction.
// It is an alias function for `SetPublicKey(nil)`.
func (tx *Tx) StripPublicKey() {
	tx.SetPublicKey(nil)
}

// IsPublicKeyStriped returns true if the public key stripped from the transaction.
func (tx *Tx) IsPublicKeyStriped() bool {
	return util.IsFlagSet(tx.data.Flags, flagStripedPublicKey)
}

// IsSigned returns true if the transaction has been signed and includes the signature.
func (tx *Tx) IsSigned() bool {
	return !util.IsFlagSet(tx.data.Flags, flagNotSigned)
}
