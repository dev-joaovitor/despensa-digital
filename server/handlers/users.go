package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/dev-joaovitor/despensa-digital/models"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (e *Env) UsersHandler(r chi.Router) {
	r.Post("/", e.CreateUserHandler)

	r.Group(func (auth chi.Router) {
		auth.Use(e.AuthRequiredMiddleware)
		auth.Patch("/", e.UpdateUserHandler)
		auth.Get("/", e.ListUsersHandler)
	})
}

func (e *Env) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var err error

	var providedUser CreateUserDTO
	ReadJSON(w, r, &providedUser)
	err = CreateUserValidator(&providedUser)
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

	var foundHousehold models.Household

	if providedUser.InvitationCode != nil {
		err = transaction.QueryRow(
			r.Context(),
			`
			SELECT id, creator_id FROM households
			WHERE invitation_code = $1
			LIMIT 1
			`,
			&providedUser.InvitationCode,
		).Scan(
			&foundHousehold.ID,
			&foundHousehold.CreatorID,
		)

		if errors.Is(err, sql.ErrNoRows) {
			WriteError(w, http.StatusUnprocessableEntity, "Residência não encontrada com o código inserido")
			return
		}
	} else if providedUser.HouseholdName != nil {
		householdQuery := `
			INSERT INTO households (creator_id, name, invitation_code)
			VALUES (NULL, $1, $2)
			RETURNING id
		`
		invitationCode, err := uuid.NewRandom()
		if err != nil {
			WriteError(
				w,
				http.StatusInternalServerError,
				"Falha ao gerar código de convite para nova residência: " + err.Error(),
			)
			return
		}

		err = transaction.QueryRow(
			r.Context(),
			householdQuery,
			&providedUser.HouseholdName,
			invitationCode.String(),
		).Scan(&foundHousehold.ID)

		if err != nil {
			WriteError(
				w,
				http.StatusBadRequest,
				"Falha ao criar uma residência: " + err.Error(),
			)
			return
		}
	} else {
		WriteError(w, http.StatusUnprocessableEntity, "É obrigatório ter uma residência")
		return
	}

	userQuery := `
		INSERT INTO users (full_name, email, password, household_id)
		VALUES ($1, $2, $3, $4)
		RETURNING id, full_name, email, household_id, created_at, updated_at
	`

	var createdUser models.User
	hashedPassword, err := e.HashPassword(providedUser.Password)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Criptografia de senhas falhou")
		return 
	}

	err = transaction.QueryRow(
		r.Context(),
		userQuery,
		&providedUser.FullName,
		&providedUser.Email,
		hashedPassword,
		&foundHousehold.ID,
	).Scan(
		&createdUser.ID,
		&createdUser.FullName,
		&createdUser.Email,
		&createdUser.HouseholdID,
		&createdUser.CreatedAt,
		&createdUser.UpdatedAt,
	)

	if err != nil {
		WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	if foundHousehold.CreatorID == nil {
		updateHouseholdQuery := `
			UPDATE households
			SET creator_id = $1
			WHERE id = $2
		`
		_, err := transaction.Exec(
			r.Context(),
			updateHouseholdQuery,
			&createdUser.ID,
			&foundHousehold.ID,
		)
		if err != nil {
			WriteError(w, http.StatusBadRequest, err.Error())
			return
		}
	}

	err = transaction.Commit(r.Context())
	if err != nil {
		WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	e.SessionState.Put(r.Context(), "userID", &createdUser.ID)
	WriteJSON(w, http.StatusCreated, &createdUser, "Conta criada com sucesso")
}

func (e *Env) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	var err error

	var providedUser UpdateUserDTO
	ReadJSON(w, r, &providedUser)
	err = UpdateUserValidator(&providedUser)

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

	userId := e.GetSessionUserId(r.Context())

	var foundUser models.User

	err = transaction.QueryRow(
		r.Context(),
		`
		SELECT id, password
		FROM users
		WHERE id = $1
		`,
		userId,
	).Scan(
		&foundUser.ID,
		&foundUser.Password,
	)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "Não foi possível atualizar o usuário")
		return
	}

	if providedUser.FullName != "" {
		transaction.QueryRow(
			r.Context(),
			`
			UPDATE users
			SET full_name = $1
			WHERE id = $2
			`,
			providedUser.FullName,
			userId,
		)
	}

	if providedUser.Email != "" {
		transaction.QueryRow(
			r.Context(),
			`
			UPDATE users
			SET email = $1
			WHERE id = $2
			`,
			providedUser.Email,
			userId,
		)
	}

	if providedUser.NewPassword != "" {
		if providedUser.Code != "" {
			_, err = transaction.Exec(
				r.Context(),
				`
				SELECT id
				FROM users
				WHERE id = $1
				AND verification_code = $2
				AND expires_at > NOW()
				LIMIT 1
				`,
				userId,
				providedUser.Code,
			)
			if err != nil {
				WriteError(w, http.StatusUnprocessableEntity, "O código expirou ou está errado")
				return 
			}
		} else {
			if !e.ComparePassword(providedUser.OldPassword, foundUser.Password) {
				WriteError(w, http.StatusUnprocessableEntity, "Senha antiga está incorreta")
				return 
			}
		}

		hashedPassword, err := e.HashPassword(providedUser.NewPassword)
		if err != nil {
			fmt.Printf("Hash error: %s\n", err)
			WriteError(w, http.StatusInternalServerError, "Criptografia de senhas falhou")
			return 
		}

		transaction.QueryRow(
			r.Context(),
			`
			UPDATE users
			SET password = $1, expires_at = NOW()
			WHERE id = $2
			`,
			hashedPassword,
			userId,
		)
	}

	transaction.QueryRow(
		r.Context(),
		`
		UPDATE users
		SET updated_at = NOW()
		WHERE id = $1
		`,
		userId,
	)
	if err != nil {
		WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	transaction.Commit(r.Context())
	e.SessionState.RenewToken(r.Context())
	WriteJSON(w, http.StatusOK, nil, "Conta atualizada com sucesso")
}

func (e *Env) ListUsersHandler(w http.ResponseWriter, r *http.Request) {
	var err error

	rows, err := e.DB.Query(
		r.Context(),
		`
		SELECT id, household_id, full_name, created_at, updated_at
		FROM users
		WHERE deleted_at IS NULL
		ORDER BY created_at, updated_at DESC
		`,
	)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Erro interno no banco de dados")
		return
	}
	defer rows.Close()

	users := []models.User{}

	for rows.Next() {
		user := models.User{}
		rows.Scan(
			&user.ID,
			&user.HouseholdID,
			&user.FullName,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		users = append(users, user)
	}

	WriteJSON(w, http.StatusOK, users, "Usuários listados com sucesso")
}

