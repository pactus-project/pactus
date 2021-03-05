package grpc

import (
	"context"
	"net"

	"github.com/zarbchain/zarb-go/logger"
	"github.com/zarbchain/zarb-go/state"
	"github.com/zarbchain/zarb-go/store"
	"github.com/zarbchain/zarb-go/sync"
	"github.com/zarbchain/zarb-go/txpool"
	zarb "github.com/zarbchain/zarb-go/www/grpc/proto"
	"google.golang.org/grpc"
)

type zarbServer struct {
	zarb.UnimplementedZarbServer

	state  state.StateReader
	store  store.StoreReader
	txPool txpool.TxPool
	sync   sync.Synchronizer
	logger *logger.Logger
}

type Server struct {
	ctx      context.Context
	config   *Config
	listener net.Listener
	grpc     *grpc.Server
	store    store.StoreReader
	state    state.StateReader
	txPool   txpool.TxPool
	sync     sync.Synchronizer
	logger   *logger.Logger
}

func NewServer(conf *Config, state state.StateReader, sync sync.Synchronizer, txPool txpool.TxPool) (*Server, error) {

	return &Server{
		ctx:    context.Background(),
		config: conf,
		store:  state.StoreReader(),
		state:  state,
		txPool: txPool,
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
		store:  s.store,
		txPool: s.txPool,
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
