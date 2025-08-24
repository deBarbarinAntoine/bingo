package binder

import (
	"github.com/debarbarinantoine/bingo/internal/enum"
)

// TextUnmarshaler is a custom interface for types that can unmarshal themselves.
type TextUnmarshaler interface {
	UnmarshalText(text []byte) error
}

// Binder is an interface for binding data from a source to a destination struct.
type Binder interface {
	
	// Fetch fetches data from the source and binds it to the destination struct.
	Fetch() error
}

type dataBind struct {
	data map[string]any
	tag  enum.Tag
	
	MaxMemory int64
	DataSrc   any
	DataDist  any
	FetchData func(bind *dataBind, src any) error
}

// BindOption is a function that configures a dataBind instance.
type BindOption func(bind *dataBind)

// WithCustomFetcher sets a custom fetcher function for the dataBind instance.
func WithCustomFetcher(fetcher func(bind *dataBind, src any) error) BindOption {
	return func(bind *dataBind) {
		bind.FetchData = fetcher
	}
}

// WithCustomMaxMemory sets a custom maximum memory limit for the multipart form data in the dataBind instance.
func WithCustomMaxMemory(maxMemory int64) BindOption {
	return func(bind *dataBind) {
		bind.MaxMemory = maxMemory
	}
}

// Fetch fetches data from the source and binds it to the destination struct.
func (b *dataBind) Fetch() error {
	
	// Fetch data from src
	err := b.FetchData(b, b.DataSrc)
	if err != nil {
		return err
	}
	
	// Bind data to dst struct
	err = bind(b.DataDist, b.tag, b.data)
	if err != nil {
		return err
	}
	
	return nil
}
