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

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
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

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
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

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
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
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = bls.PublicKeyFromBytes(pubBytes)
	}
}
