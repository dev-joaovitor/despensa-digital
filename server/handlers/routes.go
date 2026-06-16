package handlers

import (
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/dev-joaovitor/despensa-digital/mail"
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
	MailService  *mail.MailService
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

		// auth only routes
		r.Group(func (auth chi.Router) {
			auth.Use(e.AuthRequiredMiddleware)
			auth.Route("/households", e.HouseholdsHandler)
			auth.Route("/establishments", e.EstablishmentsHandler)
			auth.Route("/brands", e.BrandsHandler)
			auth.Route("/categories", e.CategoriesHandler)
			auth.Route("/price-observations", e.PriceObservationsHandler)
			auth.Route("/products", e.ProductsHandler)
		})
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

