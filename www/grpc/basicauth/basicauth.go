package basicauth

import (
	"context"
)

type BasicAuth struct {
	Username string
	Password string
}

// Credentials is a custom type implementing grpc.PerRPCCredentials.
type Credentials struct {
	Token string
}

// GetRequestMetadata gets the request metadata as a map of strings.
func (b Credentials) GetRequestMetadata(context.Context, ...string) (map[string]string, error) {
	return map[string]string{
		"authorization": b.Token,
	}, nil
}

// RequireTransportSecurity indicates whether the credentials requires transport security.
func (b Credentials) RequireTransportSecurity() bool {
	return false
}
