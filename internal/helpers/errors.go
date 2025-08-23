package helpers

import "fmt"

var (
	ErrJsonResponse = fmt.Errorf("json response error")
)

func ErrJsonResponseWith(err error) error {
	return fmt.Errorf("%w: %w", ErrJsonResponse, err)
}
