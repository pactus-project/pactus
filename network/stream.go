package network

import (
	"context"
	"time"

	lp2pcore "github.com/libp2p/go-libp2p/core"
	lp2phost "github.com/libp2p/go-libp2p/core/host"
	lp2pnetwork "github.com/libp2p/go-libp2p/core/network"
	lp2peer "github.com/libp2p/go-libp2p/core/peer"
	"github.com/pactus-project/pactus/util/logger"
)

type streamService struct {
	ctx        context.Context
	host       lp2phost.Host
	protocolID lp2pcore.ProtocolID
	eventCh    chan Event
	logger     *logger.SubLogger
}

func newStreamService(ctx context.Context, host lp2phost.Host,
	protocolID lp2pcore.ProtocolID, eventCh chan Event, log *logger.SubLogger,
) *streamService {
	s := &streamService{
		ctx:        ctx,
		host:       host,
		protocolID: protocolID,
		eventCh:    eventCh,
		logger:     log,
	}

	s.host.SetStreamHandler(protocolID, s.handleStream)

	return s
}

func (*streamService) Start() {}

func (*streamService) Stop() {}

func (s *streamService) handleStream(stream lp2pnetwork.Stream) {
	from := stream.Conn().RemotePeer()

	s.logger.Trace("receiving stream", "from", from)
	event := &StreamMessage{
		From:   from,
		Reader: stream,
	}

	s.eventCh <- event
}

// SendRequest sends a message to a specific peer.
// If a direct connection can't be established, it attempts to connect via a relay node.
// Returns an error if the sending process fails.
func (s *streamService) SendRequest(msg []byte, pid lp2peer.ID) error {
	s.logger.Trace("sending stream", "to", pid)
	_, err := s.host.Peerstore().SupportsProtocols(pid, s.protocolID)
	if err != nil {
		return LibP2PError{Err: err}
	}

	// To prevent a broken stream from being open forever.
	ctxWithTimeout, cancel := context.WithTimeout(s.ctx, 20*time.Second)
	defer cancel()

	// Attempt to open a new stream to the target peer assuming there's already direct a connection
	stream, err := s.host.NewStream(
		lp2pnetwork.WithNoDial(ctxWithTimeout, "should already have connection"), pid, s.protocolID)
	if err != nil {
		return LibP2PError{Err: err}
	}

	deadline, _ := ctxWithTimeout.Deadline()
	_ = stream.SetDeadline(deadline)

	_, err = stream.Write(msg)
	if err != nil {
		_ = stream.Reset()

		return LibP2PError{Err: err}
	}

	err = stream.CloseWrite()
	if err != nil {
		return LibP2PError{Err: err}
	}

	return nil
}
