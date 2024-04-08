package jsonrpc

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

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
	server     *http.Server
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
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return fmt.Errorf("failed to dial server: %w", err)
	}

	s.grpcClient = grpcConn

	blockchain := pactus.NewBlockchainClient(grpcConn)
	network := pactus.NewNetworkClient(grpcConn)
	transaction := pactus.NewTransactionClient(grpcConn)
	wallet := pactus.NewWalletClient(grpcConn)

	blockchainService := pactus.NewBlockchainJsonRpcService(blockchain)
	networkService := pactus.NewNetworkJsonRpcService(network)
	transactionService := pactus.NewTransactionJsonRpcService(transaction)
	walletService := pactus.NewWalletJsonRpcService(wallet)

	jgw := jrpc.NewServer()
	jgw.RegisterServices(&blockchainService, &networkService, &transactionService, &walletService)

	listener, err := net.Listen("tcp", s.config.Listen)
	if err != nil {
		s.logger.Error("unable to establish tcp connection", "error", err)
	}
	s.listener = listener

	mux := http.NewServeMux()
	mux.HandleFunc("/", jgw.HttpHandler)
	server := &http.Server{
		Addr:              listener.Addr().String(),
		ReadHeaderTimeout: 3 * time.Second,
		Handler:           mux,
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

	return nil
}

func (s *Server) StopServer() {
	s.cancel()
	s.logger.Debug("context closed", "reason", s.ctx.Err())

	if s.server != nil {
		_ = s.server.Shutdown(s.ctx)
		_ = s.server.Close()
		_ = s.listener.Close()
	}

	if s.grpcClient != nil {
		s.grpcClient.Close()
	}
}
