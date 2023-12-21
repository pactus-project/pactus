package store

import (
	"fmt"
	"os"

	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/errors"
)

type Config struct {
	Path string `toml:"path"`

	// Private configs
	TransactionToLiveInterval uint32 `toml:"-"`
	SortitionInterval         uint32 `toml:"-"`
	AccountCacheSize          int    `toml:"-"`
	PublicKeyCacheSize        int    `toml:"-"`
}

func DefaultConfig() *Config {
	return &Config{
		Path:                      "data",
		TransactionToLiveInterval: 8640,
		SortitionInterval:         17,
		AccountCacheSize:          1024,
		PublicKeyCacheSize:        1024,
	}
}

func (conf *Config) DataPath() string {
	return util.MakeAbs(conf.Path)
}

func (conf *Config) StorePath() string {
	return fmt.Sprintf("%s%c%s", conf.DataPath(), os.PathSeparator, "store.db")
}

// BasicCheck performs basic checks on the configuration.
func (conf *Config) BasicCheck() error {
	if !util.IsValidDirPath(conf.Path) {
		return errors.Errorf(errors.ErrInvalidConfig, "path is not valid")
	}
	return nil
}
