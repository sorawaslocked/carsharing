package model

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
)

type ValidationErrors map[string]error

func (ve ValidationErrors) Error() string {
	buff := bytes.NewBufferString("")
	for field, err := range ve {
		buff.WriteString(fmt.Sprintf("%s: %s", field, err))
		buff.WriteString("\n")
	}
	return strings.TrimSpace(buff.String())
}

var (
	ErrRequiredField = errors.New("required")
	ErrInvalidUUID   = errors.New("must be a valid UUID")
	ErrInvalidStatus = errors.New("must be a valid trip status")
	ErrInvalidRole   = errors.New("must be a valid role")
)
