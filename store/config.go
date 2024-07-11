package store

import (
	"path/filepath"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/util"
)

type Config struct {
	Path          string `toml:"path"`
	RetentionDays uint32 `toml:"retention_days"`

	// Private configs
	TxCacheSize        uint32                  `toml:"-"`
	SortitionCacheSize uint32                  `toml:"-"`
	AccountCacheSize   int                     `toml:"-"`
	PublicKeyCacheSize int                     `toml:"-"`
	BannedAddrs        map[crypto.Address]bool `toml:"-"`
}

func DefaultConfig() *Config {
	return &Config{
		Path:               "data",
		RetentionDays:      10,
		TxCacheSize:        1024,
		SortitionCacheSize: 1024,
		AccountCacheSize:   1024,
		PublicKeyCacheSize: 1024,
		BannedAddrs:        map[crypto.Address]bool{},
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

	if conf.RetentionDays < 10 {
		return ConfigError{
			Reason: "Retention days can't be less than 10 days",
		}
	}

	return nil
}

func (conf *Config) RetentionBlocks() uint32 {
	return conf.RetentionDays * 8640
}
