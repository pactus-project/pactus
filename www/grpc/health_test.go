package grpc

import (
	"testing"

	"github.com/stretchr/testify/require"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

func TestHealthCheck(t *testing.T) {
	td := setup(t, nil)
	client := healthpb.NewHealthClient(td.newClient(t))

	resp, err := client.Check(t.Context(), &healthpb.HealthCheckRequest{})
	require.NoError(t, err)
	require.Equal(t, healthpb.HealthCheckResponse_SERVING, resp.Status)
}
