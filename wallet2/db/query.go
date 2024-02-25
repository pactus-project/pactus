package db

import "fmt"

type QueryOption string

func WithTransactionStatus(status Status) QueryOption {
	return QueryOption(fmt.Sprintf("WHERE status = %d", status))
}

func WithTransactionAddr(addr string) QueryOption {
	return QueryOption(fmt.Sprintf("WHERE address = '%s'", addr))
}

func WithTransactionStatusAndAddr(status Status, addr string) QueryOption {
	return QueryOption(fmt.Sprintf("WHERE status = %d AND address = '%s'", status, addr))
}
