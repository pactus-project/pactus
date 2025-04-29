package rest

type Config struct {
	Enable     bool   `toml:"enable"`
	Listen     string `toml:"listen"`
	EnableCORS bool   `toml:"enable_cors"`
}

func DefaultConfig() *Config {
	return &Config{
		Enable:     false,
		Listen:     "",
		EnableCORS: false,
	}
}

func (*Config) BasicCheck() error {
	return nil
}
