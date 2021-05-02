package firewall

type Config struct {
	Enabled bool
}

func DefaultConfig() *Config {
	return &Config{
		Enabled: false,
	}
}

func TestConfig() *Config {
	return &Config{
		Enabled: false,
	}
}

// SanityCheck is a basic checks for config
func (conf *Config) SanityCheck() error {
	return nil
}
