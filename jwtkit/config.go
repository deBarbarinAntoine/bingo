// config.go
package jwtkit

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rsa"
	"fmt"
	"net/http"
	
	"github.com/go-chi/jwtauth/v5"
)

// TokenResponseOptions configures how tokens are returned in responses
type TokenResponseOptions struct {
	SetAuthorizationHeader bool
	SetCustomHeader        bool
	CustomHeaderName       string
	SetExpirationHeader    bool
	SetInCookie            bool
	CookieName             string
	CookieSecure           bool
	CookieHttpOnly         bool
	CookieSameSite         http.SameSite
}

// DefaultTokenResponseOptions returns sensible defaults
func DefaultTokenResponseOptions() *TokenResponseOptions {
	return &TokenResponseOptions{
		SetAuthorizationHeader: true,
		SetCustomHeader:        true,
		CustomHeaderName:       "X-Auth-Token",
		SetExpirationHeader:    true,
		SetInCookie:            false, // Usually tokens in cookies need special handling
		CookieName:             "auth_token",
		CookieSecure:           true,
		CookieHttpOnly:         true,
		CookieSameSite:         http.SameSiteStrictMode,
	}
}

// Config is the configuration for the JWT auth middleware
type Config struct {
	Algorithm Algorithm
	JWTAuth   *jwtauth.JWTAuth
	Options   *TokenResponseOptions
}

// NewConfigWithSecret creates a new JWT config with a secret
func NewConfigWithSecret(algorithm Algorithm, secret string, options *TokenResponseOptions) (*Config, error) {
	cfg := &Config{
		Algorithm: algorithm,
		Options:   options,
	}
	var err error
	cfg.JWTAuth, err = newAuthWithSecret(algorithm, secret)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

// NewConfigWithRSA creates a new JWT config with a RSA private key
func NewConfigWithRSA(algorithm Algorithm, privateKey *rsa.PrivateKey, options *TokenResponseOptions) (*Config, error) {
	cfg := &Config{
		Algorithm: algorithm,
		Options:   options,
	}
	var err error
	cfg.JWTAuth, err = newAuthWithRSA(algorithm, privateKey)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

// NewConfigWithECDSA creates a new JWT config with a ECDSA private key
func NewConfigWithECDSA(algorithm Algorithm, privateKey *ecdsa.PrivateKey, options *TokenResponseOptions) (*Config, error) {
	cfg := &Config{
		Algorithm: algorithm,
		Options:   options,
	}
	var err error
	cfg.JWTAuth, err = newAuthWithECDSA(algorithm, privateKey)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

// NewConfigWithEdDSA creates a new JWT config with a EdDSA private key
func NewConfigWithEdDSA(algorithm Algorithm, privateKey ed25519.PrivateKey, options *TokenResponseOptions) (*Config, error) {
	cfg := &Config{
		Algorithm: algorithm,
		Options:   options,
	}
	var err error
	cfg.JWTAuth, err = newAuthWithEdDSA(algorithm, privateKey)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

// NewConfigUnsigned creates a new JWT config with no signature
// This is useful for testing or for APIs that don't require authentication
//
// This is NOT SAFE, do NOT use this in production!
func NewConfigUnsigned(options *TokenResponseOptions) *Config {
	return &Config{
		Algorithm: AlgorithmNone,
		Options:   options,
		JWTAuth:   newAuthUnsigned(),
	}
}

func newAuthWithSecret(algorithm Algorithm, secret string) (*jwtauth.JWTAuth, error) {
	if !algorithm.IsSymmetric() {
		return nil, fmt.Errorf("algorithm %s is not a symmetric algorithm", algorithm)
	}
	
	if secret == "" {
		return nil, fmt.Errorf("secret cannot be empty for symmetric algorithms")
	}
	
	return jwtauth.New(string(algorithm), []byte(secret), nil), nil
}

func newAuthWithRSA(algorithm Algorithm, privateKey *rsa.PrivateKey) (*jwtauth.JWTAuth, error) {
	if !algorithm.IsRSA() {
		return nil, fmt.Errorf("algorithm %s is not an RSA algorithm", algorithm)
	}
	
	if privateKey == nil {
		return nil, fmt.Errorf("private key cannot be nil")
	}
	
	return jwtauth.New(string(algorithm), privateKey, nil), nil
}

func newAuthWithECDSA(algorithm Algorithm, privateKey *ecdsa.PrivateKey) (*jwtauth.JWTAuth, error) {
	if !algorithm.IsECDSA() {
		return nil, fmt.Errorf("algorithm %s is not an ECDSA algorithm", algorithm)
	}
	
	if privateKey == nil {
		return nil, fmt.Errorf("private key cannot be nil")
	}
	
	return jwtauth.New(string(algorithm), privateKey, nil), nil
}

func newAuthWithEdDSA(algorithm Algorithm, privateKey ed25519.PrivateKey) (*jwtauth.JWTAuth, error) {
	if !algorithm.IsEdDSA() {
		return nil, fmt.Errorf("algorithm %s is not an EdDSA algorithm", algorithm)
	}
	
	if privateKey == nil {
		return nil, fmt.Errorf("private key cannot be nil")
	}
	
	return jwtauth.New(string(algorithm), privateKey, nil), nil
}

func newAuthUnsigned() *jwtauth.JWTAuth {
	return jwtauth.New(string(AlgorithmNone), nil, nil)
}
