package logger

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type logger struct {
	config *Config
	subs   map[string]*SubLogger
}

var (
	globalInst *logger
)

type SubLogger struct {
	logger zerolog.Logger
	name   string
	obj    fmt.Stringer
}

func getLoggersInst() *logger {
	if globalInst == nil {
		// Only during tests the globalInst is nil

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
		if conf.Colorful {
			log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
		}

		lvl, err := zerolog.ParseLevel(conf.Levels["default"])
		if err == nil {
			zerolog.SetGlobalLevel(lvl)
		}
	}
}

func NewSubLogger(name string, obj fmt.Stringer) *SubLogger {
	sl := &SubLogger{
		logger: zerolog.New(os.Stderr).With().Timestamp().Logger(),
		name:   name,
		obj:    obj,
	}

	if getLoggersInst().config.Colorful {
		sl.logger = sl.logger.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	lvlStr := getLoggersInst().config.Levels[name]
	if lvlStr == "" {
		lvlStr = getLoggersInst().config.Levels["default"]
	}

	lvl, err := zerolog.ParseLevel(lvlStr)
	if err == nil {
		sl.logger.Level(lvl)
	}

	getLoggersInst().subs[name] = sl
	return sl
}

func (sl *SubLogger) checkEvent(event *zerolog.Event, msg string, keyvals ...interface{}) {
	if sl.obj != nil {
		event.Str(sl.name, sl.obj.String()).Fields(keyvals).Msg(msg)
	} else {
		event.Fields(keyvals).Msg(msg)
	}
}

func (sl *SubLogger) Trace(msg string, keyvals ...interface{}) {
	sl.checkEvent(sl.logger.Trace(), msg, keyvals...)
}

func (sl *SubLogger) Debug(msg string, keyvals ...interface{}) {
	sl.checkEvent(sl.logger.Debug(), msg, keyvals...)
}

func (sl *SubLogger) Info(msg string, keyvals ...interface{}) {
	sl.checkEvent(sl.logger.Info(), msg, keyvals...)
}

func (sl *SubLogger) Warn(msg string, keyvals ...interface{}) {
	sl.checkEvent(sl.logger.Warn(), msg, keyvals...)
}

func (sl *SubLogger) Error(msg string, keyvals ...interface{}) {
	sl.checkEvent(sl.logger.Error(), msg, keyvals...)
}

func (sl *SubLogger) Fatal(msg string, keyvals ...interface{}) {
	sl.checkEvent(sl.logger.Fatal(), msg, keyvals...)
}

func (sl *SubLogger) Panic(msg string, keyvals ...interface{}) {
	sl.checkEvent(sl.logger.Panic(), msg, keyvals...)
}

func Trace(msg string, keyvals ...interface{}) {
	log.Trace().Fields(keyvals).Msg(msg)
}

func Debug(msg string, keyvals ...interface{}) {
	log.Debug().Fields(keyvals).Msg(msg)
}

func Info(msg string, keyvals ...interface{}) {
	log.Info().Fields(keyvals).Msg(msg)
}

func Warn(msg string, keyvals ...interface{}) {
	log.Warn().Fields(keyvals).Msg(msg)
}

func Error(msg string, keyvals ...interface{}) {
	log.Error().Fields(keyvals).Msg(msg)
}

func Fatal(msg string, keyvals ...interface{}) {
	log.Fatal().Fields(keyvals).Msg(msg)
}

func Panic(msg string, keyvals ...interface{}) {
	log.Panic().Fields(keyvals).Msg(msg)
}
