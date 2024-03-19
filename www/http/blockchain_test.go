package http

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestBlockchainInfo(t *testing.T) {
	td := setup(t)

	td.mockState.CommitTestBlocks(10)

	w := httptest.NewRecorder()
	r := new(http.Request)

	td.httpServer.BlockchainHandler(w, r)

	assert.Equal(t, w.Code, 200)
	assert.Contains(t, w.Body.String(), "10")

	td.StopServers()
}

func TestBlock(t *testing.T) {
	td := setup(t)

	b := td.mockState.TestStore.AddTestBlock(100)

	t.Run("Shall return a block", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := new(http.Request)
		r = mux.SetURLVars(r, map[string]string{"hash": b.Hash().String()})
		td.httpServer.GetBlockByHashHandler(w, r)

		assert.Equal(t, w.Code, 200)
		assert.Contains(t, w.Body.String(), b.Hash().String())
	})

	t.Run("Shall return a block", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := new(http.Request)
		r = mux.SetURLVars(r, map[string]string{"height": "100"})
		td.httpServer.GetBlockByHeightHandler(w, r)

		assert.Equal(t, w.Code, 200)
	})

	t.Run("Shall return an error, invalid height", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := new(http.Request)
		r = mux.SetURLVars(r, map[string]string{"height": "x"})
		td.httpServer.GetBlockByHeightHandler(w, r)

		assert.Equal(t, w.Code, 400)
	})

	t.Run("Shall return an error, non exists", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := new(http.Request)
		r = mux.SetURLVars(r, map[string]string{"hash": td.RandHash().String()})
		td.httpServer.GetBlockByHashHandler(w, r)

		assert.Equal(t, w.Code, 400)
	})

	t.Run("Shall return an error, invalid hash", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := new(http.Request)
		r = mux.SetURLVars(r, map[string]string{"hash": "abc"})
		td.httpServer.GetBlockByHashHandler(w, r)
		fmt.Println(w.Body)

		assert.Equal(t, w.Code, 400)
	})

	t.Run("Shall return an error, empty hash", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := new(http.Request)
		r = mux.SetURLVars(r, map[string]string{"hash": ""})
		td.httpServer.GetBlockByHashHandler(w, r)
		fmt.Println(w.Body)

		assert.Equal(t, w.Code, 400)
	})

	t.Run("Shall return an error, no hash", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := new(http.Request)
		td.httpServer.GetBlockByHashHandler(w, r)
		fmt.Println(w.Body)

		assert.Equal(t, w.Code, 400)
	})

	td.StopServers()
}

func TestAccount(t *testing.T) {
	td := setup(t)

	acc, addr := td.mockState.TestStore.AddTestAccount()

	t.Run("Shall return an account", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := new(http.Request)
		r = mux.SetURLVars(r, map[string]string{"address": addr.String()})
		td.httpServer.GetAccountHandler(w, r)

		assert.Equal(t, w.Code, 200)
		assert.Contains(t, w.Body.String(), acc.Balance().String())
		fmt.Println(w.Body)
	})

	t.Run("Shall return nil, non exist", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := new(http.Request)
		r = mux.SetURLVars(r, map[string]string{"address": td.RandAccAddress().String()})
		td.httpServer.GetAccountHandler(w, r)

		assert.Equal(t, w.Code, 400)
	})

	t.Run("Shall return an error, invalid address", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := new(http.Request)
		r = mux.SetURLVars(r, map[string]string{"address": "invalid-address"})
		td.httpServer.GetAccountHandler(w, r)

		assert.Equal(t, w.Code, 400)
		fmt.Println(w.Body)
	})

	t.Run("Shall return an error, empty address", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := new(http.Request)
		r = mux.SetURLVars(r, map[string]string{"address": ""})
		td.httpServer.GetAccountHandler(w, r)

		assert.Equal(t, w.Code, 400)
		fmt.Println(w.Body)
	})

	t.Run("Shall return an error, no address", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := new(http.Request)
		td.httpServer.GetAccountHandler(w, r)

		assert.Equal(t, w.Code, 400)
		fmt.Println(w.Body)
	})

	td.StopServers()
}

func TestValidator(t *testing.T) {
	td := setup(t)

	val := td.mockState.TestStore.AddTestValidator()

	t.Run("Shall return an error, non exist", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := new(http.Request)
		r = mux.SetURLVars(r, map[string]string{"address": td.RandAccAddress().String()})
		td.httpServer.GetValidatorHandler(w, r)

		assert.Equal(t, w.Code, 400)
	})

	t.Run("Shall return an error, invalid address", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := new(http.Request)
		r = mux.SetURLVars(r, map[string]string{"address": "invalid-address"})
		td.httpServer.GetValidatorHandler(w, r)

		assert.Equal(t, w.Code, 400)
		fmt.Println(w.Body)
	})

	t.Run("Shall return an error, empty address", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := new(http.Request)
		r = mux.SetURLVars(r, map[string]string{"address": ""})
		td.httpServer.GetValidatorHandler(w, r)

		assert.Equal(t, w.Code, 400)
		fmt.Println(w.Body)
	})

	t.Run("Shall return an error, no address", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := new(http.Request)
		td.httpServer.GetValidatorHandler(w, r)

		assert.Equal(t, w.Code, 400)
		fmt.Println(w.Body)
	})

	t.Run("Shall return a validator", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := new(http.Request)
		r = mux.SetURLVars(r, map[string]string{"address": val.Address().String()})

		td.httpServer.GetValidatorHandler(w, r)

		assert.Equal(t, w.Code, 200)
		assert.Contains(t, w.Body.String(), "0.987")
		fmt.Println(w.Body)
	})

	td.StopServers()
}

func TestValidatorByNumber(t *testing.T) {
	td := setup(t)

	val := td.mockState.TestStore.AddTestValidator()

	t.Run("Shall return a validator", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := new(http.Request)
		fmt.Println(val.Number())
		fmt.Println(strconv.Itoa(int(val.Number())))
		r = mux.SetURLVars(r, map[string]string{"number": strconv.Itoa(int(val.Number()))})
		td.httpServer.GetValidatorByNumberHandler(w, r)

		assert.Equal(t, w.Code, 200)
		fmt.Println(w.Body)
	})

	t.Run("Shall return a error, non exist", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := new(http.Request)
		fmt.Println(val.Number())
		fmt.Println(strconv.Itoa(int(val.Number())))
		r = mux.SetURLVars(r, map[string]string{"number": strconv.Itoa(int(val.Number() + 1))})
		td.httpServer.GetValidatorByNumberHandler(w, r)

		assert.Equal(t, w.Code, 400)
		fmt.Println(w.Body)
	})

	t.Run("Shall return an error, empty number", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := new(http.Request)
		r = mux.SetURLVars(r, map[string]string{"number": ""})
		td.httpServer.GetValidatorByNumberHandler(w, r)

		assert.Equal(t, w.Code, 400)
		fmt.Println(w.Body)
	})

	t.Run("Shall return an error, invalid number", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := new(http.Request)
		r = mux.SetURLVars(r, map[string]string{"number": "not-a-number"})
		td.httpServer.GetValidatorByNumberHandler(w, r)

		assert.Equal(t, w.Code, 400)
		fmt.Println(w.Body)
	})

	t.Run("Shall return an error, no number", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := new(http.Request)
		td.httpServer.GetValidatorByNumberHandler(w, r)

		assert.Equal(t, w.Code, 400)
		fmt.Println(w.Body)
	})

	td.StopServers()
}

func TestConsensusInfo(t *testing.T) {
	td := setup(t)

	h, _ := td.mockConsMgr.HeightRound()
	v1, _ := td.GenerateTestPrepareVote(h, 1)
	v2, _ := td.GenerateTestPrecommitVote(h, 2)
	td.mockConsMgr.AddVote(v1)
	td.mockConsMgr.AddVote(v2)

	w := httptest.NewRecorder()
	r := new(http.Request)

	td.httpServer.ConsensusHandler(w, r)

	assert.Equal(t, w.Code, 200)
	assert.Contains(t, w.Body.String(), "<td>2</td>")
	assert.Contains(t, w.Body.String(), v2.Signer().String())

	td.StopServers()
}
