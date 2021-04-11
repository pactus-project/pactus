package store

import (
	"github.com/zarbchain/zarb-go/errors"
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

func TestConfig() *Config {
	return &Config{
		Path: util.TempDirPath(),
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

// SanityCheck is a basic checks for config
func (conf *Config) SanityCheck() error {
	if !util.IsValidDirPath(conf.Path) {
		return errors.Errorf(errors.ErrInvalidConfig, "Path is not valid")
	}
	return nil
}
