package grpc

import (
	"context"
	"encoding/base64"
	"testing"

	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func TestHealthCheck(t *testing.T) {
	td := setup(t, nil)
	client := healthpb.NewHealthClient(td.newClient(t))

	tests := []struct {
		name    string
		service string
	}{
		{
			name:    "OverallServer",
			service: "",
		},
		{
			name:    "BlockchainService",
			service: pactus.Blockchain_ServiceDesc.ServiceName,
		},
		{
			name:    "TransactionService",
			service: pactus.Transaction_ServiceDesc.ServiceName,
		},
		{
			name:    "NetworkService",
			service: pactus.Network_ServiceDesc.ServiceName,
		},
		{
			name:    "UtilsService",
			service: pactus.Utils_ServiceDesc.ServiceName,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := client.Check(t.Context(), &healthpb.HealthCheckRequest{
				Service: tt.service,
			})
			require.NoError(t, err)
			require.Equal(t, healthpb.HealthCheckResponse_SERVING, resp.GetStatus())
		})
	}
}

func TestHealthWatch(t *testing.T) {
	td := setup(t, nil)
	client := healthpb.NewHealthClient(td.newClient(t))

	stream, err := client.Watch(t.Context(), &healthpb.HealthCheckRequest{
		Service: pactus.Blockchain_ServiceDesc.ServiceName,
	})
	require.NoError(t, err)

	resp, err := stream.Recv()
	require.NoError(t, err)
	require.Equal(t, healthpb.HealthCheckResponse_SERVING, resp.GetStatus())
}

func TestHealthCheckWalletService(t *testing.T) {
	conf := testConfig()
	conf.EnableWallet = true

	td := setup(t, conf)
	client := healthpb.NewHealthClient(td.newClient(t))

	resp, err := client.Check(t.Context(), &healthpb.HealthCheckRequest{
		Service: pactus.Wallet_ServiceDesc.ServiceName,
	})
	require.NoError(t, err)
	require.Equal(t, healthpb.HealthCheckResponse_SERVING, resp.GetStatus())
}

func TestHealthCheckBasicAuth(t *testing.T) {
	conf := testConfig()
	conf.BasicAuth = "user:$2y$10$5Kjd955BDWLouqckHzBjKuCF6hFOUD61lhm8QpjDVHTUwMIrYUdq2"

	td := setup(t, conf)
	client := healthpb.NewHealthClient(td.newClient(t))

	_, err := client.Check(t.Context(), &healthpb.HealthCheckRequest{})
	require.Equal(t, codes.Unauthenticated, status.Code(err))

	ctx := contextWithBasicAuth(t.Context(), "user", "password")
	resp, err := client.Check(ctx, &healthpb.HealthCheckRequest{})
	require.NoError(t, err)
	require.Equal(t, healthpb.HealthCheckResponse_SERVING, resp.GetStatus())
}

func TestHealthWatchBasicAuth(t *testing.T) {
	conf := testConfig()
	conf.BasicAuth = "user:$2y$10$5Kjd955BDWLouqckHzBjKuCF6hFOUD61lhm8QpjDVHTUwMIrYUdq2"

	td := setup(t, conf)
	client := healthpb.NewHealthClient(td.newClient(t))

	stream, err := client.Watch(t.Context(), &healthpb.HealthCheckRequest{})
	require.NoError(t, err)
	_, err = stream.Recv()
	require.Equal(t, codes.Unauthenticated, status.Code(err))

	ctx := contextWithBasicAuth(t.Context(), "user", "password")
	stream, err = client.Watch(ctx, &healthpb.HealthCheckRequest{})
	require.NoError(t, err)

	resp, err := stream.Recv()
	require.NoError(t, err)
	require.Equal(t, healthpb.HealthCheckResponse_SERVING, resp.GetStatus())
}

func contextWithBasicAuth(ctx context.Context, user, password string) context.Context {
	auth := "Basic " + base64.StdEncoding.EncodeToString([]byte(user+":"+password))

	return metadata.NewOutgoingContext(ctx, metadata.Pairs("authorization", auth))
}
