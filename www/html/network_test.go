package html_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/sync/peerset/peer/service"
	"github.com/pactus-project/pactus/sync/peerset/peer/status"
	"github.com/pactus-project/pactus/version"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNodeInfo(t *testing.T) {
	td := setup(t)

	w := httptest.NewRecorder()
	r := new(http.Request)

	td.httpServer.NodeHandler(w, r)

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), version.NodeAgent.String())
	assert.Contains(t, w.Body.String(), "zmq_topic")
}

func TestNetworkInfo(t *testing.T) {
	td := setup(t)

	w := httptest.NewRecorder()
	r, err := http.NewRequestWithContext(t.Context(), http.MethodGet, "/network/info", http.NoBody)
	require.NoError(t, err)

	td.httpServer.NetworkHandler(w, r)

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "Network Name")
	assert.Contains(t, w.Body.String(), "Connected Peers Count")
	assert.Contains(t, w.Body.String(), "/network/peers")
}

func TestPeerList(t *testing.T) {
	td := setup(t)

	w := httptest.NewRecorder()
	r, err := http.NewRequestWithContext(t.Context(), http.MethodGet, "/network/peers", http.NoBody)
	require.NoError(t, err)

	td.httpServer.PeerListHandler(w, r)

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "Peers")
}

func TestPeerListMoniker(t *testing.T) {
	td := setup(t)

	addPeer := func(moniker string) {
		pid := td.RandPeerID()
		pub, _ := td.RandBLSKeyPair()
		peerSet := td.gRPCServer.FakeSync.PeerSet()
		peerSet.UpdateInfo(pid, moniker, version.NodeAgent.String(),
			[]*bls.PublicKey{pub}, service.New(service.FullNode))
		peerSet.UpdateStatus(pid, status.StatusKnown)
	}
	addPeer("name=alice;web=https://example.com")
	addPeer("Plain <b>Bob</b>\x1b[31m")

	w := httptest.NewRecorder()
	r, err := http.NewRequestWithContext(t.Context(), http.MethodGet, "/network/peers", http.NoBody)
	require.NoError(t, err)

	td.httpServer.PeerListHandler(w, r)

	assert.Equal(t, 200, w.Code)
	body := w.Body.String()

	// A structured moniker: the name, then the whole moniker as text.
	assert.Contains(t, body, "<tr><td>Moniker</td><td>alice</td></tr>")
	assert.Contains(t, body, "<tr><td>Moniker Details</td><td>name=alice;web=https://example.com</td></tr>")
	assert.NotContains(t, body, `href="https://example.com`, "the website must not be a link")

	// A plain moniker: escaped, and without its control characters.
	assert.Contains(t, body, "<tr><td>Moniker</td><td>Plain &lt;b&gt;Bob&lt;/b&gt;[31m</td></tr>")
}
