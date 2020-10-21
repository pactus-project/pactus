package network

type Config struct {
	Name      string
	Address   string
	NodeKey   string
	Bootstrap []string
}

func DefaultConfig() *Config {
	return &Config{
		Name:    "zarb-testnet",
		Address: "/ip4/0.0.0.0/tcp/0",
		NodeKey: "node_key",
	}
}
