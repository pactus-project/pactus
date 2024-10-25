package grpc

import (
	"embed"
	"fmt"
	"io/fs"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

//go:embed swagger-ui
var swaggerFS embed.FS

type GatewayConfig struct {
	Enable     bool   `toml:"enable"`
	Listen     string `toml:"listen"`
	EnableCORS bool   `toml:"enable_cors"`
}

// getOpenAPIHandler serves an OpenAPI UI.
func (*Server) getOpenAPIHandler() (http.Handler, error) {
	if _, err := swaggerFS.ReadFile("swagger-ui/index.html"); err == nil {
		swagger, err := fs.Sub(swaggerFS, "swagger-ui")
		if err != nil {
			return nil, err
		}

		return http.FileServer(http.FS(swagger)), nil
	}

	return http.FileServer(http.Dir("swagger-ui")), nil
}

func (s *Server) startGateway(grpcAddr string) error {
	if !s.config.Gateway.Enable {
		return nil
	}

	conn, err := grpc.NewClient(
		grpcAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
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
	err = pactus.RegisterWalletHandler(s.ctx, gwMux, conn)
	if err != nil {
		return err
	}
	err = pactus.RegisterUtilsHandler(s.ctx, gwMux, conn)
	if err != nil {
		return err
	}

	handler, err := s.getOpenAPIHandler()
	if err != nil {
		return err
	}

	gwServer := &http.Server{
		Addr:              s.config.Gateway.Listen,
		ReadHeaderTimeout: 3 * time.Second,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasPrefix(r.URL.Path, "/pactus/") {
				gwMux.ServeHTTP(w, r)

				return
			}
			handler.ServeHTTP(w, r)
		}),
	}

	if s.config.Gateway.EnableCORS {
		gwServer.Handler = allowCORS(gwServer.Handler)
	}

	listener, err := net.Listen("tcp", s.config.Gateway.Listen)
	if err != nil {
		return err
	}

	s.logger.Info("grpc-gateway started listening", "address", listener.Addr().String())

	go func() {
		if err := gwServer.Serve(listener); err != nil {
			s.logger.Error("error on grpc-gateway serve", "error", err)
		}
	}()

	return nil
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
func allowCORS(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if origin := r.Header.Get("Origin"); origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			if r.Method == http.MethodOptions && r.Header.Get("Access-Control-Request-Method") != "" {
				preflightHandler(w)

				return
			}
		}
		handler.ServeHTTP(w, r)
	})
}
