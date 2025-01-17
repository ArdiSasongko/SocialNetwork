package main

import (
	"log"
	"net/http"
	"time"

	"github.com/ArdiSasongko/SocialNetwork/cmd/api/v1/handlers"
	"github.com/ArdiSasongko/SocialNetwork/cmd/api/v1/middlewares"
	"github.com/ArdiSasongko/SocialNetwork/internal/env"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
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
	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{env.GetString("ALLOWED_ORIGIN", "http://*")},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
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

			// get post
			r.Group(func(r chi.Router) {
				r.Use(app.middleware.PostCTXMiddleware)
				r.Get("/{postID}", app.handler.Post.GetPostByUser)
			})

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

		// user handler
		r.Route("/users", func(r chi.Router) {
			r.Use(app.middleware.AuthMiddleware)
			r.Route("/{userID}", func(r chi.Router) {
				r.Use(app.middleware.UserProfileCTXMiddleware)
				r.Get("/", app.handler.Users.GetUserProfile)

				r.Post("/follow", app.handler.Users.FollowUser)
				r.Delete("/unfollow", app.handler.Users.UnfollowUser)
			})
		})

		// feed handler
		r.Route("/feeds", func(r chi.Router) {
			r.Use(app.middleware.AuthMiddleware)
			r.Get("/", app.handler.Feed.GetFeeds)

			r.Route("/{postID}", func(r chi.Router) {
				r.Use(app.middleware.PostCTXMiddleware)
				r.Get("/", app.handler.Feed.GetFeed)
				r.Post("/comment", app.handler.Feed.CreateComment)
				r.Put("/like", app.handler.Feed.LikedFeed)
				r.Put("/dislike", app.handler.Feed.DisikedFeed)
			})
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
