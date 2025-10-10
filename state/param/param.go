package param

import (
	_ "embed"
	"encoding/json"
	"time"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/genesis"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/protocol"
)

//go:embed foundation_testnet.json
var foundationTestnetBytes []byte

//go:embed foundation_mainnet.json
var foundationMainnetBytes []byte

// Params is the parameters of the Pactus protocol.
// These parameters are fixed among all the nodes in the network.
// TODO: Save them in DB and load them on startng the node.
type Params struct {
	BlockVersion              protocol.Version
	BlockIntervalInSecond     int
	MaxTransactionsPerBlock   int
	CommitteeSize             int
	BlockReward               amount.Amount
	TransactionToLiveInterval uint32
	BondInterval              uint32
	UnbondInterval            uint32
	SortitionInterval         uint32
	MinimumStake              amount.Amount
	MaximumStake              amount.Amount
	FoundationReward          amount.Amount
	FoundationAddress         []crypto.Address
}

func FromGenesis(genDoc *genesis.Genesis) *Params {
	params := &Params{
		// genesis parameters
		BlockVersion:              genDoc.Params().BlockVersion,
		BlockIntervalInSecond:     genDoc.Params().BlockIntervalInSecond,
		CommitteeSize:             genDoc.Params().CommitteeSize,
		BlockReward:               genDoc.Params().BlockReward,
		TransactionToLiveInterval: genDoc.Params().TransactionToLiveInterval,
		BondInterval:              genDoc.Params().BondInterval,
		UnbondInterval:            genDoc.Params().UnbondInterval,
		SortitionInterval:         genDoc.Params().SortitionInterval,
		MaximumStake:              genDoc.Params().MaximumStake,
		MinimumStake:              genDoc.Params().MinimumStake,

		// chain parameters
		MaxTransactionsPerBlock: 1000,
		FoundationAddress:       make([]crypto.Address, 0, 100),
		FoundationReward:        amount.Amount(300_000_000),
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
			buf := make([]byte, bls.PrivateKeySize)
			buf[0] = byte(i)
			prv, _ := bls.PrivateKeyFromBytes(buf)

			foundationAddressList = append(foundationAddressList, prv.PublicKeyNative().AccountAddress().String())
		}
	}

	for _, addrStr := range foundationAddressList {
		addr, err := crypto.AddressFromString(addrStr)
		if err != nil {
			panic(err)
		}
		params.FoundationAddress = append(params.FoundationAddress, addr)
	}

	return params
}

func (p *Params) BlockInterval() time.Duration {
	return time.Duration(p.BlockIntervalInSecond) * time.Second
}
