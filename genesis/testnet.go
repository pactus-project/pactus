package genesis

import "encoding/json"

const testnetJSON = `
{
   "GenesisTime":"2021-04-16T00:00:00Z",
   "Params":{
      "BlockVersion":1001,
      "BlockTimeInSecond":10,
      "CommitteeSize":11,
      "BlockReward":100000000,
      "TransactionToLiveInterval":8640,
      "UnbondInterval":181440,
      "MaximumTransactionPerBlock":1000,
      "MaximumMemoLength":1024,
      "FeeFraction":0.001,
      "MinimumFee":1000
   },
   "Accounts":[
      {
         "Address":"zrb1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqn627cy",
         "Balance":2100000000000000
      }
   ],
   "Validators":[
      {
         "PublicKey":"230e0c8723a930af757cda3ec7aab45f7a313fb8c8217e600cf1c90b4c12c1a13c0b9d9e68ef5441240e011f13658d0767d3ce405565cc51d8f3be408d594616be55d31f62167583d945f86732e2293374e1aceeb37fce0d4aacb253ceb98303"
      },
      {
         "PublicKey":"d4f6c52071b4874142089c7f258aa2baa01460660cbbf1aae6f1c7c836e5ec76bd705eb384849f0316b96d86cb59a512c953af60c761ee8dd65ba87f813e3e7e9723159e46f14b1737fb684680c3cebc9437e55f4164af5978d3b6a46f62cd98"
      },
      {
         "PublicKey":"705071243915281c0fd2c4a1186830b016d5d67bb08746fffffa589e21a92496983094996c517ab0f21940c2d635f9059ed25e495f2d384975f9cc2c999f684eee319f365689a8cd0fa6a285197213fa4c1e3c97bfa1478809246bc0fafcf096"
      },
      {
         "PublicKey":"3d03bc94da2da41622e4ce5c0bb235db6dba7002e6ed727a5fa583fd980598f74e3ea4cb74734d3e7c64b88b6bf37c132e162e4f154321344ff176cb488cea8e3baf76201167774a8e9147d5f674a8b9496f4c6cecfaa2272fc257400ece7311"
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
