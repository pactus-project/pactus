package wallet

import (
	"github.com/pactus-project/pactus/types"
	"github.com/pactus-project/pactus/types/block"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/util/logger"
	"github.com/pactus-project/pactus/wallet/provider"
	"github.com/pactus-project/pactus/wallet/storage"
	wtypes "github.com/pactus-project/pactus/wallet/types"
)

type transactions struct {
	storage  storage.IStorage
	provider provider.IBlockchainProvider
}

func newTransactions(storage storage.IStorage,
	provider provider.IBlockchainProvider,
) transactions {
	return transactions{
		storage:  storage,
		provider: provider,
	}
}

// listTransactionsConfig contains options for listing transactions.
type listTransactionsConfig struct {
	direction wtypes.TxDirection
	address   string
	count     int
	skip      int
}

var defaultListTransactionsConfig = listTransactionsConfig{
	direction: wtypes.TxDirectionAny,
	address:   "*",
	count:     10,
	skip:      0,
}

// ListTransactionsOption is a functional option for ListTransactions.
type ListTransactionsOption func(*listTransactionsConfig)

// WithDirection filters transactions by direction (incoming or outgoing).
func WithDirection(dir wtypes.TxDirection) ListTransactionsOption {
	return func(cfg *listTransactionsConfig) {
		cfg.direction = dir
	}
}

// WithAddress filters transactions by the specified address.
func WithAddress(address string) ListTransactionsOption {
	return func(cfg *listTransactionsConfig) {
		if address != "" {
			cfg.address = address
		}
	}
}

// WithCount sets the maximum number of transactions to return.
func WithCount(count int) ListTransactionsOption {
	return func(cfg *listTransactionsConfig) {
		if count > 0 {
			cfg.count = count
		}
	}
}

// WithSkip sets the number of transactions to skip.
func WithSkip(skip int) ListTransactionsOption {
	return func(cfg *listTransactionsConfig) {
		if skip > 0 {
			cfg.skip = skip
		}
	}
}

func (t *transactions) ListTransactions(opts ...ListTransactionsOption) []*wtypes.TransactionInfo {
	if t.storage.IsLegacy() {
		return nil
	}

	cfg := defaultListTransactionsConfig
	for _, opt := range opts {
		opt(&cfg)
	}

	params := storage.QueryParams{
		Address:   cfg.address,
		Direction: cfg.direction,
		Count:     cfg.count,
		Skip:      cfg.skip,
	}

	txs, _ := t.storage.QueryTransactions(params)

	return txs
}

func (t *transactions) processEvent(event any) {
	if t.storage.IsLegacy() {
		return
	}

	switch evt := event.(type) {
	case *block.Block:
		t.processBlock(evt)
	case *tx.Tx:
		t.processTransaction(evt)
	default:
		// ignore other events
	}
}

func (t *transactions) processBlock(blk *block.Block) {
	pendingTxs := t.getPendingTransaction()

	for _, trx := range blk.Transactions() {
		txID := trx.ID().String()

		if _, ok := pendingTxs[txID]; ok {
			pendingInfo := pendingTxs[txID]
			if err := t.storage.UpdateTransactionStatus(pendingInfo.No,
				wtypes.TransactionStatusConfirmed, blk.Height()); err != nil {
				logger.Warn("failed to update transaction status", "error", err, "id", txID)
			}

			logger.Info("confirmed pending transaction", "id", trx.ID())

			continue
		}

		_ = t.addTransactionWithStatus(trx, wtypes.TransactionStatusConfirmed, blk.Height())
	}
}

func (t *transactions) processTransaction(trx *tx.Tx) {
	_ = t.addTransactionWithStatus(trx, wtypes.TransactionStatusPending, 0)
}

func (t *transactions) getPendingTransaction() map[string]*wtypes.TransactionInfo {
	pendingTxs, err := t.storage.GetPendingTransactions()
	if err != nil {
		logger.Warn("failed to get pending transactions", "error", err)

		return nil
	}

	for _, pendingInfo := range pendingTxs {
		trx, err := tx.FromBytes(pendingInfo.Data)
		if err != nil {
			continue
		}

		// TODO: check for expired and failed transactions

		// Re-broadcast the transaction
		_, err = t.provider.SendTx(trx)
		if err != nil {
			continue
		}
	}

	return pendingTxs
}

func (t *transactions) AddTransactionByID(txID tx.ID) error {
	if t.storage.IsLegacy() {
		return nil
	}

	idStr := txID.String()
	if t.storage.HasTransaction(idStr) {
		return ErrTransactionExists
	}

	trx, height, err := t.provider.GetTransaction(idStr)
	if err != nil {
		return err
	}

	return t.addTransactionWithStatus(trx, wtypes.TransactionStatusConfirmed, height)
}

func (t *transactions) addTransactionWithStatus(trx *tx.Tx,
	status wtypes.TransactionStatus, blockHeight types.Height,
) error {
	txInfos, err := wtypes.MakeTransactionInfos(trx, status, blockHeight)
	if err != nil {
		return err
	}

	for _, info := range txInfos {
		if t.storage.HasAddress(info.Sender) {
			info.Direction = wtypes.TxDirectionOutgoing
		} else if t.storage.HasAddress(info.Receiver) {
			info.Direction = wtypes.TxDirectionIncoming
		} else {
			continue
		}

		if err := t.storage.InsertTransaction(info); err != nil {
			logger.Warn("failed to insert transaction into storage", "error", err, "id", trx.ID())

			return err
		}

		logger.Info("added outgoing transaction to wallet",
			"id", trx.ID(), "status", status, "blockHeight", blockHeight)
	}

	return nil
}
