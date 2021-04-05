package consensus

import (
	"time"

	"github.com/zarbchain/zarb-go/errors"
)

type Config struct {
	Timeout time.Duration
}

func DefaultConfig() *Config {
	return &Config{
		Timeout: 5 * time.Second,
	}
}

func TestConfig() *Config {
	return &Config{
		Timeout: 500 * time.Millisecond,
	}
}

func (conf *Config) SanityCheck() error {
	if conf.Timeout < 0 {
		return errors.Errorf(errors.ErrInvalidConfig, "Timeout can't be negative")
	}

	return nil
}
