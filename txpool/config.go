package txpool

import (
	"time"

	"github.com/zarbchain/zarb-go/errors"
)

type Config struct {
	WaitingTimeout time.Duration `toml:"" comment:"Query and validate transaction wait-timeout. Default is 2s"`
	MaxSize        int           `toml:"" comment:"Maximum number of unconfirmed transaction inside txpool.Default is 2000"`
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
