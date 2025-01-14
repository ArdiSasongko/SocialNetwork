package middlewares

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/ArdiSasongko/SocialNetwork/internal/auth"
	"github.com/ArdiSasongko/SocialNetwork/internal/service"
	"github.com/ArdiSasongko/SocialNetwork/utils"
	"github.com/golang-jwt/jwt/v5"
)

type userkey string

const UserCtx userkey = "user"

type Middleware struct {
	json    utils.JsonUtils
	errror  utils.ErrorUtils
	auth    auth.Authenticator
	service service.Service
}

func NewMiddleware(db *sql.DB, auth auth.Authenticator) Middleware {
	json := utils.NewJsonUtils()
	error := utils.NewErrorUtils()
	service := service.NewService(db, auth)
	return Middleware{
		json:    json,
		errror:  error,
		auth:    auth,
		service: service,
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
		user, err := m.service.Users.GetProfileByID(ctx, userID)
		if err != nil {
			m.errror.UnauthorizedError(w, r, err)
			return
		}

		log.Println(user.Fullname)
		ctx = context.WithValue(ctx, UserCtx, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
