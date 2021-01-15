package http

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"gotest.tools/assert"
)

func TestAccount(t *testing.T) {
	setup(t)

	t.Run("Shall return an account", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := new(http.Request)
		r = mux.SetURLVars(r, map[string]string{"address": tAccTestAddr.String()})
		tHTTPServer.GetAccountHandler(w, r)

		assert.Equal(t, w.Code, 200)
		fmt.Println(w.Body)
	})

	t.Run("Shall return an error", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := new(http.Request)
		r = mux.SetURLVars(r, map[string]string{"address": "invalid-address"})
		tHTTPServer.GetAccountHandler(w, r)

		assert.Equal(t, w.Code, 400)
		fmt.Println(w.Body)
	})

	t.Run("Shall return an error", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := new(http.Request)
		r = mux.SetURLVars(r, map[string]string{"address": ""})
		tHTTPServer.GetAccountHandler(w, r)

		assert.Equal(t, w.Code, 400)
		fmt.Println(w.Body)
	})

	t.Run("Shall return an error", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := new(http.Request)
		tHTTPServer.GetAccountHandler(w, r)

		assert.Equal(t, w.Code, 400)
		fmt.Println(w.Body)
	})
}
