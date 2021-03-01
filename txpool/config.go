package txpool

import (
	"time"

	"github.com/zarbchain/zarb-go/errors"
)

type Config struct {
	WaitingTimeout time.Duration
	MaxSize        int
}

func DefaultConfig() *Config {
	return &Config{
		WaitingTimeout: 2 * time.Second,
		MaxSize:        2000,
	}
}

func TestConfig() *Config {
	return &Config{
		WaitingTimeout: 100 * time.Millisecond,
		MaxSize:        10,
	}
}

// SanityCheck is a basic checks for config
func (conf *Config) SanityCheck() error {
	if conf.WaitingTimeout < 0 {
		return errors.Errorf(errors.ErrInvalidConfig, "WaitingTimeout can't be negative")
	}
	if conf.MaxSize == 0 {
		return errors.Errorf(errors.ErrInvalidConfig, "MaxSize can't be negative or zero")
	}
	return nil
}
