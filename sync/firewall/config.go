package firewall

type Config struct {
	Enabled bool `toml:"enable"`
}

func DefaultConfig() *Config {
	return &Config{
		Enabled: false,
	}
}

// SanityCheck performs basic checks on the configuration.
func (conf *Config) SanityCheck() error {
	return nil
}
