package handlers

import (
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"

)


// Env state sharing layer across handlers via receiver functions.
type Env struct {
	DB           *pgxpool.Pool
	Cache        *redis.Client
	SessionState *scs.SessionManager
}

func (e *Env) LoadRoutes() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(e.SessionState.LoadAndSave)

	r.Route("/api/v1", func (r chi.Router) {
		r.Get("/health", e.HealthCheckHandler)
		r.Route("/auth", e.AuthHandler)
		r.Route("/users", e.UsersHandler)
	})

	return r
}

func (e *Env) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	WriteJSON(
		w,
		http.StatusOK,
		map[string]string{ "status": "healthy" },
		"Tudo certo",
	)
}

