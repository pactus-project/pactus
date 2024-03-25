package basicauth

import (
	"fmt"
	"testing"
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
			// Call BasicAuth function
			result := EncodeBasicAuth(tc.username, tc.password)

			// Check if the result matches the expected output
			if result != tc.expected {
				t.Errorf("BasicAuth(%s, %s) = %s; want %s", tc.username, tc.password, result, tc.expected)
			}
		})
	}
}
