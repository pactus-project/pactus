package certificate_test

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/fxamacker/cbor/v2"
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/types/certificate"
	"github.com/pactus-project/pactus/types/validator"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/slices"
)

func TestBlockCertificate(t *testing.T) {
	expectedData, _ := hex.DecodeString(
		"04030201" + // Height
			"0100" + // Round
			"06010203040506" + // Committers
			"0102" + // Absentees
			"b53d79e156e9417e010fa21f2b2a96bee6be46fcd233295d2f697cdb9e782b6112ac01c80d0d9d64c2320664c77fa2a6") // Signature

	certHash, _ := hash.FromString("ac755295a6850b141286bde42bb8ba06ae1671f0562cbef90043924091177815")
	r := bytes.NewReader(expectedData)
	cert := new(certificate.BlockCertificate)
	err := cert.Decode(r)
	assert.NoError(t, err)
	assert.Equal(t, uint32(0x01020304), cert.Height())
	assert.Equal(t, int16(0x0001), cert.Round())
	assert.Equal(t, []int32{1, 2, 3, 4, 5, 6}, cert.Committers())
	assert.Equal(t, []int32{2}, cert.Absentees())
	assert.Equal(t, certHash, cert.Hash())

	blockHash, _ := hash.FromString("000102030405060708090a0b0c0d0e0f000102030405060708090a0b0c0d0e0f")
	expectedSignByte, _ := hex.DecodeString(
		"000102030405060708090a0b0c0d0e0f000102030405060708090a0b0c0d0e0f" + // Block hash
			"04030201" + // Height
			"0100") // Round

	assert.Equal(t, expectedSignByte, cert.SignBytes(blockHash))
}

func TestBlockCertificateCBORMarshaling(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	cert1 := ts.GenerateTestBlockCertificate(ts.RandHeight())
	bz1, err := cbor.Marshal(cert1)
	assert.NoError(t, err)
	var cert2 certificate.BlockCertificate
	err = cbor.Unmarshal(bz1, &cert2)
	assert.NoError(t, err)
	assert.NoError(t, cert2.BasicCheck())
	assert.Equal(t, cert1.Hash(), cert1.Hash())

	assert.Equal(t, cert1.Hash(), cert2.Hash())

	err = cbor.Unmarshal([]byte{1}, &cert2)
	assert.Error(t, err)
}

func TestBlockCertificateValidation(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	blockHash := ts.RandHash()
	height := ts.RandHeight()
	round := ts.RandRound()
	cert := certificate.NewBlockCertificate(height, round)
	signBytes := cert.SignBytes(blockHash)
	committers := ts.RandSlice(6)
	sigs := []*bls.Signature{}
	validators := []*validator.Validator{}

	for _, committer := range committers {
		valKey := ts.RandValKey()
		val := validator.NewValidator(valKey.PublicKey(), committer)
		sig := valKey.Sign(signBytes)

		validators = append(validators, val)
		sigs = append(sigs, sig)
	}

	t.Run("Invalid committer, should return error", func(t *testing.T) {
		invCommitters := slices.Clone(committers)
		invCommitters = append(invCommitters, ts.Rand.Int31n(10000))
		absentees := committers[4:]
		aggSig := bls.SignatureAggregate(sigs[:4]...)
		cert.SetSignature(invCommitters, absentees, aggSig)

		err := cert.Validate(validators, blockHash)
		assert.ErrorIs(t, err, certificate.UnexpectedCommittersError{
			Committers: invCommitters,
		})
	})

	t.Run("Invalid validator", func(t *testing.T) {
		absentees := committers[4:]
		aggSig := bls.SignatureAggregate(sigs[:4]...)
		cert.SetSignature(committers, absentees, aggSig)

		invValidators := slices.Clone(validators)
		invValidators[0] = ts.GenerateTestValidator()
		err := cert.Validate(invValidators, blockHash)
		assert.ErrorIs(t, err, certificate.UnexpectedCommittersError{
			Committers: committers,
		})
	})

	t.Run("Doesn't have 2f+1 majority", func(t *testing.T) {
		absentees := committers[2:]
		aggSig := bls.SignatureAggregate(sigs[:2]...)
		cert.SetSignature(committers, absentees, aggSig)

		err := cert.Validate(validators, blockHash)
		assert.ErrorIs(t, err, certificate.InsufficientPowerError{
			SignedPower:   2,
			RequiredPower: 3,
		})
	})

	t.Run("One signature short, should return an error for invalid signature", func(t *testing.T) {
		absentees := committers[4:]
		aggSig := bls.SignatureAggregate(sigs[3:]...)
		cert.SetSignature(committers, absentees, aggSig)

		err := cert.Validate(validators, blockHash)
		assert.ErrorIs(t, err, crypto.ErrInvalidSignature)
	})

	t.Run("Ok, should return no error", func(t *testing.T) {
		absentees := committers[4:]
		aggSig := bls.SignatureAggregate(sigs[:4]...)
		cert.SetSignature(committers, absentees, aggSig)

		err := cert.Validate(validators, blockHash)
		assert.NoError(t, err)
	})
}
