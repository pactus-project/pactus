package network

import (
	"bufio"
	"context"

	lp2pcore "github.com/libp2p/go-libp2p-core"
	lp2phost "github.com/libp2p/go-libp2p-core/host"
	lp2pnetwork "github.com/libp2p/go-libp2p-core/network"
	lp2peer "github.com/libp2p/go-libp2p-core/peer"
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/logger"
	"github.com/zarbchain/zarb-go/util"
)

type streamService struct {
	ctx        context.Context
	host       lp2phost.Host
	protocolID lp2pcore.ProtocolID
	eventCh    chan NetworkEvent
	logger     *logger.Logger
}

func newStreamService(ctx context.Context, host lp2phost.Host, protocolID lp2pcore.ProtocolID,
	eventCh chan NetworkEvent, logger *logger.Logger) *streamService {
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
	reader := bufio.NewReader(stream)

	s.logger.Debug("Receiving stream", "from", util.FingerprintPeerID(from))
	event := &StreamMessage{
		Source: from,
		Reader: reader,
	}

	s.eventCh <- event
}

func (s *streamService) SendRequest(msg []byte, pid lp2peer.ID) error {
	s.logger.Debug("Sending stream", "to", util.FingerprintPeerID(pid))
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

	return nil
}
