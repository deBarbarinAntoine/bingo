package binder

import (
	"net/http"
	
	"github.com/debarbarinantoine/bingo/internal/enum"
)

type Form struct {
	dataBind
}

func fetchFormData(bind *dataBind, src any) error {
	r, ok := src.(*http.Request)
	if !ok || r == nil {
		return ErrInvalidSrcType("*http.Request")
	}
	
	err := r.ParseForm()
	if err != nil {
		return ErrParseForm
	}
	
	// Unflatten the form data and store it
	bind.data = unflattenMap(r.PostForm)
	
	return nil
}

// NewForm uses functional options.
//
// Example usage:
// 		q, err := NewForm(&myStruct, &request) // uses default
// 		q, err := NewForm(&myStruct, &request, WithCustomFetcher(myCustomFetcher)) // uses custom
func NewForm(dst any, src *http.Request, opts ...BindOption) (*Form, error) {
	
	err := checkDst(dst)
	if err != nil {
		return nil, ErrInvalidDst
	}
	
	if src == nil {
		return nil, ErrInvalidSrcType("*http.Request")
	}
	
	form := &Form{
		dataBind: dataBind{
			data: make(map[string]any),
			tag:  enum.Tags.Form,
			
			DataSrc:  src,
			DataDist: dst,
			
			// Set the defaults
			MaxMemory: MaxMemoryDefault,
			FetchData: fetchFormData,
		},
	}
	
	// Apply any custom options
	for _, opt := range opts {
		opt(&form.dataBind)
	}
	
	return form, nil
}
