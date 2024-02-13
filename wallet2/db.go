package wallet2

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	//nolint:all
	_ "github.com/glebarez/go-sqlite"
)

var (
	ErrCouldNotOpenDatabase = errors.New("could not open database")
	ErrCouldNotCreateTable  = errors.New("could not create table")
)

type db struct {
	*sql.DB
}

func newDB(path string) (*db, error) {
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
		" address VARCHAR, public_key VARCHAR, label VARCHAR, path VARCHAR)"
	_, err := d.ExecContext(context.Background(), addressQuery)
	if err != nil && !strings.Contains(err.Error(), "already exists") {
		return ErrCouldNotCreateTable
	}

	return nil
}

func (d *db) createTransactionTable() error {
	transactionQuery := "CREATE TABLE transactions (id INTEGER PRIMARY KEY AUTOINCREMENT," +
		" tx_id VARCHAR, block_height INTEGER, block_time INTEGER, payload_type VARCHAR," +
		" data VARCHAR, description VARCHAR, amount BIGINT,created_at TIMESTAMP)"
	_, err := d.ExecContext(context.Background(), transactionQuery)
	if err != nil && !strings.Contains(err.Error(), "already exists") {
		return ErrCouldNotCreateTable
	}

	return nil
}

func (d *db) createPairTable() error {
	pairQuery := "CREATE TABLE pairs (id VARCHAR PRIMARY KEY, value VARCHAR)"
	_, err := d.ExecContext(context.Background(), pairQuery)
	if err != nil && !strings.Contains(err.Error(), "already exists") {
		return ErrCouldNotCreateTable
	}

	return nil
}
