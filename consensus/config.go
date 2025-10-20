package consensus

import "time"

// Config defines parameters for the legacy consensus algorithm.
type Config struct {
	ChangeProposerTimeout    time.Duration `toml:"-"`
	ChangeProposerDelta      time.Duration `toml:"-"`
	QueryVoteTimeout         time.Duration `toml:"-"`
	MinimumAvailabilityScore float64       `toml:"-"`
}

func DefaultConfig() *Config {
	return &Config{
		ChangeProposerTimeout:    5 * time.Second,
		ChangeProposerDelta:      5 * time.Second,
		QueryVoteTimeout:         5 * time.Second,
		MinimumAvailabilityScore: 0.666667,
	}
}

// BasicCheck performs basic checks on the configuration.
func (conf *Config) BasicCheck() error {
	if conf.ChangeProposerTimeout <= 0 {
		return ConfigError{
			Reason: "change proposer timeout must be greater than zero",
		}
	}
	if conf.ChangeProposerDelta <= 0 {
		return ConfigError{
			Reason: "change proposer delta must be greater than zero",
		}
	}
	if conf.MinimumAvailabilityScore < 0 || conf.MinimumAvailabilityScore > 1 {
		return ConfigError{
			Reason: "minimum availability score can't be negative or more than 1",
		}
	}

	return nil
}

func (conf *Config) CalculateChangeProposerTimeout(round int16) time.Duration {
	return conf.ChangeProposerTimeout +
		conf.ChangeProposerDelta*time.Duration(round)
}
