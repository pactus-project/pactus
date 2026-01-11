package logger

import (
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"reflect"
	"slices"

	"github.com/pactus-project/pactus/util"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/natefinch/lumberjack.v2"
)

type LogStringer interface {
	LogString() string
}

var (
	LogFilename = "pactus.log"
	MaxLogSize  = 10 // 10MB to rotate a log file
)

var globalInst *logger

type logger struct {
	config *Config
	subs   map[string]*SubLogger
	writer io.Writer
}

type SubLogger struct {
	logger zerolog.Logger
	name   string
	obj    LogStringer
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
		conf.Levels["_zmq"] = "debug"
		conf.Levels["_firewall"] = "debug"
		globalInst = &logger{
			config: conf,
			subs:   make(map[string]*SubLogger),
			writer: zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "15:04:05"},
		}
		log.Logger = zerolog.New(globalInst.writer).With().Timestamp().Logger()
	}

	return globalInst
}

func InitGlobalLogger(conf *Config) {
	if globalInst != nil {
		return
	}

	writers := []io.Writer{}

	if slices.Contains(conf.Targets, "file") {
		fileWriter := &lumberjack.Logger{
			Filename:   LogFilename,
			MaxSize:    MaxLogSize,
			MaxBackups: conf.MaxBackups,
			Compress:   conf.Compress,
			MaxAge:     conf.RotateLogAfterDays,
		}
		writers = append(writers, fileWriter)
	}

	if slices.Contains(conf.Targets, "console") {
		// console writer
		if conf.Colorful {
			consoleWriter := &zerolog.ConsoleWriter{
				Out:        os.Stderr,
				TimeFormat: "15:04:05",
			}
			writers = append(writers, consoleWriter)
		} else {
			writers = append(writers, os.Stderr)
		}
	}

	globalInst = &logger{
		config: conf,
		subs:   make(map[string]*SubLogger),
		writer: io.MultiWriter(writers...),
	}
	log.Logger = zerolog.New(globalInst.writer).With().Timestamp().Logger()

	lvl, err := zerolog.ParseLevel(conf.Levels["default"])
	if err != nil {
		Warn("invalid default log level", "error", err)
	}
	log.Logger = log.Logger.Level(lvl)
}

func addFields(event *zerolog.Event, keyvals ...any) *zerolog.Event {
	if event == nil {
		return nil
	}

	if len(keyvals)%2 != 0 {
		keyvals = append(keyvals, "!MISSING-VALUE!")
	}
	for index := 0; index < len(keyvals); index += 2 {
		key, ok := keyvals[index].(string)
		if !ok {
			key = "!INVALID-KEY!"
		}

		value := keyvals[index+1]
		switch typ := value.(type) {
		case LogStringer:
			if isNil(typ) {
				event.Any(key, typ)
			} else {
				event.Str(key, typ.LogString())
			}
		case fmt.Stringer:
			if isNil(typ) {
				event.Any(key, typ)
			} else {
				event.Stringer(key, typ)
			}
		case error:
			event.AnErr(key, typ)
		case []byte:
			event.Str(key, hex.EncodeToString(typ))
		default:
			event.Any(key, typ)
		}
	}

	return event
}

func NewSubLogger(name string, obj LogStringer) *SubLogger {
	inst := getLoggersInst()
	sub := &SubLogger{
		logger: zerolog.New(inst.writer).With().Timestamp().Logger(),
		name:   name,
		obj:    obj,
	}

	lvlStr := inst.config.Levels[name]
	if lvlStr == "" {
		lvlStr = inst.config.Levels["default"]
	}

	lvl, err := zerolog.ParseLevel(lvlStr)
	if err != nil {
		Warn("invalid log level", "error", err, "name", name)
	}
	sub.logger = sub.logger.Level(lvl)

	inst.subs[name] = sub

	return sub
}

func (sl *SubLogger) logObj(event *zerolog.Event, msg string, keyvals ...any) {
	if event == nil {
		return
	}

	if sl.obj != nil {
		event = event.Str(sl.name, sl.obj.LogString())
	}

	addFields(event, keyvals...).Msg(msg)
}

func (sl *SubLogger) Trace(msg string, keyvals ...any) {
	sl.logObj(sl.logger.Trace(), msg, keyvals...)
}

func (sl *SubLogger) Debug(msg string, keyvals ...any) {
	sl.logObj(sl.logger.Debug(), msg, keyvals...)
}

func (sl *SubLogger) Info(msg string, keyvals ...any) {
	sl.logObj(sl.logger.Info(), msg, keyvals...)
}

func (sl *SubLogger) Warn(msg string, keyvals ...any) {
	sl.logObj(sl.logger.Warn(), msg, keyvals...)
}

func (sl *SubLogger) Error(msg string, keyvals ...any) {
	sl.logObj(sl.logger.Error(), msg, keyvals...)
}

func (sl *SubLogger) Fatal(msg string, keyvals ...any) {
	sl.logObj(sl.logger.Fatal(), msg, keyvals...)
}

func (sl *SubLogger) Panic(msg string, keyvals ...any) {
	sl.logObj(sl.logger.Panic(), msg, keyvals...)
}

func Trace(msg string, keyvals ...any) {
	addFields(log.Trace(), keyvals...).Msg(msg)
}

func Debug(msg string, keyvals ...any) {
	addFields(log.Debug(), keyvals...).Msg(msg)
}

func Info(msg string, keyvals ...any) {
	addFields(log.Info(), keyvals...).Msg(msg)
}

func Warn(msg string, keyvals ...any) {
	addFields(log.Warn(), keyvals...).Msg(msg)
}

func Error(msg string, keyvals ...any) {
	addFields(log.Error(), keyvals...).Msg(msg)
}

func Fatal(msg string, keyvals ...any) {
	addFields(log.Fatal(), keyvals...).Msg(msg)
}

func Panic(msg string, keyvals ...any) {
	addFields(log.Panic(), keyvals...).Msg(msg)
}

func isNil(val any) bool {
	if val == nil {
		return true
	}
	if reflect.TypeOf(val).Kind() == reflect.Ptr {
		return reflect.ValueOf(val).IsNil()
	}

	return false
}

func (sl *SubLogger) SetObj(obj LogStringer) {
	sl.obj = obj
}
