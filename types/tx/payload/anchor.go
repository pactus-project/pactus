package payload

import (
	"fmt"
	"io"
	"unicode/utf8"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/util/encoding"
)

const (
	// AnchorActionSet writes an anchor.
	AnchorActionSet uint8 = 0
	// AnchorActionDelete clears an anchor.
	AnchorActionDelete uint8 = 1

	minAnchorHashLen = 32
	maxAnchorHashLen = 64
	maxAnchorURILen  = 128
)

// AnchorPayload is a state anchor (PIP-50).
type AnchorPayload struct {
	From        crypto.Address
	Action      uint8
	RootHash    []byte
	ManifestURI string
	AnchorType  uint8
	Deposit     amount.Amount
}

func (*AnchorPayload) Type() Type {
	return TypeAnchor
}

func (p *AnchorPayload) Signer() crypto.Address {
	return p.From
}

func (p *AnchorPayload) Value() amount.Amount {
	if p.Action != AnchorActionSet {
		return 0
	}

	return p.Deposit
}

// BasicCheck performs basic checks on the anchor payload.
func (p *AnchorPayload) BasicCheck() error {
	if !isAnchorAccount(p.From) {
		return BasicCheckError{
			Reason: "sender is not an account address: " + p.From.String(),
		}
	}

	switch p.Action {
	case AnchorActionSet:
		return p.checkSet()
	case AnchorActionDelete:
		return nil
	default:
		return BasicCheckError{
			Reason: "invalid anchor action",
		}
	}
}

func (p *AnchorPayload) SerializeSize() int {
	size := p.From.SerializeSize() + 1
	if p.Action != AnchorActionSet {
		return size
	}

	return size + 1 + len(p.RootHash) + 1 + len(p.ManifestURI) + 1 +
		encoding.VarIntSerializeSize(uint64(p.Deposit))
}

func (p *AnchorPayload) Encode(w io.Writer) error {
	if err := p.From.Encode(w); err != nil {
		return err
	}
	if err := encoding.WriteElement(w, p.Action); err != nil {
		return err
	}
	if p.Action == AnchorActionDelete {
		return nil
	}
	if p.Action != AnchorActionSet {
		return BasicCheckError{
			Reason: "invalid anchor action",
		}
	}
	if err := p.checkSet(); err != nil {
		return err
	}
	if err := encoding.WriteElement(w, uint8(len(p.RootHash))); err != nil {
		return err
	}
	if _, err := w.Write(p.RootHash); err != nil {
		return err
	}
	if err := encoding.WriteElement(w, uint8(len(p.ManifestURI))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, p.ManifestURI); err != nil {
		return err
	}
	if err := encoding.WriteElement(w, p.AnchorType); err != nil {
		return err
	}

	return p.Deposit.Encode(w)
}

func (p *AnchorPayload) Decode(_ DecodeContext, r io.Reader) error {
	decoded, err := readAnchor(r)
	if err != nil {
		*p = AnchorPayload{}

		return err
	}
	*p = *decoded

	return nil
}

// LogString returns a concise string representation intended for use in logs.
func (p *AnchorPayload) LogString() string {
	return fmt.Sprintf("{Anchor ⚓ %s %d",
		p.From.LogString(),
		p.Action)
}

func (p *AnchorPayload) checkSet() error {
	hashLen := len(p.RootHash)
	if hashLen < minAnchorHashLen || hashLen > maxAnchorHashLen {
		return BasicCheckError{
			Reason: "anchor hash length out of range",
		}
	}
	if len(p.ManifestURI) > maxAnchorURILen {
		return BasicCheckError{
			Reason: "anchor uri length out of range",
		}
	}
	if !utf8.ValidString(p.ManifestURI) {
		return BasicCheckError{
			Reason: "invalid anchor uri",
		}
	}
	if p.Deposit < 0 || p.Deposit > amount.MaxNanoPAC {
		return BasicCheckError{
			Reason: "anchor deposit out of range",
		}
	}

	return nil
}

func readAnchor(r io.Reader) (*AnchorPayload, error) {
	decoded := &AnchorPayload{}
	if err := decoded.From.Decode(r); err != nil {
		return nil, err
	}

	if err := encoding.ReadElement(r, &decoded.Action); err != nil {
		return nil, err
	}

	switch decoded.Action {
	case AnchorActionDelete:
		return decoded, nil
	case AnchorActionSet:
		if err := readAnchorSet(r, decoded); err != nil {
			return nil, err
		}

		return decoded, nil
	default:
		return nil, BasicCheckError{
			Reason: "invalid anchor action",
		}
	}
}

func readAnchorSet(r io.Reader, decoded *AnchorPayload) error {
	hashLen := uint8(0)
	if err := encoding.ReadElement(r, &hashLen); err != nil {
		return err
	}
	if hashLen < minAnchorHashLen || hashLen > maxAnchorHashLen {
		return BasicCheckError{
			Reason: "anchor hash length out of range",
		}
	}

	root := make([]byte, int(hashLen))
	if _, err := io.ReadFull(r, root); err != nil {
		return err
	}

	uriLen := uint8(0)
	if err := encoding.ReadElement(r, &uriLen); err != nil {
		return err
	}
	if uriLen > maxAnchorURILen {
		return BasicCheckError{
			Reason: "anchor uri length out of range",
		}
	}

	uri := make([]byte, int(uriLen))
	if _, err := io.ReadFull(r, uri); err != nil {
		return err
	}
	if !utf8.Valid(uri) {
		return BasicCheckError{
			Reason: "invalid anchor uri",
		}
	}

	if err := encoding.ReadElement(r, &decoded.AnchorType); err != nil {
		return err
	}

	var deposit amount.Amount
	if err := deposit.Decode(r); err != nil {
		return err
	}
	if deposit < 0 || deposit > amount.MaxNanoPAC {
		return BasicCheckError{
			Reason: "anchor deposit out of range",
		}
	}

	decoded.RootHash = root
	decoded.ManifestURI = string(uri)
	decoded.Deposit = deposit

	return nil
}

func isAnchorAccount(addr crypto.Address) bool {
	switch addr.Type() {
	case crypto.AddressTypeBLSAccount,
		crypto.AddressTypeEd25519Account,
		crypto.AddressTypeSecp256k1Account:
		return true
	case crypto.AddressTypeTreasury,
		crypto.AddressTypeValidator:
		return false
	default:
		return false
	}
}
