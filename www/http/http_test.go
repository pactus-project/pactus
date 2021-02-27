package http

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/account"
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/logger"
	"github.com/zarbchain/zarb-go/state"
	"github.com/zarbchain/zarb-go/sync"
	"github.com/zarbchain/zarb-go/tx"
	"github.com/zarbchain/zarb-go/txpool"
	"github.com/zarbchain/zarb-go/validator"
	"github.com/zarbchain/zarb-go/www/capnp"
)

var tMockState *state.MockState
var tMockPool txpool.TxPool
var tMockSync *sync.MockSync
var tCapnpServer *capnp.Server
var tHTTPServer *Server
var tAccTestAddr crypto.Address
var tValTestAddr crypto.Address
var tTxTestHash crypto.Hash

func init() {
	logger.InitLogger(logger.TestConfig())
}

func setup(t *testing.T) {
	if tHTTPServer != nil {
		return
	}

	valSet, _ := validator.GenerateTestValidatorSet()
	tMockState = state.MockingState(valSet)
	tMockPool = txpool.MockingTxPool()
	tMockSync = sync.MockingSync()

	b1, txs := block.GenerateTestBlock(nil, nil)
	b2, _ := block.GenerateTestBlock(nil, nil)
	tMockState.Store.Blocks[1] = b1
	tMockState.Store.Blocks[2] = b2

	tTxTestHash = txs[0].ID()

	tMockState.Store.Transactions[tTxTestHash] = &tx.CommittedTx{
		Tx:      txs[0],
		Receipt: txs[0].GenerateReceipt(0, b1.Hash()),
	}

	a, _ := account.GenerateTestAccount(888)
	tAccTestAddr = a.Address()
	tMockState.Store.UpdateAccount(a)

	v, _ := validator.GenerateTestValidator(88)
	tValTestAddr = v.Address()
	tMockState.Store.UpdateValidator(v)

	var err error
	tCapnpServer, err = capnp.NewServer(capnp.TestConfig(), tMockState, tMockSync, tMockPool)
	assert.NoError(t, err)
	assert.NoError(t, tCapnpServer.StartServer())

	tHTTPServer, err = NewServer(TestConfig())
	assert.NoError(t, err)
	assert.NoError(t, tHTTPServer.StartServer(tCapnpServer.Address()))
}

func TestRootHandler(t *testing.T) {
	setup(t)

	w := httptest.NewRecorder()
	r := new(http.Request)
	tHTTPServer.RootHandler(w, r)
	assert.Equal(t, w.Code, 200)
	fmt.Println(w.Body)
}
