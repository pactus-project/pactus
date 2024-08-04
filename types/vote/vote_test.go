package vote_test

import (
	"encoding/hex"
	"testing"

	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/types/vote"
	"github.com/pactus-project/pactus/util/errors"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
)

func TestVoteMarshaling(t *testing.T) {
	tests := []struct {
		data      string
		justType  string
		signBytes string
	}{
		{
			"A7" + // map(7)
				"0101" + // Type: 1 (prepare vote)
				"021832" + // Height: 50
				"0301" + // Round: 1
				"045820BBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB" + // Block Hash
				"055501AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA" + // Signer
				"06f6" + // CP_vote -> Null
				"07f6", // Signature -> Null
			"",
			"BBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB32000000010050524550415245",
		},
		{
			"A7" + // map(7)
				"0102" + // Type: 2 (precommit vote)
				"021832" + // Height: 50
				"0301" + // Round: 1
				"045820BBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB" + // Block Hash
				"055501AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA" + // Signer
				"06f6" + // CP_vote -> Null
				"07f6", // Signature -> Null
			"",
			"BBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB320000000100",
		},
		{
			"A7" + // map(7)
				"0103" + // Type: 3 (cp:pre-vote)
				"021832" + // Height: 50
				"0301" + // Round: 1
				"045820BBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB" + // Block Hash
				"055501AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA" + // Signer
				"06" + // CP_vote
				"A4" + // map(4)
				"0100" + // CP_Round: 0
				"0200" + // CP_Value: 0
				"0301" + // Just type: 1
				"045840" + // Just: JustTypeInitNo
				"A1" + // map(1)
				"01583C" + // Certificate (60 bytes)
				"32000000010004010203040094D25422904AC1D130AC981374AA4424F988" + // Certificate Data
				"61E99131078EFEFD62FC52CF072B0C08BB04E4E6496BA48DE4F3D3309AAB" +
				"07f6", // Signature -> Null
			"JustInitNo",
			"BBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB3200000001005052452d564f5445000000",
		},
		{
			"A7" + // map(7)
				"0103" + // Type: 3 (cp:pre-vote)
				"021832" + // Height: 50
				"0301" + // Round: 1
				"0458200000000000000000000000000000000000000000000000000000000000000000" + // Block Hash
				"055501AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA" + // Signer
				"06" + // CP_vote
				"A4" + // map(4)
				"0100" + // CP_Round: 0
				"0201" + // CP_Value: 1
				"0302" + // Just type: 2
				"0441" + // Just: JustTypeInitYes
				"A0" + // Empty Array
				"07f6", // Signature -> Null
			"JustInitYes",
			"00000000000000000000000000000000000000000000000000000000000000003200000001005052452d564f5445000001",
		},
		{
			"A7" + // map(7)
				"0103" + // Type: 3 (cp:pre-vote)
				"021832" + // Height: 50
				"0301" + // Round: 1
				"045820BBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB" + // Block Hash
				"055501AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA" + // Signer
				"06" + // CP_vote
				"A4" + // map(4)
				"0101" + // CP_Round: 1
				"0200" + // CP_Value: 0
				"0303" + // Just type: 3
				"045840" + // Just: JustPreVoteSoft
				"A1" + // map(1)
				"01583C" + // Certificate (60 bytes)
				"32000000010004010203040094D25422904AC1D130AC981374AA4424F988" + // Certificate Data
				"61E99131078EFEFD62FC52CF072B0C08BB04E4E6496BA48DE4F3D3309AAB" +
				"07f6", // Signature -> Null
			"JustPreVoteSoft",
			"BBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB3200000001005052452d564f5445010000",
		},
		{
			"A7" + // map(7)
				"0103" + // Type: 3 (cp:pre-vote)
				"021832" + // Height: 50
				"0301" + // Round: 1
				"045820BBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB" + // Block Hash
				"055501AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA" + // Signer
				"06" + // CP_vote
				"A4" + // map(4)
				"0101" + // CP_Round: 1
				"0200" + // CP_Value: 0
				"0304" + // Just type: 4
				"045840" + // Just: JustPreVoteHard
				"A1" + // map(1)
				"01583C" + // Certificate (60 bytes)
				"32000000010004010203040094D25422904AC1D130AC981374AA4424F988" + // Certificate Data
				"61E99131078EFEFD62FC52CF072B0C08BB04E4E6496BA48DE4F3D3309AAB" +
				"07f6", // Signature -> Null
			"JustPreVoteHard",
			"BBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB3200000001005052452d564f5445010000",
		},
		{
			"A7" + // map(7)
				"0104" + // Type: 4 (cp:main-vote)
				"021832" + // Height: 50
				"0301" + // Round: 1
				"045820BBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB" + // Block Hash
				"055501AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA" + // Signer
				"06" + // CP_vote
				"A4" + // map(4)
				"0101" + // CP_Round: 1
				"0202" + // CP_Value: 2 (abstain)
				"0305" + // Just type: 5
				"04584b" + // Just: JustTypeMainVoteConflict
				"A4" + // map(4)
				"0101" + // Just0: Type (No)
				"025840" + // Just0Data
				"A1" + // map(1)
				"01583C" + // Certificate (60 bytes)
				"32000000010004010203040094D25422904AC1D130AC981374AA4424F988" + // Certificate Data
				"61E99131078EFEFD62FC52CF072B0C08BB04E4E6496BA48DE4F3D3309AAB" +
				"0302" + // Just1: Type (JustTypeInitYes)
				"0441" + // Just1Data
				"A0" + // Empty Array
				"07f6", // Signature -> Null
			"JustMainVoteConflict",
			"BBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB3200000001004d41494e2d564f5445010002",
		},
		{
			"A7" + // map(7)
				"0104" + // Type: 4 (cp:main-vote)
				"021832" + // Height: 50
				"0301" + // Round: 1
				"0458200000000000000000000000000000000000000000000000000000000000000000" + // Block Hash
				"055501AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA" + // Signer
				"06" + // CP_vote
				"A4" + // map(4)
				"0101" + // CP_Round: 1
				"0201" + // CP_Value: 1
				"0306" + // Just type: 6
				"045840" + // Just: JustTypeMainVoteNoConflict
				"A1" + // map(1)
				"01583C" + // Certificate (60 bytes)
				"32000000010004010203040094D25422904AC1D130AC981374AA4424F988" + // Certificate Data
				"61E99131078EFEFD62FC52CF072B0C08BB04E4E6496BA48DE4F3D3309AAB" +
				"07f6", // Signature -> Null
			"JustMainVoteNoConflict",
			"00000000000000000000000000000000000000000000000000000000000000003200000001004d41494e2d564f5445010001",
		},
		{
			"A7" + // map(7)
				"0105" + // Type: 4 (cp:decided)
				"021832" + // Height: 50
				"0301" + // Round: 1
				"0458200000000000000000000000000000000000000000000000000000000000000000" + // Block Hash
				"055501AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA" + // Signer
				"06" + // CP_vote
				"A4" + // map(4)
				"0101" + // CP_Round: 1
				"0201" + // CP_Value: 1
				"0307" + // Just type: 7
				"045840" + // Just: JustTypeDecided
				"A1" + // map(1)
				"01583C" + // Certificate (60 bytes)
				"32000000010004010203040094D25422904AC1D130AC981374AA4424F988" + // Certificate Data
				"61E99131078EFEFD62FC52CF072B0C08BB04E4E6496BA48DE4F3D3309AAB" +
				"07f6", // Signature -> Null
			"JustDecided",
			"000000000000000000000000000000000000000000000000000000000000000032000000010044454349444544010001",
		},
	}

	ts := testsuite.NewTestSuite(t)
	for _, test := range tests {
		bz1, _ := hex.DecodeString(test.data)

		v := new(vote.Vote)
		err := v.UnmarshalCBOR(bz1)
		assert.NoError(t, err)

		bz2, err := v.MarshalCBOR()
		assert.NoError(t, err)

		assert.Equal(t, bz1, bz2)

		expectedHash := hash.CalcHash(bz1)
		assert.Equal(t, expectedHash, v.Hash())

		v.SetSignature(ts.RandBLSSignature())
		assert.NoError(t, v.BasicCheck())

		expectedSignBytes, _ := hex.DecodeString(test.signBytes)
		assert.Equal(t, expectedSignBytes, v.SignBytes())

		if test.justType != "" {
			assert.Equal(t, test.justType, v.CPJust().Type().String())
		}
	}
}

func TestVoteSignature(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	h1 := ts.RandHash()
	pb1, pv1 := ts.RandBLSKeyPair()
	pb2, pv2 := ts.RandBLSKeyPair()

	v1 := vote.NewPrepareVote(h1, 101, 5, pb1.ValidatorAddress())
	v2 := vote.NewPrepareVote(h1, 101, 5, pb2.ValidatorAddress())

	assert.Error(t, v1.Verify(pb1), "No signature")

	sig1 := pv1.Sign(v1.SignBytes())
	v1.SetSignature(sig1.(*bls.Signature))
	assert.NoError(t, v1.Verify(pb1), "Ok")

	sig2 := pv2.Sign(v2.SignBytes())
	v2.SetSignature(sig2.(*bls.Signature))
	assert.Error(t, v2.Verify(pb1), "invalid public key")

	sig3 := pv1.Sign(v2.SignBytes())
	v2.SetSignature(sig3.(*bls.Signature))
	assert.Error(t, v2.Verify(pb2), "invalid signature")
}

func TestCPPreVote(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	h := ts.RandHeight()
	r := ts.RandRound()
	just := &vote.JustInitYes{}

	t.Run("Invalid round", func(t *testing.T) {
		v := vote.NewCPPreVote(hash.UndefHash, h, r,
			-1, vote.CPValueYes, just, ts.RandAccAddress())

		err := v.BasicCheck()
		assert.Equal(t, errors.ErrInvalidRound, errors.Code(err))
	})

	t.Run("Invalid value", func(t *testing.T) {
		v := vote.NewCPPreVote(hash.UndefHash, h, r,
			1, 3, just, ts.RandAccAddress())

		err := v.BasicCheck()
		assert.Equal(t, errors.ErrInvalidVote, errors.Code(err))
	})

	t.Run("Ok", func(t *testing.T) {
		v := vote.NewCPPreVote(hash.UndefHash, h, r,
			1, vote.CPValueNo, just, ts.RandAccAddress())
		v.SetSignature(ts.RandBLSSignature())

		err := v.BasicCheck()
		assert.NoError(t, err)
		assert.Equal(t, int16(1), v.CPRound())
		assert.Equal(t, vote.CPValueNo, v.CPValue())
		assert.NotNil(t, v.CPJust())
	})
}

func TestCPMainVote(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	h := ts.RandHeight()
	r := ts.RandRound()
	just := &vote.JustInitYes{}

	t.Run("Invalid round", func(t *testing.T) {
		v := vote.NewCPMainVote(hash.UndefHash, h, r,
			-1, vote.CPValueNo, just, ts.RandAccAddress())

		err := v.BasicCheck()
		assert.Equal(t, errors.ErrInvalidRound, errors.Code(err))
	})

	t.Run("No CP data", func(t *testing.T) {
		data, _ := hex.DecodeString("A701040218320301045820BBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB" +
			"055501AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA06f607f6")
		v := new(vote.Vote)
		err := v.UnmarshalCBOR(data)
		assert.NoError(t, err)
		v.SetSignature(ts.RandBLSSignature())

		assert.Error(t, v.BasicCheck())
	})

	t.Run("Invalid value", func(t *testing.T) {
		v := vote.NewCPMainVote(hash.UndefHash, h, r,
			1, 3, just, ts.RandAccAddress())

		err := v.BasicCheck()
		assert.Equal(t, errors.ErrInvalidVote, errors.Code(err))
	})

	t.Run("Ok", func(t *testing.T) {
		v := vote.NewCPMainVote(hash.UndefHash, h, r,
			1, vote.CPValueAbstain, just, ts.RandAccAddress())
		v.SetSignature(ts.RandBLSSignature())

		err := v.BasicCheck()
		assert.NoError(t, err)
		assert.Equal(t, int16(1), v.CPRound())
		assert.Equal(t, vote.CPValueAbstain, v.CPValue())
		assert.NotNil(t, v.CPJust())
	})
}

func TestCPDecided(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	h := ts.RandHeight()
	r := ts.RandRound()
	just := &vote.JustInitYes{}

	t.Run("Invalid round", func(t *testing.T) {
		v := vote.NewCPDecidedVote(hash.UndefHash, h, r,
			-1, vote.CPValueNo, just, ts.RandAccAddress())

		err := v.BasicCheck()
		assert.Equal(t, errors.ErrInvalidRound, errors.Code(err))
	})

	t.Run("No CP data", func(t *testing.T) {
		data, _ := hex.DecodeString("A701050218320301045820BBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB" +
			"055501AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA06f607f6")
		v := new(vote.Vote)
		err := v.UnmarshalCBOR(data)
		assert.NoError(t, err)
		v.SetSignature(ts.RandBLSSignature())

		assert.Error(t, v.BasicCheck())
	})

	t.Run("Invalid value", func(t *testing.T) {
		v := vote.NewCPDecidedVote(hash.UndefHash, h, r,
			1, 3, just, ts.RandAccAddress())

		err := v.BasicCheck()
		assert.Equal(t, errors.ErrInvalidVote, errors.Code(err))
	})

	t.Run("Ok", func(t *testing.T) {
		v := vote.NewCPDecidedVote(hash.UndefHash, h, r,
			1, vote.CPValueAbstain, just, ts.RandAccAddress())
		v.SetSignature(ts.RandBLSSignature())

		err := v.BasicCheck()
		assert.NoError(t, err)
		assert.Equal(t, int16(1), v.CPRound())
		assert.Equal(t, vote.CPValueAbstain, v.CPValue())
		assert.NotNil(t, v.CPJust())
	})
}

func TestBasicCheck(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	t.Run("Invalid type", func(t *testing.T) {
		data, _ := hex.DecodeString("A701050218320301045820BBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB" +
			"055501AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA06f607f6")
		v := new(vote.Vote)
		err := v.UnmarshalCBOR(data)
		assert.NoError(t, err)

		err = v.BasicCheck()
		assert.Equal(t, errors.ErrInvalidVote, errors.Code(err))
	})

	t.Run("Invalid height", func(t *testing.T) {
		v := vote.NewPrepareVote(ts.RandHash(), 0, 0, ts.RandAccAddress())

		err := v.BasicCheck()
		assert.ErrorIs(t, err, vote.BasicCheckError{Reason: "invalid height"})
	})

	t.Run("Invalid round", func(t *testing.T) {
		v := vote.NewPrepareVote(ts.RandHash(), 100, -1, ts.RandAccAddress())

		err := v.BasicCheck()
		assert.ErrorIs(t, err, vote.BasicCheckError{Reason: "invalid round"})
	})

	t.Run("No signature", func(t *testing.T) {
		v := vote.NewPrepareVote(ts.RandHash(), 100, 0, ts.RandAccAddress())

		err := v.BasicCheck()
		assert.Equal(t, errors.ErrInvalidSignature, errors.Code(err))
	})

	t.Run("Has CP data", func(t *testing.T) {
		data, _ := hex.DecodeString("A701020218320301045820BBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB" +
			"055501AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA06A40100020103020441A007f6")
		v := new(vote.Vote)
		err := v.UnmarshalCBOR(data)
		assert.NoError(t, err)
		v.SetSignature(ts.RandBLSSignature())

		assert.Error(t, v.BasicCheck())
	})

	t.Run("Ok", func(t *testing.T) {
		v := vote.NewPrepareVote(ts.RandHash(), 100, 0, ts.RandAccAddress())
		v.SetSignature(ts.RandBLSSignature())

		assert.NoError(t, v.BasicCheck())
	})
}

func TestSignBytes(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	signer := ts.RandAccAddress()
	blockHash := ts.RandHash()
	height := uint32(100)
	round := int16(10)
	cpRound := int16(10)
	just := &vote.JustInitNo{}

	v1 := vote.NewPrepareVote(blockHash, height, round, signer)
	v2 := vote.NewPrecommitVote(blockHash, height, round, signer)
	v3 := vote.NewCPPreVote(blockHash, height, round, cpRound, vote.CPValueNo, just, signer)
	v4 := vote.NewCPMainVote(blockHash, height, round, cpRound, vote.CPValueAbstain, just, signer)
	v5 := vote.NewCPDecidedVote(blockHash, height, round, cpRound, vote.CPValueYes, just, signer)

	sb1 := v1.SignBytes()
	sb2 := v2.SignBytes()
	sb3 := v3.SignBytes()
	sb4 := v4.SignBytes()
	sb5 := v5.SignBytes()

	assert.Equal(t, 45, len(sb1))
	assert.Equal(t, 38, len(sb2))
	assert.Equal(t, 49, len(sb3))
	assert.Equal(t, 50, len(sb4))
	assert.Equal(t, 48, len(sb5))

	assert.Contains(t, string(sb1), "PREPARE")
	assert.Contains(t, string(sb3), "PRE-VOTE")
	assert.Contains(t, string(sb4), "MAIN-VOTE")
	assert.Contains(t, string(sb5), "DECIDED")
}

func TestLog(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	signer := ts.RandAccAddress()
	blockHash := ts.RandHash()
	height := uint32(100)
	round := int16(2)
	just := &vote.JustInitNo{}

	v1 := vote.NewPrepareVote(blockHash, height, round, signer)
	v2 := vote.NewPrecommitVote(blockHash, height, round, signer)
	v3 := vote.NewCPPreVote(blockHash, height, round, 1, vote.CPValueNo, just, signer)
	v4 := vote.NewCPMainVote(blockHash, height, round, 1, vote.CPValueAbstain, just, signer)

	assert.Contains(t, v1.String(), "100/2/PREPARE")
	assert.Contains(t, v2.String(), "100/2/PRECOMMIT")
	assert.Contains(t, v3.String(), "100/2/PRE-VOTE/1")
	assert.Contains(t, v4.String(), "100/2/MAIN-VOTE/1")
}

func TestCPValueToString(t *testing.T) {
	assert.Equal(t, "no", vote.CPValueNo.String())
	assert.Equal(t, "yes", vote.CPValueYes.String())
	assert.Equal(t, "abstain", vote.CPValueAbstain.String())
	assert.Equal(t, "unknown: -1", vote.CPValue(-1).String())
}

func TestCPInvalidJustType(t *testing.T) {
	voteData, _ := hex.DecodeString(
		"A7" + // map(7)
			"0103" + // Type: 3 (cp:pre-vote)
			"021832" + // Height: 50
			"0301" + // Round: 1
			"0458200000000000000000000000000000000000000000000000000000000000000000" + // Block Hash
			"055501AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA" + // Signer
			"06" + // CP_vote
			"A4" + // map(4)
			"0100" + // CP_Round: 0
			"0201" + // CP_Value: 1
			"0308" + // Just type: 8 <<<(Unknown Just Type)>>>
			"0441" + // Just: JustTypeInitYes
			"A0" + // Empty Array
			"07f6") // Signature -> Null

	v := new(vote.Vote)
	err := v.UnmarshalCBOR(voteData)
	assert.Error(t, err)
}
