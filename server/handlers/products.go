package handlers

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (e *Env) ProductsHandler(r chi.Router) {
	r.Get("/", e.ListProductsHandler)
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
		JOIN categories c ON p.measurement_id = c.id
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
