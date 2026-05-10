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
			return nil, err
		}

		return handler(ctx, req)
	}
}

func BasicAuthStream(storedCredential string) grpc.StreamServerInterceptor {
	return func(
		srv any, stream grpc.ServerStream, _ *grpc.StreamServerInfo, handler grpc.StreamHandler,
	) error {
		if err := checkBasicAuth(stream.Context(), storedCredential); err != nil {
			return err
		}

		return handler(srv, stream)
	}
}

func checkBasicAuth(ctx context.Context, storedCredential string) error {
	user, password, err := htpasswd.ExtractBasicAuthFromContext(ctx)
	if err != nil {
		return status.Error(codes.Unauthenticated, "failed to extract basic auth from header")
	}

	if err := htpasswd.CompareBasicAuth(storedCredential, user, password); err != nil {
		return status.Error(codes.Unauthenticated, "username or password is invalid")
	}

	return nil
}

func (s *Server) Recovery() grpc.UnaryServerInterceptor {
	opts := []rec.Option{
		rec.WithRecoveryHandler(s.recoverGRPCPanic),
	}

	return rec.UnaryServerInterceptor(opts...)
}

func (s *Server) RecoveryStream() grpc.StreamServerInterceptor {
	opts := []rec.Option{
		rec.WithRecoveryHandler(s.recoverGRPCPanic),
	}

	return rec.StreamServerInterceptor(opts...)
}

func (s *Server) recoverGRPCPanic(p any) (err error) {
	err = status.Errorf(codes.Unknown, "%v", p)
	stackTrace := debug.Stack()
	s.logger.Error(
		"recovery panic triggered in grpc server",
		"error", err,
		"stacktrace", string(stackTrace),
	)

	return err
}
