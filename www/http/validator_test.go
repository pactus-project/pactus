package http

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gorilla/mux"
	"github.com/pactus-project/pactus/crypto"
	"github.com/stretchr/testify/assert"
)

func TestValidator(t *testing.T) {
	setup(t)

	val := tMockState.TestStore.AddTestValidator()

	t.Run("Shall return a validator", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := new(http.Request)
		r = mux.SetURLVars(r, map[string]string{"address": val.Address().String()})
		tHTTPServer.GetValidatorHandler(w, r)

		assert.Equal(t, w.Code, 200)
		fmt.Println(w.Body)
	})

	t.Run("Shall return an error, non exist", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := new(http.Request)
		r = mux.SetURLVars(r, map[string]string{"address": crypto.GenerateTestAddress().String()})
		tHTTPServer.GetValidatorHandler(w, r)

		assert.Equal(t, w.Code, 400)
	})

	t.Run("Shall return an error, invalid address", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := new(http.Request)
		r = mux.SetURLVars(r, map[string]string{"address": "invalid-address"})
		tHTTPServer.GetValidatorHandler(w, r)

		assert.Equal(t, w.Code, 400)
		fmt.Println(w.Body)
	})

	t.Run("Shall return an error, empty address", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := new(http.Request)
		r = mux.SetURLVars(r, map[string]string{"address": ""})
		tHTTPServer.GetValidatorHandler(w, r)

		assert.Equal(t, w.Code, 400)
		fmt.Println(w.Body)
	})

	t.Run("Shall return an error, no address", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := new(http.Request)
		tHTTPServer.GetValidatorHandler(w, r)

		assert.Equal(t, w.Code, 400)
		fmt.Println(w.Body)
	})

	t.Run("Shall return a validator", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := new(http.Request)
		fmt.Println(val.Number())
		fmt.Println(strconv.Itoa(int(val.Number())))
		r = mux.SetURLVars(r, map[string]string{"number": strconv.Itoa(int(val.Number()))})
		tHTTPServer.GetValidatorByNumberHandler(w, r)

		assert.Equal(t, w.Code, 200)
		fmt.Println(w.Body)
	})

	t.Run("Shall return an error, empty number", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := new(http.Request)
		r = mux.SetURLVars(r, map[string]string{"number": ""})
		tHTTPServer.GetValidatorByNumberHandler(w, r)

		assert.Equal(t, w.Code, 400)
		fmt.Println(w.Body)
	})

	t.Run("Shall return an error, invalid number", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := new(http.Request)
		r = mux.SetURLVars(r, map[string]string{"number": "not-a-number"})
		tHTTPServer.GetValidatorByNumberHandler(w, r)

		assert.Equal(t, w.Code, 400)
		fmt.Println(w.Body)
	})

	t.Run("Shall return an error, no number", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := new(http.Request)
		tHTTPServer.GetValidatorByNumberHandler(w, r)

		assert.Equal(t, w.Code, 400)
		fmt.Println(w.Body)
	})
}
