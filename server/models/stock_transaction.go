package models

import "time"

// StockTransactionType defines the enum type for transaction types.
type StockTransactionType string
const (
	TransactionPurchase    StockTransactionType = "purchase"
	TransactionConsumption StockTransactionType = "consumption"
	TransactionWaste       StockTransactionType = "waste"
	TransactionCorrection  StockTransactionType = "correction"
)

// StockTransaction maps to the 'stock_transactions' table.
type StockTransaction struct {
	ID        int64                `db:"id" json:"id"`
	BatchID   int64                `db:"batch_id" json:"batch_id"`
	Quantity  int                  `db:"quantity" json:"quantity"`
	Type      StockTransactionType `db:"type" json:"type"`
	CreatedAt time.Time            `db:"created_at" json:"created_at"`
	UpdatedAt time.Time            `db:"updated_at" json:"updated_at"`
}
