package logger

import (
	"encoding/hex"
	"reflect"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type fingerprintable interface {
	Fingerprint() string
}

type loggers struct {
	config  *Config
	loggers map[string]*Logger
}

var (
	loggersInst *loggers
)

type Logger struct {
	logger *zap.Logger
	name   string
	obj    interface{}
}

func getLoggersInst() *loggers {
	if loggersInst == nil {
		// Only during tests the loggersInst is nil

		conf := &Config{
			Levels:   make(map[string]string),
			Colorful: true,
		}
		conf.Levels["default"] = "debug"
		conf.Levels["_network"] = "debug"
		conf.Levels["_consensus"] = "debug"
		conf.Levels["_state"] = "debug"
		conf.Levels["_sync"] = "debug"
		conf.Levels["_pool"] = "debug"
		conf.Levels["_http"] = "debug"
		conf.Levels["_grpc"] = "debug"
		loggersInst = &loggers{
			config:  conf,
			loggers: make(map[string]*Logger),
		}
	}

	return loggersInst
}

func InitLogger(conf *Config) {
	if loggersInst == nil {
		loggersInst = &loggers{
			config:  conf,
			loggers: make(map[string]*Logger),
		}

		var cfg zap.Config
		if conf.Colorful {
			cfg = zap.NewDevelopmentConfig()
		} else {
			cfg = zap.NewProductionConfig()
		}

		cfg.Level = zap.NewAtomicLevelAt(ParseZapLevel(conf.Levels["default"]))

		logger, err := cfg.Build()
		if err == nil {
			zap.ReplaceGlobals(logger)
		}
	}
}

func NewLogger(name string, obj interface{}) *Logger {
	l := &Logger{
		name: name,
		obj:  obj,
	}

	cfg := zap.NewDevelopmentConfig()
	if !getLoggersInst().config.Colorful {
		cfg = zap.NewProductionConfig()
	}

	cfg.Level = zap.NewAtomicLevelAt(ParseZapLevel(getLoggersInst().config.Levels[name]))

	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}

	l.logger = logger

	getLoggersInst().loggers[name] = l
	return l
}

func IsNil(i interface{}) bool {
	if i == nil {
		return true
	}
	switch reflect.TypeOf(i).Kind() {
	case reflect.Ptr:
		return reflect.ValueOf(i).IsNil()
	}
	return false
}

func KeyvalsToFields(keyvals ...interface{}) []zapcore.Field {
	if len(keyvals)%2 != 0 {
		keyvals = append(keyvals, "<MISSING VALUE>")
	}

	fields := make([]zapcore.Field, 0, len(keyvals)/2)
	for i := 0; i < len(keyvals); i += 2 {
		key, ok := keyvals[i].(string)
		if !ok {
			key = "invalid-key"
		}

		var val interface{}
		switch v := keyvals[i+1].(type) {
		case fingerprintable:
			if !IsNil(v) {
				val = v.Fingerprint()
			}
		case []byte:
			{
				val = hex.EncodeToString(v)
			}
		default:
			val = keyvals[i+1]
		}

		fields = append(fields, zap.Any(key, val))
	}

	return fields
}

func ParseZapLevel(level string) zapcore.Level {
	switch level {
	case "trace":
		return zapcore.DebugLevel
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warning":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "fatal":
		return zapcore.FatalLevel
	case "panic":
		return zapcore.PanicLevel
	default:
		return zapcore.InfoLevel
	}
}

func (l *Logger) log(level zapcore.Level, msg string, keyvals ...interface{}) {
	if ce := l.logger.Check(level, msg); ce != nil {
		ce.Write(KeyvalsToFields(keyvals...)...)
	}
}

func (l *Logger) With(keyvals ...interface{}) *zap.Logger {
	fields := KeyvalsToFields(keyvals...)
	return l.logger.With(fields...)
}

func (l *Logger) Trace(msg string, keyvals ...interface{}) {
	l.log(zapcore.DebugLevel, msg, keyvals...)
}

func (l *Logger) Debug(msg string, keyvals ...interface{}) {
	l.log(zapcore.DebugLevel, msg, keyvals...)
}

func (l *Logger) Info(msg string, keyvals ...interface{}) {
	l.log(zapcore.InfoLevel, msg, keyvals...)
}

func (l *Logger) Warn(msg string, keyvals ...interface{}) {
	l.log(zapcore.WarnLevel, msg, keyvals...)
}

func (l *Logger) Error(msg string, keyvals ...interface{}) {
	l.log(zapcore.ErrorLevel, msg, keyvals...)
}

func (l *Logger) Fatal(msg string, keyvals ...interface{}) {
	l.log(zapcore.FatalLevel, msg, keyvals...)
}

func (l *Logger) Panic(msg string, keyvals ...interface{}) {
	l.log(zapcore.PanicLevel, msg, keyvals...)
}

func log(level zapcore.Level, msg string, keyvals ...interface{}) {
	if ce := zap.L().Check(level, msg); ce != nil {
		ce.Write(KeyvalsToFields(keyvals...)...)
	}
}

func Trace(msg string, keyvals ...interface{}) {
	log(zapcore.DebugLevel, msg, keyvals...)
}

func Debug(msg string, keyvals ...interface{}) {
	log(zapcore.DebugLevel, msg, keyvals...)
}

func Info(msg string, keyvals ...interface{}) {
	log(zapcore.InfoLevel, msg, keyvals...)
}

func Warn(msg string, keyvals ...interface{}) {
	log(zapcore.WarnLevel, msg, keyvals...)
}

func Error(msg string, keyvals ...interface{}) {
	log(zapcore.ErrorLevel, msg, keyvals...)
}

func Fatal(msg string, keyvals ...interface{}) {
	log(zapcore.FatalLevel, msg, keyvals...)
}

func Panic(msg string, keyvals ...interface{}) {
	log(zapcore.PanicLevel, msg, keyvals...)
}
