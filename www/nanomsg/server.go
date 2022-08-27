package nanomsg

import (
	"context"
	"encoding/json"
	"net"

	"github.com/zarbchain/zarb-go/util/logger"
	"github.com/zarbchain/zarb-go/www/nanomsg/event"
	mangos "go.nanomsg.org/mangos/v3"
	"go.nanomsg.org/mangos/v3/protocol/pub"
)

type Server struct {
	ctx       context.Context
	config    *Config
	publisher mangos.Socket
	listener  net.Listener
	logger    *logger.Logger
	eventCh   <-chan event.Event
}

func NewServer(conf *Config, eventCh <-chan event.Event) *Server {
	return &Server{
		ctx:     context.Background(),
		config:  conf,
		logger:  logger.NewLogger("_nonomsg", nil),
		eventCh: eventCh,
	}
}

func (s *Server) StartServer() error {
	if !s.config.Enable {
		return nil
	}
	go func() {
		var publisher mangos.Socket
		var err error
		if publisher, err = pub.NewSocket(); err != nil {
			s.logger.Error("error on nanomsg creating new socket", "err", err)
		}
		if err = publisher.Listen(s.config.Listen); err != nil {
			s.logger.Error("error on nanomsg publisher binding", "err", err)
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
			err := s.publisher.Send((bs))
			if err != nil {
				s.logger.Error("error on emitting event", "err", err)
			}
		}
	}
}
