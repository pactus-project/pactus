package txpool

type Config struct {
	MaxSize int
}

func DefaultConfig() *Config {
	return &Config{
		MaxSize: 10000,
	}
}
