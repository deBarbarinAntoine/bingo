package binder

import (
	"fmt"
	"mime/multipart"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
	
	"BinGo/enum"
)

const (
	MaxMemoryDefault = 100 * 1024 * 1024 // 100MB
)

var (
	timeLayouts = []string{
		time.Layout,
		time.RFC3339,
		time.RFC3339Nano,
		time.RFC1123,
		time.RFC1123Z,
		time.RFC822,
		time.RFC822Z,
		time.RFC850,
		time.DateOnly,
		time.UnixDate,
		time.TimeOnly,
		time.Kitchen,
		time.Stamp,
		time.StampMilli,
		time.StampMicro,
		time.StampNano,
		time.ANSIC,
		time.RubyDate,
		"2006/01/02",
		"02/01/2006",
	}
)

type DataBinder interface {
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

// ###################### Common Functions ######################

func (b *dataBind) getTags() error {
	return getTagsRecursive(reflect.TypeOf(b.DataDist).Elem(), b.tag, b.data)
}

func getTagsRecursive(t reflect.Type, tag enum.Tag, data map[string]any) error {
	if t.Kind() != reflect.Struct {
		return nil
	}
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Type.Kind() == reflect.Struct && field.Type.String() != "time.Time" {
			if err := getTagsRecursive(field.Type, tag, data); err != nil {
				return err
			}
		}
		if field.Type.Kind() == reflect.Ptr && field.Type.Elem().Kind() == reflect.Struct {
			if err := getTagsRecursive(field.Type.Elem(), tag, data); err != nil {
				return err
			}
		}
		if val, ok := field.Tag.Lookup(tag.String()); ok {
			data[val] = struct{}{}
		}
	}
	return nil
}

func bind(dst any, tag enum.Tag, data map[string]any) error {
	v := reflect.ValueOf(dst).Elem() // Get the Value of the root struct
	return bindRecursive(v, tag, data)
}

// bindPrimitiveValue handles the binding of a single value to a primitive field type.
func bindPrimitiveValue(fieldValue reflect.Value, mapValue any) error {
	// If not a custom unmarshaler, handle standard types conversions from string
	stringValue := fmt.Sprintf("%v", mapValue)
	
	switch fieldValue.Kind() {
	
	case reflect.String:
		fieldValue.SetString(stringValue)
	
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, err := strconv.ParseInt(stringValue, 10, 64)
		if err != nil {
			return fmt.Errorf("failed to convert %q to int for field %s", stringValue, fieldValue.Type().Name())
		}
		if fieldValue.OverflowInt(i) {
			return fmt.Errorf("value %q overflows field %s", stringValue, fieldValue.Type().Name())
		}
		fieldValue.SetInt(i)
	
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		u, err := strconv.ParseUint(stringValue, 10, 64)
		if err != nil {
			return fmt.Errorf("failed to convert %q to uint for field %s", stringValue, fieldValue.Type().Name())
		}
		if fieldValue.OverflowUint(u) {
			return fmt.Errorf("value %q overflows field %s", stringValue, fieldValue.Type().Name())
		}
		fieldValue.SetUint(u)
	
	case reflect.Float32, reflect.Float64:
		f, err := strconv.ParseFloat(stringValue, 64)
		if err != nil {
			return fmt.Errorf("failed to convert %q to float for field %s", stringValue, fieldValue.Type().Name())
		}
		if fieldValue.OverflowFloat(f) {
			return fmt.Errorf("value %q overflows field %s", stringValue, fieldValue.Type().Name())
		}
		fieldValue.SetFloat(f)
	
	case reflect.Bool:
		b, err := strconv.ParseBool(stringValue)
		if err != nil {
			return fmt.Errorf("failed to convert %q to bool for field %s", stringValue, fieldValue.Type().Name())
		}
		fieldValue.SetBool(b)
	
	case reflect.Struct:
		// Check if the struct is a time.Time type
		if fieldValue.Type().String() == "time.Time" {
			stringValue := fmt.Sprintf("%v", mapValue)
			var parsedTime time.Time
			var err error
			
			for _, layout := range timeLayouts {
				
				// If the parsing is successful, assign the parsed time to the field and return success
				if parsedTime, err = time.Parse(layout, stringValue); err == nil {
					fieldValue.Set(reflect.ValueOf(parsedTime))
					return nil
				}
			}
			
			// If we've looped through all layouts and found no match, return the error
			return fmt.Errorf("failed to parse time %q for field %s", stringValue, fieldValue.Type().Name())
		}
	
	default:
		return fmt.Errorf("unsupported primitive type %s for field %s", fieldValue.Kind(), fieldValue.Type().Name())
	}
	return nil
}

func bindRecursive(v reflect.Value, tag enum.Tag, data map[string]any) error {
	t := v.Type()
	
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldValue := v.Field(i)
		
		tagValue, ok := field.Tag.Lookup(tag.String())
		if !ok {
			continue
		}
		
		mapValue, ok := data[tagValue]
		if !ok {
			continue
		}
		
		if !fieldValue.CanSet() {
			continue
		}
		
		// Handle structs (nested structs and special cases like time.Time)
		if fieldValue.Kind() == reflect.Struct {
			// Special case for time.Time
			if field.Type.String() == "time.Time" {
				if err := bindPrimitiveValue(fieldValue, mapValue); err != nil {
					return err
				}
				continue // Move to the next field
			}
			
			// Special case for os.File
			if field.Type.String() == "os.File" {
				// os.File is not supported, use *multipart.FileHeader instead
				return fmt.Errorf("binding to os.File is not supported. Please use *multipart.FileHeader instead")
			}
			
			// Normal nested struct recursion
			nestedMap, isMap := mapValue.(map[string]any)
			if !isMap {
				return fmt.Errorf("data for nested struct %s is not a map", field.Name)
			}
			if err := bindRecursive(fieldValue, tag, nestedMap); err != nil {
				return err
			}
			continue
		}
		
		// Main binding logic
		switch fieldValue.Kind() {
		
		case reflect.String, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64, reflect.Bool:
			if err := bindPrimitiveValue(fieldValue, mapValue); err != nil {
				return err
			}
		
		case reflect.Slice:
			elemType := fieldValue.Type().Elem()
			
			// Handle slices of file headers
			if elemType.Kind() == reflect.Ptr && elemType.Elem().String() == "multipart.FileHeader" {
				headers, ok := mapValue.([]*multipart.FileHeader)
				if !ok {
					return fmt.Errorf("data for file slice field %s is not a slice of file headers", field.Name)
				}
				newSlice := reflect.MakeSlice(fieldValue.Type(), len(headers), len(headers))
				for i, header := range headers {
					newSlice.Index(i).Set(reflect.ValueOf(header))
				}
				fieldValue.Set(newSlice)
				continue
			}
			
			// Handle slices of structs
			if elemType.Kind() == reflect.Struct {
				srcSlice, isSlice := mapValue.([]any)
				if !isSlice {
					continue
				}
				newSlice := reflect.MakeSlice(fieldValue.Type(), 0, len(srcSlice))
				for _, item := range srcSlice {
					itemMap, isMap := item.(map[string]any)
					if !isMap {
						continue
					}
					newStruct := reflect.New(elemType).Elem()
					if err := bindRecursive(newStruct, tag, itemMap); err != nil {
						return err
					}
					newSlice = reflect.Append(newSlice, newStruct)
				}
				fieldValue.Set(newSlice)
				continue
			}
			
			// Handle slices of primitives
			srcSlice, isSrcSlice := mapValue.([]string)
			if !isSrcSlice {
				continue
			}
			newSlice := reflect.MakeSlice(fieldValue.Type(), len(srcSlice), len(srcSlice))
			for i, item := range srcSlice {
				elemValue := newSlice.Index(i)
				if err := bindPrimitiveValue(elemValue, item); err != nil {
					return err
				}
			}
			fieldValue.Set(newSlice)
		
		case reflect.Map:
			srcMap, isMap := mapValue.(map[string]any)
			if !isMap {
				return fmt.Errorf("data for map field %s is not a map", field.Name)
			}
			newMap := reflect.MakeMap(fieldValue.Type())
			elemType := fieldValue.Type().Elem()
			for k, v := range srcMap {
				keyVal := reflect.ValueOf(k)
				valVal := reflect.New(elemType).Elem()
				if err := bindPrimitiveValue(valVal, v); err != nil {
					return err
				}
				newMap.SetMapIndex(keyVal, valVal)
			}
			fieldValue.Set(newMap)
		
		case reflect.Ptr:
			// Handle pointers to file headers
			if fieldValue.Type().Elem().String() == "multipart.FileHeader" {
				header, ok := mapValue.(*multipart.FileHeader)
				if !ok {
					return fmt.Errorf("data for file field %s is not a single file header", field.Name)
				}
				fieldValue.Set(reflect.ValueOf(header))
				continue
			}
			
			// Handle pointer to nested structs
			if fieldValue.Elem().Kind() == reflect.Struct {
				if fieldValue.IsNil() {
					fieldValue.Set(reflect.New(fieldValue.Type().Elem()))
				}
				nestedMap, isMap := mapValue.(map[string]any)
				if !isMap {
					return fmt.Errorf("data for nested struct pointer %s is not a map", field.Name)
				}
				if err := bindRecursive(fieldValue.Elem(), tag, nestedMap); err != nil {
					return err
				}
				continue
			}
			
			// Handle custom unmarshaling for other pointers
			ptrToUnmarshaler := reflect.New(fieldValue.Type().Elem())
			if unmarshaler, ok := ptrToUnmarshaler.Interface().(TextUnmarshaler); ok {
				strVal, strOk := mapValue.(string)
				if !strOk {
					continue
				}
				if err := unmarshaler.UnmarshalText([]byte(strVal)); err == nil {
					fieldValue.Set(ptrToUnmarshaler)
					continue
				}
			}
		
		default:
			return fmt.Errorf("unsupported type %s for field %s", fieldValue.Kind(), field.Name)
		}
	}
	return nil
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

// Custom interface for types that can unmarshal themselves.
type TextUnmarshaler interface {
	UnmarshalText(text []byte) error
}
