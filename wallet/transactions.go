package wallet

import (
	"github.com/pactus-project/pactus/types/block"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/tx/payload"
	"github.com/pactus-project/pactus/util/logger"
	"github.com/pactus-project/pactus/wallet/provider"
	"github.com/pactus-project/pactus/wallet/storage"
	"github.com/pactus-project/pactus/wallet/types"
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

func (t *transactions) AddTransaction(txID tx.ID) error {
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

	return t.addTransactionWithStatus(trx, types.TransactionStatusConfirmed, height)
}

func (t *transactions) addTransactionWithStatus(trx *tx.Tx, status types.TransactionStatus,
	blockHeight block.Height,
) error {
	txInfos, err := types.MakeTransactionInfos(trx, status, blockHeight)
	if err != nil {
		return err
	}

	for _, info := range txInfos {
		if t.storage.HasAddress(info.Sender) {
			info.Direction = types.TxDirectionOutgoing
		} else if t.storage.HasAddress(info.Receiver) {
			info.Direction = types.TxDirectionIncoming
		} else {
			continue
		}

		if err := t.storage.InsertTransaction(info); err != nil {
			return err
		}
	}

	return nil
}

// listTransactionsConfig contains options for listing transactions.
type listTransactionsConfig struct {
	direction types.TxDirection
	address   string
	count     int
	skip      int
}

var defaultListTransactionsConfig = listTransactionsConfig{
	direction: types.TxDirectionAny,
	address:   "*",
	count:     10,
	skip:      0,
}

// ListTransactionsOption is a functional option for ListTransactions.
type ListTransactionsOption func(*listTransactionsConfig)

// WithDirection filters transactions by direction (incoming or outgoing).
func WithDirection(dir types.TxDirection) ListTransactionsOption {
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

func (t *transactions) ListTransactions(opts ...ListTransactionsOption) []*types.TransactionInfo {
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
	default:
		// ignore other events
	}
}

func collectReceivers(trx *tx.Tx) []string {
	switch pld := trx.Payload().(type) {
	case *payload.TransferPayload:
		return []string{pld.To.String()}
	case *payload.BondPayload:
		return []string{pld.To.String()}
	case *payload.UnbondPayload:
		return []string{pld.Validator.String()}
	case *payload.WithdrawPayload:
		return []string{pld.To.String()}
	case *payload.SortitionPayload:
		return nil
	case *payload.BatchTransferPayload:
		receivers := make([]string, 0, len(pld.Recipients))
		for _, recipient := range pld.Recipients {
			receivers = append(receivers, recipient.To.String())
		}

		return receivers
	default:
		return nil
	}
}

func (t *transactions) addIncomingIfWallet(trx *tx.Tx, receivers []string, blockHeight block.Height) bool {
	for _, receiver := range receivers {
		if !t.storage.HasAddress(receiver) {
			continue
		}

		if err := t.addTransactionWithStatus(trx, types.TransactionStatusConfirmed, blockHeight); err != nil {
			logger.Warn("failed to add incoming transaction to wallet", "error", err, "id", trx.ID())

			continue
		}

		logger.Info("added incoming transaction to wallet", "id", trx.ID())

		return true
	}

	return false
}

func (t *transactions) processBlock(blk *block.Block) {
	pendingTxs, err := t.storage.GetPendingTransactions()
	if err != nil {
		logger.Warn("failed to get pending transactions", "error", err)

		return
	}

	for txID, pendingInfo := range pendingTxs {
		trx, err := tx.FromBytes(pendingInfo.Data)
		if err != nil {
			logger.Warn("failed to deserialize transaction", "error", err, "id", txID)

			continue
		}

		// TODO: cehck for expired and failed transactions

		// Re-broadcast the transaction
		_, err = t.provider.SendTx(trx)
		if err != nil {
			logger.Warn("failed to broadcast transaction", "error", err, "id", txID, "fee", trx.Fee())

			continue
		}
	}

	for _, trx := range blk.Transactions() {
		txID := trx.ID().String()

		if _, ok := pendingTxs[txID]; ok {
			pendingInfo := pendingTxs[txID]
			if err := t.storage.UpdateTransactionStatus(pendingInfo.No,
				types.TransactionStatusConfirmed, blk.Height()); err != nil {
				logger.Warn("failed to update transaction status", "error", err, "id", txID)
			}

			logger.Info("confirmed pending transaction", "id", trx.ID())

			continue
		}

		if t.addIncomingIfWallet(trx, collectReceivers(trx), block.Height(blk.Height())) {
			continue
		}

		if trx.IsSubsidyTx() {
			continue
		}

		signer := trx.Payload().Signer().String()
		if !t.storage.HasAddress(signer) {
			continue
		}

		if err := t.addTransactionWithStatus(trx, types.TransactionStatusConfirmed, block.Height(blk.Height())); err != nil {
			logger.Warn("failed to add outgoing transaction to wallet", "error", err, "id", trx.ID())

			continue
		}

		logger.Info("added outgoing transaction to wallet", "id", trx.ID())
	}
}
