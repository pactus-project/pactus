package consensus

import (
	"time"

	"github.com/zarbchain/zarb-go/errors"
)

type Config struct {
	TimeoutPropose          time.Duration
	TimeoutPrevote          time.Duration
	TimeoutPrecommit        time.Duration
	NewRoundDeltaDuration   time.Duration
	PeerGossipSleepDuration time.Duration
	FuzzTesting             bool
}

func DefaultConfig() *Config {
	return &Config{
		TimeoutPropose:          3 * time.Second,
		TimeoutPrevote:          2 * time.Second,
		TimeoutPrecommit:        2 * time.Second,
		NewRoundDeltaDuration:   1 * time.Second,
		PeerGossipSleepDuration: 100 * time.Millisecond,
		FuzzTesting:             false,
	}
}

func (conf *Config) SanityCheck() error {
	if conf.TimeoutPropose < 0 {
		return errors.Errorf(errors.ErrInvalidConfig, "timeout_propose can't be negative")
	}
	if conf.TimeoutPrevote < 0 {
		return errors.Errorf(errors.ErrInvalidConfig, "timeout_prevote can't be negative")
	}
	if conf.TimeoutPrecommit < 0 {
		return errors.Errorf(errors.ErrInvalidConfig, "timeout_precommit can't be negative")
	}
	if conf.NewRoundDeltaDuration < 0 {
		return errors.Errorf(errors.ErrInvalidConfig, "new_round_delta_duration can't be negative")
	}
	if conf.PeerGossipSleepDuration < 0 {
		return errors.Errorf(errors.ErrInvalidConfig, "peer_gossip_sleep_duration can't be negative")
	}

	return nil
}

func (conf *Config) Propose(round int) time.Duration {
	return time.Duration(
		conf.TimeoutPropose.Milliseconds()+conf.NewRoundDeltaDuration.Milliseconds()*int64(round),
	) * time.Millisecond
}

func (conf *Config) Prevote(round int) time.Duration {
	return time.Duration(
		conf.TimeoutPrevote.Milliseconds()+conf.NewRoundDeltaDuration.Milliseconds()*int64(round),
	) * time.Millisecond
}

func (conf *Config) Precommit(round int) time.Duration {
	return time.Duration(
		conf.TimeoutPrecommit.Milliseconds()+conf.NewRoundDeltaDuration.Milliseconds()*int64(round),
	) * time.Millisecond
}
