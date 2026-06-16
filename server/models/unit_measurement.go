package models

// UnitMeasurement maps to the 'unit_measurements' table.
type UnitMeasurement struct {
	ID        int64     `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	Acronym   string    `db:"acronym" json:"acronym"`
}
