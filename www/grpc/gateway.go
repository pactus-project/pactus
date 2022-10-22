package grpc

import (
	"fmt"
	"mime"
	"net/http"
	"strings"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
	"github.com/rakyll/statik/fs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	// Static files
	_ "github.com/pactus-project/pactus/www/grpc/statik"
)

type GatewayConfig struct {
	Enable     bool   `toml:"enable"`
	Listen     string `toml:"listen"`
	EnableCORS bool   `toml:"enable_cors"`
}

// getOpenAPIHandler serves an OpenAPI UI.
// https://github.com/philips/grpc-gateway-example/blob/master/cmd/serve.go
func (s *Server) getOpenAPIHandler() (http.Handler, error) {
	err := mime.AddExtensionType(".svg", "image/svg+xml")
	if err != nil {
		return nil, err
	}

	statikFS, err := fs.New()
	if err != nil {
		return nil, err
	}

	return http.FileServer(statikFS), nil
}

func (s *Server) startGateway() error {
	if !s.config.Gateway.Enable {
		return nil
	}

	conn, err := grpc.DialContext(
		s.ctx,
		s.config.Listen,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return fmt.Errorf("failed to dial server: %w", err)
	}

	gwMux := runtime.NewServeMux()
	err = pactus.RegisterBlockchainHandler(s.ctx, gwMux, conn)
	if err != nil {
		return err
	}
	err = pactus.RegisterTransactionHandler(s.ctx, gwMux, conn)
	if err != nil {
		return err
	}
	err = pactus.RegisterNetworkHandler(s.ctx, gwMux, conn)
	if err != nil {
		return err
	}

	oa, err := s.getOpenAPIHandler()
	if err != nil {
		return err
	}

	gwServer := &http.Server{
		Addr:              s.config.Gateway.Listen,
		ReadHeaderTimeout: 3 * time.Second,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasPrefix(r.URL.Path, "/v1") {
				gwMux.ServeHTTP(w, r)
				return
			}
			oa.ServeHTTP(w, r)
		}),
	}
	if s.config.Gateway.EnableCORS {
		gwServer.Handler = allowCORS(gwServer.Handler)
	}
	return gwServer.ListenAndServe()
}

// preflightHandler adds the necessary headers in order to serve
// CORS from any origin using the methods "GET", "HEAD", "POST", "PUT", "DELETE"
// We insist, don't do this without consideration in production systems.
func preflightHandler(w http.ResponseWriter) {
	headers := []string{"Content-Type", "Accept", "Authorization"}
	w.Header().Set("Access-Control-Allow-Headers", strings.Join(headers, ","))
	methods := []string{"GET", "HEAD", "POST", "PUT", "DELETE"}
	w.Header().Set("Access-Control-Allow-Methods", strings.Join(methods, ","))
}

// allowCORS allows Cross Origin Resource Sharing from any origin.
// Don't do this without consideration in production systems.
func allowCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if origin := r.Header.Get("Origin"); origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			if r.Method == "OPTIONS" && r.Header.Get("Access-Control-Request-Method") != "" {
				preflightHandler(w)
				return
			}
		}
		h.ServeHTTP(w, r)
	})
}
