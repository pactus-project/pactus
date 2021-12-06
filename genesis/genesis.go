package genesis

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/fxamacker/cbor/v2"
	"github.com/zarbchain/zarb-go/account"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/crypto/bls"
	"github.com/zarbchain/zarb-go/crypto/hash"
	"github.com/zarbchain/zarb-go/param"
	"github.com/zarbchain/zarb-go/util"
	"github.com/zarbchain/zarb-go/validator"
)

type genAccount struct {
	Address crypto.Address `cbor:"1,keyasint"`
	Balance int64          `cbor:"2,keyasint"`
}

type genValidator struct {
	PublicKey *bls.BLSPublicKey `cbor:"1,keyasint"`
}

// Genesis is stored in the state database
type Genesis struct {
	data genesisData
}

type genesisData struct {
	GenesisTime time.Time      `cbor:"1,keyasint"`
	Params      param.Params   `cbor:"2,keyasint"`
	Accounts    []genAccount   `cbor:"3,keyasint"`
	Validators  []genValidator `cbor:"4,keyasint"`
}

func (gen *Genesis) Hash() hash.Hash {
	bs, err := cbor.Marshal(gen.data)
	if err != nil {
		panic(fmt.Errorf("could not create hash of Genesis: %v", err))
	}
	return hash.HashH(bs)
}

func (gen *Genesis) GenesisTime() time.Time {
	return gen.data.GenesisTime
}

func (gen *Genesis) Params() param.Params {
	return gen.data.Params
}

func (gen *Genesis) Accounts() []*account.Account {
	accs := make([]*account.Account, 0)
	for i, genAcc := range gen.data.Accounts {
		acc := account.NewAccount(genAcc.Address, i)
		acc.AddToBalance(genAcc.Balance)
		accs = append(accs, acc)
	}

	return accs
}

func (gen *Genesis) Validators() []*validator.Validator {
	vals := make([]*validator.Validator, 0, len(gen.data.Validators))
	for i, genVal := range gen.data.Validators {
		val := validator.NewValidator(genVal.PublicKey, i)
		vals = append(vals, val)
	}

	return vals
}

func (gen Genesis) MarshalJSON() ([]byte, error) {
	return json.Marshal(&gen.data)
}

func (gen *Genesis) UnmarshalJSON(bs []byte) error {
	return json.Unmarshal(bs, &gen.data)
}

func makeGenesisAccount(acc *account.Account) genAccount {
	return genAccount{
		Address: acc.Address(),
		Balance: acc.Balance(),
	}
}

func makeGenesisValidator(val *validator.Validator) genValidator {
	return genValidator{
		PublicKey: val.PublicKey(),
	}
}

func MakeGenesis(genesisTime time.Time,
	accounts []*account.Account,
	validators []*validator.Validator, params param.Params) *Genesis {

	genAccs := make([]genAccount, 0, len(accounts))
	for _, acc := range accounts {
		genAcc := makeGenesisAccount(acc)
		genAccs = append(genAccs, genAcc)
	}

	genVals := make([]genValidator, 0, len(validators))
	for _, val := range validators {
		genVal := makeGenesisValidator(val)
		genVals = append(genVals, genVal)
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

// LoadFromFile loads genesis object from a JSON file
func LoadFromFile(file string) (*Genesis, error) {
	dat, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	var gen Genesis
	if err := json.Unmarshal(dat, &gen); err != nil {
		return nil, err
	}
	return &gen, nil
}

// SaveToFile saves the genesis info a JSON file
func (gen *Genesis) SaveToFile(file string) error {
	json, err := gen.MarshalJSON()
	if err != nil {
		return err
	}

	// write  dataContent to file
	if err := util.WriteFile(file, json); err != nil {
		return err
	}

	return nil
}
