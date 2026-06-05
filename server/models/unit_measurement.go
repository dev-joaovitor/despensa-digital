package models

import "time"

// UnitMeasurement maps to the 'unit_measurements' table.
type UnitMeasurement struct {
	ID        int64     `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	Acronym   string    `db:"acronym" json:"acronym"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
