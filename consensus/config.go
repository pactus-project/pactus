package consensus

import (
	"time"

	"github.com/pactus-project/pactus/util/errors"
)

type Config struct {
	ChangeProposerTimeout time.Duration `toml:"change_proposer_timeout"`
	ChangeProposerDelta   time.Duration `toml:"change_proposer_delta"`
}

func DefaultConfig() *Config {
	return &Config{
		ChangeProposerTimeout: 6 * time.Second,
		ChangeProposerDelta:   2 * time.Second,
	}
}

// BasicCheck performs basic checks on the configuration.
func (conf *Config) BasicCheck() error {
	if conf.ChangeProposerTimeout <= 0 {
		return errors.Errorf(errors.ErrInvalidConfig, "timeout for change proposer can't be negative")
	}
	if conf.ChangeProposerDelta <= 0 {
		return errors.Errorf(errors.ErrInvalidConfig, "change proposer delta can't be negative")
	}

	return nil
}

func (conf *Config) CalculateChangeProposerTimeout(round int16) time.Duration {
	return time.Duration(
		conf.ChangeProposerTimeout.Milliseconds()+conf.ChangeProposerDelta.Milliseconds()*int64(round),
	) * time.Millisecond
}
