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
	setup(t)

	t.Run("Shall return a block", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := new(http.Request)
		r = mux.SetURLVars(r, map[string]string{"height": "2"})
		tHTTPServer.GetBlockHandler(w, r)

		assert.Equal(t, w.Code, 200)
		//fmt.Println(w.Body)
	})

	t.Run("Shall return an error", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := new(http.Request)
		r = mux.SetURLVars(r, map[string]string{"height": "5"})
		tHTTPServer.GetBlockHandler(w, r)

		assert.Equal(t, w.Code, 400)
		//fmt.Println(w.Body)
	})

	t.Run("Shall return an error", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := new(http.Request)
		r = mux.SetURLVars(r, map[string]string{"height": "-1"})
		tHTTPServer.GetBlockHandler(w, r)
		fmt.Println(w.Body)

		assert.Equal(t, w.Code, 400)
	})

	t.Run("Shall return an error", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := new(http.Request)
		r = mux.SetURLVars(r, map[string]string{"height": "abc"})
		tHTTPServer.GetBlockHandler(w, r)
		fmt.Println(w.Body)

		assert.Equal(t, w.Code, 400)
	})

	t.Run("Shall return an error", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := new(http.Request)
		tHTTPServer.GetBlockHandler(w, r)
		fmt.Println(w.Body)

		assert.Equal(t, w.Code, 400)
	})

}
