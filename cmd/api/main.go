package main

import (
	"log"
	"time"

	"github.com/ArdiSasongko/SocialNetwork/cmd/api/v1/handlers"
	"github.com/ArdiSasongko/SocialNetwork/cmd/api/v1/middlewares"
	"github.com/ArdiSasongko/SocialNetwork/internal/auth"
	"github.com/ArdiSasongko/SocialNetwork/internal/db"
	"github.com/ArdiSasongko/SocialNetwork/internal/env"
	"github.com/ArdiSasongko/SocialNetwork/internal/storage/cldnary"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err.Error())
	}

	cfg := Config{
		addr: env.GetString("SERVER_ADDR", ":3000"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgres://root:mypassword@localhost:5432/socialnetwork?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 15),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 5),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "10m"),
		},
		auth: authConfig{
			secret: env.GetString("JWT_SECRET", "mostsecretvalue"),
			iss:    env.GetString("JWT_ISS", "SocialNetwork"),
			exp:    time.Hour * 24 * 3,
		},
		cloudinary: cldConfig{
			url:    env.GetString("CLOUDINARY_URL", ""),
			folder: env.GetString("CLOUDINARY_FOLDER", ""),
		},
	}

	// connection to database
	conn, err := db.New(
		cfg.db.addr,
		cfg.db.maxOpenConns,
		cfg.db.maxIdleConns,
		cfg.db.maxIdleTime,
	)

	if err != nil {
		log.Fatal(err.Error())
	}

	auth := auth.NewJWT(
		cfg.auth.secret,
		cfg.auth.iss,
		cfg.auth.iss,
	)

	cld, err := cldnary.NewCloudinary(
		cfg.cloudinary.url,
		cfg.cloudinary.folder,
	)

	if err != nil {
		log.Fatal(err.Error())
	}

	handler := handlers.NewHandler(conn, auth, *cld)
	middleware := middlewares.NewMiddleware(conn, auth)

	app := application{
		config:     cfg,
		handler:    handler,
		middleware: middleware,
	}

	mux := app.mount()
	log.Fatal(app.run(mux))
}
