package txpool

import (
	"github.com/pactus-project/pactus/types/amount"
)

type Config struct {
	MaxSize     int     `toml:"max_size"`
	MinValuePAC float64 `toml:"min_value"`
}

func DefaultConfig() *Config {
	return &Config{
		MaxSize:     1000,
		MinValuePAC: 0.1,
	}
}

// BasicCheck performs basic checks on the configuration.
func (conf *Config) BasicCheck() error {
	if conf.MaxSize < 10 {
		return ConfigError{
			Reason: "maxSize can't be less than 10",
		}
	}

	if conf.MinValuePAC > 1 {
		return ConfigError{
			Reason: "minVale can't be greater than 1 PAC",
		}
	}

	return nil
}

func (conf *Config) minValue() amount.Amount {
	amt, _ := amount.NewAmount(conf.MinValuePAC)

	return amt
}

func (conf *Config) sortitionPoolSize() int {
	return int(float32(conf.MaxSize) * 0.1)
}

func (conf *Config) bondPoolSize() int {
	return int(float32(conf.MaxSize) * 0.1)
}

func (conf *Config) unbondPoolSize() int {
	return int(float32(conf.MaxSize) * 0.1)
}

func (conf *Config) withdrawPoolSize() int {
	return int(float32(conf.MaxSize) * 0.1)
}

func (conf *Config) transferPoolSize() int {
	return int(float32(conf.MaxSize) * 0.6)
}
