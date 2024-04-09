package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServicesString(t *testing.T) {
	assert.Equal(t, New(None).String(), "")
	assert.Equal(t, New(Network).String(), "NETWORK")
	assert.Equal(t, New(Foo).String(), "FOO")
	assert.Equal(t, New(Network, Foo).String(), "NETWORK | FOO")
	assert.Equal(t, New(5).String(), "NETWORK | 4")
}

func TestAppend(t *testing.T) {
	s := New(Network)
	assert.True(t, s.IsNetwork())
	assert.False(t, s.IsFoo())

	s.Append(Foo)
	assert.True(t, s.IsNetwork())
	assert.True(t, s.IsFoo())
}

func TestIsNetwork(t *testing.T) {
	assert.False(t, New(None).IsNetwork())
	assert.True(t, New(Network).IsNetwork())
	assert.False(t, New(Foo).IsNetwork())
	assert.True(t, New(Foo, Network).IsNetwork())
}

func TestIsFoo(t *testing.T) {
	assert.False(t, New(None).IsFoo())
	assert.False(t, New(Network).IsFoo())
	assert.True(t, New(Foo).IsFoo())
	assert.True(t, New(Foo, Network).IsNetwork())
}
