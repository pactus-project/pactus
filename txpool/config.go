package txpool

import "time"

type Config struct {
	WaitingTimeout time.Duration
	MaxSize        int
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

// SanityCheck is a basic hecks for config
func (conf *Config) SanityCheck() error {
	return nil
}
