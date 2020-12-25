package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSigner(t *testing.T) {
	_, _, priv := GenerateTestKeyPair()
	s := NewSigner(priv)
	pub := s.PublicKey()
	assert.True(t, pub.EqualsTo(priv.PublicKey()))
	assert.True(t, s.Address().EqualsTo(priv.PublicKey().Address()))
	assert.True(t, s.Sign([]byte("zarb")).EqualsTo(*priv.Sign([]byte("zarb"))))
}
