package handler

import (
	"fmt"
)

func errNotImplemented(msg string) error {
	return fmt.Errorf("not implemented: %s", msg)
}
