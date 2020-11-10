package state

import (
	"github.com/zarbchain/zarb-go/store"
)

type Config struct {
	Store *store.Config
}

func DefaultConfig() *Config {
	return &Config{
		Store: store.DefaultConfig(),
	}
}

func TestConfig() *Config {
	return &Config{
		Store: store.TestConfig(),
	}
}
