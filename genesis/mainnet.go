package genesis

import "encoding/json"

const mainnetJSON = `
{
	"GenesisTime":"2021-04-07T00:00:00Z",
	"Params":{
	   "BlockVersion":1,
	   "BlockTimeInSecond":10,
	   "CommitteeSize":21,
	   "BlockReward":100000000,
	   "TransactionToLiveInterval":8640,
	   "UnbondInterval":181440,
	   "MaximumTransactionPerBlock":1000,
	   "MaximumMemoLength":64,
	   "FeeFraction":0.001,
	   "MinimumFee":1000
	},
	"Accounts":[
	   {
		  "Address":"0000000000000000000000000000000000000000",
		  "Balance":2100000000000000
	   }
	],
	"Validators":[
	   {
		  "PublicKey":"594ac38ee38949356e139340cd9668f48d908e76b44781e7013e3f70b738a9b6b53e95dfcba23bd1bbe923d2df354815986643467f25b755d76a908c0dca20327cc111e16d30f37041a23417f8d7cb446cc891c551176df641f07c1f4e1e068b"
	   },
	   {
		  "PublicKey":"332f2f3a6250b7ff955cd73a0b43e567e82b1e6f4e5ace219b74408deefe995b96481d673ce99b20ce62c2177c05880b37b42d3d63f6e7a951492166e74cec3625870582f4a8b8b135abeb4dd171455a2a4a413b79a50b7ace4f8a3123b1ed8f"
	   },
	   {
		  "PublicKey":"0fe092c870d0cee720a30388e40f14ec2df38526e3db040efd30d2b59df1afd5b25568b87806799c829cd65659a84e193f1dfbb67e9aea6eefd4fbf9dd6ddaac694d59efba0df6aba336c1e373d0228514481edf9cce376933a05a9d8e60830f"
	   },
	   {
		  "PublicKey":"2a59438cb5790fa9d0d3c584e54dfc6f41f998b0bbd5297ea74d5a1b62b1022f222ee1c966fcc060ef199bd867d0d80b416a10423b070d3dd5d4c8d32678b9a686703fe818ba662162416389965579162134c622a9d90d10fb508eef03c38d92"
	   }
	]
}
`

func Mainnet() *Genesis {
	var gen Genesis
	if err := json.Unmarshal([]byte(mainnetJSON), &gen); err != nil {
		panic(err)
	}
	return &gen
}
