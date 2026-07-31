package database

import (
	"fmt"
)

func ErrNotFound(data any) error {
	return fmt.Errorf("%v not found", data)
}
