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

type MultiBinderOption func(binder *MultiBinder)

func WithCustomBinder(customBinder *DataBinder) MultiBinderOption {
	myCustomBinder := any(customBinder)
	return func(binder *MultiBinder) {
		switch myCustomBinder.(type) {
		case *Query:
			binder.QueryBinder = myCustomBinder.(*Query)
		case *Form:
			binder.FormBinder = myCustomBinder.(*Form)
		case *MultipartForm:
			binder.MultipartFormBinder = myCustomBinder.(*MultipartForm)
		case *UrlParam:
			binder.UrlParamBinder = myCustomBinder.(*UrlParam)
		case *Header:
			binder.HeaderBinder = myCustomBinder.(*Header)
		case *Cookie:
			binder.CookieBinder = myCustomBinder.(*Cookie)
		case *JSON:
			binder.JSONBinder = myCustomBinder.(*JSON)
		default:
			// do nothing
		}
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
// 		q, err := NewMultiBinder(&myStruct, &request, WithCustomBinder(myCustomBinder)) // uses custom
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

func (binder *MultiBinder) createBinders(dst any, src *http.Request) error {
	presentTags := hasTags(dst)
	var err error
	if _, ok := presentTags[enum.Tags.Query]; ok {
		if binder.QueryBinder, err = NewQuery(dst, src.URL.Query()); err != nil {
			return err
		}
	}
	if _, ok := presentTags[enum.Tags.Form]; ok {
		if binder.FormBinder, err = NewForm(dst, src); err != nil {
			return err
		}
	}
	if _, ok := presentTags[enum.Tags.MultipartForm]; ok {
		if binder.MultipartFormBinder, err = NewMultipartForm(dst, src); err != nil {
			return err
		}
	}
	if _, ok := presentTags[enum.Tags.UrlParam]; ok {
		if binder.UrlParamBinder, err = NewUrlParam(dst, src); err != nil {
			return err
		}
	}
	if _, ok := presentTags[enum.Tags.Header]; ok {
		if binder.HeaderBinder, err = NewHeader(dst, src); err != nil {
			return err
		}
	}
	if _, ok := presentTags[enum.Tags.Cookie]; ok {
		if binder.CookieBinder, err = NewCookie(dst, src); err != nil {
			return err
		}
	}
	if _, ok := presentTags[enum.Tags.Json]; ok {
		if binder.JSONBinder, err = NewJSON(dst, src); err != nil {
			return err
		}
	}
	return nil
}

func (binder *MultiBinder) Fetch() error {
	if binder.UrlParamBinder != nil {
		if err := binder.UrlParamBinder.Fetch(); err != nil {
			return err
		}
	}
	if binder.HeaderBinder != nil {
		if err := binder.HeaderBinder.Fetch(); err != nil {
			return err
		}
	}
	if binder.CookieBinder != nil {
		if err := binder.CookieBinder.Fetch(); err != nil {
			return err
		}
	}
	if binder.QueryBinder != nil {
		if err := binder.QueryBinder.Fetch(); err != nil {
			return err
		}
	}
	if binder.JSONBinder != nil {
		if err := binder.JSONBinder.Fetch(); err != nil {
			return err
		}
	}
	if binder.FormBinder != nil {
		if err := binder.FormBinder.Fetch(); err != nil {
			return err
		}
	}
	if binder.MultipartFormBinder != nil {
		if err := binder.MultipartFormBinder.Fetch(); err != nil {
			return err
		}
	}
	return nil
}
