package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/dev-joaovitor/despensa-digital/models"
	"github.com/go-chi/chi/v5"
)

func (e *Env) BrandsHandler(r chi.Router) {
	r.Get("/{id}", e.ListOneBrandHandler)
	r.Get("/", e.ListBrandsHandler)
	r.Post("/", e.CreateBrandHandler)
	r.Patch("/{id}", e.UpdateBrandHandler)
	r.Delete("/{id}", e.DeleteBrandHandler)
}

func (e *Env) CreateBrandHandler(w http.ResponseWriter, r *http.Request) {
	var err error

	var providedBrand CreateBrandDTO
	ReadJSON(w, r, &providedBrand)
	err = CreateBrandValidator(&providedBrand)

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

	var createdBrand models.Brand
	userHousehold, err := e.GetSessionUserHousehold(r.Context())
	if err != nil {
		transaction.Rollback(r.Context())
		WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = e.DB.QueryRow(
		r.Context(),
		`
		INSERT INTO brands (household_id, name)
		VALUES ($1, $2)
		RETURNING id, household_id, name, created_at, updated_at
		`,
		&userHousehold.ID,
		&providedBrand.Name,
	).Scan(
		&createdBrand.ID,
		&createdBrand.HouseholdID,
		&createdBrand.Name,
		&createdBrand.CreatedAt,
		&createdBrand.UpdatedAt,
	)

	if err != nil {
		transaction.Rollback(r.Context())
		WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	transaction.Commit(r.Context())
	WriteJSON(w, http.StatusCreated, &createdBrand, "Marca criada com sucesso")
}

func (e *Env) UpdateBrandHandler(w http.ResponseWriter, r *http.Request) {
	var err error

	var providedBrand UpdateBrandDTO
	ReadJSON(w, r, &providedBrand)
	err = UpdateBrandValidator(&providedBrand)
	if err != nil {
		WriteError(w, http.StatusUnprocessableEntity, "Erro de validação: " + err.Error())
		return
	}

	if providedBrand.Name == "" {
		WriteJSON(w, http.StatusNotModified, nil, "Nada mudou")
		return
	}

	brandId := r.PathValue("id")
	if brandId == "" {
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
		UPDATE brands 
		SET name = $1, updated_at = NOW()
		WHERE id = $2
		`,
		&providedBrand.Name,
		brandId,
	)
	if err != nil {
		transaction.Rollback(r.Context())
		WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	transaction.Commit(r.Context())
	WriteJSON(w, http.StatusOK, nil, "Marca atualizada com sucesso")
}

func (e *Env) ListBrandsHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := e.DB.Query(
		r.Context(),
		`
		SELECT id, name, household_id, created_at, updated_at
		FROM brands
		WHERE deleted_at IS NULL
		ORDER BY created_at, updated_at DESC
		`,
	)
	if err != nil {
		fmt.Printf("Database error: %v\n", err)
		WriteError(w, http.StatusInternalServerError, "Erro interno no banco de dados")
		return
	}
	defer rows.Close()

	brands := []models.Brand{}

	for rows.Next() {
		brand := models.Brand{}
		rows.Scan(
			&brand.ID,
			&brand.Name,
			&brand.HouseholdID,
			&brand.CreatedAt,
			&brand.UpdatedAt,
		)
		brands = append(brands, brand)
	}

	WriteJSON(w, http.StatusOK, brands, "Marcas listadas com sucesso")
}

func (e *Env) ListOneBrandHandler(w http.ResponseWriter, r *http.Request) {
	var foundBrand models.Brand
	
	brandId := r.PathValue("id")
	if brandId == "" {
		WriteError(w, http.StatusBadRequest, "Insira um ID")
		return
	}

	err := e.DB.QueryRow(
		r.Context(),
		`
		SELECT id, name, household_id, created_at, updated_at
		FROM brands
		WHERE id = $1
		AND deleted_at IS NULL
		`,
		brandId,
	).Scan(
		&foundBrand.ID,
		&foundBrand.Name,
		&foundBrand.HouseholdID,
		&foundBrand.CreatedAt,
		&foundBrand.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			WriteError(w, http.StatusNotFound, "Marca não encontrada")
			return
		}

		fmt.Printf("Database error: %v\n", err)
		WriteError(w, http.StatusInternalServerError, "Erro interno no banco de dados")
		return
	}

	WriteJSON(w, http.StatusOK, foundBrand, "Marca listada com sucesso")
}

func (e *Env) DeleteBrandHandler(w http.ResponseWriter, r *http.Request) {
	brandId := r.PathValue("id")
	if brandId == "" {
		WriteError(w, http.StatusBadRequest, "Insira um ID")
		return
	}

	err := e.DB.QueryRow(
		r.Context(),
		`
		UPDATE brands
		SET deleted_at = NOW()
		WHERE id = $1
		AND deleted_at IS NULL
		RETURNING id
		`,
		brandId,
	).Scan(nil)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			WriteError(w, http.StatusNotFound, "Marca não encontrada")
			return
		}

		fmt.Printf("Database error: %v\n", err)
		WriteError(w, http.StatusInternalServerError, "Erro interno no banco de dados")
		return
	}

	WriteJSON(w, http.StatusOK, nil, "Marca deletada com sucesso")
}

