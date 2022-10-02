package http

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestTransaction(t *testing.T) {
	setup(t)

	testBlock := tMockState.TestStore.AddTestBlock(1)
	testTx := testBlock.Transactions()[0]

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
