package zmq

import (
	"context"
	"log"
	"net"

	"github.com/gorilla/mux"
	"github.com/zarbchain/zarb-go/util/logger"
	zmq "github.com/zeromq/goczmq"
)

type Server struct {
	ctx      context.Context
	config   *Config
	router   *mux.Router
	listener net.Listener
	logger   *logger.Logger
}
func NewServer(conf *Config) *Server {
	return &Server{
		ctx:    context.Background(),
		config: conf,
		logger: logger.NewLogger("_zmq", nil),
	}
}

func (s *Server) Address() string {
	return s.listener.Addr().String()
}

func (s *Server) StartServer() error {
	if !s.config.Enable {
		return nil
	}
	con, err := net.Dial("tcp", s.config.Listen)
	if err != nil {
		return err
	}
	s.logger.Info("zmq started listening", "address", con)
	go func() {
		router,err := zmq.NewRouter(con.LocalAddr().String());
		if  err != nil {
			s.logger.Error("error on zmq serve", "err", err)
		}
		defer router.Destroy()
		log.Println("router created and bound")
	}()

	return nil
}
func (s *Server) StopServer() {
	s.ctx.Done()

	if s.listener != nil {
		s.listener.Close()
	}
}
