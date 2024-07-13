package genesis

import (
	_ "embed"
	"encoding/json"
)

//go:embed mainnet.json
var mainnetJSON []byte

func MainnetGenesis() *Genesis {
	gen := new(Genesis)
	_ = json.Unmarshal(mainnetJSON, &gen)

	return gen
}
