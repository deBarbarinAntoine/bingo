// helpers.go
package jwtkit

import (
	"fmt"
	"net/http"
	"time"
	
	"github.com/lestrrat-go/jwx/v2/jwt"
)

func EncodeJWT(r *http.Request, claims map[string]interface{}) (jwt.Token, string, error) {
	jwtConfig, err := GetJWT(r)
	if err != nil {
		return nil, "", err
	}
	return jwtConfig.JWTAuth.Encode(claims)
}

func setTokenInResponse(r *http.Request, w http.ResponseWriter, token jwt.Token, tokenString string, opts *TokenResponseOptions) {
	if token == nil || tokenString == "" {
		return
	}
	
	if opts == nil {
		opts = DefaultTokenResponseOptions()
	}
	
	// Set Authorization header
	if opts.SetAuthorizationHeader {
		w.Header().Set("Authorization", fmt.Sprintf("Bearer %s", tokenString))
	}
	
	// Set custom header
	if opts.SetCustomHeader && opts.CustomHeaderName != "" {
		w.Header().Set(opts.CustomHeaderName, tokenString)
	}
	
	// Set expiration header
	if opts.SetExpirationHeader {
		w.Header().Set("X-Auth-Expires", token.Expiration().Format(time.RFC3339))
	}
	
	// Set cookie
	if opts.SetInCookie {
		cookie := &http.Cookie{
			Name:     opts.CookieName,
			Value:    tokenString,
			Expires:  token.Expiration(),
			Secure:   opts.CookieSecure,
			HttpOnly: opts.CookieHttpOnly,
			SameSite: opts.CookieSameSite,
			Path:     "/",
		}
		http.SetCookie(w, cookie)
	}
}

func SetTokenInResponse(r *http.Request, w http.ResponseWriter, token jwt.Token, tokenString string) error {
	jwtConfig, err := GetJWT(r)
	if err != nil {
		return err
	}
	setTokenInResponse(r, w, token, tokenString, jwtConfig.Options)
	return nil
}
