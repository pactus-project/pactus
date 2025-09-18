package http

import (
	"bytes"
	"context"
	"embed"
	"fmt"
	"io/fs"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/logger"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Server struct {
	ctx      context.Context
	config   *Config
	listener net.Listener
	server   *http.Server
	grpcConn *grpc.ClientConn
	logger   *logger.SubLogger
}

//go:embed swagger-ui
var swaggerFS embed.FS

// getOpenAPIHandler serves an OpenAPI UI.
func (s *Server) getOpenAPIHandler() (http.Handler, error) {
	swaggerTree, err := fs.Sub(swaggerFS, "swagger-ui")
	if err != nil {
		return nil, err
	}
	handler := http.FileServer(http.FS(swaggerTree))

	// Modify basePath in Swagger JSON file
	origContent, err := swaggerFS.ReadFile("swagger-ui/pactus.swagger.json")
	if err != nil {
		return nil, err
	}

	modifiedContent := bytes.Replace(
		origContent,
		[]byte(`"basePath": "/http/api"`),
		[]byte(fmt.Sprintf(`"basePath": %q`, s.patternToPrefix(s.config.apiPattern()))),
		1,
	)

	handler = s.changeBasePath(handler, modifiedContent)

	return handler, nil
}

func NewServer(ctx context.Context, conf *Config) *Server {
	return &Server{
		ctx:    ctx,
		config: conf,
		logger: logger.NewSubLogger("_http", nil),
	}
}

func (s *Server) StartServer(grpcAddr string) error {
	if !s.config.Enable {
		return nil
	}

	grpcConn, err := grpc.NewClient(
		grpcAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return fmt.Errorf("failed to dial server: %w", err)
	}

	s.grpcConn = grpcConn

	// gRPC-Gateway multiplexer
	gatewayMux := runtime.NewServeMux()
	if err := pactus.RegisterBlockchainHandler(s.ctx, gatewayMux, grpcConn); err != nil {
		return err
	}
	if err := pactus.RegisterTransactionHandler(s.ctx, gatewayMux, grpcConn); err != nil {
		return err
	}
	if err := pactus.RegisterNetworkHandler(s.ctx, gatewayMux, grpcConn); err != nil {
		return err
	}
	if err := pactus.RegisterWalletHandler(s.ctx, gatewayMux, grpcConn); err != nil {
		return err
	}
	if err := pactus.RegisterUtilsHandler(s.ctx, gatewayMux, grpcConn); err != nil {
		return err
	}

	// Swagger UI
	swaggerHandler, err := s.getOpenAPIHandler()
	if err != nil {
		return err
	}

	httpMux := http.NewServeMux()

	// Register gRPC-Gateway Handler at `/http/api`
	httpMux.Handle(s.config.apiPattern(),
		http.StripPrefix(s.patternToPrefix(s.config.apiPattern()), gatewayMux))

	// Register Swagger Handler at `/http/ui`
	httpMux.Handle(s.config.swaggerPattern(),
		http.StripPrefix(s.patternToPrefix(s.config.swaggerPattern()), swaggerHandler))

	// Redirect `/http` to `/http/ui`
	httpMux.HandleFunc(s.config.rootPattern(),
		func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, s.config.swaggerPattern(), http.StatusFound)
		})

	gwServer := &http.Server{
		Addr:              s.config.Listen,
		ReadHeaderTimeout: 3 * time.Second,
		Handler:           httpMux,
	}

	if s.config.EnableCORS {
		gwServer.Handler = allowCORS(gwServer.Handler)
	}

	listener, err := util.NetworkListen(s.ctx, "tcp", s.config.Listen)
	if err != nil {
		return err
	}

	s.server = gwServer
	s.listener = listener

	go func() {
		s.logger.Info("HTTP-API server start listening", "address", listener.Addr().String())
		if err := s.server.Serve(listener); err != nil && err != http.ErrServerClosed {
			s.logger.Debug("error on HTTP-API server", "error", err)
		}
	}()

	return nil
}

func (*Server) changeBasePath(handler http.Handler, modifiedContent []byte) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/pactus.swagger.json" {
			w.Header().Set("Content-Type", "application/json")
			_, err := w.Write(modifiedContent)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)

				return
			}
		} else {
			handler.ServeHTTP(w, r)
		}
	})
}

func (s *Server) StopServer() {
	if s.server != nil {
		_ = s.server.Close()
		_ = s.listener.Close()
	}

	if s.grpcConn != nil {
		_ = s.grpcConn.Close()
	}
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

// patternToPrefix removes the trailing '/' from the given pattern.
// Example: "/http/ui/" becomes "/http/ui".
func (*Server) patternToPrefix(pattern string) string {
	return pattern[:len(pattern)-1]
}
