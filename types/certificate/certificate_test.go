package certificate_test

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/fxamacker/cbor/v2"
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/types/certificate"
	"github.com/pactus-project/pactus/types/validator"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/slices"
)

func TestCertificate(t *testing.T) {
	data, _ := hex.DecodeString(
		"04030201" + // Height
			"0100" + // Round
			"06010203040506" + // Committers
			"0102" + // Absentees
			"b53d79e156e9417e010fa21f2b2a96bee6be46fcd233295d2f697cdb9e782b6112ac01c80d0d9d64c2320664c77fa2a6") // Signature

	certHash, _ := hash.FromString("ac755295a6850b141286bde42bb8ba06ae1671f0562cbef90043924091177815")
	r := bytes.NewReader(data)
	cert := new(certificate.Certificate)
	err := cert.Decode(r)
	assert.NoError(t, err)
	assert.Equal(t, uint32(0x01020304), cert.Height())
	assert.Equal(t, int16(0x0001), cert.Round())
	assert.Equal(t, []int32{1, 2, 3, 4, 5, 6}, cert.Committers())
	assert.Equal(t, []int32{2}, cert.Absentees())
	assert.Equal(t, certHash, cert.Hash())

	blockHash, _ := hash.FromString("000102030405060708090a0b0c0d0e0f000102030405060708090a0b0c0d0e0f")
	expectedPrepareSignByte, _ := hex.DecodeString(
		"000102030405060708090a0b0c0d0e0f000102030405060708090a0b0c0d0e0f" + // Block hash
			"04030201" + // Height
			"0100" + // Round
			"50524550415245") // PREPARE
	expectedPrecommitSignByte, _ := hex.DecodeString(
		"000102030405060708090a0b0c0d0e0f000102030405060708090a0b0c0d0e0f" + // Block hash
			"04030201" + // Height
			"0100") // Round
	expectedCPPreVoteSignByte, _ := hex.DecodeString(
		"000102030405060708090a0b0c0d0e0f000102030405060708090a0b0c0d0e0f" + // Block hash
			"04030201" + // Height
			"0100" + // Round
			"5052452d564f5445" + // PRE-VOTE
			"0100" + // CP Round
			"02") // CP Value

	expectedCPMainVoteSignByte, _ := hex.DecodeString(
		"000102030405060708090a0b0c0d0e0f000102030405060708090a0b0c0d0e0f" + // Block hash
			"04030201" + // Height
			"0100" + // Round
			"4d41494e2d564f5445" + // MAIN-VOTE
			"0100" + // CP Round
			"02") // CP Value

	expectedCPDecidedSignByte, _ := hex.DecodeString(
		"000102030405060708090a0b0c0d0e0f000102030405060708090a0b0c0d0e0f" + // Block hash
			"04030201" + // Height
			"0100" + // Round
			"44454349444544" + // DECIDED
			"0100" + // CP Round
			"02") // CP Value

	assert.Equal(t, expectedPrepareSignByte, cert.SignBytesPrepare(blockHash))
	assert.Equal(t, expectedPrecommitSignByte, cert.SignBytesPrecommit(blockHash))
	assert.Equal(t, expectedCPPreVoteSignByte, cert.SignBytesCPPreVote(blockHash, 1, 2))
	assert.Equal(t, expectedCPMainVoteSignByte, cert.SignBytesCPMainVote(blockHash, 1, 2))
	assert.Equal(t, expectedCPDecidedSignByte, cert.SignBytesCPDecided(blockHash, 1, 2))
}

func TestCertificateCBORMarshaling(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	cert1 := ts.GenerateTestCertificate(ts.RandHeight())
	bz1, err := cbor.Marshal(cert1)
	assert.NoError(t, err)
	var cert2 certificate.Certificate
	err = cbor.Unmarshal(bz1, &cert2)
	assert.NoError(t, err)
	assert.NoError(t, cert2.BasicCheck())
	assert.Equal(t, cert1.Hash(), cert1.Hash())

	assert.True(t, cert1.Signature().EqualsTo(cert2.Signature()))
	assert.Equal(t, cert1.Hash(), cert2.Hash())

	err = cbor.Unmarshal([]byte{1}, &cert2)
	assert.Error(t, err)
}

func TestInvalidCertificate(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	t.Run("Invalid height", func(t *testing.T) {
		cert := certificate.NewCertificate(0, 0)

		err := cert.BasicCheck()
		assert.ErrorIs(t, err, certificate.BasicCheckError{
			Reason: "height is not positive: 0",
		})
	})

	t.Run("Invalid round", func(t *testing.T) {
		cert := certificate.NewCertificate(1, -1)

		err := cert.BasicCheck()
		assert.ErrorIs(t, err, certificate.BasicCheckError{
			Reason: "round is negative: -1",
		})
	})

	t.Run("Committers is nil", func(t *testing.T) {
		cert := certificate.NewCertificate(ts.RandHeight(), ts.RandRound())
		cert.SetSignature(nil, []int32{1}, ts.RandBLSSignature())

		err := cert.BasicCheck()
		assert.ErrorIs(t, err, certificate.BasicCheckError{
			Reason: "committers is missing",
		})
	})

	t.Run("Absentees is nil", func(t *testing.T) {
		cert := certificate.NewCertificate(ts.RandHeight(), ts.RandRound())
		cert.SetSignature([]int32{1, 2, 3, 4}, nil, ts.RandBLSSignature())

		err := cert.BasicCheck()
		assert.ErrorIs(t, err, certificate.BasicCheckError{
			Reason: "absentees is missing",
		})
	})

	t.Run("Signature is nil", func(t *testing.T) {
		cert := certificate.NewCertificate(ts.RandHeight(), ts.RandRound())
		cert.SetSignature([]int32{1, 2, 3, 4, 5, 6}, []int32{1}, nil)

		err := cert.BasicCheck()
		assert.ErrorIs(t, err, certificate.BasicCheckError{
			Reason: "signature is missing",
		})
	})

	t.Run("Invalid Absentees ", func(t *testing.T) {
		cert := certificate.NewCertificate(ts.RandHeight(), ts.RandRound())
		cert.SetSignature([]int32{11, 2, 3, 4, 5, 6}, []int32{66}, ts.RandBLSSignature())

		err := cert.BasicCheck()
		assert.ErrorIs(t, err, certificate.BasicCheckError{
			Reason: fmt.Sprintf("absentees are not a subset of committers: %v, %v",
				cert.Committers(), []int32{66}),
		})
	})

	t.Run("Invalid Absentees ", func(t *testing.T) {
		cert := certificate.NewCertificate(ts.RandHeight(), ts.RandRound())
		cert.SetSignature([]int32{1, 2, 3, 4, 5, 6}, []int32{2, 1}, ts.RandBLSSignature())

		err := cert.BasicCheck()
		assert.ErrorIs(t, err, certificate.BasicCheckError{
			Reason: fmt.Sprintf("absentees are not a subset of committers: %v, %v",
				cert.Committers(), []int32{2, 1}),
		})
	})
}

func TestEncodingCertificate(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	cert1 := ts.GenerateTestCertificate(ts.RandHeight())
	length := cert1.SerializeSize()

	for i := 0; i < length; i++ {
		w := util.NewFixedWriter(i)
		assert.Error(t, cert1.Encode(w), "encode test %v failed", i)
	}
	writer := util.NewFixedWriter(length)
	assert.NoError(t, cert1.Encode(writer))

	for i := 0; i < length; i++ {
		cert := new(certificate.Certificate)
		r := util.NewFixedReader(i, writer.Bytes())
		assert.Error(t, cert.Decode(r), "decode test %v failed", i)
	}

	cert2 := new(certificate.Certificate)
	reader := util.NewFixedReader(length, writer.Bytes())
	assert.NoError(t, cert2.Decode(reader))
	assert.Equal(t, cert1.Hash(), cert2.Hash())
}

func TestAddSignature(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	blockHash := ts.RandHash()
	height := ts.RandHeight()
	round := ts.RandRound()
	cert := certificate.NewCertificate(height, round)
	signBytes := cert.SignBytesPrecommit(blockHash)
	committers := ts.RandSlice(4)
	sigs := []*bls.Signature{}
	validators := []*validator.Validator{}

	for _, committer := range committers {
		valKey := ts.RandValKey()
		val := validator.NewValidator(valKey.PublicKey(), committer)
		sig := valKey.Sign(signBytes)

		validators = append(validators, val)
		sigs = append(sigs, sig)
	}

	absentees := committers[3:]
	aggSig := bls.SignatureAggregate(sigs[:3]...)
	cert.SetSignature(committers, absentees, aggSig)

	err := cert.ValidatePrecommit(validators, blockHash)
	assert.NoError(t, err)

	numAbsentees := len(cert.Absentees())

	t.Run("Add an existing signature", func(t *testing.T) {
		cert.AddSignature(validators[0].Number(), sigs[0])

		err := cert.ValidatePrecommit(validators, blockHash)
		assert.NoError(t, err)

		assert.Len(t, cert.Absentees(), numAbsentees)
	})

	t.Run("Add non existing signature", func(t *testing.T) {
		cert.AddSignature(validators[3].Number(), sigs[3])

		err := cert.ValidatePrecommit(validators, blockHash)
		assert.NoError(t, err)

		assert.Len(t, cert.Absentees(), numAbsentees-1)
	})
}

// Deprecated test.
func TestCertificateValidatePrepare(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	blockHash := ts.RandHash()
	height := ts.RandHeight()
	round := ts.RandRound()
	cert := certificate.NewCertificate(height, round)
	signBytes := cert.SignBytesPrepare(blockHash)
	committers := ts.RandSlice(4)
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

func TestCertificateValidatePrecommit(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	blockHash := ts.RandHash()
	height := ts.RandHeight()
	round := ts.RandRound()
	cert := certificate.NewCertificate(height, round)
	signBytes := cert.SignBytesPrecommit(blockHash)
	committers := ts.RandSlice(4)
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

		err := cert.ValidatePrecommit(validators, blockHash)
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
		err := cert.ValidatePrecommit(invValidators, blockHash)
		assert.ErrorIs(t, err, certificate.UnexpectedCommittersError{
			Committers: committers,
		})
	})

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

	t.Run("One signature short, should return an error for invalid signature", func(t *testing.T) {
		absentees := committers[4:]
		aggSig := bls.SignatureAggregate(sigs[3:]...)
		cert.SetSignature(committers, absentees, aggSig)

		err := cert.ValidatePrecommit(validators, blockHash)
		assert.ErrorIs(t, err, crypto.ErrInvalidSignature)
	})

	t.Run("Invalid block hash", func(t *testing.T) {
		absentees := committers[3:]
		aggSig := bls.SignatureAggregate(sigs[:3]...)
		cert.SetSignature(committers, absentees, aggSig)

		err := cert.ValidatePrecommit(validators, ts.RandHash())
		assert.ErrorIs(t, err, crypto.ErrInvalidSignature)
	})

	t.Run("Ok, should return no error", func(t *testing.T) {
		absentees := committers[3:]
		aggSig := bls.SignatureAggregate(sigs[:3]...)
		cert.SetSignature(committers, absentees, aggSig)

		err := cert.ValidatePrecommit(validators, blockHash)
		assert.NoError(t, err)
	})
}

func TestCertificateValidateCPPreVote(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	blockHash := ts.RandHash()
	height := ts.RandHeight()
	round := ts.RandRound()
	cpRound := ts.RandRound()
	cpValue := byte(2)
	cert := certificate.NewCertificate(height, round)
	signBytes := cert.SignBytesCPPreVote(blockHash, cpRound, cpValue)
	committers := ts.RandSlice(4)
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
		assert.ErrorIs(t, err, crypto.ErrInvalidSignature)
	})

	t.Run("Doesn't have 2f+1 majority", func(t *testing.T) {
		absentees := committers[2:]
		aggSig := bls.SignatureAggregate(sigs[:2]...)
		cert.SetSignature(committers, absentees, aggSig)

		err := cert.ValidateCPPreVote(validators, blockHash, cpRound, cpValue)
		assert.ErrorIs(t, err, certificate.InsufficientPowerError{
			SignedPower:   2,
			RequiredPower: 3,
		})
	})

	t.Run("Ok, should return no error", func(t *testing.T) {
		absentees := committers[3:]
		aggSig := bls.SignatureAggregate(sigs[:3]...)
		cert.SetSignature(committers, absentees, aggSig)

		err := cert.ValidateCPPreVote(validators, blockHash, cpRound, cpValue)
		assert.NoError(t, err)
	})
}

func TestCertificateValidateCPMainVote(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	blockHash := ts.RandHash()
	height := ts.RandHeight()
	round := ts.RandRound()
	cpRound := ts.RandRound()
	cpValue := byte(2)
	cert := certificate.NewCertificate(height, round)
	signBytes := cert.SignBytesCPMainVote(blockHash, cpRound, cpValue)
	committers := ts.RandSlice(4)
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
		assert.ErrorIs(t, err, crypto.ErrInvalidSignature)
	})

	t.Run("Doesn't have 2f+1 majority", func(t *testing.T) {
		absentees := committers[2:]
		aggSig := bls.SignatureAggregate(sigs[:2]...)
		cert.SetSignature(committers, absentees, aggSig)

		err := cert.ValidateCPMainVote(validators, blockHash, cpRound, cpValue)
		assert.ErrorIs(t, err, certificate.InsufficientPowerError{
			SignedPower:   2,
			RequiredPower: 3,
		})
	})

	t.Run("Ok, should return no error", func(t *testing.T) {
		absentees := committers[3:]
		aggSig := bls.SignatureAggregate(sigs[:3]...)
		cert.SetSignature(committers, absentees, aggSig)

		err := cert.ValidateCPMainVote(validators, blockHash, cpRound, cpValue)
		assert.NoError(t, err)
	})
}

func TestClone(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	cert1 := ts.GenerateTestCertificate(ts.RandHeight())
	cert2 := cert1.Clone()
	cert2.AddSignature(cert2.Absentees()[0], ts.RandBLSSignature())
	assert.NotEqual(t, cert1.Absentees(), cert2.Absentees())
	assert.NotEqual(t, cert1, cert2)
}
