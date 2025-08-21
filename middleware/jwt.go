// jwt.go
package middleware

import (
	"net/http"
	
	"github.com/debarbarinantoine/bingo/context"
	
	"github.com/go-chi/jwtauth/v5"
	"github.com/rs/zerolog/hlog"
)

const (
	JwtContext = "jwt"
)

func SetJWT(jwt *jwtauth.JWTAuth) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r = context.SetCtxData(r, JwtContext, jwt)
			next.ServeHTTP(w, r)
		})
	}
}

func VerifyAndAuthenticateJWT() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			jwt, ok := context.GetCtxData(r.Context(), JwtContext).(*jwtauth.JWTAuth)
			if !ok {
				hlog.FromRequest(r).Error().Msg("JWT Auth configuration not found in context")
				http.Error(w, "JWT Auth not found", http.StatusInternalServerError)
				return
			}
			jwtauth.Verifier(jwt)(jwtauth.Authenticator(jwt)(next)).ServeHTTP(w, r)
		})
	}
}
