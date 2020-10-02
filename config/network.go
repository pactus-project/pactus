package config

type NetworkConfig struct {
	Name      string
	Address   string
	NodeKey   string
	Bootstrap []string
}

func DefaultNetworkConfig() *NetworkConfig {
	return &NetworkConfig{
		Name:    "zarb-testnet",
		Address: "/ip4/0.0.0.0/tcp/0",
		NodeKey: "node_key",
	}
}
