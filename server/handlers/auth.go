package handlers

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
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
	r.Post("/send-recovery-code", e.SendRecoveryCodeHandler)
	r.Post("/verify-recovery-code", e.VerifyRecoveryCodeHandler)

	r.Group(func (auth chi.Router) {
		auth.Use(e.AuthRequiredMiddleware)
		auth.Post("/change-password", e.ChangePasswordHandler)
	})
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
		AND deleted_at IS NULL
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

type SendRecoveryCodeDTO struct {
	Email string `json:"email"`
}

func (e *Env) SendRecoveryCodeHandler(w http.ResponseWriter, r *http.Request) {
	var err error

	var providedPayload SendRecoveryCodeDTO
	ReadJSON(w, r, &providedPayload)

	email := strings.TrimSpace(providedPayload.Email)
	if email == "" {
		WriteError(w, http.StatusUnprocessableEntity, "Email é obrigatório")
		return
	}

	err = e.DB.QueryRow(
		r.Context(),
		`
		SELECT id
		FROM users
		WHERE email = $1
		AND deleted_at IS NULL
		LIMIT 1
		`,
		email,
	).Scan(nil)
	if err != nil {
		fmt.Printf("User not found with email %s: %s\n", email, err)
	}

	if err == nil {
		generatedCode, err := Generate4DigitString()
		if err != nil {
			fmt.Printf("Failed to generate code: %s\n", err)
		}

		if err == nil {
			err = e.MailService.SendPasswordReset(email, generatedCode)
			if err == nil {
				_, err = e.DB.Query(
					r.Context(),
					`
					UPDATE users
					SET verification_code = $1,
					expires_at = NOW() + INTERVAL '5 minutes'
					WHERE email = $2
					`,
					generatedCode,
					email,
				)
			} else {
				fmt.Printf("Email service is not working: %s; To: %s / Recovery Code: %s\n", err, email, generatedCode)
			}
		}
	}

	WriteJSON(w, http.StatusOK, nil, "Se a conta existir, enviaremos um código de recuperação para ele")
}

type VerifyRecoveryCodeDTO struct {
	Email string `json:"email"`
	Code string `json:"code"`
}

func (e *Env) VerifyRecoveryCodeHandler(w http.ResponseWriter, r *http.Request) {
	var err error

	var providedPayload VerifyRecoveryCodeDTO
	ReadJSON(w, r, &providedPayload)

	email := strings.TrimSpace(providedPayload.Email)
	if email == "" {
		WriteError(w, http.StatusUnprocessableEntity, "Email é obrigatório")
		return
	}

	code := strings.TrimSpace(providedPayload.Code)
	if code == "" {
		WriteError(w, http.StatusUnprocessableEntity, "O código é obrigatório")
		return
	}

	if len(code) != 4 {
		WriteError(w, http.StatusUnprocessableEntity, "O código inserido deve ter 4 dígitos")
		return
	}

	var foundUser models.User

	err = e.DB.QueryRow(
		r.Context(),
		`
		SELECT id
		FROM users
		WHERE email = $1
		AND verification_code = $2
		AND expires_at > NOW()
		AND deleted_at IS NULL
		`,
		email,
		code,
	).Scan(&foundUser.ID)
	if err != nil {
		WriteError(w, http.StatusUnauthorized, "O código está errado ou já expirou")
		return
	}

	e.SessionState.Put(r.Context(), "userID", &foundUser.ID)
	WriteJSON(w, http.StatusOK, nil, "Sucesso ao validar o código")
}

type ChangePasswordDTO struct {
	NewPassword string `json:"new_password"`
	NewPasswordConfirmation string `json:"new_password_confirmation"`
}

func (e *Env) ChangePasswordHandler(w http.ResponseWriter, r *http.Request) {
	var err error

	var providedPayload ChangePasswordDTO
	ReadJSON(w, r, &providedPayload)

	newPassword := providedPayload.NewPassword
	if newPassword == "" {
		WriteError(w, http.StatusUnprocessableEntity, "Senha é obrigatória")
		return
	}

	if len(newPassword) < 6 {
		WriteError(w, http.StatusUnprocessableEntity, "Senha deve ter pelo menos 6 caracteres.")
		return
	}

	passwordConfirmation := providedPayload.NewPasswordConfirmation
	if newPassword != passwordConfirmation {
		WriteError(w, http.StatusUnprocessableEntity, "As senhas devem ser iguais")
		return
	}

	userId := e.GetSessionUserId(r.Context())
	newPassword, err = e.HashPassword(newPassword)

	err = e.DB.QueryRow(
		r.Context(),
		`
		UPDATE users
		SET password = $1
		WHERE id = $2, updated_at = NOW()
		`,
		newPassword,
		userId,
	).Scan(nil)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			fmt.Printf("user id: %d / new password: %s / error: %s\n", userId, newPassword, err)
			WriteError(w, http.StatusBadRequest, "Não foi possível trocar a senha, tente novamente mais tarde")
			return
		}
	}

	e.SessionState.RenewToken(r.Context())
	WriteJSON(w, http.StatusOK, nil, "A senha foi trocada com sucesso")
}

func (e *Env) AuthRequiredMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !e.SessionState.Exists(r.Context(), "userID") {
			WriteError(w, http.StatusUnauthorized, "Você deve estar logado para fazer isso")
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

