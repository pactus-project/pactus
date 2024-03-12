package grpc

import (
	"context"

	"github.com/pactus-project/pactus/util/htpasswd"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func BasicAuth(basicAuthCredential string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (
		any, error,
	) {
		user, password, err := htpasswd.ExtractBasicAuthFromContext(ctx)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "failed to extract basic auth from header")
		}

		if err := htpasswd.CompareBasicAuth(basicAuthCredential, user, password); err != nil {
			return nil, status.Error(codes.Unauthenticated, "username or password is invalid")
		}

		return handler(ctx, req)
	}
}
