package http

import (
	"fmt"
	"strings"
)

type Config struct {
	Enable     bool   `toml:"enable"`
	Listen     string `toml:"listen"`
	BasePath   string `toml:"base_path"`
	EnableCORS bool   `toml:"enable_cors"`
}

func DefaultConfig() *Config {
	return &Config{
		Enable:     false,
		Listen:     "",
		BasePath:   "/http",
		EnableCORS: false,
	}
}

func (*Config) BasicCheck() error {
	return nil
}

func (c *Config) swaggerPattern() string {
	return fmt.Sprintf("%sui/", c.rootPattern())
}

func (c *Config) apiPattern() string {
	return fmt.Sprintf("%sapi/", c.rootPattern())
}

func (c *Config) rootPattern() string {
	path := fmt.Sprintf("/%s/", c.BasePath)
	path = strings.ReplaceAll(path, "//", "/")
	path = strings.ReplaceAll(path, "//", "/")

	return path
}
