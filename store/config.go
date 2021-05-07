package store

import (
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/util"
)

type Config struct {
	Path string `toml:"" comment:"Path contains database directory. Default is ./store.db"`
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

func (conf *Config) StorePath() string {
	return util.MakeAbs(conf.Path + "/store.db")
}

// SanityCheck is a basic checks for config
func (conf *Config) SanityCheck() error {
	if !util.IsValidDirPath(conf.Path) {
		return errors.Errorf(errors.ErrInvalidConfig, "path is not valid")
	}
	return nil
}
