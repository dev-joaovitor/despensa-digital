package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/dev-joaovitor/despensa-digital/models"
	"github.com/go-chi/chi/v5"
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
		if transaction != nil {
			transaction.Rollback(r.Context())
		}

		fmt.Printf("Database error: %v\n", err)
		WriteError(w, http.StatusInternalServerError, "Erro interno no banco de dados")
		return
	}

	sessionHousehold, err := e.GetSessionUserHousehold(r.Context())
	if err != nil {
		transaction.Rollback(r.Context())
		fmt.Printf("Session error: %v\n", err)
		WriteError(w, http.StatusInternalServerError, "Residência não encontrada")
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
		&providedShoppingItem.ProductID,
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
		FROM shopping_list_items
		WHERE product_id = $1
		AND deleted_at IS NULL
		`,
		&providedShoppingItem.ProductID,
	).Scan(nil)
	if err == nil {
		transaction.Rollback(r.Context())
		WriteError(w, http.StatusForbidden, "O produto já existe na lista")
		return
	} else {
		if !errors.Is(err, sql.ErrNoRows) {
			transaction.Rollback(r.Context())
			fmt.Printf("Database error: %v\n", err)
			WriteError(w, http.StatusInternalServerError, "Erro interno no banco de dados")
			return
		}
	}

	var createdShoppingItem models.ShoppingListItem
	err = e.DB.QueryRow(
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
		transaction.Rollback(r.Context())
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
		if transaction != nil {
			transaction.Rollback(r.Context())
		}

		fmt.Printf("Database error: %v\n", err)
		WriteError(w, http.StatusInternalServerError, "Erro interno no banco de dados")
		return
	}

	err = e.DB.QueryRow(
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
		transaction.Rollback(r.Context())
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
		if transaction != nil {
			transaction.Rollback(r.Context())
		}

		fmt.Printf("Database error: %v\n", err)
		WriteError(w, http.StatusInternalServerError, "Erro interno no banco de dados")
		return
	}

	err = e.DB.QueryRow(
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
		transaction.Rollback(r.Context())
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
		if transaction != nil {
			transaction.Rollback(r.Context())
		}

		fmt.Printf("Database error: %v\n", err)
		WriteError(w, http.StatusInternalServerError, "Erro interno no banco de dados")
		return
	}

	err = e.DB.QueryRow(
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
		transaction.Rollback(r.Context())
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

func (e *Env) SubmitShoppingListHandler(w http.ResponseWriter, r *http.Request) {}
