package txpool

import "github.com/pactus-project/pactus/types/amount"

// Config defines parameters for the transaction pool.
type Config struct {
	MaxSize int        `toml:"max_size"`
	Fee     *FeeConfig `toml:"fee"`

	// Private configuration
	ConsumptionWindow uint32 `toml:"-"`
}

// FeeConfig holds fee-related settings used to estimate and validate
// transaction fees.
type FeeConfig struct {
	FixedFee   float64 `toml:"fixed_fee"`
	DailyLimit int     `toml:"daily_limit"`
	UnitPrice  float64 `toml:"unit_price"`
}

// DefaultConfig returns the default transaction pool configuration.
func DefaultConfig() *Config {
	return &Config{
		MaxSize:           1000,
		Fee:               DefaultFeeConfig(),
		ConsumptionWindow: 8640,
	}
}

// DefaultFeeConfig returns the default fee configuration.
func DefaultFeeConfig() *FeeConfig {
	return &FeeConfig{
		FixedFee:   0.01,
		DailyLimit: 360,
		UnitPrice:  0,
	}
}

// BasicCheck performs basic checks on the configuration.
func (conf *Config) BasicCheck() error {
	if conf.MaxSize < 10 {
		return ConfigError{
			Reason: "maxSize cannot be less than 10",
		}
	}

	if conf.Fee.DailyLimit <= 0 {
		return ConfigError{
			Reason: "dailyLimit must be positive",
		}
	}

	return nil
}

func (conf *Config) calculateConsumption() bool {
	return conf.Fee.UnitPrice > 0
}

func (conf *Config) fixedFee() amount.Amount {
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
	return int(float32(conf.MaxSize) * 0.5)
}

func (conf *Config) batchTransferPoolSize() int {
	return int(float32(conf.MaxSize) * 0.1)
}
