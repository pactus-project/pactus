package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type Foo struct{}

func (f Foo) Fingerprint() string {
	return "foo"
}

type Bar struct{}

func (b *Bar) Fingerprint() string {
	return "bar"
}
func TestFingerprint(t *testing.T) {
	f1 := Foo{}
	f2 := &Foo{}
	b1 := Bar{}
	b2 := &Bar{}

	assert.Equal(t, keyvalsToFields("key", f1)["key"], "foo")
	assert.Equal(t, keyvalsToFields("key", &f1)["key"], "foo")
	assert.Equal(t, keyvalsToFields("key", f2)["key"], "foo")
	assert.Equal(t, keyvalsToFields("key", b1)["key"], "{}")
	assert.Equal(t, keyvalsToFields("key", &b1)["key"], "bar")
	assert.Equal(t, keyvalsToFields("key", b2)["key"], "bar")
}

func TestNilFingerprint(t *testing.T) {
	var f1 Foo
	var f2 *Foo
	var b1 Bar
	var b2 *Bar

	assert.Equal(t, keyvalsToFields("key", f1)["key"], "foo")
	assert.Equal(t, keyvalsToFields("key", &f1)["key"], "foo")
	assert.Equal(t, keyvalsToFields("key", f2)["key"], "nil")
	assert.Equal(t, keyvalsToFields("key", b1)["key"], "{}")
	assert.Equal(t, keyvalsToFields("key", &b1)["key"], "bar")
	assert.Equal(t, keyvalsToFields("key", b2)["key"], "nil")
}
