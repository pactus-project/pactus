package wallet2

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	//nolint:all
	_ "github.com/glebarez/go-sqlite"
)

var (
	ErrCouldNotOpenDatabase    = errors.New("could not open database")
	ErrCouldNotCreateTable     = errors.New("could not create table")
	ErrCouldNotInsertIntoTable = errors.New("could not insert record into table")
	ErrCouldNotFindRecord      = errors.New("could not find record")
)

type DB interface {
	CreateTables() error

	InsertIntoAddress(addr *Address) (*Address, error)
	InsertIntoTransaction(t *Transaction) (*Transaction, error)
	InsertIntoPair(key string, value string) (*Pair, error)

	GetAddressByID(id int) (*Address, error)
	GetTransactionByID(id int) (*Transaction, error)
	GetPairByKey(key string) (*Pair, error)

	GetAllAddresses(pageIndex, pageSize int) ([]*Address, error)
	GetAllTransactions(pageIndex, pageSize int) ([]*Transaction, error)
}

type db struct {
	*sql.DB
}

type Address struct {
	ID        int       `json:"id"`         // id of wallet
	Address   string    `json:"address"`    // Address in the wallet
	PublicKey string    `json:"public_key"` // Public key associated with the address
	Label     string    `json:"label"`      // Label for the address
	Path      string    `json:"path"`       // Path for the address
	CreatedAt time.Time `json:"created_at"`
}

type Transaction struct {
	ID          int       `json:"id"`
	TxID        string    `json:"tx_id"`
	BlockHeight uint32    `json:"block_height"`
	BlockTime   uint32    `json:"block_time"`
	PayloadType string    `json:"payload_type"`
	Data        string    `json:"data"`
	Description string    `json:"description"`
	Amount      int64     `json:"amount"`
	Status      int       `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}

type Pair struct {
	Key       string
	Value     string
	CreatedAt time.Time
}

func newDB(path string) (DB, error) {
	dbInstance, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, ErrCouldNotOpenDatabase
	}

	return &db{
		dbInstance,
	}, nil
}

func (d *db) CreateTables() error {
	if err := d.createAddressTable(); err != nil {
		return err
	}

	if err := d.createTransactionTable(); err != nil {
		return err
	}

	return d.createPairTable()
}

func (d *db) createAddressTable() error {
	addressQuery := "CREATE TABLE addresses (id INTEGER PRIMARY KEY AUTOINCREMENT," +
		" address VARCHAR, public_key VARCHAR, label VARCHAR, path VARCHAR, created_at TIMESTAMP)"
	_, err := d.ExecContext(context.Background(), addressQuery)
	if err != nil && !strings.Contains(err.Error(), "already exists") {
		return ErrCouldNotCreateTable
	}

	return nil
}

func (d *db) createTransactionTable() error {
	transactionQuery := "CREATE TABLE transactions (id INTEGER PRIMARY KEY AUTOINCREMENT," +
		" tx_id VARCHAR, block_height INTEGER, block_time INTEGER, payload_type VARCHAR," +
		" data VARCHAR, description VARCHAR, amount BIGINT,status INTEGER, created_at TIMESTAMP)"
	_, err := d.ExecContext(context.Background(), transactionQuery)
	if err != nil && !strings.Contains(err.Error(), "already exists") {
		return ErrCouldNotCreateTable
	}

	return nil
}

func (d *db) createPairTable() error {
	pairQuery := "CREATE TABLE pairs (key VARCHAR PRIMARY KEY, value VARCHAR, created_at TIMESTAMP)"
	_, err := d.ExecContext(context.Background(), pairQuery)
	if err != nil && !strings.Contains(err.Error(), "already exists") {
		return ErrCouldNotCreateTable
	}

	return nil
}

func (d *db) InsertIntoAddress(addr *Address) (*Address, error) {
	insertQuery := "INSERT INTO addresses (address, public_key, label, path, created_at) VALUES (?,?,?,?,?)"

	addr.CreatedAt = time.Now().UTC()
	r, err := d.ExecContext(context.Background(), insertQuery, addr.Address,
		addr.PublicKey, addr.Label, addr.Path, addr.CreatedAt)
	if err != nil {
		return nil, ErrCouldNotInsertIntoTable
	}

	rowID, err := r.LastInsertId()
	if err != nil {
		return nil, ErrCouldNotInsertIntoTable
	}

	return &Address{
		ID:        int(rowID),
		Address:   addr.Address,
		PublicKey: addr.PublicKey,
		Label:     addr.Label,
		Path:      addr.Path,
		CreatedAt: addr.CreatedAt,
	}, nil
}

func (d *db) InsertIntoTransaction(t *Transaction) (*Transaction, error) {
	insertQuery := "INSERT INTO transactions (tx_id, block_height, block_time," +
		" payload_type, data, description, amount, status, created_at) VALUES" +
		" (?,?,?,?,?,?,?,?,?)"

	t.CreatedAt = time.Now().UTC()
	r, err := d.ExecContext(context.Background(), insertQuery, t.TxID, t.BlockHeight, t.BlockTime,
		t.PayloadType, t.Data, t.Description, t.Amount, t.Status, t.CreatedAt)
	if err != nil {
		return nil, ErrCouldNotInsertIntoTable
	}

	rowID, err := r.LastInsertId()
	if err != nil {
		return nil, ErrCouldNotInsertIntoTable
	}

	return &Transaction{
		ID:          int(rowID),
		TxID:        t.TxID,
		BlockHeight: t.BlockHeight,
		BlockTime:   t.BlockTime,
		PayloadType: t.PayloadType,
		Data:        t.Data,
		Description: t.Description,
		Amount:      t.Amount,
		Status:      t.Status,
		CreatedAt:   t.CreatedAt,
	}, nil
}

func (d *db) InsertIntoPair(key, value string) (*Pair, error) {
	createdAt := time.Now().UTC()
	insertQuery := "INSERT INTO pairs (key, value, created_at) VALUES (?,?,?)"
	_, err := d.ExecContext(context.Background(), insertQuery, key, value, createdAt)
	if err != nil {
		return nil, ErrCouldNotInsertIntoTable
	}

	return &Pair{
		Key:       key,
		Value:     value,
		CreatedAt: createdAt,
	}, nil
}

func (d *db) GetAddressByID(id int) (*Address, error) {
	getQuery := "SELECT * FROM addresses WHERE id = ?"
	row := d.QueryRowContext(context.Background(), getQuery, id)
	if row.Err() != nil {
		return nil, ErrCouldNotFindRecord
	}

	addr := &Address{}
	err := row.Scan(&addr.ID, &addr.Address, &addr.PublicKey, &addr.Label, &addr.Path, &addr.CreatedAt)
	if err != nil {
		return nil, ErrCouldNotFindRecord
	}

	return addr, nil
}

func (d *db) GetTransactionByID(id int) (*Transaction, error) {
	getQuery := "SELECT * FROM transactions WHERE id = ?"
	row := d.QueryRowContext(context.Background(), getQuery, id)
	if row.Err() != nil {
		return nil, ErrCouldNotFindRecord
	}

	t := &Transaction{}
	err := row.Scan(&t.ID, &t.TxID, &t.BlockHeight, &t.BlockTime, &t.PayloadType,
		&t.Data, &t.Description, &t.Amount, &t.Status, &t.CreatedAt)
	if err != nil {
		return nil, ErrCouldNotFindRecord
	}

	return t, nil
}

func (d *db) GetPairByKey(key string) (*Pair, error) {
	getQuery := "SELECT * FROM pairs WHERE key = ?"
	row := d.QueryRowContext(context.Background(), getQuery, key)
	if row.Err() != nil {
		return nil, ErrCouldNotFindRecord
	}

	p := &Pair{}
	err := row.Scan(&p.Key, &p.Value, &p.CreatedAt)
	if err != nil {
		return nil, ErrCouldNotFindRecord
	}

	return p, nil
}

func (d *db) GetAllAddresses(pageIndex, pageSize int) ([]*Address, error) {
	getAllQuery := "SELECT * FROM addresses ORDER BY id DESC LIMIT ? OFFSET ?"
	rows, err := d.QueryContext(context.Background(), getAllQuery, pageSize, calcOffset(pageIndex, pageSize))
	if err != nil || rows.Err() != nil {
		return nil, ErrCouldNotFindRecord
	}
	defer rows.Close()

	addrs := make([]*Address, 0, pageSize)
	for {
		if !rows.Next() {
			break
		}

		addr := &Address{}
		err := rows.Scan(&addr.ID, &addr.Address, &addr.PublicKey, &addr.Label, &addr.Path, &addr.CreatedAt)
		if err != nil {
			return nil, ErrCouldNotFindRecord
		}

		addrs = append(addrs, addr)
	}

	return addrs, nil
}

func (d *db) GetAllTransactions(pageIndex, pageSize int) ([]*Transaction, error) {
	getAllQuery := "SELECT * FROM transactions ORDER BY id DESC LIMIT ? OFFSET ?"
	rows, err := d.QueryContext(context.Background(), getAllQuery, pageSize, calcOffset(pageIndex, pageSize))
	if err != nil || rows.Err() != nil {
		return nil, ErrCouldNotFindRecord
	}
	defer rows.Close()

	transactions := make([]*Transaction, 0, pageSize)
	for {
		if !rows.Next() {
			break
		}

		t := &Transaction{}
		err := rows.Scan(&t.ID, &t.TxID, &t.BlockHeight, &t.BlockTime, &t.PayloadType,
			&t.Data, &t.Description, &t.Amount, &t.Status, &t.CreatedAt)
		if err != nil {
			return nil, ErrCouldNotFindRecord
		}

		transactions = append(transactions, t)
	}

	return transactions, nil
}

func calcOffset(pageIndex, pageSize int) int {
	return (pageIndex - 1) * pageSize
}
