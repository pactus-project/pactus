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
	Version          uint           `cbor:"1,keyasint"`
	UnixTime         int64          `cbor:"2,keyasint"`
	LastBlockHash    crypto.Hash    `cbor:"3,keyasint"`
	StateHash        crypto.Hash    `cbor:"4,keyasint"`
	TxsHash          crypto.Hash    `cbor:"5,keyasint"`
	LastReceiptsHash crypto.Hash    `cbor:"6,keyasint"`
	CommitersHash    crypto.Hash    `cbor:"7,keyasint"`
	ProposerAddress  crypto.Address `cbor:"8,keyasint"`
	LastCommit       *Commit        `cbor:"9,keyasint"`
}

func (h *Header) Version() uint                   { return h.data.Version }
func (h *Header) Time() time.Time                 { return time.Unix(h.data.UnixTime, 0) }
func (h *Header) TxsHash() crypto.Hash            { return h.data.TxsHash }
func (h *Header) StateHash() crypto.Hash          { return h.data.StateHash }
func (h *Header) LastBlockHash() crypto.Hash      { return h.data.LastBlockHash }
func (h *Header) LastReceiptsHash() crypto.Hash   { return h.data.LastReceiptsHash }
func (h *Header) CommitersHash() crypto.Hash      { return h.data.CommitersHash }
func (h *Header) ProposerAddress() crypto.Address { return h.data.ProposerAddress }
func (h *Header) LastCommit() *Commit             { return h.data.LastCommit }

func NewHeader(version uint,
	time time.Time,
	txsHash, lastBlockHash, CommitersHash, stateHash, lastReceiptsHash crypto.Hash,
	proposerAddress crypto.Address,
	lastCommit *Commit) Header {

	return Header{
		data: headerData{
			Version:          version,
			UnixTime:         time.Unix(),
			TxsHash:          txsHash,
			LastBlockHash:    lastBlockHash,
			CommitersHash:    CommitersHash,
			StateHash:        stateHash,
			LastReceiptsHash: lastReceiptsHash,
			ProposerAddress:  proposerAddress,
			LastCommit:       lastCommit,
		},
	}
}

func (h *Header) SanityCheck() error {
	if err := h.data.TxsHash.SanityCheck(); err != nil {
		return errors.Errorf(errors.ErrInvalidBlock, err.Error())
	}
	if err := h.data.ProposerAddress.SanityCheck(); err != nil {
		return errors.Errorf(errors.ErrInvalidBlock, err.Error())
	}
	// TODO: fix ma later
	// if err := h.data.CommitersHash.SanityCheck(); err != nil {
	// 	return errors.Errorf(errors.ErrInvalidBlock, err.Error())
	// }

	if h.data.LastCommit != nil {
		if err := h.data.LastBlockHash.SanityCheck(); err != nil {
			return errors.Errorf(errors.ErrInvalidBlock, err.Error())
		}
		if err := h.data.LastReceiptsHash.SanityCheck(); err != nil {
			return errors.Errorf(errors.ErrInvalidBlock, err.Error())
		}
		if err := h.data.LastCommit.SanityCheck(); err != nil {
			return err
		}
	} else {
		// Check for genesis block
		if !h.data.LastBlockHash.IsUndef() ||
			!h.data.LastReceiptsHash.IsUndef() {
			return errors.Errorf(errors.ErrInvalidBlock, "Invalid Genesis block")
		}
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
