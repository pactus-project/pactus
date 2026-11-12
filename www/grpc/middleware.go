package grpc

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"fmt"
	"runtime/debug"
	"sync/atomic"

	rec "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/pactus-project/pactus/util/htpasswd"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func BasicAuth(storedCredential string) grpc.UnaryServerInterceptor {
	verifier := newCredentialVerifier(storedCredential)

	return func(ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (
		any, error,
	) {
		user, password, err := htpasswd.ExtractBasicAuthFromContext(ctx)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "failed to extract basic auth from header")
		}

		if !verifier.verify(user, password) {
			return nil, status.Error(codes.Unauthenticated, "username or password is invalid")
		}

		return handler(ctx, req)
	}
}

// credentialVerifier checks the basic auth credential of the requests.
//
// A bcrypt comparison takes tens of milliseconds by design.
// The verifier remembers the credential that passed it,
// so the following requests of the same client don't pay it again.
// It keeps a keyed digest of the credential, and the key never leaves the process.
type credentialVerifier struct {
	storedCredential string
	key              [32]byte
	verified         atomic.Pointer[[sha256.Size]byte]
}

func newCredentialVerifier(storedCredential string) *credentialVerifier {
	verifier := &credentialVerifier{storedCredential: storedCredential}
	if _, err := rand.Read(verifier.key[:]); err != nil {
		panic(err) // unreachable: crypto/rand never fails
	}

	return verifier
}

func (v *credentialVerifier) verify(user, password string) bool {
	mac := hmac.New(sha256.New, v.key[:])
	_, _ = fmt.Fprintf(mac, "%d:%s:%s", len(user), user, password)
	var digest [sha256.Size]byte
	mac.Sum(digest[:0])

	if verified := v.verified.Load(); verified != nil &&
		subtle.ConstantTimeCompare(verified[:], digest[:]) == 1 {
		return true
	}

	if err := htpasswd.CompareBasicAuth(v.storedCredential, user, password); err != nil {
		return false
	}
	v.verified.Store(&digest)

	return true
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
