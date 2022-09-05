package grpc

import (
	"context"
	"net"

	"github.com/zarbchain/zarb-go/state"
	"github.com/zarbchain/zarb-go/sync"
	"github.com/zarbchain/zarb-go/util/logger"
	zarb "github.com/zarbchain/zarb-go/www/grpc/proto"
	"google.golang.org/grpc"
)

type Server struct {
	ctx      context.Context
	config   *Config
	listener net.Listener
	grpc     *grpc.Server
	state    state.Facade
	sync     sync.Synchronizer
	logger   *logger.Logger
}

func NewServer(conf *Config, state state.Facade, sync sync.Synchronizer) *Server {
	return &Server{
		ctx:    context.Background(),
		config: conf,
		state:  state,
		sync:   sync,
		logger: logger.NewLogger("_grpc", nil),
	}
}

func (s *Server) StartServer() error {
	if !s.config.Enable {
		return nil
	}

	grpc := grpc.NewServer()
	blockchainServer := &blockchainServer{
		state:  s.state,
		logger: s.logger,
	}
	transactionServer := &transactionServer{
		state:  s.state,
		logger: s.logger,
	}
	networkServer := &networkServer{
		sync:   s.sync,
		logger: s.logger,
	}
	zarb.RegisterBlockchainServer(grpc, blockchainServer)
	zarb.RegisterTransactionServer(grpc, transactionServer)
	zarb.RegisterNetworkServer(grpc, networkServer)

	listener, err := net.Listen("tcp", s.config.Listen)
	if err != nil {
		return err
	}

	s.listener = listener
	s.grpc = grpc
	go func() {
		if err := s.grpc.Serve(listener); err != nil {
			s.logger.Error("error on grpc serve", "err", err)
		}
	}()

	go func() {
		if err := s.startGateway(); err != nil {
			s.logger.Error("error on grpc-gateway serve", "err", err)
		}
	}()

	return nil
}

func (s *Server) StopServer() {
	s.ctx.Done()

	if s.grpc != nil {
		s.grpc.Stop()
	}

	if s.listener != nil {
		s.listener.Close()
	}
}
