package logger

import (
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

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

type testFingerprintable struct{}

func (t *testFingerprintable) Fingerprint() string {
	return "test-fingerprint"
}
