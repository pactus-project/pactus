package grpc_test

import (
	"testing"

	pactusgrpc "github.com/pactus-project/pactus/www/grpc"
	"github.com/stretchr/testify/require"
)

func TestConfigBasicCheck(t *testing.T) {
	t.Run("the default config is valid", func(t *testing.T) {
		require.NoError(t, pactusgrpc.DefaultConfig().BasicCheck())
	})

	t.Run("a user with a bcrypt hash is valid", func(t *testing.T) {
		conf := pactusgrpc.DefaultConfig()
		conf.BasicAuth = "user:$2y$10$q6I6fxG2c79jBSXJ8L2jde15czipSRpu/uhW5Le.ooJLyfXiaPDZG"

		require.NoError(t, conf.BasicCheck())
	})

	t.Run("a basic auth without a password is not valid", func(t *testing.T) {
		conf := pactusgrpc.DefaultConfig()
		conf.BasicAuth = "user-without-password"

		require.Error(t, conf.BasicCheck())
	})
}
