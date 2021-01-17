package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBlockchainInfo(t *testing.T) {
	setup(t)

	w := httptest.NewRecorder()
	r := new(http.Request)

	tHTTPServer.GetBlockchainHandler(w, r)

	assert.Equal(t, w.Code, 200)
	assert.Contains(t, w.Body.String(), "\"Height\": 2")

}
