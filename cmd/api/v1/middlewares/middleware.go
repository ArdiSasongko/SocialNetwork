package middlewares

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/ArdiSasongko/SocialNetwork/internal/auth"
	"github.com/ArdiSasongko/SocialNetwork/internal/storage/postgresql"
	"github.com/ArdiSasongko/SocialNetwork/utils"
	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
)

type userkey string
type postKey string

const UserCtx userkey = "user"
const PostCtx postKey = "post"

type Middleware struct {
	json    utils.JsonUtils
	errror  utils.ErrorUtils
	auth    auth.Authenticator
	storage postgresql.Storage
}

func NewMiddleware(db *sql.DB, auth auth.Authenticator) Middleware {
	json := utils.NewJsonUtils()
	error := utils.NewErrorUtils()
	storage := postgresql.NewStorage(db)
	return Middleware{
		json:    json,
		errror:  error,
		auth:    auth,
		storage: storage,
	}
}

func (m *Middleware) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			m.errror.UnauthorizedError(w, r, fmt.Errorf("missing authorizaton header"))
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			m.errror.UnauthorizedError(w, r, fmt.Errorf("header authorization is malformed"))
			return
		}

		token := parts[1]
		jwtToken, err := m.auth.ValidateToken(token)
		if err != nil {
			m.errror.UnauthorizedError(w, r, err)
			return
		}

		claims, _ := jwtToken.Claims.(jwt.MapClaims)

		userID, err := strconv.ParseInt(fmt.Sprintf("%.f", claims["sub"]), 10, 64)
		if err != nil {
			m.errror.UnauthorizedError(w, r, err)
			return
		}

		ctx := r.Context()
		// get user from database
		user, err := m.storage.Users.GetByID(ctx, userID)
		if err != nil {
			m.errror.UnauthorizedError(w, r, err)
			return
		}

		ctx = context.WithValue(ctx, UserCtx, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *Middleware) PostCTXMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idParam := chi.URLParam(r, "postID")
		postID, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			m.errror.InternalServerError(w, r, err)
			return
		}

		ctx := r.Context()
		post, err := m.storage.Posts.GetPostByID(ctx, postID)
		if err != nil {
			switch {
			case errors.Is(err, postgresql.ErrNotFound):
				m.errror.NotFoundError(w, r, err)
			default:
				m.errror.InternalServerError(w, r, err)
			}
			return
		}

		log.Println(post.Title)
		ctx = context.WithValue(ctx, PostCtx, post)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
