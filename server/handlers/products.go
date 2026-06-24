package handlers

import (
	"fmt"
	"net/http"

	"github.com/dev-joaovitor/despensa-digital/models"
	"github.com/go-chi/chi/v5"
)

func (e *Env) ProductsHandler(r chi.Router) {
	r.Post("/", e.CreateProductHandler)
	r.Get("/", e.ListProductsHandler)
	r.Patch("/{id}", e.UpdateProductHandler)
}

func (e *Env) CreateProductHandler(w http.ResponseWriter, r *http.Request) {
	var err error

	var providedProduct CreateProductDTO
	ReadJSON(w, r, &providedProduct)
	err = CreateProductValidator(&providedProduct)

	if err != nil {
		WriteError(w, http.StatusUnprocessableEntity, "Erro de validação: " + err.Error())
		return
	}

	transaction, err := e.DB.Begin(r.Context())
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Erro interno no banco de dados")
		return
	}
	defer transaction.Rollback(r.Context())

	sessionHousehold, err := e.GetSessionUserHousehold(r.Context())
	if err != nil {
		WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	var createdProduct models.Product
	err = transaction.QueryRow(
		r.Context(),
		`
		INSERT INTO products (household_id, brand_id, measurement_id,
			category_id, name, unit_size)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, household_id, brand_id, measurement_id,
			category_id, name, unit_size, created_at, updated_at
		`,
		&sessionHousehold.ID,
		&providedProduct.BrandID,
		&providedProduct.MeasurementID,
		&providedProduct.CategoryID,
		&providedProduct.Name,
		&providedProduct.UnitSize,
	).Scan(
		&createdProduct.ID,
		&createdProduct.HouseholdID,
		&createdProduct.BrandID,
		&createdProduct.MeasurementID,
		&createdProduct.CategoryID,
		&createdProduct.Name,
		&createdProduct.UnitSize,
		&createdProduct.CreatedAt,
		&createdProduct.UpdatedAt,
	)
	if err != nil {
		WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	transaction.Commit(r.Context())
	WriteJSON(w, http.StatusCreated, &createdProduct, "Produto criado com sucesso")
}

func (e *Env) UpdateProductHandler(w http.ResponseWriter, r *http.Request) {
	var err error

	var providedProduct UpdateProductDTO
	ReadJSON(w, r, &providedProduct)
	err = UpdateProductValidator(&providedProduct)
	if err != nil {
		WriteError(w, http.StatusUnprocessableEntity, "Erro de validação: " + err.Error())
		return
	}

	productId := r.PathValue("id")
	if productId == "" {
		WriteError(w, http.StatusBadRequest, "Insira um ID")
		return
	}

	transaction, err := e.DB.Begin(r.Context())
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Erro interno no banco de dados")
		return
	}
	defer transaction.Rollback(r.Context())

	sessionHousehold, err := e.GetSessionUserHousehold(r.Context())
	if err != nil {
		WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	_, err = transaction.Exec(
		r.Context(),
		`
		UPDATE products
		SET brand_id = $1, measurement_id = $2, category_id = $3,
			name = $4, unit_size = $5, updated_at = NOW()
		WHERE id = $6
		AND household_id = $7
		`,
		&providedProduct.BrandID,
		&providedProduct.MeasurementID,
		&providedProduct.CategoryID,
		&providedProduct.Name,
		&providedProduct.UnitSize,
		productId,
		&sessionHousehold.ID,
	)
	if err != nil {
		WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	transaction.Commit(r.Context())
	WriteJSON(w, http.StatusOK, nil, "Produto atualizado com sucesso")
}

func (e *Env) ListProductsHandler(w http.ResponseWriter, r *http.Request) {
	sessionHousehold, err := e.GetSessionUserHousehold(r.Context())
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Residência não encontrada")
		return
	}

	rows, err := e.DB.Query(
		r.Context(),
		`
		SELECT p.id, p.name, p.created_at, p.updated_at, b.name, p.unit_size,
			um.acronym, c.name
		FROM products p
		JOIN brands b ON p.brand_id = b.id
		JOIN unit_measurements um ON p.measurement_id = um.id
		JOIN categories c ON p.category_id = c.id
		WHERE p.deleted_at IS NULL
		AND p.household_id = $1
		ORDER BY p.created_at, p.updated_at DESC
		`,
		&sessionHousehold.ID,
	)
	if err != nil {
		fmt.Printf("Database error: %v\n", err)
		WriteError(w, http.StatusInternalServerError, "Erro interno no banco de dados")
		return
	}
	defer rows.Close()

	products := []ListProductsDTO{}

	for rows.Next() {
		product := ListProductsDTO{}
		rows.Scan(
			&product.ID,
			&product.Name,
			&product.CreatedAt,
			&product.UpdatedAt,
			&product.Brand.Name,
			&product.Measurement.Size,
			&product.Measurement.Acronym,
			&product.Category.Name,
		)
		products = append(products, product)
	}

	WriteJSON(w, http.StatusOK, products, "Produtos listados com sucesso")
}
