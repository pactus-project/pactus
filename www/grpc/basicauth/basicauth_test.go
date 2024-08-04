package basicauth

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeCredentials(t *testing.T) {
	testCases := []struct {
		username string
		password string
		expected string
	}{
		{"user", "pass", "Basic dXNlcjpwYXNz"},
		{"admin", "admin123", "Basic YWRtaW46YWRtaW4xMjM="},
		{"", "", "Basic Og=="},
		{"test", "123£", "Basic dGVzdDoxMjPCow=="},
	}

	// Iterate over test cases
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Username: %s, Password: %s", tc.username, tc.password), func(t *testing.T) {
			// Call basicAuth function
			result := EncodeBasicAuth(tc.username, tc.password)

			// Check if the result matches the expected output
			assert.Equal(t, tc.expected, result)
		})
	}
}
