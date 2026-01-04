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
			public_key 		TEXT NOT NULL,
			path 			TEXT NOT NULL DEFAULT '',
			label 			TEXT NOT NULL DEFAULT '',
			created_at 		DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at 		DATETIME DEFAULT CURRENT_TIMESTAMP
		)`

	createTransactionsTableSQL = `
		CREATE TABLE transactions (
			id 				TEXT NOT NULL,
			sender 			TEXT NOT NULL,
			receiver 		TEXT NOT NULL,
			amount 			INTEGER NOT NULL,
			fee 			INTEGER NOT NULL,
			memo 			TEXT NOT NULL DEFAULT '',
			status 			INTEGER NOT NULL,
			block_height 	INTEGER NOT NULL,
			payload_type 	INTEGER NOT NULL,
			data 			BLOB NOT NULL DEFAULT X'',
			comment 		TEXT NOT NULL DEFAULT '',
			created_at 		DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at 		DATETIME DEFAULT CURRENT_TIMESTAMP,

			PRIMARY KEY (id, receiver)
		)`

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
		INSERT INTO transactions (id, sender, receiver, amount, fee, memo, status, block_height, payload_type, data, comment)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	updateTransactionStatusSQL = `
		UPDATE transactions SET status = ?, block_height = ?
		WHERE id = ?`

	selectTransactionByIDSQL = `
		SELECT
			id, sender, receiver, amount, fee, memo, status, block_height, payload_type,
			data, comment, created_at, updated_at
		FROM transactions
		WHERE id = ?
		LIMIT 1`

	countTransactionByIDSQL = `
		SELECT COUNT(*) FROM transactions WHERE id = ?`

	selectPendingTransactionsSQL = `
		SELECT
			id, sender, receiver, amount, fee, memo, status, block_height, payload_type,
			data, comment, created_at, updated_at
		FROM transactions WHERE status = ?
		ORDER BY created_at DESC`
)
