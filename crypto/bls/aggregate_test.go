package bls_test

import (
	"testing"

	bls12381 "github.com/consensys/gnark-crypto/ecc/bls12-381"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAggregateMatchesAffineSum(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	msg := []byte("pactus")
	pubs := make([]*bls.PublicKey, 0, 52)
	sigs := make([]*bls.Signature, 0, 52)
	for range 51 {
		pub, prv := ts.RandBLSKeyPair()
		pubs = append(pubs, pub)
		sigs = append(sigs, prv.SignNative(msg))
	}
	// A repeated point makes the sum double a point.
	pubs = append(pubs, pubs[0])
	sigs = append(sigs, sigs[0])

	wantPub := new(bls12381.G2Affine)
	wantSig := new(bls12381.G1Affine)
	for i := range pubs {
		pointG2, err := pubs[i].PointG2()
		require.NoError(t, err)
		pointG1, err := sigs[i].PointG1()
		require.NoError(t, err)
		if i == 0 {
			wantPub.Set(pointG2)
			wantSig.Set(pointG1)

			continue
		}
		wantPub.Add(wantPub, pointG2)
		wantSig.Add(wantSig, pointG1)
	}

	aggPub, err := bls.PublicKeyAggregate(pubs...)
	require.NoError(t, err)
	wantPubData := wantPub.Bytes()
	assert.Equal(t, wantPubData[:], aggPub.Bytes())

	aggSig, err := bls.SignatureAggregate(sigs...)
	require.NoError(t, err)
	wantSigData := wantSig.Bytes()
	assert.Equal(t, wantSigData[:], aggSig.Bytes())

	gotPoint, err := aggPub.PointG2()
	require.NoError(t, err)
	assert.True(t, gotPoint.Equal(wantPub))
}

func BenchmarkPublicKeyAggregate(b *testing.B) {
	// The size of the committee.
	pubs := make([]*bls.PublicKey, 51)
	for i := range pubs {
		ikm := make([]byte, 32)
		ikm[0] = byte(i)
		prv, err := bls.KeyGen(ikm, nil)
		require.NoError(b, err)
		pubs[i] = prv.PublicKeyNative()
		_, err = pubs[i].PointG2()
		require.NoError(b, err)
	}

	b.ResetTimer()
	for b.Loop() {
		_, _ = bls.PublicKeyAggregate(pubs...)
	}
}
