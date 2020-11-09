package sync

import "time"

type Config struct {
	StartingTimeout  time.Duration
	HeartBeatTimeout time.Duration
	BlockPerMessage  int
}

func DefaultConfig() *Config {
	return &Config{
		StartingTimeout:  time.Second * 10,
		HeartBeatTimeout: time.Second * 10,
		BlockPerMessage:  10,
	}
}

func TestConfig() *Config {
	return &Config{
		StartingTimeout:  time.Second * 1,
		HeartBeatTimeout: time.Second * 5,
		BlockPerMessage:  10,
	}
}
