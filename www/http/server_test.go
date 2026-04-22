package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// passThroughHandler is a minimal handler that records that it was called and
// optionally sets a body so tests can assert whether the CORS wrapper invoked it.
func passThroughHandler(t *testing.T, called *bool) http.Handler {
	t.Helper()

	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		if called != nil {
			*called = true
		}
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("ok"))
		assert.NoError(t, err)
	})
}

func TestAllowCORS(t *testing.T) {
	tests := []struct {
		name              string
		method            string
		origin            string
		preflightMethod   string
		wantAllowOrigin   string
		wantPreflight     bool
		wantHandlerCalled bool
	}{
		{
			name:              "SimpleGETWithOriginSetsAllowOrigin",
			method:            http.MethodGet,
			origin:            "https://example.com",
			wantAllowOrigin:   "https://example.com",
			wantHandlerCalled: true,
		},
		{
			name:              "SimplePOSTWithOriginSetsAllowOrigin",
			method:            http.MethodPost,
			origin:            "https://app.example.com",
			wantAllowOrigin:   "https://app.example.com",
			wantHandlerCalled: true,
		},
		{
			name:              "RequestWithoutOriginGetsNoCORSHeaders",
			method:            http.MethodGet,
			origin:            "",
			wantAllowOrigin:   "",
			wantHandlerCalled: true,
		},
		{
			name:              "PreflightOPTIONSSetsPreflightHeaders",
			method:            http.MethodOptions,
			origin:            "https://example.com",
			preflightMethod:   http.MethodPost,
			wantAllowOrigin:   "https://example.com",
			wantPreflight:     true,
			wantHandlerCalled: false,
		},
		{
			name:              "OPTIONSWithoutPreflightHeaderFallsThrough",
			method:            http.MethodOptions,
			origin:            "https://example.com",
			preflightMethod:   "",
			wantAllowOrigin:   "https://example.com",
			wantPreflight:     false,
			wantHandlerCalled: true,
		},
		{
			name:              "PreflightWithoutOriginDoesNothing",
			method:            http.MethodOptions,
			origin:            "",
			preflightMethod:   http.MethodPost,
			wantAllowOrigin:   "",
			wantPreflight:     false,
			wantHandlerCalled: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			called := false
			handler := allowCORS(passThroughHandler(t, &called))

			req := httptest.NewRequestWithContext(
				t.Context(),
				tt.method,
				"/http/api/blockchain/get_block_hash",
				http.NoBody,
			)
			if tt.origin != "" {
				req.Header.Set("Origin", tt.origin)
			}
			if tt.preflightMethod != "" {
				req.Header.Set("Access-Control-Request-Method", tt.preflightMethod)
			}

			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantAllowOrigin, rec.Header().Get("Access-Control-Allow-Origin"))
			assert.Equal(t, tt.wantHandlerCalled, called, "handler invocation expectation mismatch")

			if tt.wantPreflight {
				allowHeaders := rec.Header().Get("Access-Control-Allow-Headers")
				allowMethods := rec.Header().Get("Access-Control-Allow-Methods")

				assert.Contains(t, allowHeaders, "Authorization")
				assert.Contains(t, allowHeaders, "Content-Type")
				assert.Contains(t, allowHeaders, "Accept")

				for _, m := range []string{"GET", "HEAD", "POST", "PUT", "DELETE"} {
					assert.Contains(t, allowMethods, m, "preflight missing method %q", m)
				}
			} else {
				assert.Empty(t, rec.Header().Get("Access-Control-Allow-Headers"))
				assert.Empty(t, rec.Header().Get("Access-Control-Allow-Methods"))
			}
		})
	}
}

func TestPreflightHandlerSetsAllMethodsAndHeaders(t *testing.T) {
	rec := httptest.NewRecorder()
	preflightHandler(rec)

	allowHeaders := rec.Header().Get("Access-Control-Allow-Headers")
	allowMethods := rec.Header().Get("Access-Control-Allow-Methods")

	// Headers must include everything documented in the preflightHandler comment.
	for _, h := range []string{"Content-Type", "Accept", "Authorization"} {
		assert.Containsf(t, allowHeaders, h, "preflight Access-Control-Allow-Headers %q missing %q", allowHeaders, h)
	}

	// Methods must include everything documented in the preflightHandler comment.
	for _, m := range []string{"GET", "HEAD", "POST", "PUT", "DELETE"} {
		assert.Containsf(t, allowMethods, m, "preflight Access-Control-Allow-Methods %q missing %q", allowMethods, m)
	}
}
