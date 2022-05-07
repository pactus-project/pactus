package firewall

type Config struct {
	Enabled bool `toml:"enable"`
}

func DefaultConfig() *Config {
	return &Config{
		Enabled: false,
	}
}

// SanityCheck is a basic checks for config.
func (conf *Config) SanityCheck() error {
	return nil
}
