package binder

import (
	"encoding/json"
	"net/http"
	
	"github.com/debarbarinantoine/bingo/internal/enum"
)

type JSON struct {
	dataBind
}

func fetchJSONData(bind *dataBind, src any) error {
	r, ok := src.(*http.Request)
	if !ok || r == nil {
		return ErrInvalidSrcType("*http.Request")
	}
	
	if err := json.NewDecoder(r.Body).Decode(bind.DataDist); err != nil {
		return err
	}
	
	return nil
}

// NewJSON uses functional options.
//
// Example usage:
// 		q, err := NewJSON(&myStruct, &request) // uses default
// 		q, err := NewJSON(&myStruct, &request, WithCustomFetcher(myCustomFetcher)) // uses custom
func NewJSON(dst any, src *http.Request, opts ...BindOption) (*JSON, error) {
	
	err := checkDst(dst)
	if err != nil {
		return nil, ErrInvalidDst
	}
	
	if src == nil {
		return nil, ErrInvalidSrcType("*http.Request")
	}
	
	form := &JSON{
		dataBind: dataBind{
			tag: enum.Tags.Json,
			
			DataSrc:  src,
			DataDist: dst,
			
			// Set the defaults
			MaxMemory: MaxMemoryDefault,
			FetchData: fetchJSONData,
		},
	}
	
	// Apply any custom options
	for _, opt := range opts {
		opt(&form.dataBind)
	}
	
	return form, nil
}

func (b *JSON) Fetch() error {
	return b.FetchData(&b.dataBind, b.DataSrc)
}
