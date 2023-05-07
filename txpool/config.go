package txpool

import (
	"github.com/pactus-project/pactus/util/errors"
)

type Config struct {
	MaxSize int `toml:"max_size"`
}

func DefaultConfig() *Config {
	return &Config{
		MaxSize: 2000,
	}
}

// SanityCheck performs basic checks on the configuration.
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
