package sessions

import (
	"fmt"
	"net/http"
	
	"github.com/debarbarinantoine/bingo/internal/ctx"
	"github.com/debarbarinantoine/bingo/internal/enum"
	"github.com/debarbarinantoine/bingo/internal/helpers"
	"github.com/debarbarinantoine/bingo/middleware"
	
	"github.com/alexedwards/scs/v2"
	"github.com/rs/zerolog/hlog"
)

const (
	// SessionManagerContext is the key used to store the session manager in the request context.
	SessionManagerContext = "sessionManager"
	
	// IsAuthenticatedContext is the key used to store the authentication status in the request context.
	IsAuthenticatedContext = "isAuthenticated"
	
	// AuthenticatedUserIDSessionManager is the key used to store the authenticated user ID in the session.
	AuthenticatedUserIDSessionManager = "authenticatedUserID"
)

var (
	// Stores is a list of supported session stores.
	Stores = enum.SessionStores
)

// SetSessionManager sets the session manager in the request context.
//
// It is automatically called when setting the router with session handling.
func SetSessionManager(sessionManager *scs.SessionManager) middleware.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r = ctx.SetData(r, SessionManagerContext, sessionManager)
			next.ServeHTTP(w, r)
		})
	}
}

// GetSession retrieves the session manager from the request context.
func GetSession(r *http.Request) (*scs.SessionManager, error) {
	sessionManager, ok := ctx.GetData(r.Context(), SessionManagerContext).(*scs.SessionManager)
	if !ok {
		hlog.FromRequest(r).Error().Msg("Session Manager not found in ctx")
		return nil, fmt.Errorf("session Manager not found in ctx")
	}
	return sessionManager, nil
}

// Session returns a middleware that loads and saves session data for each request.
//
// This middleware is needed to be used before any other middleware that requires session data.
func Session(sessionManager *scs.SessionManager) middleware.Middleware {
	return sessionManager.LoadAndSave
}

// Login logs in a user by setting their ID in the session.
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

// Logout logs out a user by clearing their session.
func Logout(r *http.Request) error {
	
	sessionManager, err := GetSession(r)
	if err != nil {
		hlog.FromRequest(r).Error().Err(err).Msg("Session Manager not found in ctx")
		return err
	}
	
	err = sessionManager.Clear(r.Context())
	if err != nil {
		return err
	}
	
	err = sessionManager.RenewToken(r.Context())
	if err != nil {
		return err
	}
	
	return nil
}

// Authenticate checks if a user is authenticated by verifying their ID in the session.
//
// This middleware is not blocking, it only sets a flag in the context.
func Authenticate(userExists func(id int) (bool, error)) middleware.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			
			sessionManager, err := GetSession(r)
			if err != nil {
				hlog.FromRequest(r).Error().Err(err).Msg("Session Manager not found in ctx")
				helpers.ServerError(r, w, err, "Session Manager not found")
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
				helpers.ServerError(r, w, err, "User existence check failed")
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

// RequireAuthentication redirects unauthenticated users to the specified URL.
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
