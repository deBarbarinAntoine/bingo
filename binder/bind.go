// bind.go
package binder

import (
	"fmt"
	"mime/multipart"
	"reflect"
	"strconv"
	"time"
	
	"github.com/debarbarinantoine/bingo/internal/enum"
)

const (
	intType   = "int"
	uintType  = "uint"
	floatType = "float"
	boolType  = "bool"
	timeType  = "time"
	
	mapType             = "map"
	nestedStructType    = "nested struct"
	nestedStructPtrType = "nested struct pointer"
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

func bind(dst any, tag enum.Tag, data map[string]any) error {
	// Get the Value of the root struct
	v := reflect.ValueOf(dst).Elem()
	return bindRecursive(v, tag, data)
}

// #################################################################################
// Primitive types
// #################################################################################

// bindPrimitiveValue handles the binding of a single value to a primitive field type.
func bindPrimitiveValue(fieldValue reflect.Value, mapValue any) error {
	// If not a custom unmarshaler, handle standard types conversions from string
	stringValue := fmt.Sprintf("%v", mapValue)
	
	switch fieldValue.Kind() {
	
	case reflect.String:
		fieldValue.SetString(stringValue)
		return nil
	
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return bindInt(stringValue, fieldValue)
	
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return bindUint(stringValue, fieldValue)
	
	case reflect.Float32, reflect.Float64:
		return bindFloat(stringValue, fieldValue)
	
	case reflect.Bool:
		return bindBool(stringValue, fieldValue)
	
	case reflect.Struct:
		// If the struct is a time.Time type, parse the string value and set the field value
		if err := bindTime(stringValue, fieldValue); err != nil {
			return err
		}
		// DEBUG: Return nil anyway -> may cause silent errors
		return nil
	
	default:
	}
	return ErrBindUnsupported(fieldValue.Kind(), fieldValue.Type().Name())
}

// ################ Specific binding functions for primitive types ################

func bindInt(stringValue string, fieldValue reflect.Value) error {
	i, err := strconv.ParseInt(stringValue, 10, 64)
	if err != nil {
		return ErrBindConversion(stringValue, intType, fieldValue.Type().Name())
	}
	if fieldValue.OverflowInt(i) {
		return ErrBindOverflow(stringValue, fieldValue.Type().Name())
	}
	fieldValue.SetInt(i)
	
	return nil
}

func bindUint(stringValue string, fieldValue reflect.Value) error {
	u, err := strconv.ParseUint(stringValue, 10, 64)
	if err != nil {
		return ErrBindConversion(stringValue, uintType, fieldValue.Type().Name())
	}
	if fieldValue.OverflowUint(u) {
		return ErrBindOverflow(stringValue, fieldValue.Type().Name())
	}
	fieldValue.SetUint(u)
	
	return nil
}

func bindFloat(stringValue string, fieldValue reflect.Value) error {
	f, err := strconv.ParseFloat(stringValue, 64)
	if err != nil {
		return ErrBindConversion(stringValue, floatType, fieldValue.Type().Name())
	}
	if fieldValue.OverflowFloat(f) {
		return ErrBindOverflow(stringValue, fieldValue.Type().Name())
	}
	fieldValue.SetFloat(f)
	
	return nil
}

func bindBool(stringValue string, fieldValue reflect.Value) error {
	b, err := strconv.ParseBool(stringValue)
	if err != nil {
		return ErrBindConversion(stringValue, boolType, fieldValue.Type().Name())
	}
	fieldValue.SetBool(b)
	
	return nil
}

func bindTime(stringValue string, fieldValue reflect.Value) error {
	// Check if the struct is a time.Time type
	if fieldValue.Type().String() == "time.Time" {
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
		return ErrBindConversion(stringValue, timeType, fieldValue.Type().Name())
	}
	
	return nil
}

// #################################################################################
// Complex types
// #################################################################################

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
			if err := bindStruct(field, fieldValue, tag, mapValue); err != nil {
				return err
			}
			continue
		}
		// Handle primitive types
		if isPrimitive(fieldValue.Kind()) {
			if err := bindPrimitiveValue(fieldValue, mapValue); err != nil {
				return err
			}
			continue
		}
		// Main binding logic
		switch fieldValue.Kind() {
		case reflect.Slice:
			if err := bindSlice(field, fieldValue, tag, mapValue); err != nil {
				return err
			}
		case reflect.Map:
			if err := bindMap(field, fieldValue, mapValue); err != nil {
				return err
			}
		case reflect.Ptr:
			if err := bindPtr(field, fieldValue, tag, mapValue); err != nil {
				return err
			}
		default:
			return ErrBindUnsupported(fieldValue.Kind(), field.Name)
		}
	}
	return nil
}

// ################ Specific binding functions for complex types ################

func bindMap(field reflect.StructField, fieldValue reflect.Value, mapValue any) error {
	srcMap, isMap := mapValue.(map[string]any)
	if !isMap {
		return ErrBindNotAMap(mapType, field.Name)
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
	
	return nil
}

func bindSlice(field reflect.StructField, fieldValue reflect.Value, tag enum.Tag, mapValue any) error {
	elemType := fieldValue.Type().Elem()
	
	// Handle slices of file headers
	if elemType.Kind() == reflect.Ptr && elemType.Elem().String() == "multipart.FileHeader" {
		headers, ok := mapValue.([]*multipart.FileHeader)
		if !ok {
			return ErrBindNotSliceFileHeaders(field.Name)
		}
		newSlice := reflect.MakeSlice(fieldValue.Type(), len(headers), len(headers))
		for i, header := range headers {
			newSlice.Index(i).Set(reflect.ValueOf(header))
		}
		fieldValue.Set(newSlice)
		return nil
	}
	
	// Handle slices of structs
	if elemType.Kind() == reflect.Struct {
		srcSlice, isSlice := mapValue.([]any)
		if !isSlice {
			return nil
		}
		newSlice := reflect.MakeSlice(fieldValue.Type(), 0, len(srcSlice))
		for _, item := range srcSlice {
			itemMap, isMap := item.(map[string]any)
			if !isMap {
				return nil
			}
			newStruct := reflect.New(elemType).Elem()
			if err := bindRecursive(newStruct, tag, itemMap); err != nil {
				return err
			}
			newSlice = reflect.Append(newSlice, newStruct)
		}
		fieldValue.Set(newSlice)
		return nil
	}
	
	// Handle slices of primitives
	srcSlice, isSrcSlice := mapValue.([]string)
	if !isSrcSlice {
		return nil
	}
	newSlice := reflect.MakeSlice(fieldValue.Type(), len(srcSlice), len(srcSlice))
	for i, item := range srcSlice {
		elemValue := newSlice.Index(i)
		if err := bindPrimitiveValue(elemValue, item); err != nil {
			return err
		}
	}
	fieldValue.Set(newSlice)
	
	return nil
}

func bindStruct(field reflect.StructField, fieldValue reflect.Value, tag enum.Tag, mapValue any) error {
	
	// Special case for time.Time
	if field.Type.String() == "time.Time" {
		stringValue := fmt.Sprintf("%v", mapValue)
		if err := bindTime(stringValue, fieldValue); err != nil {
			return err
		}
		return nil
	}
	
	// Special case for os.File
	if field.Type.String() == "os.File" {
		// os.File is not supported, use *multipart.FileHeader instead
		return ErrBindOsFile(field.Name)
	}
	
	// Normal nested struct recursion
	nestedMap, isMap := mapValue.(map[string]any)
	if !isMap {
		return ErrBindNotAMap(nestedStructType, field.Name)
	}
	if err := bindRecursive(fieldValue, tag, nestedMap); err != nil {
		return err
	}
	
	return nil
}

// Pointers

func bindPtr(field reflect.StructField, fieldValue reflect.Value, tag enum.Tag, mapValue any) error {
	// Handle pointers to file headers
	if fieldValue.Type().Elem().String() == "multipart.FileHeader" {
		return bindPtrFileHeader(field, fieldValue, mapValue)
	}
	
	// Handle pointer to nested structs
	if fieldValue.Elem().Kind() == reflect.Struct {
		return bindPtrNestedStruct(field, fieldValue, tag, mapValue)
	}
	
	// Handle custom unmarshalling for other pointers
	ptrToUnmarshaler := reflect.New(fieldValue.Type().Elem())
	if unmarshaler, ok := ptrToUnmarshaler.Interface().(TextUnmarshaler); ok {
		strVal, strOk := mapValue.(string)
		if !strOk {
			return nil
		}
		if err := unmarshaler.UnmarshalText([]byte(strVal)); err == nil {
			fieldValue.Set(ptrToUnmarshaler)
			return nil
		}
	}
	
	return nil
}

func bindPtrFileHeader(field reflect.StructField, fieldValue reflect.Value, mapValue any) error {
	header, ok := mapValue.(*multipart.FileHeader)
	if !ok {
		return ErrBindNotSingleFileHeader(field.Name)
	}
	fieldValue.Set(reflect.ValueOf(header))
	return nil
}

func bindPtrNestedStruct(field reflect.StructField, fieldValue reflect.Value, tag enum.Tag, mapValue any) error {
	if fieldValue.IsNil() {
		fieldValue.Set(reflect.New(fieldValue.Type().Elem()))
	}
	nestedMap, isMap := mapValue.(map[string]any)
	if !isMap {
		return ErrBindNotAMap(nestedStructPtrType, field.Name)
	}
	if err := bindRecursive(fieldValue.Elem(), tag, nestedMap); err != nil {
		return err
	}
	return nil
}
