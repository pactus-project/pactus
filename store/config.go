package store

import (
	"path/filepath"

	"github.com/pactus-project/pactus/util"
)

type Config struct {
	Path string `toml:"path"`

	// Private configs
	TxCacheSize        uint32 `toml:"-"`
	SortitionCacheSize uint32 `toml:"-"`
	AccountCacheSize   int    `toml:"-"`
	PublicKeyCacheSize int    `toml:"-"`
}

func DefaultConfig() *Config {
	return &Config{
		Path:               "data",
		TxCacheSize:        1024,
		SortitionCacheSize: 1024,
		AccountCacheSize:   1024,
		PublicKeyCacheSize: 1024,
	}
}

func (conf *Config) DataPath() string {
	return util.MakeAbs(conf.Path)
}

func (conf *Config) StorePath() string {
	return filepath.Join(conf.DataPath(), "store.db")
}

// BasicCheck performs basic checks on the configuration.
func (conf *Config) BasicCheck() error {
	if !util.IsValidDirPath(conf.Path) {
		return ConfigError{
			Reason: "path is not valid",
		}
	}

	if conf.TxCacheSize == 0 ||
		conf.SortitionCacheSize == 0 ||
		conf.AccountCacheSize == 0 ||
		conf.PublicKeyCacheSize == 0 {
		return ConfigError{
			Reason: "cache size set to zero",
		}
	}

	return nil
}
