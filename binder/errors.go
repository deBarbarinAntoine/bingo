package binder

import (
	"errors"
)

var (
	ErrInvalidSrc = errors.New("invalid source")
	ErrInvalidDst = errors.New("invalid destination: must be a struct pointer")
	
	ErrParseForm          = errors.New("failed to parse form")
	ErrParseMultipartForm = errors.New("failed to parse multipart form")
)
