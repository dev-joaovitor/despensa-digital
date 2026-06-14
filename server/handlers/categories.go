package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/dev-joaovitor/despensa-digital/models"
	"github.com/go-chi/chi/v5"
)

func (e *Env) CategoriesHandler(r chi.Router) {
	r.Get("/{id}", e.ListOneCategoryHandler)
	r.Get("/", e.ListCategorysHandler)
	r.Post("/", e.CreateCategoryHandler)
	r.Patch("/{id}", e.UpdateCategoryHandler)
	r.Delete("/{id}", e.DeleteCategoryHandler)
}

func (e *Env) CreateCategoryHandler(w http.ResponseWriter, r *http.Request) {
	var err error

	var providedCategory CreateCategoryDTO
	ReadJSON(w, r, &providedCategory)
	err = CreateCategoryValidator(&providedCategory)

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

	var createdCategory models.Category
	userHousehold, err := e.GetSessionUserHousehold(r.Context())
	if err != nil {
		transaction.Rollback(r.Context())
		WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = e.DB.QueryRow(
		r.Context(),
		`
		INSERT INTO categories (household_id, name)
		VALUES ($1, $2)
		RETURNING id, household_id, name, created_at, updated_at
		`,
		&userHousehold.ID,
		&providedCategory.Name,
	).Scan(
		&createdCategory.ID,
		&createdCategory.HouseholdID,
		&createdCategory.Name,
		&createdCategory.CreatedAt,
		&createdCategory.UpdatedAt,
	)

	if err != nil {
		transaction.Rollback(r.Context())
		WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	transaction.Commit(r.Context())
	WriteJSON(w, http.StatusCreated, &createdCategory, "Categoria criada com sucesso")
}

func (e *Env) UpdateCategoryHandler(w http.ResponseWriter, r *http.Request) {
	var err error

	var providedCategory UpdateCategoryDTO
	ReadJSON(w, r, &providedCategory)
	err = UpdateCategoryValidator(&providedCategory)
	if err != nil {
		WriteError(w, http.StatusUnprocessableEntity, "Erro de validação: " + err.Error())
		return
	}

	if providedCategory.Name == "" {
		WriteJSON(w, http.StatusNotModified, nil, "Nada mudou")
		return
	}

	categoryId := r.PathValue("id")
	if categoryId == "" {
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
		UPDATE categories 
		SET name = $1, updated_at = NOW()
		WHERE id = $2
		`,
		&providedCategory.Name,
		categoryId,
	)
	if err != nil {
		transaction.Rollback(r.Context())
		WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	transaction.Commit(r.Context())
	WriteJSON(w, http.StatusOK, nil, "Categoria atualizada com sucesso")
}

func (e *Env) ListCategorysHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := e.DB.Query(
		r.Context(),
		`
		SELECT id, name, household_id, created_at, updated_at
		FROM categories
		WHERE deleted_at IS NULL
		ORDER BY created_at, updated_at DESC
		`,
	)
	if err != nil {
		fmt.Printf("Database error: %v\n", err)
		WriteError(w, http.StatusInternalServerError, "Erro interno no banco de dados")
		return
	}

	categories := []models.Category{}

	for rows.Next() == true {
		category := models.Category{}
		rows.Scan(
			&category.ID,
			&category.Name,
			&category.HouseholdID,
			&category.CreatedAt,
			&category.UpdatedAt,
		)
		categories = append(categories, category)
	}

	WriteJSON(w, http.StatusOK, categories, "Categorias listadas com sucesso")
}

func (e *Env) ListOneCategoryHandler(w http.ResponseWriter, r *http.Request) {
	var foundCategory models.Category
	
	categoryId := r.PathValue("id")
	if categoryId == "" {
		WriteError(w, http.StatusBadRequest, "Insira um ID")
		return
	}

	err := e.DB.QueryRow(
		r.Context(),
		`
		SELECT id, name, household_id, created_at, updated_at
		FROM categories
		WHERE id = $1
		AND deleted_at IS NULL
		`,
		categoryId,
	).Scan(
		&foundCategory.ID,
		&foundCategory.Name,
		&foundCategory.HouseholdID,
		&foundCategory.CreatedAt,
		&foundCategory.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			WriteError(w, http.StatusNotFound, "Categoria não encontrada")
			return
		}

		fmt.Printf("Database error: %v\n", err)
		WriteError(w, http.StatusInternalServerError, "Erro interno no banco de dados")
		return
	}

	WriteJSON(w, http.StatusOK, foundCategory, "Categoria listada com sucesso")
}

func (e *Env) DeleteCategoryHandler(w http.ResponseWriter, r *http.Request) {
	categoryId := r.PathValue("id")
	if categoryId == "" {
		WriteError(w, http.StatusBadRequest, "Insira um ID")
		return
	}

	err := e.DB.QueryRow(
		r.Context(),
		`
		UPDATE categories
		SET deleted_at = NOW()
		WHERE id = $1
		AND deleted_at IS NULL
		RETURNING id
		`,
		categoryId,
	).Scan(nil)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			WriteError(w, http.StatusNotFound, "Categoria não encontrada")
			return
		}

		fmt.Printf("Database error: %v\n", err)
		WriteError(w, http.StatusInternalServerError, "Erro interno no banco de dados")
		return
	}

	WriteJSON(w, http.StatusOK, nil, "Categoria deletada com sucesso")
}

