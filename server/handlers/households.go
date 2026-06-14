package handlers

import (
	"context"

	"github.com/dev-joaovitor/despensa-digital/models"
)

func (e *Env) GetSessionUserHousehold(ctx context.Context) (*models.Household, error) {
	userId := e.GetSessionUserId(ctx)

	var foundHousehold models.Household

	err := e.DB.QueryRow(
		ctx, 
		`
		SELECT h.id, h.name, h.creator_id
		FROM households h
		LEFT JOIN users u
		ON u.household_id = h.id
		AND u.id = $1
		WHERE h.deleted_at IS NULL
		`,
		userId,
	).Scan(
		&foundHousehold.ID,
		&foundHousehold.Name,
		&foundHousehold.CreatorID,
	)

	if err != nil {
		return nil, err
	}

	return &foundHousehold, nil
}
