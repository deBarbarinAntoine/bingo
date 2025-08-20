// binder.go
package binder

import (
	"BinGo/enum"
)

// TextUnmarshaler is a custom interface for types that can unmarshal themselves.
type TextUnmarshaler interface {
	UnmarshalText(text []byte) error
}

type Binder interface {
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

type BindOption func(bind *dataBind)

func WithCustomFetcher(fetcher func(bind *dataBind, src any) error) BindOption {
	return func(bind *dataBind) {
		bind.FetchData = fetcher
	}
}

func WithCustomMaxMemory(maxMemory int64) BindOption {
	return func(bind *dataBind) {
		bind.MaxMemory = maxMemory
	}
}

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
