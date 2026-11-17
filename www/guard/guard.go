// Package guard checks the HTTP requests that the gateways pass to the gRPC server.
package guard

import (
	"mime"
	"net/http"
)

// MaxRequestBodySize is the maximum size of the body of a request.
// It is the maximum size of a message that the gRPC server accepts.
const MaxRequestBodySize = 4 << 20

// Requests rejects the requests that a gateway must not pass to the gRPC server.
//
// A request with a body must be a JSON request. A web page can make a browser send
// a "simple" request (like a text/plain POST) to a local node, without a CORS preflight
// and whatever the CORS settings are. A JSON content type always needs a preflight.
//
// The body is limited as well: the gateway would decode it entirely
// before the gRPC server rejects a message that is too large.
func Requests(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
			if r.ContentLength != 0 {
				mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
				if err != nil || mediaType != "application/json" {
					http.Error(w, "content type must be application/json", http.StatusUnsupportedMediaType)

					return
				}
			}
		}

		r.Body = http.MaxBytesReader(w, r.Body, MaxRequestBodySize)

		handler.ServeHTTP(w, r)
	})
}
