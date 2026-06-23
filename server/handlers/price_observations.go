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

func (e *Env) PriceObservationsHandler(r chi.Router) {
	r.Post("/", e.CreatePriceObservationHandler)
	r.Get("/", e.ListPriceObservationsHandler)
	r.Get("/history/{productid}", e.PriceObservationsHistoryHandler)
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
		WriteError(w, http.StatusInternalServerError, "Erro interno no banco de dados")
		return
	}
	defer transaction.Rollback(r.Context())

	err = transaction.QueryRow(
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
		if errors.Is(err, sql.ErrNoRows) {
			WriteError(w, http.StatusNotFound, "Produto não encontrado")
			return
		}
		fmt.Printf("Database error: %v\n", err)
		WriteError(w, http.StatusInternalServerError, "Erro interno no banco de dados")
		return
	}

	err = transaction.QueryRow(
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

	err = transaction.QueryRow(
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
		fmt.Printf("Database error: %v\n", err)
		WriteError(w, http.StatusInternalServerError, "Erro interno no banco de dados")
		return
	}

	transaction.Commit(r.Context())
	WriteJSON(w, http.StatusCreated, createdPriceObservation, "Observação de preço registrada")
}

func (e *Env) PriceObservationsHistoryHandler(w http.ResponseWriter, r *http.Request) {
	productId, err := strconv.ParseInt(r.PathValue("productid"), int(10), 64)
	if err != nil {
		WriteError(w, http.StatusUnprocessableEntity, "ID do produto deve ser numérico")
		return
	}

	params := r.URL.Query()
	establishmentId, err := strconv.ParseInt(params.Get("establishment_id"), int(10), 64)
	if err != nil {
		WriteError(w, http.StatusUnprocessableEntity, "establishment_id: ID do estabelecimento deve ser numérico")
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

	err = e.DB.QueryRow(
		r.Context(),
		`
		SELECT id
		FROM establishments
		WHERE id = $1
		AND deleted_at IS NULL
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
		AND po.establishment_id = $4
		`,
		productId,
		from.Format(time.DateOnly),
		to.Format(time.DateOnly),
		establishmentId,
	)
	if err != nil {
		fmt.Printf("Database error: %v\n", err)
		WriteError(w, http.StatusInternalServerError, "Erro interno no banco de dados")
		return
	}
	defer rows.Close()

	priceObservations := []HistoryPriceObservationsDTO{}
	for rows.Next() {
		priceObservation := HistoryPriceObservationsDTO{}
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

	WriteJSON(w, http.StatusOK, priceObservations, "Histórico de Observações de preço listados com sucesso")
}

func (e *Env) ListPriceObservationsHandler(w http.ResponseWriter, r *http.Request) {
	sessionHousehold, err := e.GetSessionUserHousehold(r.Context())
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Residência não encontrada")
		return
	}

	query := r.URL.Query()
	searchQuery := "AND $2 = $2"

	search := query.Get("search")
	if search != "" {
		search = "%" + search + "%"
		searchQuery = `
		AND (p.name ILIKE $2 OR b.name ILIKE $2)
		`
	}

	priceObservations := []ListPriceObservationsDTO{}
	rows, err := e.DB.Query(
		r.Context(),
		`
		WITH windowed_observations AS (
			SELECT
				observed_at,
				product_id,
				observed_price,
				establishment_id,
				ROW_NUMBER() OVER(PARTITION BY product_id ORDER BY observed_price ASC) AS lowest_row,
				ROW_NUMBER() OVER(PARTITION BY product_id ORDER BY observed_at DESC) AS latest_row,
				MIN(observed_price) OVER(PARTITION BY product_id) AS lowest_price,
				AVG(observed_price) OVER(PARTITION BY product_id) AS average_price
			FROM price_observations
			WHERE deleted_at IS NULL
			AND household_id = $1
		)
		SELECT
			p.id,
			p.name product_name,
			b.name brand_name,
			p.unit_size,
			um.acronym unit_acronym,
			latest_obs.average_price,

			latest_obs.lowest_price,
			lowest_obs.observed_at lowest_observed_at,
			lowest_est.name lowest_establishment_name,

			latest_obs.observed_price current_price,
			latest_obs.observed_at current_observed_at,
			latest_est.name latest_establishment_name
		FROM products p

		JOIN brands b ON p.brand_id = b.id

		JOIN unit_measurements um ON p.measurement_id = um.id

		JOIN windowed_observations latest_obs
		ON p.id = latest_obs.product_id
		AND latest_obs.latest_row = 1

		JOIN windowed_observations lowest_obs
		ON p.id = lowest_obs.product_id
		AND lowest_obs.lowest_row = 1

		JOIN establishments latest_est
		ON latest_est.id = latest_obs.establishment_id

		JOIN establishments lowest_est
		ON lowest_est.id = lowest_obs.establishment_id

		WHERE p.deleted_at IS NULL
		AND p.household_id = $1
		` + searchQuery,
		&sessionHousehold.ID,
		search,
	)
	if err != nil {
		fmt.Printf("Database error: %v\n", err)
		WriteError(w, http.StatusInternalServerError, "Erro interno no banco de dados")
		return
	}
	defer rows.Close()

	for rows.Next() {
		priceObservation := ListPriceObservationsDTO{}

		err = rows.Scan(
			&priceObservation.Product.ID,
			&priceObservation.Product.Name,
			&priceObservation.Product.Brand.Name,
			&priceObservation.Product.Measurement.Size,
			&priceObservation.Product.Measurement.Acronym,
			&priceObservation.AverageObservedPrice,
			&priceObservation.Lowest.ObservedPrice,
			&priceObservation.Lowest.ObservedAt,
			&priceObservation.Lowest.Establishment.Name,
			&priceObservation.Current.ObservedPrice,
			&priceObservation.Current.ObservedAt,
			&priceObservation.Current.Establishment.Name,
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

