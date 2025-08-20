// multibinder.go
package binder

import (
	"net/http"
	
	"BinGo/enum"
)

type MultiBinder struct {
	QueryBinder         *Query
	FormBinder          *Form
	MultipartFormBinder *MultipartForm
	UrlParamBinder      *UrlParam
	HeaderBinder        *Header
	CookieBinder        *Cookie
	JSONBinder          *JSON
}

type MultiBinderOption func(mb *MultiBinder)

func WithBodyBinder(tag enum.Tag) []MultiBinderOption {
	binderOptions := make([]MultiBinderOption, 0, 3)
	switch tag {
	case enum.Tags.Form:
		binderOptions = append(binderOptions, WithoutMultipartFormBinder(), WithoutJSONBinder())
	case enum.Tags.MultipartForm:
		binderOptions = append(binderOptions, WithoutFormBinder(), WithoutJSONBinder())
	case enum.Tags.Json:
		binderOptions = append(binderOptions, WithoutMultipartFormBinder(), WithoutFormBinder())
	default:
		return nil
	}
	return binderOptions
}

func WithCustomQueryBinder(customBinder *Query) MultiBinderOption {
	return func(mb *MultiBinder) {
		mb.QueryBinder = customBinder
	}
}

func WithCustomMultipartBinder(customBinder *MultipartForm) MultiBinderOption {
	return func(mb *MultiBinder) {
		mb.MultipartFormBinder = customBinder
	}
}

func WithCustomFormBinder(customBinder *Form) MultiBinderOption {
	return func(mb *MultiBinder) {
		mb.FormBinder = customBinder
	}
}

func WithCustomUrlParamBinder(customBinder *UrlParam) MultiBinderOption {
	return func(mb *MultiBinder) {
		mb.UrlParamBinder = customBinder
	}
}

func WithCustomHeaderBinder(customBinder *Header) MultiBinderOption {
	return func(mb *MultiBinder) {
		mb.HeaderBinder = customBinder
	}
}

func WithCustomCookieBinder(customBinder *Cookie) MultiBinderOption {
	return func(mb *MultiBinder) {
		mb.CookieBinder = customBinder
	}
}

func WithCustomJSONBinder(customBinder *JSON) MultiBinderOption {
	return func(mb *MultiBinder) {
		mb.JSONBinder = customBinder
	}
}

func WithoutJSONBinder() MultiBinderOption {
	return func(mb *MultiBinder) {
		mb.JSONBinder = nil
	}
}

func WithoutFormBinder() MultiBinderOption {
	return func(mb *MultiBinder) {
		mb.FormBinder = nil
	}
}

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
