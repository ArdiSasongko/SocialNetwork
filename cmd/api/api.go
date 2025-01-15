package main

import (
	"log"
	"net/http"
	"time"

	"github.com/ArdiSasongko/SocialNetwork/cmd/api/v1/handlers"
	"github.com/ArdiSasongko/SocialNetwork/cmd/api/v1/middlewares"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type application struct {
	config     Config
	handler    handlers.Handler
	middleware middlewares.Middleware
}

type Config struct {
	addr       string
	db         dbConfig
	auth       authConfig
	cloudinary cldConfig
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

type cldConfig struct {
	url    string
	folder string
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
			r.Post("/register", app.handler.Auth.RegisterUser)
			r.Post("/login", app.handler.Auth.LoginUser)
		})

		// profile handler
		r.Route("/profile", func(r chi.Router) {
			r.Use(app.middleware.AuthMiddleware)
			r.Get("/", app.handler.Users.GetProfile)

			r.Patch("/", app.handler.Users.UpdateUser)
			r.Put("/image", app.handler.Users.UpdateImages)
		})

		// post handler
		r.Route("/posts", func(r chi.Router) {
			r.Use(app.middleware.AuthMiddleware)
			r.Post("/", app.handler.Post.CreatePost)

			r.Route("/{postID}", func(r chi.Router) {
				r.Use(app.middleware.PostCTXMiddleware)
				r.Get("/", app.handler.Post.GetPostByID)

				// middleware authorization
				r.Patch("/", app.handler.Post.CheckOwnerPost("moderator", app.handler.Post.UpdatePost))
				r.Delete("/", app.handler.Post.CheckOwnerPost("admin", app.handler.Post.DeletePost))
			})
		})

		// feed handler

		// user handler
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
