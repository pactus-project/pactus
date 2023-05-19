package network

import (
	"context"
	"fmt"

	lp2pcore "github.com/libp2p/go-libp2p/core"
	lp2phost "github.com/libp2p/go-libp2p/core/host"
	lp2pnetwork "github.com/libp2p/go-libp2p/core/network"
	lp2peer "github.com/libp2p/go-libp2p/core/peer"
	ma "github.com/multiformats/go-multiaddr"
	"github.com/pactus-project/pactus/util/errors"
	"github.com/pactus-project/pactus/util/logger"
)

type streamService struct {
	ctx        context.Context
	host       lp2phost.Host
	protocolID lp2pcore.ProtocolID
	relayAddrs []ma.Multiaddr
	eventCh    chan Event
	logger     *logger.Logger
}

func newStreamService(ctx context.Context, host lp2phost.Host,
	protocolID lp2pcore.ProtocolID, relayAddrs []ma.Multiaddr,
	eventCh chan Event, logger *logger.Logger) *streamService {
	s := &streamService{
		ctx:        ctx,
		host:       host,
		protocolID: protocolID,
		relayAddrs: relayAddrs,
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
	_, err := s.host.Peerstore().SupportsProtocols(pid, s.protocolID)
	if err != nil {
		return errors.Errorf(errors.ErrNetwork, err.Error())
	}
	stream, err := s.host.NewStream(
		lp2pnetwork.WithNoDial(s.ctx, "should already have connection"), pid, s.protocolID)
	if err != nil {
		s.logger.Warn("unable to open direct stream", "pid", pid, "err", err)

		circuitAddrs := make([]ma.Multiaddr, len(s.relayAddrs))
		for i, addr := range s.relayAddrs {
			circuitAddr, err := ma.NewMultiaddr(fmt.Sprintf("%s/p2p-circuit/p2p/%s", addr.String(), pid))
			if err != nil {
				return errors.Errorf(errors.ErrNetwork, err.Error())
			}
			//fmt.Println(circuitAddr)
			circuitAddrs[i] = circuitAddr
		}

		// Open a connection to the previously unreachable host via the relay address
		unreachableRelayInfo := lp2peer.AddrInfo{
			ID:    pid,
			Addrs: circuitAddrs,
		}

		if err := s.host.Connect(s.ctx, unreachableRelayInfo); err != nil {
			s.logger.Warn("unable to connect to peer using relay", "pid", pid, "err", err)
			return errors.Errorf(errors.ErrNetwork, err.Error())
		}

		stream, err = s.host.NewStream(
			lp2pnetwork.WithUseTransient(s.ctx, string(s.protocolID)), pid, s.protocolID)
		if err != nil {
			s.logger.Warn("unable to open relay stream", "pid", pid, "err", err)
			return errors.Errorf(errors.ErrNetwork, err.Error())
		}
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
