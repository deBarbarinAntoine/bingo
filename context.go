package bingo

import (
	"context"
	"net/http"
	
	"github.com/debarbarinantoine/bingo/internal/ctx"
)

func SetCtxData(r *http.Request, key string, value any) *http.Request {
	return ctx.SetData(r, key, value)
}
func GetCtxData(context context.Context, key string) any {
	return ctx.GetData(context, key)
}
