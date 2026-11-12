package grpc

import (
	"fmt"

	"github.com/pactus-project/pactus/util/htpasswd"
	"golang.org/x/crypto/bcrypt"
)

// Config defines parameters for the gRPC server.
type Config struct {
	Enable       bool   `toml:"enable"`
	EnableWallet bool   `toml:"enable_wallet"`
	Listen       string `toml:"listen"`
	BasicAuth    string `toml:"basic_auth"`
}

func DefaultConfig() *Config {
	return &Config{
		Enable: false,
		Listen: "",
	}
}

func (c *Config) BasicCheck() error {
	if c.BasicAuth != "" {
		_, passwordHash, err := htpasswd.ExtractBasicAuth(c.BasicAuth)
		if err != nil {
			return err
		}

		// A password that is not a bcrypt hash can never be matched:
		// the server would start with an authentication that rejects everyone.
		if _, err := bcrypt.Cost([]byte(passwordHash)); err != nil {
			return fmt.Errorf("invalid basic auth: the password must be a bcrypt hash: %w", err)
		}
	}

	return nil
}
