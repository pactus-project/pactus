package capnp

import (
	"context"
	"net"

	"github.com/zarbchain/zarb-go/state"
	"github.com/zarbchain/zarb-go/sync"
	"github.com/zarbchain/zarb-go/util/logger"
	"zombiezen.com/go/capnproto2/rpc"
)

type zarbServer struct {
	state  state.Facade
	sync   sync.Synchronizer
	logger *logger.Logger
}

type Server struct {
	ctx      context.Context
	config   *Config
	address  string
	listener net.Listener
	state    state.Facade
	sync     sync.Synchronizer
	logger   *logger.Logger
}

func NewServer(conf *Config, state state.Facade, sync sync.Synchronizer) *Server {
	return &Server{
		ctx:    context.Background(),
		state:  state,
		sync:   sync,
		config: conf,
		logger: logger.NewLogger("_capnp", nil),
	}
}

func (s *Server) Address() string {
	return s.address
}

func (s *Server) StartServer() error {
	if !s.config.Enable {
		return nil
	}

	l, err := net.Listen("tcp", s.config.Listen)
	if err != nil {
		return err
	}

	s.address = l.Addr().String()

	s.logger.Info("capnp started listening", "address", l.Addr())
	s.listener = l
	go func() {
		for {
			// Wait for a connection.
			conn, err := l.Accept()
			if err != nil {
				s.logger.Debug("error on accepting a connection", "err", err)
			} else {
				//
				go func(c net.Conn) {
					s2c := ZarbServer_ServerToClient(&zarbServer{s.state, s.sync, s.logger})
					conn := rpc.NewConn(rpc.StreamTransport(c), rpc.MainInterface(s2c.Client))
					err := conn.Wait()
					if err != nil {
						s.logger.Error("error on  a connection", "err", err)
					}
				}(conn)
			}

			// TODO:
			// handle close signal/channel
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
