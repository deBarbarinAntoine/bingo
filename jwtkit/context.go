// context.go
package jwtkit

import (
	"fmt"
	"net/http"
	
	"github.com/debarbarinantoine/bingo/context"
	"github.com/rs/zerolog/hlog"
)

const (
	ContextKey = "jwtkit"
)

func GetJWT(r *http.Request) (*Config, error) {
	jwtConfig, ok := context.GetCtxData(r.Context(), ContextKey).(*Config)
	if !ok {
		hlog.FromRequest(r).Error().Msg("jwtConfig not found in context")
		return nil, fmt.Errorf("jwtConfig not found in context")
	}
	return jwtConfig, nil
}
