package bls_test

import (
	"crypto/rand"
	"testing"

	"github.com/pactus-project/pactus/crypto/bls"
)

func BenchmarkEncode(b *testing.B) {
	b.ReportAllocs()

	buf := make([]byte, bls.PrivateKeySize)
	_, _ = rand.Read(buf)

	prv, _ := bls.PrivateKeyFromBytes(buf)
	pub := prv.PublicKeyNative()

	for b.Loop() {
		_ = pub.Bytes()
	}
}

func BenchmarkDecodeSign(b *testing.B) {
	b.ReportAllocs()

	buf := make([]byte, bls.PrivateKeySize)
	_, _ = rand.Read(buf)

	prv, _ := bls.PrivateKeyFromBytes(buf)
	bufMsg := []byte("pactus")
	sig := prv.Sign(bufMsg)
	sigBytes := sig.Bytes()

	for b.Loop() {
		_, _ = bls.SignatureFromBytes(sigBytes)
	}
}

func BenchmarkVerify(b *testing.B) {
	b.ReportAllocs()

	buf := make([]byte, bls.PrivateKeySize)
	_, _ = rand.Read(buf)
	prv, _ := bls.PrivateKeyFromBytes(buf)
	pub := prv.PublicKeyNative()

	bufMsg := []byte("pactus")
	sig1 := prv.Sign(bufMsg)

	for b.Loop() {
		_ = pub.Verify(bufMsg, sig1)
	}
}

func BenchmarkDecode(b *testing.B) {
	b.ReportAllocs()

	buf := make([]byte, bls.PrivateKeySize)
	_, _ = rand.Read(buf)
	prv, _ := bls.PrivateKeyFromBytes(buf)
	pub := prv.PublicKeyNative()
	pubBytes := pub.Bytes()

	for b.Loop() {
		_, _ = bls.PublicKeyFromBytes(pubBytes)
	}
}
