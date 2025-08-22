// middleware.go
package jwtkit

import (
	"net/http"
	
	"github.com/debarbarinantoine/bingo/context"
	"github.com/debarbarinantoine/bingo/middleware"
	"github.com/go-chi/jwtauth/v5"
	"github.com/rs/zerolog/hlog"
)

func SetJWT(jwtConfig *Config) middleware.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r = context.SetCtxData(r, ContextKey, jwtConfig)
			next.ServeHTTP(w, r)
		})
	}
}

func VerifyAndAuthenticateJWT() middleware.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			jwtConfig, err := GetJWT(r)
			if err != nil {
				hlog.FromRequest(r).Error().Err(err).Msg("JwtConfig configuration not found in context")
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			jwtauth.Verifier(jwtConfig.JWTAuth)(jwtauth.Authenticator(jwtConfig.JWTAuth)(next)).ServeHTTP(w, r)
		})
	}
}
