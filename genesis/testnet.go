package genesis

import "encoding/json"

const testnetJSON = `
{
   "ChainName":"zarb-testnet",
   "GenesisTime":"2020-12-20T12:00:00.0+03:30",
   "Params":{
      "BlockTimeInSecond":10,
      "MaximumTransactionPerBlock":1000,
      "MaximumPower":5,
      "SubsidyReductionInterval":2100000,
      "MaximumMemoLength":1024,
      "FeeFraction":0.001,
      "MinimumFee":1000,
      "TransactionToLiveInterval":500
   },
   "Accounts":[
      {
         "Address":"0000000000000000000000000000000000000000",
         "Balance":2100000000000000
      }
   ],
   "Validators":[
      {
         "PublicKey":"b8388bde49b17b62b63f7660435d1480904ea83cd2d4e8758d9ce487dd3d45b88884373f0c1c29afbb1bff3959216f15690bb5da1f1f0857a7dd64999f81c2be17917b468058126883fff3ba0a5cc789cce90134c79a372e5ed2a4d9fbd80b8c"
      },
      {
         "PublicKey":"38b258b8f5be33ebb77bdb3126fd43db0bec88dfcfb6357d9c52a39220b218d744ccfe58926d738de881afe35af4a406fc181b104a17f4a4dba4474a6a450407162a8b6deb5c9282af777505ebc9ccbd1b3a107ef2e77b5d7558a43a8100ae18"
      },
      {
         "PublicKey":"12fe7b4eafee633ec476380afbc87e48b40536a353cd00380bdd18a46b3aeefff5685b54908b9d617949eb6ee66c98178f114b7c747b422e0bc08903a23ba0e27ca2dfe6ed4dc4fa4ae0e8f44c3e1e9ebfadbf98e03a28c5334c5dc3dd26b683"
      },
      {
         "PublicKey":"e2087bfaa5dd4681a2691fdb173f2a7dbe3c3beb9597288f3d5e9825a1d099b3e5322e5269f451ef40ecb8e7cfda0208027d90153deda6a474254ab9abf54066f4e0cb67ebf1dc2427bb21f8a1f08623a36024a37d3904c28fa1fd9e7c4a9b95"
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
