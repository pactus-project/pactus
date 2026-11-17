package guard

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGuardRequests(t *testing.T) {
	called := false
	handler := Requests(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		if _, err := io.ReadAll(r.Body); err != nil {
			http.Error(w, err.Error(), http.StatusRequestEntityTooLarge)

			return
		}
		w.WriteHeader(http.StatusOK)
	}))

	serve := func(method, contentType string, body io.Reader) int {
		called = false
		req := httptest.NewRequestWithContext(t.Context(), method, "/pactus/wallet/create_wallet", body)
		if contentType != "" {
			req.Header.Set("Content-Type", contentType)
		}
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		return rec.Code
	}

	t.Run("a JSON request passes", func(t *testing.T) {
		code := serve(http.MethodPost, "application/json; charset=utf-8", strings.NewReader(`{}`))
		assert.Equal(t, http.StatusOK, code)
		assert.True(t, called)
	})

	t.Run("a request without body passes", func(t *testing.T) {
		assert.Equal(t, http.StatusOK, serve(http.MethodGet, "", http.NoBody))
		assert.Equal(t, http.StatusOK, serve(http.MethodPost, "", http.NoBody))
	})

	t.Run("a simple cross-site request is rejected", func(t *testing.T) {
		// A web page can send these content types without a CORS preflight.
		for _, contentType := range []string{"text/plain", "application/x-www-form-urlencoded", "multipart/form-data", ""} {
			code := serve(http.MethodPost, contentType, strings.NewReader(`{"wallet_name":"x"}`))
			assert.Equal(t, http.StatusUnsupportedMediaType, code, "content type %q", contentType)
			assert.False(t, called, "content type %q", contentType)
		}
	})

	t.Run("a body that is too large is rejected", func(t *testing.T) {
		body := bytes.NewReader(make([]byte, MaxRequestBodySize+1))
		code := serve(http.MethodPost, "application/json", body)
		assert.Equal(t, http.StatusRequestEntityTooLarge, code)
	})
}
