package models

import "time"

// Category maps to the 'categories' table.
type Category struct {
	ID          int64      `db:"id" json:"id"`
	HouseholdID *int64     `db:"household_id" json:"household_id,omitempty"` // Nullable if category is global.
	Name        string     `db:"name" json:"name"`
	CreatedAt   time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt   *time.Time `db:"deleted_at" json:"deleted_at,omitempty"`
}
