package http

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/account"
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/logger"
	"github.com/zarbchain/zarb-go/store"
	"github.com/zarbchain/zarb-go/tx"
	"github.com/zarbchain/zarb-go/txpool"
	"github.com/zarbchain/zarb-go/util"
	"github.com/zarbchain/zarb-go/validator"
	"github.com/zarbchain/zarb-go/www/capnp"
)

type MockStore struct {
	Blocks       map[int]*block.Block
	Accounts     map[crypto.Address]*account.Account
	Validators   map[crypto.Address]*validator.Validator
	Transactions map[crypto.Hash]*tx.CommittedTx
}

type MockState struct {
	Store *MockStore
}

func (m *MockState) StoreReader() store.StoreReader {
	return m.Store
}
func (m *MockState) ValidatorSet() *validator.ValidatorSet {
	return nil
}
func (m *MockState) LastBlockHeight() int {
	return -1
}
func (m *MockState) GenesisHash() crypto.Hash {
	return crypto.UndefHash
}
func (m *MockState) LastBlockHash() crypto.Hash {
	return crypto.UndefHash
}
func (m *MockState) LastBlockTime() time.Time {
	return util.Now()
}
func (m *MockState) LastCommit() *block.Commit {
	return nil
}
func (m *MockState) BlockTime() time.Duration {
	return time.Second
}
func (m *MockState) UpdateLastCommit(blockHash crypto.Hash, commit block.Commit) {

}
func (m *MockState) Fingerprint() string {
	return ""
}

func (m *MockStore) Block(height int) (*block.Block, error) {
	b, ok := m.Blocks[height]
	if ok {
		return b, nil
	}
	return nil, fmt.Errorf("Not found")
}
func (m *MockStore) BlockHeight(hash crypto.Hash) (int, error) {
	return 0, nil
}
func (m *MockStore) Transaction(hash crypto.Hash) (*tx.CommittedTx, error) {
	b, ok := m.Transactions[hash]
	if ok {
		return b, nil
	}
	return nil, fmt.Errorf("Not found")
}
func (m *MockStore) HasAccount(addr crypto.Address) bool {
	_, ok := m.Accounts[addr]
	return ok
}
func (m *MockStore) Account(addr crypto.Address) (*account.Account, error) {
	a, ok := m.Accounts[addr]
	if ok {
		return a, nil
	}
	return nil, fmt.Errorf("Not found")
}
func (m *MockStore) TotalAccounts() int {
	return len(m.Accounts)
}
func (m *MockStore) HasValidator(addr crypto.Address) bool {
	_, ok := m.Validators[addr]
	return ok
}
func (m *MockStore) Validator(addr crypto.Address) (*validator.Validator, error) {
	v, ok := m.Validators[addr]
	if ok {
		return v, nil
	}
	return nil, fmt.Errorf("Not found")
}
func (m *MockStore) TotalValidators() int {
	return len(m.Validators)
}

var mockState *MockState
var mockPool txpool.TxPoolReader
var capnpServer *capnp.Server
var httpServer *Server
var accTestAddr crypto.Address
var valTestAddr crypto.Address
var txTestHash crypto.Hash

func setup(t *testing.T) {
	if httpServer != nil {
		return
	}

	mockState = &MockState{
		Store: &MockStore{
			Blocks:       make(map[int]*block.Block),
			Accounts:     make(map[crypto.Address]*account.Account),
			Validators:   make(map[crypto.Address]*validator.Validator),
			Transactions: make(map[crypto.Hash]*tx.CommittedTx),
		},
	}
	mockPool = txpool.NewMockTxPool()

	b1, txs := block.GenerateTestBlock(nil)
	b2, _ := block.GenerateTestBlock(nil)
	mockState.Store.Blocks[1] = &b1
	mockState.Store.Blocks[2] = &b2

	txTestHash = txs[0].Hash()

	mockState.Store.Transactions[txTestHash] = &tx.CommittedTx{
		Tx:      txs[0],
		Receipt: txs[0].GenerateReceipt(0, b1.Hash()),
	}

	a, _ := account.GenerateTestAccount(88)
	accTestAddr = a.Address()
	mockState.Store.Accounts[accTestAddr] = a

	v, _ := validator.GenerateTestValidator(88)
	valTestAddr = v.Address()
	mockState.Store.Validators[valTestAddr] = v

	loggerConfig := logger.TestConfig()
	logger.InitLogger(loggerConfig)

	capnpConf := capnp.TestConfig()
	s, err := capnp.NewServer(capnpConf, mockState, mockPool)
	assert.NoError(t, err)
	assert.NoError(t, s.StartServer())
	capnpServer = s

	httpConf := TestConfig()
	httpServer, err = NewServer(httpConf)
	assert.NoError(t, err)
	assert.NoError(t, httpServer.StartServer(capnpServer.Address()))
}
