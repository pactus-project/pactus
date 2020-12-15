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
	Blocks map[int]*block.Block
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
	return nil, nil
}
func (m *MockStore) HasAccount(crypto.Address) bool {
	return false
}
func (m *MockStore) Account(addr crypto.Address) (*account.Account, error) {
	return nil, nil
}
func (m *MockStore) TotalAccounts() int {
	return 0
}
func (m *MockStore) HasValidator(crypto.Address) bool {
	return false
}
func (m *MockStore) Validator(addr crypto.Address) (*validator.Validator, error) {
	return nil, nil
}
func (m *MockStore) TotalValidators() int {
	return 0
}

var mockState *MockState
var mockPool txpool.TxPoolReader
var capnpServer *capnp.Server
var httpServer *Server

func setup(t *testing.T) {
	mockState = &MockState{
		Store: &MockStore{
			Blocks: make(map[int]*block.Block),
		},
	}
	mockPool = txpool.NewMockTxPool()

	b1, _ := block.GenerateTestBlock(nil)
	b2, _ := block.GenerateTestBlock(nil)
	mockState.Store.Blocks[1] = &b1
	mockState.Store.Blocks[2] = &b2

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
