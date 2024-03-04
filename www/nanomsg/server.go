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
	_ "go.nanomsg.org/mangos/v3/transport/all" // register nano ports transports.
)

type Server struct {
	ctx       context.Context
	cancel    context.CancelFunc
	config    *Config
	publisher mangos.Socket
	listener  net.Listener
	logger    *logger.SubLogger
	eventCh   <-chan event.Event
	seqNum    uint32
}

func NewServer(conf *Config, eventCh <-chan event.Event) *Server {
	ctx, cancel := context.WithCancel(context.Background())

	return &Server{
		ctx:     ctx,
		cancel:  cancel,
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
	publisher, err := pub.NewSocket()
	if err != nil {
		return err
	}
	listener, err := publisher.NewListener(s.config.Listen, nil)
	if err != nil {
		return err
	}
	err = listener.Listen()
	if err != nil {
		return err
	}

	s.publisher = publisher

	s.logger.Info("nanomsg started listening", "address", listener.Address())

	go s.eventLoop()

	return nil
}

func (s *Server) StopServer() {
	s.cancel()
	s.logger.Debug("context closed", "reason", s.ctx.Err())

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
