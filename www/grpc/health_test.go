package grpc

import (
	"testing"

	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
	"github.com/stretchr/testify/assert"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

func TestHealthCheck(t *testing.T) {
	td := setup(t, testConfig())
	client := td.healthClient(t)

	tests := []struct {
		name    string
		service string
	}{
		{
			name:    "overall server status",
			service: "",
		},
		{
			name:    "blockchain service status",
			service: pactus.Blockchain_ServiceDesc.ServiceName,
		},
		{
			name:    "network service status",
			service: pactus.Network_ServiceDesc.ServiceName,
		},
		{
			name:    "transaction service status",
			service: pactus.Transaction_ServiceDesc.ServiceName,
		},
		{
			name:    "utils service status",
			service: pactus.Utils_ServiceDesc.ServiceName,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := client.Check(t.Context(), &healthpb.HealthCheckRequest{
				Service: tt.service,
			})

			assert.NoError(t, err)
			assert.Equal(t, healthpb.HealthCheckResponse_SERVING, res.Status)
		})
	}
}

func TestHealthCheckWalletService(t *testing.T) {
	conf := testConfig()
	conf.EnableWallet = true

	td := setup(t, conf)
	client := td.healthClient(t)

	res, err := client.Check(t.Context(), &healthpb.HealthCheckRequest{
		Service: pactus.Wallet_ServiceDesc.ServiceName,
	})

	assert.NoError(t, err)
	assert.Equal(t, healthpb.HealthCheckResponse_SERVING, res.Status)
}
