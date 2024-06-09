package jsonrpc

import (
	"context"
	"fmt"
	"net"

	"github.com/pactus-project/pactus/util/logger"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
	"github.com/pacviewer/jrpc-gateway/jrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Server struct {
	ctx        context.Context
	cancel     context.CancelFunc
	config     *Config
	server     *jrpc.Server
	listener   net.Listener
	grpcClient *grpc.ClientConn
	logger     *logger.SubLogger
}

func NewServer(conf *Config) *Server {
	ctx, cancel := context.WithCancel(context.Background())

	return &Server{
		ctx:    ctx,
		cancel: cancel,
		config: conf,
		logger: logger.NewSubLogger("_jsonrpc", nil),
	}
}

func (s *Server) StartServer(grpcServer string) error {
	if !s.config.Enable {
		return nil
	}

	grpcConn, err := grpc.NewClient(
		grpcServer,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return fmt.Errorf("failed to dial server: %w", err)
	}

	s.grpcClient = grpcConn

	blockchainService := pactus.RegisterBlockchainJsonRPC(grpcConn)
	networkService := pactus.RegisterNetworkJsonRPC(grpcConn)
	transactionService := pactus.RegisterTransactionJsonRPC(grpcConn)
	walletService := pactus.RegisterWalletJsonRPC(grpcConn)

	server := jrpc.NewServer()
	server.RegisterServices(blockchainService, networkService, transactionService, walletService)

	listener, err := net.Listen("tcp", s.config.Listen)
	if err != nil {
		s.logger.Error("unable to establish tcp connection", "error", err)
	}

	go func() {
		for {
			select {
			case <-s.ctx.Done():
				return
			default:
				if err := server.Serve(listener); err != nil {
					s.logger.Error("error while establishing JSON-RPC connection", "error", err)
				}
			}
		}
	}()

	s.logger.Info("json-rpc started listening", "address", listener.Addr().String())
	s.server = server
	s.listener = listener

	return nil
}

func (s *Server) StopServer() {
	s.cancel()
	s.logger.Debug("context closed", "reason", s.ctx.Err())

	if s.server != nil {
		_ = s.server.GracefulStop(s.ctx)
		_ = s.listener.Close()
	}

	if s.grpcClient != nil {
		_ = s.grpcClient.Close()
	}
}
