package genesis

import (
	_ "embed"
	"encoding/json"
)

//go:embed testnet.json
var testnetJSON []byte

func TestnetGenesis() *Genesis {
	gen := new(Genesis)
	_ = json.Unmarshal(testnetJSON, &gen)

	return gen
}
