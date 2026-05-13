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
	ErrRequiredField         = errors.New("required")
	ErrPasswordsDoNotMatch   = errors.New("passwords do not match")
	ErrNotAlphaNum           = errors.New("must only contain ascii letters and numbers")
	ErrNotAlphaUnicode       = errors.New("must only contain letters")
	ErrNotUppercase          = errors.New("must only contain uppercase letters")
	ErrInvalidEmail          = errors.New("must be a valid email address")
	ErrInvalidPhoneNumber    = errors.New("must be a valid e.164 phone number")
	ErrInvalidDateFormat     = errors.New("must be a valid date format (YYYY-MM-DD)")
	ErrNotComplexPassword    = errors.New("must contain uppercase, lowercase, numbers, and special characters (!@#)")
	ErrDuplicateEmail        = errors.New("user with this email already exists")
	ErrDuplicatePhone        = errors.New("user with this phone number already exists")
	ErrInvalidRole           = errors.New("must be a valid role")
	ErrInvalidActivationCode = errors.New("invalid or expired activation code")
	ErrInvalidDocumentStatus = errors.New("must be \"approved\" or \"rejected\"")
	ErrInvalidImageType      = errors.New("must be a valid image type")
)
