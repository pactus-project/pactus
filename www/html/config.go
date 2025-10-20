package html

// Config defines parameters for the HTML UI server.
type Config struct {
	Enable      bool   `toml:"enable"`
	Listen      string `toml:"listen"`
	EnablePprof bool   `toml:"enable_pprof"`
}

func DefaultConfig() *Config {
	return &Config{
		Enable:      false,
		Listen:      "",
		EnablePprof: false,
	}
}

// BasicCheck performs basic checks on the configuration.
func (*Config) BasicCheck() error {
	return nil
}
