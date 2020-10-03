package config

import (
	"github.com/zarbchain/zarb-go/utils"
)

type StoreConfig struct {
	Path string
}

func DefaultStoreConfig() *StoreConfig {
	return &StoreConfig{
		Path: "data",
	}
}

func (conf *StoreConfig) BlockStorePath() string {
	return utils.MakeAbs(conf.Path + "/block.db")
}

func (conf *StoreConfig) TxStorePath() string {
	return utils.MakeAbs(conf.Path + "/tx.db")
}

func (conf *StoreConfig) AccountStorePath() string {
	return utils.MakeAbs(conf.Path + "/account.db")
}

func (conf *StoreConfig) ValidatorStorePath() string {
	return utils.MakeAbs(conf.Path + "/validator.db")
}

func (conf *StoreConfig) StateStorePath() string {
	return utils.MakeAbs(conf.Path + "/state.db")
}
