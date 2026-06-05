package models

import "time"

// PriceObservation maps to the 'price_observations' table.
type PriceObservation struct {
	ID              int64      `db:"id" json:"id"`
	HouseholdID     int64      `db:"household_id" json:"household_id"`
	ProductID       int64      `db:"product_id" json:"product_id"`
	EstablishmentID int64      `db:"establishment_id" json:"establishment_id"`
	ObservedPrice   float64    `db:"observed_price" json:"observed_price"`
	ObservedAt      time.Time  `db:"observed_at" json:"observed_at"`
	DeletedAt       *time.Time `db:"deleted_at" json:"deleted_at,omitempty"`
}
