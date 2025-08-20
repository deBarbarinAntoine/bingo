// utils.go
package binder

import (
	"reflect"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

var (
	primitives = []reflect.Kind{
		reflect.String,
		
		reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64,
		
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
		
		reflect.Float32,
		reflect.Float64,
		
		reflect.Bool,
	}
)

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

func isPrimitive(kind reflect.Kind) bool {
	return slices.Contains(primitives, kind)
}

// unflattenMap takes a map with dot-separated keys and converts it into a nested map or slice.
func unflattenMap(data map[string][]string) map[string]any {
	unflattened := make(map[string]any)
	
	re := regexp.MustCompile(`^(.+?)\[(\d+)]$`)
	
	for key, values := range data {
		keys := strings.Split(key, ".")
		currentContainer := any(unflattened)
		
		for i, k := range keys {
			// Check for array index notation (e.g., "publications[0]")
			matches := re.FindStringSubmatch(k)
			isIndexed := len(matches) == 3
			
			if isIndexed {
				// We have an array key like "publications[0]"
				arrayKey := matches[1]
				index, _ := strconv.Atoi(matches[2])
				
				// Get or create the slice at the current level
				var currentSlice []any
				if val, ok := currentContainer.(map[string]any)[arrayKey]; ok {
					currentSlice, _ = val.([]any)
				} else {
					currentSlice = make([]any, index+1)
					currentContainer.(map[string]any)[arrayKey] = currentSlice
				}
				
				// Expand the slice if the index is out of bounds
				if index >= len(currentSlice) {
					newSlice := make([]any, index+1)
					copy(newSlice, currentSlice)
					currentSlice = newSlice
					currentContainer.(map[string]any)[arrayKey] = currentSlice
				}
				
				// Move down into the slice element
				if i < len(keys)-1 {
					if currentSlice[index] == nil {
						currentSlice[index] = make(map[string]any)
					}
					currentContainer = currentSlice[index]
				} else {
					// Last key, assign the value
					if len(values) == 1 {
						currentSlice[index] = values[0]
					} else {
						currentSlice[index] = values
					}
				}
				
			} else {
				// No array index, continue with standard map traversal
				if i == len(keys)-1 {
					// Last key, assign the value
					if len(values) == 1 {
						currentContainer.(map[string]any)[k] = values[0]
					} else {
						currentContainer.(map[string]any)[k] = values
					}
				} else {
					// Not the last key, create a new map if it doesn't exist
					if _, ok := currentContainer.(map[string]any)[k]; !ok {
						currentContainer.(map[string]any)[k] = make(map[string]any)
					}
					// Move down the map
					currentContainer = currentContainer.(map[string]any)[k]
				}
			}
		}
	}
	return unflattened
}
