package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/dev-joaovitor/despensa-digital/models"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
)

func (e *Env) ShoppingListHandler(r chi.Router) {
	r.Get("/", e.ListShoppingItemsHandler)
	r.Post("/", e.CreateShoppingItemHandler)
	r.Patch("/{id}", e.UpdateItemHandler)
	r.Delete("/{id}", e.DeleteShoppingItemHandler)
	r.Post("/{id}/tick", e.TickShoppingItemHandler)
	r.Post("/submit", e.SubmitShoppingListHandler)
}

func (e *Env) CreateShoppingItemHandler(w http.ResponseWriter, r *http.Request) {
	var err error

	var providedShoppingItem CreateShoppingItemDTO
	ReadJSON(w, r, &providedShoppingItem)
	err = CreateShoppingItemValidator(&providedShoppingItem)

	if err != nil {
		WriteError(w, http.StatusUnprocessableEntity, "Erro de validação: " + err.Error())
		return
	}

	transaction, err := e.DB.Begin(r.Context())
	if err != nil {
		fmt.Printf("Database error: %v\n", err)
		WriteError(w, http.StatusInternalServerError, "Erro interno no banco de dados")
		return
	}
	defer transaction.Rollback(r.Context())

	sessionHousehold, err := e.GetSessionUserHousehold(r.Context())
	if err != nil {
		fmt.Printf("Session error: %v\n", err)
		WriteError(w, http.StatusInternalServerError, "Residência não encontrada")
		return
	}

	err = transaction.QueryRow(
		r.Context(),
		`
		SELECT id
		FROM products
		WHERE id = $1
		AND deleted_at IS NULL
		`,
		&providedShoppingItem.ProductID,
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
		FROM shopping_list_items
		WHERE product_id = $1
		AND deleted_at IS NULL
		`,
		&providedShoppingItem.ProductID,
	).Scan(nil)
	if err == nil {
		WriteError(w, http.StatusForbidden, "O produto já existe na lista")
		return
	} else {
		if !errors.Is(err, sql.ErrNoRows) {
			fmt.Printf("Database error: %v\n", err)
			WriteError(w, http.StatusInternalServerError, "Erro interno no banco de dados")
			return
		}
	}

	var createdShoppingItem models.ShoppingListItem
	err = transaction.QueryRow(
		r.Context(),
		`
		INSERT INTO shopping_list_items (household_id, product_id,
			quantity, is_checked)
		VALUES ($1, $2, $3, false)
		RETURNING id, household_id, product_id, quantity,
			is_checked, created_at, updated_at
		`,
		&sessionHousehold.ID,
		&providedShoppingItem.ProductID,
		&providedShoppingItem.Quantity,
	).Scan(
		&createdShoppingItem.ID,
		&createdShoppingItem.HouseholdID,
		&createdShoppingItem.ProductID,
		&createdShoppingItem.Quantity,
		&createdShoppingItem.IsChecked,
		&createdShoppingItem.CreatedAt,
		&createdShoppingItem.UpdatedAt,
	)
	if err != nil {
		WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	transaction.Commit(r.Context())
	WriteJSON(w, http.StatusCreated, &createdShoppingItem, "Item adicionado com sucesso")
}

func (e *Env) ListShoppingItemsHandler(w http.ResponseWriter, r *http.Request) {
	sessionHousehold, err := e.GetSessionUserHousehold(r.Context())
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Residência não encontrada")
		return
	}

	rows, err := e.DB.Query(
		r.Context(),
		`
		SELECT i.id, p.id, p.name, p.unit_size,
			um.acronym, b.name, i.quantity, i.is_checked,
			i.created_at, i.updated_at
		FROM shopping_list_items i
		JOIN products p ON i.product_id = p.id
		JOIN brands b ON p.brand_id = b.id
		JOIN unit_measurements um ON p.measurement_id = um.id
		WHERE i.deleted_at IS NULL
		AND i.household_id = $1
		ORDER BY i.created_at, i.updated_at DESC
		`,
		&sessionHousehold.ID,
	)
	if err != nil {
		fmt.Printf("Database error: %v\n", err)
		WriteError(w, http.StatusInternalServerError, "Erro interno no banco de dados")
		return
	}
	defer rows.Close()

	items := []ListShoppingItemsDTO{}

	for rows.Next() {
		item := ListShoppingItemsDTO{}
		rows.Scan(
			&item.ID,
			&item.Product.ID,
			&item.Product.Name,
			&item.Product.Measurement.Size,
			&item.Product.Measurement.Acronym,
			&item.Product.Brand.Name,
			&item.Quantity,
			&item.IsChecked,
			&item.CreatedAt,
			&item.UpdatedAt,
		)
		items = append(items, item)
	}

	WriteJSON(w, http.StatusOK, items, "Items da lista de compras listados com sucesso")
}

func (e *Env) TickShoppingItemHandler(w http.ResponseWriter, r *http.Request) {
	itemId := r.PathValue("id")
	if itemId == "" {
		WriteError(w, http.StatusBadRequest, "Insira um ID")
		return
	}

	transaction, err := e.DB.Begin(r.Context())
	if err != nil {
		fmt.Printf("Database error: %v\n", err)
		WriteError(w, http.StatusInternalServerError, "Erro interno no banco de dados")
		return
	}
	defer transaction.Rollback(r.Context())

	err = transaction.QueryRow(
		r.Context(),
		`
		UPDATE shopping_list_items
		SET is_checked = NOT is_checked
		WHERE id = $1
		AND deleted_at IS NULL
		RETURNING id
		`,
		itemId,
	).Scan(nil)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			WriteError(w, http.StatusNotFound, "Item não encontrado")
			return
		}
		fmt.Printf("Database error: %v\n", err)
		WriteError(w, http.StatusInternalServerError, "Erro interno no banco de dados")
		return
	}

	transaction.Commit(r.Context())
	WriteJSON(w, http.StatusOK, nil, "Item atualizado com sucesso")
}

func (e *Env) UpdateItemHandler(w http.ResponseWriter, r *http.Request) {
	itemId := r.PathValue("id")
	if itemId == "" {
		WriteError(w, http.StatusBadRequest, "Insira um ID")
		return
	}

	var providedShoppingItem UpdateShoppingItemDTO
	ReadJSON(w, r, &providedShoppingItem)
	err := UpdateShoppingItemValidator(&providedShoppingItem)

	if err != nil {
		WriteError(w, http.StatusUnprocessableEntity, "Erro de validação: " + err.Error())
		return
	}

	transaction, err := e.DB.Begin(r.Context())
	if err != nil {
		fmt.Printf("Database error: %v\n", err)
		WriteError(w, http.StatusInternalServerError, "Erro interno no banco de dados")
		return
	}
	defer transaction.Rollback(r.Context())

	err = transaction.QueryRow(
		r.Context(),
		`
		UPDATE shopping_list_items
		SET quantity = $1
		WHERE id = $2
		AND deleted_at IS NULL
		RETURNING id
		`,
		&providedShoppingItem.Quantity,
		itemId,
	).Scan(nil)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			WriteError(w, http.StatusNotFound, "Item não encontrado")
			return
		}
		fmt.Printf("Database error: %v\n", err)
		WriteError(w, http.StatusInternalServerError, "Erro interno no banco de dados")
		return
	}

	transaction.Commit(r.Context())
	WriteJSON(w, http.StatusOK, nil, "Item atualizado com sucesso")
}

func (e *Env) DeleteShoppingItemHandler(w http.ResponseWriter, r *http.Request) {
	itemId := r.PathValue("id")
	if itemId == "" {
		WriteError(w, http.StatusBadRequest, "Insira um ID")
		return
	}

	transaction, err := e.DB.Begin(r.Context())
	if err != nil {
		fmt.Printf("Database error: %v\n", err)
		WriteError(w, http.StatusInternalServerError, "Erro interno no banco de dados")
		return
	}
	defer transaction.Rollback(r.Context())

	err = transaction.QueryRow(
		r.Context(),
		`
		UPDATE shopping_list_items
		SET deleted_at = NOW()
		WHERE id = $1
		AND deleted_at IS NULL
		RETURNING id
		`,
		itemId,
	).Scan(nil)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			WriteError(w, http.StatusNotFound, "Item não encontrado")
			return
		}
		fmt.Printf("Database error: %v\n", err)
		WriteError(w, http.StatusInternalServerError, "Erro interno no banco de dados")
		return
	}

	transaction.Commit(r.Context())
	WriteJSON(w, http.StatusOK, nil, "Item removido com sucesso")
}

func (e *Env) SubmitShoppingListHandler(w http.ResponseWriter, r *http.Request) {
	sessionHousehold, err := e.GetSessionUserHousehold(r.Context())
	if err != nil {
		fmt.Printf("Session error: %v\n", err)
		WriteError(w, http.StatusInternalServerError, "Residência não encontrada")
		return
	}

	var providedShoppingList SubmitShoppingListDTO
	ReadJSON(w, r, &providedShoppingList)
	err = SubmitShoppingListValidator(&providedShoppingList)

	if err != nil {
		WriteError(w, http.StatusUnprocessableEntity, "Erro de validação: " + err.Error())
		return
	}
	
	var productsIds []int64
	var establishmentIds []int64
	tempStockBatchRows := make([][]any, len(providedShoppingList.Items))
	priceObservationRows := make([][]any, len(providedShoppingList.Items))
	for i, item  := range providedShoppingList.Items {
		productsIds = append(productsIds, item.ProductID)
		establishmentIds = append(establishmentIds, item.EstablishmentID)
		tempStockBatchRows[i] = []any{
			sessionHousehold.ID,
			item.ProductID,
			item.EstablishmentID,
			item.Quantity, // initial
			item.Quantity, // remaining
			item.Price,
			item.ExpirationDate,
		}
		priceObservationRows[i] = []any{
			sessionHousehold.ID,
			item.ProductID,
			item.EstablishmentID,
			item.Price,
		}
	}

	var count int
	err = e.DB.QueryRow(
		r.Context(),
		`SELECT COUNT(1) FROM products p
		JOIN shopping_list_items i
		ON p.id = i.product_id
		AND i.is_checked = true
		AND i.deleted_at IS NULL
		WHERE p.id = ANY($1)
		AND p.deleted_at IS NULL`,
		productsIds,
	).Scan(&count)
	if err != nil {
		fmt.Printf("Database error count products: %v\n", err)
		WriteError(w, http.StatusInternalServerError, "Validação de produtos falhou")
		return
	}
	if count != len(productsIds) {
		WriteError(w, http.StatusBadRequest, "Um ou mais produtos não existem, não foram riscados ou já foram enviados")
		return
	}

	err = e.DB.QueryRow(
		r.Context(),
		`SELECT COUNT(1) FROM establishments
		WHERE id = ANY($1)
		AND deleted_at IS NULL`,
		establishmentIds,
	).Scan(&count)
	if err != nil {
		fmt.Printf("Database error count establishments: %v\n", err)
		WriteError(w, http.StatusInternalServerError, "Validação de estabelecimentos falhou")
		return
	}
	if count != len(establishmentIds) {
		WriteError(w, http.StatusBadRequest, "Um ou mais estabelecimentos não existem")
		return
	}

	transaction, err := e.DB.Begin(r.Context())
	if err != nil {
		fmt.Printf("Database error begin transaction: %v\n", err)
		WriteError(w, http.StatusInternalServerError, "Erro interno no banco de dados")
		return
	}
	defer transaction.Rollback(r.Context())

	_, err = transaction.Exec(
		r.Context(),
		`
		CREATE TEMP TABLE temp_stock_batches
		(LIKE stock_batches INCLUDING ALL)
		ON COMMIT DROP;
		`,
	)
	if err != nil {
		fmt.Printf("Database error create temp stock batches: %v\n", err)
		WriteError(w, http.StatusInternalServerError, "Erro interno no banco de dados")
		return
	}

	_, err = transaction.CopyFrom(
		r.Context(),
		pgx.Identifier{"temp_stock_batches"},
		[]string{
			"household_id",
			"product_id",
			"establishment_id",
			"initial_quantity",
			"remaining_quantity",
			"unit_price",
			"expiration_date",
		},
		pgx.CopyFromRows(tempStockBatchRows),
	)
	if err != nil {
		fmt.Printf("Database error insert into temp stock batches: %v\n", err)
		WriteError(w, http.StatusInternalServerError, "Erro interno no banco de dados")
		return
	}

	stockBatchesRows, err := transaction.Query(
		r.Context(),
		`
		INSERT INTO stock_batches (
			household_id,
			product_id,
			establishment_id,
			initial_quantity,
			remaining_quantity,
			unit_price,
			expiration_date
		)
		SELECT
			household_id,
			product_id,
			establishment_id,
			initial_quantity,
			remaining_quantity,
			unit_price,
			expiration_date
		FROM temp_stock_batches
		RETURNING id, initial_quantity, $1
		`,
		models.TransactionPurchase,
	)
	if err != nil {
		fmt.Printf("Database error stock batches: %v\n", err)
		WriteError(w, http.StatusInternalServerError, "Erro interno no banco de dados")
		return
	}
	
	stockTransactionRows := make([][]any, len(productsIds))
	i := 0
	for stockBatchesRows.Next() {
		stockTransactionValues := make([]any, 3)
		err = stockBatchesRows.Scan(
			&stockTransactionValues[0],
			&stockTransactionValues[1],
			&stockTransactionValues[2],
		)
		if err != nil {
			fmt.Printf("Database error: %v\n", err)
			WriteError(w, http.StatusInternalServerError, "Erro interno no banco de dados")
			return
		}

		stockTransactionRows[i] = stockTransactionValues
		i++
	}

	_, err = transaction.CopyFrom(
		r.Context(),
		pgx.Identifier{"stock_transactions"},
		[]string{
			"batch_id",
			"quantity",
			"type",
		},
		pgx.CopyFromRows(stockTransactionRows),
	)
	if err != nil {
		fmt.Printf("Database error temp stock batches: %v\n", err)
		WriteError(w, http.StatusInternalServerError, "Erro interno no banco de dados")
		return
	}

	err = transaction.QueryRow(
		r.Context(),
		`
		UPDATE shopping_list_items
		SET deleted_at = NOW()
		WHERE id = ANY($1)
		AND deleted_at IS NULL
		RETURNING id
		`,
		productsIds,
	).Scan(nil)
	if err != nil {
		fmt.Printf("Database error update shopping list: %v\n", err)
		WriteError(w, http.StatusInternalServerError, "Erro interno no banco de dados")
		return
	}

	_, err = transaction.CopyFrom(
		r.Context(),
		pgx.Identifier{"price_observations"},
		[]string{
			"household_id",
			"product_id",
			"establishment_id",
			"observed_price",
		},
		pgx.CopyFromRows(priceObservationRows),
	)
	if err != nil {
		fmt.Printf("Database error insert price observations: %v\n", err)
		WriteError(w, http.StatusInternalServerError, "Erro interno no banco de dados")
		return
	}

	transaction.Commit(r.Context())
	WriteJSON(w, http.StatusOK, nil, "Lista enviada com sucesso")
}
