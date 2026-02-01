package sqlitestorage

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	_ "github.com/glebarez/go-sqlite" // sqlite driver
	"github.com/google/uuid"
	"github.com/pactus-project/pactus/genesis"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/tx/payload"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/wallet/storage"
	"github.com/pactus-project/pactus/wallet/types"
	"github.com/pactus-project/pactus/wallet/vault"
)

// Wallet metadata keys.
const (
	keyVersion    = "version"
	keyUUID       = "uuid"
	keyCreatedAt  = "created_at"
	keyNetwork    = "network"
	keyDefaultFee = "default_fee"
	keyVault      = "vault"
	keyLastUpdate = "last_update"
)

// Storage represents the SQLite-based wallet storage implementing IStorage interface.
type Storage struct {
	ctx  context.Context
	db   *sql.DB
	path string
	info *types.WalletInfo
	vlt  *vault.Vault

	addressMap map[string]types.AddressInfo
}

func dbPath(path string) string {
	return filepath.Join(path, "wallet.db")
}

func configurePragmas(ctx context.Context, db *sql.DB) error {
	pragmas := []string{
		"PRAGMA locking_mode=EXCLUSIVE;",
		"PRAGMA synchronous=NORMAL;",
	}

	for _, pragma := range pragmas {
		if _, err := db.ExecContext(ctx, pragma); err != nil {
			return fmt.Errorf("failed to set pragma %q: %w", pragma, err)
		}
	}

	return nil
}

func openDB(ctx context.Context, path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", dbPath(path))
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := configurePragmas(ctx, db); err != nil {
		_ = db.Close()

		return nil, err
	}

	return db, nil
}

// Create creates a new SQLite storage instance and initializes the schema.
func Create(ctx context.Context, path string, network genesis.ChainType, vlt *vault.Vault) (*Storage, error) {
	if err := util.Mkdir(path); err != nil {
		return nil, err
	}

	db, err := openDB(ctx, path)
	if err != nil {
		return nil, err
	}

	// Initialize database schema
	tables := []string{
		createWalletTableSQL,
		createAddressesTableSQL,
		createTransactionsTableSQL,

		createAddressesUpdatedAtTriggerSQL,
		createTransactionsUpdatedAtTriggerSQL,

		createPendingNoIdxSQL,
		createTxSenderNoIdxSQL,
		createTxReceiverNoIdxSQL,

		createAddressesLastUpdateTriggerSQL,
		createTransactionsLastUpdateTriggerSQL,
	}
	for _, query := range tables {
		if _, err := db.ExecContext(ctx, query); err != nil {
			_ = db.Close()

			return nil, fmt.Errorf("failed to create table: %w", err)
		}
	}

	// Marshal vault to JSON
	vaultJSON, err := json.Marshal(vlt)
	if err != nil {
		_ = db.Close()

		return nil, fmt.Errorf("failed to marshal vault: %w", err)
	}

	// Store wallet metadata
	// KeyValue represents a key-value pair for wallet entries.
	type KeyValue struct {
		Key   string
		Value string
	}

	entries := []KeyValue{
		{Key: keyVersion, Value: fmt.Sprintf("%d", VersionLatest)},
		{Key: keyUUID, Value: uuid.New().String()},
		{Key: keyCreatedAt, Value: fmt.Sprintf("%d", util.RoundNow(1).Unix())},
		{Key: keyNetwork, Value: fmt.Sprintf("%d", network)},
		{Key: keyDefaultFee, Value: fmt.Sprintf("%d", amount.Amount(10_000_000))},
		{Key: keyVault, Value: string(vaultJSON)},
		{Key: keyLastUpdate, Value: util.RoundNow(1).Format("2006-01-02 15:04:05")},
	}

	for _, entry := range entries {
		if _, err := db.ExecContext(ctx, insertWalletEntrySQL, entry.Key, entry.Value); err != nil {
			_ = db.Close()

			return nil, fmt.Errorf("failed to insert wallet entry %s: %w", entry.Key, err)
		}
	}

	return open(ctx, db, path)
}

// Open opens an existing SQLite storage instance without creating schema.
func Open(ctx context.Context, path string) (*Storage, error) {
	db, err := openDB(ctx, path)
	if err != nil {
		return nil, err
	}

	return open(ctx, db, path)
}

// open loads wallet info and returns a Storage instance.
func open(ctx context.Context, db *sql.DB, path string) (*Storage, error) {
	// Ensure triggers and last_update exist for older wallets.
	// We use Exec instead of ExecContext here for simplicity as these are idempotent or "OR IGNORE".
	_, _ = db.Exec(createAddressesLastUpdateTriggerSQL)
	_, _ = db.Exec(createTransactionsLastUpdateTriggerSQL)
	_, _ = db.Exec(insertWalletEntrySQL, keyLastUpdate, util.RoundNow(1).Format("2006-01-02 15:04:05"))

	strg := &Storage{
		ctx:  ctx,
		db:   db,
		path: path,
	}

	// Load wallet info into memory
	if err := strg.loadWalletInfo(); err != nil {
		_ = db.Close()

		return nil, fmt.Errorf("failed to load wallet info: %w", err)
	}

	// Load addresses into memory
	if err := strg.loadAddresses(); err != nil {
		_ = db.Close()

		return nil, fmt.Errorf("failed to load addresses: %w", err)
	}

	return strg, nil
}

// Close closes the database connection.
func (s *Storage) Close() error {
	if s.db != nil {
		return s.db.Close()
	}

	return nil
}

// WalletInfo returns the wallet information.
func (s *Storage) WalletInfo() *types.WalletInfo {
	return s.info
}

// loadWalletInfo loads wallet information from the database into memory.
func (s *Storage) loadWalletInfo() error {
	// Fetch all wallet entries at once
	rows, err := s.db.QueryContext(s.ctx, selectAllWalletEntriesSQL)
	if err != nil {
		return fmt.Errorf("failed to query wallet entries: %w", err)
	}
	defer func() { _ = rows.Close() }()

	entries := make(map[string]string)
	for rows.Next() {
		var name, value string
		if err := rows.Scan(&name, &value); err != nil {
			return fmt.Errorf("failed to scan wallet entry: %w", err)
		}
		entries[name] = value
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("failed to iterate wallet entries: %w", err)
	}

	version := 0
	if _, err := fmt.Sscanf(entries[keyVersion], "%d", &version); err != nil {
		return fmt.Errorf("failed to parse version: %w", err)
	}

	var defaultFee amount.Amount
	if _, err := fmt.Sscanf(entries[keyDefaultFee], "%d", &defaultFee); err != nil {
		return fmt.Errorf("failed to parse default fee: %w", err)
	}

	var createdAtUnix int64
	if _, err := fmt.Sscanf(entries[keyCreatedAt], "%d", &createdAtUnix); err != nil {
		return fmt.Errorf("failed to parse created_at: %w", err)
	}
	createdAt := time.Unix(createdAtUnix, 0)

	// Parse network type
	var network genesis.ChainType
	if _, err := fmt.Sscanf(entries[keyNetwork], "%d", &network); err != nil {
		return fmt.Errorf("failed to parse network: %w", err)
	}

	var vlt vault.Vault
	if vaultJSON, ok := entries[keyVault]; ok {
		if err := json.Unmarshal([]byte(vaultJSON), &vlt); err != nil {
			return fmt.Errorf("failed to unmarshal vault: %w", err)
		}
	}
	s.vlt = &vlt

	lastUpdate := time.Time{}
	if val, ok := entries[keyLastUpdate]; ok {
		t, err := time.Parse("2006-01-02 15:04:05", val)
		if err == nil {
			lastUpdate = t
		}
	}

	s.info = &types.WalletInfo{
		Path:       s.path,
		Driver:     "SQLite",
		Version:    version,
		Network:    network,
		DefaultFee: defaultFee,
		UUID:       entries[keyUUID],
		Encrypted:  s.vlt.IsEncrypted(),
		Neutered:   s.vlt.IsNeutered(),
		CreatedAt:  createdAt,
		LastUpdate: lastUpdate,
	}

	return nil
}

// loadAddresses loads addresses into memory.
func (s *Storage) loadAddresses() error {
	rows, err := s.db.QueryContext(s.ctx, selectAllAddressesSQL)
	if err != nil {
		return fmt.Errorf("failed to query addresses: %w", err)
	}
	defer func() { _ = rows.Close() }()

	s.addressMap = make(map[string]types.AddressInfo)
	for rows.Next() {
		addr, err := scanAddress(rows)
		if err != nil {
			return fmt.Errorf("failed to scan address: %w", err)
		}

		s.addressMap[addr.Address] = *addr
	}

	return rows.Err()
}

// Vault returns the vault.
func (s *Storage) Vault() *vault.Vault {
	return s.vlt
}

// UpdateVault updates the vault in storage.
func (s *Storage) UpdateVault(vlt *vault.Vault) error {
	if err := s.saveVault(vlt); err != nil {
		return err
	}
	s.vlt = vlt
	s.info.Encrypted = vlt.IsEncrypted()
	s.info.Neutered = vlt.IsNeutered()

	return nil
}

// SetDefaultFee sets the default fee.
func (s *Storage) SetDefaultFee(fee amount.Amount) error {
	if err := s.updateWalletEntry(keyDefaultFee, fmt.Sprintf("%d", fee)); err != nil {
		return err
	}
	s.info.DefaultFee = fee

	return nil
}

// AllAddresses returns all addresses in the wallet.
func (s *Storage) AllAddresses() []types.AddressInfo {
	addresses := make([]types.AddressInfo, 0, len(s.addressMap))
	for _, addr := range s.addressMap {
		addresses = append(addresses, addr)
	}

	return addresses
}

// AddressCount returns the number of addresses in the wallet.
func (s *Storage) AddressCount() int {
	return len(s.addressMap)
}

// AddressInfo returns the address information for the given address.
func (s *Storage) AddressInfo(address string) (*types.AddressInfo, error) {
	info, exists := s.addressMap[address]
	if !exists {
		return nil, storage.ErrNotFound
	}

	return &info, nil
}

// InsertAddress inserts a new address.
func (s *Storage) InsertAddress(info *types.AddressInfo) error {
	_, err := s.db.ExecContext(s.ctx, insertAddressSQL,
		info.Address, info.PublicKey, info.Label, info.Path)
	if err != nil {
		return err
	}

	return s.loadAddresses()
}

// UpdateAddress updates an existing address.
func (s *Storage) UpdateAddress(info *types.AddressInfo) error {
	_, err := s.db.ExecContext(s.ctx, updateAddressSQL,
		info.Label, info.PublicKey, info.Path, info.Address)
	if err != nil {
		return err
	}

	return s.loadAddresses()
}

// HasAddress checks if an address exists.
func (s *Storage) HasAddress(address string) bool {
	_, exists := s.addressMap[address]

	return exists
}

// InsertTransaction inserts a new transaction.
func (s *Storage) InsertTransaction(info *types.TransactionInfo) error {
	result, err := s.db.ExecContext(s.ctx, insertTransactionSQL,
		info.TxID,
		info.Sender,
		info.Receiver,
		info.Direction,
		info.Amount,
		info.Fee,
		info.Memo,
		info.Status,
		info.BlockHeight,
		int(info.PayloadType),
		info.Data,
		info.Comment,
	)
	if err != nil {
		return err
	}

	info.No, _ = result.LastInsertId()

	return nil
}

// UpdateTransactionStatus updates the status and block height for all transactions with the given primary key.
func (s *Storage) UpdateTransactionStatus(no int64, status types.TransactionStatus, blockHeight uint32) error {
	_, err := s.db.ExecContext(s.ctx, updateTransactionStatusSQL, int(status), blockHeight, no)

	return err
}

// HasTransaction checks if a transaction exists by transaction ID.
func (s *Storage) HasTransaction(txID string) bool {
	var count int
	err := s.db.QueryRowContext(s.ctx, countTransactionByTxIDSQL, txID).Scan(&count)
	if err != nil {
		return false
	}

	return count > 0
}

// scanner is an interface that both *sql.Row and *sql.Rows satisfy.
type scanner interface {
	Scan(dest ...any) error
}

// scanAddress scans a row into an AddressInfo struct.
func scanAddress(s scanner) (*types.AddressInfo, error) {
	var info types.AddressInfo

	err := s.Scan(
		&info.Address,
		&info.PublicKey,
		&info.Label,
		&info.Path,
		&info.CreatedAt,
		&info.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &info, nil
}

// scanTransaction scans a row into a TransactionInfo struct.
func scanTransaction(s scanner) (*types.TransactionInfo, error) {
	var info types.TransactionInfo
	var status, payloadType int

	err := s.Scan(
		&info.No,
		&info.TxID,
		&info.Sender,
		&info.Receiver,
		&info.Direction,
		&info.Amount,
		&info.Fee,
		&info.Memo,
		&status,
		&info.BlockHeight,
		&payloadType,
		&info.Data,
		&info.Comment,
		&info.CreatedAt,
		&info.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	info.Status = types.TransactionStatus(status)
	info.PayloadType = payload.Type(payloadType)

	return &info, nil
}

// GetTransaction retrieves a transaction by primary key.
func (s *Storage) GetTransaction(no int64) (*types.TransactionInfo, error) {
	info, err := scanTransaction(s.db.QueryRowContext(s.ctx, selectTransactionByNoSQL, no))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storage.ErrNotFound
		}

		return nil, err
	}

	return info, nil
}

// QueryTransactions returns transactions matching the provided filters with pagination.
// Empty or "*" address value is treated as no filter. Filters are combined with AND.
func (s *Storage) QueryTransactions(params storage.QueryParams) ([]*types.TransactionInfo, error) {
	conditions := make([]string, 0, 3)
	args := make([]any, 0, 6)

	if params.Direction != types.TxDirectionAny {
		conditions = append(conditions, "direction = ?")
		args = append(args, params.Direction)
	}

	if params.Address != "*" {
		switch params.Direction {
		case types.TxDirectionAny:
			conditions = append(conditions, "(sender = ? OR receiver = ?)")
			args = append(args, params.Address, params.Address)
		case types.TxDirectionIncoming:
			conditions = append(conditions, "receiver = ?")
			args = append(args, params.Address)
		case types.TxDirectionOutgoing:
			conditions = append(conditions, "sender = ?")
			args = append(args, params.Address)
		}
	}

	where := ""
	if len(conditions) > 0 {
		where = "WHERE " + strings.Join(conditions, " AND ")
	}

	// Apply pagination
	args = append(args, params.Count, params.Skip)

	var queryBuilder strings.Builder
	queryBuilder.WriteString(`
		SELECT
			no, tx_id, sender, receiver, direction, amount, fee, memo, status, block_height, payload_type,
			data, comment, created_at, updated_at
		FROM transactions
	`)
	if where != "" {
		queryBuilder.WriteString(where)
		queryBuilder.WriteString("\n")
	}
	queryBuilder.WriteString(`ORDER BY created_at DESC
		LIMIT ? OFFSET ?`)

	query := queryBuilder.String()

	rows, err := s.db.QueryContext(s.ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	transactions := make([]*types.TransactionInfo, 0)
	for rows.Next() {
		info, err := scanTransaction(rows)
		if err != nil {
			return nil, err
		}

		transactions = append(transactions, info)
	}

	return transactions, rows.Err()
}

// GetPendingTransactions returns pending transactions keyed by transaction ID.
func (s *Storage) GetPendingTransactions() (map[string]*types.TransactionInfo, error) {
	rows, err := s.db.QueryContext(s.ctx, selectPendingTransactionsSQL, int(types.TransactionStatusPending))
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	pending := make(map[string]*types.TransactionInfo)
	for rows.Next() {
		info, err := scanTransaction(rows)
		if err != nil {
			return nil, err
		}

		pending[info.TxID] = info
	}

	return pending, rows.Err()
}

// Clone creates a copy of the storage at a new path.
func (s *Storage) Clone(path string) (storage.IStorage, error) {
	if err := util.Mkdir(path); err != nil {
		return nil, err
	}

	// Use VACUUM INTO to create a backup
	_, err := s.db.ExecContext(s.ctx, fmt.Sprintf("VACUUM INTO '%s'", dbPath(path)))
	if err != nil {
		return nil, fmt.Errorf("failed to backup database: %w", err)
	}

	// Open the cloned database
	clonedStorage, err := Open(s.ctx, path)
	if err != nil {
		return nil, err
	}

	// Update UUID and CreatedAt for the cloned wallet
	newUUID := uuid.New().String()
	if err := clonedStorage.updateWalletEntry(keyUUID, newUUID); err != nil {
		_ = clonedStorage.Close()

		return nil, err
	}
	clonedStorage.info.UUID = newUUID

	newCreatedAt := util.RoundNow(1)
	if err := clonedStorage.updateWalletEntry(keyCreatedAt, fmt.Sprintf("%d", newCreatedAt.Unix())); err != nil {
		_ = clonedStorage.Close()

		return nil, err
	}
	clonedStorage.info.CreatedAt = newCreatedAt
	clonedStorage.info.Path = path

	return clonedStorage, nil
}

func (s *Storage) updateWalletEntry(key, value string) error {
	_, err := s.db.ExecContext(s.ctx, updateWalletEntrySQL, value, key)
	if err != nil {
		return err
	}

	_, err = s.db.ExecContext(s.ctx, incrementLastUpdateSQL)

	return err
}

func (s *Storage) saveVault(vlt *vault.Vault) error {
	vaultJSON, err := json.Marshal(vlt)
	if err != nil {
		return fmt.Errorf("failed to marshal vault: %w", err)
	}

	return s.updateWalletEntry(keyVault, string(vaultJSON))
}

func (*Storage) IsLegacy() bool {
	return false
}

func (s *Storage) LastUpdate() time.Time {
	var val string
	err := s.db.QueryRowContext(s.ctx, "SELECT value FROM wallet WHERE name = ?", keyLastUpdate).Scan(&val)
	if err != nil {
		return time.Time{}
	}

	t, err := time.Parse("2006-01-02 15:04:05", val)
	if err != nil {
		return time.Time{}
	}

	return t
}
