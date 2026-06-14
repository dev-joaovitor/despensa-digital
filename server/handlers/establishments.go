package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/dev-joaovitor/despensa-digital/models"
	"github.com/go-chi/chi/v5"
)

func (e *Env) EstablishmentsHandler(r chi.Router) {
	r.Get("/{id}", e.ListOneEstablishmentHandler)
	r.Get("/", e.ListEstablishmentsHandler)
	r.Post("/", e.CreateEstablishmentHandler)
	r.Patch("/{id}", e.UpdateEstablishmentHandler)
	r.Delete("/{id}", e.DeleteEstablishmentHandler)
}

func (e *Env) CreateEstablishmentHandler(w http.ResponseWriter, r *http.Request) {
	var err error

	var providedEstablishment CreateEstablishmentDTO
	ReadJSON(w, r, &providedEstablishment)
	err = CreateEstablishmentValidator(&providedEstablishment)

	if err != nil {
		WriteError(w, http.StatusUnprocessableEntity, "Erro de validação: " + err.Error())
		return
	}

	transaction, err := e.DB.Begin(r.Context())

	if err != nil {
		if transaction != nil {
			transaction.Rollback(r.Context())
		}

		WriteError(w, http.StatusInternalServerError, "Erro interno no banco de dados")
		return
	}

	var createdEstablishment models.Establishment
	userHousehold, err := e.GetSessionUserHousehold(r.Context())
	if err != nil {
		transaction.Rollback(r.Context())
		WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = e.DB.QueryRow(
		r.Context(),
		`
		INSERT INTO establishments (household_id, name)
		VALUES ($1, $2)
		RETURNING id, household_id, name, created_at, updated_at
		`,
		&userHousehold.ID,
		&providedEstablishment.Name,
	).Scan(
		&createdEstablishment.ID,
		&createdEstablishment.HouseholdID,
		&createdEstablishment.Name,
		&createdEstablishment.CreatedAt,
		&createdEstablishment.UpdatedAt,
	)

	if err != nil {
		transaction.Rollback(r.Context())
		WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	transaction.Commit(r.Context())
	WriteJSON(w, http.StatusCreated, &createdEstablishment, "Estabelecimento criado com sucesso")
}

func (e *Env) UpdateEstablishmentHandler(w http.ResponseWriter, r *http.Request) {
	var err error

	var providedEstablishment UpdateEstablishmentDTO
	ReadJSON(w, r, &providedEstablishment)
	err = UpdateEstablishmentValidator(&providedEstablishment)
	if err != nil {
		WriteError(w, http.StatusUnprocessableEntity, "Erro de validação: " + err.Error())
		return
	}

	if providedEstablishment.Name == "" {
		WriteJSON(w, http.StatusNotModified, nil, "Nada mudou")
		return
	}

	establishmentId := r.PathValue("id")
	if establishmentId == "" {
		WriteError(w, http.StatusBadRequest, "Insira um ID")
		return
	}

	transaction, err := e.DB.Begin(r.Context())

	if err != nil {
		if transaction != nil {
			transaction.Rollback(r.Context())
		}

		WriteError(w, http.StatusInternalServerError, "Erro interno no banco de dados")
		return
	}

	_, err = e.DB.Query(
		r.Context(),
		`
		UPDATE establishments 
		SET name = $1, updated_at = NOW()
		WHERE id = $2
		`,
		&providedEstablishment.Name,
		establishmentId,
	)
	if err != nil {
		transaction.Rollback(r.Context())
		WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	transaction.Commit(r.Context())
	WriteJSON(w, http.StatusOK, nil, "Estabelecimento atualizado com sucesso")
}

func (e *Env) ListEstablishmentsHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := e.DB.Query(
		r.Context(),
		`
		SELECT id, name, household_id, created_at, updated_at
		FROM establishments
		WHERE deleted_at IS NULL
		ORDER BY created_at, updated_at DESC
		`,
	)
	if err != nil {
		fmt.Printf("Database error: %v\n", err)
		WriteError(w, http.StatusInternalServerError, "Erro interno no banco de dados")
		return
	}

	establishments := []models.Establishment{}

	for rows.Next() == true {
		establishment := models.Establishment{}
		rows.Scan(
			&establishment.ID,
			&establishment.Name,
			&establishment.HouseholdID,
			&establishment.CreatedAt,
			&establishment.UpdatedAt,
		)
		establishments = append(establishments, establishment)
	}

	WriteJSON(w, http.StatusOK, establishments, "Estabelecimentos listados com sucesso")
}

func (e *Env) ListOneEstablishmentHandler(w http.ResponseWriter, r *http.Request) {
	var foundEstablishment models.Establishment
	
	establishmentId := r.PathValue("id")
	if establishmentId == "" {
		WriteError(w, http.StatusBadRequest, "Insira um ID")
		return
	}

	err := e.DB.QueryRow(
		r.Context(),
		`
		SELECT id, name, household_id, created_at, updated_at
		FROM establishments
		WHERE id = $1
		AND deleted_at IS NULL
		`,
		establishmentId,
	).Scan(
		&foundEstablishment.ID,
		&foundEstablishment.Name,
		&foundEstablishment.HouseholdID,
		&foundEstablishment.CreatedAt,
		&foundEstablishment.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			WriteError(w, http.StatusNotFound, "Estabelecimento não encontrado")
			return
		}

		fmt.Printf("Database error: %v\n", err)
		WriteError(w, http.StatusInternalServerError, "Erro interno no banco de dados")
		return
	}

	WriteJSON(w, http.StatusOK, foundEstablishment, "Estabelecimento listado com sucesso")
}

func (e *Env) DeleteEstablishmentHandler(w http.ResponseWriter, r *http.Request) {
	establishmentId := r.PathValue("id")
	if establishmentId == "" {
		WriteError(w, http.StatusBadRequest, "Insira um ID")
		return
	}

	err := e.DB.QueryRow(
		r.Context(),
		`
		UPDATE establishments
		SET deleted_at = NOW()
		WHERE id = $1
		AND deleted_at IS NULL
		RETURNING id
		`,
		establishmentId,
	).Scan(nil)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			WriteError(w, http.StatusNotFound, "Estabelecimento não encontrado")
			return
		}

		fmt.Printf("Database error: %v\n", err)
		WriteError(w, http.StatusInternalServerError, "Erro interno no banco de dados")
		return
	}

	WriteJSON(w, http.StatusOK, nil, "Estabelecimento deletado com sucesso")
}

