package jwtkit

import (
	"fmt"
	"net/http"
	
	"github.com/debarbarinantoine/bingo/internal/ctx"
	"github.com/rs/zerolog/hlog"
)

const (
	ContextKey = "jwtkit"
)

func GetJWT(r *http.Request) (*Config, error) {
	jwtConfig, ok := ctx.GetData(r.Context(), ContextKey).(*Config)
	if !ok {
		hlog.FromRequest(r).Error().Msg("jwtConfig not found in ctx")
		return nil, fmt.Errorf("jwtConfig not found in ctx")
	}
	return jwtConfig, nil
}
