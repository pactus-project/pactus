package txpool

import (
	"github.com/pactus-project/pactus/types/amount"
)

type Config struct {
	MaxSize int        `toml:"max_size"`
	Fee     *FeeConfig `toml:"fee"`

	// Private configs
	ConsumptionWindow uint32 `toml:"-"`
}

type FeeConfig struct {
	FixedFee   float64 `toml:"fixed_fee"`
	DailyLimit uint32  `toml:"daily_limit"`
	UnitPrice  float64 `toml:"unit_price"`
}

func DefaultConfig() *Config {
	return &Config{
		MaxSize:           1000,
		Fee:               DefaultFeeConfig(),
		ConsumptionWindow: 8640,
	}
}

func DefaultFeeConfig() *FeeConfig {
	return &FeeConfig{
		FixedFee:   0.01,
		DailyLimit: 280,
		UnitPrice:  0,
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
	amt, _ := amount.NewAmount(conf.Fee.FixedFee)

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
