package network

import (
	"context"
	"time"

	lp2pcore "github.com/libp2p/go-libp2p/core"
	lp2phost "github.com/libp2p/go-libp2p/core/host"
	lp2pnetwork "github.com/libp2p/go-libp2p/core/network"
	lp2peer "github.com/libp2p/go-libp2p/core/peer"
	"github.com/pactus-project/pactus/util/logger"
	"github.com/pactus-project/pactus/util/pipeline"
)

type streamService struct {
	ctx         context.Context
	host        lp2phost.Host
	protocolID  lp2pcore.ProtocolID
	timeout     time.Duration
	networkPipe pipeline.Pipeline[Event]
	logger      *logger.SubLogger
}

func newStreamService(ctx context.Context, host lp2phost.Host, conf *Config,
	protocolID lp2pcore.ProtocolID, networkPipe pipeline.Pipeline[Event], log *logger.SubLogger,
) *streamService {
	service := &streamService{
		ctx:         ctx,
		host:        host,
		protocolID:  protocolID,
		timeout:     conf.StreamTimeout,
		networkPipe: networkPipe,
		logger:      log,
	}

	service.host.SetStreamHandler(protocolID, service.handleStream)

	return service
}

func (*streamService) Start() {}

func (*streamService) Stop() {}

func (s *streamService) handleStream(stream lp2pnetwork.Stream) {
	from := stream.Conn().RemotePeer()

	s.logger.Debug("receiving stream", "from", from)
	event := &StreamMessage{
		From:   from,
		Reader: stream,
	}

	s.networkPipe.Send(event)
}

// SendRequest sends a message to a specific peer, assuming there is already a direct connection.
//
// For simplicity, we do not use bi-directional streams.
// Each time a peer wants to send a message, it creates a new stream.
//
// For more details on stream multiplexing, refer to: https://docs.libp2p.io/concepts/multiplex/overview/
func (s *streamService) SendRequest(msg []byte, pid lp2peer.ID) (lp2pnetwork.Stream, error) {
	s.logger.Trace("sending stream", "to", pid)
	_, err := s.host.Peerstore().SupportsProtocols(pid, s.protocolID)
	if err != nil {
		return nil, LibP2PError{Err: err}
	}

	// To prevent a broken stream from being open forever.
	ctxWithTimeout, cancel := context.WithTimeout(s.ctx, 5*time.Second)
	defer cancel()

	// Attempt to open a new stream to the peer, assuming there's already a direct connection.
	stream, err := s.host.NewStream(
		lp2pnetwork.WithNoDial(ctxWithTimeout, "should already have connection"), pid, s.protocolID)
	if err != nil {
		return nil, LibP2PError{Err: err}
	}

	_, err = stream.Write(msg)
	if err != nil {
		_ = stream.Reset()

		return nil, LibP2PError{Err: err}
	}

	err = stream.CloseWrite()
	if err != nil {
		return nil, LibP2PError{Err: err}
	}

	// We need to close the stream once it is read by the receiver.
	// If, for any reason, the receiver doesn't close the stream, we need to close it after a timeout.
	go func() {
		timer := time.NewTimer(s.timeout)
		closed := make(chan bool)

		go func() {
			// We need only one byte to read the EOF.
			buf := make([]byte, 1)
			_, _ = stream.Read(buf)
			closed <- true
		}()

		select {
		case <-timer.C:
			s.logger.Warn("stream timeout", "to", pid)
			_ = stream.Close()

		case <-closed:
			s.logger.Debug("stream closed", "to", pid)
			_ = stream.Close()
		}
	}()

	return stream, nil
}
