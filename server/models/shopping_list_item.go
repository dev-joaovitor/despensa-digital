package models

import "time"

// ShoppingListItem maps to the 'shopping_list_items' table.
type ShoppingListItem struct {
	ID          int64      `db:"id" json:"id"`
	HouseholdID int64      `db:"household_id" json:"household_id"`
	ProductID   int64      `db:"product_id" json:"product_id"`
	Quantity    int        `db:"quantity" json:"quantity"`
	IsChecked   bool       `db:"is_checked" json:"is_checked"`
	CreatedAt   time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt   *time.Time `db:"deleted_at" json:"deleted_at,omitempty"`
}
