package txpool

import (
	"time"

	"github.com/zarbchain/zarb-go/errors"
)

type Config struct {
	WaitingTimeout time.Duration `toml:"WaitingTimeout" comment:"WaitingTimeout is block interval time. Default is 2 second."`
	MaxSize        int           `toml:"MaxSize" comment:"MaxSize descirbe the txpool maximum memory size."`
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
		return errors.Errorf(errors.ErrInvalidConfig, "waitingTimeout can't be negative")
	}
	if conf.MaxSize == 0 {
		return errors.Errorf(errors.ErrInvalidConfig, "maxSize can't be negative or zero")
	}
	return nil
}
