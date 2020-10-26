package store

import (
	"github.com/zarbchain/zarb-go/util"
)

type Config struct {
	Path string
}

func DefaultConfig() *Config {
	return &Config{
		Path: "data",
	}
}

func (conf *Config) BlockStorePath() string {
	return util.MakeAbs(conf.Path + "/block.db")
}

func (conf *Config) TxStorePath() string {
	return util.MakeAbs(conf.Path + "/tx.db")
}

func (conf *Config) AccountStorePath() string {
	return util.MakeAbs(conf.Path + "/account.db")
}

func (conf *Config) ValidatorStorePath() string {
	return util.MakeAbs(conf.Path + "/validator.db")
}

func (conf *Config) StateStorePath() string {
	return util.MakeAbs(conf.Path + "/state.db")
}
