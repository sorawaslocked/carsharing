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

var (
	ErrRequiredField = errors.New("required")
)

var (
	ErrNotAlphaNum     = errors.New("must only contain ascii letters and numbers")
	ErrNotAlphaUnicode = errors.New("must only contain letters")
	ErrNotUppercase    = errors.New("must only contain uppercase letters")
)

var (
	ErrInvalidID         = errors.New("must be a valid UUID")
	ErrInvalidTripStatus = errors.New("must be a valid trip status")
)
