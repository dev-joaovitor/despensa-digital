package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/dev-joaovitor/despensa-digital/models"
	"github.com/go-chi/chi/v5"
)

func (e *Env) PriceObservationHandler(r chi.Router) {
	r.Post("/", e.CreatePriceObservationHandler)
	r.Get("/", e.ListPriceObservationsHandler)
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

func (e *Env) ListPriceObservationsHandler(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	productId, err := strconv.ParseInt(params.Get("product_id"), int(10), 64)
	if err != nil {
		WriteError(w, http.StatusUnprocessableEntity, "product_id: ID do produto deve ser numérico")
		return
	}

	from := time.Now().UTC().AddDate(0, 0, -7)
	if params.Get("from") != ""  {
		from, err = time.Parse(time.DateOnly, params.Get("from"))
		if err != nil {
			WriteError(w, http.StatusUnprocessableEntity, "from: Data inválida")
			return
		}
	}

	to := time.Now().UTC()
	if params.Get("to") != ""  {
		to, err = time.Parse(time.DateOnly, params.Get("to"))
		if err != nil {
			WriteError(w, http.StatusUnprocessableEntity, "to: Data inválida")
			return
		}
	}

	if from.Compare(to) == 1 {
		WriteError(w, http.StatusUnprocessableEntity, "`from` deve ser anterior ao `to`")
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
		productId,
	).Scan(nil)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			WriteError(w, http.StatusNotFound, "Produto não encontrado")
			return
		}
		fmt.Printf("Database error: %v\n", err)
		WriteError(w, http.StatusInternalServerError, "Erro interno no banco de dados")
		return
	}

	rows, err := e.DB.Query(
		r.Context(),
		`
		SELECT po.id, po.product_id, e.id, e.name, e.created_at,
			e.updated_at, po.observed_price, po.observed_at
		FROM price_observations po
		LEFT JOIN establishments e
			ON e.id = po.establishment_id
			AND e.deleted_at IS NULL
		WHERE po.deleted_at IS NULL
		AND po.product_id = $1
		AND (
			DATE(po.observed_at) >= $2
			AND DATE(po.observed_at) <= $3
		)
		`,
		productId,
		from.Format(time.DateOnly),
		to.Format(time.DateOnly),
	)
	if err != nil {
		fmt.Printf("Database error: %v\n", err)
		WriteError(w, http.StatusInternalServerError, "Erro interno no banco de dados")
		return
	}
	defer rows.Close()

	priceObservations := []ListPriceObservationsDTO{}
	for rows.Next() {
		priceObservation := ListPriceObservationsDTO{}
		err = rows.Scan(
			&priceObservation.ID,
			&priceObservation.ProductID,
			&priceObservation.Establishment.ID,
			&priceObservation.Establishment.Name,
			&priceObservation.Establishment.CreatedAt,
			&priceObservation.Establishment.UpdatedAt,
			&priceObservation.ObservedPrice,
			&priceObservation.ObservedAt,
		)
		if err != nil {
			fmt.Printf("Database error: %v\n", err)
			WriteError(w, http.StatusInternalServerError, "Erro interno no banco de dados")
			return
		}
		priceObservations = append(priceObservations, priceObservation)
	}

	WriteJSON(w, http.StatusOK, priceObservations, "Observações de preço listados com sucesso")
}
