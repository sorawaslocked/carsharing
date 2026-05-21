package validation

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
)

type Errors map[string]error

func (e Errors) Error() string {
	buff := bytes.NewBufferString("")

	for field, err := range e {
		buff.WriteString(fmt.Sprintf("%s: %s", field, err))
		buff.WriteString("\n")
	}

	return strings.TrimSpace(buff.String())
}

// General
var (
	ErrRequiredField = errors.New("required")
)

// String format
var (
	ErrNotAlphaNum     = errors.New("must only contain ascii letters and numbers")
	ErrNotAlphaUnicode = errors.New("must only contain letters")
	ErrNotUppercase    = errors.New("must only contain uppercase letters")
)

// Auth/User
var (
	ErrInvalidEmail          = errors.New("must be a valid email address")
	ErrInvalidPhoneNumber    = errors.New("must be a valid e.164 phone number")
	ErrInvalidDateFormat     = errors.New("must be a valid date format (YYYY-MM-DD)")
	ErrNotComplexPassword    = errors.New("must contain uppercase, lowercase, numbers, and special characters (!@#)")
	ErrInvalidRole           = errors.New("must be a valid role")
	ErrInvalidID             = errors.New("must be a valid UUID")
	ErrInvalidActivationCode = errors.New("invalid or expired activation code")
)

// Document
var (
	ErrInvalidDocumentStatus = errors.New("must be a valid document status")
	ErrInvalidImageType      = errors.New("must be a valid image type")
)
