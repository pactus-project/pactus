package consensus

import (
	"time"

	"github.com/zarbchain/zarb-go/errors"
)

type Config struct {
	ChangeProposerTimeout time.Duration
	ChangeProposerDelta   time.Duration
}

func DefaultConfig() *Config {
	return &Config{
		ChangeProposerTimeout: 6 * time.Second,
		ChangeProposerDelta:   2 * time.Second,
	}
}

func TestConfig() *Config {
	return &Config{
		ChangeProposerTimeout: 1 * time.Second,
		ChangeProposerDelta:   200 * time.Millisecond,
	}
}

func (conf *Config) SanityCheck() error {
	if conf.ChangeProposerTimeout <= 0 {
		return errors.Errorf(errors.ErrInvalidConfig, "changeProposerTimeout can't be negative")
	}

	if conf.ChangeProposerDelta <= 0 {
		return errors.Errorf(errors.ErrInvalidConfig, "changeProposerDelta can't be negative")
	}

	return nil
}

func (conf *Config) CalculateChangeProposerTimeout(round int) time.Duration {
	return time.Duration(
		conf.ChangeProposerTimeout.Milliseconds()+conf.ChangeProposerDelta.Milliseconds()*int64(round),
	) * time.Millisecond
}
