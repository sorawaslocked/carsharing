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
)

// IDs and roles
var (
	ErrInvalidID   = errors.New("must be a valid UUID")
	ErrInvalidRole = errors.New("must be a valid role")
)

// Booking
var (
	ErrInvalidBookingStatus = errors.New("must be a valid booking status")
)

// Pricing rule
var (
	ErrInvalidPricingRuleType = errors.New("must be a valid pricing rule type")
)

// Car
var (
	ErrInvalidCarStatus = errors.New("must be a valid car status")
	ErrInvalidCarClass  = errors.New("must be a valid car class")
)
