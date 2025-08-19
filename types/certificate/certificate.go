package certificate

import (
	"bytes"
	"fmt"
	"io"

	"github.com/fxamacker/cbor/v2"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/types/validator"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/encoding"
)

// baseCertificate represents a base structure for both BlockCertificate and VoteCertificate.
// As a BlockCertificate, it verifies if a block is signed by a majority of validators.
// As a VoteCertificate, it checks whether a majority of validators have voted in the consensus step.
type baseCertificate struct {
	height     uint32
	round      int16
	committers []int32
	absentees  []int32
	signature  *bls.Signature
}

func (cert *baseCertificate) Height() uint32 {
	return cert.height
}

func (cert *baseCertificate) Round() int16 {
	return cert.round
}

func (cert *baseCertificate) Committers() []int32 {
	return cert.committers
}

func (cert *baseCertificate) Absentees() []int32 {
	return cert.absentees
}

func (cert *baseCertificate) Signature() *bls.Signature {
	return cert.signature
}

func (cert *baseCertificate) BasicCheck() error {
	if cert.height <= 0 {
		return BasicCheckError{
			Reason: fmt.Sprintf("height is not positive: %d", cert.height),
		}
	}
	if cert.round < 0 {
		return BasicCheckError{
			Reason: fmt.Sprintf("round is negative: %d", cert.round),
		}
	}
	if cert.signature == nil {
		return BasicCheckError{
			Reason: "signature is missing",
		}
	}
	if cert.committers == nil {
		return BasicCheckError{
			Reason: "committers is missing",
		}
	}
	if cert.absentees == nil {
		return BasicCheckError{
			Reason: "absentees is missing",
		}
	}
	if !util.IsSubset(cert.committers, cert.absentees) {
		return BasicCheckError{
			Reason: fmt.Sprintf("absentees are not a subset of committers: %v, %v",
				cert.committers, cert.absentees),
		}
	}

	return nil
}

func (cert *baseCertificate) Hash() hash.Hash {
	buf := bytes.NewBuffer(make([]byte, 0, cert.SerializeSize()))
	if err := cert.Encode(buf); err != nil {
		return hash.UndefHash
	}

	return hash.CalcHash(buf.Bytes())
}

func (cert *baseCertificate) SetSignature(committers, absentees []int32, signature *bls.Signature) {
	cert.committers = committers
	cert.absentees = absentees
	cert.signature = signature
}

// SerializeSize returns the number of bytes it would take to serialize the block.
func (cert *baseCertificate) SerializeSize() int {
	size := 6 + // height (4) + round(2)
		encoding.VarIntSerializeSize(uint64(len(cert.committers))) +
		encoding.VarIntSerializeSize(uint64(len(cert.absentees))) +
		bls.SignatureSize

	for _, n := range cert.committers {
		size += encoding.VarIntSerializeSize(uint64(n))
	}

	for _, n := range cert.absentees {
		size += encoding.VarIntSerializeSize(uint64(n))
	}

	return size
}

func (cert *baseCertificate) MarshalCBOR() ([]byte, error) {
	buf := bytes.NewBuffer(make([]byte, 0, cert.SerializeSize()))
	if err := cert.Encode(buf); err != nil {
		return nil, err
	}

	return cbor.Marshal(buf.Bytes())
}

func (cert *baseCertificate) UnmarshalCBOR(bs []byte) error {
	data := make([]byte, 0, cert.SerializeSize())
	err := cbor.Unmarshal(bs, &data)
	if err != nil {
		return err
	}
	buf := bytes.NewBuffer(data)

	return cert.Decode(buf)
}

func (cert *baseCertificate) Encode(w io.Writer) error {
	if err := encoding.WriteElements(w, cert.height, cert.round); err != nil {
		return err
	}
	if err := encoding.WriteVarInt(w, uint64(len(cert.committers))); err != nil {
		return err
	}
	for _, n := range cert.committers {
		if err := encoding.WriteVarInt(w, uint64(n)); err != nil {
			return err
		}
	}
	if err := encoding.WriteVarInt(w, uint64(len(cert.absentees))); err != nil {
		return err
	}
	for _, n := range cert.absentees {
		if err := encoding.WriteVarInt(w, uint64(n)); err != nil {
			return err
		}
	}

	return cert.signature.Encode(w)
}

func (cert *baseCertificate) Decode(r io.Reader) error {
	err := encoding.ReadElements(r, &cert.height, &cert.round)
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

	sig := new(bls.Signature)
	if err := sig.Decode(r); err != nil {
		return err
	}

	cert.committers = committers
	cert.absentees = absentees
	cert.signature = sig

	return nil
}

type requiredPowerFn func(int64) int64

var require2FPower = func(committeePower int64) int64 {
	f := (committeePower - 1) / 3
	p := (2 * f) + 1

	return p
}

func (cert *baseCertificate) validate(validators []*validator.Validator,
	signBytes []byte, requiredPowerFn requiredPowerFn,
) error {
	if len(validators) != len(cert.committers) {
		return UnexpectedCommittersError{
			Committers: cert.committers,
		}
	}

	pubs := make([]*bls.PublicKey, 0, len(cert.committers))
	committeePower := int64(0)
	signedPower := int64(0)

	for index, num := range cert.committers {
		val := validators[index]
		if val.Number() != num {
			return UnexpectedCommittersError{
				Committers: cert.committers,
			}
		}

		if !util.Contains(cert.absentees, num) {
			pubs = append(pubs, val.PublicKey())
			signedPower += val.Power()
		}
		committeePower += val.Power()
	}

	requiredPower := requiredPowerFn(committeePower)

	// Check if signers have enough power
	if signedPower < requiredPower {
		return InsufficientPowerError{
			SignedPower:   signedPower,
			RequiredPower: requiredPower,
		}
	}

	aggPub := bls.PublicKeyAggregate(pubs...)

	return aggPub.Verify(signBytes, cert.signature)
}

// AddSignature adds a new signature to the certificate.
// It does not check the validity of the signature.
// The caller should ensure that the signature is valid.
func (cert *baseCertificate) AddSignature(valNum int32, sig *bls.Signature) {
	absentees, removed := util.RemoveFirstOccurrenceOf(cert.absentees, valNum)
	if removed {
		cert.signature = bls.SignatureAggregate(cert.signature, sig)
		cert.absentees = absentees
	}
}
