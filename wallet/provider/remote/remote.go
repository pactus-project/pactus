package remoteprovider

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net"
	"time"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/genesis"
	"github.com/pactus-project/pactus/types/account"
	"github.com/pactus-project/pactus/types/block"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/validator"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/wallet/provider"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var _ provider.IBlockchainProvider = (*RemoteBlockchainProvider)(nil)

type remoteProviderConfig struct {
	network genesis.ChainType
	timeout time.Duration
	servers []string
}

var defaultOpenWalletConfig = remoteProviderConfig{
	network: genesis.Mainnet,
	timeout: 5 * time.Second,
	servers: make([]string, 0),
}

type RemoteProviderOption func(*remoteProviderConfig)

func WithNetwork(network genesis.ChainType) RemoteProviderOption {
	return func(cfg *remoteProviderConfig) {
		cfg.network = network
	}
}

func WithTimeout(timeout time.Duration) RemoteProviderOption {
	return func(cfg *remoteProviderConfig) {
		cfg.timeout = timeout
	}
}

func WithCustomServers(servers []string) RemoteProviderOption {
	return func(cfg *remoteProviderConfig) {
		cfg.servers = servers
	}
}

// RemoteBlockchainProvider is a blockchain provider that connects to a remote gRPC server.
// It randomly selects a server from a predefined list.
type RemoteBlockchainProvider struct {
	ctx               context.Context
	servers           []string
	conn              *grpc.ClientConn
	timeout           time.Duration
	blockchainClient  pactus.BlockchainClient
	transactionClient pactus.TransactionClient
}

func NewRemoteBlockchainProvider(ctx context.Context, opts ...RemoteProviderOption) (*RemoteBlockchainProvider, error) {
	cfg := defaultOpenWalletConfig
	for _, opt := range opts {
		opt(&cfg)
	}

	serversData := map[string][]ServerInfo{}
	err := json.Unmarshal(serversJSON, &serversData)
	if err != nil {
		return nil, err
	}

	var servers []string
	switch cfg.network {
	case genesis.Mainnet:
		for _, srv := range serversData["mainnet"] {
			servers = append(servers, srv.Address)
		}

	case genesis.Testnet:
		crypto.ToTestnetHRP()

		for _, srv := range serversData["testnet"] {
			servers = append(servers, srv.Address)
		}

	case genesis.Localnet:
		crypto.ToTestnetHRP()

		servers = []string{"localhost:50052"}

	default:
		return nil, ErrInvalidNetwork
	}

	util.Shuffle(servers)

	return &RemoteBlockchainProvider{
		ctx:     ctx,
		servers: servers,
		timeout: cfg.timeout,
	}, nil
}

func (p *RemoteBlockchainProvider) connect() error {
	if p.conn != nil {
		return nil
	}

	for _, server := range p.servers {
		conn, err := grpc.NewClient(server,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithContextDialer(func(ctx context.Context, address string) (net.Conn, error) {
				return util.NetworkDialTimeout(ctx, "tcp", address, p.timeout)
			}))
		if err != nil {
			continue
		}

		blockchainClient := pactus.NewBlockchainClient(conn)
		transactionClient := pactus.NewTransactionClient(conn)

		// Check if client is responding
		// TODO: Use Ping API in version 1.11.0
		_, err = blockchainClient.GetBlockchainInfo(p.ctx,
			&pactus.GetBlockchainInfoRequest{})
		if err != nil {
			_ = conn.Close()

			continue
		}

		p.conn = conn
		p.blockchainClient = blockchainClient
		p.transactionClient = transactionClient

		return nil
	}

	return errors.New("unable to connect to the servers")
}

func (p *RemoteBlockchainProvider) Close() error {
	if p.conn != nil {
		return p.conn.Close()
	}

	return nil
}

func (p *RemoteBlockchainProvider) LastBlockHeight() (block.Height, error) {
	if err := p.connect(); err != nil {
		return 0, err
	}

	ctx, cancel := context.WithTimeout(p.ctx, p.timeout)
	defer cancel()

	res, err := p.blockchainClient.GetBlockchainInfo(ctx,
		&pactus.GetBlockchainInfoRequest{})
	if err != nil {
		return 0, err
	}

	return block.Height(res.LastBlockHeight), nil
}

func (p *RemoteBlockchainProvider) GetAccount(addrStr string) (*account.Account, error) {
	if err := p.connect(); err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(p.ctx, p.timeout)
	defer cancel()

	res, err := p.blockchainClient.GetAccount(ctx,
		&pactus.GetAccountRequest{Address: addrStr})
	if err != nil {
		return nil, err
	}

	data, err := hex.DecodeString(res.Account.Data)
	if err != nil {
		return nil, err
	}

	return account.FromBytes(data)
}

func (p *RemoteBlockchainProvider) GetValidator(addrStr string) (*validator.Validator, error) {
	if err := p.connect(); err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(p.ctx, p.timeout)
	defer cancel()

	res, err := p.blockchainClient.GetValidator(ctx,
		&pactus.GetValidatorRequest{Address: addrStr})
	if err != nil {
		return nil, err
	}

	data, err := hex.DecodeString(res.Validator.Data)
	if err != nil {
		return nil, err
	}

	return validator.FromBytes(data)
}

func (p *RemoteBlockchainProvider) SendTx(trx *tx.Tx) (string, error) {
	if err := p.connect(); err != nil {
		return "", err
	}

	data, err := trx.Bytes()
	if err != nil {
		return "", err
	}

	ctx, cancel := context.WithTimeout(p.ctx, p.timeout)
	defer cancel()

	res, err := p.transactionClient.BroadcastTransaction(ctx,
		&pactus.BroadcastTransactionRequest{SignedRawTransaction: hex.EncodeToString(data)})
	if err != nil {
		return "", err
	}

	return res.Id, nil
}

func (p *RemoteBlockchainProvider) GetTransaction(txID string) (*tx.Tx, block.Height, error) {
	if err := p.connect(); err != nil {
		return nil, 0, err
	}

	ctx, cancel := context.WithTimeout(p.ctx, p.timeout)
	defer cancel()

	res, err := p.transactionClient.GetTransaction(ctx,
		&pactus.GetTransactionRequest{
			Id:        txID,
			Verbosity: pactus.TransactionVerbosity_TRANSACTION_VERBOSITY_DATA,
		})
	if err != nil {
		return nil, 0, err
	}

	data, err := hex.DecodeString(res.Transaction.Data)
	if err != nil {
		return nil, 0, err
	}

	tx, err := tx.FromBytes(data)
	if err != nil {
		return nil, 0, err
	}

	return tx, block.Height(res.BlockHeight), nil
}
