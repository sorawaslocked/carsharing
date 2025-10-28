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
	ErrNotFound       = errors.New("resource not found")
	ErrNoUpdateFields = errors.New("no update fields set")
	ErrEmptyFilter    = errors.New("filter is empty")

	ErrRequiredField       = errors.New("required")
	ErrPasswordsDoNotMatch = errors.New("passwords do not match")
	ErrNotAlphaUnicode     = errors.New("must only contain letters")
	ErrInvalidEmail        = errors.New("must be a valid email address")
	ErrInvalidPhoneNumber  = errors.New("must be a valid 164 phone number")
	ErrInvalidDateFormat   = errors.New("must be a valid date format")
	ErrInvalidJwtToken     = errors.New("must be a valid jwt token")
	ErrNotComplexPassword  = errors.New("must contain uppercase characters, lowercase characters, numbers, and special characters(!@#)")
	ErrDuplicateEmail      = errors.New("user with this email already exists")

	ErrSqlTransaction = errors.New("sql transaction error")
	ErrSql            = errors.New("sql error")
	ErrJwt            = errors.New("jwt error")
	ErrBcrypt         = errors.New("bcrypt error")
)
