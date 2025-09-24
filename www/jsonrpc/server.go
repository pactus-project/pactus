package jsonrpc

import (
	"context"
	"fmt"
	"net"

	ret "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/logger"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
	"github.com/pacviewer/jrpc-gateway/jrpc"
	"github.com/rs/cors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Server struct {
	ctx      context.Context
	config   *Config
	listener net.Listener
	server   *jrpc.Server
	grpcConn *grpc.ClientConn
	logger   *logger.SubLogger
}

func NewServer(ctx context.Context, conf *Config) *Server {
	return &Server{
		ctx:    ctx,
		config: conf,
		logger: logger.NewSubLogger("_jsonrpc", nil),
	}
}

func (s *Server) StartServer(grpcServer string) error {
	if !s.config.Enable {
		return nil
	}

	dialOpts := make([]grpc.DialOption, 0)
	dialOpts = append(dialOpts,
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

	listener, err := util.NetworkListen(s.ctx, "tcp", s.config.Listen)
	if err != nil {
		return err
	}

	s.server = server
	s.listener = listener

	go func() {
		s.logger.Info("JSON-RPC server start listening", "address", listener.Addr())
		if err := server.Serve(listener); err != nil {
			s.logger.Debug("error on JSON-RPC server", "error", err)
		}
	}()

	return nil
}

func (s *Server) StopServer() {
	if s.server != nil {
		_ = s.server.GracefulStop(s.ctx)
		_ = s.listener.Close()
	}

	if s.grpcConn != nil {
		_ = s.grpcConn.Close()
	}
}
