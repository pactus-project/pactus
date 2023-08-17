package http

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

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
		// fmt.Println(w.Body)
	})

	t.Run("Shall return an error, invalid height", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := new(http.Request)
		r = mux.SetURLVars(r, map[string]string{"height": "x"})
		td.httpServer.GetBlockByHeightHandler(w, r)

		assert.Equal(t, w.Code, 400)
		// fmt.Println(w.Body)
	})

	t.Run("Shall return an error, non exists", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := new(http.Request)
		r = mux.SetURLVars(r, map[string]string{"hash": td.RandomHash().String()})
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
}
