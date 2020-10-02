package config

type TxPoolConfig struct {
	MaxSize int
}

func DefaultTxPoolConfig() *TxPoolConfig {
	return &TxPoolConfig{
		MaxSize: 10000,
	}
}
