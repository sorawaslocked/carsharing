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
	ErrRequiredField = errors.New("this field is required")
)

// Car
var (
	ErrInvalidCarFuelType     = errors.New("must be a valid fuel type")
	ErrInvalidCarTransmission = errors.New("must be a valid car transmission")
	ErrInvalidCarBodyType     = errors.New("must be a valid car body type")
	ErrInvalidCarClass        = errors.New("must be a valid car class")
	ErrInvalidCarStatus       = errors.New("must be a valid car status")
)

// Location
var (
	ErrInvalidLatitudeRange  = errors.New("latitude must be between -90 and 90")
	ErrInvalidLongitudeRange = errors.New("longitude must be between -180 and 180")
	ErrInvalidRadiusRange    = fmt.Errorf("radius must be between %f and %f km", minRadiusKM, maxRadiusKM)
)
