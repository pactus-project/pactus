package zmq

import (
	"context"
	"net"

	"github.com/zarbchain/zarb-go/state"
	"github.com/zarbchain/zarb-go/sync"
	"github.com/zarbchain/zarb-go/util/logger"
	zarb "github.com/zarbchain/zarb-go/www/grpc/proto"
	"github.com/zeromq/goczmq"
)

type zarbServer struct {
	zarb.UnimplementedZarbServer
	state  state.Facade
	sync   sync.Synchronizer
	logger *logger.Logger
}

type Server struct {
	ctx      context.Context
	config   *Config
	listener net.Listener
	zmq      *goczmq.Sock
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
		logger: logger.NewLogger("_zeromq", nil),
	}
}

func (s *Server) StartServer() error {
	if !s.config.Enable {
		return nil
	}

	listener, err := net.Listen("tcp", s.config.Listen)
	if err != nil {
		return err
	}

	s.listener = listener
	go func() {
		if err := s.zmq.Connect(listener.Addr().String()); err != nil {
			s.logger.Error("error on grpc serve", "err", err)
		}
	}()

	return nil
}
func (s *Server) StopServer() {
	s.ctx.Done()
	if s.listener != nil {
		s.listener.Close()
	}
}
