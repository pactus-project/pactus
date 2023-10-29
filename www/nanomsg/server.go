package nanomsg

import (
	"bytes"
	"context"
	"net"

	"github.com/pactus-project/pactus/util/encoding"
	"github.com/pactus-project/pactus/util/logger"
	"github.com/pactus-project/pactus/www/nanomsg/event"
	mangos "go.nanomsg.org/mangos/v3"
	"go.nanomsg.org/mangos/v3/protocol/pub"

	// register nano ports transports.
	_ "go.nanomsg.org/mangos/v3/transport/all"
)

type Server struct {
	ctx       context.Context
	config    *Config
	publisher mangos.Socket
	listener  net.Listener
	logger    *logger.SubLogger
	eventCh   <-chan event.Event
	seqNum    uint32
}

func NewServer(conf *Config, eventCh <-chan event.Event) *Server {
	return &Server{
		ctx:     context.Background(),
		config:  conf,
		logger:  logger.NewSubLogger("_nonomsg", nil),
		eventCh: eventCh,
		seqNum:  0,
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
			s.logger.Error("error on nanomsg creating new socket", "error", err)
		}

		if err = publisher.Listen(s.config.Listen); err != nil {
			s.logger.Error("error on nanomsg publisher binding", "error", err)
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
			w := bytes.NewBuffer(e)

			err := encoding.WriteElement(w, s.seqNum)
			if err != nil {
				s.logger.Error("error on encoding event", "error", err)
				return
			}

			err = s.publisher.Send(w.Bytes())
			if err != nil {
				s.logger.Error("error on emitting event", "error", err)
				return
			}
			s.seqNum++
		}
	}
}
