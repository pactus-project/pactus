package sync

import "time"

type Config struct {
	StartingTimeout  time.Duration
	HeartBeatTimeout time.Duration
	BlockPerMessage  int
}

func DefaultConfig() *Config {
	return &Config{
		StartingTimeout:  time.Second * 5,
		HeartBeatTimeout: time.Second * 10,
		BlockPerMessage:  10,
	}
}
