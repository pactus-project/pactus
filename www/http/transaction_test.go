package http

import (
	"encoding/hex"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/tx"
)

func TestTransaction(t *testing.T) {
	setup(t)

	testTx := tMockState.TestStore.AddTestTransaction()

	t.Run("Shall return a transaction", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := new(http.Request)
		r = mux.SetURLVars(r, map[string]string{"id": testTx.ID().String()})
		tHTTPServer.GetTransactionHandler(w, r)

		assert.Equal(t, w.Code, 200)
		fmt.Println(w.Body)
	})

	t.Run("Shall return an error", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := new(http.Request)
		tHTTPServer.GetTransactionHandler(w, r)

		assert.Equal(t, w.Code, 400)
		fmt.Println(w.Body)
	})

}

func TestSendTransaction(t *testing.T) {
	setup(t)

	t.Run("Send valid transaction", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := new(http.Request)
		trx, _ := tx.GenerateTestSendTx()
		data, _ := trx.Bytes()
		r = mux.SetURLVars(r, map[string]string{"data": hex.EncodeToString(data)})
		tHTTPServer.SendRawTransactionHandler(w, r)

		assert.Equal(t, w.Code, 200)
		fmt.Println(w.Body)
	})

	t.Run("Send invalid transaction", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := new(http.Request)
		trx, _ := tx.GenerateTestSendTx()
		trx.SetSignature(nil)
		data, _ := trx.Bytes()
		r = mux.SetURLVars(r, map[string]string{"data": hex.EncodeToString(data)})
		tHTTPServer.SendRawTransactionHandler(w, r)

		assert.Equal(t, w.Code, 400)
		fmt.Println(w.Body)
	})

	t.Run("Send invalid input data", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := new(http.Request)
		r = mux.SetURLVars(r, map[string]string{"data": "invalid data"})
		tHTTPServer.SendRawTransactionHandler(w, r)
		assert.Equal(t, w.Code, 400)
		fmt.Println(w.Body)
	})

	t.Run("Send invalid marshal", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := new(http.Request)
		r = mux.SetURLVars(r, map[string]string{"data": "010203"})
		tHTTPServer.SendRawTransactionHandler(w, r)
		assert.Equal(t, w.Code, 400)
		fmt.Println(w.Body)
	})

}
