package consensus

import (
	"time"

	"github.com/zarbchain/zarb-go/errors"
)

type Config struct {
	TimeoutPropose   time.Duration
	TimeoutPrepare   time.Duration
	TimeoutPrecommit time.Duration
	DeltaDuration    time.Duration
}

func DefaultConfig() *Config {
	return &Config{
		TimeoutPropose:   3 * time.Second,
		TimeoutPrepare:   2 * time.Second,
		TimeoutPrecommit: 2 * time.Second,
		DeltaDuration:    1 * time.Second,
	}
}

func TestConfig() *Config {
	return &Config{
		TimeoutPropose:   300 * time.Millisecond,
		TimeoutPrepare:   200 * time.Millisecond,
		TimeoutPrecommit: 200 * time.Millisecond,
		DeltaDuration:    100 * time.Millisecond,
	}
}

func (conf *Config) SanityCheck() error {
	if conf.TimeoutPropose < 0 {
		return errors.Errorf(errors.ErrInvalidConfig, "TimeoutPropose can't be negative")
	}
	if conf.TimeoutPrepare < 0 {
		return errors.Errorf(errors.ErrInvalidConfig, "TimeoutPrepare can't be negative")
	}
	if conf.TimeoutPrecommit < 0 {
		return errors.Errorf(errors.ErrInvalidConfig, "TimeoutPrecommit can't be negative")
	}
	if conf.DeltaDuration < 0 {
		return errors.Errorf(errors.ErrInvalidConfig, "DeltaDuration can't be negative")
	}

	return nil
}

func (conf *Config) ProposeTimeout(round int) time.Duration {
	return time.Duration(
		conf.TimeoutPropose.Milliseconds()+conf.DeltaDuration.Milliseconds()*int64(round),
	) * time.Millisecond
}

func (conf *Config) PrepareTimeout(round int) time.Duration {
	return time.Duration(
		conf.TimeoutPrepare.Milliseconds()+conf.DeltaDuration.Milliseconds()*int64(round),
	) * time.Millisecond
}

func (conf *Config) PrecommitTimeout(round int) time.Duration {
	return time.Duration(
		conf.TimeoutPrecommit.Milliseconds()+conf.DeltaDuration.Milliseconds()*int64(round),
	) * time.Millisecond
}
