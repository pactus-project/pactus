package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pactus-project/pactus/version"
	"github.com/stretchr/testify/assert"
)

func TestNodeInfo(t *testing.T) {
	td := setup(t)

	w := httptest.NewRecorder()
	r := new(http.Request)

	td.httpServer.NodeHandler(w, r)

	assert.Equal(t, w.Code, 200)
	assert.Contains(t, w.Body.String(), version.Agent())

	td.StopServers()
}

func TestNetworkInfo(t *testing.T) {
	td := setup(t)

	w := httptest.NewRecorder()
	r := new(http.Request)

	td.httpServer.NetworkHandler(w, r)

	assert.Equal(t, w.Code, 200)
	assert.Contains(t, w.Body.String(), "Network Name")
	assert.Contains(t, w.Body.String(), "Connected Peers Count")

	td.StopServers()
}
