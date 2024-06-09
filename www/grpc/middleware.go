package grpc

import (
	"context"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/pactus-project/pactus/util/htpasswd"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
)

func BasicAuth(storedCredential string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (
		any, error,
	) {
		user, password, err := htpasswd.ExtractBasicAuthFromContext(ctx)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "failed to extract basic auth from header")
		}

		if err := htpasswd.CompareBasicAuth(storedCredential, user, password); err != nil {
			return nil, status.Error(codes.Unauthenticated, "username or password is invalid")
		}

		return handler(ctx, req)
	}
}

func Recovery() grpc.UnaryServerInterceptor {
	rec := func(p interface{}) (err error) {
		err = status.Errorf(codes.Unknown, "%v", p)
		log.Println("recovery: panic triggered in grpc server", "error", err)
		return
	}
	opts := []grpc_recovery.Option{
		grpc_recovery.WithRecoveryHandler(rec),
	}
	return grpc_recovery.UnaryServerInterceptor(opts...)
}
