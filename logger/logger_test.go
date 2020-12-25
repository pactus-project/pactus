package logger

import (
	"bytes"
	"log"
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

func TestObjLogger(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
		}
	}()
	l := NewLogger("test", Foo{})
	var buf bytes.Buffer
	log.SetOutput(&buf)

	l.Trace("a")
	l.Debug("b")
	l.Info("c")
	l.Warn("d")
	l.Error("e")
	l.Fatal("f")
	l.Panic("g")

	out := buf.String()

	assert.Contains(t, out, "foo")
	assert.Contains(t, out, "Trace")
	assert.Contains(t, out, "Debug")
	assert.Contains(t, out, "Info")
	assert.Contains(t, out, "Warn")
	assert.Contains(t, out, "Error")
	assert.Contains(t, out, "Fatal")
	assert.Contains(t, out, "Panic")
}

func TestLogger(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
		}
	}()

	var buf bytes.Buffer
	log.SetOutput(&buf)

	Trace("a")
	Debug("b")
	Info("c")
	Warn("d")
	Error("e")
	//Fatal("f")
	Panic("g")

	out := buf.String()

	assert.Contains(t, out, "foo")
	assert.Contains(t, out, "Trace")
	assert.Contains(t, out, "Debug")
	assert.Contains(t, out, "Info")
	assert.Contains(t, out, "Warn")
	assert.Contains(t, out, "Error")
	assert.Contains(t, out, "Fatal")
	assert.Contains(t, out, "Panic")
}
