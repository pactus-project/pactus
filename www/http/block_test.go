package http

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"gotest.tools/assert"
)

func TestBlock(t *testing.T) {
	setup(t)

	t.Run("Shal return a block", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := new(http.Request)
		r = mux.SetURLVars(r, map[string]string{"height": "2"})
		httpServer.GetBlockHandler(w, r)

		assert.Equal(t, w.Code, 200)
		//fmt.Println(w.Body)
	})

	// TODO: handle errors in better way
	t.Run("???Shal return an error???", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := new(http.Request)
		r = mux.SetURLVars(r, map[string]string{"height": "-1"})
		httpServer.GetBlockHandler(w, r)
		fmt.Println(w.Body)

		assert.Equal(t, w.Code, 200)
	})

}
