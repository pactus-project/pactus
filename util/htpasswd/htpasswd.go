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

// CompareBasicAuth compare basic auth with bcrypt algorithm.
func CompareBasicAuth(basicAuthCredential, user, password string) error {
	parsedUser, parsedHashedPass, err := ParseHtpasswdAuth(basicAuthCredential)
	if err != nil {
		return err
	}

	if parsedUser != user {
		return ErrInvalidUser
	}

	if err := bcrypt.CompareHashAndPassword([]byte(parsedHashedPass), []byte(password)); err != nil {
		return ErrInvalidPassword
	}

	return nil
}

// ParseHtpasswdAuth parse htpasswd auth.
func ParseHtpasswdAuth(basicAuthCredential string) (string, string, error) {
	parts := strings.SplitN(basicAuthCredential, passwordSeparator, 2)
	if len(parts) != 2 {
		return "", "", ErrFailedToParseBasicAuth
	}

	user := parts[0]
	encodedPassword := parts[1]

	return user, encodedPassword, nil
}

// ExtractBasicAuthFromContext extract basic auth from incoming context in grpc request.
func ExtractBasicAuthFromContext(ctx context.Context) (string, string, error) {
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
