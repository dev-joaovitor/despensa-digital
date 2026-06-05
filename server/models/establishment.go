package models

import "time"

// Establishment maps to the 'establishments' table.
type Establishment struct {
	ID          int64      `db:"id" json:"id"`
	HouseholdID *int64     `db:"household_id" json:"household_id,omitempty"`
	Name        string     `db:"name" json:"name"`
	CreatedAt   time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt   *time.Time `db:"deleted_at" json:"deleted_at,omitempty"`
}
