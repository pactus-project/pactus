package sqlitestorage

// SQL queries and schema definitions.

const (
	// Schema creation queries.
	createWalletTableSQL = `
		CREATE TABLE wallet (
			name 			TEXT PRIMARY KEY,
			value 			TEXT NOT NULL
		)`

	createAddressesTableSQL = `
		CREATE TABLE addresses (
			address 		TEXT PRIMARY KEY,
			public_key 		TEXT NOT NULL DEFAULT '',
			path 			TEXT NOT NULL DEFAULT '',
			label 			TEXT NOT NULL DEFAULT '',
			created_at 		DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at 		DATETIME DEFAULT CURRENT_TIMESTAMP
		)`

	createTransactionsTableSQL = `
		CREATE TABLE transactions (
			no 				INTEGER PRIMARY KEY,   -- alias for rowid, auto-incremented starting at 1
			tx_id 			TEXT NOT NULL,         -- transaction identifier
			sender 			TEXT NOT NULL,
			receiver 		TEXT NOT NULL,
			direction 		INTEGER NOT NULL,
			amount 			INTEGER NOT NULL,
			fee 			INTEGER NOT NULL,
			memo 			TEXT NOT NULL DEFAULT '',
			status 			INTEGER NOT NULL,      -- status: failed=-1, pending=0, confirmed=1
			block_height 	INTEGER NOT NULL,
			payload_type 	INTEGER NOT NULL,
			data 			BLOB NOT NULL DEFAULT X'',
			comment 		TEXT NOT NULL DEFAULT '',
			created_at 		DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at 		DATETIME DEFAULT CURRENT_TIMESTAMP,

			UNIQUE (tx_id, receiver)
			CHECK (direction IN (1, 2)) -- incoming=1, outgoing=2
		)`

	createAddressesUpdatedAtTriggerSQL = `
		CREATE TRIGGER IF NOT EXISTS trg_addresses_updated_at
		AFTER UPDATE ON addresses
		FOR EACH ROW
		BEGIN
			UPDATE addresses
			SET updated_at = CURRENT_TIMESTAMP
			WHERE address = OLD.address;
		END`

	createTransactionsUpdatedAtTriggerSQL = `
		CREATE TRIGGER IF NOT EXISTS trg_transactions_updated_at
		AFTER UPDATE ON transactions
		FOR EACH ROW
		BEGIN
			UPDATE transactions
			SET updated_at = CURRENT_TIMESTAMP
			WHERE no = OLD.no;
		END`

	// Partial index to speed up pending-transaction queries ordered by creation time.
	createPendingNoIdxSQL = `
		CREATE INDEX IF NOT EXISTS idx_transactions_pending_no
		ON transactions (no DESC)
		WHERE status = 0`

	// Indexes to speed up lookups by sender/receiver and ordering by created_at.
	createTxSenderNoIdxSQL = `
		CREATE INDEX IF NOT EXISTS idx_tx_sender_no
		ON transactions(sender, no DESC)`

	createTxReceiverNoIdxSQL = `
		CREATE INDEX IF NOT EXISTS idx_tx_receiver_no
		ON transactions(receiver, no DESC)`

	// Wallet table operations.
	insertWalletEntrySQL = `
		INSERT INTO wallet (name, value) VALUES (?, ?)`

	updateWalletEntrySQL = `
		UPDATE wallet SET value = ? WHERE name = ?`

	selectAllWalletEntriesSQL = `
		SELECT name, value FROM wallet`

	// Address table operations.
	insertAddressSQL = `
		INSERT INTO addresses (address, public_key, label, path)
		VALUES (?, ?, ?, ?)`

	updateAddressSQL = `
		UPDATE addresses SET label = ?, public_key = ?, path = ?
		WHERE address = ?`

	selectAllAddressesSQL = `
		SELECT address, public_key, label, path, created_at, updated_at
		FROM addresses ORDER BY created_at ASC`

	// Transaction table operations.
	insertTransactionSQL = `
		INSERT INTO transactions (tx_id, sender, receiver, direction, amount, fee, memo,
			status, block_height, payload_type, data, comment)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	updateTransactionStatusSQL = `
		UPDATE transactions SET status = ?, block_height = ?
		WHERE no = ?`

	selectTransactionByNoSQL = `
		SELECT
			no, tx_id, sender, receiver, direction, amount, fee, memo, status, block_height, payload_type,
			data, comment, created_at, updated_at
		FROM transactions
		WHERE no = ?
		LIMIT 1`

	countTransactionByTxIDSQL = `
		SELECT COUNT(*) FROM transactions WHERE tx_id = ?`

	selectPendingTransactionsSQL = `
		SELECT
		no, tx_id, sender, receiver, direction, amount, fee, memo, status, block_height, payload_type,
		data, comment, created_at, updated_at
		FROM transactions WHERE status = ?
		ORDER BY created_at DESC`
)
