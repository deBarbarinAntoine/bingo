package binder

import (
	"net/http"
	
	"github.com/debarbarinantoine/bingo/internal/enum"
)

const (
	MaxMemoryDefault = 100 * 1024 * 1024 // 100MB
)

type MultipartForm struct {
	dataBind
}

func fetchMultipartFormData(bind *dataBind, src any) error {
	r, ok := src.(*http.Request)
	if !ok || r == nil {
		return ErrInvalidSrcType("*http.Request")
	}
	
	// r.ParseMultipartForm populates both Form and PostForm.
	// We need to use MultipartForm because it contains the file data.
	err := r.ParseMultipartForm(bind.MaxMemory)
	if err != nil {
		return ErrParseMultipartForm
	}
	
	// First, unflatten the text values and store them in bind.data.
	bind.data = unflattenMap(r.MultipartForm.Value)
	
	// Now, iterate through the files and add them to the same map.
	// The keys from both maps must be unique.
	for key, headers := range r.MultipartForm.File {
		if len(headers) == 1 {
			// If it's a single file, bind it to a single *multipart.FileHeader.
			bind.data[key] = headers[0]
		} else {
			// If it's multiple files from the same field, bind to a slice.
			bind.data[key] = headers
		}
	}
	
	return nil
}

// NewMultipartForm uses functional options.
//
// Example usage:
// 		// with default fetcher with default 100MB max memory
// 		q, err := NewMultipartForm(&myStruct, &request)
//
// 		// uses custom fetcher and 1GB max memory
// 		q, err := NewMultipartForm(&myStruct, &request, WithCustomFetcher(myCustomFetcher), WithMaxMemory(1024*1024*1024))
func NewMultipartForm(dst any, src *http.Request, opts ...BindOption) (*MultipartForm, error) {
	
	err := checkDst(dst)
	if err != nil {
		return nil, ErrInvalidDst
	}
	
	if src == nil {
		return nil, ErrInvalidSrcType("*http.Request")
	}
	
	form := &MultipartForm{
		dataBind: dataBind{
			data: make(map[string]any),
			tag:  enum.Tags.MultipartForm,
			
			DataSrc:  src,
			DataDist: dst,
			
			// Set the defaults
			MaxMemory: MaxMemoryDefault,
			FetchData: fetchMultipartFormData,
		},
	}
	
	// Apply any custom options
	for _, opt := range opts {
		opt(&form.dataBind)
	}
	
	return form, nil
}
