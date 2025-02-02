package config

// Modifier override some configuration after load configuration.
type Modifier func(cfg *Config)

func EnableGRPCTransport(listen string, enableWallet bool) Modifier {
	return func(cfg *Config) {
		if listen != "" {
			cfg.GRPC.Enable = true
			cfg.GRPC.Listen = listen
		}

		cfg.GRPC.EnableWallet = enableWallet
	}
}

func EnableZMQBlockInfoPub(listen string) Modifier {
	return func(cfg *Config) {
		if listen != "" {
			cfg.ZeroMq.ZmqPubBlockInfo = listen
		}
	}
}

func EnableZMQTxInfoPub(listen string) Modifier {
	return func(cfg *Config) {
		if listen != "" {
			cfg.ZeroMq.ZmqPubTxInfo = listen
		}
	}
}

func EnableZMQRawBlockPub(listen string) Modifier {
	return func(cfg *Config) {
		if listen != "" {
			cfg.ZeroMq.ZmqPubRawBlock = listen
		}
	}
}

func EnableZMQRawTxPub(listen string) Modifier {
	return func(cfg *Config) {
		if listen != "" {
			cfg.ZeroMq.ZmqPubRawTx = listen
		}
	}
}
