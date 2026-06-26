package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/dev-joaovitor/despensa-digital/models"
	"github.com/go-chi/chi/v5"
)

// lowStockRestockThreshold: a product is auto-added to the shopping list once its
// total remaining stock drops below this fraction of its latest batch's initial
// quantity. Provisional until a per-product threshold field exists.
const lowStockRestockThreshold = 0.2

func (e *Env) StockHandler(r chi.Router) {
	r.Get("/products", e.ListStockProductsHandler)
	r.Get("/products/{productid}", e.GetStockProductHandler)
	r.Get("/products/{productid}/batches", e.ListStockProductBatchesHandler)
	r.Post("/transact", e.TransactStockBatchHandler)
}

func (e *Env) ListStockProductsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	sessionHousehold, err := e.GetSessionUserHousehold(ctx)
	if err != nil {
		fmt.Printf("Error getting household: %v\n", err)
		WriteError(w, http.StatusInternalServerError, "Erro interno")
		return
	}

	query := r.URL.Query()
	searchQuery := "AND $2 = $2"

	search := query.Get("search")
	if search != "" {
		search = "%" + search + "%"
		searchQuery = `
		AND (
			p.name ILIKE $2 OR b.name ILIKE $2
			OR c.name ILIKE $2
		)
		`
	}

	rows, err := e.DB.Query(
		r.Context(),
		`
		SELECT p.id, p.name, p.created_at, p.updated_at, b.id, b.name, p.unit_size,
			um.id, um.acronym, c.id, c.name, COALESCE(SUM(sb.initial_quantity), 0) initial,
			COALESCE(SUM(sb.remaining_quantity), 0) remaining
		FROM products p
		JOIN brands b ON p.brand_id = b.id
		JOIN unit_measurements um ON p.measurement_id = um.id
		JOIN categories c ON p.category_id = c.id
		LEFT JOIN stock_batches sb ON sb.product_id = p.id AND sb.remaining_quantity > 0
		WHERE p.deleted_at IS NULL
		AND p.household_id = $1
		` + searchQuery + `
		GROUP BY p.id, b.id, b.name, um.id, um.acronym, c.id, c.name
		ORDER BY p.created_at, p.updated_at DESC
		`,
		&sessionHousehold.ID,
		search,
	)
	if err != nil {
		fmt.Printf("Database error: %v\n", err)
		WriteError(w, http.StatusInternalServerError, "Erro interno no banco de dados")
		return
	}
	defer rows.Close()

	products := []ListStockProductsDTO{}

	for rows.Next() {
		product := ListStockProductsDTO{}
		rows.Scan(
			&product.ID,
			&product.Name,
			&product.CreatedAt,
			&product.UpdatedAt,
			&product.Brand.ID,
			&product.Brand.Name,
			&product.Measurement.Size,
			&product.Measurement.ID,
			&product.Measurement.Acronym,
			&product.Category.ID,
			&product.Category.Name,
			&product.Stock.Initial,
			&product.Stock.Remaining,
		)
		products = append(products, product)
	}
	
	WriteJSON(w, http.StatusOK, products, "Produtos da despensa listados com sucesso")
}

func (e *Env) GetStockProductHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	sessionHousehold, err := e.GetSessionUserHousehold(ctx)
	if err != nil {
		fmt.Printf("Error getting household: %v\n", err)
		WriteError(w, http.StatusInternalServerError, "Erro interno")
		return
	}

	productId := r.PathValue("productid")
	if productId == "" {
		WriteError(w, http.StatusBadRequest, "Insira o ID do produto")
		return
	}

	product := ListStockProductsDTO{}
	err = e.DB.QueryRow(
		r.Context(),
		`
		SELECT p.id, p.name, p.created_at, p.updated_at, b.id, b.name, p.unit_size,
			um.id, um.acronym, c.id, c.name, COALESCE(SUM(sb.initial_quantity), 0) initial,
			COALESCE(SUM(sb.remaining_quantity), 0) remaining
		FROM products p
		JOIN brands b ON p.brand_id = b.id
		JOIN unit_measurements um ON p.measurement_id = um.id
		JOIN categories c ON p.category_id = c.id
		LEFT JOIN stock_batches sb ON sb.product_id = p.id AND sb.remaining_quantity > 0
		WHERE p.deleted_at IS NULL
		AND p.household_id = $1
		AND p.id = $2
		GROUP BY p.id, b.id, b.name, um.id, um.acronym, c.id, c.name
		`,
		&sessionHousehold.ID,
		productId,
	).Scan(
		&product.ID,
		&product.Name,
		&product.CreatedAt,
		&product.UpdatedAt,
		&product.Brand.ID,
		&product.Brand.Name,
		&product.Measurement.Size,
		&product.Measurement.ID,
		&product.Measurement.Acronym,
		&product.Category.ID,
		&product.Category.Name,
		&product.Stock.Initial,
		&product.Stock.Remaining,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			WriteError(w, http.StatusNotFound, "Produto não encontrado")
			return
		}
		fmt.Printf("Database error: %v\n", err)
		WriteError(w, http.StatusInternalServerError, "Erro interno no banco de dados")
		return
	}

	WriteJSON(w, http.StatusOK, product, "Produto da despensa listado com sucesso")
}

func (e *Env) ListStockProductBatchesHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	sessionHousehold, err := e.GetSessionUserHousehold(ctx)
	if err != nil {
		fmt.Printf("Error getting household: %v\n", err)
		WriteError(w, http.StatusInternalServerError, "Erro interno")
		return
	}

	productId := r.PathValue("productid")
	if productId == "" {
		WriteError(w, http.StatusBadRequest, "Insira o ID do produto")
		return
	}

	rows, err := e.DB.Query(
		r.Context(),
		`
		SELECT sb.id, sb.unit_price, sb.initial_quantity, sb.remaining_quantity,
			sb.expiration_date, sb.created_at, sb.updated_at, e.name
		FROM stock_batches sb
		JOIN establishments e
		ON e.id = sb.establishment_id
		WHERE sb.product_id = $2
		AND sb.household_id = $1
		ORDER BY sb.created_at, sb.updated_at DESC
		`,
		&sessionHousehold.ID,
		productId,
	)
	if err != nil {
		fmt.Printf("Database error: %v\n", err)
		WriteError(w, http.StatusInternalServerError, "Erro interno no banco de dados")
		return
	}
	defer rows.Close()

	var batches []ListStockProductBatchesDTO
	for rows.Next() {
		var batch ListStockProductBatchesDTO
		rows.Scan(
			&batch.ID,
			&batch.UnitPrice,
			&batch.InitialQuantity,
			&batch.RemainingQuantity,
			&batch.ExpirationDate,
			&batch.CreatedAt,
			&batch.UpdatedAt,
			&batch.Establishment.Name,
		)
		batches = append(batches, batch)
	}

	WriteJSON(w, http.StatusOK, batches, "Lotes listados com sucesso")
}

func (e *Env) TransactStockBatchHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	sessionHousehold, err := e.GetSessionUserHousehold(ctx)
	if err != nil {
		fmt.Printf("Error getting household: %v\n", err)
		WriteError(w, http.StatusInternalServerError, "Erro interno")
		return
	}

	var providedTransaction TransactStockBatchDTO
	ReadJSON(w, r, &providedTransaction)
	err = TransactStockBatchValidator(&providedTransaction)
	if err != nil {
		WriteError(w, http.StatusUnprocessableEntity, "Erro de validação: " + err.Error())
		return
	}

	transaction, err := e.DB.Begin(ctx)
	if err != nil {
		fmt.Printf("Database error: %v\n", err)
		WriteError(w, http.StatusInternalServerError, "Erro interno no banco de dados")
		return
	}
	defer transaction.Rollback(ctx)


	var batchId int64
	var batchQuery string
	var batchQueryArgs []any

	switch providedTransaction.Type {
	case models.TransactionPurchase:
		batchQuery = `
		INSERT INTO stock_batches (household_id, product_id,
			establishment_id, initial_quantity, remaining_quantity,
			expiration_date, unit_price)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
		`
		batchQueryArgs = []any{
			sessionHousehold.ID,
			providedTransaction.ProductID,
			providedTransaction.EstablishmentID,
			providedTransaction.Quantity,
			providedTransaction.Quantity,
			providedTransaction.ExpirationDate,
			providedTransaction.UnitPrice,
		}
	case models.TransactionCorrection:
		batchQuery = `
		UPDATE stock_batches
		SET establishment_id = $2, initial_quantity = $3,
			remaining_quantity = $4, unit_price = $5,
			expiration_date = $6
		WHERE id = $1
		RETURNING id
		`
		batchQueryArgs = []any{
			providedTransaction.BatchID,
			providedTransaction.EstablishmentID,
			providedTransaction.Quantity,
			providedTransaction.Quantity,
			providedTransaction.UnitPrice,
			providedTransaction.ExpirationDate,
		}
	case models.TransactionConsumption, models.TransactionWaste:
		batchQuery = `
		UPDATE stock_batches
		SET remaining_quantity = remaining_quantity - $2
		WHERE id = $1
		AND (remaining_quantity - $2) >= 0
		RETURNING id
		`
		batchQueryArgs = []any{
			providedTransaction.BatchID,
			providedTransaction.Quantity,
		}
	}

	err = transaction.QueryRow(
		ctx,
		batchQuery,
		batchQueryArgs...,
	).Scan(&batchId)
	if err != nil {
		fmt.Printf("Database error update stock batch: %v\n", err)
		if errors.Is(err, sql.ErrNoRows) {
			WriteError(w, http.StatusBadRequest, "Dados inválidos")
			return
		}
		WriteError(w, http.StatusInternalServerError, "Erro interno no banco de dados")
		return
	}

	err = transaction.QueryRow(
		ctx,
		`
		INSERT INTO stock_transactions (type, batch_id, quantity)
		VALUES ($1, $2, $3)
		RETURNING id
		`,
		providedTransaction.Type,
		batchId,
		providedTransaction.Quantity,
	).Scan(nil)
	if err != nil {
		fmt.Printf("Database error insert stock transactions: %v\n", err)
		WriteError(w, http.StatusInternalServerError, "Erro interno no banco de dados")
		return
	}

	// When a consumption/waste drops a product's live stock low, silently add it to
	// the shopping list so the user is prompted to restock.
	if providedTransaction.Type == models.TransactionConsumption ||
		providedTransaction.Type == models.TransactionWaste {
		_, err = transaction.Exec(
			ctx,
			`
			INSERT INTO shopping_list_items (household_id, product_id, quantity, is_checked)
			SELECT $1, latest.product_id, latest.initial_quantity, false
			FROM (
				SELECT product_id, initial_quantity
				FROM stock_batches
				WHERE product_id = (SELECT product_id FROM stock_batches WHERE id = $2)
				AND household_id = $1
				ORDER BY created_at DESC, id DESC
				LIMIT 1
			) latest
			WHERE (
				SELECT COALESCE(SUM(remaining_quantity), 0)
				FROM stock_batches
				WHERE product_id = latest.product_id
				AND household_id = $1
				AND remaining_quantity > 0
			) <= GREATEST(1, latest.initial_quantity * $3)
			AND NOT EXISTS (
				SELECT 1 FROM shopping_list_items
				WHERE product_id = latest.product_id
				AND household_id = $1
				AND deleted_at IS NULL
			)
			`,
			sessionHousehold.ID,
			batchId,
			lowStockRestockThreshold,
		)
		if err != nil {
			fmt.Printf("Database error auto-add shopping item: %v\n", err)
			WriteError(w, http.StatusInternalServerError, "Erro interno no banco de dados")
			return
		}
	}

	transaction.Commit(ctx)
	WriteJSON(w, http.StatusOK, nil, "Estoque atualizado com sucesso")
}
