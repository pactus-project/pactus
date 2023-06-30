package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pactus-project/pactus/version"
	"github.com/stretchr/testify/assert"
)

func TestBlockchainInfo(t *testing.T) {
	td := setup(t)

	td.mockState.CommitTestBlocks(10)

	w := httptest.NewRecorder()
	r := new(http.Request)

	td.httpServer.BlockchainHandler(w, r)

	assert.Equal(t, w.Code, 200)
	assert.Contains(t, w.Body.String(), "10")
}

func TestNetworkInfo(t *testing.T) {
	td := setup(t)

	w := httptest.NewRecorder()
	r := new(http.Request)

	td.httpServer.NetworkHandler(w, r)

	assert.Equal(t, w.Code, 200)
	assert.Contains(t, w.Body.String(), "Peers")
	assert.Contains(t, w.Body.String(), "ID")
	//	fmt.Println(w.Body.String())
}

func TestNodeInfo(t *testing.T) {
	td := setup(t)

	w := httptest.NewRecorder()
	r := new(http.Request)

	td.httpServer.NodeHandler(w, r)

	assert.Equal(t, w.Code, 200)
	assert.Contains(t, w.Body.String(), version.Agent())
	//	fmt.Println(w.Body.String())
}
