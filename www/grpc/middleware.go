package grpc

import (
	"context"
	"runtime/debug"

	rec "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/pactus-project/pactus/util/htpasswd"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func (s *Server) Recovery() grpc.UnaryServerInterceptor {
	recovery := func(p any) (err error) {
		err = status.Errorf(codes.Unknown, "%v", p)
		stackTrace := debug.Stack()
		s.logger.Error(
			"recovery panic triggered in grpc server",
			"error", err,
			"stacktrace", string(stackTrace),
		)

		return err
	}
	opts := []rec.Option{
		rec.WithRecoveryHandler(recovery),
	}

	return rec.UnaryServerInterceptor(opts...)
}
