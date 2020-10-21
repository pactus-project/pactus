package capnp

type Config struct {
	Enable  bool
	Address string
}

func DefaultConfig() *Config {
	return &Config{
		Enable:  true,
		Address: "0.0.0.0:0",
	}
}
