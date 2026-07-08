package integration

import (
	"testing"

	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNetworkInfo(t *testing.T) {
	res, err := tNetwork.GetNodeInfo(t.Context(), &pactus.GetNodeInfoRequest{})
	require.NoError(t, err)

	assert.Equal(t, "test", res.NetworkName)
}
