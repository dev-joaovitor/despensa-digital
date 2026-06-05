package models

import "time"

// StockBatch maps to the 'stock_batches' table.

type StockBatch struct {
	ID                int64      `db:"id" json:"id"`
	HouseholdID       int64      `db:"household_id" json:"household_id"`
	ProductID         int64      `db:"product_id" json:"product_id"`
	EstablishmentID   int64      `db:"establishment_id" json:"establishment_id"`
	InitialQuantity   int        `db:"initial_quantity" json:"initial_quantity"`
	RemainingQuantity int        `db:"remaining_quantity" json:"remaining_quantity"`
	ExpirationDate    *time.Time `db:"expiration_date" json:"expiration_date,omitempty"` // Maps to PG DATE type.
	UnitPrice         float64    `db:"unit_price" json:"unit_price"`
	CreatedAt         time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt         time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt         *time.Time `db:"deleted_at" json:"deleted_at,omitempty"`
}
