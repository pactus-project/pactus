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

// TODO: try to memorize hash for better performance
type Header struct {
	data headerData
}
type headerData struct {
	Version             int            `cbor:"1,keyasint"`
	UnixTime            int64          `cbor:"2,keyasint"`
	LastBlockHash       hash.Hash      `cbor:"3,keyasint"`
	StateHash           hash.Hash      `cbor:"4,keyasint"`
	TxIDsHash           hash.Hash      `cbor:"5,keyasint"`
	LastCertificateHash hash.Hash      `cbor:"6,keyasint"`
	SortitionSeed       sortition.Seed `cbor:"7,keyasint"`
	ProposerAddress     crypto.Address `cbor:"8,keyasint"`
}

func (h *Header) Version() int                    { return h.data.Version }
func (h *Header) Time() time.Time                 { return time.Unix(h.data.UnixTime, 0) }
func (h *Header) TxIDsHash() hash.Hash            { return h.data.TxIDsHash }
func (h *Header) StateHash() hash.Hash            { return h.data.StateHash }
func (h *Header) LastBlockHash() hash.Hash        { return h.data.LastBlockHash }
func (h *Header) LastCertificateHash() hash.Hash  { return h.data.LastCertificateHash }
func (h *Header) SortitionSeed() sortition.Seed   { return h.data.SortitionSeed }
func (h *Header) ProposerAddress() crypto.Address { return h.data.ProposerAddress }

func NewHeader(version int,
	time time.Time,
	txIDsHash, lastBlockHash, stateHash, lastCertificateHash hash.Hash,
	sortitionSeed sortition.Seed, proposerAddress crypto.Address) Header {

	return Header{
		data: headerData{
			Version:             version,
			UnixTime:            time.Unix(),
			TxIDsHash:           txIDsHash,
			LastBlockHash:       lastBlockHash,
			StateHash:           stateHash,
			LastCertificateHash: lastCertificateHash,
			ProposerAddress:     proposerAddress,
			SortitionSeed:       sortitionSeed,
		},
	}
}

func (h *Header) SanityCheck() error {
	if err := h.data.StateHash.SanityCheck(); err != nil {
		return errors.Errorf(errors.ErrInvalidBlock, err.Error())
	}
	if err := h.data.TxIDsHash.SanityCheck(); err != nil {
		return errors.Errorf(errors.ErrInvalidBlock, err.Error())
	}
	if err := h.data.ProposerAddress.SanityCheck(); err != nil {
		return errors.Errorf(errors.ErrInvalidBlock, err.Error())
	}

	if h.data.LastCertificateHash.IsUndef() {
		// Check for genesis block
		if !h.data.LastBlockHash.IsUndef() {
			return errors.Errorf(errors.ErrInvalidBlock, "invalid Last Block hash")
		}
	} else {
		if err := h.data.LastBlockHash.SanityCheck(); err != nil {
			return errors.Errorf(errors.ErrInvalidBlock, err.Error())
		}
	}

	return nil
}

func (h Header) Hash() hash.Hash {
	bs, err := h.MarshalCBOR()
	if err != nil {
		return hash.UndefHash
	}
	return hash.CalcHash(bs)
}

func (h *Header) MarshalCBOR() ([]byte, error) {
	return cbor.Marshal(h.data)
}

func (h *Header) UnmarshalCBOR(bs []byte) error {
	return cbor.Unmarshal(bs, &h.data)
}

func (h Header) MarshalJSON() ([]byte, error) {
	return json.Marshal(h.data)
}
