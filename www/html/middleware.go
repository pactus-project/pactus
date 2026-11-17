package html

import (
	"net/http"

	"github.com/pactus-project/pactus/util/htpasswd"
	"github.com/pactus-project/pactus/www/grpc/basicauth"
	"google.golang.org/grpc/metadata"
)

// BasicAuth checks the credential of the request against the stored credential (username:password_hash).
// Some pages, like the metrics and the profiles, make no gRPC call,
// so the gRPC server can't be the one that rejects a wrong credential.
func BasicAuth(storedCredential string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, password, ok := r.BasicAuth()
		if !ok || htpasswd.CompareBasicAuth(storedCredential, user, password) != nil {
			w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
			http.Error(w, "unauthorized", http.StatusUnauthorized)

			return
		}

		ba := basicauth.New(user, password)
		tokens, _ := ba.GetRequestMetadata(r.Context())
		md := metadata.New(tokens)

		r = r.WithContext(metadata.NewOutgoingContext(r.Context(), md))

		next.ServeHTTP(w, r)
	})
}
