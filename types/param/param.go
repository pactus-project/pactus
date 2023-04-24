package param

import "time"

type Params struct {
	BlockVersion              uint8   `cbor:"1,keyasint"`
	BlockTimeInSecond         int     `cbor:"2,keyasint"`
	CommitteeSize             int     `cbor:"3,keyasint"`
	BlockReward               int64   `cbor:"4,keyasint"`
	TransactionToLiveInterval uint32  `cbor:"5,keyasint"`
	BondInterval              uint32  `cbor:"6,keyasint"`
	UnbondInterval            uint32  `cbor:"7,keyasint"`
	SortitionInterval         uint32  `cbor:"8,keyasint"`
	FeeFraction               float64 `cbor:"9,keyasint"`
	MinimumFee                int64   `cbor:"10,keyasint"`
	MaximumFee                int64   `cbor:"11,keyasint"`
	MaximumStake              int64   `cbor:"12,keyasint"`
}

func DefaultParams() Params {
	return Params{
		BlockVersion:              1,
		BlockTimeInSecond:         10,
		CommitteeSize:             21,
		BlockReward:               1000000000,
		TransactionToLiveInterval: 8640,   // one day
		BondInterval:              360,    // one hour
		UnbondInterval:            181440, // 21 days
		SortitionInterval:         7,
		FeeFraction:               0.0001,
		MinimumFee:                1000,
		MaximumFee:                100000000,
		MaximumStake:              12381000000000,
	}
}

func (p Params) BlockTime() time.Duration {
	return time.Duration(p.BlockTimeInSecond) * time.Second
}

func (p Params) IsMainnet() bool {
	return p.BlockVersion == 0x01
}

func (p Params) IsTestnet() bool {
	return p.BlockVersion == 0x3f // 63
}
