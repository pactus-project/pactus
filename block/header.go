package block

import (
	"encoding/json"
	"time"

	"github.com/fxamacker/cbor/v2"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/crypto/hash"
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/sortition"
)

type Header struct {
	memorizedHash hash.Hash
	data          headerData
}
type headerData struct {
	Version         int                      `cbor:"1,keyasint"`
	UnixTime        int64                    `cbor:"2,keyasint"`
	PrevBlockHash   hash.Hash                `cbor:"3,keyasint"`
	PrevCertHash    hash.Hash                `cbor:"4,keyasint"`
	StateRoot       hash.Hash                `cbor:"5,keyasint"`
	TxsRoot         hash.Hash                `cbor:"6,keyasint"`
	SortitionSeed   sortition.VerifiableSeed `cbor:"7,keyasint"`
	ProposerAddress crypto.Address           `cbor:"8,keyasint"`
}

func (h *Header) Version() int                            { return h.data.Version }
func (h *Header) Time() time.Time                         { return time.Unix(h.data.UnixTime, 0) }
func (h *Header) TxsRoot() hash.Hash                      { return h.data.TxsRoot }
func (h *Header) StateRoot() hash.Hash                    { return h.data.StateRoot }
func (h *Header) PrevBlockHash() hash.Hash                { return h.data.PrevBlockHash }
func (h *Header) PrevCertificateHash() hash.Hash          { return h.data.PrevCertHash }
func (h *Header) SortitionSeed() sortition.VerifiableSeed { return h.data.SortitionSeed }
func (h *Header) ProposerAddress() crypto.Address         { return h.data.ProposerAddress }

func NewHeader(version int, time time.Time,
	txsRoot, stateRoot, prevBlockHash, prevCertHash hash.Hash,
	sortitionSeed sortition.VerifiableSeed, proposerAddress crypto.Address) Header {

	h := Header{
		data: headerData{
			Version:         version,
			UnixTime:        time.Unix(),
			TxsRoot:         txsRoot,
			PrevBlockHash:   prevBlockHash,
			StateRoot:       stateRoot,
			PrevCertHash:    prevCertHash,
			ProposerAddress: proposerAddress,
			SortitionSeed:   sortitionSeed,
		},
	}
	h.memorizedHash = h.calcHash()
	return h
}

func (h *Header) SanityCheck() error {
	if err := h.data.StateRoot.SanityCheck(); err != nil {
		return errors.Errorf(errors.ErrInvalidBlock, "invalid state root")
	}
	if err := h.data.TxsRoot.SanityCheck(); err != nil {
		return errors.Errorf(errors.ErrInvalidBlock, "invalid transactions root")
	}
	if err := h.data.ProposerAddress.SanityCheck(); err != nil {
		return errors.Errorf(errors.ErrInvalidBlock, "invalid proposer address")
	}
	if h.data.SortitionSeed.IsUndef() {
		return errors.Errorf(errors.ErrInvalidBlock, "invalid sortition seed")
	}

	if h.data.PrevCertHash.IsUndef() {
		// Genesis block checks
		if !h.data.PrevBlockHash.IsUndef() {
			return errors.Errorf(errors.ErrInvalidBlock, "invalid previous block hash")
		}
	} else {
		if err := h.data.PrevBlockHash.SanityCheck(); err != nil {
			return errors.Errorf(errors.ErrInvalidBlock, err.Error())
		}
	}

	return nil
}

func (h *Header) calcHash() hash.Hash {
	bs, _ := h.Encode()
	return hash.CalcHash(bs)
}

func (h *Header) Hash() hash.Hash {
	return h.memorizedHash
}

func (h *Header) MarshalCBOR() ([]byte, error) {
	return h.Encode()
}

func (h *Header) UnmarshalCBOR(bs []byte) error {
	return h.Decode(bs)
}

func (h *Header) Encode() ([]byte, error) {
	return cbor.Marshal(h.data)
}

func (h *Header) Decode(bs []byte) error {
	if err := cbor.Unmarshal(bs, &h.data); err != nil {
		return err
	}

	h.memorizedHash = h.calcHash()
	return nil
}

func (h Header) MarshalJSON() ([]byte, error) {
	return json.Marshal(h.data)
}
