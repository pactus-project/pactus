package db

const (
	EmptyQuery = ""
)

type QueryOption string

func WithTransactionStatus() string {
	return "WHERE status = ?"
}

func WithTransactionAddr() string {
	return "WHERE address = ?"
}

func WithTransactionStatusAndAddr() string {
	return "WHERE status = ? AND address = ?"
}
