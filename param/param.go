package param

import "time"

type Params struct {
	BlockVersion              uint8   `cbor:"1,keyasint"`
	BlockTimeInSecond         int     `cbor:"2,keyasint"`
	CommitteeSize             int     `cbor:"3,keyasint"`
	BlockReward               int64   `cbor:"4,keyasint"`
	TransactionToLiveInterval int32   `cbor:"5,keyasint"`
	BondInterval              int32   `cbor:"6,keyasint"`
	UnbondInterval            int32   `cbor:"7,keyasint"`
	FeeFraction               float64 `cbor:"8,keyasint"`
	MinimumFee                int64   `cbor:"9,keyasint"`
}

func DefaultParams() Params {
	return Params{
		BlockVersion:              1,
		BlockTimeInSecond:         10,
		CommitteeSize:             21,
		BlockReward:               100000000,
		TransactionToLiveInterval: 8640,   // one day
		BondInterval:              360,    // one hour
		UnbondInterval:            181440, // 21 days
		FeeFraction:               0.001,
		MinimumFee:                1000,
	}
}

func (p Params) BlockTime() time.Duration {
	return time.Duration(p.BlockTimeInSecond) * time.Second
}

func (p Params) IsMainnet() bool {
	return p.BlockVersion == 1
}

func (p Params) IsTestnet() bool {
	return p.BlockVersion == 77
}
