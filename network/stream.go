package network

import (
	"context"
	"fmt"

	lp2pcore "github.com/libp2p/go-libp2p/core"
	lp2phost "github.com/libp2p/go-libp2p/core/host"
	lp2pnetwork "github.com/libp2p/go-libp2p/core/network"
	lp2peer "github.com/libp2p/go-libp2p/core/peer"
	ma "github.com/multiformats/go-multiaddr"
	"github.com/pactus-project/pactus/util/logger"
)

type streamService struct {
	ctx        context.Context
	host       lp2phost.Host
	protocolID lp2pcore.ProtocolID
	relayAddrs []ma.Multiaddr
	eventCh    chan Event
	logger     *logger.SubLogger
}

func newStreamService(ctx context.Context, host lp2phost.Host,
	protocolID lp2pcore.ProtocolID, relayAddrs []ma.Multiaddr,
	eventCh chan Event, logger *logger.SubLogger,
) *streamService {
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

	s.logger.Debug("receiving stream", "from", from.ShortString())
	event := &StreamMessage{
		Source: from,
		Reader: stream,
	}

	s.eventCh <- event
}

// SendRequest sends a message to a specific peer.
// If a direct connection can't be established, it attempts to connect via a relay node.
// Returns an error if the sending process fails.
func (s *streamService) SendRequest(msg []byte, pid lp2peer.ID) error {
	s.logger.Trace("sending stream", "to", pid.ShortString())
	_, err := s.host.Peerstore().SupportsProtocols(pid, s.protocolID)
	if err != nil {
		return LibP2PError{Err: err}
	}

	// Attempt to open a new stream to the target peer assuming there's already direct a connection
	stream, err := s.host.NewStream(
		lp2pnetwork.WithNoDial(s.ctx, "should already have connection"), pid, s.protocolID)
	if err != nil {
		s.logger.Debug("unable to open direct stream", "pid", pid.ShortString(), "error", err)
		if len(s.relayAddrs) == 0 {
			return err
		}

		// We don't have a direct connection to the destination node,
		// so we try to connect via a relay node.
		// An example of a relay connection is described here:
		// https://github.com/libp2p/go-libp2p/blob/master/examples/relay/main.go
		circuitAddrs := make([]ma.Multiaddr, len(s.relayAddrs))
		for i, addr := range s.relayAddrs {
			// To connect a peer over relay, we need a relay address.
			// The format for the relay address is defined here:
			// https://docs.libp2p.io/concepts/nat/circuit-relay/#relay-addresses
			circuitAddr, err := ma.NewMultiaddr(fmt.Sprintf("%s/p2p-circuit/p2p/%s", addr.String(), pid))
			if err != nil {
				return LibP2PError{Err: err}
			}
			// fmt.Println(circuitAddr)
			circuitAddrs[i] = circuitAddr
		}

		// Open a connection to the previously unreachable host via the relay address
		unreachableRelayInfo := lp2peer.AddrInfo{
			ID:    pid,
			Addrs: circuitAddrs,
		}

		if err := s.host.Connect(s.ctx, unreachableRelayInfo); err != nil {
			// There is no relay connection to peer as well
			s.logger.Warn("unable to connect to peer using relay", "pid", pid.ShortString(), "error", err)
			return LibP2PError{Err: err}
		}
		s.logger.Debug("connected to peer using relay", "pid", pid.ShortString())

		// Try to open a new stream to the target peer using the relay connection.
		// The connection is marked as transient.
		stream, err = s.host.NewStream(
			lp2pnetwork.WithUseTransient(s.ctx, string(s.protocolID)), pid, s.protocolID)
		if err != nil {
			s.logger.Warn("unable to open relay stream", "pid", pid.ShortString(), "error", err)
			return LibP2PError{Err: err}
		}
	}

	_, err = stream.Write(msg)
	if err != nil {
		return LibP2PError{Err: err}
	}
	err = stream.Close()
	if err != nil {
		return LibP2PError{Err: err}
	}

	return nil
}
