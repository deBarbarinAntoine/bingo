package sessions

import (
	"fmt"
	"net/http"
	"time"
)

// Get returns the value for a given key from the session data. The return value has the type interface{} so will usually need to be type asserted before you can use it. For example:
// 	fooVal, err := session.Get(r, "foo")
// 	if err != nil {
//		return fmt.Errorf("failed to get 'foo' from session: %w", err)
// 	}
// 	foo, ok := fooVal.(string)
// 	if !ok {
//		return fmt.Errorf("type assertion to string failed")
// 	}
// 	fmt.Println(foo)
//
// Also see the GetString(), GetInt(), GetBytes() and other helper methods which wrap the type conversion for common types.
func Get(r *http.Request, key string) (any, error) {
	sessionManager, err := GetSession(r)
	if err != nil {
		return false, fmt.Errorf("failed to get session: %w", err)
	}
	
	return sessionManager.Get(r.Context(), key), nil
}

// GetString returns the value for a given key from the session data as a string. If the value is not found or is not a string, an error is returned.
func GetString(r *http.Request, key string) (string, error) {
	value, err := Get(r, key)
	if err != nil {
		return "", err
	}
	
	str, ok := value.(string)
	if !ok {
		return "", fmt.Errorf("type assertion to string failed")
	}
	
	return str, nil
}

// GetInt returns the value for a given key from the session data as an integer. If the value is not found or is not an integer, an error is returned.
func GetInt(r *http.Request, key string) (int, error) {
	value, err := Get(r, key)
	if err != nil {
		return 0, err
	}
	
	i, ok := value.(int)
	if !ok {
		return 0, fmt.Errorf("type assertion to int failed")
	}
	
	return i, nil
}

// GetBytes returns the value for a given key from the session data as a byte slice. If the value is not found or is not a byte slice, an error is returned.
func GetBytes(r *http.Request, key string) ([]byte, error) {
	value, err := Get(r, key)
	if err != nil {
		return nil, err
	}
	
	b, ok := value.([]byte)
	if !ok {
		return nil, fmt.Errorf("type assertion to []byte failed")
	}
	
	return b, nil
}

// GetBool returns the value for a given key from the session data as a boolean. If the value is not found or is not a boolean, an error is returned.
func GetBool(r *http.Request, key string) (bool, error) {
	value, err := Get(r, key)
	if err != nil {
		return false, err
	}
	
	b, ok := value.(bool)
	if !ok {
		return false, fmt.Errorf("type assertion to bool failed")
	}
	
	return b, nil
}

// GetFloat64 returns the value for a given key from the session data as a float64. If the value is not found or is not a float64, an error is returned.
func GetFloat64(r *http.Request, key string) (float64, error) {
	value, err := Get(r, key)
	if err != nil {
		return 0, err
	}
	
	f, ok := value.(float64)
	if !ok {
		return 0, fmt.Errorf("type assertion to float64 failed")
	}
	
	return f, nil
}

// GetFloat32 returns the value for a given key from the session data as a float32. If the value is not found or is not a float32, an error is returned.
func GetFloat32(r *http.Request, key string) (float32, error) {
	value, err := Get(r, key)
	if err != nil {
		return 0, err
	}
	
	f, ok := value.(float32)
	if !ok {
		return 0, fmt.Errorf("type assertion to float32 failed")
	}
	
	return f, nil
}

// GetInt64 returns the value for a given key from the session data as an int64. If the value is not found or is not an int64, an error is returned.
func GetInt64(r *http.Request, key string) (int64, error) {
	value, err := Get(r, key)
	if err != nil {
		return 0, err
	}
	
	i, ok := value.(int64)
	if !ok {
		return 0, fmt.Errorf("type assertion to int64 failed")
	}
	
	return i, nil
}

// GetInt32 returns the value for a given key from the session data as an int32. If the value is not found or is not an int32, an error is returned.
func GetInt32(r *http.Request, key string) (int32, error) {
	value, err := Get(r, key)
	if err != nil {
		return 0, err
	}
	
	i, ok := value.(int32)
	if !ok {
		return 0, fmt.Errorf("type assertion to int32 failed")
	}
	
	return i, nil
}

// GetInt16 returns the value for a given key from the session data as an int16. If the value is not found or is not an int16, an error is returned.
func GetInt16(r *http.Request, key string) (int16, error) {
	value, err := Get(r, key)
	if err != nil {
		return 0, err
	}
	
	i, ok := value.(int16)
	if !ok {
		return 0, fmt.Errorf("type assertion to int16 failed")
	}
	
	return i, nil
}

// GetInt8 returns the value for a given key from the session data as an int8. If the value is not found or is not an int8, an error is returned.
func GetInt8(r *http.Request, key string) (int8, error) {
	value, err := Get(r, key)
	if err != nil {
		return 0, err
	}
	
	i, ok := value.(int8)
	if !ok {
		return 0, fmt.Errorf("type assertion to int8 failed")
	}
	
	return i, nil
}

// GetUint returns the value for a given key from the session data as a uint. If the value is not found or is not a uint, an error is returned.
func GetUint(r *http.Request, key string) (uint, error) {
	value, err := Get(r, key)
	if err != nil {
		return 0, err
	}
	
	i, ok := value.(uint)
	if !ok {
		return 0, fmt.Errorf("type assertion to uint failed")
	}
	
	return i, nil
}

// GetUint64 returns the value for a given key from the session data as a uint64. If the value is not found or is not a uint64, an error is returned.
func GetUint64(r *http.Request, key string) (uint64, error) {
	value, err := Get(r, key)
	if err != nil {
		return 0, err
	}
	
	i, ok := value.(uint64)
	if !ok {
		return 0, fmt.Errorf("type assertion to uint64 failed")
	}
	
	return i, nil
}

// GetUint32 returns the value for a given key from the session data as a uint32. If the value is not found or is not a uint32, an error is returned.
func GetUint32(r *http.Request, key string) (uint32, error) {
	value, err := Get(r, key)
	if err != nil {
		return 0, err
	}
	
	i, ok := value.(uint32)
	if !ok {
		return 0, fmt.Errorf("type assertion to uint32 failed")
	}
	
	return i, nil
}

// GetUint16 returns the value for a given key from the session data as a uint16. If the value is not found or is not a uint16, an error is returned.
func GetUint16(r *http.Request, key string) (uint16, error) {
	value, err := Get(r, key)
	if err != nil {
		return 0, err
	}
	
	i, ok := value.(uint16)
	if !ok {
		return 0, fmt.Errorf("type assertion to uint16 failed")
	}
	
	return i, nil
}

// GetUint8 returns the value for a given key from the session data as a uint8. If the value is not found or is not a uint8, an error is returned.
func GetUint8(r *http.Request, key string) (uint8, error) {
	value, err := Get(r, key)
	if err != nil {
		return 0, err
	}
	
	i, ok := value.(uint8)
	if !ok {
		return 0, fmt.Errorf("type assertion to uint8 failed")
	}
	
	return i, nil
}

// GetTime returns the value for a given key from the session data as a time.Time. If the value is not found or is not a time.Time, an error is returned.
func GetTime(r *http.Request, key string) (time.Time, error) {
	value, err := Get(r, key)
	if err != nil {
		return time.Time{}, err
	}
	
	t, ok := value.(time.Time)
	if !ok {
		return time.Time{}, fmt.Errorf("type assertion to time.Time failed")
	}
	
	return t, nil
}

// GetDuration returns the value for a given key from the session data as a time.Duration. If the value is not found or is not a time.Duration, an error is returned.
func GetDuration(r *http.Request, key string) (time.Duration, error) {
	value, err := Get(r, key)
	if err != nil {
		return 0, err
	}
	
	t, ok := value.(time.Duration)
	if !ok {
		return 0, fmt.Errorf("type assertion to time.Duration failed")
	}
	
	return t, nil
}

// Pop removes a key and corresponding value from the session data and returns the value. If the key is not found, nil is returned.
func Pop(r *http.Request, key string) (any, error) {
	sessionManager, err := GetSession(r)
	if err != nil {
		return nil, err
	}
	
	value := sessionManager.Pop(r.Context(), key)
	
	return value, nil
}

// PopString removes a key and corresponding value from the session data and returns the value as a string. If the key is not found, an empty string is returned.
func PopString(r *http.Request, key string) (string, error) {
	value, err := Pop(r, key)
	if err != nil {
		return "", err
	}
	
	s, ok := value.(string)
	if !ok {
		return "", fmt.Errorf("type assertion to string failed")
	}
	
	return s, nil
}

// PopInt removes a key and corresponding value from the session data and returns the value as an int. If the key is not found, 0 is returned.
func PopInt(r *http.Request, key string) (int64, error) {
	value, err := Pop(r, key)
	if err != nil {
		return 0, err
	}
	
	i, ok := value.(int64)
	if !ok {
		return 0, fmt.Errorf("type assertion to int64 failed")
	}
	
	return i, nil
}

// PopBytes removes a key and corresponding value from the session data and returns the value as a byte slice. If the key is not found, nil is returned.
func PopBytes(r *http.Request, key string) ([]byte, error) {
	value, err := Pop(r, key)
	if err != nil {
		return nil, err
	}
	
	b, ok := value.([]byte)
	if !ok {
		return nil, fmt.Errorf("type assertion to []byte failed")
	}
	
	return b, nil
}

// PopBool removes a key and corresponding value from the session data and returns the value as a bool. If the key is not found, false is returned.
func PopBool(r *http.Request, key string) (bool, error) {
	value, err := Pop(r, key)
	if err != nil {
		return false, err
	}
	
	b, ok := value.(bool)
	if !ok {
		return false, fmt.Errorf("type assertion to bool failed")
	}
	
	return b, nil
}

// PopInt64 removes a key and corresponding value from the session data and returns the value as an int64. If the key is not found, 0 is returned.
func PopInt64(r *http.Request, key string) (int64, error) {
	value, err := Pop(r, key)
	if err != nil {
		return 0, err
	}
	
	i, ok := value.(int64)
	if !ok {
		return 0, fmt.Errorf("type assertion to int64 failed")
	}
	
	return i, nil
}

// PopInt32 removes a key and corresponding value from the session data and returns the value as an int32. If the key is not found, 0 is returned.
func PopInt32(r *http.Request, key string) (int32, error) {
	value, err := Pop(r, key)
	if err != nil {
		return 0, err
	}
	
	i, ok := value.(int32)
	if !ok {
		return 0, fmt.Errorf("type assertion to int32 failed")
	}
	
	return i, nil
}

// PopInt16 removes a key and corresponding value from the session data and returns the value as an int16. If the key is not found, 0 is returned.
func PopInt16(r *http.Request, key string) (int16, error) {
	value, err := Pop(r, key)
	if err != nil {
		return 0, err
	}
	
	i, ok := value.(int16)
	if !ok {
		return 0, fmt.Errorf("type assertion to int16 failed")
	}
	return i, nil
}

// PopInt8 removes a key and corresponding value from the session data and returns the value as an int8. If the key is not found, 0 is returned.
func PopInt8(r *http.Request, key string) (int8, error) {
	value, err := Pop(r, key)
	if err != nil {
		return 0, err
	}
	
	i, ok := value.(int8)
	if !ok {
		return 0, fmt.Errorf("type assertion to int8 failed")
	}
	
	return i, nil
}

// PopUint removes a key and corresponding value from the session data and returns the value as a uint. If the key is not found, 0 is returned.
func PopUint(r *http.Request, key string) (uint, error) {
	value, err := Pop(r, key)
	if err != nil {
		return 0, err
	}
	
	u, ok := value.(uint)
	if !ok {
		return 0, fmt.Errorf("type assertion to uint failed")
	}
	
	return u, nil
}

// PopUint64 removes a key and corresponding value from the session data and returns the value as a uint64. If the key is not found, 0 is returned.
func PopUint64(r *http.Request, key string) (uint64, error) {
	value, err := Pop(r, key)
	if err != nil {
		return 0, err
	}
	
	u, ok := value.(uint64)
	if !ok {
		return 0, fmt.Errorf("type assertion to uint64 failed")
	}
	
	return u, nil
}

// PopUint32 removes a key and corresponding value from the session data and returns the value as a uint32. If the key is not found, 0 is returned.
func PopUint32(r *http.Request, key string) (uint32, error) {
	value, err := Pop(r, key)
	if err != nil {
		return 0, err
	}
	
	u, ok := value.(uint32)
	if !ok {
		return 0, fmt.Errorf("type assertion to uint32 failed")
	}
	
	return u, nil
}

// PopUint16 removes a key and corresponding value from the session data and returns the value as a uint16. If the key is not found, 0 is returned.
func PopUint16(r *http.Request, key string) (uint16, error) {
	value, err := Pop(r, key)
	if err != nil {
		return 0, err
	}
	
	u, ok := value.(uint16)
	if !ok {
		return 0, fmt.Errorf("type assertion to uint16 failed")
	}
	
	return u, nil
}

// PopUint8 removes a key and corresponding value from the session data and returns the value as a uint8. If the key is not found, 0 is returned.
func PopUint8(r *http.Request, key string) (uint8, error) {
	value, err := Pop(r, key)
	if err != nil {
		return 0, err
	}
	
	u, ok := value.(uint8)
	if !ok {
		return 0, fmt.Errorf("type assertion to uint8 failed")
	}
	
	return u, nil
}

// PopFloat64 removes a key and corresponding value from the session data and returns the value as a float64. If the key is not found, 0 is returned.
func PopFloat64(r *http.Request, key string) (float64, error) {
	value, err := Pop(r, key)
	if err != nil {
		return 0, err
	}
	
	f, ok := value.(float64)
	if !ok {
		return 0, fmt.Errorf("type assertion to float64 failed")
	}
	
	return f, nil
}

// PopFloat32 removes a key and corresponding value from the session data and returns the value as a float32. If the key is not found, 0 is returned.
func PopFloat32(r *http.Request, key string) (float32, error) {
	value, err := Pop(r, key)
	if err != nil {
		return 0, err
	}
	
	f, ok := value.(float32)
	if !ok {
		return 0, fmt.Errorf("type assertion to float32 failed")
	}
	
	return f, nil
}

// PopTime removes a key and corresponding value from the session data and returns the value as a time.Time. If the key is not found, zero time is returned.
func PopTime(r *http.Request, key string) (time.Time, error) {
	value, err := Pop(r, key)
	if err != nil {
		return time.Time{}, err
	}
	
	t, ok := value.(time.Time)
	if !ok {
		return time.Time{}, fmt.Errorf("type assertion to time.Time failed")
	}
	
	return t, nil
}

// PopDuration removes a key and corresponding value from the session data and returns the value as a time.Duration. If the key is not found, zero duration is returned.
func PopDuration(r *http.Request, key string) (time.Duration, error) {
	value, err := Pop(r, key)
	if err != nil {
		return 0, err
	}
	
	d, ok := value.(time.Duration)
	if !ok {
		return 0, fmt.Errorf("type assertion to time.Duration failed")
	}
	
	return d, nil
}
