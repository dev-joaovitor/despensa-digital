package handlers

import (
	"database/sql"
	"errors"
	"net/http"
	"strings"

	"github.com/dev-joaovitor/despensa-digital/models"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (e *Env) UsersHandler(r chi.Router) {
	r.Post("/", e.CreateUser)

	r.Group(func (auth chi.Router) {
		auth.Use(e.AuthRequiredMiddleware)
		auth.Get("/", e.ListUsers)
	})
}

type UserDTO struct {
	FullName string `json:"full_name"`
	Email string `json:"email"`
	Password string `json:"password"`

	InvitationCode *string `json:"invitation_code"`
	HouseholdName *string `json:"household_name"`
}

func ValidateCreateUser(user *UserDTO) error {
	validationErrors := []string{}

	fullname := strings.TrimSpace(user.FullName)
	if fullname == "" {
		validationErrors = append(validationErrors, "Nome completo é obrigatório.")
	}

	if len(fullname) < 4 || len(fullname) > 100 {
		validationErrors = append(validationErrors, "Nome completo deve ter entre 4 a 100 caracteres.")
	}

	password := user.Password
	if password == "" {
		validationErrors = append(validationErrors, "Senha é obrigatória.")
	}

	if len(password) < 6 {
		validationErrors = append(validationErrors, "Senha deve ter pelo menos 6 caracteres.")
	}

	email := strings.TrimSpace(user.Email)
	if email == "" {
		validationErrors = append(validationErrors, "Email é obrigatório.")
	}

	if len(email) > 254 {
		validationErrors = append(validationErrors, "Email é muito grande.")
	}

	if user.HouseholdName != nil && user.InvitationCode == nil {
		if strings.TrimSpace(*user.HouseholdName) == "" {
			validationErrors = append(validationErrors, "É necessário criar uma residência quando não se tem um convite.")
		}
	}

	if user.InvitationCode != nil && user.HouseholdName == nil {
		if strings.TrimSpace(*user.InvitationCode) == "" {
			validationErrors = append(validationErrors, "É necessário ter um convite se não for criar uma residência.")
		}
	}

	if len(validationErrors) == 0 {
		return nil
	}

	return errors.New(strings.Join(validationErrors, " "))
}

func (e *Env) CreateUser(w http.ResponseWriter, r *http.Request) {
	var err error

	var providedUser UserDTO
	ReadJSON(w, r, &providedUser)
	err = ValidateCreateUser(&providedUser)

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

	var foundHousehold models.Household

	if providedUser.InvitationCode != nil {
		householdQuery := `
			SELECT id FROM households
			WHERE invitation_code = $1
			LIMIT 1
		`
		err = e.DB.QueryRow(
			r.Context(),
			householdQuery,
			&providedUser.InvitationCode,
		).Scan(&foundHousehold.ID)

		if errors.Is(err, sql.ErrNoRows) {
			transaction.Rollback(r.Context())
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
			transaction.Rollback(r.Context())
			WriteError(
				w,
				http.StatusInternalServerError,
				"Falha ao gerar código de convite para nova residência: " + err.Error(),
			)
			return
		}

		err = e.DB.QueryRow(
			r.Context(),
			householdQuery,
			&providedUser.HouseholdName,
			invitationCode.String(),
		).Scan(&foundHousehold.ID)

		if err != nil {
			transaction.Rollback(r.Context())
			WriteError(
				w,
				http.StatusBadRequest,
				"Falha ao criar uma residência: " + err.Error(),
			)
			return
		}
	} else {
		transaction.Rollback(r.Context())
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
		transaction.Rollback(r.Context())
		WriteError(w, http.StatusInternalServerError, "Criptografia de senhas falhou")
		return 
	}

	err = e.DB.QueryRow(
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
		transaction.Rollback(r.Context())
		WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	if foundHousehold.CreatorID == nil {
		updateHouseholdQuery := `
			UPDATE households
			SET creator_id = $1
			WHERE id = $2
		`
		_, err = e.DB.Query(
			r.Context(),
			updateHouseholdQuery,
			&createdUser.ID,
			&foundHousehold.ID,
		)

		if err != nil {
			transaction.Rollback(r.Context())
			WriteError(w, http.StatusBadRequest, err.Error())
			return
		}
	}

	transaction.Commit(r.Context())
	WriteJSON(w, http.StatusCreated, &createdUser, "Conta criada com sucesso")
}

func (e *Env) ListUsers(w http.ResponseWriter, r *http.Request) {
	var err error

	query := `
		SELECT id, household_id, full_name, created_at, updated_at
		FROM users
		WHERE deleted_at IS NULL
		ORDER BY created_at, updated_at DESC
	`

	rows, err := e.DB.Query(
		r.Context(),
		query,
	)

	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Erro interno no banco de dados")
		return
	}

	users := []models.User{}

	for rows.Next() == true {
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
