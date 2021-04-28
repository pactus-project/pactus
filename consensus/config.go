package consensus

import (
	"time"

	"github.com/zarbchain/zarb-go/errors"
)

type Config struct {
	QueryProposalTimeout  time.Duration
	ChangeProposerTimeout time.Duration
	ChangeProposerDelta   time.Duration
}

func DefaultConfig() *Config {
	return &Config{
		QueryProposalTimeout:  1 * time.Second,
		ChangeProposerTimeout: 6 * time.Second,
		ChangeProposerDelta:   2 * time.Second,
	}
}

func TestConfig() *Config {
	return &Config{
		QueryProposalTimeout:  200 * time.Millisecond,
		ChangeProposerTimeout: 1 * time.Second,
		ChangeProposerDelta:   200 * time.Millisecond,
	}
}

func (conf *Config) SanityCheck() error {
	if conf.QueryProposalTimeout <= 0 {
		return errors.Errorf(errors.ErrInvalidConfig, "QueryProposalTimeout can't be negative")
	}
	if conf.ChangeProposerTimeout <= 0 {
		return errors.Errorf(errors.ErrInvalidConfig, "ChangeProposerTimeout can't be negative")
	}
	if conf.ChangeProposerDelta <= 0 {
		return errors.Errorf(errors.ErrInvalidConfig, "ChangeProposerDelta can't be negative")
	}

	return nil
}

func (conf *Config) CalculateChangeProposerTimeout(round int) time.Duration {
	return time.Duration(
		conf.ChangeProposerTimeout.Milliseconds()+conf.ChangeProposerDelta.Milliseconds()*int64(round),
	) * time.Millisecond
}
