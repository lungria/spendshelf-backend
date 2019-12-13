package db

// Repository define i/o methods for MongoDB
type Repository interface {
	GetTransactionByID(transactionID string) (Transaction, error)
	GetAllTransactions(accountID string) ([]Transaction, error)
	SaveOneTransaction(transaction *Transaction) error
}
