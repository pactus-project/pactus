package consensus

import (
	"time"

	"github.com/zarbchain/zarb-go/util/errors"
)

type Config struct {
	ChangeProposerTimeout time.Duration `toml:"change_proposer_timeout"`
	ChangeProposerDelta   time.Duration `toml:"change_proposer_delta"`
	QueryProposalTimeout  time.Duration `toml:"query_proposal_timeout"`
}

func DefaultConfig() *Config {
	return &Config{
		QueryProposalTimeout:  1 * time.Second,
		ChangeProposerTimeout: 6 * time.Second,
		ChangeProposerDelta:   2 * time.Second,
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

func (conf *Config) CalculateChangeProposerTimeout(round int16) time.Duration {
	return time.Duration(
		conf.ChangeProposerTimeout.Milliseconds()+conf.ChangeProposerDelta.Milliseconds()*int64(round),
	) * time.Millisecond
}
