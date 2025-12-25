package wallet

import (
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/wallet/storage"
	"github.com/pactus-project/pactus/wallet/types"
)

type transactions struct {
	storage    storage.IStorage
	grpcClient *grpcClient
}

func newTransactions(storage storage.IStorage,
	grpcClient *grpcClient,
) transactions {
	return transactions{
		storage:    storage,
		grpcClient: grpcClient,
	}
}

func (t *transactions) AddTransaction(txID tx.ID) error {
	idStr := txID.String()
	if t.storage.HasTransaction(idStr) {
		return ErrTransactionExists
	}

	res, err := t.grpcClient.getTransaction(txID)
	if err != nil {
		return err
	}

	trx, err := tx.FromString(res.Transaction.Data)
	if err != nil {
		return err
	}

	return t.addTransactionWithStatus(trx, types.TransactionStatusConfirmed)
}

func (t *transactions) addTransactionWithStatus(trx *tx.Tx, status types.TransactionStatus) error {
	txInfos, err := types.MakeTransactionInfos(trx)
	if err != nil {
		return err
	}

	for _, info := range txInfos {
		info.Status = status

		if t.storage.HasAddress(info.Sender) ||
			t.storage.HasAddress(info.Receiver) {
			if err := t.storage.InsertTransaction(info); err != nil {
				return err
			}
		}
	}

	return nil
}

// TxDirection indicates whether to include incoming or outgoing transactions.
type TxDirection int

const (
	// TxDirectionIncoming includes only transactions where the wallet receives funds.
	TxDirectionIncoming = 0
	// TxDirectionOutgoing includes only transactions where the wallet sends funds.
	TxDirectionOutgoing = 1
)

// listTransactionsConfig contains options for listing transactions.
type listTransactionsConfig struct {
	direction TxDirection
	address   string
	count     int
	skip      int
}

var defaultListTransactionsConfig = listTransactionsConfig{
	direction: TxDirectionIncoming,
	address:   "",
	count:     10,
	skip:      0,
}

// ListTransactionsOption is a functional option for ListTransactions.
type ListTransactionsOption func(*listTransactionsConfig)

// WithDirection filters transactions by direction (incoming or outgoing).
func WithDirection(dir TxDirection) ListTransactionsOption {
	return func(cfg *listTransactionsConfig) {
		cfg.direction = dir
	}
}

// WithAddress filters transactions by the specified address.
func WithAddress(address string) ListTransactionsOption {
	return func(cfg *listTransactionsConfig) {
		cfg.address = address
	}
}

// WithCount sets the maximum number of transactions to return.
func WithCount(count int) ListTransactionsOption {
	return func(cfg *listTransactionsConfig) {
		cfg.count = count
	}
}

// WithSkip sets the number of transactions to skip.
func WithSkip(skip int) ListTransactionsOption {
	return func(cfg *listTransactionsConfig) {
		cfg.skip = skip
	}
}

func (t *transactions) ListTransactions(addr string, opts ...ListTransactionsOption) []*types.TransactionInfo {
	cfg := defaultListTransactionsConfig
	for _, opt := range opts {
		opt(&cfg)
	}

	txs, _ := t.storage.ListTransactions(addr, cfg.count, cfg.skip)

	return txs
}
