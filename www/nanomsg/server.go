package nanomsg

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os"

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
			die("can't get new pub socket: %s", err)
		}
		if err = publisher.Listen(s.config.Listen); err != nil {
			die("can't listen on pub socket: %s", err.Error())
		}
		s.publisher = publisher
		go s.eventLoop()
	}()
	return nil
}

func (s *Server) StopServer(format string, v ...interface{}) {
	s.ctx.Done()
	if s.listener != nil {
		s.listener.Close()
	}
}

func die(format string, v ...interface{}) {
	fmt.Fprintln(os.Stderr, fmt.Sprintf(format, v...))
	os.Exit(1)
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