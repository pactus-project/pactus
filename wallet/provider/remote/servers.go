package remoteprovider

import _ "embed"

//go:embed servers.json
var serversJSON []byte

type ServerInfo struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Website string `json:"website"`
	Address string `json:"address"`
}
