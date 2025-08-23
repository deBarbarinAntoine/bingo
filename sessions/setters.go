package sessions

import "net/http"

// Set adds a key and corresponding value to the session data.
// Any existing value for the key will be replaced.
// The session data status will be set to Modified.
func Set(r *http.Request, key string, value any) error {
	sessionManager, err := GetSession(r)
	if err != nil {
		return err
	}
	
	sessionManager.Put(r.Context(), key, value)
	
	return nil
}

// Put calls Set to save a new session key/value in the session.
//
// It exists for those who are used to the SessionManager.Put() from alexedwards/scs
func Put(r *http.Request, key string, value any) error {
	return Set(r, key, value)
}
