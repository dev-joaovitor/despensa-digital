package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/dev-joaovitor/despensa-digital/models"
	"github.com/go-chi/chi/v5"
)

func (e *Env) PriceObservationHandler(r chi.Router) {
	r.Post("/", e.CreatePriceObservationHandler)
}

func (e *Env) CreatePriceObservationHandler(w http.ResponseWriter, r *http.Request) {
	var providedPriceObservation CreatePriceObservationDTO
	ReadJSON(w, r, &providedPriceObservation)
	err := CreatePriceObservationValidator(&providedPriceObservation)
	if err != nil {
		WriteError(w, http.StatusUnprocessableEntity, err.Error())
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

	err = e.DB.QueryRow(
		r.Context(),
		`
		SELECT id
		FROM products
		WHERE id = $1
		AND deleted_at IS NULL
		`,
		providedPriceObservation.ProductID,
	).Scan(nil)
	if err != nil {
		transaction.Rollback(r.Context())
		if errors.Is(err, sql.ErrNoRows) {
			WriteError(w, http.StatusNotFound, "Produto não encontrado")
			return
		}
		fmt.Printf("Database error: %v\n", err)
		WriteError(w, http.StatusInternalServerError, "Erro interno no banco de dados")
		return
	}

	err = e.DB.QueryRow(
		r.Context(),
		`
		SELECT id
		FROM establishments
		WHERE id = $1
		AND deleted_at IS NULL
		`,
		providedPriceObservation.EstablishmentID,
	).Scan(nil)
	if err != nil {
		transaction.Rollback(r.Context())
		if errors.Is(err, sql.ErrNoRows) {
			WriteError(w, http.StatusNotFound, "Estabelecimento não encontrado")
			return
		}
		fmt.Printf("Database error: %v\n", err)
		WriteError(w, http.StatusInternalServerError, "Erro interno no banco de dados")
		return
	}

	var createdPriceObservation models.PriceObservation
	sessionHousehold, err := e.GetSessionUserHousehold(r.Context())

	err = e.DB.QueryRow(
		r.Context(),
		`
		INSERT INTO price_observations (household_id, product_id, establishment_id, observed_price, observed_at)
		VALUES ($1, $2, $3, $4, NOW())
		RETURNING id, household_id, product_id, establishment_id, observed_price, observed_at
		`,
		&sessionHousehold.ID,
		&providedPriceObservation.ProductID,
		&providedPriceObservation.EstablishmentID,
		&providedPriceObservation.Price,
	).Scan(
		&createdPriceObservation.ID,
		&createdPriceObservation.HouseholdID,
		&createdPriceObservation.ProductID,
		&createdPriceObservation.EstablishmentID,
		&createdPriceObservation.ObservedPrice,
		&createdPriceObservation.ObservedAt,
	)
	if err != nil {
		transaction.Rollback(r.Context())
		fmt.Printf("Database error: %v\n", err)
		WriteError(w, http.StatusInternalServerError, "Erro interno no banco de dados")
		return
	}

	transaction.Commit(r.Context())
	WriteJSON(w, http.StatusCreated, createdPriceObservation, "Observação de preço registrada")
}
