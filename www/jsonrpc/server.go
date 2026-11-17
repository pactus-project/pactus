package jsonrpc

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	ret "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"github.com/pactus-project/gopkg/logger"
	"github.com/pactus-project/gopkg/netutil"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
	"github.com/pactus-project/pactus/www/guard"
	"github.com/pacviewer/jrpc-gateway/jrpc"
	"github.com/rs/cors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Server struct {
	ctx      context.Context
	config   *Config
	listener net.Listener
	address  string
	server   *jrpc.Server
	pipe     *pipeListener
	front    *http.Server
	grpcConn *grpc.ClientConn
	logger   *logger.SubLogger
}

func (s *Server) Address() string {
	return s.address
}

func NewServer(ctx context.Context, conf *Config) *Server {
	return &Server{
		ctx:    ctx,
		config: conf,
		logger: logger.NewSubLogger(ctx, "_jsonrpc", nil),
	}
}

func (s *Server) StartServer(grpcServer string) error {
	if !s.config.Enable {
		return nil
	}

	dialOpts := make([]grpc.DialOption, 0, 2)
	dialOpts = append(
		dialOpts,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(ret.UnaryClientInterceptor()),
	)
	grpcConn, err := grpc.NewClient(
		grpcServer,
		dialOpts...,
	)
	if err != nil {
		return fmt.Errorf("failed to dial server: %w", err)
	}

	s.grpcConn = grpcConn

	blockchainService := pactus.RegisterBlockchainJsonRPC(grpcConn)
	networkService := pactus.RegisterNetworkJsonRPC(grpcConn)
	transactionService := pactus.RegisterTransactionJsonRPC(grpcConn)
	walletService := pactus.RegisterWalletJsonRPC(grpcConn)
	utilsService := pactus.RegisterUtilsJsonRPC(grpcConn)

	opts := make([]jrpc.Option, 0)
	if len(s.config.Origins) > 0 {
		opts = append(opts, jrpc.WithCorsOrigins(&cors.Options{
			AllowedOrigins:   s.config.Origins,
			AllowedMethods:   []string{"POST"},
			AllowedHeaders:   []string{"*"},
			AllowCredentials: true,
		}))
	}

	server := jrpc.NewServer(opts...)
	server.RegisterServices(blockchainService, networkService, transactionService, walletService, utilsService)

	listener, err := netutil.NetworkListen(s.ctx, "tcp", s.config.Listen)
	if err != nil {
		return err
	}

	// The JSON-RPC library serves on a listener in memory. The front server listens on the
	// network: it checks the requests, with the same rules as the HTTP gateway, and it has
	// the timeouts that the library doesn't set. Then it passes the requests to the library.
	pipe := newPipeListener()
	proxy := httputil.NewSingleHostReverseProxy(&url.URL{Scheme: "http", Host: "jsonrpc"})
	proxy.Transport = &http.Transport{DialContext: pipe.dial}
	proxy.ErrorHandler = func(w http.ResponseWriter, _ *http.Request, err error) {
		var tooLarge *http.MaxBytesError
		if errors.As(err, &tooLarge) {
			http.Error(w, "request body too large", http.StatusRequestEntityTooLarge)

			return
		}
		http.Error(w, "bad gateway", http.StatusBadGateway)
	}

	front := &http.Server{
		Handler:           guard.Requests(proxy),
		ReadHeaderTimeout: 3 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      60 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	s.server = server
	s.pipe = pipe
	s.front = front
	s.listener = listener
	s.address = listener.Addr().String()

	go func() {
		if err := server.Serve(pipe); err != nil {
			s.logger.Debug("error on JSON-RPC backend", "error", err)
		}
	}()

	go func() {
		s.logger.Info("JSON-RPC server start listening", "address", listener.Addr())
		if err := front.Serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.logger.Debug("error on JSON-RPC server", "error", err)
		}
	}()

	return nil
}

func (s *Server) StopServer() {
	if s.front != nil {
		_ = s.front.Close()
	}

	if s.server != nil {
		_ = s.server.GracefulStop(s.ctx)
		_ = s.pipe.Close()
		_ = s.listener.Close()
	}

	if s.grpcConn != nil {
		_ = s.grpcConn.Close()
	}
}
