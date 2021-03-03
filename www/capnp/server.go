package capnp

import (
	"context"
	"net"

	"github.com/zarbchain/zarb-go/logger"
	"github.com/zarbchain/zarb-go/state"
	"github.com/zarbchain/zarb-go/store"
	"github.com/zarbchain/zarb-go/sync"
	"github.com/zarbchain/zarb-go/txpool"
	"zombiezen.com/go/capnproto2/rpc"
)

type Server struct {
	ctx      context.Context
	config   *Config
	address  string
	listener net.Listener
	store    store.StoreReader
	state    state.StateReader
	txPool   txpool.TxPool
	sync     sync.Synchronizer
	logger   *logger.Logger
}

func NewServer(conf *Config, state state.StateReader, sync sync.Synchronizer, txPool txpool.TxPool) (*Server, error) {
	return &Server{
		ctx:    context.Background(),
		store:  state.StoreReader(),
		state:  state,
		txPool: txPool,
		sync:   sync,
		config: conf,
		logger: logger.NewLogger("_capnp", nil),
	}, nil
}
func (s *Server) Address() string {
	return s.address
}

func (s *Server) StartServer() error {
	if !s.config.Enable {
		return nil
	}

	l, err := net.Listen("tcp", s.config.Address)
	if err != nil {
		return err
	}

	s.address = l.Addr().String()

	s.logger.Info("Capnp started listening", "address", l.Addr())
	s.listener = l
	go func() {
		for {
			// Wait for a connection.
			conn, err := l.Accept()
			if err != nil {
				s.logger.Debug("Error on accepting a connection", "err", err)
			} else {
				//
				go func(c net.Conn) {
					s2c := ZarbServer_ServerToClient(zarbServer{s.state, s.store, s.txPool, s.sync, s.logger})
					conn := rpc.NewConn(rpc.StreamTransport(conn), rpc.MainInterface(s2c.Client))
					err := conn.Wait()
					if err != nil {
						s.logger.Error("Error on  a connection", "err", err)
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
