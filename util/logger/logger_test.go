package logger_test

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/pactus-project/pactus/util/logger"
)

type mockFingerprintable struct{}

func (m *mockFingerprintable) Fingerprint() string {
	return "mock-fingerprint"
}

func TestFingerprintable(t *testing.T) {
	f := &mockFingerprintable{}

	assert.Equal(t, "mock-fingerprint", f.Fingerprint())
}

func TestIsNil(t *testing.T) {
	assert.True(t, logger.IsNil(nil))

	obj := struct{}{}
	assert.False(t, logger.IsNil(obj))
}

func TestKeyvalsToFields(t *testing.T) {
	fields := logger.KeyvalsToFields("key1", "value1", "key2", 123, "key3", []byte{0x01, 0x02})
	expectedFields := []zapcore.Field{
		zap.Any("key1", "value1"),
		zap.Any("key2", 123),
		zap.Any("key3", hex.EncodeToString([]byte{0x01, 0x02})),
	}
	assert.Equal(t, expectedFields, fields)
}

func TestParseZapLevel(t *testing.T) {
	assert.Equal(t, zapcore.DebugLevel, logger.ParseZapLevel("debug"))
	assert.Equal(t, zapcore.InfoLevel, logger.ParseZapLevel("info"))
	assert.Equal(t, zapcore.WarnLevel, logger.ParseZapLevel("warning"))
	assert.Equal(t, zapcore.ErrorLevel, logger.ParseZapLevel("error"))
	assert.Equal(t, zapcore.FatalLevel, logger.ParseZapLevel("fatal"))
	assert.Equal(t, zapcore.PanicLevel, logger.ParseZapLevel("panic"))

	assert.Equal(t, zapcore.InfoLevel, logger.ParseZapLevel("invalid-level"))
}
