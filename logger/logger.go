//go:build !debug
// +build !debug

package logger

import (
	"encoding/hex"
	"fmt"
	"reflect"

	"github.com/sirupsen/logrus"
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
	logger *logrus.Logger
	name   string
	obj    interface{}
}

func InitLogger(conf *Config) {
	if loggersInst == nil {
		loggersInst = &loggers{
			config:  conf,
			loggers: make(map[string]*Logger),
		}
		if conf.Colorfull {
			logrus.SetFormatter(&logrus.TextFormatter{ForceColors: true})
		} else {
			logrus.SetFormatter(&logrus.TextFormatter{DisableColors: true})
		}

		lvl, err := logrus.ParseLevel(conf.Levels["default"])
		if err == nil {
			logrus.SetLevel(lvl)
		}
	}
}

func NewLogger(name string, obj interface{}) *Logger {
	l := &Logger{
		logger: logrus.New(),
		name:   name,
		obj:    obj,
	}
	if loggersInst.config.Colorfull {
		l.logger.SetFormatter(&logrus.TextFormatter{ForceColors: true})

	} else {
		l.logger.SetFormatter(&logrus.TextFormatter{DisableColors: true})
	}

	lvl := loggersInst.config.Levels[name]
	if lvl == "" {
		lvl = loggersInst.config.Levels["default"]
	}

	level, err := logrus.ParseLevel(lvl)
	if err == nil {
		l.SetLevel(level)
	}

	loggersInst.loggers[name] = l
	return l
}

func isNil(i interface{}) bool {
	if i == nil {
		return true
	}
	switch reflect.TypeOf(i).Kind() {
	case reflect.Ptr, reflect.Map, reflect.Array, reflect.Chan, reflect.Slice:
		return reflect.ValueOf(i).IsNil()
	}
	return false
}

func keyvalsToFields(keyvals ...interface{}) logrus.Fields {
	if len(keyvals)%2 != 0 {
		keyvals = append(keyvals, "<MISSING VALUE>")
	}
	fields := make(logrus.Fields)
	for i := 0; i < len(keyvals); i += 2 {
		key, ok := keyvals[i].(string)
		if !ok {
			key = "invalid-key"
		}
		///
		val := "nil"
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
			val = fmt.Sprintf("%v", keyvals[i+1])
		}
		fields[key] = val
	}

	return fields
}

func (l *Logger) SetLevel(level logrus.Level) {
	l.logger.SetLevel(level)
}

func (l *Logger) withDefaultFields() *logrus.Entry {
	var fields logrus.Fields
	if f, ok := l.obj.(fingerprintable); ok {
		fields = keyvalsToFields(l.name, f.Fingerprint())
	} else {
		fields = keyvalsToFields("module", l.name)
	}
	return l.logger.WithFields(fields)
}

func (l *Logger) log(level logrus.Level, msg string, keyvals ...interface{}) {
	if l.logger.IsLevelEnabled(level) {
		l.withDefaultFields().WithFields(keyvalsToFields(keyvals...)).Log(level, msg)
	}
}

func (l *Logger) With(keyvals ...interface{}) *logrus.Entry {
	fields := keyvalsToFields(keyvals...)
	return l.logger.WithFields(fields)
}

func (l *Logger) Trace(msg string, keyvals ...interface{}) {
	l.log(logrus.TraceLevel, msg, keyvals...)
}

func (l *Logger) Debug(msg string, keyvals ...interface{}) {
	l.log(logrus.DebugLevel, msg, keyvals...)
}

func (l *Logger) Info(msg string, keyvals ...interface{}) {
	l.log(logrus.InfoLevel, msg, keyvals...)
}

func (l *Logger) Warn(msg string, keyvals ...interface{}) {
	l.log(logrus.WarnLevel, msg, keyvals...)
}

func (l *Logger) Error(msg string, keyvals ...interface{}) {
	l.log(logrus.ErrorLevel, msg, keyvals...)
}

func (l *Logger) Fatal(msg string, keyvals ...interface{}) {
	l.log(logrus.FatalLevel, msg, keyvals...)
}

func (l *Logger) Panic(msg string, keyvals ...interface{}) {
	l.log(logrus.PanicLevel, msg, keyvals...)
}

//---
func log(level logrus.Level, msg string, keyvals ...interface{}) {
	if logrus.IsLevelEnabled(level) {
		logrus.WithFields(keyvalsToFields(keyvals...)).Log(level, msg)
	}
}

func Trace(msg string, keyvals ...interface{}) {
	log(logrus.TraceLevel, msg, keyvals...)
}

func Debug(msg string, keyvals ...interface{}) {
	log(logrus.DebugLevel, msg, keyvals...)
}

func Info(msg string, keyvals ...interface{}) {
	log(logrus.InfoLevel, msg, keyvals...)
}

func Warn(msg string, keyvals ...interface{}) {
	log(logrus.WarnLevel, msg, keyvals...)
}

func Error(msg string, keyvals ...interface{}) {
	log(logrus.ErrorLevel, msg, keyvals...)
}

func Fatal(msg string, keyvals ...interface{}) {
	log(logrus.FatalLevel, msg, keyvals...)
}

func Panic(msg string, keyvals ...interface{}) {
	log(logrus.PanicLevel, msg, keyvals...)
}
