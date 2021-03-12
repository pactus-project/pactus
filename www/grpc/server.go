package grpc

import (
	"context"
	"net"

	"github.com/zarbchain/zarb-go/logger"
	"github.com/zarbchain/zarb-go/state"
	"github.com/zarbchain/zarb-go/sync"
	zarb "github.com/zarbchain/zarb-go/www/grpc/proto"
	"google.golang.org/grpc"
)

type zarbServer struct {
	zarb.UnimplementedZarbServer

	state  state.StateFacade
	sync   sync.Synchronizer
	logger *logger.Logger
}

type Server struct {
	ctx      context.Context
	config   *Config
	listener net.Listener
	grpc     *grpc.Server
	state    state.StateFacade
	sync     sync.Synchronizer
	logger   *logger.Logger
}

func NewServer(conf *Config, state state.StateFacade, sync sync.Synchronizer) (*Server, error) {

	return &Server{
		ctx:    context.Background(),
		config: conf,
		state:  state,
		sync:   sync,
		logger: logger.NewLogger("_grpc", nil),
	}, nil
}

func (s *Server) StartServer() error {
	if !s.config.Enable {
		return nil
	}

	grpc := grpc.NewServer()
	server := &zarbServer{
		state:  s.state,
		sync:   s.sync,
		logger: s.logger,
	}
	zarb.RegisterZarbServer(grpc, server)

	listener, err := net.Listen("tcp", s.config.Address)
	if err != nil {
		return err
	}

	s.listener = listener
	s.grpc = grpc
	go func() {
		if err := s.grpc.Serve(listener); err != nil {
			s.logger.Error("Error on grpc serve", "err", err)
		}
	}()

	go func() {
		if err := s.startGateway(); err != nil {
			s.logger.Error("Error on grpc-gateway serve", "err", err)
		}
	}()

	return nil
}

func (s *Server) StopServer() {
	s.ctx.Done()

	if s.listener != nil {
		s.listener.Close()
	}

	if s.grpc != nil {
		s.grpc.Stop()
	}
}
