package html

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestTransaction(t *testing.T) {
	td := setup(t)

	testBlock := td.mockState.TestStore.AddTestBlock(1)
	testTx := testBlock.Transactions()[0]

	t.Run("Shall return a transaction", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := new(http.Request)
		r = mux.SetURLVars(r, map[string]string{"id": testTx.ID().String()})
		td.httpServer.GetTransactionHandler(w, r)

		assert.Equal(t, 200, w.Code)
		assert.Contains(t, w.Body.String(), testTx.Signature().String())
		assert.Contains(t, w.Body.String(), testTx.Signature().String())
		fmt.Println(w.Body)
	})

	t.Run("Shall return an error", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := new(http.Request)
		td.httpServer.GetTransactionHandler(w, r)

		assert.Equal(t, 400, w.Code)
		fmt.Println(w.Body)
	})

	td.StopServers()
}
