package main

import (
	"log"
	"net/http"
	"time"

	"github.com/ArdiSasongko/SocialNetwork/cmd/api/v1/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type application struct {
	config  Config
	handler handlers.Handler
}

type Config struct {
	addr string
	db   dbConfig
	auth authConfig
}

type dbConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

type authConfig struct {
	secret string
	exp    time.Duration
	iss    string
}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/v1", func(r chi.Router) {
		r.Get("/health", app.handler.Health.Get)

		// auth handler
		r.Route("/authentication", func(r chi.Router) {
			r.Post("/user", app.handler.Users.RegisterUser)
			r.Post("/login", app.handler.Users.LoginUser)
		})
	})

	return r
}

func (app *application) run(mux http.Handler) error {
	server := &http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 30,
		IdleTimeout:  time.Minute,
	}

	log.Printf("server has running, address %v", app.config.addr)

	return server.ListenAndServe()
}
