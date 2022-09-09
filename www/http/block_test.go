package http

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/stretchr/testify/assert"
)

func TestBlock(t *testing.T) {
	setup(t)

	b := tMockState.TestStore.AddTestBlock(100)

	t.Run("Shall return a block", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := new(http.Request)
		r = mux.SetURLVars(r, map[string]string{"hash": b.Hash().String()})
		tHTTPServer.GetBlockByHashHandler(w, r)

		assert.Equal(t, w.Code, 200)
		assert.Contains(t, w.Body.String(), b.Hash().String())
	})

	t.Run("Shall return a block", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := new(http.Request)
		r = mux.SetURLVars(r, map[string]string{"height": "100"})
		tHTTPServer.GetBlockByHeightHandler(w, r)

		assert.Equal(t, w.Code, 200)
		//fmt.Println(w.Body)
	})

	t.Run("Shall return an error, non exists", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := new(http.Request)
		r = mux.SetURLVars(r, map[string]string{"hash": hash.GenerateTestHash().String()})
		tHTTPServer.GetBlockByHashHandler(w, r)

		assert.Equal(t, w.Code, 400)
		//fmt.Println(w.Body)
	})

	t.Run("Shall return an error, invalid hash", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := new(http.Request)
		r = mux.SetURLVars(r, map[string]string{"hash": "abc"})
		tHTTPServer.GetBlockByHashHandler(w, r)
		fmt.Println(w.Body)

		assert.Equal(t, w.Code, 400)
	})

	t.Run("Shall return an error, empty hash", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := new(http.Request)
		r = mux.SetURLVars(r, map[string]string{"hash": ""})
		tHTTPServer.GetBlockByHashHandler(w, r)
		fmt.Println(w.Body)

		assert.Equal(t, w.Code, 400)
	})

	t.Run("Shall return an error, no hash", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := new(http.Request)
		tHTTPServer.GetBlockByHashHandler(w, r)
		fmt.Println(w.Body)

		assert.Equal(t, w.Code, 400)
	})
}

func TestBlockHash(t *testing.T) {
	setup(t)

	b := tMockState.TestStore.AddTestBlock(100)

	t.Run("Shall return the block hash", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := new(http.Request)
		r = mux.SetURLVars(r, map[string]string{"height": "100"})
		tHTTPServer.GetBlockHashHandler(w, r)

		assert.Equal(t, w.Code, 200)
		assert.Contains(t, w.Body.String(), b.Hash().String())
		//fmt.Println(w.Body)
	})
}
