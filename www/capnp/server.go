package capnp

import (
	"context"
	"net"

	"gitlab.com/zarb-chain/zarb-go/config"
	"gitlab.com/zarb-chain/zarb-go/logger"
	"gitlab.com/zarb-chain/zarb-go/store"
	"gitlab.com/zarb-chain/zarb-go/txpool"
	"zombiezen.com/go/capnproto2/rpc"
)

type Server struct {
	ctx      context.Context
	listener net.Listener
	store    store.StoreReader
	txPool   txpool.TxPoolReader
	config   *config.Config
	logger   *logger.Logger
}

func NewServer(store store.StoreReader, conf *config.Config) (*Server, error) {
	return &Server{
		ctx:    context.Background(),
		store:  store,
		config: conf,
		logger: logger.NewLogger("capnp", nil),
	}, nil
}

func (s *Server) StartServer() error {
	if !s.config.Capnp.Enable {
		return nil
	}

	l, err := net.Listen("tcp", s.config.Capnp.Address)
	if err != nil {
		return err
	}

	s.config.Capnp.Address = l.Addr().String()
	s.logger.Info("Capnp started listening", "address", s.config.Capnp.Address)
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
