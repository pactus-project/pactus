package logger

import (
	"bytes"
	"testing"

	"github.com/rs/zerolog"
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
	assert.Equal(t, keyvalsToFields("key", &b1)["key"], "bar")
	assert.Equal(t, keyvalsToFields("key", b2)["key"], "bar")
	assert.Nil(t, keyvalsToFields(1)["key"])
	assert.Nil(t, keyvalsToFields(nil, 1)["key"])
	assert.Nil(t, keyvalsToFields(1, nil)["key"])
}

func TestNilFingerprint(t *testing.T) {
	var f1 Foo
	var f2 *Foo
	var b1 Bar
	var b2 *Bar

	assert.Equal(t, keyvalsToFields("key", f1)["key"], "foo")
	assert.Equal(t, keyvalsToFields("key", &f1)["key"], "foo")
	assert.Equal(t, keyvalsToFields("key", f2)["key"], nil)
	assert.Equal(t, keyvalsToFields("key", &b1)["key"], "bar")
	assert.Equal(t, keyvalsToFields("key", b2)["key"], nil)
	assert.Nil(t, keyvalsToFields(1)["key"])
	assert.Nil(t, keyvalsToFields(nil, 1)["key"])
	assert.Nil(t, keyvalsToFields(1, nil)["key"])
}

func TestNewLogger(t *testing.T) {
	logger := NewLogger("test", nil)
	assert.NotNil(t, logger)
	assert.Equal(t, "test", logger.name)
	assert.Nil(t, logger.obj)

	loggersMap := getLoggersInst().loggers
	assert.Contains(t, loggersMap, "test")
	assert.Equal(t, logger, loggersMap["test"])
}

func TestInitLogger(t *testing.T) {
	inst := getLoggersInst()
	assert.NoError(t, InitLogger(inst.config), "no error")
}

func TestSetLevel(t *testing.T) {
	logger := NewLogger("test", nil)

	logger.SetLevel(zerolog.ErrorLevel)
	assert.Equal(t, zerolog.ErrorLevel, logger.logger.GetLevel())
}

func TestLog(t *testing.T) {
	logger := NewLogger("test", nil)

	assert.Panics(t, func() { logger.Panic("should panic") }, "Expected Panic to panic")
}
func TestWith(t *testing.T) {
	logger := NewLogger("test", nil)

	event := logger.With("key1", "value1", "key2", "value2")
	assert.NotNil(t, event)
}

func TestPackageLoggingFunctions(t *testing.T) {
	Trace("Trace message")
	Debug("Debug message")
	Info("Info message")
	Warn("Warn message")
	Error("Error message")
	Fatal("Fatal message")
	assert.Panics(t, func() { Panic("should panic") }, "Expected Panic to panic")
}

func TestKeyvalsToFields(t *testing.T) {
	fields := keyvalsToFields("key1", "value1", "key2", "value2", "fingerprint", &testFingerprintable{})
	assert.NotNil(t, fields)
	assert.Equal(t, "value1", fields["key1"])
	assert.Equal(t, "value2", fields["key2"])
	assert.Equal(t, "test-fingerprint", fields["fingerprint"])
}

func TestLogger(t *testing.T) {
	l := NewLogger("test", nil)

	var buf bytes.Buffer
	l2 := l.logger.Output(&buf)
	l.logger = l2

	l.Debug("debug log")
	l.Info("info log")
	l.Error("info log")

	out := buf.String()

	assert.Contains(t, out, "error")
	assert.Contains(t, out, "debug")
	assert.Contains(t, out, "info")
	assert.NotContains(t, out, "trace")
	assert.NotContains(t, out, "panci")
}

type testFingerprintable struct{}

func (t *testFingerprintable) Fingerprint() string {
	return "test-fingerprint"
}
