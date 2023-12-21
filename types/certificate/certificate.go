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

// TODO: proofread this comment later
// A certificate is a multipurpose object.
// It can be used either as a block certificate or as a quorum vote certificate within consensus.
// As a block certificate, it verifies if a block is signed by a majority of validators.
// As a quorum vote certificate, it checks whether a majority of validators have voted in the consensus step.
type Certificate struct {
	data certificateData
}

type certificateData struct {
	Height     uint32
	Round      int16
	FastPath   bool
	Committers []int32
	Absentees  []int32
	Signature  *bls.Signature
}

func NewCertificate(height uint32, round int16,
	committers, absentees []int32, signature *bls.Signature,
) *Certificate {
	cert := &Certificate{
		data: certificateData{
			Height:     height,
			Round:      round,
			FastPath:   false,
			Committers: committers,
			Absentees:  absentees,
			Signature:  signature,
		},
	}

	return cert
}

func NewCertificateFastPath(height uint32, round int16,
	committers, absentees []int32, signature *bls.Signature,
) *Certificate {
	cert := &Certificate{
		data: certificateData{
			Height:     height,
			Round:      round,
			FastPath:   true,
			Committers: committers,
			Absentees:  absentees,
			Signature:  signature,
		},
	}

	return cert
}

func (cert *Certificate) Height() uint32 {
	return cert.data.Height
}

func (cert *Certificate) Round() int16 {
	return cert.data.Round
}

func (cert *Certificate) FastPath() bool {
	return cert.data.FastPath
}

func (cert *Certificate) Committers() []int32 {
	return cert.data.Committers
}

func (cert *Certificate) Absentees() []int32 {
	return cert.data.Absentees
}

func (cert *Certificate) Signature() *bls.Signature {
	return cert.data.Signature
}

func (cert *Certificate) BasicCheck() error {
	if cert.Height() <= 0 {
		return BasicCheckError{
			Reason: fmt.Sprintf("height is not positive: %d", cert.Height()),
		}
	}
	if cert.Round() < 0 {
		return BasicCheckError{
			Reason: fmt.Sprintf("round is negative: %d", cert.Round()),
		}
	}
	if cert.Signature() == nil {
		return BasicCheckError{
			Reason: "signature is missing",
		}
	}
	if cert.Committers() == nil {
		return BasicCheckError{
			Reason: "committers is missing",
		}
	}
	if cert.Absentees() == nil {
		return BasicCheckError{
			Reason: "absentees is missing",
		}
	}
	if !util.IsSubset(cert.Committers(), cert.Absentees()) {
		return BasicCheckError{
			Reason: fmt.Sprintf("absentees are not a subset of committers: %v, %v",
				cert.Committers(), cert.Absentees()),
		}
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

func (cert *Certificate) Clone() *Certificate {
	cloned := &Certificate{
		data: certificateData{
			Height:     cert.Height(),
			Round:      cert.Round(),
			Committers: make([]int32, len(cert.data.Committers)),
			Absentees:  make([]int32, len(cert.data.Absentees)),
			Signature:  new(bls.Signature),
		},
	}

	copy(cloned.data.Committers, cert.data.Committers)
	copy(cloned.data.Absentees, cert.data.Absentees)
	*cloned.data.Signature = *cert.data.Signature

	return cloned
}

// SerializeSize returns the number of bytes it would take to serialize the block.
func (cert *Certificate) SerializeSize() int {
	sz := 6 + // height (4) + round(2)
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
	roundAndPath := uint16(cert.data.Round)
	if cert.FastPath() {
		roundAndPath |= 0x8000
	}

	if err := encoding.WriteElements(w, cert.data.Height, roundAndPath); err != nil {
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

	return cert.data.Signature.Encode(w)
}

func (cert *Certificate) Decode(r io.Reader) error {
	roundAndPath := uint16(0)
	err := encoding.ReadElements(r, &cert.data.Height, &roundAndPath)
	if err != nil {
		return err
	}

	cert.data.Round = int16(roundAndPath & 0x7FFF)
	if roundAndPath&0x8000 == 0x8000 {
		cert.data.FastPath = true
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

	cert.data.Committers = committers
	cert.data.Absentees = absentees
	cert.data.Signature = sig

	return nil
}

func (cert *Certificate) Validate(blockHash hash.Hash, height uint32,
	validators []*validator.Validator, extraSignBytes []byte,
) error {
	if cert.Height() != height {
		return UnexpectedHeightError{
			Expected: height,
			Got:      cert.Height(),
		}
	}

	if len(validators) != len(cert.Committers()) {
		return UnexpectedCommittersError{
			Committers: cert.Committers(),
		}
	}

	pubs := make([]*bls.PublicKey, 0, len(cert.Committers()))
	committeePower := int64(0)
	signedPower := int64(0)

	for index, num := range cert.Committers() {
		val := validators[index]
		if val.Number() != num {
			return UnexpectedCommittersError{
				Committers: cert.Committers(),
			}
		}

		if !util.Contains(cert.Absentees(), num) {
			pubs = append(pubs, val.PublicKey())
			signedPower += val.Power()
		}
		committeePower += val.Power()
	}

	t := (committeePower - 1) / 3 // TODO => t := (committeePower - 1) / 5
	p := int64(0)
	if cert.data.FastPath {
		p = 3*t + 1 // TODO ???? => 4t+1
	} else {
		p = 2*t + 1 // TODO ???? => 3t+1
	}

	// Check if signers have enough power
	if signedPower < p {
		return InsufficientPowerError{
			SignedPower:   signedPower,
			RequiredPower: p,
		}
	}

	// Check signature
	var signBytes []byte
	if cert.FastPath() {
		signBytes = BlockCertificateSignBytesFastPath(blockHash, cert.Height(), cert.Round())
	} else {
		signBytes = BlockCertificateSignBytes(blockHash, cert.Height(), cert.Round())
	}
	signBytes = append(signBytes, extraSignBytes...)
	return bls.VerifyAggregated(cert.Signature(), pubs, signBytes)
}

// AddSignature adds a new signature to the certificate.
// It does not check the validity of the signature.
// The caller should ensure that the signature is valid.
func (cert *Certificate) AddSignature(valNum int32, sig *bls.Signature) {
	absentees, removed := util.RemoveFirstOccurrenceOf(cert.data.Absentees, valNum)
	if removed {
		cert.data.Signature = bls.SignatureAggregate(cert.data.Signature, sig)
		cert.data.Absentees = absentees
	}
}

func BlockCertificateSignBytes(blockHash hash.Hash, height uint32, round int16) []byte {
	sb := blockHash.Bytes()
	sb = append(sb, util.Uint32ToSlice(height)...)
	sb = append(sb, util.Int16ToSlice(round)...)

	return sb
}

func BlockCertificateSignBytesFastPath(blockHash hash.Hash, height uint32, round int16) []byte {
	sb := BlockCertificateSignBytes(blockHash, height, round)
	sb = append(sb, util.StringToBytes("PREPARE")...)

	return sb
}
