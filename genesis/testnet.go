package genesis

import (
	_ "embed"
	"encoding/json"
)

//go:embed testnet.json
var testnetJSON []byte

func TestnetGenesis() *Genesis {
	var gen Genesis
	if err := json.Unmarshal(testnetJSON, &gen); err != nil {
		panic(err)
	}
	return &gen
}
