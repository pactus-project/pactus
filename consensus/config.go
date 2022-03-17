package consensus

import (
	"time"

	"github.com/zarbchain/zarb-go/errors"
)

type Config struct {
	QueryProposalTimeout  time.Duration `toml:"" comment:"QueryProposalTimeout which query the network if propsal does not exist.Default is 1 second."`
	ChangeProposerTimeout time.Duration `toml:"" comment:"ChangeProposerTimeout if current proposer failed to create the block .Default is 6 second."`
	ChangeProposerDelta   time.Duration `toml:"" comment:"ChangeProposerDelta which increase proposer timeout by round.Default is 2 second."`
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
		return errors.Errorf(errors.ErrInvalidConfig, "timeout for query proposal can't be negative")
	}
	if conf.ChangeProposerTimeout <= 0 {
		return errors.Errorf(errors.ErrInvalidConfig, "timeout for change proposer can't be negative")
	}
	if conf.ChangeProposerDelta <= 0 {
		return errors.Errorf(errors.ErrInvalidConfig, "change proposer delta can't be negative")
	}

	return nil
}

func (conf *Config) CalculateChangeProposerTimeout(round int) time.Duration {
	return time.Duration(
		conf.ChangeProposerTimeout.Milliseconds()+conf.ChangeProposerDelta.Milliseconds()*int64(round),
	) * time.Millisecond
}
