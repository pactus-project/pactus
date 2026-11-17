package grpc

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfigBasicAuth(t *testing.T) {
	conf := DefaultConfig()
	require.NoError(t, conf.BasicCheck())

	conf.BasicAuth = "user:$2y$10$5Kjd955BDWLouqckHzBjKuCF6hFOUD61lhm8QpjDVHTUwMIrYUdq2"
	require.NoError(t, conf.BasicCheck())

	// The password is given as it is, not as a bcrypt hash: nobody could authenticate.
	conf.BasicAuth = "user:password"
	require.ErrorContains(t, conf.BasicCheck(), "bcrypt")

	conf.BasicAuth = "no-separator"
	require.Error(t, conf.BasicCheck())
}
