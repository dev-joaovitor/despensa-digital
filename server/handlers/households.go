package handlers

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/dev-joaovitor/despensa-digital/models"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (e *Env) HouseholdsHandler(r chi.Router) {
	r.Get("/", e.GetHouseholdHandler)
	r.Patch("/", e.UpdateHouseholdHandler)
	r.Post("/code", e.GenerateCodeHandler)
}

func (e *Env) UpdateHouseholdHandler(w http.ResponseWriter, r *http.Request) {
	var providedHousehold UpdateHouseholdDTO
	ReadJSON(w, r, &providedHousehold)

	err := UpdateHouseholdValidator(&providedHousehold)
	if err != nil {
		WriteError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	if providedHousehold.Name == "" {
		WriteJSON(w, http.StatusOK, nil, "Nada mudou")
		return
	}

	userId := e.GetSessionUserId(r.Context())
	sessionHousehold, err := e.GetSessionUserHousehold(r.Context())
	if err != nil {
		fmt.Printf("Error getting household: %v\n", err)
		WriteError(w, http.StatusInternalServerError, "Erro interno")
		return
	}

	if *sessionHousehold.CreatorID != userId {
		WriteError(w, http.StatusUnauthorized, "Apenas o criador da residência pode alterar seus dados")
		return
	}

	householdId := sessionHousehold.ID

	transaction, err := e.DB.Begin(r.Context())
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Erro interno no banco de dados")
		return
	}
	defer transaction.Rollback(r.Context())

	err = transaction.QueryRow(
		r.Context(),
		`
		UPDATE households
		SET name = $1, updated_at = NOW()
		WHERE id = $2
		RETURNING id
		`,
		&providedHousehold.Name,
		householdId,
	).Scan(nil)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			WriteError(w, http.StatusNotFound, "Residência não encontrada")
			return
		}
		fmt.Printf("Database error: %v\n", err)
		WriteError(w, http.StatusInternalServerError, "Erro interno no banco de dados")
		return
	}

	transaction.Commit(r.Context())
	WriteJSON(w, http.StatusOK, nil, "Residência atualizada com sucesso")
}

func (e *Env) GenerateCodeHandler(w http.ResponseWriter, r *http.Request) {
	userId := e.GetSessionUserId(r.Context())
	sessionHousehold, err := e.GetSessionUserHousehold(r.Context())
	if err != nil {
		fmt.Printf("Error getting household: %v\n", err)
		WriteError(w, http.StatusInternalServerError, "Erro interno")
		return
	}

	if *sessionHousehold.CreatorID != userId {
		WriteError(w, http.StatusUnauthorized, "Apenas o criador da residência pode alterar seus dados")
		return
	}

	householdId := sessionHousehold.ID

	invitationCode, err := uuid.NewRandom()
	if err != nil {
		WriteError(
			w,
			http.StatusInternalServerError,
			"Falha ao gerar código de convite para nova residência: " + err.Error(),
			)
		return
	}

	transaction, err := e.DB.Begin(r.Context())
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Erro interno no banco de dados")
		return
	}
	defer transaction.Rollback(r.Context())

	err = transaction.QueryRow(
		r.Context(),
		`
		UPDATE households
		SET invitation_code = $1, updated_at = NOW()
		WHERE id = $2
		RETURNING id
		`,
		invitationCode,
		householdId,
	).Scan(nil)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			WriteError(w, http.StatusNotFound, "Residência não encontrada")
			return
		}
		fmt.Printf("Database error: %v\n", err)
		WriteError(w, http.StatusInternalServerError, "Erro interno no banco de dados")
		return
	}

	transaction.Commit(r.Context())
	WriteJSON(w, http.StatusOK, map[string]string{"new_code": invitationCode.String()}, "Novo código gerado com sucesso")
}

func (e *Env) GetHouseholdHandler(w http.ResponseWriter, r *http.Request) {
	sessionHousehold, err := e.GetSessionUserHousehold(r.Context())
	if err != nil {
		fmt.Printf("Error getting household: %v\n", err)
		WriteError(w, http.StatusInternalServerError, "Erro interno")
		return
	}
	WriteJSON(w, http.StatusOK, sessionHousehold, "Residência encontrada com sucesso")
}

func (e *Env) GetSessionUserHousehold(ctx context.Context) (*models.Household, error) {
	userId := e.GetSessionUserId(ctx)

	var foundHousehold models.Household

	err := e.DB.QueryRow(
		ctx, 
		`
		SELECT h.id, h.name, h.creator_id,
			h.invitation_code, h.created_at, h.updated_at
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
		&foundHousehold.InvitationCode,
		&foundHousehold.CreatedAt,
		&foundHousehold.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	if *foundHousehold.CreatorID != userId {
		foundHousehold.InvitationCode = ""
	}

	return &foundHousehold, nil
}
