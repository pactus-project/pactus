package cmd

import "context"

var Auth = &auth{}

type auth struct {
	Username string
	Password string
}

// BasicAuthCredentials is a custom type implementing grpc.PerRPCCredentials.
type basicAuthCredentials struct {
	Token string
}

// GetRequestMetadata gets the request metadata as a map of strings.
func (b basicAuthCredentials) GetRequestMetadata(context.Context, ...string) (map[string]string, error) {
	return map[string]string{
		"authorization": b.Token,
	}, nil
}

// RequireTransportSecurity indicates whether the credentials requires transport security.
func (b basicAuthCredentials) RequireTransportSecurity() bool {
	return false
}
