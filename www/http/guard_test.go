package http

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// The requests go through the guard: the server only accepts JSON requests.
func TestServerGuardsRequests(t *testing.T) {
	td := setupAnchorHTTP(t)

	post := func(contentType string, data []byte) int {
		req, err := http.NewRequestWithContext(t.Context(), http.MethodPost,
			td.endpoint("/pactus/no_such_method"), bytes.NewReader(data))
		require.NoError(t, err)
		req.Header.Set("Content-Type", contentType)

		res, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		defer func() { _ = res.Body.Close() }()

		return res.StatusCode
	}

	// A web page can send a text/plain request without a CORS preflight.
	// The guard checks it before the route: it is refused even for a route that doesn't exist.
	// The guard also bounds the size of the body: it is tested on its own.
	assert.Equal(t, http.StatusUnsupportedMediaType, post("text/plain", []byte(`{}`)))
}
