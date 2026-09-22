// Package account provides functionality for managing account information.
package account

import (
	"bytes"
	"io"
	"unicode/utf8"

	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/types"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/util/encoding"
)

const (
	accountHeaderSize = 12 // number (int32) + balance (int64)

	anchorFlagSet     = 1
	minAnchorHashLen  = 32
	maxAnchorHashLen  = 64
	maxAnchorURILen   = 128
	anchorFixedSuffix = 1 + 1 + 1 + 1 + 8 + 16 // flag, hashLen, uriLen, type, deposit, timestamps
)

// The Account struct represents an account object.
type Account struct {
	data accountData
}

// accountData contains the data associated with an account.
type accountData struct {
	Number  int32
	Balance amount.Amount

	// Optional anchor (PIP-50). Nil means the account has no anchor.
	Anchor *AnchorData
}

// AnchorData holds the optional anchor stored on an account.
type AnchorData struct {
	RootHash        []byte
	ManifestURI     string
	AnchorType      uint8
	LockedDeposit   amount.Amount
	CreatedAtHeight types.Height
	CreatedAtTime   uint32
	UpdatedAtHeight types.Height
	UpdatedAtTime   uint32
}

// NewAccount constructs a new account from the given number.
func NewAccount(number int32) *Account {
	return &Account{
		data: accountData{
			Number: number,
		},
	}
}

// FromBytes constructs a new account from raw byte data.
// Bytes after the 12-byte header are an anchor suffix.
func FromBytes(data []byte) (*Account, error) {
	acc := new(Account)
	r := bytes.NewReader(data)
	err := encoding.ReadElements(r,
		&acc.data.Number,
		&acc.data.Balance)
	if err != nil {
		return nil, err
	}

	if r.Len() == 0 {
		return acc, nil
	}

	if err := acc.decodeAnchor(r); err != nil {
		return nil, err
	}

	return acc, nil
}

func (acc *Account) decodeAnchor(r *bytes.Reader) error {
	var flag uint8
	if err := encoding.ReadElement(r, &flag); err != nil {
		return err
	}
	if flag != anchorFlagSet {
		return ErrInvalidAnchorFlag
	}

	var hashLen uint8
	if err := encoding.ReadElement(r, &hashLen); err != nil {
		return err
	}
	if hashLen < minAnchorHashLen || hashLen > maxAnchorHashLen {
		return ErrAnchorHashLength
	}

	root := make([]byte, int(hashLen))
	if _, err := io.ReadFull(r, root); err != nil {
		return err
	}

	var uriLen uint8
	if err := encoding.ReadElement(r, &uriLen); err != nil {
		return err
	}
	if uriLen > maxAnchorURILen {
		return ErrAnchorURILength
	}

	uriBytes := make([]byte, int(uriLen))
	if _, err := io.ReadFull(r, uriBytes); err != nil {
		return err
	}
	if !utf8.Valid(uriBytes) {
		return ErrInvalidAnchorURI
	}

	anchor := &AnchorData{
		RootHash:    root,
		ManifestURI: string(uriBytes),
	}
	if err := encoding.ReadElements(r,
		&anchor.AnchorType,
		&anchor.LockedDeposit,
		&anchor.CreatedAtHeight,
		&anchor.CreatedAtTime,
		&anchor.UpdatedAtHeight,
		&anchor.UpdatedAtTime,
	); err != nil {
		return err
	}

	if err := validateAnchorDeposit(anchor.LockedDeposit); err != nil {
		return err
	}
	if r.Len() != 0 {
		return ErrAnchorTrailingBytes
	}

	acc.data.Anchor = anchor

	return nil
}

// Number returns the number of the account.
func (acc Account) Number() int32 {
	return acc.data.Number
}

// Balance returns the balance of the account, excluding the anchor deposit.
func (acc Account) Balance() amount.Amount {
	return acc.data.Balance
}

// HasAnchor returns true if the account has an anchor.
func (acc Account) HasAnchor() bool {
	return acc.data.Anchor != nil
}

// RootHash returns a copy of the anchor root hash.
func (acc Account) RootHash() []byte {
	if acc.data.Anchor == nil {
		return nil
	}

	return append([]byte(nil), acc.data.Anchor.RootHash...)
}

// ManifestURI returns the anchor manifest URI.
func (acc Account) ManifestURI() string {
	if acc.data.Anchor == nil {
		return ""
	}

	return acc.data.Anchor.ManifestURI
}

// AnchorType returns the anchor type.
func (acc Account) AnchorType() uint8 {
	if acc.data.Anchor == nil {
		return 0
	}

	return acc.data.Anchor.AnchorType
}

// LockedDeposit returns the locked anchor deposit.
func (acc Account) LockedDeposit() amount.Amount {
	if acc.data.Anchor == nil {
		return 0
	}

	return acc.data.Anchor.LockedDeposit
}

// CreatedAtHeight returns the height at which the anchor was created.
func (acc Account) CreatedAtHeight() types.Height {
	if acc.data.Anchor == nil {
		return 0
	}

	return acc.data.Anchor.CreatedAtHeight
}

// CreatedAtTime returns the time at which the anchor was created.
func (acc Account) CreatedAtTime() uint32 {
	if acc.data.Anchor == nil {
		return 0
	}

	return acc.data.Anchor.CreatedAtTime
}

// UpdatedAtHeight returns the height at which the anchor was last updated.
func (acc Account) UpdatedAtHeight() types.Height {
	if acc.data.Anchor == nil {
		return 0
	}

	return acc.data.Anchor.UpdatedAtHeight
}

// UpdatedAtTime returns the time at which the anchor was last updated.
func (acc Account) UpdatedAtTime() uint32 {
	if acc.data.Anchor == nil {
		return 0
	}

	return acc.data.Anchor.UpdatedAtTime
}

// SetAnchor sets the anchor fields (PIP-50).
func (acc *Account) SetAnchor(anchor AnchorData) error {
	if err := validateAnchor(&anchor); err != nil {
		return err
	}

	copied := anchor
	copied.RootHash = append([]byte(nil), anchor.RootHash...)
	acc.data.Anchor = &copied

	return nil
}

// ClearAnchor clears the anchor.
func (acc *Account) ClearAnchor() {
	acc.data.Anchor = nil
}

func validateAnchor(anchor *AnchorData) error {
	hashLen := len(anchor.RootHash)
	if hashLen < minAnchorHashLen || hashLen > maxAnchorHashLen {
		return ErrAnchorHashLength
	}
	if len(anchor.ManifestURI) > maxAnchorURILen {
		return ErrAnchorURILength
	}
	if !utf8.ValidString(anchor.ManifestURI) {
		return ErrInvalidAnchorURI
	}

	return validateAnchorDeposit(anchor.LockedDeposit)
}

func validateAnchorDeposit(deposit amount.Amount) error {
	if deposit <= 0 || deposit > amount.MaxNanoPAC {
		return ErrInvalidAnchorDeposit
	}

	return nil
}

// SubtractFromBalance subtracts the given amount from the account's balance.
func (acc *Account) SubtractFromBalance(amt amount.Amount) {
	acc.data.Balance -= amt
}

// AddToBalance adds the given amount to the account's balance.
func (acc *Account) AddToBalance(amt amount.Amount) {
	acc.data.Balance += amt
}

// Hash calculates and returns the hash of the account.
func (acc *Account) Hash() hash.Hash {
	bs, err := acc.Bytes()
	if err != nil {
		panic(err)
	}

	return hash.CalcHash(bs)
}

// SerializeSize returns the size in bytes required to serialize the account.
func (acc *Account) SerializeSize() int {
	if !acc.HasAnchor() {
		return accountHeaderSize
	}

	suffix := anchorFixedSuffix + len(acc.data.Anchor.RootHash) + len(acc.data.Anchor.ManifestURI)

	return accountHeaderSize + suffix
}

// Bytes returns the serialized byte representation of the account.
func (acc *Account) Bytes() ([]byte, error) {
	buf := bytes.NewBuffer(make([]byte, 0, acc.SerializeSize()))
	err := encoding.WriteElements(buf,
		acc.data.Number,
		acc.data.Balance)
	if err != nil {
		return nil, err
	}

	if acc.data.Anchor == nil {
		return buf.Bytes(), nil
	}

	if err := writeAnchor(buf, acc.data.Anchor); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func writeAnchor(buf *bytes.Buffer, anchor *AnchorData) error {
	err := encoding.WriteElements(buf, uint8(anchorFlagSet), uint8(len(anchor.RootHash)))
	if err != nil {
		return err
	}
	if _, err := buf.Write(anchor.RootHash); err != nil {
		return err
	}
	if err := encoding.WriteElement(buf, uint8(len(anchor.ManifestURI))); err != nil {
		return err
	}
	if _, err := buf.WriteString(anchor.ManifestURI); err != nil {
		return err
	}

	return encoding.WriteElements(buf,
		anchor.AnchorType,
		anchor.LockedDeposit,
		anchor.CreatedAtHeight,
		anchor.CreatedAtTime,
		anchor.UpdatedAtHeight,
		anchor.UpdatedAtTime,
	)
}

// Clone creates a deep copy of the account.
func (acc *Account) Clone() *Account {
	cloned := new(Account)
	*cloned = *acc
	if acc.data.Anchor == nil {
		return cloned
	}

	anchor := *acc.data.Anchor
	anchor.RootHash = append([]byte(nil), acc.data.Anchor.RootHash...)
	cloned.data.Anchor = &anchor

	return cloned
}
