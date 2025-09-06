package certificate_test

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/types/certificate"
	"github.com/pactus-project/pactus/types/validator"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
)

// TestVoteCertificate tests the general properties of the Vote certificate.
// The data for the Vote certificate is the same as the block certificate in the slow path,
// but the sign bytes are different.
func TestVoteCertificate(t *testing.T) {
	expectedData, _ := hex.DecodeString(
		"04030201" + // Height
			"0100" + // Round
			"06010203040506" + // Committers
			"0102" + // Absentees
			"b53d79e156e9417e010fa21f2b2a96bee6be46fcd233295d2f697cdb9e782b6112ac01c80d0d9d64c2320664c77fa2a6") // Signature

	certHash, _ := hash.FromString("ac755295a6850b141286bde42bb8ba06ae1671f0562cbef90043924091177815")
	r := bytes.NewReader(expectedData)
	cert := new(certificate.VoteCertificate)
	err := cert.Decode(r)
	assert.NoError(t, err)
	assert.Equal(t, uint32(0x01020304), cert.Height())
	assert.Equal(t, int16(0x0001), cert.Round())
	assert.Equal(t, []int32{1, 2, 3, 4, 5, 6}, cert.Committers())
	assert.Equal(t, []int32{2}, cert.Absentees())
	assert.Equal(t, certHash, cert.Hash())
}

// Deprecated test.
func TestVoteCertificateValidatePrepare(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	blockHash := ts.RandHash()
	height := ts.RandHeight()
	round := ts.RandRound()
	cert := certificate.NewVoteCertificate(height, round)
	signBytes := cert.SignBytes(blockHash,
		util.StringToBytes("PREPARE"))
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

	t.Run("Doesn't have 2f+1 majority", func(t *testing.T) {
		absentees := committers[2:]
		aggSig := bls.SignatureAggregate(sigs[:2]...)
		cert.SetSignature(committers, absentees, aggSig)

		err := cert.ValidatePrepare(validators, blockHash)
		assert.ErrorIs(t, err, certificate.InsufficientPowerError{
			SignedPower:   2,
			RequiredPower: 3,
		})
	})

	t.Run("Ok, should return no error", func(t *testing.T) {
		absentees := committers[3:]
		aggSig := bls.SignatureAggregate(sigs[:3]...)
		cert.SetSignature(committers, absentees, aggSig)

		err := cert.ValidatePrepare(validators, blockHash)
		assert.NoError(t, err)
	})
}

func TestVoteCertificateValidatePrecommit(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	blockHash := ts.RandHash()
	height := ts.RandHeight()
	round := ts.RandRound()
	cert := certificate.NewVoteCertificate(height, round)
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

	t.Run("Doesn't have 2f+1 majority", func(t *testing.T) {
		absentees := committers[2:]
		aggSig := bls.SignatureAggregate(sigs[:2]...)
		cert.SetSignature(committers, absentees, aggSig)

		err := cert.ValidatePrecommit(validators, blockHash)
		assert.ErrorIs(t, err, certificate.InsufficientPowerError{
			SignedPower:   2,
			RequiredPower: 3,
		})
	})

	t.Run("Ok, should return no error", func(t *testing.T) {
		absentees := committers[3:]
		aggSig := bls.SignatureAggregate(sigs[:3]...)
		cert.SetSignature(committers, absentees, aggSig)

		err := cert.ValidatePrecommit(validators, blockHash)
		assert.NoError(t, err)
	})
}

func TestVoteCertificateValidateCPPreVote(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	blockHash := ts.RandHash()
	height := ts.RandHeight()
	round := ts.RandRound()
	cpRound := ts.RandRound()
	cpValue := 2
	cert := certificate.NewVoteCertificate(height, round)
	signBytes := cert.SignBytes(blockHash,
		util.StringToBytes("PRE-VOTE"),
		util.Int16ToSlice(cpRound),
		[]byte{byte(cpValue)})
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

	t.Run("Invalid cpValue", func(t *testing.T) {
		absentees := committers[3:]
		aggSig := bls.SignatureAggregate(sigs[:3]...)
		cert.SetSignature(committers, absentees, aggSig)

		err := cert.ValidateCPPreVote(validators, blockHash, cpRound, byte(0))
		assert.Error(t, err)
	})

	t.Run("Doesn't have 2f+1 majority", func(t *testing.T) {
		absentees := committers[2:]
		aggSig := bls.SignatureAggregate(sigs[:2]...)
		cert.SetSignature(committers, absentees, aggSig)

		err := cert.ValidateCPPreVote(validators, blockHash, cpRound, byte(cpValue))
		assert.ErrorIs(t, err, certificate.InsufficientPowerError{
			SignedPower:   2,
			RequiredPower: 3,
		})
	})

	t.Run("Ok, should return no error", func(t *testing.T) {
		absentees := committers[3:]
		aggSig := bls.SignatureAggregate(sigs[:3]...)
		cert.SetSignature(committers, absentees, aggSig)

		err := cert.ValidateCPPreVote(validators, blockHash, cpRound, byte(cpValue))
		assert.NoError(t, err)
	})
}

func TestVoteCertificateValidateCPMainVote(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	blockHash := ts.RandHash()
	height := ts.RandHeight()
	round := ts.RandRound()
	cpRound := ts.RandRound()
	cpValue := 2
	cert := certificate.NewVoteCertificate(height, round)
	signBytes := cert.SignBytes(blockHash,
		util.StringToBytes("MAIN-VOTE"),
		util.Int16ToSlice(cpRound),
		[]byte{byte(cpValue)})
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

	t.Run("Invalid cpValue", func(t *testing.T) {
		absentees := committers[3:]
		aggSig := bls.SignatureAggregate(sigs[:3]...)
		cert.SetSignature(committers, absentees, aggSig)

		err := cert.ValidateCPMainVote(validators, blockHash, cpRound, byte(0))
		assert.Error(t, err)
	})

	t.Run("Doesn't have 2f+1 majority", func(t *testing.T) {
		absentees := committers[2:]
		aggSig := bls.SignatureAggregate(sigs[:2]...)
		cert.SetSignature(committers, absentees, aggSig)

		err := cert.ValidateCPMainVote(validators, blockHash, cpRound, byte(cpValue))
		assert.ErrorIs(t, err, certificate.InsufficientPowerError{
			SignedPower:   2,
			RequiredPower: 3,
		})
	})

	t.Run("Ok, should return no error", func(t *testing.T) {
		absentees := committers[3:]
		aggSig := bls.SignatureAggregate(sigs[:3]...)
		cert.SetSignature(committers, absentees, aggSig)

		err := cert.ValidateCPMainVote(validators, blockHash, cpRound, byte(cpValue))
		assert.NoError(t, err)
	})
}
