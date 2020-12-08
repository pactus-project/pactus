package genesis

import "encoding/json"

const mainnetJSON = `
{
    "chainName": "zarb-mainnet"
}
`

func Mainnet() *Genesis {
	var gen Genesis
	if err := json.Unmarshal([]byte(mainnetJSON), &gen); err != nil {
		panic(err)
	}
	return &gen
}
