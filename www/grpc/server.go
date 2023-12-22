package grpc

import (
	"context"
	"net"

	"github.com/pactus-project/pactus/consensus"
	"github.com/pactus-project/pactus/network"
	"github.com/pactus-project/pactus/state"
	"github.com/pactus-project/pactus/sync"
	"github.com/pactus-project/pactus/util/logger"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
	"google.golang.org/grpc"
)

type Server struct {
	ctx      context.Context
	config   *Config
	listener net.Listener
	address  string
	grpc     *grpc.Server
	state    state.Facade
	net      network.Network
	sync     sync.Synchronizer
	consMgr  consensus.ManagerReader
	logger   *logger.SubLogger
}

func NewServer(conf *Config, st state.Facade, syn sync.Synchronizer,
	n network.Network, consMgr consensus.ManagerReader,
) *Server {
	return &Server{
		ctx:     context.Background(),
		config:  conf,
		state:   st,
		sync:    syn,
		net:     n,
		consMgr: consMgr,
		logger:  logger.NewSubLogger("_grpc", nil),
	}
}

func (s *Server) Address() string {
	return s.address
}

func (s *Server) StartServer() error {
	if !s.config.Enable {
		return nil
	}

	listener, err := net.Listen("tcp", s.config.Listen)
	if err != nil {
		return err
	}

	return s.startListening(listener)
}

func (s *Server) startListening(listener net.Listener) error {
	grpcServer := grpc.NewServer()
	blockchainServer := &blockchainServer{s}
	transactionServer := &transactionServer{s}
	networkServer := &networkServer{s}
	walletServer := &walletServer{
		Server:    s,
		wallets:   make(map[string]*loadedWallet),
		chainType: s.state.Genesis().ChainType(),
	}
	pactus.RegisterBlockchainServer(grpcServer, blockchainServer)
	pactus.RegisterTransactionServer(grpcServer, transactionServer)
	pactus.RegisterNetworkServer(grpcServer, networkServer)
	pactus.RegisterWalletServer(grpcServer, walletServer)

	s.listener = listener
	s.address = listener.Addr().String()
	s.grpc = grpcServer
	go func() {
		s.logger.Info("grpc server started", "addr", listener.Addr())
		if err := s.grpc.Serve(listener); err != nil {
			s.logger.Error("error on grpc serve", "error", err)
		}
	}()

	go func() {
		if err := s.startGateway(listener.Addr().String()); err != nil {
			s.logger.Error("error on grpc-gateway serve", "error", err)
		}
	}()

	return nil
}

func (s *Server) StopServer() {
	s.ctx.Done()

	if s.grpc != nil {
		s.grpc.Stop()
	}

	if s.listener != nil {
		s.listener.Close()
	}
}
