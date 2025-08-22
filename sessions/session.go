package sessions

import (
	"fmt"
	"net/http"
	
	"github.com/debarbarinantoine/bingo/internal/ctx"
	"github.com/debarbarinantoine/bingo/internal/enum"
	"github.com/debarbarinantoine/bingo/middleware"
	
	"github.com/alexedwards/scs/v2"
	"github.com/rs/zerolog/hlog"
)

const (
	SessionManagerContext             = "sessionManager"
	IsAuthenticatedContext            = "isAuthenticated"
	AuthenticatedUserIDSessionManager = "authenticatedUserID"
)

var (
	Stores = enum.SessionStores
)

func SetSessionManager(sessionManager *scs.SessionManager) middleware.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r = ctx.SetData(r, SessionManagerContext, sessionManager)
			next.ServeHTTP(w, r)
		})
	}
}

func GetSession(r *http.Request) (*scs.SessionManager, error) {
	sessionManager, ok := ctx.GetData(r.Context(), SessionManagerContext).(*scs.SessionManager)
	if !ok {
		hlog.FromRequest(r).Error().Msg("Session Manager not found in ctx")
		return nil, fmt.Errorf("session Manager not found in ctx")
	}
	return sessionManager, nil
}

func Session(sessionManager *scs.SessionManager) middleware.Middleware {
	return sessionManager.LoadAndSave
}

func Login(r *http.Request, id int) error {
	
	// Prevent user from using a null ID
	if id == 0 {
		return fmt.Errorf("invalid user ID")
	}
	
	sessionManager, err := GetSession(r)
	if err != nil {
		hlog.FromRequest(r).Error().Err(err).Msg("Session Manager not found in ctx")
		return err
	}
	sessionManager.Put(r.Context(), AuthenticatedUserIDSessionManager, id)
	return nil
}

func Authenticate(userExists func(id int) (bool, error)) middleware.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			
			sessionManager, err := GetSession(r)
			if err != nil {
				hlog.FromRequest(r).Error().Err(err).Msg("Session Manager not found in ctx")
				http.Error(w, "Session Manager not found", http.StatusInternalServerError)
				return
			}
			
			// getting the userID from the session
			id := sessionManager.GetInt(r.Context(), AuthenticatedUserIDSessionManager)
			
			// if user is not authenticated
			if id == 0 {
				next.ServeHTTP(w, r)
				return
			}
			
			exists, err := userExists(id)
			if err != nil {
				hlog.FromRequest(r).Error().Err(err).Msg("User existence check failed")
				http.Error(w, "User existence check failed", http.StatusInternalServerError)
				return
			}
			
			if exists {
				// setting the user as authenticated in the ctx
				r = ctx.SetData(r, IsAuthenticatedContext, true)
			}
			
			next.ServeHTTP(w, r)
		})
	}
}

func RequireAuthentication(redirectionURL string) middleware.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			
			isAuthenticated, ok := ctx.GetData(r.Context(), IsAuthenticatedContext).(bool)
			if !ok || !isAuthenticated {
				http.Redirect(w, r, redirectionURL, http.StatusSeeOther)
				return
			}
			
			w.Header().Add("Cache-Control", "no-store")
			
			next.ServeHTTP(w, r)
		})
	}
}
