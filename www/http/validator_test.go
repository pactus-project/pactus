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

func TestValidator(t *testing.T) {
	td := setup(t)

	val := td.mockState.TestStore.AddTestValidator()

	t.Run("Shall return a validator", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := new(http.Request)
		r = mux.SetURLVars(r, map[string]string{"address": val.Address().String()})
		td.httpServer.GetValidatorHandler(w, r)

		assert.Equal(t, w.Code, 200)
		fmt.Println(w.Body)
	})

	t.Run("Shall return an error, non exist", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := new(http.Request)
		r = mux.SetURLVars(r, map[string]string{"address": td.RandomAddress().String()})
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
}
