package basicauth

import (
	"encoding/base64"
	"fmt"
)

func MakeCredentials(username, password string) string {
	authString := fmt.Sprintf("%s:%s", username, password)
	encodedAuth := base64.StdEncoding.EncodeToString([]byte(authString))

	return fmt.Sprintf("Basic %s", encodedAuth)
}
