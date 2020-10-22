package capnp

import (
	"context"
	"net"

	"github.com/zarbchain/zarb-go/logger"
	"github.com/zarbchain/zarb-go/store"
	"github.com/zarbchain/zarb-go/txpool"
	"zombiezen.com/go/capnproto2/rpc"
)

type Server struct {
	ctx      context.Context
	config   *Config
	listener net.Listener
	store    store.StoreReader
	txPool   txpool.TxPoolReader
	logger   *logger.Logger
}

func NewServer(conf *Config, store store.StoreReader) (*Server, error) {
	return &Server{
		ctx:    context.Background(),
		store:  store,
		config: conf,
		logger: logger.NewLogger("_capnp", nil),
	}, nil
}

func (s *Server) StartServer() error {
	if !s.config.Enable {
		return nil
	}

	l, err := net.Listen("tcp", s.config.Address)
	if err != nil {
		return err
	}

	s.logger.Info("Capnp started listening", "address", l.Addr())
	s.listener = l
	go func() {
		for {
			defer func() {
				if r := recover(); r != nil {
					s.logger.Error("Recovered from a panic", r)
				}
			}()
			// Wait for a connection.
			conn, err := l.Accept()
			if err != nil {
				s.logger.Error("Error on accepting a connection", "error", err)
			}

			go func(c net.Conn) {
				s2c := ZarbServer_ServerToClient(factory{s.store, s.logger})
				conn := rpc.NewConn(rpc.StreamTransport(conn), rpc.MainInterface(s2c.Client))
				err := conn.Wait()
				if err != nil {
					s.logger.Error("Error on  a connection", "error", err)
				}

			}(conn)

			// TODO:
			// handle close signal/channel
		}
	}()

	return nil
}

func (s *Server) StopServer() error {
	s.listener.Close()

	return nil
}
