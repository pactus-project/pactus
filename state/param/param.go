package param

import (
	_ "embed"
	"encoding/json"
	"time"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/genesis"
	"github.com/pactus-project/pactus/types"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/protocol"
)

//go:embed foundation_testnet.json
var foundationTestnetBytes []byte

//go:embed foundation_mainnet.json
var foundationMainnetBytes []byte

// MaxDelegateOwnerRewardShare is the maximum stake-owner share of the block reward
// under PIP-49 validator delegation (0.7 PAC in nano PAC).
const MaxDelegateOwnerRewardShare = amount.Amount(0.7 * amount.NanoPACPerPAC)

// Params is the parameters of the Pactus protocol.
// These parameters are fixed among all the nodes in the network.
// TODO: Save them in DB and load them on startng the node.
type Params struct {
	BlockVersion              protocol.Version
	BlockIntervalInSecond     int
	MaxTransactionsPerBlock   int
	CommitteeSize             int
	TransactionToLiveInterval uint32
	BondInterval              uint32
	UnbondInterval            uint32
	SortitionInterval         uint32
	MinimumStake              amount.Amount
	MaximumStake              amount.Amount
	FoundationAddresses       []crypto.Address

	baseBlockReward      amount.Amount
	baseFoundationReward amount.Amount
}

func FromGenesis(genDoc *genesis.Genesis) *Params {
	params := &Params{
		// genesis parameters
		BlockVersion:              genDoc.Params().BlockVersion,
		BlockIntervalInSecond:     genDoc.Params().BlockIntervalInSecond,
		CommitteeSize:             genDoc.Params().CommitteeSize,
		TransactionToLiveInterval: genDoc.Params().TransactionToLiveInterval,
		BondInterval:              genDoc.Params().BondInterval,
		UnbondInterval:            genDoc.Params().UnbondInterval,
		SortitionInterval:         genDoc.Params().SortitionInterval,
		MaximumStake:              genDoc.Params().MaximumStake,
		MinimumStake:              genDoc.Params().MinimumStake,

		// chain parameters
		MaxTransactionsPerBlock: 1000,
		FoundationAddresses:     make([]crypto.Address, 0, 100),

		baseBlockReward:      genDoc.Params().BlockReward,
		baseFoundationReward: amount.Amount(300_000_000),
	}

	foundationAddressList := make([]string, 0)
	switch genDoc.ChainType() {
	case genesis.Mainnet:
		if err := json.Unmarshal(foundationMainnetBytes, &foundationAddressList); err != nil {
			panic(err)
		}

	case genesis.Testnet:
		if err := json.Unmarshal(foundationTestnetBytes, &foundationAddressList); err != nil {
			panic(err)
		}

	case genesis.Localnet:
		for i := 0; i < 100; i++ {
			buf := hash.CalcHash([]byte{byte(i)}).Bytes()
			prv, _ := bls.PrivateKeyFromBytes(buf)

			foundationAddressList = append(foundationAddressList, prv.PublicKeyNative().AccountAddress().String())
		}
	}

	for _, addrStr := range foundationAddressList {
		addr, err := crypto.AddressFromString(addrStr)
		if err != nil {
			panic(err)
		}
		params.FoundationAddresses = append(params.FoundationAddresses, addr)
	}

	return params
}

func (p *Params) BlockInterval() time.Duration {
	return time.Duration(p.BlockIntervalInSecond) * time.Second
}

// RewardCoefficient returns the block reward multiplier based on block height.
// This implements the halving schedule defined in PIP-55:
//
//	Blocks 1 – 8,000,000:    1.000
//	Blocks 8,000,001 – 24M:  0.500
//	Blocks 24,000,001 – 56M: 0.250
//	Blocks 56,000,001+:      0.125
func (*Params) RewardCoefficient(height types.Height) float64 {
	switch {
	case height <= 8_000_000:
		return 1.0
	case height <= 24_000_000:
		return 0.5
	case height <= 56_000_000:
		return 0.25
	default:
		return 0.125
	}
}

func (p *Params) FoundationAddress(height types.Height) crypto.Address {
	addressIndex := int(height) % len(p.FoundationAddresses)

	return p.FoundationAddresses[addressIndex]
}

func (p *Params) BlockReward(height types.Height) amount.Amount {
	return p.baseBlockReward.MulF64(p.RewardCoefficient(height))
}

func (p *Params) FoundationReward(height types.Height) amount.Amount {
	return p.baseFoundationReward.MulF64(p.RewardCoefficient(height))
}
