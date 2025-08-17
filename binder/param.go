package binder

import (
	"net/http"
	
	"BinGo/enum"
)

type UrlParam struct {
	dataBind
}

func fetchUrlParamData(bind *dataBind, src any) error {
	r, ok := src.(*http.Request)
	if !ok || r == nil {
		return ErrInvalidSrc
	}
	
	err := bind.getTags()
	if err != nil {
		return err
	}
	
	for k := range bind.data {
		if v := r.PathValue(k); v != "" {
			bind.data[k] = v
		}
	}
	
	return nil
}

// NewUrlParam now uses functional options.
//
// Example usage:
// 		q, err := NewUrlParam(&myStruct, &request) // uses default
// 		q, err := NewUrlParam(&myStruct, &request, WithCustomFetcher(myCustomFetcher)) // uses custom
func NewUrlParam(dst any, src *http.Request, opts ...BindOption) (*UrlParam, error) {
	err := checkDst(dst)
	if err != nil {
		return nil, ErrInvalidDst
	}
	if src == nil {
		return nil, ErrInvalidSrc
	}
	form := &UrlParam{
		dataBind: dataBind{
			data: make(map[string]any),
			tag:  enum.Tags.UrlParam,
			
			DataSrc:  src,
			DataDist: dst,
			
			// Set the defaults
			MaxMemory: MaxMemoryDefault,
			FetchData: fetchUrlParamData,
		},
	}
	
	// Apply any custom options
	for _, opt := range opts {
		opt(&form.dataBind)
	}
	
	return form, nil
}
