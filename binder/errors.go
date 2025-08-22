package binder

import (
	"errors"
	"fmt"
	"reflect"
)

var (
	ErrInvalidSrc = errors.New("invalid source")
	ErrInvalidDst = errors.New("invalid destination: must be a struct pointer")
	
	ErrParseForm          = errors.New("failed to parse form")
	ErrParseMultipartForm = errors.New("failed to parse multipart form")
	
	ErrBind                 = errors.New("bind error")
	ErrFileTypeNotSupported = errors.New("binding to os.File is not supported. Please use *multipart.FileHeader instead")
)

func ErrInvalidSrcType(srcType string) error {
	return fmt.Errorf("%w: must be a %s", ErrInvalidSrc, srcType)
}

func ErrBindConversion(value, fromType, toType string) error {
	actionVerb := "convert"
	if fromType == timeType {
		actionVerb = "parse"
	}
	return fmt.Errorf("%w: failed to %s %q to %s for field %s", ErrBind, actionVerb, value, fromType, toType)
}

func ErrBindOverflow(value, toType string) error {
	return fmt.Errorf("%w: value %q overflows field %s", ErrBind, value, toType)
}

func ErrBindUnsupported(kind reflect.Kind, field string) error {
	return fmt.Errorf("%w: unsupported type %s for field %s", ErrBind, kind, field)
}

func ErrBindNotAMap(fieldType, field string) error {
	return fmt.Errorf("%w: data for %s field %s is not a map", ErrBind, fieldType, field)
}

func ErrBindOsFile(field string) error {
	return fmt.Errorf("%w for field %s: %w", ErrBind, field, ErrFileTypeNotSupported)
}

func ErrBindNotSingleFileHeader(field string) error {
	return fmt.Errorf("%w: data for file field %s is not a single file header", ErrBind, field)
}

func ErrBindNotSliceFileHeaders(field string) error {
	return fmt.Errorf("%w: data for file slice field %s is not a slice of file headers", ErrBind, field)
}
