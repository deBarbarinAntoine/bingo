package binder

import "reflect"

func checkDst(dst any) error {
	// Check if dst is a struct pointer
	v := reflect.ValueOf(dst)
	if v.Kind() != reflect.Ptr {
		return ErrInvalidDst
	}
	
	// Check if the element the pointer points to is a struct.
	// Elem() returns the value that the pointer points to.
	if v.Elem().Kind() != reflect.Struct {
		return ErrInvalidDst
	}
	
	return nil
}
