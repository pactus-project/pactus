package consensus

import (
	"time"

	"github.com/pactus-project/pactus/util/errors"
)

type Config struct {
	ChangeProposerTimeout    time.Duration `toml:"-"`
	ChangeProposerDelta      time.Duration `toml:"-"`
	MinimumAvailabilityScore float64       `toml:"-"`
}

func DefaultConfig() *Config {
	return &Config{
		ChangeProposerTimeout:    8 * time.Second,
		ChangeProposerDelta:      4 * time.Second,
		MinimumAvailabilityScore: 0.8,
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
	if conf.MinimumAvailabilityScore < 0 || conf.MinimumAvailabilityScore > 1 {
		return errors.Errorf(errors.ErrInvalidConfig, "minimum availability score can't be negative or more than 1")
	}

	return nil
}

func (conf *Config) CalculateChangeProposerTimeout(round int16) time.Duration {
	return time.Duration(
		conf.ChangeProposerTimeout.Milliseconds()+conf.ChangeProposerDelta.Milliseconds()*int64(round),
	) * time.Millisecond
}
