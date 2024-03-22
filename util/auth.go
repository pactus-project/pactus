package util

import (
	"encoding/base64"
	"fmt"
)

func BasicAuth(username, password string) string {
	return fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(username+":"+password)))
}
