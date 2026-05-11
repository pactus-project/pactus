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
		if err := checkBasicAuth(ctx, storedCredential); err != nil {
			return nil, status.Error(codes.Unauthenticated, "username or password is invalid")
		}

		return handler(ctx, req)
	}
}

func BasicAuthStream(storedCredential string) grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, _ *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		if err := checkBasicAuth(ss.Context(), storedCredential); err != nil {
			return status.Error(codes.Unauthenticated, "username or password is invalid")
		}

		return handler(srv, ss)
	}
}

func checkBasicAuth(ctx context.Context, storedCredential string) error {
	user, password, err := htpasswd.ExtractBasicAuthFromContext(ctx)
	if err != nil {
		return err
	}

	return htpasswd.CompareBasicAuth(storedCredential, user, password)
}

func (s *Server) Recovery() grpc.UnaryServerInterceptor {
	return rec.UnaryServerInterceptor(s.recoveryOptions()...)
}

func (s *Server) RecoveryStream() grpc.StreamServerInterceptor {
	return rec.StreamServerInterceptor(s.recoveryOptions()...)
}

func (s *Server) recoveryOptions() []rec.Option {
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

	return []rec.Option{
		rec.WithRecoveryHandler(recovery),
	}
}
