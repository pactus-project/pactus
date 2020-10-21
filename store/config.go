package store

import (
	"github.com/zarbchain/zarb-go/utils"
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
	return utils.MakeAbs(conf.Path + "/block.db")
}

func (conf *Config) TxStorePath() string {
	return utils.MakeAbs(conf.Path + "/tx.db")
}

func (conf *Config) AccountStorePath() string {
	return utils.MakeAbs(conf.Path + "/account.db")
}

func (conf *Config) ValidatorStorePath() string {
	return utils.MakeAbs(conf.Path + "/validator.db")
}

func (conf *Config) StateStorePath() string {
	return utils.MakeAbs(conf.Path + "/state.db")
}
