package txpool

import (
	"github.com/pactus-project/pactus/types/amount"
)

type Config struct {
	MaxSize   int     `toml:"max_size"`
	MinFeePAC float64 `toml:"min_fee"`
}

func DefaultConfig() *Config {
	return &Config{
		MaxSize:   1000,
		MinFeePAC: 0.000001,
	}
}

// BasicCheck performs basic checks on the configuration.
func (conf *Config) BasicCheck() error {
	if conf.MaxSize < 10 {
		return ConfigError{
			Reason: "maxSize can't be less than 10",
		}
	}

	return nil
}

func (conf *Config) minFee() amount.Amount {
	amt, _ := amount.NewAmount(conf.MinFeePAC)

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
