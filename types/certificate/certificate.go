package certificate

import (
	"bytes"
	"fmt"
	"io"
	"slices"

	"github.com/fxamacker/cbor/v2"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/types/validator"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/encoding"
)

// Certificate represents a base structure for both Certificate and Certificate.
// As a Certificate, it verifies if a block is signed by a majority of validators.
// As a Certificate, it checks whether a majority of validators have voted in the consensus step.
type Certificate struct {
	height     uint32
	round      int16
	committers []int32
	absentees  []int32
	signature  *bls.Signature
}

// NewCertificate creates a new Certificate instance.
func NewCertificate(height uint32, round int16) *Certificate {
	return &Certificate{
		height: height,
		round:  round,
	}
}

func (cert *Certificate) Height() uint32 {
	return cert.height
}

func (cert *Certificate) Round() int16 {
	return cert.round
}

func (cert *Certificate) Committers() []int32 {
	return cert.committers
}

func (cert *Certificate) Absentees() []int32 {
	return cert.absentees
}

func (cert *Certificate) Signature() *bls.Signature {
	return cert.signature
}

func (cert *Certificate) BasicCheck() error {
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

func (cert *Certificate) Hash() hash.Hash {
	buf := bytes.NewBuffer(make([]byte, 0, cert.SerializeSize()))
	if err := cert.Encode(buf); err != nil {
		return hash.UndefHash
	}

	return hash.CalcHash(buf.Bytes())
}

func (cert *Certificate) SetSignature(committers, absentees []int32, signature *bls.Signature) {
	cert.committers = committers
	cert.absentees = absentees
	cert.signature = signature
}

// SerializeSize returns the number of bytes it would take to serialize the block.
func (cert *Certificate) SerializeSize() int {
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

func (cert *Certificate) MarshalCBOR() ([]byte, error) {
	buf := bytes.NewBuffer(make([]byte, 0, cert.SerializeSize()))
	if err := cert.Encode(buf); err != nil {
		return nil, err
	}

	return cbor.Marshal(buf.Bytes())
}

func (cert *Certificate) UnmarshalCBOR(bs []byte) error {
	data := make([]byte, 0, cert.SerializeSize())
	err := cbor.Unmarshal(bs, &data)
	if err != nil {
		return err
	}
	buf := bytes.NewBuffer(data)

	return cert.Decode(buf)
}

func (cert *Certificate) Encode(w io.Writer) error {
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

func (cert *Certificate) Decode(r io.Reader) error {
	err := encoding.ReadElements(r, &cert.height, &cert.round)
	if err != nil {
		return err
	}

	lenCommitters, err := encoding.ReadVarInt(r)
	if err != nil {
		return err
	}

	// We set 51 as a hardcoded value,
	// matching the maximum number of committers allowed in a block.
	if lenCommitters > 51 {
		return ErrTooManyCommitters
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
	// We set 51 as a hardcoded value,
	// matching the maximum number of committers allowed in a block.
	if lenAbsentees > 51 {
		return ErrTooManyAbsentees
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

func (cert *Certificate) SignBytesPrepare(blockHash hash.Hash) []byte {
	return cert.signBytes(blockHash,
		util.StringToBytes("PREPARE"))
}

func (cert *Certificate) SignBytesPrecommit(blockHash hash.Hash) []byte {
	return cert.signBytes(blockHash)
}

func (cert *Certificate) SignBytesCPPreVote(blockHash hash.Hash, cpRound int16, cpValue byte) []byte {
	return cert.signBytes(blockHash,
		util.StringToBytes("PRE-VOTE"),
		util.Int16ToSlice(cpRound),
		[]byte{cpValue})
}

func (cert *Certificate) SignBytesCPMainVote(blockHash hash.Hash, cpRound int16, cpValue byte) []byte {
	return cert.signBytes(blockHash,
		util.StringToBytes("MAIN-VOTE"),
		util.Int16ToSlice(cpRound),
		[]byte{cpValue})
}

func (cert *Certificate) SignBytesCPDecided(blockHash hash.Hash, cpRound int16, cpValue byte) []byte {
	return cert.signBytes(blockHash,
		util.StringToBytes("DECIDED"),
		util.Int16ToSlice(cpRound),
		[]byte{cpValue})
}

// signBytes returns the sign bytes for the vote certificate.
func (cert *Certificate) signBytes(blockHash hash.Hash, extraData ...[]byte) []byte {
	signBytes := blockHash.Bytes()
	signBytes = append(signBytes, util.Uint32ToSlice(cert.height)...)
	signBytes = append(signBytes, util.Int16ToSlice(cert.round)...)
	for _, data := range extraData {
		signBytes = append(signBytes, data...)
	}

	return signBytes
}

func (cert *Certificate) ValidatePrepare(validators []*validator.Validator,
	blockHash hash.Hash,
) error {
	signBytes := cert.SignBytesPrepare(blockHash)

	return cert.validate(validators, signBytes, Required2FP1Power)
}

func (cert *Certificate) ValidatePrecommit(validators []*validator.Validator,
	blockHash hash.Hash,
) error {
	signBytes := cert.SignBytesPrecommit(blockHash)

	return cert.validate(validators, signBytes, Required2FP1Power)
}

func (cert *Certificate) ValidateCPPreVote(validators []*validator.Validator,
	blockHash hash.Hash, cpRound int16, cpValue byte,
) error {
	signBytes := cert.SignBytesCPPreVote(blockHash, cpRound, cpValue)

	return cert.validate(validators, signBytes, Required2FP1Power)
}

func (cert *Certificate) ValidateCPMainVote(validators []*validator.Validator,
	blockHash hash.Hash, cpRound int16, cpValue byte,
) error {
	signBytes := cert.SignBytesCPMainVote(blockHash, cpRound, cpValue)

	return cert.validate(validators, signBytes, Required2FP1Power)
}

func (cert *Certificate) validate(validators []*validator.Validator,
	signBytes []byte, requiredPowerFn RequiredPowerFn,
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

		if !slices.Contains(cert.absentees, num) {
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

	aggPub, _ := bls.PublicKeyAggregate(pubs...)

	return aggPub.Verify(signBytes, cert.signature)
}

// AddSignature adds a new signature to the certificate.
// It does not check the validity of the signature.
// The caller should ensure that the signature is valid.
func (cert *Certificate) AddSignature(valNum int32, sig *bls.Signature) {
	absentees, removed := util.RemoveFirstOccurrenceOf(cert.absentees, valNum)
	if removed {
		aggSig, _ := bls.SignatureAggregate(cert.signature, sig)
		cert.signature = aggSig
		cert.absentees = absentees
	}
}

func (cert *Certificate) Clone() *Certificate {
	cloned := &Certificate{
		height:     cert.height,
		round:      cert.round,
		committers: make([]int32, len(cert.committers)),
		absentees:  make([]int32, len(cert.absentees)),
		signature:  new(bls.Signature),
	}

	copy(cloned.committers, cert.committers)
	copy(cloned.absentees, cert.absentees)
	*cloned.signature = *cert.signature

	return cloned
}
