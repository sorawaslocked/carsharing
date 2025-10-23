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
	ErrNotFound            = errors.New("resource not found")
	ErrRequiredField       = errors.New("required")
	ErrInvalidDateFormat   = errors.New("invalid date format")
	ErrPasswordsDoNotMatch = errors.New("passwords do not match")
	ErrInvalidToken        = errors.New("invalid token")
	ErrJwt                 = errors.New("jwt error")
	ErrBcrypt              = errors.New("bcrypt error")
	ErrSql                 = errors.New("sql error")
)
