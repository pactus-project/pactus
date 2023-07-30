package logger

// import (
// 	"bytes"
// 	// "fmt"
// 	"testing"

// 	// "github.com/rs/zerolog"
// 	"github.com/stretchr/testify/assert"
// )

// type Foo struct{}

// func (f Foo) String() string {
// 	return "foo"
// }

// type Bar struct{}

// func (b *Bar) String() string {
// 	return "bar"
// }

// func TestFingerprint(t *testing.T) {
// 	f1 := Foo{}
// 	f2 := &Foo{}
// 	b1 := Bar{}
// 	b2 := &Bar{}

// 	assert.Equal(t, keyvalsToFields("key", f1)["key"], "foo")
// 	assert.Equal(t, keyvalsToFields("key", &f1)["key"], "foo")
// 	assert.Equal(t, keyvalsToFields("key", f2)["key"], "foo")
// 	assert.Equal(t, keyvalsToFields("key", b1)["key"], "{}")
// 	assert.Equal(t, keyvalsToFields("key", &b1)["key"], "bar")
// 	assert.Equal(t, keyvalsToFields("key", b2)["key"], "bar")
// 	assert.Nil(t, keyvalsToFields(1)["key"])
// 	assert.Nil(t, keyvalsToFields(nil, 1)["key"])
// 	assert.Nil(t, keyvalsToFields(1, nil)["key"])
// }

// func TestNilFingerprint(t *testing.T) {
// 	var f1 Foo
// 	var f2 *Foo
// 	var b1 Bar
// 	var b2 *Bar

// 	assert.Equal(t, keyvalsToFields("key", f1)["key"], "foo")
// 	assert.Equal(t, keyvalsToFields("key", &f1)["key"], "foo")
// 	assert.Equal(t, keyvalsToFields("key", f2)["key"], "nil")
// 	assert.Equal(t, keyvalsToFields("key", b1)["key"], "{}")
// 	assert.Equal(t, keyvalsToFields("key", &b1)["key"], "bar")
// 	assert.Equal(t, keyvalsToFields("key", b2)["key"], "nil")
// 	assert.Nil(t, keyvalsToFields(1)["key"])
// 	assert.Nil(t, keyvalsToFields(nil, 1)["key"])
// 	assert.Nil(t, keyvalsToFields(1, nil)["key"])
// }

// func TestObjLogger(t *testing.T) {
// 	globalInst = nil
// 	c := DefaultConfig()
// 	c.Colorful = false
// 	InitGlobalLogger(c)

// 	l := NewSubLogger("test", Foo{})
// 	var buf bytes.Buffer
// 	l.logger.Output(&buf)

// 	l.Trace("a")
// 	l.Debug("b")
// 	l.Info("c")
// 	l.Warn("d")
// 	l.Error("e")

// 	out := buf.String()

// 	assert.Contains(t, out, "foo")
// 	assert.NotContains(t, out, "trace")
// 	assert.NotContains(t, out, "debug")
// 	assert.Contains(t, out, "info")
// 	assert.Contains(t, out, "warn")
// 	assert.Contains(t, out, "err")
// }

// func TestLogger(t *testing.T) {
// 	globalInst = nil
// 	c := DefaultConfig()
// 	c.Colorful = true
// 	InitGlobalLogger(c)

// 	var buf bytes.Buffer
// 	zerolog.DefaultContextLogger.Output(&buf)

// 	Trace("a")
// 	Debug("b", "a", nil)
// 	Info("c", "b", []byte{1, 2, 3})
// 	Warn("d", "x")
// 	Error("e", "y", Foo{})

// 	out := buf.String()

// 	fmt.Println(out)
// 	assert.Contains(t, out, "foo")
// 	assert.Contains(t, out, "010203")
// 	assert.Contains(t, out, "<MISSING VALUE>")
// 	assert.NotContains(t, out, "TRACE")
// 	assert.NotContains(t, out, "DEBU")
// 	assert.Contains(t, out, "INFO")
// 	assert.Contains(t, out, "WARN")
// 	assert.Contains(t, out, "ERR")
// }
