package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/goredisstore"
	"github.com/alexedwards/scs/v2"
	"github.com/dev-joaovitor/despensa-digital/config"
	"github.com/dev-joaovitor/despensa-digital/database"
	"github.com/dev-joaovitor/despensa-digital/handlers"
	"github.com/dev-joaovitor/despensa-digital/mail"
)

func main() {
	// context
	ctx, cancel := context.WithTimeout(
		context.Background(),
		15 * time.Second,
	)
	defer cancel()

	// env vars
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Initialization error: %v", err)
	}

	// persistence
	dbPool, err := database.ConnectPostgres(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Postgres database isolation layer failed: %v", err)
	}
	defer dbPool.Close()

	redisClient, err := database.ConnectRedis(ctx, cfg.RedisURL)
	if err != nil {
		log.Fatalf("Redis execution state engine connection failed: %v", err)
	}
	defer redisClient.Close()

	// session manager
	sessionManager := scs.New()
	sessionManager.Lifetime = 24 * time.Hour
	sessionManager.Cookie.Persist = true
	sessionManager.Cookie.Secure = true
	sessionManager.Cookie.SameSite = http.SameSiteStrictMode
	sessionManager.Store = goredisstore.New(redisClient)

	// mailer
	mailer := mail.NewEmailService(
		cfg.MailSMTPHost,
		cfg.MailSMTPPort,
		cfg.MailSMTPUser,
		cfg.MailSMTPPassword,
		"noreply@intellistock.com",
	)

	// migrations
	if err := database.RunMigrations(cfg.DatabaseURL); err != nil {
		log.Fatalf("Database schema migration halted system boot: %v", err)
	}

	appEnv := &handlers.Env{
		DB: dbPool,
		Cache: redisClient,
		SessionState: sessionManager,
		MailService: mailer,
	}

	server := &http.Server{
		Addr: ":" + cfg.Port,
		Handler: appEnv.LoadRoutes(),
		ReadTimeout: 5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout: 120 * time.Second,
	}

	log.Printf("Server initializing execution engine boundary on %s", cfg.Port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Fatal network lifecycle failure: %v", err)
	}
}

