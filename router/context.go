// context.go
package router

import (
	"context"
	"fmt"
	"net/http"
)

type userKey struct{}
type userContextMap map[string]any

func getUserCtxMap(ctx context.Context) userContextMap {
	userCtxMap, ok := ctx.Value(userKey{}).(userContextMap)
	if !ok {
		// TODO -> Handle error better if necessary
		panic(fmt.Errorf("userCtxMap not found in context"))
	}
	return userCtxMap
}

func SetCtxData(r *http.Request, key string, value any) *http.Request {
	ctx := r.Context()
	userCtxMap := getUserCtxMap(ctx)
	userCtxMap[key] = value
	return r.WithContext(context.WithValue(ctx, userKey{}, userCtxMap))
}

func GetCtxData(ctx context.Context, key string) any {
	userCtxMap := getUserCtxMap(ctx)
	data, ok := userCtxMap[key]
	if !ok {
		return nil
	}
	return data
}

func userCtxData(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userCtxMap := make(userContextMap)
		next.ServeHTTP(w, r.WithContext(context.WithValue(ctx, userKey{}, userCtxMap)))
	})
}
