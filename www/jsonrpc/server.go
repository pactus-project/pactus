package jsonrpc

import (
	"context"
	"encoding/json"
	"fmt"
	"maps"
	"net"

	"github.com/pactus-project/pactus/util/logger"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
	"github.com/sourcegraph/jsonrpc2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Server struct {
	ctx        context.Context
	cancel     context.CancelFunc
	config     *Config
	grpcClient *grpc.ClientConn
	listener   net.Listener
	handlers   map[string]func(ctx context.Context, params json.RawMessage) (interface{}, error)
	logger     *logger.SubLogger
}

func NewServer(conf *Config, enableAuth bool) *Server {
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

	s.handlers = make(map[string]func(ctx context.Context, params json.RawMessage) (interface{}, error))
	blockchainService := pactus.NewBlockchainJsonRpcService(blockchain)
	networkService := pactus.NewNetworkJsonRpcService(network)

	maps.Copy(s.handlers, blockchainService.Methods())
	maps.Copy(s.handlers, networkService.Methods())

	listener, err := net.Listen("tcp", s.config.Listen)
	if err != nil {
		return err
	}

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				s.logger.Error("error accepting a connection", "error", err)

				continue
			}

			opts := []jsonrpc2.ConnOpt{
				jsonrpc2.LogMessages(s),
			}
			jsonrpc2.NewConn(
				s.ctx,
				jsonrpc2.NewBufferedStream(conn, jsonrpc2.VSCodeObjectCodec{}),
				jsonrpc2.AsyncHandler(s),
				opts...,
			)
		}
	}()

	s.logger.Info("jsonrpc started listening", "address", listener.Addr())
	s.listener = listener

	return nil
}

func (s *Server) StopServer() {
	s.cancel()
	s.logger.Debug("context closed", "reason", s.ctx.Err())

	if s.listener != nil {
		s.listener.Close()
	}

	if s.grpcClient != nil {
		s.grpcClient.Close()
	}
}

func (s *Server) Handle(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) {
	defer conn.Close()

	handler, ok := s.handlers[req.Method]
	if !ok {
		respErr := &jsonrpc2.Error{
			Code:    jsonrpc2.CodeMethodNotFound,
			Message: fmt.Sprintf("method not found: %v", req.Method)}
		replyErr := conn.ReplyWithError(ctx, req.ID, respErr)
		if replyErr != nil {
			s.logger.Error("failed to send error response", "error", replyErr)
		}

		return
	}
	params := json.RawMessage{}
	if req.Params != nil {
		params = *req.Params
	}
	res, err := handler(ctx, params)
	if err != nil {
		respErr := &jsonrpc2.Error{
			Code:    jsonrpc2.CodeInternalError,
			Message: err.Error()}
		replyErr := conn.ReplyWithError(ctx, req.ID, respErr)
		if replyErr != nil {
			s.logger.Error("failed to send error response", "error", replyErr)
		}

		return
	}

	replyErr := conn.Reply(ctx, req.ID, res)
	if replyErr != nil {
		s.logger.Error("failed to send response", "error", replyErr)
	}
}

func (s *Server) Printf(format string, v ...interface{}) {
	s.logger.Debug(format, v...)
}
