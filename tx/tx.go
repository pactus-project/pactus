package tx

import (
	"encoding/json"

	"github.com/fxamacker/cbor/v2"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/errors"
)

type Type int

const (
	TypeSend = Type(0)
	TypeCall = Type(1)
)

type Tx struct {
	data txData

	memorizedHash *crypto.Hash
}

type txData struct {
	Stamp     crypto.Hash       `cbor:"1,keyasint"`
	Version   int               `cbor:"2,keyasint"`
	Sender    crypto.Address    `cbor:"3,keyasint"`
	Receiver  crypto.Address    `cbor:"4,keyasint"`
	Amount    int64             `cbor:"5,keyasint"`
	Fee       int64             `cbor:"6,keyasint"`
	Memo      string            `cbor:"7,keyasint"`
	Data      []byte            `cbor:"16,keyasint,omitempty"`
	GasWants  int               `cbor:"17,keyasint,omitempty"`
	PublicKey *crypto.PublicKey `cbor:"20,keyasint,omitempty"`
	Signature *crypto.Signature `cbor:"21,keyasint,omitempty"`
}

func (tx *Tx) Stamp() crypto.Hash       { return tx.data.Stamp }
func (tx *Tx) Sender() crypto.Address   { return tx.data.Sender }
func (tx *Tx) Receiver() crypto.Address { return tx.data.Receiver }
func (tx *Tx) Amount() int64            { return tx.data.Amount }
func (tx *Tx) Fee() int64               { return tx.data.Fee }
func (tx *Tx) GasWants() int            { return tx.data.GasWants }
func (tx *Tx) Data() []byte             { return tx.data.Data }
func (tx *Tx) Memo() string             { return tx.data.Memo }
func (tx *Tx) IsCallTx() bool           { return len(tx.data.Data) > 0 }
func (tx *Tx) IsMintbaseTx() bool       { return tx.data.Sender.EqualsTo(crypto.MintbaseAddress) }

func NewMintbaseTx(stamp crypto.Hash, receiver crypto.Address, amount int64, memo string) *Tx {
	return &Tx{
		data: txData{
			Stamp:    stamp,
			Version:  1,
			Sender:   crypto.MintbaseAddress,
			Receiver: receiver,
			Amount:   amount,
			Fee:      0,
		},
	}
}

func NewSendTx(stamp crypto.Hash,
	sender, receiver crypto.Address,
	amount, fee int64, memo string,
	publicKey crypto.PublicKey, signature crypto.Signature) *Tx {
	return &Tx{
		data: txData{
			Stamp:     stamp,
			Version:   1,
			Sender:    sender,
			Receiver:  receiver,
			Amount:    amount,
			Fee:       fee,
			PublicKey: &publicKey,
			Signature: &signature,
		},
	}
}

func (tx *Tx) SanityCheck() error {
	if tx.data.Version != 1 {
		return errors.Errorf(errors.ErrInvalidTx, "Invalid version")
	}
	if tx.data.Amount < 0 {
		return errors.Errorf(errors.ErrInvalidTx, "Invalid amount")
	}
	if len(tx.data.Memo) > 256 {
		return errors.Errorf(errors.ErrInvalidTx, "Invalid memo")
	}
	if err := tx.data.Stamp.SanityCheck(); err != nil {
		return errors.Errorf(errors.ErrInvalidTx, "Invalid stamp")
	}
	if err := tx.data.Sender.SanityCheck(); err != nil {
		return errors.Errorf(errors.ErrInvalidTx, "Invalid sender address")
	}
	if err := tx.data.Receiver.SanityCheck(); err != nil {
		return errors.Errorf(errors.ErrInvalidTx, "Invalid receiver address")
	}

	if tx.IsCallTx() {
		if tx.data.GasWants <= 0 {
			return errors.Errorf(errors.ErrInvalidTx, "Invalid gas wants")
		}
	} else {
		if tx.data.GasWants != 0 {
			return errors.Errorf(errors.ErrInvalidTx, "Send tx has no gas")
		}
	}

	if tx.IsMintbaseTx() {
		if tx.data.Data != nil {
			return errors.Errorf(errors.ErrInvalidTx, "Mintbase transaction should not have data")
		}
		if tx.data.PublicKey != nil {
			return errors.Errorf(errors.ErrInvalidTx, "Mintbase transaction should not have public key")
		}
		if tx.data.Signature != nil {
			return errors.Errorf(errors.ErrInvalidTx, "Mintbase transaction should not have signature")
		}
	} else {
		if tx.data.PublicKey == nil {
			return errors.Errorf(errors.ErrInvalidTx, "Transaction should not have public key")
		}
		if tx.data.Signature == nil {
			return errors.Errorf(errors.ErrInvalidTx, "Transaction should not have signature")
		}
		if tx.data.Fee <= 0 {
			return errors.Errorf(errors.ErrInvalidTx, "Invalid fee")
		}
		if err := tx.data.PublicKey.SanityCheck(); err != nil {
			return errors.Errorf(errors.ErrInvalidTx, "Invalid pubic key")
		}
		if err := tx.data.Signature.SanityCheck(); err != nil {
			return errors.Errorf(errors.ErrInvalidTx, "Invalid signature")
		}
		if tx.data.Sender.Verify(*tx.data.PublicKey) {
			return errors.Errorf(errors.ErrInvalidTx, "Invalid cryptographic public key")
		}
		bs, _ := tx.SignBytes()
		if tx.data.PublicKey.Verify(bs, *tx.data.Signature) {
			return errors.Errorf(errors.ErrInvalidTx, "Invalid cryptographic signature")
		}
	}

	return nil
}

func (tx *Tx) Hash() crypto.Hash {
	if tx.memorizedHash == nil {
		bz, _ := tx.Encode()
		hash := crypto.HashH(bz)
		tx.memorizedHash = &hash
	}

	return *tx.memorizedHash
}

func (tx *Tx) String() string {
	bz, _ := json.Marshal(tx.data)
	return string(bz)
}

func (tx *Tx) GenerateReceipt(status int) *Receipt {
	return &Receipt{
		data: receiptData{
			TxHash: tx.Hash(),
			Status: status,
		},
	}
}

func (tx Tx) SignBytes() ([]byte, error) {
	tx2 := tx
	tx2.data.PublicKey = nil
	tx2.data.Signature = nil

	return cbor.Marshal(tx.data)
}

func (tx *Tx) Encode() ([]byte, error) {
	return cbor.Marshal(tx.data)
}

func (tx *Tx) Decode(bs []byte) error {
	return cbor.Unmarshal(bs, &tx.data)
}

func (tx *Tx) MarshalJSON() ([]byte, error) {
	return json.Marshal(tx.data)
}

func (tx *Tx) UnmarshalJSON(bs []byte) error {
	return json.Unmarshal(bs, &tx.data)
}
