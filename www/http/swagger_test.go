package http

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func httpGet(t *testing.T, client *http.Client, url string) (int, http.Header, []byte) {
	t.Helper()

	req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, url, http.NoBody)
	require.NoError(t, err)

	resp, err := client.Do(req)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp.StatusCode, resp.Header, body
}

// The gateway serves the OpenAPI file under the base path of the config,
// and the file gives that base path to the clients.
func TestSwaggerFollowsTheBasePath(t *testing.T) {
	td := setupAnchorHTTP(t)
	base := "http://" + td.srv.listener.Addr().String()
	noRedirect := &http.Client{
		CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse },
	}

	t.Run("the root redirects to the Swagger UI", func(t *testing.T) {
		status, header, _ := httpGet(t, noRedirect, base+"/http/")
		assert.Equal(t, http.StatusFound, status)
		assert.Equal(t, "/http/ui/", header.Get("Location"))
	})

	t.Run("the Swagger UI is served", func(t *testing.T) {
		status, _, body := httpGet(t, http.DefaultClient, base+"/http/ui/")
		assert.Equal(t, http.StatusOK, status)
		assert.NotEmpty(t, body)
	})

	t.Run("the OpenAPI file carries the base path of the API", func(t *testing.T) {
		status, header, body := httpGet(t, http.DefaultClient, base+"/http/ui/pactus.swagger.json")
		require.Equal(t, http.StatusOK, status)
		assert.Equal(t, "application/json", header.Get("Content-Type"))

		var doc map[string]any
		require.NoError(t, json.Unmarshal(body, &doc))
		assert.Equal(t, "/http/api", doc["basePath"])
		assert.NotEmpty(t, doc["paths"])
	})

	t.Run("an unknown path is not found", func(t *testing.T) {
		status, _, _ := httpGet(t, http.DefaultClient, base+"/unknown")
		assert.Equal(t, http.StatusNotFound, status)
	})
}

// A server that is not enabled doesn't listen, and stopping it is harmless.
func TestDisabledServer(t *testing.T) {
	conf := DefaultConfig()
	require.NoError(t, conf.BasicCheck())
	require.False(t, conf.Enable)

	srv := NewServer(t.Context(), conf)
	require.NoError(t, srv.StartServer("127.0.0.1:0"))
	assert.Nil(t, srv.listener)

	srv.StopServer()
}
