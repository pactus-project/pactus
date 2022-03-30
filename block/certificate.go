package block

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/zarbchain/zarb-go/crypto/bls"
	"github.com/zarbchain/zarb-go/crypto/hash"
	"github.com/zarbchain/zarb-go/encoding"
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/util"
)

type Certificate struct {
	data certificateData
}
type certificateData struct {
	Round      int16
	Committers []int32
	Absentees  []int32
	Signature  *bls.Signature
}

func NewCertificate(round int16, committers, absentees []int32, signature *bls.Signature) *Certificate {
	cert := &Certificate{
		data: certificateData{
			Round:      round,
			Committers: committers,
			Absentees:  absentees,
			Signature:  signature,
		},
	}

	return cert
}

func (cert *Certificate) Round() int16              { return cert.data.Round }
func (cert *Certificate) Committers() []int32       { return cert.data.Committers }
func (cert *Certificate) Absentees() []int32        { return cert.data.Absentees }
func (cert *Certificate) Signature() *bls.Signature { return cert.data.Signature }

func (cert *Certificate) SanityCheck() error {
	if cert.Round() < 0 {
		return errors.Error(errors.ErrInvalidRound)
	}
	if cert.Signature() == nil {
		return errors.Errorf(errors.ErrInvalidSignature, "no signature")
	}
	if err := cert.Signature().SanityCheck(); err != nil {
		return err
	}
	if cert.Committers() == nil {
		return errors.Errorf(errors.ErrInvalidBlock, "invalid committers")
	}
	if cert.Absentees() == nil {
		return errors.Errorf(errors.ErrInvalidBlock, "invalid absentees")
	}
	signedBy := util.Subtracts(cert.Committers(), cert.Absentees())
	if !util.Equal(util.Subtracts(cert.Committers(), signedBy), cert.Absentees()) {
		return errors.Errorf(errors.ErrInvalidBlock, "absentees is not subset of committers")
	}

	return nil
}

func (cert *Certificate) Hash() hash.Hash {
	w := bytes.NewBuffer(make([]byte, 0, cert.SerializeSize()))
	if err := cert.Encode(w); err != nil {
		return hash.UndefHash
	}
	return hash.CalcHash(w.Bytes())
}

// SerializeSize returns the number of bytes it would take to serialize the block
func (cert *Certificate) SerializeSize() int {
	sz := encoding.VarIntSerializeSize(uint64(cert.Round())) +
		encoding.VarIntSerializeSize(uint64(len(cert.Committers()))) +
		encoding.VarIntSerializeSize(uint64(len(cert.Absentees()))) +
		bls.SignatureSize

	for _, n := range cert.Committers() {
		sz += encoding.VarIntSerializeSize(uint64(n))
	}

	for _, n := range cert.Absentees() {
		sz += encoding.VarIntSerializeSize(uint64(n))
	}
	return sz
}

// Encode encodes the receiver to w.
func (cert *Certificate) Encode(w io.Writer) error {
	if err := encoding.WriteVarInt(w, uint64(cert.Round())); err != nil {
		return err
	}
	if err := encoding.WriteVarInt(w, uint64(len(cert.data.Committers))); err != nil {
		return err
	}
	for _, n := range cert.data.Committers {
		if err := encoding.WriteVarInt(w, uint64(n)); err != nil {
			return err
		}
	}
	if err := encoding.WriteVarInt(w, uint64(len(cert.data.Absentees))); err != nil {
		return err
	}
	for _, n := range cert.data.Absentees {
		if err := encoding.WriteVarInt(w, uint64(n)); err != nil {
			return err
		}
	}
	if err := cert.data.Signature.Encode(w); err != nil {
		return err
	}

	return nil
}

func (cert *Certificate) Decode(r io.Reader) error {
	round, err := encoding.ReadVarInt(r)
	if err != nil {
		return err
	}

	lenCommitters, err := encoding.ReadVarInt(r)
	if err != nil {
		return err
	}
	committers := make([]int32, lenCommitters)
	for i := 0; i < int(lenCommitters); i++ {
		n, err := encoding.ReadVarInt(r)
		if err != nil {
			return err
		}
		committers[i] = int32(n)
	}

	lenAbsentees, err := encoding.ReadVarInt(r)
	if err != nil {
		return err
	}
	absentees := make([]int32, lenAbsentees)
	for i := 0; i < int(lenAbsentees); i++ {
		n, err := encoding.ReadVarInt(r)
		if err != nil {
			return err
		}
		absentees[i] = int32(n)
	}
	if err != nil {
		return err
	}

	sig := new(bls.Signature)
	if err := sig.Decode(r); err != nil {
		return err
	}

	cert.data.Round = int16(round)
	cert.data.Committers = committers
	cert.data.Absentees = absentees
	cert.data.Signature = sig

	return nil
}

func (cert *Certificate) MarshalJSON() ([]byte, error) {
	return json.Marshal(cert.data)
}

func CertificateSignBytes(blockHash hash.Hash, round int16) []byte {
	sb := blockHash.Bytes()
	sb = append(sb, util.Int16ToSlice(round)...)

	return sb
}

func GenerateTestCertificate(blockHash hash.Hash) *Certificate {
	_, priv2 := bls.GenerateTestKeyPair()
	_, priv3 := bls.GenerateTestKeyPair()
	_, priv4 := bls.GenerateTestKeyPair()

	sigs := []*bls.Signature{
		priv2.Sign(blockHash.Bytes()).(*bls.Signature),
		priv3.Sign(blockHash.Bytes()).(*bls.Signature),
		priv4.Sign(blockHash.Bytes()).(*bls.Signature),
	}
	sig := bls.Aggregate(sigs)

	c1 := util.RandInt32(10)
	c2 := util.RandInt32(10) + 10
	c3 := util.RandInt32(10) + 20
	c4 := util.RandInt32(10) + 30
	return NewCertificate(
		util.RandInt16(10),
		[]int32{c1, c2, c3, c4},
		[]int32{c2},
		sig)
}
