package handlers

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/dev-joaovitor/despensa-digital/config"
	"github.com/dev-joaovitor/despensa-digital/models"
	"github.com/go-chi/chi/v5"
	"golang.org/x/crypto/bcrypt"
)

func (e *Env) AuthHandler(r chi.Router) {
	r.Post("/login", e.LoginHandler)
	r.Post("/logout", e.LogoutHandler)
}

type LoginDTO struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

func ValidateLogin(login *LoginDTO) error {
	validationErrors := []string{}

	password := login.Password
	if password == "" {
		validationErrors = append(validationErrors, "Senha é obrigatória.")
	}

	email := strings.TrimSpace(login.Email)
	if email == "" {
		validationErrors = append(validationErrors, "Email é obrigatório.")
	}

	if len(email) > 254 {
		validationErrors = append(validationErrors, "Email é muito grande.")
	}

	if len(validationErrors) == 0 {
		return nil
	}

	return errors.New(strings.Join(validationErrors, " "))
}

func (e *Env) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var err error

	var providedLogin LoginDTO
	ReadJSON(w, r, &providedLogin)
	err = ValidateLogin(&providedLogin)

	if err != nil {
		WriteError(w, http.StatusUnprocessableEntity, "Erro de validação: " + err.Error())
		return
	}

	var foundUser models.User

	err = e.DB.QueryRow(
		r.Context(),
		`
			SELECT id, password
			FROM users
			WHERE email = $1
			LIMIT 1
		`,
		providedLogin.Email,
	).Scan(&foundUser.ID, &foundUser.Password)

	if !e.ComparePassword(providedLogin.Password, foundUser.Password) {
		WriteError(w, http.StatusUnauthorized, "Senha ou email inválidos")
		return
	}

	e.SessionState.Put(r.Context(), "userID", &foundUser.ID)
	WriteJSON(w, http.StatusOK, nil, "Bem vindo de volta")
}

func (e *Env) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	e.SessionState.RenewToken(r.Context())
	e.SessionState.Remove(r.Context(), "userID")
	WriteJSON(w, http.StatusOK, nil, "Até breve")
}

func (e *Env) AuthRequiredMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !e.SessionState.Exists(r.Context(), "userID") {
			WriteError(w, http.StatusUnauthorized, "Token inválido. Faça login e tente novamente")
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (e *Env) HashPassword(password string) (string, error) {
	cfg, err := config.Load()
	if err != nil {
		return "", err
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), cfg.BcryptRounds)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func (e *Env) ComparePassword(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (e *Env) GetSessionUserId(c context.Context) int64 {
	return e.SessionState.GetInt64(c, "userID")
}

