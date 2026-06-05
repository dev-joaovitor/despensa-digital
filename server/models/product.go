package models

import "time"

// Product maps to the 'products' table.
type Product struct {
	ID            int64      `db:"id" json:"id"`
	HouseholdID   int64      `db:"household_id" json:"household_id"`
	BrandID       int64      `db:"brand_id" json:"brand_id"`
	MeasurementID int64      `db:"measurement_id" json:"measurement_id"`
	CategoryID    int64      `db:"category_id" json:"category_id"`
	Name          string     `db:"name" json:"name"`
	UnitSize      int        `db:"unit_size" json:"unit_size"`
	CreatedAt     time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt     *time.Time `db:"deleted_at" json:"deleted_at,omitempty"`
}
