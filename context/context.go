// context.go
package context

import (
	"context"
	"fmt"
	"net/http"
)

type contextKey struct{}
type contextMap map[string]any

func getCtxMap(ctx context.Context) contextMap {
	ctxMap, ok := ctx.Value(contextKey{}).(contextMap)
	if !ok {
		// TODO -> Handle error better if necessary
		panic(fmt.Errorf("contextMap not found in context"))
	}
	return ctxMap
}

func SetCtxData(r *http.Request, key string, value any) *http.Request {
	ctx := r.Context()
	ctxMap := getCtxMap(ctx)
	ctxMap[key] = value
	return r.WithContext(context.WithValue(ctx, contextKey{}, ctxMap))
}

func GetCtxData(ctx context.Context, key string) any {
	ctxMap := getCtxMap(ctx)
	data, ok := ctxMap[key]
	if !ok {
		return nil
	}
	return data
}

func CtxData(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctxMap := make(contextMap)
		next.ServeHTTP(w, r.WithContext(context.WithValue(ctx, contextKey{}, ctxMap)))
	})
}
