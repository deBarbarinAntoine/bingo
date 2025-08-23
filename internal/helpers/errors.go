package helpers

import "fmt"

var (
	ErrJsonEncode = fmt.Errorf("failed to encode JSON response")
	ErrJsonWrite  = fmt.Errorf("failed to write JSON response")
)
