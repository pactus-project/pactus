package wallet

import "time"

type walletOpt struct {
	timeout time.Duration
	servers []string
}

type Option func(*walletOpt)

var defaultWalletOpt = &walletOpt{
	timeout: 5 * time.Second,
	servers: make([]string, 0),
}

func WithTimeout(timeout time.Duration) Option {
	return func(opt *walletOpt) {
		opt.timeout = timeout
	}
}

func WithCustomServers(servers []string) Option {
	return func(opt *walletOpt) {
		opt.servers = servers
	}
}
