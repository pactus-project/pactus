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
	// TxDirectionAny Include both incoming and outgoing transactions.
	TxDirectionAny TxDirection = 0
	// TxDirectionIncoming includes only incoming transactions where the wallet receives funds.
	TxDirectionIncoming = 1
	// TxDirectionOutgoing includes only outgoing transactions where the wallet sends funds.
	TxDirectionOutgoing = 2
)

// listTransactionsConfig contains options for listing transactions.
type listTransactionsConfig struct {
	direction TxDirection
	address   string
	count     int
	skip      int
}

var defaultListTransactionsConfig = listTransactionsConfig{
	direction: TxDirectionAny,
	address:   "*",
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

func (t *transactions) ListTransactions(opts ...ListTransactionsOption) []*types.TransactionInfo {
	if t.storage.IsLegacy() {
		return nil
	}

	cfg := defaultListTransactionsConfig
	for _, opt := range opts {
		opt(&cfg)
	}

	params := storage.QueryParams{
		Count: cfg.count,
		Skip:  cfg.skip,
	}

	switch cfg.direction {
	case TxDirectionAny:
		params.Sender = cfg.address
		params.Receiver = cfg.address

	case TxDirectionIncoming:
		params.Sender = "*"
		params.Receiver = cfg.address

	case TxDirectionOutgoing:
		params.Sender = cfg.address
		params.Receiver = "*"
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

	for _, trx := range blk.Transactions() {
		txID := trx.ID().String()

		if _, ok := pendingTxs[txID]; ok {
			if err := t.storage.UpdateTransactionStatus(txID, types.TransactionStatusConfirmed, blk.Height()); err != nil {
				logger.Warn("failed to update transaction status", "error", err, "id", txID)
			}

			logger.Info("confirmed pending transaction", "id", txID)

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
			logger.Warn("failed to add outgoing transaction to wallet", "error", err, "id", txID)

			continue
		}

		logger.Info("added outgoing transaction to wallet", "id", txID)
	}
}
