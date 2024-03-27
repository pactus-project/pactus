package db

const (
	EmptyQuery = ""
)

func WithTransactionStatus() string {
	return "WHERE status = ?"
}

func WithTransactionAddr() string {
	return "WHERE address = ?"
}

func WithTransactionStatusAndAddr() string {
	return "WHERE status = ? AND address = ?"
}
