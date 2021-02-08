package consensus

import (
	"time"

	"github.com/zarbchain/zarb-go/errors"
)

type Config struct {
	TimeoutPrepare   time.Duration
	TimeoutPrecommit time.Duration
	DeltaDuration    time.Duration
}

func DefaultConfig() *Config {
	return &Config{
		TimeoutPrepare:   3 * time.Second,
		TimeoutPrecommit: 2 * time.Second,
		DeltaDuration:    1 * time.Second,
	}
}

func TestConfig() *Config {
	return &Config{
		TimeoutPrepare:   300 * time.Millisecond,
		TimeoutPrecommit: 200 * time.Millisecond,
		DeltaDuration:    100 * time.Millisecond,
	}
}

func (conf *Config) SanityCheck() error {
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
