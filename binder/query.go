package binder

import (
	"net/url"
	
	"BinGo/enum"
)

type Query struct {
	dataBind
}

func fetchQueryData(bind *dataBind, src any) error {
	query, ok := src.(url.Values)
	if !ok {
		return ErrInvalidSrc
	}
	
	for k, v := range query {
		if _, ok := bind.data[k]; ok {
			if len(v) == 1 {
				bind.data[k] = v[0]
			} else {
				bind.data[k] = v
			}
		}
	}
	
	return nil
}

// NewQuery now uses functional options.
//
// Example usage:
// 		q, err := NewQuery(&myStruct, urlValues) // uses default
// 		q, err := NewQuery(&myStruct, urlValues, WithCustomFetcher(myCustomFetcher)) // uses custom
func NewQuery(dst any, src url.Values, opts ...BindOption) (*Query, error) {
	
	err := checkDst(dst)
	if err != nil {
		return nil, ErrInvalidDst
	}
	
	if src == nil {
		return nil, ErrInvalidSrc
	}
	
	query := &Query{
		dataBind: dataBind{
			data: make(map[string]any),
			tag:  enum.Tags.Query,
			
			DataSrc:  src,
			DataDist: dst,
			
			// Set the defaults
			MaxMemory: MaxMemoryDefault,
			FetchData: fetchQueryData,
		},
	}
	
	// Apply any custom options
	for _, opt := range opts {
		opt(&query.dataBind)
	}
	
	return query, nil
}
