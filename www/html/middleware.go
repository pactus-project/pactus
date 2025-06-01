package html

import (
	"net/http"

	"github.com/pactus-project/pactus/www/grpc/basicauth"
	"google.golang.org/grpc/metadata"
)

func basicAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, password, ok := r.BasicAuth()
		if !ok {
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
