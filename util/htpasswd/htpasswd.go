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

// CompareBasicAuth compare basic auth with bcrypt algorithm
func CompareBasicAuth(basicAuthCredential, user, password string) error {
	parsedUser, parsedHashedPass, err := ParseHtpasswdAuth(basicAuthCredential)
	if err != nil {
		return err
	}

	if parsedUser != user {
		return errors.New("user is invalid")
	}

	return bcrypt.CompareHashAndPassword([]byte(parsedHashedPass), []byte(password))
}

// ParseHtpasswdAuth parse htpasswd auth
func ParseHtpasswdAuth(auth string) (user, encodedPassword string, err error) {
	parts := strings.SplitN(auth, passwordSeparator, 2)
	if len(parts) != 2 {
		return "", "", errors.New("auth is invalid for parse")
	}

	user = parts[0]
	encodedPassword = parts[1]
	return
}

// ExtractBasicAuthFromContext extract basic auth from incoming context in grpc request
func ExtractBasicAuthFromContext(ctx context.Context) (user, password string, err error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", "", errors.New("metadata not found")
	}

	authHeader, ok := md["authorization"]
	if !ok || len(authHeader) == 0 {
		return "", "", errors.New("authorization header not found")
	}

	auth := strings.TrimPrefix(authHeader[0], "Basic ")
	decoded, err := base64.StdEncoding.DecodeString(auth)
	if err != nil {
		return "", "", errors.New("failed to decode authorization header")
	}

	parts := strings.SplitN(string(decoded), ":", 2)
	if len(parts) != 2 {
		return "", "", errors.New("invalid authorization header format")
	}

	return parts[0], parts[1], nil
}
