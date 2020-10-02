package config

type CapnpConfig struct {
	Enable  bool
	Address string
}

func DefaultCapnpConfig() *CapnpConfig {
	return &CapnpConfig{
		Enable:  true,
		Address: "0.0.0.0:0",
	}
}
