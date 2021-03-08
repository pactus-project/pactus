package block

import (
	"encoding/json"

	"github.com/fxamacker/cbor/v2"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/util"
)

type Certificate struct {
	data certificateData
}
type certificateData struct {
	BlockHash  crypto.Hash      `cbor:"1,keyasint"`
	Round      int              `cbor:"2,keyasint"`
	Committers []int            `cbor:"3,keyasint"`
	Absences   []int            `cbor:"4,keyasint"`
	Signature  crypto.Signature `cbor:"8,keyasint"`
}

func NewCertificate(blockHash crypto.Hash, round int, committers, absences []int, signature crypto.Signature) *Certificate {
	return &Certificate{
		data: certificateData{
			BlockHash:  blockHash,
			Round:      round,
			Committers: committers,
			Absences:   absences,
			Signature:  signature,
		},
	}
}

func (cert *Certificate) BlockHash() crypto.Hash      { return cert.data.BlockHash }
func (cert *Certificate) Round() int                  { return cert.data.Round }
func (cert *Certificate) Committers() []int           { return cert.data.Committers }
func (cert *Certificate) Absences() []int             { return cert.data.Absences }
func (cert *Certificate) Signature() crypto.Signature { return cert.data.Signature }

func (cert *Certificate) SanityCheck() error {
	if err := cert.data.BlockHash.SanityCheck(); err != nil {
		return errors.Errorf(errors.ErrInvalidBlock, err.Error())
	}
	if cert.data.Round < 0 {
		return errors.Errorf(errors.ErrInvalidBlock, "Invalid Round")
	}
	if err := cert.data.Signature.SanityCheck(); err != nil {
		return errors.Errorf(errors.ErrInvalidBlock, err.Error())
	}
	if cert.data.Committers == nil {
		return errors.Errorf(errors.ErrInvalidBlock, "Invalid committers")
	}
	if cert.data.Absences == nil {
		return errors.Errorf(errors.ErrInvalidBlock, "Invalid absences")
	}
	signedBy := util.Subtracts(cert.data.Committers, cert.data.Absences)
	if !util.Equal(util.Subtracts(cert.data.Committers, signedBy), cert.data.Absences) {
		return errors.Errorf(errors.ErrInvalidBlock, "Absences is not subset of committers")
	}

	return nil
}

func (cert *Certificate) Hash() crypto.Hash {
	if cert == nil {
		return crypto.UndefHash
	}
	bs, err := cert.MarshalCBOR()
	if err != nil {
		return crypto.UndefHash
	}
	return crypto.HashH(bs)
}

func (cert *Certificate) CommitteeHash() crypto.Hash {
	bz, _ := cbor.Marshal(cert.data.Committers)
	return crypto.HashH(bz)
}

func (cert *Certificate) MarshalCBOR() ([]byte, error) {
	return cbor.Marshal(cert.data)
}

func (cert *Certificate) UnmarshalCBOR(bs []byte) error {
	return cbor.Unmarshal(bs, &cert.data)
}

func (cert *Certificate) MarshalJSON() ([]byte, error) {
	return json.Marshal(cert.data)
}

func (cert *Certificate) UnmarshalJSON(bz []byte) error {
	return json.Unmarshal(bz, &cert.data)
}

type signVote struct {
	BlockHash crypto.Hash `cbor:"1,keyasint"`
	Round     int         `cbor:"2,keyasint"`
}

func (cert *Certificate) SignBytes() []byte {
	return CertificateSignBytes(cert.data.BlockHash, cert.data.Round)
}

func CertificateSignBytes(blockHash crypto.Hash, round int) []byte {
	bz, _ := cbor.Marshal(signVote{
		Round:     round,
		BlockHash: blockHash,
	})

	return bz
}

func GenerateTestCertificate(blockhash crypto.Hash) *Certificate {
	_, _, priv2 := crypto.GenerateTestKeyPair()
	_, _, priv3 := crypto.GenerateTestKeyPair()
	_, _, priv4 := crypto.GenerateTestKeyPair()

	sigs := []crypto.Signature{
		priv2.Sign(blockhash.RawBytes()),
		priv3.Sign(blockhash.RawBytes()),
		priv4.Sign(blockhash.RawBytes()),
	}
	sig := crypto.Aggregate(sigs)

	return NewCertificate(
		blockhash,
		util.RandInt(10),
		[]int{0, 1, 2, 3},
		[]int{0},
		sig)
}
