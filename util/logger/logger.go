package logger

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"reflect"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
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
	logger zerolog.Logger
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

func InitLogger(conf *Config) error {
	if loggersInst == nil {
		loggersInst = &loggers{
			config:  conf,
			loggers: make(map[string]*Logger),
		}
		lvl, err := zerolog.ParseLevel(conf.Levels["default"])
		if err != nil {
			return err
		}
		zerolog.SetGlobalLevel(lvl)
	}
	return nil
}

func NewLogger(name string, obj interface{}) *Logger {
	l := &Logger{
		logger: zerolog.New(&bytes.Buffer{}),
		name:   name,
		obj:    obj,
	}

	lvl := getLoggersInst().config.Levels[name]
	if lvl == "" {
		lvl = getLoggersInst().config.Levels["default"]
	}

	level, err := zerolog.ParseLevel(lvl)
	if err == nil {
		l.logger.Level(level)
	}

	getLoggersInst().loggers[name] = l
	return l
}

func isNil(i interface{}) bool {
	if i == nil {
		return true
	}
	switch reflect.TypeOf(i).Kind() {
	case reflect.Ptr:
		return reflect.ValueOf(i).IsNil()
	}
	return false
}

func keyvalsToFields(keyvals ...interface{}) map[string]interface{} {
	if len(keyvals)%2 != 0 {
		keyvals = append(keyvals, "<MISSING VALUE>")
	}
	fields := make(map[string]interface{})
	for i := 0; i < len(keyvals); i += 2 {
		key, ok := keyvals[i].(string)
		if !ok {
			key = "invalid-key"
		}

		var val interface{}
		switch v := keyvals[i+1].(type) {
		case fingerprintable:
			if !isNil(v) {
				val = fmt.Sprintf("%v", v.Fingerprint())
			}
		case []byte:
			{
				val = fmt.Sprintf("%v", hex.EncodeToString(v))
			}
		default:
			val = keyvals[i+1]
		}
		fields[key] = val
	}

	return fields
}

func (l *Logger) SetLevel(level zerolog.Level) {
	l.logger = l.logger.Level(level)
}

func (l *Logger) Log(level zerolog.Level, msg string, keyvals ...interface{}) {
	switch level {
	case zerolog.DebugLevel:
		l.logger.Debug().Msgf(msg, keyvalsToFields(keyvals...))
	case zerolog.InfoLevel:
		l.logger.Info().Msgf(msg, keyvalsToFields(keyvals...))
	case zerolog.WarnLevel:
		l.logger.Warn().Msgf(msg, keyvalsToFields(keyvals...))
	case zerolog.ErrorLevel:
		l.logger.Error().Msgf(msg, keyvalsToFields(keyvals...))
	case zerolog.PanicLevel:
		l.logger.Panic().Msgf(msg, keyvalsToFields(keyvals...))
	case zerolog.TraceLevel:
		l.logger.Trace().Msgf(msg, keyvalsToFields(keyvals...))
	default:
		l.logger.Log().Msgf(msg, keyvalsToFields(keyvals...))
	}
}

func (l *Logger) With(keyvals ...interface{}) *zerolog.Event {
	fields := keyvalsToFields(keyvals...)
	return l.logger.Log().Fields(fields)
}

func (l *Logger) Trace(msg string, keyvals ...interface{}) {
	l.Log(zerolog.TraceLevel, msg, keyvals...)
}

func (l *Logger) Debug(msg string, keyvals ...interface{}) {
	l.Log(zerolog.DebugLevel, msg, keyvals...)
}

func (l *Logger) Info(msg string, keyvals ...interface{}) {
	l.Log(zerolog.InfoLevel, msg, keyvals...)
}

func (l *Logger) Warn(msg string, keyvals ...interface{}) {
	l.Log(zerolog.WarnLevel, msg, keyvals...)
}

func (l *Logger) Error(msg string, keyvals ...interface{}) {
	l.Log(zerolog.ErrorLevel, msg, keyvals...)
}

func (l *Logger) Fatal(msg string, keyvals ...interface{}) {
	l.Log(zerolog.FatalLevel, msg, keyvals...)
}

func (l *Logger) Panic(msg string, keyvals ...interface{}) {
	l.Log(zerolog.PanicLevel, msg, keyvals...)
}

func pLog(level zerolog.Level, msg string, keyvals ...interface{}) {
	switch level {
	case zerolog.DebugLevel:
		log.Debug().Msgf(msg, keyvalsToFields(keyvals...))
	case zerolog.InfoLevel:
		log.Info().Msgf(msg, keyvalsToFields(keyvals...))
	case zerolog.WarnLevel:
		log.Warn().Msgf(msg, keyvalsToFields(keyvals...))
	case zerolog.ErrorLevel:
		log.Error().Msgf(msg, keyvalsToFields(keyvals...))
	case zerolog.PanicLevel:
		log.Panic().Msgf(msg, keyvalsToFields(keyvals...))
	case zerolog.TraceLevel:
		log.Trace().Msgf(msg, keyvalsToFields(keyvals...))
	default:
		log.Log().Msgf(msg, keyvalsToFields(keyvals...))
	}
}
func Trace(msg string, keyvals ...interface{}) {
	pLog(zerolog.TraceLevel, msg, keyvals...)
}

func Debug(msg string, keyvals ...interface{}) {
	pLog(zerolog.DebugLevel, msg, keyvals...)
}

func Info(msg string, keyvals ...interface{}) {
	pLog(zerolog.InfoLevel, msg, keyvals...)
}

func Warn(msg string, keyvals ...interface{}) {
	pLog(zerolog.WarnLevel, msg, keyvals...)
}

func Error(msg string, keyvals ...interface{}) {
	pLog(zerolog.ErrorLevel, msg, keyvals...)
}

func Fatal(msg string, keyvals ...interface{}) {
	pLog(zerolog.FatalLevel, msg, keyvals...)
}

func Panic(msg string, keyvals ...interface{}) {
	pLog(zerolog.PanicLevel, msg, keyvals...)
}
