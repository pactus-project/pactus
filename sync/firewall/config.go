package firewall

type Config struct {
	Enabled bool `toml:"enable"`
}

func DefaultConfig() *Config {
	return &Config{
		Enabled: false,
	}
}

// BasicCheck performs basic checks on the configuration.
func (conf *Config) BasicCheck() error {
	return nil
}
