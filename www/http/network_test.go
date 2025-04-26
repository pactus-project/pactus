package http

import (
	"context"
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

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), version.NodeAgent.String())
	assert.Contains(t, w.Body.String(), "zmq_topic")

	td.StopServers()
}

func TestNetworkInfo(t *testing.T) {
	td := setup(t)

	w := httptest.NewRecorder()
	r, err := http.NewRequestWithContext(context.Background(), http.MethodGet,
		"localhost:80?onlyConnected=false", http.NoBody)
	assert.NoError(t, err)

	td.httpServer.NetworkHandler(w, r)

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "Network Name")
	assert.Contains(t, w.Body.String(), "Connected Peers Count")

	td.StopServers()
}
