package network

import (
	"context"

	lp2pcore "github.com/libp2p/go-libp2p-core"
	lp2phost "github.com/libp2p/go-libp2p-core/host"
	lp2pnetwork "github.com/libp2p/go-libp2p-core/network"
	lp2peer "github.com/libp2p/go-libp2p-core/peer"
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/logger"
)

type streamService struct {
	ctx        context.Context
	host       lp2phost.Host
	protocolID lp2pcore.ProtocolID
	eventCh    chan Event
	logger     *logger.Logger
}

func newStreamService(ctx context.Context, host lp2phost.Host, protocolID lp2pcore.ProtocolID,
	eventCh chan Event, logger *logger.Logger) *streamService {
	s := &streamService{
		ctx:        ctx,
		host:       host,
		protocolID: protocolID,
		eventCh:    eventCh,
		logger:     logger,
	}

	s.host.SetStreamHandler(protocolID, s.handleStream)
	return s
}

func (s *streamService) Start() {
}
func (s *streamService) Stop() {
}

func (s *streamService) handleStream(stream lp2pnetwork.Stream) {
	from := stream.Conn().RemotePeer()

	s.logger.Debug("receiving stream", "from", from)
	event := &StreamMessage{
		Source: from,
		Reader: stream,
	}

	s.eventCh <- event
}

func (s *streamService) SendRequest(msg []byte, pid lp2peer.ID) error {
	s.logger.Debug("sending stream", "to", pid)
	_, err := s.host.Peerstore().SupportsProtocols(pid, string(s.protocolID))
	if err != nil {
		return errors.Errorf(errors.ErrNetwork, err.Error())
	}
	stream, err := s.host.NewStream(
		lp2pnetwork.WithNoDial(s.ctx, "should already have connection"),
		pid,
		s.protocolID)
	if err != nil {
		return errors.Errorf(errors.ErrNetwork, err.Error())
	}

	_, err = stream.Write(msg)
	if err != nil {
		return errors.Errorf(errors.ErrNetwork, err.Error())
	}
	err = stream.Close()
	if err != nil {
		return errors.Errorf(errors.ErrNetwork, err.Error())
	}

	return nil
}
