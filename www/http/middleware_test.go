package http

import (
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/metadata"
)

func TestBasicAuthMiddleware(t *testing.T) {
	handler := basicAuth(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("authorized"))
		assert.NoError(t, err)
	}))

	t.Run("NoAuth", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
		assert.Equal(t, `Basic realm="restricted", charset="UTF-8"`, rr.Header().Get("WWW-Authenticate"))
	})

	t.Run("WithAuth", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
		req.SetBasicAuth("username", "password")
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, "authorized", rr.Body.String())
	})

	t.Run("CheckMetadata", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
		req.SetBasicAuth("username", "password")
		rr := httptest.NewRecorder()

		checkMetadataHandler := basicAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			md, ok := metadata.FromOutgoingContext(r.Context())
			if !ok {
				t.Errorf("No metadata in context")
			}

			auth := md["authorization"][0]

			const prefix = "Basic "
			c, err := base64.StdEncoding.DecodeString(auth[len(prefix):])
			assert.NoError(t, err)
			cs := string(c)
			username, password, ok := strings.Cut(cs, ":")
			assert.True(t, ok)

			assert.Equal(t, "username", username)
			assert.Equal(t, "password", password)

			w.WriteHeader(http.StatusOK)
		}))

		checkMetadataHandler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	})
}
