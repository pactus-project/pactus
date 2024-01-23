package genesis

import (
	_ "embed"
	"encoding/json"
)

//go:embed mainnet.json
var mainnetJSON []byte

func MainnetGenesis() *Genesis {
	var gen Genesis
	if err := json.Unmarshal(mainnetJSON, &gen); err != nil {
		panic(err)
	}

	return &gen
}
