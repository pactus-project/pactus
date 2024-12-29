package zmq

import (
	"context"

	"github.com/go-zeromq/zmq4"
	"github.com/pactus-project/pactus/types/block"
	"github.com/pactus-project/pactus/util/logger"
)

type Server struct {
	sockets    map[string]zmq4.Socket
	publishers []Publisher
	config     *Config
	eventCh    <-chan any
	logger     *logger.SubLogger
}

func New(ctx context.Context, conf *Config, eventCh <-chan any) (*Server, error) {
	svr := &Server{
		eventCh:    eventCh,
		logger:     logger.NewSubLogger("_zmq", nil),
		publishers: make([]Publisher, 0),
		sockets:    make(map[string]zmq4.Socket),
		config:     conf,
	}

	publisherOpts := []zmq4.Option{
		zmq4.WithAutomaticReconnect(conf.ZmqAutomaticReconnect),
		zmq4.WithDialerRetry(conf.ZmqDialerRetryTime),
		zmq4.WithDialerMaxRetries(conf.ZmqDialerMaxRetries),
	}

	makePublisher := func(addr string, newPublisher func(socket zmq4.Socket, logger *logger.SubLogger) Publisher) error {
		if addr == "" {
			return nil
		}

		socket, ok := svr.sockets[addr]
		if !ok {
			socket = zmq4.NewPub(ctx, publisherOpts...)

			if err := socket.SetOption(zmq4.OptionHWM, conf.ZmqPubHWM); err != nil {
				return err
			}

			if err := socket.Listen(addr); err != nil {
				return err
			}
		}

		pub := newPublisher(socket, svr.logger)

		svr.publishers = append(svr.publishers, pub)

		svr.sockets[addr] = socket

		svr.logger.Info("publisher initialized", "topic", pub.TopicName(), "socket", addr)

		return nil
	}

	if err := makePublisher(conf.ZmqPubBlockInfo, newBlockInfoPub); err != nil {
		return nil, err
	}

	if err := makePublisher(conf.ZmqPubTxInfo, newTxInfoPub); err != nil {
		return nil, err
	}

	if err := makePublisher(conf.ZmqPubRawBlock, newRawBlockPub); err != nil {
		return nil, err
	}

	if err := makePublisher(conf.ZmqPubRawTx, newRawTxPub); err != nil {
		return nil, err
	}

	go svr.receivedEventLoop(ctx)

	return svr, nil
}

func (s *Server) Publishers() []Publisher {
	return s.publishers
}

func (s *Server) HWM() int {
	return s.config.ZmqPubHWM
}

func (s *Server) Close() {
	for _, sock := range s.sockets {
		if err := sock.Close(); err != nil {
			s.logger.Error("failed to close socket", "err", err)
		}
	}
}

func (s *Server) receivedEventLoop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case event, ok := <-s.eventCh:
			if !ok {
				s.logger.Warn("event channel closed")

				return
			}

			switch ev := event.(type) {
			case *block.Block:
				for _, pub := range s.publishers {
					pub.onNewBlock(ev)
				}
			default:
				panic("invalid event type")
			}
		}
	}
}
