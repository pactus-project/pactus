package store

import (
	"path/filepath"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/util"
)

const blockPerDay = 8640

// XeggexAccount defines the required information to
// freeze Xeggex Deposit account based on PIP-38.
type XeggexAccount struct {
	DepositAddrs crypto.Address
	WatcherAddrs crypto.Address
	AccountHash  hash.Hash
	Balance      amount.Amount
	FreezeHeight uint32
}

type Config struct {
	Path          string `toml:"path"`
	RetentionDays uint32 `toml:"retention_days"`

	// Private configs
	TxCacheWindow      uint32                  `toml:"-"`
	SeedCacheWindow    uint32                  `toml:"-"`
	AccountCacheSize   int                     `toml:"-"`
	PublicKeyCacheSize int                     `toml:"-"`
	BannedAddrs        map[crypto.Address]bool `toml:"-"`
	XeggexAccount      XeggexAccount           `toml:"-"`
}

func DefaultConfig() *Config {
	xeggexDepositAddrs, _ := crypto.AddressFromString("pc1z2wtq43p8fnueya9qufq9hkutyr899npk2suleu")
	xeggexAccountHash, _ := hash.FromString("31b863dcb0bb82184fb7357ae54c1e50c8984f5641eae1aff7a7b1f39284b9f5")
	xeggexBalance, _ := amount.NewAmount(500_000)
	xeggexWatcherAddrs, _ := crypto.AddressFromString("pc1rqy07rwx7kdesnens3e5mc2ngk745q44wyndyc4")

	return &Config{
		Path:               "data",
		RetentionDays:      10,
		TxCacheWindow:      1024,
		SeedCacheWindow:    1024,
		AccountCacheSize:   1024,
		PublicKeyCacheSize: 1024,
		BannedAddrs:        map[crypto.Address]bool{},

		// https://pacviewer.com/address/pc1z2wtq43p8fnueya9qufq9hkutyr899npk2suleu
		XeggexAccount: XeggexAccount{
			DepositAddrs: xeggexDepositAddrs,
			WatcherAddrs: xeggexWatcherAddrs,
			Balance:      xeggexBalance,
			AccountHash:  xeggexAccountHash,
			FreezeHeight: 3_164_119,
		},
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

	if conf.TxCacheWindow == 0 ||
		conf.SeedCacheWindow == 0 {
		return ConfigError{
			Reason: "cache window set to zero",
		}
	}

	if conf.AccountCacheSize == 0 ||
		conf.PublicKeyCacheSize == 0 {
		return ConfigError{
			Reason: "cache size set to zero",
		}
	}

	if conf.RetentionDays < 10 {
		return ConfigError{
			Reason: "retention days can't be less than 10 days",
		}
	}

	return nil
}

func (conf *Config) RetentionBlocks() uint32 {
	return conf.RetentionDays * blockPerDay
}
