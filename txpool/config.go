package txpool

import "time"

type Config struct {
	WaitingTimeout time.Duration
	MaxSize        int
}

func DefaultConfig() *Config {
	return &Config{
		WaitingTimeout: 2 * time.Second,
		MaxSize:        10000,
	}
}

func TestConfig() *Config {
	return &Config{
		WaitingTimeout: 1 * time.Second,
		MaxSize:        10000,
	}
}
