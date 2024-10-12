package txpool

import (
	"github.com/pactus-project/pactus/types/amount"
)

type Config struct {
	MaxSize int       `toml:"max_size"`
	Fee     FeeConfig `toml:"fee"`
}

type FeeConfig struct {
	DailyLimit uint32  `toml:"daily_limit"`
	UnitPrice  float64 `toml:"unit_price"`
	FixedPrice float64 `toml:"fixed_price"`
}

func DefaultConfig() *Config {
	return &Config{
		MaxSize: 1000,
		Fee:     DefaultFeeConfig(),
	}
}

func DefaultFeeConfig() FeeConfig {
	return FeeConfig{
		DailyLimit: 280,
		UnitPrice:  0,
		FixedPrice: 0.01,
	}
}

// BasicCheck performs basic checks on the configuration.
func (conf *Config) BasicCheck() error {
	if conf.MaxSize < 10 {
		return ConfigError{
			Reason: "maxSize can't be less than 10",
		}
	}

	if conf.Fee.DailyLimit == 0 {
		return ConfigError{
			Reason: "dailyLimit can't be zero",
		}
	}

	return nil
}

func (conf *Config) minFee() amount.Amount {
	amt, _ := amount.NewAmount(conf.Fee.FixedPrice)

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
