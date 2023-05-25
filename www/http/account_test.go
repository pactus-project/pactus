package http

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/util"
	"github.com/stretchr/testify/assert"
)

func TestAccount(t *testing.T) {
	setup(t)

	acc, signer := tMockState.TestStore.AddTestAccount()

	t.Run("Shall return an account", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := new(http.Request)
		r = mux.SetURLVars(r, map[string]string{"address": signer.Address().String()})
		tHTTPServer.GetAccountHandler(w, r)

		assert.Equal(t, w.Code, 200)
		assert.Contains(t, w.Body.String(), util.ChangeToStringWithTrailingZeros(acc.Balance()))
		fmt.Println(w.Body)
	})

	t.Run("Shall return nil, non exist", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := new(http.Request)
		r = mux.SetURLVars(r, map[string]string{"address": crypto.GenerateTestAddress().String()})
		tHTTPServer.GetAccountHandler(w, r)

		assert.Equal(t, w.Code, 400)
	})

	t.Run("Shall return an error, invalid address", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := new(http.Request)
		r = mux.SetURLVars(r, map[string]string{"address": "invalid-address"})
		tHTTPServer.GetAccountHandler(w, r)

		assert.Equal(t, w.Code, 400)
		fmt.Println(w.Body)
	})

	t.Run("Shall return an error, empty address", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := new(http.Request)
		r = mux.SetURLVars(r, map[string]string{"address": ""})
		tHTTPServer.GetAccountHandler(w, r)

		assert.Equal(t, w.Code, 400)
		fmt.Println(w.Body)
	})

	t.Run("Shall return an error, no address", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := new(http.Request)
		tHTTPServer.GetAccountHandler(w, r)

		assert.Equal(t, w.Code, 400)
		fmt.Println(w.Body)
	})
}
