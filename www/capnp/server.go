package capnp

import (
	"context"
	"net"

	"github.com/pactus-project/pactus/consensus"
	"github.com/pactus-project/pactus/state"
	"github.com/pactus-project/pactus/sync"
	"github.com/pactus-project/pactus/util/logger"
	"zombiezen.com/go/capnproto2/rpc"
)

type pactusServer struct {
	state     state.Facade
	sync      sync.Synchronizer
	consensus consensus.Reader
	logger    *logger.Logger
}

type Server struct {
	ctx       context.Context
	config    *Config
	address   string
	listener  net.Listener
	state     state.Facade
	sync      sync.Synchronizer
	consensus consensus.Reader
	logger    *logger.Logger
}

func NewServer(conf *Config, state state.Facade, sync sync.Synchronizer,
	consensus consensus.Reader) *Server {
	return &Server{
		ctx:       context.Background(),
		state:     state,
		sync:      sync,
		consensus: consensus,
		config:    conf,
		logger:    logger.NewLogger("_capnp", nil),
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
					s2c := PactusServer_ServerToClient(
						&pactusServer{s.state, s.sync, s.consensus, s.logger})
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
