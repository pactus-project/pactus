package grpc

import (
	"context"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// mockUnaryHandler simulates a gRPC method handler.
func mockUnaryHandler(_ context.Context, _ interface{}) (interface{}, error) {
	return "response", nil
}

func TestBasicAuth(t *testing.T) {
	auth := "Basic dXNlcjpwYXNzd29yZA=="        // This is base64 for "user:password"
	invalidAuth := "Basic aW52YWxpZDppbnZhbGlk" // This is base64 for "invalid:invalid"

	tests := []struct {
		name          string
		authHeader    string
		expectedError codes.Code
	}{
		{
			name:          "Valid credentials",
			authHeader:    auth,
			expectedError: codes.OK,
		},
		{
			name:          "Invalid credentials",
			authHeader:    invalidAuth,
			expectedError: codes.Unauthenticated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			if tt.authHeader != "" {
				md := metadata.New(map[string]string{"authorization": tt.authHeader})
				ctx = metadata.NewIncomingContext(ctx, md)
			}

			interceptor := BasicAuth("user:$2y$10$5Kjd955BDWLouqckHzBjKuCF6hFOUD61lhm8QpjDVHTUwMIrYUdq2")

			_, err := interceptor(ctx, nil, &grpc.UnaryServerInfo{}, mockUnaryHandler)

			if got, want := status.Code(err), tt.expectedError; got != want {
				t.Errorf("expected error code %v, got %v", want, got)
			}
		})
	}
}
