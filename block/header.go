package block

import (
	"encoding/json"
	"time"

	"github.com/fxamacker/cbor/v2"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/errors"
)

// TODO: try to memorize hash for better performance
type Header struct {
	data headerData
}
type headerData struct {
	Version            uint           `cbor:"1,keyasint"`
	UnixTime           int64          `cbor:"2,keyasint"`
	TxsHash            crypto.Hash    `cbor:"3,keyasint"`
	StateHash          crypto.Hash    `cbor:"4,keyasint"`
	LastBlockHash      crypto.Hash    `cbor:"5,keyasint"`
	LastCommitHash     crypto.Hash    `cbor:"6,keyasint"`
	LastReceiptsHash   crypto.Hash    `cbor:"7,keyasint"`
	NextValidatorsHash crypto.Hash    `cbor:"8,keyasint"`
	ProposerAddress    crypto.Address `cbor:"9,keyasint"`
}

func (h *Header) Version() uint                   { return h.data.Version }
func (h *Header) Time() time.Time                 { return time.Unix(h.data.UnixTime, 0) }
func (h *Header) TxsHash() crypto.Hash            { return h.data.TxsHash }
func (h *Header) StateHash() crypto.Hash          { return h.data.StateHash }
func (h *Header) LastBlockHash() crypto.Hash      { return h.data.LastBlockHash }
func (h *Header) LastCommitHash() crypto.Hash     { return h.data.LastCommitHash }
func (h *Header) LastReceiptsHash() crypto.Hash   { return h.data.LastReceiptsHash }
func (h *Header) NextValidatorsHash() crypto.Hash { return h.data.NextValidatorsHash }
func (h *Header) ProposerAddress() crypto.Address { return h.data.ProposerAddress }

func NewHeader(version uint,
	time time.Time,
	txsHash, lastBlockHash, nextValHash, stateHash, lastReceiptsHash, lastCommitHash crypto.Hash,
	proposerAddress crypto.Address) Header {

	return Header{
		data: headerData{
			Version:            version,
			UnixTime:           time.Unix(),
			TxsHash:            txsHash,
			LastBlockHash:      lastBlockHash,
			LastCommitHash:     lastCommitHash,
			NextValidatorsHash: nextValHash,
			StateHash:          stateHash,
			LastReceiptsHash:   lastReceiptsHash,
			ProposerAddress:    proposerAddress,
		},
	}
}

func (h *Header) SanityCheck() error {
	if err := h.data.LastBlockHash.SanityCheck(); err != nil {
		return errors.Errorf(errors.ErrInvalidBlock, err.Error())
	}
	if err := h.data.LastCommitHash.SanityCheck(); err != nil {
		return errors.Errorf(errors.ErrInvalidBlock, err.Error())
	}
	if err := h.data.TxsHash.SanityCheck(); err != nil {
		return errors.Errorf(errors.ErrInvalidBlock, err.Error())
	}
	if err := h.data.ProposerAddress.SanityCheck(); err != nil {
		return errors.Errorf(errors.ErrInvalidBlock, err.Error())
	}
	if err := h.data.NextValidatorsHash.SanityCheck(); err != nil {
		return errors.Errorf(errors.ErrInvalidBlock, err.Error())
	}
	if err := h.data.LastReceiptsHash.SanityCheck(); err != nil {
		return errors.Errorf(errors.ErrInvalidBlock, err.Error())
	}

	return nil
}

func (h *Header) Hash() crypto.Hash {
	bs, err := h.MarshalCBOR()
	if err != nil {
		return crypto.UndefHash
	}
	return crypto.HashH(bs)
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

func (h *Header) UnmarshalJSON(bz []byte) error {
	return json.Unmarshal(bz, &h.data)
}
