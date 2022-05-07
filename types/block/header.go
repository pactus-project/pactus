package block

import (
	"io"
	"time"

	"github.com/zarbchain/zarb-go/sortition"
	"github.com/zarbchain/zarb-go/types/crypto"
	"github.com/zarbchain/zarb-go/types/crypto/hash"
	"github.com/zarbchain/zarb-go/util/encoding"
	"github.com/zarbchain/zarb-go/util/errors"
)

type Header struct {
	data headerData
}

type headerData struct {
	Version         uint8
	UnixTime        uint32
	PrevBlockHash   hash.Hash
	StateRoot       hash.Hash
	SortitionSeed   sortition.VerifiableSeed
	ProposerAddress crypto.Address
}

func (h *Header) Version() uint8                          { return h.data.Version }
func (h *Header) Time() time.Time                         { return time.Unix(int64(h.data.UnixTime), 0) }
func (h *Header) StateRoot() hash.Hash                    { return h.data.StateRoot }
func (h *Header) PrevBlockHash() hash.Hash                { return h.data.PrevBlockHash }
func (h *Header) SortitionSeed() sortition.VerifiableSeed { return h.data.SortitionSeed }
func (h *Header) ProposerAddress() crypto.Address         { return h.data.ProposerAddress }

func NewHeader(version uint8, time time.Time, stateRoot, prevBlockHash hash.Hash,
	sortitionSeed sortition.VerifiableSeed, proposerAddress crypto.Address) Header {
	h := Header{
		data: headerData{
			Version:         version,
			UnixTime:        uint32(time.Unix()),
			PrevBlockHash:   prevBlockHash,
			StateRoot:       stateRoot,
			ProposerAddress: proposerAddress,
			SortitionSeed:   sortitionSeed,
		},
	}
	return h
}

func (h *Header) SanityCheck() error {
	if err := h.data.StateRoot.SanityCheck(); err != nil {
		return errors.Errorf(errors.ErrInvalidBlock, "invalid state root")
	}
	if err := h.data.ProposerAddress.SanityCheck(); err != nil {
		return errors.Errorf(errors.ErrInvalidBlock, "invalid proposer address")
	}

	return nil
}

// SerializeSize returns the number of bytes it would take to serialize the header.
func (h *Header) SerializeSize() int {
	return 138 // 5 + (2 * 32) + 48 + 21
}

// Encode encodes the receiver to w.
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
