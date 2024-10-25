package logger

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

type Foo struct{}

func (Foo) String() string {
	return "foo"
}

type Bar struct{}

func (Bar) ShortString() string {
	return "bar"
}

func TestNilObjLogger(t *testing.T) {
	l := NewSubLogger("test", nil)
	var buf bytes.Buffer
	l.logger = l.logger.Output(&buf)

	l.Info("hello", "error", fmt.Errorf("error"))
	assert.Contains(t, buf.String(), "hello")
	assert.Contains(t, buf.String(), "error")
}

func TestObjLogger(t *testing.T) {
	globalInst = nil
	c := DefaultConfig()
	c.Colorful = false
	InitGlobalLogger(c)

	globalInst.config.Levels["test"] = "warn"
	subLogger := NewSubLogger("test", Foo{})
	var buf bytes.Buffer
	subLogger.logger = subLogger.logger.Output(&buf)

	subLogger.Trace("msg")
	subLogger.Debug("msg")
	subLogger.Info("msg")
	subLogger.Warn("msg")
	subLogger.Error("msg")

	out := buf.String()

	assert.Contains(t, out, "foo")
	assert.NotContains(t, out, "trace")
	assert.NotContains(t, out, "debug")
	assert.NotContains(t, out, "info")
	assert.Contains(t, out, "warn")
	assert.Contains(t, out, "error")
}

func TestLogger(t *testing.T) {
	globalInst = nil
	c := DefaultConfig()
	c.Colorful = true
	InitGlobalLogger(c)

	var buf bytes.Buffer
	log.Logger = log.Output(&buf)

	Trace("msg", "trace", "trace")
	Debug("msg", "trace", "trace")
	Info("msg", nil)
	Info("msg", "a", nil)
	Info("msg", "b", []byte{1, 2, 3})
	Warn("msg", "x")
	Error("msg", "y", Foo{})

	out := buf.String()

	fmt.Println(out)
	assert.NotContains(t, out, "trace")
	assert.NotContains(t, out, "debug")
	assert.Contains(t, out, "foo")
	assert.Contains(t, out, "010203")
	assert.Contains(t, out, "!INVALID-KEY!")
	assert.Contains(t, out, "!MISSING-VALUE!")
	assert.Contains(t, out, "null")
	assert.NotContains(t, out, "trace")
	assert.NotContains(t, out, "debug")
	assert.Contains(t, out, "info")
	assert.Contains(t, out, "warn")
	assert.Contains(t, out, "error")
}

func TestNilValue(t *testing.T) {
	var buf bytes.Buffer
	log.Logger = log.Output(&buf)

	foo := new(Foo)
	if true {
		// to avoid some linting errors
		foo = nil
	}

	Info("msg", "null", nil)
	Info("msg", "error", error(nil))
	Info("msg", "stringer", foo)

	out := buf.String()

	fmt.Println(out)
	assert.Contains(t, out, "null")
	assert.Contains(t, out, "error")
	assert.Contains(t, out, "stringer")
}

func TestInvalidLevel(t *testing.T) {
	globalInst = nil
	c := DefaultConfig()
	c.Colorful = false
	InitGlobalLogger(c)

	var buf bytes.Buffer
	log.Logger = log.Logger.Output(&buf)

	globalInst.config.Levels["test"] = "invalid"
	l := NewSubLogger("test", Foo{})

	var buf2 bytes.Buffer
	l.logger = l.logger.Output(&buf2)

	l.Error("message", "key", "val")

	out := buf.String()
	out2 := buf2.String()

	fmt.Println(out)
	assert.Contains(t, out, "Unknown Level String")
	assert.NotContains(t, out2, "error")
	assert.NotContains(t, out2, "message")
}

func TestShortStringer(t *testing.T) {
	globalInst = nil
	c := DefaultConfig()
	c.Colorful = false
	InitGlobalLogger(c)

	l := NewSubLogger("test", &Foo{})
	var buf bytes.Buffer
	l.logger = l.logger.Output(&buf)

	l.Info("msg", "f", Bar{})

	out := buf.String()

	assert.Contains(t, out, "bar")
}
