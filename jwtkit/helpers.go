package jwtkit

import (
	"fmt"
	"net/http"
	"time"
	
	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/v3/jwt"
	"github.com/rs/zerolog/hlog"
)

type encodingOption struct {
	IssuedAt time.Time
	Subject  string
	JwtID    string
	Audience []string
	Issuer   string
	TTL      time.Duration
}

// EncodingOptions is a functional option type for encoding JWT tokens.
type EncodingOptions func(*encodingOption)

// WithIssuedAt sets the issued at time for the JWT token.
func WithIssuedAt(issuedAt time.Time) EncodingOptions {
	return func(o *encodingOption) {
		o.IssuedAt = issuedAt
	}
}

// WithSubject sets the subject for the JWT token.
func WithSubject(subject string) EncodingOptions {
	return func(o *encodingOption) {
		o.Subject = subject
	}
}

// WithJwtID sets the JWT ID for the JWT token.
func WithJwtID(jwtID string) EncodingOptions {
	return func(o *encodingOption) {
		o.JwtID = jwtID
	}
}

// WithAudience sets the audience for the JWT token.
func WithAudience(audience []string) EncodingOptions {
	return func(o *encodingOption) {
		o.Audience = audience
	}
}

// WithIssuer sets the issuer for the JWT token.
func WithIssuer(issuer string) EncodingOptions {
	return func(o *encodingOption) {
		o.Issuer = issuer
	}
}

// WithTTL sets the time-to-live duration for the JWT token.
func WithTTL(ttl time.Duration) EncodingOptions {
	return func(o *encodingOption) {
		o.TTL = ttl
	}
}

// EncodeJWT encodes a JWT token with the provided claims and options.
func EncodeJWT(r *http.Request, claims map[string]any, opts ...EncodingOptions) (jwt.Token, string, error) {
	jwtConfig, err := GetJWT(r)
	if err != nil {
		return nil, "", err
	}
	
	options := &encodingOption{
		IssuedAt: time.Now(),
		JwtID:    uuid.NewString(),
		Audience: jwtConfig.Audience,
		Issuer:   jwtConfig.Issuer,
		TTL:      jwtConfig.TTL,
	}
	
	for _, opt := range opts {
		opt(options)
	}
	
	builder := jwt.NewBuilder()
	if options.Subject != "" {
		builder = builder.Subject(options.Subject)
	}
	if options.JwtID != "" {
		builder = builder.JwtID(options.JwtID)
	}
	if options.Audience != nil {
		builder = builder.Audience(options.Audience)
	}
	if options.Issuer != "" {
		builder = builder.Issuer(options.Issuer)
	}
	
	// Build the token with the provided options
	token, err := builder.
		IssuedAt(options.IssuedAt).
		Expiration(time.Now().Add(options.TTL)).
		Build()
	if err != nil {
		return nil, "", err
	}
	
	// Set claims in the token
	for key, value := range claims {
		switch key {
		// Skip reserved keys
		case jwt.ExpirationKey, jwt.IssuedAtKey, jwt.JwtIDKey, jwt.AudienceKey, jwt.IssuerKey, jwt.SubjectKey:
			continue
		}
		err = token.Set(key, value)
		if err != nil {
			return nil, "", err
		}
	}
	
	var key any
	if jwtConfig.Algorithm.IsSymmetric() {
		if jwtConfig.secret == "" {
			return nil, "", fmt.Errorf("secret is required for symmetric algorithm")
		}
		key = []byte(jwtConfig.secret)
	} else if jwtConfig.Algorithm.IsAsymmetric() {
		if jwtConfig.privateKey == nil {
			return nil, "", fmt.Errorf("private key is required for asymmetric algorithm")
		}
		key = jwtConfig.privateKey
	} else {
		key = nil
	}
	
	tokenStr, err := jwt.Sign(token, jwt.WithKey(jwtConfig.Algorithm.toJwaAlgo(), key))
	if err != nil {
		hlog.FromRequest(r).Error().Err(err).Msg("failed to sign token")
		return nil, "", fmt.Errorf("failed to sign token: %w", err)
	}
	
	return token, string(tokenStr), nil
}

// func EncodeJWT(r *http.Request, claims map[string]interface{}) (jwt.Token, string, error) {
// 	jwtConfig, err := GetJWT(r)
// 	if err != nil {
// 		return nil, "", err
// 	}
//
// 	// Set TTL in the token if necessary
// 	if _, ok := claims[jwt.ExpirationKey]; jwtConfig.Options.TTL != 0 && !ok {
// 		claims[jwt.ExpirationKey] = time.Now().Add(jwtConfig.Options.TTL).Unix()
// 	}
//
// 	return jwtConfig.JWTAuth.Encode(claims)
// }

func setTokenInResponse(w http.ResponseWriter, token jwt.Token, tokenString string, opts *TokenResponseOptions) error {
	if token == nil || tokenString == "" {
		return fmt.Errorf("token or tokenString is empty")
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
	exp, ok := token.Expiration()
	if ok && opts.SetExpirationHeader {
		w.Header().Set("X-Auth-Expires", exp.Format(time.RFC3339))
	}
	
	// Set cookie
	if opts.SetInCookie {
		exp, hasExp := token.Expiration()
		cookie := &http.Cookie{
			Name:     opts.CookieName,
			Value:    tokenString,
			Secure:   opts.CookieSecure,
			HttpOnly: opts.CookieHttpOnly,
			SameSite: opts.CookieSameSite,
			Path:     "/",
		}
		
		// Only set expiration if the token has one
		if hasExp {
			cookie.Expires = exp
		}
		
		http.SetCookie(w, cookie)
	}
	return nil
}

func SetTokenInResponse(r *http.Request, w http.ResponseWriter, token jwt.Token, tokenString string) error {
	jwtConfig, err := GetJWT(r)
	if err != nil {
		return err
	}
	return setTokenInResponse(w, token, tokenString, jwtConfig.Options)
}
