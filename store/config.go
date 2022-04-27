package store

import (
	"fmt"
	"os"

	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/util"
)

type Config struct {
	Path string `toml:"path"`
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

func (conf *Config) DataPath() string {
	return util.MakeAbs(conf.Path)
}

func (conf *Config) StorePath() string {
	return fmt.Sprintf("%s%c%s", conf.DataPath(), os.PathSeparator, "store.db")
}

// SanityCheck is a basic checks for config
func (conf *Config) SanityCheck() error {
	if !util.IsValidDirPath(conf.Path) {
		return errors.Errorf(errors.ErrInvalidConfig, "path is not valid")
	}
	return nil
}
