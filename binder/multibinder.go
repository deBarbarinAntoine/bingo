package binder

import (
	"net/http"
	
	"github.com/debarbarinantoine/bingo/internal/enum"
)

// MultiBinder is a struct that contains multiple binders for different data sources.
type MultiBinder struct {
	QueryBinder         *Query
	FormBinder          *Form
	MultipartFormBinder *MultipartForm
	UrlParamBinder      *UrlParam
	HeaderBinder        *Header
	CookieBinder        *Cookie
	JSONBinder          *JSON
}

// MultiBinderOption is a function that configures a MultiBinder instance.
type MultiBinderOption func(mb *MultiBinder)

// WithJsonBodyBinder only keeps the JSON binder for the MultiBinder instance.
func WithJsonBodyBinder() MultiBinderOption {
	return func(mb *MultiBinder) {
		mb.FormBinder = nil
		mb.MultipartFormBinder = nil
	}
}

// WithFormBodyBinder only keeps the form binder for the MultiBinder instance.
func WithFormBodyBinder() MultiBinderOption {
	return func(mb *MultiBinder) {
		mb.JSONBinder = nil
		mb.MultipartFormBinder = nil
	}
}

// WithMultipartFormBodyBinder only keeps the multipart form binder for the MultiBinder instance.
func WithMultipartFormBodyBinder() MultiBinderOption {
	return func(mb *MultiBinder) {
		mb.JSONBinder = nil
		mb.FormBinder = nil
	}
}

// WithoutBodyBinder sets the body binders for the MultiBinder instance to nil.
func WithoutBodyBinder() MultiBinderOption {
	return func(mb *MultiBinder) {
		mb.JSONBinder = nil
		mb.FormBinder = nil
		mb.MultipartFormBinder = nil
	}
}

// WithCustomQueryBinder sets a custom query binder for the MultiBinder instance.
func WithCustomQueryBinder(customBinder *Query) MultiBinderOption {
	return func(mb *MultiBinder) {
		mb.QueryBinder = customBinder
	}
}

// WithCustomMultipartBinder sets a custom multipart form binder for the MultiBinder instance.
func WithCustomMultipartBinder(customBinder *MultipartForm) MultiBinderOption {
	return func(mb *MultiBinder) {
		mb.MultipartFormBinder = customBinder
	}
}

// WithCustomFormBinder sets a custom form binder for the MultiBinder instance.
func WithCustomFormBinder(customBinder *Form) MultiBinderOption {
	return func(mb *MultiBinder) {
		mb.FormBinder = customBinder
	}
}

// WithCustomUrlParamBinder sets a custom URL param binder for the MultiBinder instance.
func WithCustomUrlParamBinder(customBinder *UrlParam) MultiBinderOption {
	return func(mb *MultiBinder) {
		mb.UrlParamBinder = customBinder
	}
}

// WithCustomHeaderBinder sets a custom header binder for the MultiBinder instance.
func WithCustomHeaderBinder(customBinder *Header) MultiBinderOption {
	return func(mb *MultiBinder) {
		mb.HeaderBinder = customBinder
	}
}

// WithCustomCookieBinder sets a custom cookie binder for the MultiBinder instance.
func WithCustomCookieBinder(customBinder *Cookie) MultiBinderOption {
	return func(mb *MultiBinder) {
		mb.CookieBinder = customBinder
	}
}

// WithCustomJSONBinder sets a custom JSON binder for the MultiBinder instance.
func WithCustomJSONBinder(customBinder *JSON) MultiBinderOption {
	return func(mb *MultiBinder) {
		mb.JSONBinder = customBinder
	}
}

// WithoutJSONBinder removes the JSON binder from the MultiBinder instance.
func WithoutJSONBinder() MultiBinderOption {
	return func(mb *MultiBinder) {
		mb.JSONBinder = nil
	}
}

// WithoutFormBinder removes the form binder from the MultiBinder instance.
func WithoutFormBinder() MultiBinderOption {
	return func(mb *MultiBinder) {
		mb.FormBinder = nil
	}
}

// WithoutMultipartFormBinder removes the multipart form binder from the MultiBinder instance.
func WithoutMultipartFormBinder() MultiBinderOption {
	return func(mb *MultiBinder) {
		mb.MultipartFormBinder = nil
	}
}

// NewMultiBinder now uses functional options.
//
// Example usage:
// 		q, err := NewMultiBinder(&myStruct, &request) // uses default
//
// 		myCustomBinder, err := NewUrlParam(&myStruct, &request, WithCustomFetcher(myCustomFetcher))
// 		if err != nil {
// 			panic(err)
// 		}
// 		q, err := NewMultiBinder(&myStruct, &request, WithCustomUrlParamBinder(myCustomBinder)) // uses custom
func NewMultiBinder(dst any, src *http.Request, opts ...MultiBinderOption) (*MultiBinder, error) {
	
	err := checkDst(dst)
	if err != nil {
		return nil, ErrInvalidDst
	}
	
	if src == nil {
		return nil, ErrInvalidSrcType("*http.Request")
	}
	
	binder := &MultiBinder{}
	
	// Creates necessary Binders
	if err = binder.createBinders(dst, src); err != nil {
		return nil, err
	}
	
	// Apply any custom options
	for _, opt := range opts {
		opt(binder)
	}
	
	return binder, nil
}

func (mb *MultiBinder) createBinders(dst any, src *http.Request) error {
	presentTags := hasTags(dst)
	var err error
	if _, ok := presentTags[enum.Tags.Query]; ok {
		if mb.QueryBinder, err = NewQuery(dst, src.URL.Query()); err != nil {
			return err
		}
	}
	if _, ok := presentTags[enum.Tags.Form]; ok {
		if mb.FormBinder, err = NewForm(dst, src); err != nil {
			return err
		}
	}
	if _, ok := presentTags[enum.Tags.MultipartForm]; ok {
		if mb.MultipartFormBinder, err = NewMultipartForm(dst, src); err != nil {
			return err
		}
	}
	if _, ok := presentTags[enum.Tags.UrlParam]; ok {
		if mb.UrlParamBinder, err = NewUrlParam(dst, src); err != nil {
			return err
		}
	}
	if _, ok := presentTags[enum.Tags.Header]; ok {
		if mb.HeaderBinder, err = NewHeader(dst, src); err != nil {
			return err
		}
	}
	if _, ok := presentTags[enum.Tags.Cookie]; ok {
		if mb.CookieBinder, err = NewCookie(dst, src); err != nil {
			return err
		}
	}
	if _, ok := presentTags[enum.Tags.Json]; ok {
		if mb.JSONBinder, err = NewJSON(dst, src); err != nil {
			return err
		}
	}
	return nil
}

// Fetch fetches data from the source and binds it to the destination struct.
func (mb *MultiBinder) Fetch() error {
	if mb.JSONBinder != nil {
		if err := mb.JSONBinder.Fetch(); err != nil {
			return err
		}
	}
	if mb.FormBinder != nil {
		if err := mb.FormBinder.Fetch(); err != nil {
			return err
		}
	}
	if mb.MultipartFormBinder != nil {
		if err := mb.MultipartFormBinder.Fetch(); err != nil {
			return err
		}
	}
	if mb.QueryBinder != nil {
		if err := mb.QueryBinder.Fetch(); err != nil {
			return err
		}
	}
	if mb.UrlParamBinder != nil {
		if err := mb.UrlParamBinder.Fetch(); err != nil {
			return err
		}
	}
	if mb.HeaderBinder != nil {
		if err := mb.HeaderBinder.Fetch(); err != nil {
			return err
		}
	}
	if mb.CookieBinder != nil {
		if err := mb.CookieBinder.Fetch(); err != nil {
			return err
		}
	}
	return nil
}
