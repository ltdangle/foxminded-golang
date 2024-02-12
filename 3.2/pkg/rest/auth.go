package rest

import (
	"jwt/pkg/infra"
	"jwt/pkg/model"
	"net/http"

	"github.com/golang-jwt/jwt"
)

type Auth struct {
	jwt JwtServiceInterface
	log infra.LoggerInterface
}

func NewAuth(jwt JwtServiceInterface, log infra.LoggerInterface) *Auth {
	return &Auth{jwt: jwt, log: log}
}

func (a *Auth) AdminJwtAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")

		claims, err := a.parseToken(tokenString)
		if err != nil {
			a.log.Warn("Error parsing token: " + err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if claims.Role !=model.ROLE_ADMIN {
			a.log.Warn("User is not admin, denied access.")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		a.log.Info("Token validated for user: " + claims.Email)
		next.ServeHTTP(w, r)
	})
}

func (a *Auth) RegisteredUserJwtAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")

		claims, err := a.parseToken(tokenString)
		if err != nil {
			a.log.Warn("Error parsing token: " + err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if claims.Role == "" {
			a.log.Warn("User is not authorized, denied access.")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		a.log.Info("Token validated for user: " + claims.Email)
		next.ServeHTTP(w, r)
	})
}

func (a *Auth) parseToken(token string) (*Claims, error) {
	claims, err := a.jwt.ParseToken(token)
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			a.log.Warn("Invalid token signature")
			return nil, err
		}
		a.log.Warn("Error parsing token: " + err.Error())
		return nil, err
	}
	return claims, nil
}
