package genesis

import "encoding/json"

const testnetJSON = `
{
	"GenesisTime":"2022-02-21T00:00:00Z",
	"Params":{
	   "BlockVersion":77,
	   "BlockTimeInSecond":10,
	   "CommitteeSize":7,
	   "BlockReward":100000000,
	   "TransactionToLiveInterval":2880,
	   "BondInterval":120,
	   "UnbondInterval":60480,
	   "FeeFraction":0.0001,
	   "MinimumFee":10000
	},
	"Accounts":[
	   {
		  "Address":"000000000000000000000000000000000000000000",
		  "Balance":2100000000000000
	   }
	],
	"Validators":[
	   {
		  "PublicKey":"af5db675ec83d0c87d33cd033ef7561b3a8e175bcd6ed5dddb72da85f11ce17b7442bc41ab842d2dd1769b2108e1b79d0bb3be03dae1c9ae658948a9206a0a8b147ff5c1847b290d27b6cfb3585ad737cac21362c41fd02b20a1a58c99a1ff67"
	   },
	   {
		  "PublicKey":"85956a2b952b5434c94767169b2b0c9d029e195fa6030ef802aca34d39d27988f56919f9fac1d7196c88bb3ad34eafdd08b2cf7ff5c1241a3f16507a8a35423cf990b63abfb5372f46990e5994e17329111bcfc5510cd74a57b654f6ea49f27e"
	   },
	   {
		  "PublicKey":"85a2aef96a1c983ed245b9f45dc0b9fe9614dbe27e021c9ddeeb5a8e06ab1e2b75420de1d88dc6634a067b7cd8af15a40d000fc786e4c753f91026fc6e643f6482a6fde6a571381ade5fe03429c2e57af5dfc52d8740b97d40c140ccb5325e78"
	   },
	   {
		  "PublicKey":"92fccbd7654deaec256b6e048c331198adc7ad9256677f50195d3150d191d60ae936cb506520370b3734de64162976f20fa0a4e6e540c0b1532f0f0d6e6b0021372e52fb1408379a7d12c474744d868a6b304c8a466ec252fa7bb3b8b77c7183"
	   }
	]
}
`

func Testnet() *Genesis {
	var gen Genesis
	if err := json.Unmarshal([]byte(testnetJSON), &gen); err != nil {
		panic(err)
	}
	return &gen
}
