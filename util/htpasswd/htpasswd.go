package htpasswd

import (
	"context"
	"encoding/base64"
	"errors"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/metadata"
)

const (
	// passwordSeparator defines the separator used in basic auth credentials (username:password).
	passwordSeparator = ":"
)

var (
	ErrInvalidUser             = errors.New("user is invalid")
	ErrFailedToParseBasicAuth  = errors.New("the provided basic authentication credentials are invalid")
	ErrMetadataNotFound        = errors.New("metadata not found")
	ErrAuthHeaderNotFound      = errors.New("authorization header not found")
	ErrFailedToDecodeBasicAuth = errors.New("failed to decode authorization header")
	ErrAuthHeaderInvalidFormat = errors.New("invalid authorization header format")
	ErrInvalidPassword         = errors.New("password is invalid")
)

// CompareBasicAuth compares a stored credential (username:password_hash) with a provided username and password.
// It uses bcrypt to securely compare the password hash stored in the credential with the provided password.
func CompareBasicAuth(storedCredential, user, password string) error {
	storedUser, storedPasswordHash, err := ExtractBasicAuth(storedCredential)
	if err != nil {
		return err
	}

	if storedUser != user {
		return ErrInvalidUser
	}

	if err := bcrypt.CompareHashAndPassword([]byte(storedPasswordHash), []byte(password)); err != nil {
		return ErrInvalidPassword
	}

	return nil
}

// ExtractBasicAuth extracts the user and password or password hash from the given basic auth credential.
// The credential should be in the form "user:password" or "user:password_hahs".
func ExtractBasicAuth(basicAuthCredential string) (user, password string, err error) {
	parts := strings.SplitN(basicAuthCredential, passwordSeparator, 2)
	if len(parts) != 2 {
		return "", "", ErrFailedToParseBasicAuth
	}

	user = parts[0]
	password = parts[1]

	return user, password, nil
}

// ExtractBasicAuthFromContext extracts the user and password from the incoming context in gRPC request.
func ExtractBasicAuthFromContext(ctx context.Context) (user, password string, err error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", "", ErrMetadataNotFound
	}

	authHeader, ok := md["authorization"]
	if !ok || len(authHeader) == 0 {
		return "", "", ErrAuthHeaderNotFound
	}

	auth := strings.TrimPrefix(authHeader[0], "Basic ")
	decoded, err := base64.StdEncoding.DecodeString(auth)
	if err != nil {
		return "", "", ErrFailedToDecodeBasicAuth
	}

	parts := strings.SplitN(string(decoded), ":", 2)
	if len(parts) != 2 {
		return "", "", ErrAuthHeaderInvalidFormat
	}

	return parts[0], parts[1], nil
}
