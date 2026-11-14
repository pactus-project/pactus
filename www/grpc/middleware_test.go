package grpc_test

import (
	"context"
	"encoding/base64"
	"testing"
	"time"

	pactusgrpc "github.com/pactus-project/pactus/www/grpc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// mockUnaryHandler simulates a gRPC method handler.
func mockUnaryHandler(_ context.Context, _ any) (any, error) {
	return "response", nil
}

func mockUnaryPanicHandler(_ context.Context, _ any) (any, error) {
	panic("panic happen!!!")
}

func TestBasicAuth(t *testing.T) {
	auth := "Basic " + base64.StdEncoding.EncodeToString([]byte("user:password"))
	invalidAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte("invalid:invalid"))
	malformedAuth := "Malformed"

	tests := []struct {
		name          string
		authHeader    string
		expectedError codes.Code
	}{
		{
			name:          "ValidCredentials",
			authHeader:    auth,
			expectedError: codes.OK,
		},
		{
			name:          "InvalidCredentials",
			authHeader:    invalidAuth,
			expectedError: codes.Unauthenticated,
		},
		{
			name:          "NoMetadata",
			authHeader:    "",
			expectedError: codes.Unauthenticated,
		},
		{
			name:          "MalformedAuthHeader",
			authHeader:    malformedAuth,
			expectedError: codes.Unauthenticated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := t.Context()
			if tt.authHeader != "" {
				md := metadata.New(map[string]string{"authorization": tt.authHeader})
				ctx = metadata.NewIncomingContext(ctx, md)
			}

			interceptor := pactusgrpc.BasicAuth("user:$2y$10$5Kjd955BDWLouqckHzBjKuCF6hFOUD61lhm8QpjDVHTUwMIrYUdq2")

			_, err := interceptor(ctx, nil, &grpc.UnaryServerInfo{}, mockUnaryHandler)

			got, want := status.Code(err), tt.expectedError
			assert.Equal(t, want, got)
		})
	}
}

func TestGrpcRecovery(t *testing.T) {
	td := setup(t, nil)

	interceptor := td.Server.Recovery()

	_, err := interceptor(t.Context(), nil, &grpc.UnaryServerInfo{}, mockUnaryPanicHandler)
	assert.Equal(t, codes.Unknown, status.Code(err))
}

func TestBasicAuthRemembersVerifiedCredential(t *testing.T) {
	interceptor := pactusgrpc.BasicAuth("user:$2y$10$5Kjd955BDWLouqckHzBjKuCF6hFOUD61lhm8QpjDVHTUwMIrYUdq2")

	call := func(user, password string) (time.Duration, codes.Code) {
		auth := "Basic " + base64.StdEncoding.EncodeToString([]byte(user+":"+password))
		ctx := metadata.NewIncomingContext(t.Context(), metadata.New(map[string]string{"authorization": auth}))

		start := time.Now()
		_, err := interceptor(ctx, nil, &grpc.UnaryServerInfo{}, mockUnaryHandler)

		return time.Since(start), status.Code(err)
	}

	first, code := call("user", "password")
	require.Equal(t, codes.OK, code)

	// The following requests of the same client don't pay the bcrypt comparison again.
	second, code := call("user", "password")
	require.Equal(t, codes.OK, code)
	assert.Less(t, second, first/4)

	// A wrong credential is still rejected, whatever was verified before.
	_, code = call("user", "wrong")
	assert.Equal(t, codes.Unauthenticated, code)

	_, code = call("other", "password")
	assert.Equal(t, codes.Unauthenticated, code)

	_, code = call("user", "password")
	assert.Equal(t, codes.OK, code)
}
