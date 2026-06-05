package models

import "time"

// User maps to the 'users' table.
type User struct {
	ID               int64      `db:"id" json:"id"`
	HouseholdID      *int64      `db:"household_id" json:"household_id,omitempty"`
	Household		 *Household  `json:"household,omitempty"`
	FullName         string     `db:"full_name" json:"full_name"`
	Password         string     `db:"password" json:"-"` // Omitted from JSON serialization for security.
	Email            string     `db:"email" json:"email"`
	VerificationCode *string    `db:"verification_code" json:"verification_code,omitempty"`
	ExpiresAt        *time.Time `db:"expires_at" json:"expires_at,omitempty"`
	CreatedAt        time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt        time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt        *time.Time `db:"deleted_at" json:"deleted_at,omitempty"`
}
