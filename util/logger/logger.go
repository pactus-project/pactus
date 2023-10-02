package logger

import (
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"reflect"

	"github.com/pactus-project/pactus/util"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	LogFilename = "pactus.log"
	MaxLogSize  = 10 // 10MB to rotate a log file
)

var globalInst *logger

type logger struct {
	config *Config
	subs   map[string]*SubLogger
}

type SubLogger struct {
	logger zerolog.Logger
	name   string
	obj    fmt.Stringer
}

func getLoggersInst() *logger {
	if globalInst == nil {
		// Only during tests the globalInst is nil

		LogFilename = util.TempFilePath()

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
		globalInst = &logger{
			config: conf,
			subs:   make(map[string]*SubLogger),
		}
	}

	return globalInst
}

func InitGlobalLogger(conf *Config) {
	if globalInst == nil {
		globalInst = &logger{
			config: conf,
			subs:   make(map[string]*SubLogger),
		}
		log.Logger = zerolog.New(globalInst.writers()).With().Timestamp().Logger()

		lvl, err := zerolog.ParseLevel(conf.Levels["default"])
		if err == nil {
			zerolog.SetGlobalLevel(lvl)
		}
	}
}

func addFields(event *zerolog.Event, keyvals ...interface{}) *zerolog.Event {
	if len(keyvals)%2 != 0 {
		keyvals = append(keyvals, "!MISSING-VALUE!")
	}
	for i := 0; i < len(keyvals); i += 2 {
		key, ok := keyvals[i].(string)
		if !ok {
			key = "!INVALID-KEY!"
		}
		///
		value := keyvals[i+1]
		switch v := value.(type) {
		case fmt.Stringer:
			if isNil(v) {
				event.Any(key, v)
			} else {
				event.Stringer(key, v)
			}
		case error:
			event.AnErr(key, v)
		case []byte:
			event.Str(key, fmt.Sprintf("%v", hex.EncodeToString(v)))
		default:
			event.Any(key, v)
		}
	}
	return event
}

func NewSubLogger(name string, obj fmt.Stringer) *SubLogger {
	inst := getLoggersInst()

	sl := &SubLogger{
		logger: zerolog.New(inst.writers()).With().Timestamp().Logger(),
		name:   name,
		obj:    obj,
	}

	lvlStr := inst.config.Levels[name]
	if lvlStr == "" {
		lvlStr = inst.config.Levels["default"]
	}

	lvl, err := zerolog.ParseLevel(lvlStr)
	if err == nil {
		sl.logger = sl.logger.Level(lvl)
	}

	inst.subs[name] = sl
	return sl
}

func (l *logger) writers() io.Writer {
	writers := []io.Writer{}
	// console writer
	if l.config.Colorful {
		writers = append(writers, zerolog.ConsoleWriter{Out: os.Stderr})
	} else {
		writers = append(writers, os.Stderr)
	}

	// file writer
	fw := &lumberjack.Logger{
		Filename:   LogFilename,
		MaxSize:    MaxLogSize,
		MaxBackups: 0,
		Compress:   true,
	}
	writers = append(writers, fw)

	return io.MultiWriter(writers...)
}

func (sl *SubLogger) logObj(event *zerolog.Event, msg string, keyvals ...interface{}) {
	if sl.obj != nil {
		addFields(event.Str(sl.name, sl.obj.String()), keyvals...).Msg(msg)
	} else {
		addFields(event, keyvals...).Msg(msg)
	}
}

func (sl *SubLogger) Trace(msg string, keyvals ...interface{}) {
	sl.logObj(sl.logger.Trace(), msg, keyvals...)
}

func (sl *SubLogger) Debug(msg string, keyvals ...interface{}) {
	sl.logObj(sl.logger.Debug(), msg, keyvals...)
}

func (sl *SubLogger) Info(msg string, keyvals ...interface{}) {
	sl.logObj(sl.logger.Info(), msg, keyvals...)
}

func (sl *SubLogger) Warn(msg string, keyvals ...interface{}) {
	sl.logObj(sl.logger.Warn(), msg, keyvals...)
}

func (sl *SubLogger) Error(msg string, keyvals ...interface{}) {
	sl.logObj(sl.logger.Error(), msg, keyvals...)
}

func (sl *SubLogger) Fatal(msg string, keyvals ...interface{}) {
	sl.logObj(sl.logger.Fatal(), msg, keyvals...)
}

func (sl *SubLogger) Panic(msg string, keyvals ...interface{}) {
	sl.logObj(sl.logger.Panic(), msg, keyvals...)
}

func Trace(msg string, keyvals ...interface{}) {
	addFields(log.Trace(), keyvals...).Msg(msg)
}

func Debug(msg string, keyvals ...interface{}) {
	addFields(log.Debug(), keyvals...).Msg(msg)
}

func Info(msg string, keyvals ...interface{}) {
	addFields(log.Info(), keyvals...).Msg(msg)
}

func Warn(msg string, keyvals ...interface{}) {
	addFields(log.Warn(), keyvals...).Msg(msg)
}

func Error(msg string, keyvals ...interface{}) {
	addFields(log.Error(), keyvals...).Msg(msg)
}

func Fatal(msg string, keyvals ...interface{}) {
	addFields(log.Fatal(), keyvals...).Msg(msg)
}

func Panic(msg string, keyvals ...interface{}) {
	addFields(log.Panic(), keyvals...).Msg(msg)
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
