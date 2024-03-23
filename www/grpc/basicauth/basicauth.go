package basicauth

import (
	"context"
	"encoding/base64"
	"fmt"
)

// EncodeBasicAuth generates a Basic Authentication header value using the provided username and password
// according to RFC 7617. It formats the authString as "username:password", encodes it in base64,
// and returns the Basic Auth header in the format "Basic <base64EncodedString>".
func EncodeBasicAuth(username, password string) string {
	authString := fmt.Sprintf("%s:%s", username, password)
	encodedAuth := base64.StdEncoding.EncodeToString([]byte(authString))

	return fmt.Sprintf("Basic %s", encodedAuth)
}

// BasicAuth is an implementation of grpc.PerRPCCredentials based on the Basic HTTP Authentication Schema.
type BasicAuth struct {
	Username string
	Password string
}

func New(username, password string) *BasicAuth {
	return &BasicAuth{
		Username: username,
		Password: password,
	}
}

// GetRequestMetadata gets the request metadata as a map of strings.
func (b *BasicAuth) GetRequestMetadata(context.Context, ...string) (map[string]string, error) {
	return map[string]string{
		"authorization": EncodeBasicAuth(b.Username, b.Password),
	}, nil
}

// RequireTransportSecurity indicates whether the credentials requires transport security.
func (b *BasicAuth) RequireTransportSecurity() bool {
	return false
}
