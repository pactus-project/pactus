package zmq

import (
	"context"
	"encoding/json"
	"net"

	zmq "github.com/pebbe/zmq4"
	"github.com/zarbchain/zarb-go/util/logger"
	"github.com/zarbchain/zarb-go/www/zmq/event"
)

type Server struct {
	ctx       context.Context
	config    *Config
	publisher *zmq.Socket
	listener  net.Listener
	logger    *logger.Logger
	eventCh   <-chan event.Event
}

func NewServer(conf *Config, eventCh <-chan event.Event) *Server {
	return &Server{
		ctx:     context.Background(),
		config:  conf,
		logger:  logger.NewLogger("_zmq", nil),
		eventCh: eventCh,
	}
}

func (s *Server) Address() string {
	return s.listener.Addr().String()
}

func (s *Server) StartServer() error {
	if !s.config.Enable {
		return nil
	}
	go func() {
		ctx, err := zmq.NewContext()
		if err != nil {
			s.logger.Error("error on zmq context", "err", err)
		}
		publisher, err := ctx.NewSocket(zmq.PUB)
		if err != nil {
			s.logger.Error("error on creating new socket", "err", err)
		}
		err = publisher.Bind("tcp://*:5555")
		if err != nil {
			s.logger.Error("error on zmq publisher binding", "err", err)
		}
		s.publisher = publisher
		go s.eventLoop()
	}()
	return nil
}
func (s *Server) StopServer() {
	s.ctx.Done()

	if s.listener != nil {
		s.listener.Close()
	}
}

func (s *Server) eventLoop() {
	for {
		select {
		case <-s.ctx.Done():
			return

		case e := <-s.eventCh:
			bs, _ := json.Marshal(e)
			_, err := s.publisher.Send(string(bs), 0)
			if err != nil {
				s.logger.Error("error on emitting event", "err", err)
			}
		}
	}
}
