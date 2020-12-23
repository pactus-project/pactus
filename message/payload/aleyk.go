package payload

import (
	"fmt"

	peer "github.com/libp2p/go-libp2p-peer"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/util"
	"github.com/zarbchain/zarb-go/version"
)

const SalamResponseCodeOK = 0
const SalamResponseCodeRejected = 1

type SalamResponse struct {
	Status  int    `cbor:"1,keyasint"`
	Message string `cbor:"2,keyasint,omitempty"`
}

type AleykPayload struct {
	NodeVersion version.Version  `cbor:"1,keyasint"`
	Moniker     string           `cbor:"2,keyasint"`
	PublicKey   crypto.PublicKey `cbor:"3,keyasint"`
	PeerID      peer.ID          `cbor:"4,keyasint"`
	GenesisHash crypto.Hash      `cbor:"5,keyasint"`
	Height      int              `cbor:"6,keyasint"`
	Flags       int              `cbor:"7,keyasint"`
	Response    SalamResponse    `cbor:"8,keyasint"`
}

func (p *AleykPayload) SanityCheck() error {
	if p.Height < 0 {
		return errors.Errorf(errors.ErrInvalidMessage, "invalid Height")
	}
	if err := p.PublicKey.SanityCheck(); err != nil {
		return err
	}
	if err := p.PeerID.Validate(); err != nil {
		return err
	}
	return nil
}

func (p *AleykPayload) Type() PayloadType {
	return PayloadTypeAleyk
}

func (p *AleykPayload) Fingerprint() string {
	return fmt.Sprintf("{%v %v}", util.FingerprintPeerID(p.PeerID), p.Height)
}
