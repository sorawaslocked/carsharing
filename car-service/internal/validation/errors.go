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
	ErrInvalidID     = errors.New("must be a valid UUID")
)

// Car
var (
	ErrInvalidCarFuelType     = errors.New("must be a valid fuel type")
	ErrInvalidCarTransmission = errors.New("must be a valid car transmission")
	ErrInvalidCarBodyType     = errors.New("must be a valid car body type")
	ErrInvalidCarClass        = errors.New("must be a valid car class")
	ErrInvalidCarStatus       = errors.New("must be a valid car status")
)
