package block

import (
	"encoding/json"
	"time"

	"github.com/fxamacker/cbor/v2"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/sortition"
)

// TODO: try to memorize hash for better performance
type Header struct {
	data headerData
}
type headerData struct {
	Version          int            `cbor:"1,keyasint"`
	UnixTime         int64          `cbor:"2,keyasint"`
	LastBlockHash    crypto.Hash    `cbor:"3,keyasint"`
	StateHash        crypto.Hash    `cbor:"4,keyasint"`
	TxIDsHash        crypto.Hash    `cbor:"5,keyasint"`
	LastReceiptsHash crypto.Hash    `cbor:"6,keyasint"`
	LastCommitHash   crypto.Hash    `cbor:"7,keyasint"`
	CommitteeHash    crypto.Hash    `cbor:"8,keyasint"`
	SortitionSeed    sortition.Seed `cbor:"9,keyasint"`
	ProposerAddress  crypto.Address `cbor:"10,keyasint"`
}

func (h Header) Version() int                    { return h.data.Version }
func (h Header) Time() time.Time                 { return time.Unix(h.data.UnixTime, 0) }
func (h Header) TxIDsHash() crypto.Hash          { return h.data.TxIDsHash }
func (h Header) StateHash() crypto.Hash          { return h.data.StateHash }
func (h Header) LastBlockHash() crypto.Hash      { return h.data.LastBlockHash }
func (h Header) LastReceiptsHash() crypto.Hash   { return h.data.LastReceiptsHash }
func (h Header) LastCommitHash() crypto.Hash     { return h.data.LastCommitHash }
func (h Header) CommitteeHash() crypto.Hash      { return h.data.CommitteeHash }
func (h Header) SortitionSeed() sortition.Seed   { return h.data.SortitionSeed }
func (h Header) ProposerAddress() crypto.Address { return h.data.ProposerAddress }

func NewHeader(version int,
	time time.Time,
	txIDsHash, lastBlockHash, committeeHash, stateHash, lastReceiptsHash, lastCommitHash crypto.Hash,
	sortitionSeed sortition.Seed, proposerAddress crypto.Address) Header {

	return Header{
		data: headerData{
			Version:          version,
			UnixTime:         time.Unix(),
			TxIDsHash:        txIDsHash,
			LastBlockHash:    lastBlockHash,
			CommitteeHash:    committeeHash,
			StateHash:        stateHash,
			LastReceiptsHash: lastReceiptsHash,
			LastCommitHash:   lastCommitHash,
			ProposerAddress:  proposerAddress,
			SortitionSeed:    sortitionSeed,
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
	if err := h.data.CommitteeHash.SanityCheck(); err != nil {
		return errors.Errorf(errors.ErrInvalidBlock, err.Error())
	}

	if h.data.LastCommitHash.IsUndef() {
		// Check for genesis block
		if !h.data.LastBlockHash.IsUndef() ||
			!h.data.LastReceiptsHash.IsUndef() {
			return errors.Errorf(errors.ErrInvalidBlock, "Invalid genesis block hash")
		}
	} else {
		if err := h.data.LastBlockHash.SanityCheck(); err != nil {
			return errors.Errorf(errors.ErrInvalidBlock, err.Error())
		}
		if err := h.data.LastReceiptsHash.SanityCheck(); err != nil {
			return errors.Errorf(errors.ErrInvalidBlock, err.Error())
		}
	}

	return nil
}

func (h Header) Hash() crypto.Hash {
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
