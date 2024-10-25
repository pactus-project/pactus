package vote_test

import (
	"encoding/hex"
	"testing"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/types/vote"
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
	for _, tt := range tests {
		bz1, _ := hex.DecodeString(tt.data)

		vote := new(vote.Vote)
		err := vote.UnmarshalCBOR(bz1)
		assert.NoError(t, err)

		bz2, err := vote.MarshalCBOR()
		assert.NoError(t, err)

		assert.Equal(t, bz1, bz2)

		expectedHash := hash.CalcHash(bz1)
		assert.Equal(t, expectedHash, vote.Hash())

		vote.SetSignature(ts.RandBLSSignature())
		assert.NoError(t, vote.BasicCheck())

		expectedSignBytes, _ := hex.DecodeString(tt.signBytes)
		assert.Equal(t, expectedSignBytes, vote.SignBytes())

		if tt.justType != "" {
			assert.Equal(t, tt.justType, vote.CPJust().Type().String())
		}
	}
}

func TestVoteSignature(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	hash1 := ts.RandHash()
	pub1, prv1 := ts.RandBLSKeyPair()
	pub2, prv2 := ts.RandBLSKeyPair()

	vote1 := vote.NewPrepareVote(hash1, 101, 5, pub1.ValidatorAddress())
	vote2 := vote.NewPrepareVote(hash1, 101, 5, pub2.ValidatorAddress())

	assert.Error(t, vote1.BasicCheck(), "No signature")

	sig1 := prv1.SignNative(vote1.SignBytes())
	vote1.SetSignature(sig1)
	err1 := vote1.Verify(pub1)
	assert.NoError(t, err1, "Ok")

	sig2 := prv2.SignNative(vote2.SignBytes())
	vote2.SetSignature(sig2)
	err2 := vote2.Verify(pub1)
	assert.ErrorIs(t, err2, vote.InvalidSignerError{
		Expected: pub1.ValidatorAddress(),
		Got:      pub2.ValidatorAddress(),
	})

	sig3 := prv1.SignNative(vote2.SignBytes())
	vote2.SetSignature(sig3)
	err3 := vote2.Verify(pub2)
	assert.ErrorIs(t, err3, crypto.ErrInvalidSignature)
}

func TestCPPreVote(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	height := ts.RandHeight()
	round := ts.RandRound()
	just := &vote.JustInitYes{}

	t.Run("Invalid CP round", func(t *testing.T) {
		invalidCPRound := int16(-1)
		cpVote := vote.NewCPPreVote(hash.UndefHash, height, round,
			invalidCPRound, vote.CPValueYes, just, ts.RandAccAddress())

		err := cpVote.BasicCheck()
		assert.ErrorIs(t, err, vote.BasicCheckError{Reason: "invalid CP round"})
	})

	t.Run("invalid CP value", func(t *testing.T) {
		invalidCPValue := vote.CPValue(3)
		cpVote := vote.NewCPPreVote(hash.UndefHash, height, round,
			1, invalidCPValue, just, ts.RandAccAddress())

		err := cpVote.BasicCheck()
		assert.ErrorIs(t, err, vote.BasicCheckError{Reason: "invalid CP value"})
	})

	t.Run("Ok", func(t *testing.T) {
		cpVote := vote.NewCPPreVote(hash.UndefHash, height, round,
			1, vote.CPValueNo, just, ts.RandAccAddress())
		cpVote.SetSignature(ts.RandBLSSignature())

		err := cpVote.BasicCheck()
		assert.NoError(t, err)
		assert.Equal(t, int16(1), cpVote.CPRound())
		assert.Equal(t, vote.CPValueNo, cpVote.CPValue())
		assert.NotNil(t, cpVote.CPJust())
	})
}

func TestCPMainVote(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	height := ts.RandHeight()
	round := ts.RandRound()
	just := &vote.JustInitYes{}

	t.Run("Invalid CP round", func(t *testing.T) {
		invalidCPRound := int16(-1)
		invVote := vote.NewCPMainVote(hash.UndefHash, height, round,
			invalidCPRound, vote.CPValueNo, just, ts.RandAccAddress())

		err := invVote.BasicCheck()
		assert.ErrorIs(t, err, vote.BasicCheckError{Reason: "invalid CP round"})
	})

	t.Run("No CP data", func(t *testing.T) {
		data, _ := hex.DecodeString(
			"A701040218320301045820BBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB" +
				"055501AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA06f607f6")
		decVote := new(vote.Vote)
		_ = decVote.UnmarshalCBOR(data)
		decVote.SetSignature(ts.RandBLSSignature())

		err := decVote.BasicCheck()
		assert.ErrorIs(t, err, vote.BasicCheckError{Reason: "should have CP data"})
	})

	t.Run("Invalid CP value", func(t *testing.T) {
		invalidCPValue := vote.CPValue(3)
		cpVote := vote.NewCPMainVote(hash.UndefHash, height, round,
			1, invalidCPValue, just, ts.RandAccAddress())

		err := cpVote.BasicCheck()
		assert.ErrorIs(t, err, vote.BasicCheckError{Reason: "invalid CP value"})
	})

	t.Run("Ok", func(t *testing.T) {
		cpVote := vote.NewCPMainVote(hash.UndefHash, height, round,
			1, vote.CPValueAbstain, just, ts.RandAccAddress())
		cpVote.SetSignature(ts.RandBLSSignature())

		err := cpVote.BasicCheck()
		assert.NoError(t, err)
		assert.Equal(t, int16(1), cpVote.CPRound())
		assert.Equal(t, vote.CPValueAbstain, cpVote.CPValue())
		assert.NotNil(t, cpVote.CPJust())
	})
}

func TestCPDecided(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	height := ts.RandHeight()
	round := ts.RandRound()
	just := &vote.JustInitYes{}

	t.Run("Invalid round", func(t *testing.T) {
		invalidCPRound := int16(-1)
		v := vote.NewCPDecidedVote(hash.UndefHash, height, round,
			invalidCPRound, vote.CPValueNo, just, ts.RandAccAddress())

		err := v.BasicCheck()
		assert.ErrorIs(t, err, vote.BasicCheckError{Reason: "invalid CP round"})
	})

	t.Run("No CP data", func(t *testing.T) {
		data, _ := hex.DecodeString("A701050218320301045820BBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB" +
			"055501AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA06f607f6")
		v := new(vote.Vote)
		_ = v.UnmarshalCBOR(data)
		v.SetSignature(ts.RandBLSSignature())

		err := v.BasicCheck()
		assert.ErrorIs(t, err, vote.BasicCheckError{Reason: "should have CP data"})
	})

	t.Run("Invalid CP value", func(t *testing.T) {
		invalidCPValue := vote.CPValue(3)
		v := vote.NewCPDecidedVote(hash.UndefHash, height, round,
			1, invalidCPValue, just, ts.RandAccAddress())

		err := v.BasicCheck()
		assert.ErrorIs(t, err, vote.BasicCheckError{Reason: "invalid CP value"})
	})

	t.Run("Ok", func(t *testing.T) {
		vte := vote.NewCPDecidedVote(hash.UndefHash, height, round,
			1, vote.CPValueAbstain, just, ts.RandAccAddress())
		vte.SetSignature(ts.RandBLSSignature())

		err := vte.BasicCheck()
		assert.NoError(t, err)
		assert.Equal(t, int16(1), vte.CPRound())
		assert.Equal(t, vote.CPValueAbstain, vte.CPValue())
		assert.NotNil(t, vte.CPJust())
	})
}

func TestBasicCheck(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	t.Run("Should have CP data", func(t *testing.T) {
		data, _ := hex.DecodeString("A701050218320301045820BBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB" +
			"055501AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA06f607f6")
		v := new(vote.Vote)
		err := v.UnmarshalCBOR(data)
		assert.NoError(t, err)

		err = v.BasicCheck()
		assert.ErrorIs(t, err, vote.BasicCheckError{Reason: "should have CP data"})
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
		assert.ErrorIs(t, err, vote.BasicCheckError{Reason: "no signature"})
	})

	t.Run("Should not have CP data", func(t *testing.T) {
		data, _ := hex.DecodeString("A701020218320301045820BBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB" +
			"055501AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA06A40100020103020441A007f6")
		v := new(vote.Vote)
		_ = v.UnmarshalCBOR(data)
		v.SetSignature(ts.RandBLSSignature())

		err := v.BasicCheck()
		assert.ErrorIs(t, err, vote.BasicCheckError{Reason: "should not have CP data"})
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

	vote1 := vote.NewPrepareVote(blockHash, height, round, signer)
	vote2 := vote.NewPrecommitVote(blockHash, height, round, signer)
	vote3 := vote.NewCPPreVote(blockHash, height, round, cpRound, vote.CPValueNo, just, signer)
	vote4 := vote.NewCPMainVote(blockHash, height, round, cpRound, vote.CPValueAbstain, just, signer)
	vote5 := vote.NewCPDecidedVote(blockHash, height, round, cpRound, vote.CPValueYes, just, signer)

	sby1 := vote1.SignBytes()
	sby2 := vote2.SignBytes()
	sby3 := vote3.SignBytes()
	sby4 := vote4.SignBytes()
	sby5 := vote5.SignBytes()

	assert.Equal(t, 45, len(sby1))
	assert.Equal(t, 38, len(sby2))
	assert.Equal(t, 49, len(sby3))
	assert.Equal(t, 50, len(sby4))
	assert.Equal(t, 48, len(sby5))

	assert.Contains(t, string(sby1), "PREPARE")
	assert.Contains(t, string(sby3), "PRE-VOTE")
	assert.Contains(t, string(sby4), "MAIN-VOTE")
	assert.Contains(t, string(sby5), "DECIDED")
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
