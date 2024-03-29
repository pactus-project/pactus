package jsonrpc

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/pacviewer/jrpc-gateway/jrpc"

	"github.com/pactus-project/pactus/util/logger"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Server struct {
	ctx        context.Context
	cancel     context.CancelFunc
	config     *Config
	server     *http.Server
	grpcClient *grpc.ClientConn
	logger     *logger.SubLogger
}

func NewServer(conf *Config, _ bool) *Server {
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

	grpcConn, err := grpc.DialContext(
		s.ctx,
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

	mux := http.NewServeMux()
	mux.HandleFunc("/", jgw.HttpHandler)
	server := &http.Server{
		Addr:              s.config.Listen,
		ReadHeaderTimeout: 3 * time.Second,
		Handler:           mux,
	}
	go func() {
		for {
			select {
			case <-s.ctx.Done():
				return
			default:
				if err := server.ListenAndServe(); err != nil {
					s.logger.Error("error while establishing JSON-RPC connection", "error", err)
				}
			}
		}
	}()
	s.logger.Info("jsonrpc started listening", "address", s.config.Listen)
	s.server = server

	return nil
}

func (s *Server) StopServer() {
	s.cancel()
	s.logger.Debug("context closed", "reason", s.ctx.Err())

	if s.server != nil {
		s.server.Close()
	}

	if s.grpcClient != nil {
		s.grpcClient.Close()
	}
}

func (s *Server) Printf(format string, v ...interface{}) {
	s.logger.Debug(format, v...)
}
