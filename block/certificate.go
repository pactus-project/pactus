package block

import (
	"encoding/json"

	"github.com/fxamacker/cbor/v2"
	"github.com/zarbchain/zarb-go/crypto/bls"
	"github.com/zarbchain/zarb-go/crypto/hash"
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/util"
)

type Certificate struct {
	memorizedHash hash.Hash
	data          certificateData
}
type certificateData struct {
	BlockHash  hash.Hash      `cbor:"1,keyasint"`
	Round      int            `cbor:"2,keyasint"`
	Committers []int          `cbor:"3,keyasint"`
	Absentees  []int          `cbor:"4,keyasint"`
	Signature  *bls.Signature `cbor:"5,keyasint"`
}

func NewCertificate(blockHash hash.Hash, round int, committers, absentees []int, signature *bls.Signature) *Certificate {
	cert := &Certificate{
		data: certificateData{
			BlockHash:  blockHash,
			Round:      round,
			Committers: committers,
			Absentees:  absentees,
			Signature:  signature,
		},
	}

	cert.memorizedHash = cert.calcHash()
	return cert
}

func (cert *Certificate) BlockHash() hash.Hash      { return cert.data.BlockHash }
func (cert *Certificate) Round() int                { return cert.data.Round }
func (cert *Certificate) Committers() []int         { return cert.data.Committers }
func (cert *Certificate) Absentees() []int          { return cert.data.Absentees }
func (cert *Certificate) Signature() *bls.Signature { return cert.data.Signature }

func (cert *Certificate) SanityCheck() error {
	if err := cert.BlockHash().SanityCheck(); err != nil {
		return errors.Errorf(errors.ErrInvalidBlock, err.Error())
	}
	if cert.Round() < 0 {
		return errors.Errorf(errors.ErrInvalidBlock, "Invalid Round")
	}
	if err := cert.Signature().SanityCheck(); err != nil {
		return errors.Errorf(errors.ErrInvalidBlock, err.Error())
	}
	if cert.Committers() == nil {
		return errors.Errorf(errors.ErrInvalidBlock, "Invalid committers")
	}
	if cert.data.Absentees == nil {
		return errors.Errorf(errors.ErrInvalidBlock, "invalid absentees")
	}
	signedBy := util.Subtracts(cert.Committers(), cert.Absentees())
	if !util.Equal(util.Subtracts(cert.Committers(), signedBy), cert.Absentees()) {
		return errors.Errorf(errors.ErrInvalidBlock, "absentees is not subset of committers")
	}

	return nil
}

func (cert *Certificate) calcHash() hash.Hash {
	bs, _ := cert.Encode()
	return hash.CalcHash(bs)
}

func (cert *Certificate) Hash() hash.Hash {
	return cert.memorizedHash
}

func (cert *Certificate) MarshalCBOR() ([]byte, error) {
	return cert.Encode()
}

func (cert *Certificate) UnmarshalCBOR(bs []byte) error {
	return cert.Decode(bs)
}
func (cert *Certificate) Encode() ([]byte, error) {
	return cbor.Marshal(cert.data)
}

func (cert *Certificate) Decode(bs []byte) error {
	if err := cbor.Unmarshal(bs, &cert.data); err != nil {
		return err
	}

	cert.memorizedHash = cert.calcHash()
	return nil
}

func (cert *Certificate) MarshalJSON() ([]byte, error) {
	return json.Marshal(cert.data)
}

type signVote struct {
	BlockHash hash.Hash `cbor:"1,keyasint"`
	Round     int       `cbor:"2,keyasint"`
}

func (cert *Certificate) SignBytes() []byte {
	return CertificateSignBytes(cert.data.BlockHash, cert.data.Round)
}

func CertificateSignBytes(blockHash hash.Hash, round int) []byte {
	bz, _ := cbor.Marshal(signVote{
		Round:     round,
		BlockHash: blockHash,
	})

	return bz
}

func GenerateTestCertificate(blockHash hash.Hash) *Certificate {
	_, priv2 := bls.GenerateTestKeyPair()
	_, priv3 := bls.GenerateTestKeyPair()
	_, priv4 := bls.GenerateTestKeyPair()

	sigs := []*bls.Signature{
		priv2.Sign(blockHash.RawBytes()).(*bls.Signature),
		priv3.Sign(blockHash.RawBytes()).(*bls.Signature),
		priv4.Sign(blockHash.RawBytes()).(*bls.Signature),
	}
	sig := bls.Aggregate(sigs)

	return NewCertificate(
		blockHash,
		util.RandInt(10),
		[]int{10, 18, 12, 16},
		[]int{18},
		sig)
}
