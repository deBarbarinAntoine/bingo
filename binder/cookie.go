package binder

import (
	"net/http"
	
	"BinGo/enum"
)

const (
	CookieHeader          = "Cookie"
	CookieSeparator       = ";"
	CookieAttributionSign = "="
)

type Cookie struct {
	dataBind
}

func fetchCookieData(bind *dataBind, src any) error {
	r, ok := src.(*http.Request)
	if !ok || r == nil {
		return ErrInvalidSrcType("*http.Request")
	}
	// Get all tags from the destination struct
	bind.getTags()
	
	cookies := r.Cookies()
	
	// Create a temporary map for easier lookup
	cookieMap := make(map[string]string)
	for _, cookie := range cookies {
		cookieMap[cookie.Name] = cookie.Value
	}
	
	// Iterate through the collected tags and populate the binder's data map
	for k := range bind.data {
		if val, ok := cookieMap[k]; ok {
			bind.data[k] = val
		}
	}
	
	return nil
}

// NewCookie uses functional options.
//
// Example usage:
// 		q, err := NewCookie(&myStruct, &request) // uses default
// 		q, err := NewCookie(&myStruct, &request, WithCustomFetcher(myCustomFetcher)) // uses custom
func NewCookie(dst any, src *http.Request, opts ...BindOption) (*Cookie, error) {
	err := checkDst(dst)
	if err != nil {
		return nil, ErrInvalidDst
	}
	if src == nil {
		return nil, ErrInvalidSrcType("*http.Request")
	}
	form := &Cookie{
		dataBind: dataBind{
			data: make(map[string]any),
			tag:  enum.Tags.Cookie,
			
			DataSrc:  src,
			DataDist: dst,
			
			// Set the defaults
			MaxMemory: MaxMemoryDefault,
			FetchData: fetchCookieData,
		},
	}
	
	// Apply any custom options
	for _, opt := range opts {
		opt(&form.dataBind)
	}
	
	return form, nil
}
