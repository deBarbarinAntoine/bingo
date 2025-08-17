package binder

import (
	"net/http"
	
	"BinGo/enum"
)

type Header struct {
	dataBind
}

func fetchHeaderData(bind *dataBind, src any) error {
	r, ok := src.(*http.Request)
	if !ok || r == nil {
		return ErrInvalidSrc
	}
	
	err := bind.getTags()
	if err != nil {
		return err
	}
	
	for k := range bind.data {
		if v := r.Header.Values(k); v != nil {
			if len(v) == 1 {
				bind.data[k] = v[0]
			} else {
				bind.data[k] = v
			}
		}
	}
	
	return nil
}

// NewHeader now uses functional options.
//
// Example usage:
// 		q, err := NewHeader(&myStruct, &request) // uses default
// 		q, err := NewHeader(&myStruct, &request, WithCustomFetcher(myCustomFetcher)) // uses custom
func NewHeader(dst any, src *http.Request, opts ...BindOption) (*Header, error) {
	err := checkDst(dst)
	if err != nil {
		return nil, ErrInvalidDst
	}
	if src == nil {
		return nil, ErrInvalidSrc
	}
	form := &Header{
		dataBind: dataBind{
			data: make(map[string]any),
			tag:  enum.Tags.Header,
			
			DataSrc:  src,
			DataDist: dst,
			
			// Set the defaults
			MaxMemory: MaxMemoryDefault,
			FetchData: fetchHeaderData,
		},
	}
	
	// Apply any custom options
	for _, opt := range opts {
		opt(&form.dataBind)
	}
	
	return form, nil
}
