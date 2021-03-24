package proposal

import (
	"encoding/json"
	"fmt"

	"github.com/fxamacker/cbor/v2"
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/errors"
)

type Proposal struct {
	data proposalData
}
type proposalData struct {
	Height    int               `cbor:"1,keyasint"`
	Round     int               `cbor:"2,keyasint"`
	Block     block.Block       `cbor:"3,keyasint"`
	Signature *crypto.Signature `cbor:"4,keyasint"`
}

func NewProposal(height int, round int, block block.Block) *Proposal {
	return &Proposal{
		data: proposalData{
			Height: height,
			Round:  round,
			Block:  block,
		},
	}
}
func (p *Proposal) Height() int                  { return p.data.Height }
func (p *Proposal) Round() int                   { return p.data.Round }
func (p *Proposal) Block() *block.Block          { return &p.data.Block }
func (p *Proposal) Signature() *crypto.Signature { return p.data.Signature }

func (p *Proposal) SanityCheck() error {
	if err := p.data.Block.SanityCheck(); err != nil {
		return errors.Errorf(errors.ErrInvalidProposal, err.Error())
	}
	if p.data.Height <= 0 {
		return errors.Errorf(errors.ErrInvalidProposal, "Invalid round")
	}
	if p.data.Round < 0 {
		return errors.Errorf(errors.ErrInvalidProposal, "Invalid round")
	}
	if p.data.Signature == nil {
		return errors.Errorf(errors.ErrInvalidProposal, "No signature")
	}
	if p.data.Signature.SanityCheck() != nil {
		return errors.Errorf(errors.ErrInvalidProposal, "Invalid signature")
	}
	return nil
}

func (p *Proposal) SetSignature(sig crypto.Signature) {
	p.data.Signature = &sig
}

// SetPublicKey is doing nothing and just satisfies SignableMsg interface
func (p *Proposal) SetPublicKey(crypto.PublicKey) {}

func (p *Proposal) SignBytes() []byte {
	type signProposal struct {
		Height    int         `cbor:"1,keyasint"`
		Round     int         `cbor:"2,keyasint"`
		BlockHash crypto.Hash `cbor:"3,keyasint"`
	}
	bz, _ := cbor.Marshal(signProposal{
		Height:    p.data.Height,
		Round:     p.data.Round,
		BlockHash: p.data.Block.Hash(),
	})
	return bz
}

func (p *Proposal) MarshalCBOR() ([]byte, error) {
	return cbor.Marshal(p.data)
}

func (p *Proposal) UnmarshalCBOR(bs []byte) error {
	return cbor.Unmarshal(bs, &p.data)
}

func (p Proposal) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.data)
}

func (p *Proposal) UnmarshalJSON(bs []byte) error {
	return json.Unmarshal(bs, &p.data)
}

func (p *Proposal) Verify(pubKey crypto.PublicKey) error {
	if p.data.Signature == nil {
		return errors.Errorf(errors.ErrInvalidProposal, "No signature")
	}
	if !pubKey.Address().EqualsTo(p.data.Block.Header().ProposerAddress()) {
		return errors.Errorf(errors.ErrInvalidProposal, "Invalid proposer")
	}
	if !pubKey.Verify(p.SignBytes(), *p.data.Signature) {
		return errors.Errorf(errors.ErrInvalidProposal, "Invalid signature")
	}
	return nil
}
func (p *Proposal) Hash() crypto.Hash {
	return crypto.HashH(p.SignBytes())
}

func (p *Proposal) IsForBlock(hash *crypto.Hash) bool {
	if hash == nil {
		return false
	}
	return p.Block().HashesTo(*hash)
}

func (p Proposal) Fingerprint() string {
	b := p.Block()
	return fmt.Sprintf("{%v/%v 🗃 %v}", p.data.Height, p.data.Round, b.Fingerprint())
}

// ---------
// For tests
func GenerateTestProposal(height, round int) (*Proposal, crypto.PrivateKey) {
	addr, _, pv := crypto.GenerateTestKeyPair()
	b, _ := block.GenerateTestBlock(&addr, nil)
	p := NewProposal(height, round, *b)
	sig := pv.Sign(p.SignBytes())
	p.SetSignature(sig)
	return p, pv
}
