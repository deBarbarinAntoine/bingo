// errors.go
package bingo

import "fmt"

var (
	ErrInvalidDBPool = fmt.Errorf("invalid db pool")
	
	ErrJsonResponse = fmt.Errorf("json response error")
)

func ErrJsonResponseWith(err error) error {
	return fmt.Errorf("%w: %w", ErrJsonResponse, err)
}
