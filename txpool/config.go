package txpool

import (
	"time"

	"github.com/zarbchain/zarb-go/errors"
)

type Config struct {
	MaxSize int `toml:"" comment:"Maximum number of unconfirmed transaction inside pool. Default is 2000"`
}

func DefaultConfig() *Config {
	return &Config{
		MaxSize: 2000,
	}
}

func TestConfig() *Config {
	return &Config{
		MaxSize: 20,
	}
}

// SanityCheck is a basic checks for config
func (conf *Config) SanityCheck() error {
	if conf.MaxSize == 0 {
		return errors.Errorf(errors.ErrInvalidConfig, "maxSize can't be negative or zero")
	}
	return nil
}

func (conf *Config) sortitionPoolSize() int {
	return int(float32(conf.MaxSize) * 0.05)
}

func (conf *Config) bondPoolSize() int {
	return int(float32(conf.MaxSize) * 0.05)
}

func (conf *Config) unbondPoolSize() int {
	return int(float32(conf.MaxSize) * 0.05)
}

func (conf *Config) withdrawPoolSize() int {
	return int(float32(conf.MaxSize) * 0.05)
}

func (conf *Config) sendPoolSize() int {
	return int(float32(conf.MaxSize) * 0.8)
}

func (conf *Config) queryTimeout() time.Duration {
	return time.Second * 2
}
