package consensus

import (
	"time"

	"github.com/zarbchain/zarb-go/errors"
)

type Config struct {
	ChangeProposerTimeout time.Duration
}

func DefaultConfig() *Config {
	return &Config{
		ChangeProposerTimeout: 1 * time.Second,
	}
}

func TestConfig() *Config {
	return &Config{
		ChangeProposerTimeout: 1 * time.Second,
	}
}

func (conf *Config) SanityCheck() error {
	if conf.ChangeProposerTimeout < 0 {
		return errors.Errorf(errors.ErrInvalidConfig, "ChangeProposerTimeout can't be negative")
	}

	return nil
}
