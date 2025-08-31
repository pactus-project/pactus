package block

import (
	"fmt"
	"io"
	"time"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/sortition"
	"github.com/pactus-project/pactus/types/protocol"
	"github.com/pactus-project/pactus/util/encoding"
)

type Header struct {
	data headerData
}

type headerData struct {
	Version         protocol.Version
	UnixTime        uint32
	PrevBlockHash   hash.Hash
	StateRoot       hash.Hash
	SortitionSeed   sortition.VerifiableSeed
	ProposerAddress crypto.Address
}

// Version returns the block version.
func (h *Header) Version() protocol.Version {
	return h.data.Version
}

// Time returns the block time.
func (h *Header) Time() time.Time {
	return time.Unix(int64(h.data.UnixTime), 0)
}

// UnixTime returns the block time in Unix value.
func (h *Header) UnixTime() uint32 {
	return h.data.UnixTime
}

// StateRoot returns the state root hash.
func (h *Header) StateRoot() hash.Hash {
	return h.data.StateRoot
}

// PrevBlockHash returns the previous block hash.
func (h *Header) PrevBlockHash() hash.Hash {
	return h.data.PrevBlockHash
}

// SortitionSeed returns the sortition seed.
func (h *Header) SortitionSeed() sortition.VerifiableSeed {
	return h.data.SortitionSeed
}

// ProposerAddress returns the proposer address.
func (h *Header) ProposerAddress() crypto.Address {
	return h.data.ProposerAddress
}

func NewHeader(version protocol.Version, tme time.Time, stateRoot, prevBlockHash hash.Hash,
	sortitionSeed sortition.VerifiableSeed, proposerAddress crypto.Address,
) *Header {
	return &Header{
		data: headerData{
			Version:         version,
			UnixTime:        uint32(tme.Unix()),
			PrevBlockHash:   prevBlockHash,
			StateRoot:       stateRoot,
			ProposerAddress: proposerAddress,
			SortitionSeed:   sortitionSeed,
		},
	}
}

func (h *Header) BasicCheck() error {
	if h.data.Version == 0 {
		return BasicCheckError{
			Reason: "invalid block version: 0",
		}
	}

	if !h.data.ProposerAddress.IsValidatorAddress() {
		return BasicCheckError{
			Reason: fmt.Sprintf("invalid proposer address: %s",
				h.data.ProposerAddress.String()),
		}
	}

	return nil
}

// SerializeSize returns the number of bytes it would take to serialize the header.
func (*Header) SerializeSize() int {
	return 138 // 5 + (2 * 32) + 48 + 21
}

func (h *Header) Encode(w io.Writer) error {
	return encoding.WriteElements(w,
		h.data.Version,
		h.data.UnixTime,
		h.data.PrevBlockHash,
		h.data.StateRoot,
		h.data.SortitionSeed,
		h.data.ProposerAddress)
}

func (h *Header) Decode(r io.Reader) error {
	return encoding.ReadElements(r,
		&h.data.Version,
		&h.data.UnixTime,
		&h.data.PrevBlockHash,
		&h.data.StateRoot,
		&h.data.SortitionSeed,
		&h.data.ProposerAddress)
}
