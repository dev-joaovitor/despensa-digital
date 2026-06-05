package models

import "time"

// Household maps to the 'households' table.
type Household struct {
	ID             int64      `db:"id" json:"id"`
	Name           string     `db:"name" json:"name"`
	CreatorID      *int64     `db:"creator_id" json:"creator_id"`
	InvitationCode string     `db:"invitation_code" json:"invitation_code"`
	CreatedAt      time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt      *time.Time `db:"deleted_at" json:"deleted_at,omitempty"`
}
