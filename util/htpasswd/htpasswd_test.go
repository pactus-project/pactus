package htpasswd

import (
	"context"
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/metadata"
)

func TestParseHtpasswdAuth(t *testing.T) {
	tests := []struct {
		input           string
		user            string
		encodedPassword string
	}{
		{
			input:           "user:$2y$10$q6I6fxG2c79jBSXJ8L2jde15czipSRpu/uhW5Le.ooJLyfXiaPDZG",
			user:            "user",
			encodedPassword: "$2y$10$q6I6fxG2c79jBSXJ8L2jde15czipSRpu/uhW5Le.ooJLyfXiaPDZG",
		},
		{
			input:           "user1:$2y$10$/4EcZtrJUgivhcTJPGOz/uhQEUAQP.zvThFwIHwdjQT97iL4gWMri",
			user:            "user1",
			encodedPassword: "$2y$10$/4EcZtrJUgivhcTJPGOz/uhQEUAQP.zvThFwIHwdjQT97iL4gWMri",
		},
		{
			input:           "user2:$2y$10$xXmx6BQv6re3P2sOAoPGNu/MJOwWxDtxtNzlEJ2qkUVRK6SqAXD9m",
			user:            "user2",
			encodedPassword: "$2y$10$xXmx6BQv6re3P2sOAoPGNu/MJOwWxDtxtNzlEJ2qkUVRK6SqAXD9m",
		},
		{
			input:           "user3:$2y$10$eKLWzld7iMPrcyDqam8.Y.R1deeSUBWFD3P6eQHJ0Iqa1qR4yBxaq",
			user:            "user3",
			encodedPassword: "$2y$10$eKLWzld7iMPrcyDqam8.Y.R1deeSUBWFD3P6eQHJ0Iqa1qR4yBxaq",
		},
	}

	for _, tt := range tests {
		t.Run(tt.user, func(t *testing.T) {
			user, encodedPass, err := ParseHtpasswdAuth(tt.input)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, user+":"+encodedPass, tt.input)
		})
	}
}

func TestCompareBasicAuth(t *testing.T) {
	tests := []struct {
		input    string
		user     string
		password string
	}{
		{
			input:    "user:$2y$10$q6I6fxG2c79jBSXJ8L2jde15czipSRpu/uhW5Le.ooJLyfXiaPDZG",
			user:     "user",
			password: "foobar",
		},
		{
			input:    "user1:$2y$05$y9dWO1FBS34D7RSZSNZ6S.NjE3LMNBvSAwidgTrER/AHBNN9cBeR.",
			user:     "user1",
			password: "foobar1",
		},
		{
			input:    "user2:$2y$11$RuWzAY2N57m.iZuT9bUh2ufOj2nNd02BviZSVx2Hbid8PvonjPWRi",
			user:     "user2",
			password: "foobar2",
		},
		{
			input:    "user3:$2y$09$866UNklDooeXGSd6MI/XPu1Fg9.2nTX6dFnPsEdgtBY6HMF5.NhPq",
			user:     "user3",
			password: "foobar3",
		},
	}

	for _, tt := range tests {
		if err := CompareBasicAuth(tt.input, tt.user, tt.password); err != nil {
			t.Error(err)
		}
	}
}

func BenchmarkParseHtpasswd(b *testing.B) {
	auth := []string{
		"user:$2y$10$q6I6fxG2c79jBSXJ8L2jde15czipSRpu/uhW5Le.ooJLyfXiaPDZG",
		"user1:$2y$05$y9dWO1FBS34D7RSZSNZ6S.NjE3LMNBvSAwidgTrER/AHBNN9cBeR.",
		"user2:$2y$11$RuWzAY2N57m.iZuT9bUh2ufOj2nNd02BviZSVx2Hbid8PvonjPWRi",
		"user3:$2y$09$866UNklDooeXGSd6MI/XPu1Fg9.2nTX6dFnPsEdgtBY6HMF5.NhPq",
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for _, a := range auth {
			_, _, err := ParseHtpasswdAuth(a)
			if err != nil {
				b.Fatal(err)
			}
		}
	}
}

func BenchmarkCompareBasicAuth(b *testing.B) {
	tests := []struct {
		input    string
		user     string
		password string
	}{
		{
			input:    "user:$2y$10$q6I6fxG2c79jBSXJ8L2jde15czipSRpu/uhW5Le.ooJLyfXiaPDZG",
			user:     "user",
			password: "foobar",
		},
		{
			input:    "user1:$2y$05$y9dWO1FBS34D7RSZSNZ6S.NjE3LMNBvSAwidgTrER/AHBNN9cBeR.",
			user:     "user1",
			password: "foobar1",
		},
		{
			input:    "user2:$2y$11$RuWzAY2N57m.iZuT9bUh2ufOj2nNd02BviZSVx2Hbid8PvonjPWRi",
			user:     "user2",
			password: "foobar2",
		},
		{
			input:    "user3:$2y$09$866UNklDooeXGSd6MI/XPu1Fg9.2nTX6dFnPsEdgtBY6HMF5.NhPq",
			user:     "user3",
			password: "foobar3",
		},
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for _, tt := range tests {
			_ = CompareBasicAuth(tt.input, tt.user, tt.password)
		}
	}
}

func TestExtractBasicAuthFromContext(t *testing.T) {
	tests := []struct {
		name     string
		ctx      context.Context
		wantUser string
		wantPass string
		wantErr  bool
	}{
		{
			name:     "Valid credentials",
			ctx:      createTestContext("Basic " + base64.StdEncoding.EncodeToString([]byte("user:password"))),
			wantUser: "user",
			wantPass: "password",
			wantErr:  false,
		},
		{
			name:     "Invalid encoding",
			ctx:      createTestContext("Basic user:password"),
			wantUser: "",
			wantPass: "",
			wantErr:  true,
		},
		{
			name:     "No authorization header",
			ctx:      metadata.NewIncomingContext(context.Background(), metadata.MD{}),
			wantUser: "",
			wantPass: "",
			wantErr:  true,
		},
		{
			name:     "Incorrect format",
			ctx:      createTestContext("Basic " + base64.StdEncoding.EncodeToString([]byte("userpassword"))),
			wantUser: "",
			wantPass: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, pass, err := ExtractBasicAuthFromContext(tt.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExtractBasicAuthFromContext() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if user != tt.wantUser || pass != tt.wantPass {
				t.Errorf("ExtractBasicAuthFromContext() got = %v, %v, want %v, %v", user, pass, tt.wantUser, tt.wantPass)
			}
		})
	}
}

func createTestContext(authValue string) context.Context {
	md := metadata.New(map[string]string{"authorization": authValue})

	return metadata.NewIncomingContext(context.Background(), md)
}
