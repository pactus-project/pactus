package genesis

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/fxamacker/cbor/v2"
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/types/account"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/validator"
	"github.com/pactus-project/pactus/util"
)

type ChainType uint8

const (
	Mainnet  ChainType = 0
	Testnet  ChainType = 1
	Localnet ChainType = 2
)

func (n ChainType) IsMainnet() bool {
	return n == Mainnet
}

func (n ChainType) IsTestnet() bool {
	return n == Testnet
}

func (n ChainType) String() string {
	switch n {
	case Mainnet:
		return "Mainnet"
	case Testnet:
		return "Testnet"
	case Localnet:
		return "Localnet"
	default:
		return "Unknown"
	}
}

type genAccount struct {
	Address string        `cbor:"1,keyasint" json:"address"`
	Balance amount.Amount `cbor:"2,keyasint" json:"balance"`
}

type genValidator struct {
	PublicKey string `cbor:"1,keyasint" json:"public_key"`
}

// Genesis is stored in the state database.
type Genesis struct {
	data genesisData
}

type genesisData struct {
	GenesisTime time.Time      `cbor:"1,keyasint" json:"genesis_time"`
	Params      *GenesisParams `cbor:"2,keyasint" json:"params"`
	Accounts    []genAccount   `cbor:"3,keyasint" json:"accounts"`
	Validators  []genValidator `cbor:"4,keyasint" json:"validators"`
}

func (gen *Genesis) Hash() hash.Hash {
	bs, _ := cbor.Marshal(gen.data)

	return hash.CalcHash(bs)
}

func (gen *Genesis) GenesisTime() time.Time {
	return gen.data.GenesisTime
}

func (gen *Genesis) Params() *GenesisParams {
	return gen.data.Params
}

func (gen *Genesis) Accounts() map[crypto.Address]*account.Account {
	accs := make(map[crypto.Address]*account.Account)
	for i, genAcc := range gen.data.Accounts {
		addr, err := crypto.AddressFromString(genAcc.Address)
		if err != nil {
			panic(err)
		}
		acc := account.NewAccount(int32(i))
		acc.AddToBalance(genAcc.Balance)
		accs[addr] = acc
	}

	return accs
}

func (gen *Genesis) Validators() []*validator.Validator {
	vals := make([]*validator.Validator, 0, len(gen.data.Validators))
	for i, genVal := range gen.data.Validators {
		pub, _ := bls.PublicKeyFromString(genVal.PublicKey)
		val := validator.NewValidator(pub, int32(i))
		vals = append(vals, val)
	}

	return vals
}

//nolint:gocritic // copying an 80-byte Genesis while marshaling JSON is negligible.
func (gen Genesis) MarshalJSON() ([]byte, error) {
	return json.MarshalIndent(&gen.data, "  ", "  ")
}

func (gen *Genesis) UnmarshalJSON(bs []byte) error {
	return json.Unmarshal(bs, &gen.data)
}

func makeGenesisAccount(addr crypto.Address, acc *account.Account) genAccount {
	return genAccount{
		Address: addr.String(),
		Balance: acc.Balance(),
	}
}

func makeGenesisValidator(val *validator.Validator) genValidator {
	return genValidator{
		PublicKey: val.PublicKey().String(),
	}
}

func MakeGenesis(genesisTime time.Time, accounts map[crypto.Address]*account.Account,
	validators []*validator.Validator, params *GenesisParams,
) *Genesis {
	genAccs := make([]genAccount, len(accounts))
	for addr, acc := range accounts {
		genAcc := makeGenesisAccount(addr, acc)
		genAccs[acc.Number()] = genAcc
	}

	genVals := make([]genValidator, len(validators))
	for _, val := range validators {
		genVal := makeGenesisValidator(val)
		genVals[val.Number()] = genVal
	}

	return &Genesis{
		data: genesisData{
			GenesisTime: genesisTime,
			Accounts:    genAccs,
			Validators:  genVals,
			Params:      params,
		},
	}
}

// LoadFromFile loads genesis object from a JSON file.
func LoadFromFile(file string) (*Genesis, error) {
	dat, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	var gen Genesis
	if err := json.Unmarshal(dat, &gen); err != nil {
		return nil, err
	}

	return &gen, nil
}

// Validate checks that the genesis document is well-formed: the parameters are
// present and usable, and every account address and validator public key can be
// parsed and appears only once.
// Addresses and keys are parsed with the currently active HRP, so call it after the
// network type has been selected.
func (gen *Genesis) Validate() error {
	params := gen.data.Params
	if params == nil {
		return InvalidGenesisError{Reason: "params are missing"}
	}

	if params.BlockIntervalInSecond <= 0 {
		return InvalidGenesisError{Reason: "block interval should be positive"}
	}

	// Block times are rounded to the interval inside the minute.
	if 60%params.BlockIntervalInSecond != 0 {
		return InvalidGenesisError{Reason: "block interval should be a divisor of 60 seconds"}
	}

	if params.CommitteeSize <= 0 {
		return InvalidGenesisError{Reason: "committee size should be positive"}
	}

	addrs := make(map[crypto.Address]struct{}, len(gen.data.Accounts))
	for i, genAcc := range gen.data.Accounts {
		addr, err := crypto.AddressFromString(genAcc.Address)
		if err != nil {
			return InvalidGenesisError{Reason: fmt.Sprintf("account %d has an invalid address: %v", i, err)}
		}

		if _, ok := addrs[addr]; ok {
			return InvalidGenesisError{Reason: fmt.Sprintf("account %d has a duplicated address: %s", i, genAcc.Address)}
		}
		addrs[addr] = struct{}{}
	}

	pubKeys := make(map[string]struct{}, len(gen.data.Validators))
	for i, genVal := range gen.data.Validators {
		pub, err := bls.PublicKeyFromString(genVal.PublicKey)
		if err != nil {
			return InvalidGenesisError{Reason: fmt.Sprintf("validator %d has an invalid public key: %v", i, err)}
		}

		if _, ok := pubKeys[pub.String()]; ok {
			return InvalidGenesisError{Reason: fmt.Sprintf("validator %d has a duplicated public key", i)}
		}
		pubKeys[pub.String()] = struct{}{}
	}

	return nil
}

// SaveToFile saves the genesis into a JSON file.
func (gen *Genesis) SaveToFile(file string) error {
	data, err := gen.MarshalJSON()
	if err != nil {
		return err
	}

	// write  dataContent to file
	return util.WriteFile(file, data)
}

func (gen *Genesis) TotalSupply() amount.Amount {
	totalSuppyly := amount.Amount(0)
	for _, acc := range gen.data.Accounts {
		totalSuppyly += acc.Balance
	}

	return totalSuppyly
}

// mainnetHash and testnetHash are the hashes of the published genesis documents.
// They never change, so they are computed once: ChainType is called for every
// gossip bundle, and parsing both documents each time is costly.
var (
	mainnetHash = sync.OnceValue(func() hash.Hash { return MainnetGenesis().Hash() })
	testnetHash = sync.OnceValue(func() hash.Hash { return TestnetGenesis().Hash() })
)

func (gen *Genesis) ChainType() ChainType {
	switch gen.Hash() {
	case mainnetHash():
		return Mainnet
	case testnetHash():
		return Testnet
	default:
		return Localnet
	}
}
