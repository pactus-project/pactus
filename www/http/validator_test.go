package http

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestValidator(t *testing.T) {
	setup(t)

	t.Run("Shall return a validator", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := new(http.Request)
		r = mux.SetURLVars(r, map[string]string{"address": tValTestAddr.String()})
		tHTTPServer.GetValidatorHandler(w, r)

		assert.Equal(t, w.Code, 200)
		fmt.Println(w.Body)
	})

	t.Run("Shall return an error", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := new(http.Request)
		r = mux.SetURLVars(r, map[string]string{"address": "invalid-address"})
		tHTTPServer.GetValidatorHandler(w, r)

		assert.Equal(t, w.Code, 400)
		fmt.Println(w.Body)
	})

	t.Run("Shall return an error", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := new(http.Request)
		r = mux.SetURLVars(r, map[string]string{"address": ""})
		tHTTPServer.GetValidatorHandler(w, r)

		assert.Equal(t, w.Code, 400)
		fmt.Println(w.Body)
	})

	t.Run("Shall return an error", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := new(http.Request)
		tHTTPServer.GetValidatorHandler(w, r)

		assert.Equal(t, w.Code, 400)
		fmt.Println(w.Body)
	})
}
